package worker

import (
	"context"
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"github.com/coreos/etcd/clientv3"
	"net"
	"time"
)

//服务注册功能 /cron/workers

/* Register模块的疑问
loopback  本地回环 , 需要看网络方面的知识
租约自动续租   : register.etcdLease.KeepAlive开启了一个goroutine进行续租
如何设置cancelFunc的 : 注册context
<-chan 类型 , 只能接受 , 不能写入的channel
getLocalIP 重看
*/

type Register struct {
	etcdClient *clientv3.Client
	etcdKvAPI  clientv3.KV    // 键值对操作API集合
	etcdLease  clientv3.Lease // 租约API集合
	localIP    string
}

// 实际注册并保活的实现
func (register *Register) keepOnline() {
	var (
		regKey             string
		leaseGrantResponse *clientv3.LeaseGrantResponse
		err                error
		keepAliveChannel   <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResponse  *clientv3.LeaseKeepAliveResponse
		cancelContext      context.Context
		cancelFunc         context.CancelFunc
	)

	//循环 , 失败的时候, 跳到RETRY后重试
	for {
		regKey = common.JOB_WORKER_DIR + register.localIP

		cancelFunc = nil

		//生成租约实体
		if leaseGrantResponse, err = register.etcdLease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		//启动 自动续租 , 获得续租response的channel
		if keepAliveChannel, err = register.etcdLease.KeepAlive(context.TODO(), leaseGrantResponse.ID); err != nil {
			goto RETRY
		}

		cancelContext, cancelFunc = context.WithCancel(context.TODO())

		//注册到etcd
		if _, err = register.etcdKvAPI.Put(cancelContext, regKey, "", clientv3.WithLease(leaseGrantResponse.ID)); err != nil {
			goto RETRY
		}

		//消费续租response
		for {
			select {
			case keepAliveResponse = <-keepAliveChannel:
				//如果续租响应失败 , 重试
				if keepAliveResponse != nil {
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		//如果注册到etcd失败 , 或者续租失败 , 则调用cancelFunc() 取消自动续租
		if cancelFunc != nil {
			cancelFunc()
		}

	}

}

var (
	G_register *Register
)

func getLocalIP() (ipv4 string, err error) {
	var (
		addrArr []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet
		isIpNet bool
	)

	//获取所有网卡
	if addrArr, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡ip (lo = loopback)
	for _, addr = range addrArr {
		//addr可能是ipv4也可能是ipv6 ,  先强制转换为*net.IPNet类型
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 只有IPV4的才能转换成功
			if ipNet.IP.To4() != nil {
				ipv4 = ipNet.IP.String()
				return
			}
		}
	}

	err = common.ERR_NO_LOCAL_IP_FOUND
	return
}

func InitRegister() (err error) {
	var (
		etcdConfig clientv3.Config
		etcdClient *clientv3.Client
		etcdKvAPI  clientv3.KV
		etcdLease  clientv3.Lease
		localIP    string
	)

	//etcd配置
	etcdConfig = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,                                     //集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, //超时时间
	}

	//新建连接
	if etcdClient, err = clientv3.New(etcdConfig); err != nil {
		return
	}

	if localIP, err = getLocalIP(); err != nil {
		return
	}

	// 得到KV和Lease的API子集
	etcdKvAPI = clientv3.NewKV(etcdClient)
	etcdLease = clientv3.NewLease(etcdClient)

	G_register = &Register{
		etcdClient: etcdClient,
		etcdKvAPI:  etcdKvAPI,
		etcdLease:  etcdLease,
		localIP:    localIP,
	}

	go G_register.keepOnline()

	return
}

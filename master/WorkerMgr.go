package master

import "github.com/coreos/etcd/clientv3"

// 职责 : 和etcd交互 , 获取 , 管理worker信息
type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_workerMgr *WorkerMgr
)

func (workerMgr *WorkerMgr) ListWorkers() (workerArr []string, err error) {

}

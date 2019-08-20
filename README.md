# Distributed Crontab

> An Distributed Crontab Project based on Go ,built with etcd,mongodb,docker,nginx





## Usage
```$xslt
docker run master.go
docker run worker.go

visit localhost:8070 in browser

add 

```

## Feature

- 具有高度伸缩性的多worker部署
- 精确到秒、分、时、日、月、周、年
- 执行日志记录, 查看

## 原理

etcd用于服务注册,发现/ master,worker之间的信息同步
mongodb 存储日志,查询日志

## 结构


## 待深入知识点
Snow Flake 算法 (Distributed Increment ID Generator)

CAP理论

Raft协议

systemctl模板

docker-composer up

systemctl/etcd/mongodb usage

architecture design (master/worker)

availability design (leader election)



## 基于
- Bootstrap
- etcd
- mongoDB
- golang
- cronexpr,etcd库


## 遇到的难点

日志写入的改进 , 并发问题

协程的调试

服务器机器时间不同步导致的锁抢占机会不均等

上下文的使用

线程和协程的关系 , Processor M Goroutine的关系

runtime.GOMAXPROCS(runtime.NumCPU())的作用

用到的库ioutil , time , strconv , string, json , http , etcd

etcd的集群和搭建, kv , lease租约概念 , put兼具新增和更新 , 更新的prekv会有值

杀任务还没写

go和协程的调试技巧 

mongodb和go的mongo库的使用总结
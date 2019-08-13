package worker

import "github.com/coreos/etcd/clientv3"

type JobLock struct {
	etcdClient *clientv3.Client
}

func (jobLock *JobLock) TryLock() (err error) {

}

func (jobLock *JobLock) UnLock() {

}

//初始化一把分布式锁并返回
func InitJobLock() (jobLock *JobLock) {

}

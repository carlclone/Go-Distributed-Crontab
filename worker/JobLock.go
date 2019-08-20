package worker

import (
	"context"
	"github.com/coreos/etcd/clientv3"
)

type JobLock struct {
	//etcdClient *clientv3.Client
	kv    clientv3.KV
	lease clientv3.Lease

	jobName    string
	cancelFunc context.CancelFunc //用于取消自动续租的协程
	leaseID    clientv3.LeaseID
	isLocked   bool
}

func (jobLock *JobLock) TryLock() (err error) {

}

func (jobLock *JobLock) UnLock() {
	if jobLock.isLocked {
		jobLock.cancelFunc()
		jobLock.lease.Revoke(context.TODO(), jobLock.leaseID)
	}
}

//初始化一把分布式锁并返回
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		kv:      kv,
		lease:   lease,
		jobName: jobName,
	}
	return
}

func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}

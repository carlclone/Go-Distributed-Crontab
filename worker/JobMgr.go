package worker

import (
	"context"
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

// Responsibility : listen for job changes and push change event to scheduler
type JobWatcher struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	// single instance
	G_jobMgr *JobWatcher
)

//获取etcd中任务列表和监听任务的变化并更新到worker中
func (jobMgr *JobWatcher) watchJobs() (err error) {
	var (
		getResp *clientv3.GetResponse
		kvpair  *mvccpb.KeyValue

		job      *common.Job
		jobName  string
		jobEvent *common.JobEvent

		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
	)

	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}

	for _, kvpair = range getResp.Kvs {
		// 反序列化json得到Job
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			// 同步给scheduler(调度协程)
			G_scheduler.PushJobEvent(jobEvent)
		}
	}

	//TODO;监听变化的协程

	return
}

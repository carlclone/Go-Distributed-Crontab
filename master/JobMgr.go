package master

import (
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"github.com/coreos/etcd/clientv3"
)

//连接etcd , 增删改查,管理任务

type JobMgr struct {
	//3个etcd的接口实例
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_jobMgr *JobMgr
)

//连接并初始化JobMgr
func InitJobMgr() (err error) {

}

//保存任务到etcd的 /cron/jobs/任务名 , 值为json , 并返回旧任务(已存在的情况)
func (jobMgr *JobMgr) SaveJob(job *common.Job) (oldJob *common.Job, err error) {

}

//删除任务
func (jobMgr *JobMgr) DeleteJob(job *common.Job) (oldJob *common.Job, err error) {

}

//任务列表,显示在web console里
func (jobMgr *JobMgr) ListJobs() (jobList []*common.Job, err error) {

}

//杀死任务 , 写入一个key到/cron/killer/任务名 处 , 如果该任务正在进行,worker会监听到并杀死子进程
func (jobMgr *JobMgr) KillJob(name string) (err error) {

}

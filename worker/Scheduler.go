package worker

import (
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"time"
)

//调度 模拟 crontab 的功能
//调度这个词太模糊了 , 调度员要负责的事情太多了 , 比如空管调度 ,  处于资源方和飞行员的中间 , 起飞,降落,航线,
type PlaneScheduler struct {
	jobEventChan     chan *common.JobEvent
	jobReadyList     map[string]*common.ReadyJob
	jobExecutingList map[string]*common.JobExecutingInfo
	jobResultChan    chan *common.JobExecutionResult
}

var (
	G_scheduler *PlaneScheduler
)

//比 linux 的 crontab 多出的功能,  动态增删,杀进程
//和资源方对接
func (scheduler *PlaneScheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		readyJob         *common.ReadyJob
		jobExecutingInfo *common.JobExecutingInfo
		jobExecuting     bool
		jobExisted       bool // ??
		err              error
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE:
		if readyJob, err = common.BuildReadyJob(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobReadyList[jobEvent.Job.Name] = readyJob
	case common.JOB_EVENT_DEL:
		if readyJob, jobExisted = scheduler.jobReadyList[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobReadyList, jobEvent.Job.Name)
		}
	case common.JOB_EVENT_KILL:
		if jobExecutingInfo, jobExecuting = scheduler.jobExecutingList[jobEvent.Job.Name]; jobExecuting {
			jobExecutingInfo.CancelFunc() // 杀死shell子进程
		}
	}
}

//模拟执行某个命令的功能
func (scheduler *PlaneScheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {

}

// 模拟了 linux 系统 crontab 的功能
func (scheduler *PlaneScheduler) TrySchedule() (scheduleAfter time.Duration) {

}

// 模拟 crontab 保存 log 功能
func (scheduler *PlaneScheduler) handleJobResult(result *common.JobExecutionResult) {

}

func (scheduler *PlaneScheduler) scheduleLoop() {

}

//动态增删杀推送
func (scheduler *PlaneScheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//模拟执行器推送结果给 crontab
func (scheduler *PlaneScheduler) PushJobResult(jobResult *common.JobExecutionResult) {
	scheduler.jobResultChan <- jobResult
}

package worker

import (
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"time"
)

//调度 模拟 crontab 的功能
type Scheduler struct {
}

var (
	G_scheduler *Scheduler
)

//比 linux 的 crontab 多出的功能,  动态增删,杀进程
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {

}

//模拟执行某个命令的功能
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) {

}

// 模拟了 linux 系统 crontab 的功能
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {

}

// 模拟 crontab 保存 log 功能
func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {

}

func (scheduler *Scheduler) scheduleLoop() {

}

//动态增删杀推送
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//模拟执行器推送结果给 crontab
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}

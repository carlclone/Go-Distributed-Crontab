package worker

import "github.com/carlclone/Go-Distributed-Crontab/common"

//从Scheduler分离出来的一些方法 还不知道怎么命名好

//动态增删杀推送
func (scheduler *PlaneScheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//模拟执行器推送结果给 crontab
func (scheduler *PlaneScheduler) PushJobResult(jobResult *common.JobExecutedResult) {
	scheduler.jobResultChan <- jobResult
}

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

// 1从执行表中删除  2记录日志
func (scheduler *PlaneScheduler) handlerJobResult(result *common.JobExecutedResult) {
	var (
		jobLog *common.JobLog
	)

	delete(scheduler.jobExecutingList, result.ExecutingInfo.Job.Name)

	if result.Err != common.ERR_LOCK_ALREADY_REQUIRED {
		jobLog = &common.JobLog{
			JobName:      result.ExecutingInfo.Job.Name,
			Command:      result.ExecutingInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecutingInfo.PlanTime.UnixNano() / 1000 / 1000,
			ScheduleTime: result.ExecutingInfo.RealTime.UnixNano() / 1000 / 1000,
			StartTime:    result.StartTime.UnixNano() / 1000 / 1000,
			EndTime:      result.EndTime.UnixNano() / 1000 / 1000,
		}

		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
		G_logsink.Append(jobLog)
	}
}

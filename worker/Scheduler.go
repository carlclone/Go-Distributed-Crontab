package worker

import (
	"fmt"
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"time"
)

//调度 模拟 crontab 的功能
//调度这个词太模糊了 , 调度员要负责的事情太多了 , 比如空管调度 ,  处于资源方和飞行员的中间 , 起飞,降落,航线,
type PlaneScheduler struct {
	jobEventChan     chan *common.JobEvent
	jobReadyList     map[string]*common.ReadyJob
	jobExecutingList map[string]*common.JobExecutingInfo
	jobResultChan    chan *common.JobExecutedResult
}

var (
	G_scheduler *PlaneScheduler
)

//模拟执行某个命令的功能
//TODO; cron表达式? 新的cron表达式? 对比
//TODO; 周边生态 crontab , etcd , mongodb 都得深入了解啊
func (scheduler *PlaneScheduler) TryStartJob(readyJob *common.ReadyJob) {
	var (
		jobExecutingInfo *common.JobExecutingInfo
		jobExecuting     bool
	)

	// 这里是怎么防止一个任务并发执行的 , 这里只防止了单个worker中的并发问题 , 多worker并发还得靠加锁
	// 执行的任务可能运行很久, 1分钟会调度60次，但是只能执行1次, 防止并发！
	if jobExecutingInfo, jobExecuting = scheduler.jobExecutingList[readyJob.Job.Name]; jobExecuting {
		// 尚未执行完,跳过这次定时计划
		return
	}

	scheduler.jobExecutingList[readyJob.Job.Name] = jobExecutingInfo

	fmt.Println("尝试执行任务:", jobExecutingInfo.Job.Name)
	G_executor.ExecuteJob(jobExecutingInfo)
}

// 模拟了 linux 系统 crontab 的功能
func (scheduler *PlaneScheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		readyJob   *common.ReadyJob
		now        time.Time
		recentTime *time.Time
	)

	//优化 (可忽略)
	if len(scheduler.jobReadyList) == 0 {
		scheduleAfter = 1 * time.Second
		return
	}

	now = time.Now()

	for _, readyJob = range scheduler.jobReadyList {
		//核心start
		if readyJob.NextTime.Before(now) || readyJob.NextTime.Equal(now) {
			scheduler.TryStartJob(readyJob)
			readyJob.NextTime = readyJob.Expr.Next(now)
		}
		//核心end

		//TODO;这没搞懂 , 懂了! 下次遍历肯定是先遍历到这个 , 遍历的耗时忽略不计 , 那下次遍历时间就是这个任务的时间 (为了减少遍历次数的优化 , 会不会出现执行时间不准确的问题啊)
		// 统计最近一个要过期的任务时间
		if recentTime == nil || readyJob.NextTime.Before(*recentTime) {
			recentTime = &readyJob.NextTime
		}
	}

	//优化 可忽略
	scheduleAfter = (*recentTime).Sub(now)
	return

}

func (scheduler *PlaneScheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecutedResult
	)

	scheduleAfter = scheduler.TrySchedule()

	scheduleTimer = time.NewTimer(scheduleAfter)

	for {
		select {
		case jobEvent = <-scheduler.jobEventChan:
			scheduler.handleJobEvent(jobEvent)
		case jobResult = <-scheduler.jobResultChan:
			scheduler.handlerJobResult(jobResult)
		case <-scheduleTimer.C:
			scheduleAfter = scheduler.TrySchedule()
			scheduleTimer.Reset(scheduleAfter)
		}
	}
}

// 初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &PlaneScheduler{
		jobEventChan:     make(chan *common.JobEvent, 1000),
		jobReadyList:     make(map[string]*common.ReadyJob),
		jobExecutingList: make(map[string]*common.JobExecutingInfo),
		jobResultChan:    make(chan *common.JobExecutedResult, 1000),
	}
	// 启动调度协程
	go G_scheduler.scheduleLoop()
	return
}

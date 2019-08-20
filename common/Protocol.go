package common

//master和worker之间的全局Job抽象协议
type Job struct {
	Name     string `json:"name"`
	Command  string `json:"command"`
	CronExpr string `json:"cronExpr"`
}

//事件抽象 , worker的JobWatcher和Scheduler之间的协商
type JobEvent struct {
	EventType int
	Job       *Job
}

const (
	JOB_EVENT_SAVE = 1 // 保存任务事件
	JOB_EVENT_DEL  = 2 // 删除任务事件
	JOB_EVENT_KILL = 3 // 强杀任务事件
)

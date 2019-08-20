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

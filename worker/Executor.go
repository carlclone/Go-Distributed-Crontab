package worker

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {

}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}

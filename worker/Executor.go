package worker

import (
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"math/rand"
	"os/exec"
	"time"
)

type Executor struct {
}

var (
	G_executor *Executor
)

func (executor *Executor) ExecuteJob(info *common.JobExecutingInfo) {
	go func() {
		var (
			cmd     *exec.Cmd
			err     error
			output  []byte
			result  *common.JobExecutedResult
			jobLock *JobLock
		)

		result = &common.JobExecutedResult{
			ExecutingInfo: info,
			Output:        make([]byte, 0),
		}

		//TODO;如何调试协程?
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)

		//随机睡眠0-1秒 , 保证每台机器上的协程都有随机机会抢到锁 , 而不会因为时间不同步问题被一台机器一直抢到
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		//记录尝试上锁时间
		result.TryLockTime = time.Now()

		err = jobLock.TryLock()
		defer jobLock.UnLock()

		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			result.StartTime = time.Now()

			// -c 是干嘛?
			//-c        If the -c option is present, then commands are read from the first non-option argument command_string.
			// 			If there are arguments after the command_string, they are
			//           assigned to the positional parameters, starting with $0.
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)

			output, err = cmd.CombinedOutput()

			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}

		//推送结果给调度者
		G_scheduler.PushJobResult(result)

	}()
}

func InitExecutor() (err error) {
	G_executor = &Executor{}
	return
}

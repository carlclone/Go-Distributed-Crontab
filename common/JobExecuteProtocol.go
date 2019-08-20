package common

import (
	"context"
	"github.com/gorhill/cronexpr"
	"time"
)

// 任务执行状态
type JobExecutingInfo struct {
	Job        *Job               // 任务信息
	PlanTime   time.Time          // 理论上的调度时间
	RealTime   time.Time          // 实际的调度时间
	CancelCtx  context.Context    // 任务command的context
	CancelFunc context.CancelFunc //  用于取消command执行的cancel函数
}

// 任务执行结果
type JobExecutionResult struct {
	ExecuteInfo *JobExecutingInfo // 执行状态
	Output      []byte            // 脚本输出
	Err         error             // 脚本错误原因
	StartTime   time.Time         // 启动时间
	EndTime     time.Time         // 结束时间
	TryLockTime time.Time
}

type ReadyJob struct {
	Job      *Job                 // 要调度的任务信息
	Expr     *cronexpr.Expression // 解析好的cronexpr表达式
	NextTime time.Time            // 下次调度时间
}

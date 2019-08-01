package main

import (
	"github.com/gorhill/cronexpr"
	"time"
)

type CronJob struct {
	expr     *cronexpr.Expression
	nextTime time.Time //expr.Next(now)
}

func main() {
	/*


		var (
			cronJob *CronJob
			expr *cronexpr.Expression
			now time.Time
			scheduleTable map[string]*CronJob
		)

		scheduleTable = make(map[string]*CronJob)

		now = time.Now()

		expr = cronexpr.MustParse()

	*/
}

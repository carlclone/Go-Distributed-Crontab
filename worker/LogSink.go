package worker

import (
	"context"
	"github.com/carlclone/Go-Distributed-Crontab/common"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

/*
疑问

Append里的写入channel的select写法, 队列满了则走default吗
*/

type LogSink struct {
	mongoClient       *mongo.Client
	logCollection     *mongo.Collection
	logChannel        chan *common.JobLog
	autoCommitChannel chan *common.LogBatch
}

var (
	G_logsink *LogSink
)

//批量写入日志到mongodb
func (logSink *LogSink) saveLogs(batch *common.LogBatch) {

}

func (logSink *LogSink) writeLoop() {
	var (
		log             *common.JobLog
		logBatch        *common.LogBatch
		autoCommitTimer *time.Timer      // ?
		timeoutBatch    *common.LogBatch // ?
	)

	for {
		select {
		case log = <-logSink.logChannel:
			// TODO; process log
			if logBatch == nil {
				logBatch = &common.LogBatch{}

				//设定超时自动提交
				autoCommitTimer = time.AfterFunc(
					time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() { //为什么要返回一个func()变量
							logSink.autoCommitChannel <- batch
						}
					}(logBatch), //这里是指针吧 , 如果被清空了? 那指针就不相等了
				)
			}
			//把新日志追加到batch中
			logBatch.Logs = append(logBatch.Logs, log)

			//如果满了则发送到MongoDB
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				//发日志

				//清空batch桶 , 可以构造一个桶的数据结构 , 发送的时候自动清空
				//取消自动提交定时器
			}

		case timeoutBatch = <-logSink.autoCommitChannel:
			if timeoutBatch != logBatch {
				//pass , 否则就存了两份log了
			}
			logSink.saveLogs(timeoutBatch)
			logBatch = nil
		}
	}

}

func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChannel <- jobLog:
	default:
		//队列满了,什么都不做 (丢弃)
	}
}

func InitLogSink() (err error) {
	var (
		mongoClient *mongo.Client
	)

	if mongoClient, err = mongo.Connect(
		context.TODO(),
		G_config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond)); err != nil {
		return
	}

	G_logsink = &LogSink{
		mongoClient:       mongoClient,
		logCollection:     mongoClient.Database("cron").Collection("log"),
		logChannel:        make(chan *common.JobLog, 1000),
		autoCommitChannel: make(chan *common.LogBatch, 1000),
	}

	//启动日志处理协程
	go G_logsink.writeLoop()

	return
}

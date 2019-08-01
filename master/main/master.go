package main

import (
	"flag"
	"fmt"
	"github.com/carlclone/Go-Distributed-Crontab/master"
	"runtime"
)

var (
	confFile string //配置文件路径
)

func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "配置文件路径")
	flag.Parse()
}

func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	initArgs()
	initEnv()

	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}

ERR:
	fmt.Println(err)
}

package main

import (
	"flag"
	"fmt"
	"github.com/carlclone/Go-Distributed-Crontab/master"
	"runtime"
	"time"
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

	initArgs() //可自定义配置文件 路径
	initEnv()  // 初始化环境

	//读取配置文件到全局变量
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}

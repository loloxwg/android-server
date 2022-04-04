package main

import (
	"androidServer/api/apibase"
	"androidServer/app"
	"androidServer/app/log"
	"androidServer/app/utils"
	"androidServer/handler"
	"flag"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	config  = flag.String("c", "", "configuration file, json format")
	version = flag.Bool("v", false, "show version")
)

//版本信息
var (
	BuildVersion string
	BuildTime    string
	BuildName    string
)

func runHttpServer(listenIP string, httpPort uint16) {
	if listenIP == "" {
		fmt.Printf("No http listen ip can be found")
		panic(listenIP)
	}
	if httpPort == 0 {
		fmt.Printf("No http listen port can be found")
		panic(httpPort)
	}

	httpAddr := fmt.Sprintf("%s:%d", listenIP, httpPort)
	err := handler.HttpRun(httpAddr)
	fmt.Printf("The tcp server is running faild: %v\n", err)
	panic(err)
}

func redirectCrashDump() (err error) {
	redirectpath := utils.RedirectPath(app.Conf.LogConfig.Dir, app.Conf.CrashLogConfig.LogPrefix, app.Conf.CrashLogConfig.LogSuffix)
	err = utils.RedirectOutput(redirectpath)
	return
}

func initLog() (err error) {

	// 初始化 log 模块
	err = log.InitGlobalLogger(&app.Conf.LogConfig, zap.AddCallerSkip(1))
	if err != nil {
		fmt.Printf("Init log module failed: %v\n", err)
		return err
	}
	fmt.Printf("Init log module successfully\n")
	return
}

//初始化mysql
func initMysql() error {
	err := app.InitMysql()
	if err != nil {
		fmt.Printf("Init mysql module failed: %v\n", err)
		return err
	}
	return nil
}

func initUserMap() error {
	err := apibase.InitCache()
	if err != nil {
		fmt.Printf("Init usermap module failed: %v\n", err)
		return err
	}
	return nil
}

func showVersion() {
	fmt.Printf("Build Name:\t%s\n", BuildName)
	fmt.Printf("Build Verison:\t%s\n", BuildVersion)
	fmt.Printf("Build Time:\t%s\n", BuildTime)
}

func main() {

	flag.Parse()
	if len(os.Args) == 2 && *version {
		showVersion()
		return
	}

	if *config == "" {
		_, err := fmt.Fprintln(os.Stderr, "missing config file!")
		if err != nil {
			return
		}
		flag.PrintDefaults()
		return
	}

	//加载配置文件
	err := app.LoadConf(*config)
	if err != nil {
		fmt.Printf("load senseptrel-backend config file %v failed: %v\n", config, err)
		panic(err)
	}

	//初始化日志
	err = initLog()
	if err != nil {
		fmt.Printf("Init log module failed: %v\n", err)
		panic(err)
	}

	//初始化mysql
	err = initMysql()
	if err != nil {
		fmt.Printf("Init mysql module failed: %v\n", err)
		panic(err)
	}

	err = initUserMap()
	if err != nil {
		fmt.Printf("Init usermap module failed: %v\n", err)
		panic(err)
	}

	//重定向crash输出
	//err = redirectCrashDump()
	//if err != nil {
	//	fmt.Fprintln(os.Stderr, err)
	//	return
	//}

	// 启动 http 模块
	go runHttpServer(app.Conf.ListenIP, app.Conf.HttpPort)

	//开启db-stat定时器模块
	//go timer.BucketTimer()
	//go timer.CheckClearTask() //task任务重启生效
	//go timer.TaskTimer()
	//go timer.GroupTrustTimer() //托管权限任务

	// 在程序退出时记录一下
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	for {
		<-c
		log.Error("service will be down")
		time.Sleep(1 * time.Second)
		os.Exit(1)
	}
}

package main

import (
	"EntryTask/config"
	"EntryTask/database"
	"EntryTask/internal/service"
	"EntryTask/rpc/server"
	rpcService "EntryTask/rpc/service"
	"github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

func main() {
	// 初始化日志
	logrus.SetLevel(logrus.TraceLevel)
	// 连接Mysql
	database.MysqlInit()
	// 连接Redis
	database.RedisInit()
	// 初始化RPC Server
	rpcServer := server.MakeServer()
	userService := &service.UserService{}
	rpcService := rpcService.MakeService(userService)
	rpcServer.Register(rpcService)
	go rpcServer.Accept(config.RpcAddr)
	logrus.Infoln("tcpserver start...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		logrus.Info("tcpServer shutdown")
		return
	}
}

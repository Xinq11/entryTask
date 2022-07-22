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
	// Mysql
	database.MysqlInit()
	// Redis
	database.RedisInit()
	// RPC Server
	rpcServer := server.MakeServer()
	userService := &service.UserService{}
	rpcService := rpcService.MakeService(userService)
	rpcServer.Register(rpcService)
	go rpcServer.Accept(config.RpcAddr)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for {
		<-c
		logrus.Info("tcpServer shutdown")
		return
	}
}

package main

import (
	"EntryTask/config"
	"EntryTask/internal/controller"
	"EntryTask/logger"
	"EntryTask/rpc/client"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

// 路由
func route() {
	http.HandleFunc("/api/entrytask/user/signup", controller.SignUpHandler)
	http.HandleFunc("/api/entrytask/user/signin", controller.SignInHandler)
	http.HandleFunc("/api/entrytask/user/signout", controller.SignOutHandler)
	http.HandleFunc("/api/entrytask/user/get_user_info", controller.GetUserInfoHandler)
	http.HandleFunc("/api/entrytask/user/update_profile_pic", controller.UpdateProfilePicHandler)
	http.HandleFunc("/api/entrytask/user/update_nickname", controller.UpdateNicknameHandler)
}

func main() {
	// 初始化日志
	logger.Init()
	// 初始化RPC Client
	err := client.MakeClient(config.RpcAddr)
	if err != nil {
		logrus.Panic("HttpServer MakeClient error: ", err.Error())
	}
	// 启动HTTP Server
	route()
	logrus.Infoln("httpserver start...")
	err = http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		logrus.Panic("HttpServer ListenAndServe error: ", err.Error())
	}
}

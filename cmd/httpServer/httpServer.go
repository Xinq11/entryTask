package main

import (
	"EntryTask/config"
	"EntryTask/internal/controller"
	"EntryTask/rpc/client"
	"github.com/sirupsen/logrus"
	"net/http"
	_ "net/http/pprof"
)

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
	logrus.SetLevel(logrus.TraceLevel)
	// RPC Client
	err := client.MakeClient(config.RpcAddr)
	if err != nil {
		logrus.Panic("RPC MakeClient error: ", err.Error())
	}
	// HTTP Server
	route()
	err = http.ListenAndServe(":9090", nil) // 设置监听的端口
	if err != nil {
		logrus.Panic("HttpServer ListenAndServe error: ", err.Error())
	}
}

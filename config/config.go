package config

import "time"

const (
	// Mysql配置参数
	MysqlUsername     = "root:123456@tcp(127.0.0.1:3306)/entrytask?charset=utf8"
	MysqlMaxOpenConns = 500
	MysqlMaxIdleConns = 500

	// redis配置参数
	RedisAddr       = "127.0.0.1:6379"
	RedisPassword   = "123456"
	RedisPoolNum    = 200
	DialTimeout     = time.Duration(30) * time.Second
	DialReadTimeout = time.Duration(30) * time.Second
	// 图片保存配置参数
	FilePath = "img/"
	LogPath  = "logger/"

	// RPC配置
	RpcAddr = "127.0.0.1:20000"
	ConnNum = 2000

	// salt生成字符串
	Letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 用户信息默认值
	DefaultNickname    = "default"
	DefaultProfilePath = "default-2022-07-26.jpg"
)

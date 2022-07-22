package config

const (
	// Mysql配置参数
	MysqlUsername     = "root:123456@tcp(127.0.0.1:3306)/entrytask?charset=utf8"
	MysqlMaxOpenConns = 500
	MysqlMaxIdleConns = 500

	// redis配置参数
	RedisAddr     = "127.0.0.1:6379"
	RedisPassword = "123456"
	RedisPoolNum  = 30

	// 图片保存配置参数

	// RPC
	RpcAddr = "127.0.0.1:20000"
)

package logger

import (
	"EntryTask/config"
	"log"
	"os"
)

func Init() {
	file := config.LogPath + "log.txt"
	logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		panic(err.Error())
	}
	log.SetOutput(logFile) // 将文件设置为log输出的文件
	log.SetFlags(log.Lshortfile | log.LUTC)
	return
}

func Info(msg string) {
	log.SetPrefix("[INFO]")     // 前置字符串加上特定标记
	log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式
	log.Println(msg)
}

func Error(msg string) {
	log.SetPrefix("[ERROR]")    // 前置字符串加上特定标记
	log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式
	log.Println(msg)
}

func Warn(msg string) {
	log.SetPrefix("[WARN]")     // 前置字符串加上特定标记
	log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式
	log.Println(msg)
}

func Panic(msg string) {
	log.SetPrefix("[PANIC]")    // 前置字符串加上特定标记
	log.SetFlags(log.LstdFlags) // 设置成日期+时间 格式
	log.Panic(msg)
}

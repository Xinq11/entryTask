package database

import (
	"EntryTask/config"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

var MySqlDB *sql.DB

func MysqlInit() {
	db, err := sql.Open("mysql", config.MysqlUsername)
	if err != nil {
		logrus.Panic("Mysql open error: ", err.Error())
	}
	db.SetMaxOpenConns(config.MysqlMaxOpenConns)
	db.SetMaxIdleConns(config.MysqlMaxIdleConns)
	if err = db.Ping(); err != nil {
		logrus.Panic("Mysql connect error: ", err.Error())
	}
	MySqlDB = db
	logrus.Infoln("Mysql connect success")
}

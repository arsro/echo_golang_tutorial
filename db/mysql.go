package db

import (
	"github.com/Sirupsen/logrus"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	
	"../config"
)


func Init() *dbr.Session {

	session := GetSession()

	return session
}

func GetSession() *dbr.Session {
	var db_setting string = config.USER + ":" + config.PASSWD + "@tcp(" + config.HOST + ":" + config.PORT + ")/" + config.DB_NAME

	conn, err := dbr.Open("mysql", db_setting, nil)
	if err != nil {
		logrus.Error(err)
	}else{
		sess := conn.NewSession(nil)
		return sess
	}
	return nil
}

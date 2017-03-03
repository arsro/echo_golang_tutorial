package controllers

import (
	"github.com/labstack/echo"

	"net/http"

	"../db"
)


func ConnectDB( c echo.Context ) error{
	// mysqlへの接続設定
	// dbr.Open("mysql", "[mysql_username]:[mysql_passwd]@tcp( [接続先のmysqlコンテナのホストネーム] :3306)/[接続DB名]", nil)
	sess := db.Init()
	
	var user []db.RC_User
	sess.Select("*").From("user_info").Where("id = ?", 1).Load(&user)
	
	return c.JSON(http.StatusCreated, user)
}

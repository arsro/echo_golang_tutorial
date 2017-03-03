package controllers

import (
	"github.com/labstack/echo"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gocraft/dbr"
	
	"strconv"
	"net/http"
	
	"../models"
)

func BindJson(c echo.Context) error {
	// 構造体の初期化1
	// user := &models.User{ ID: models.Seq }
	
	// 構造体の初期化2
	user := new(models.User)
	// user.ID = models.Seq
	
	
	// リクエストをUser構造体にbindする
	// cf. https://echo.labstack.com/guide/request
	if err := c.Bind(user); err != nil {
		return err
	}
	
	models.Users[user.ID] = user
	return c.JSON(http.StatusCreated, models.Users)
}

/**
 * リクエストURLのpostパラメータを取得してjsonで表示
 */
func GetUser(c echo.Context) error {
	
	// string -> int
	/**
	 *  Package strconv
	 *  see https://golang.org/pkg/strconv/
	 */
	id,_ := strconv.Atoi(c.QueryParam("id"))
	var name string = c.QueryParam("name")
	age,_ := strconv.Atoi(c.QueryParam("age"))

	var user *models.User = newUser(id, name, age)
		
	models.Users[user.ID] = user
	return c.JSON(http.StatusCreated, models.Users)
}

func PostUser(c echo.Context) error {
	conn, err := dbr.Open("mysql", "tomorrow:root@tcp(mysql:3306)/test_db", nil)
	if err != nil{
		return err
	}
	sess := conn.NewSession(nil)
	
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	c.JSON(http.StatusCreated, user)
	
	result, err := sess.InsertInto("users").
					Columns("id", "name", "age").
					Record(user).Exec()

	if err != nil {
		return err
	} else {
		count, _ := result.RowsAffected()
		return c.JSON(http.StatusCreated, count)
	}
}

func DeleteUser(c echo.Context) error {
	conn, err := dbr.Open("mysql", "tomorrow:root@tcp(mysql:3306)/test_db", nil)
	if err != nil{
		return err
	}
	sess := conn.NewSession(nil)
	
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return err
	}
	c.JSON(http.StatusCreated, user)
	
	result, err := sess.DeleteFrom("users").
					Where("id = ?", user.ID).
					Exec()
					
	if err != nil {
		return c.JSON(http.StatusCreated, err)
	}
	return c.JSON(http.StatusCreated, result)
}

func PutUser(c echo.Context) error {
	
	conn, err := dbr.Open("mysql", "tomorrow:root@tcp(mysql:3306)/test_db", nil)
	if err != nil{
		return err
	}
	sess := conn.NewSession(nil)
	
	user := new(models.User)
		c.JSON(http.StatusCreated, user)
	if err := c.Bind(user); err != nil {
		return err
	}
	c.JSON(http.StatusCreated, user)
	result, err := sess.Update("users").
					Set("name", user.Name).
					Set("age", 33).
					Where("id = ?", user.ID).
					Exec()

	if err != nil {
		c.JSON(http.StatusCreated, "error")
		c.JSON(http.StatusCreated, err)
	} else {
		// 変更されたレコード数を取得
		count, _ := result.RowsAffected()
		c.JSON(http.StatusCreated, count)// => 1
	}
	return c.JSON(http.StatusCreated, user)
}

//構造体は下記のような初期化関数を作るのが一般的らしい...
func newUser(id int, name string, age int) *models.User{
	user := new(models.User)
	user.ID = id
	user.Name = name
	user.Age = age
	return user
}

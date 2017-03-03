package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	
	"./controllers"
)

func main() {
	e := echo.New()
	
	/**
	 * Middleware
	 */
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	//CORSを許可する設定???
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	
	/**
	 * Routing
	 */
	e.GET("/", controllers.ConnectDB)
	
	e.GET("/users", controllers.GetUser)
	// e.POST("/users", controllers.BindJson)
	e.POST("/users", controllers.PostUser)
	e.PUT("/users", controllers.PutUser)
	e.DELETE("/users", controllers.DeleteUser)
	
	// Start server
	e.Logger.Fatal(e.Start(":3000"))
}

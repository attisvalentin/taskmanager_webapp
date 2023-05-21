package main

import (
	"github.com/attisvalentin/tasklist/handler"
	"github.com/attisvalentin/tasklist/middleware"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func init() {
	handler.Firstinit()
}

func main() {
	app := gin.Default()
	app.Use(sessions.Sessions("sessions", middleware.SessionStore()))

	app.Static("/static", "./static")
	app.LoadHTMLGlob("templates/*")

	app.GET("/", handler.LoadLogin)
	app.GET("/logout", handler.Logout)
	app.GET("/registering", handler.LoadRegisterForm)
	app.POST("/register", handler.Register)
	app.GET("/index", handler.ReturnTaskList)
	app.POST("/login", handler.Login)
	app.POST("/task", handler.AddTask)
	app.POST("/changestatus", handler.ChangeStatus)
	app.POST("/adduser", handler.AddUser)

	app.Run()
}

package main

import (
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/routers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func init() {
	database.EnvLoad()
	database.DbConnect()
}

func main() {
	router := gin.Default()

	store := cookie.NewStore([]byte("secret"))
	router.Use(sessions.Sessions("mysession", store))

	router.LoadHTMLGlob("template/*")

	user := router.Group("/user")
	routers.UserRouters(user)

	admin := router.Group("/admin")
	routers.AdminRouters(admin)

	router.Run(":8080")
}

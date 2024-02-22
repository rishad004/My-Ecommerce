package main

import (
	"project/database"
	"project/routers"

	"github.com/gin-gonic/gin"
)

func init() {
	database.EnvLoad()
	database.DbConnect()
}

func main() {
	router := gin.Default()

	user := router.Group("/user")
	routers.UserRouters(user)

	admin := router.Group("/admin")
	routers.AdminRouters(admin)

	router.Run(":8080")
}

package main

import (
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/routers"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/rishad004/My-Ecommerce/docs"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/files"
)

//	@title			B Y E C O M  LTD
//	@version		1.0
//	@description	B Y E C O M  is your ecom. shop solution, where you can sell anything online with ease and at an affordable price. We provide an eCommerce platform for businesses to sell their products online and connect with customers worldwide.
// @host      https://byecom.shop
// @BasePath  /

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

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run(":8080")
}

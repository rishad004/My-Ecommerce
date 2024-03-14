package routers

import (
	controllers "project/controllers/user"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouters(r *gin.RouterGroup) {

	//---------------Login&SignUp
	r.POST("/signup", controllers.PostSignupU)
	r.POST("/signup/otp", controllers.PostOtpU)
	r.POST("/login", controllers.PostLoginU)
	r.DELETE("/logout", middleware.Auth, controllers.LogoutU)

	//---------------Home & Product
	r.GET("/home", controllers.UserHome)
	r.GET("/product/:Id", controllers.UserShowP)

	//---------------Profile
	r.GET("/profile", middleware.Auth, controllers.UserProfile)

	//---------------Address
	r.POST("/address", middleware.Auth, controllers.AddAddress)

	//---------------Rating
	r.POST("/rating/:Id", middleware.Auth, controllers.AddRating)

	//---------------Cart
	r.POST("/cart/:Id/:Color", middleware.Auth, controllers.AddCart)
	r.GET("/cart", middleware.Auth, controllers.ShowCart)
}

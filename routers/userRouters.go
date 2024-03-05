package routers

import (
	controllers "project/controllers/user"

	"github.com/gin-gonic/gin"
)

func UserRouters(r *gin.RouterGroup) {

	//---------------Login&SignUp
	r.POST("/signup", controllers.PostSignupU)
	r.POST("/signup/otp", controllers.PostOtpU)
	r.POST("/login", controllers.PostLoginU)

	//---------------Product
	r.GET("/home", controllers.UserHome)
	r.GET("/product/:Id", controllers.UserShowP)

	//---------------Home
	r.GET("/profile", controllers.UserProfile)

	//---------------Address
	r.POST("/address", controllers.AddAddress)

	//---------------Rating
	r.POST("/rating/:Id", controllers.AddRating)

	//---------------Cart
	r.POST("/cart/:Id/:Color", controllers.AddCart)
	r.GET("/cart", controllers.ShowCart)
}

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
	r.GET("/home/sort", controllers.SortProduct)
	r.GET("/home/search", controllers.UserSearchP)

	//---------------Profile
	r.GET("/profile", middleware.Auth, controllers.UserProfile)
	r.PATCH("/password", middleware.Auth, controllers.UpdatePass)
	r.PATCH("/profile", middleware.Auth, controllers.EditProfile)

	//---------------Address
	r.POST("/address", middleware.Auth, controllers.AddAddress)
	r.PUT("/address/:Id", middleware.Auth, controllers.EditAddress)
	r.DELETE("/address/:Id", middleware.Auth, controllers.DeleteAddress)

	//---------------Rating
	r.POST("/rating/:Id", middleware.Auth, controllers.AddRating)
	r.PUT("/rating/:Id", middleware.Auth, controllers.EditRating)

	//---------------Cart
	r.POST("/cart/:Id/:Color", middleware.Auth, controllers.AddCart)
	r.GET("/cart", middleware.Auth, controllers.ShowCart)
	r.PATCH("/cart/:Id", middleware.Auth, controllers.LessCart)
	r.DELETE("/cart/:Id", middleware.Auth, controllers.DeleteCart)

	//---------------Order
	r.POST("/cart/checkout", middleware.Auth, controllers.CheckoutCart)
	r.PATCH("/order", middleware.Auth, controllers.CancelOrder)
	r.GET("/order", middleware.Auth, controllers.ShowOrder)

	//---------------Payment
	r.GET("/payment/:payment", middleware.Auth, controllers.RazorPay)
	r.POST("/payment", middleware.Auth, controllers.RazorPayVerify)
}

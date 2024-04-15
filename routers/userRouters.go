package routers

import (
	controllers "github.com/rishad004/My-Ecommerce/controllers/user"
	"github.com/rishad004/My-Ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouters(r *gin.RouterGroup) {

	//---------------Login&SignUp
	r.POST("/signup", controllers.PostSignupU)
	r.POST("/signup/otp", controllers.PostOtpU)
	r.POST("/login", controllers.PostLoginU)
	r.DELETE("/logout", middleware.Auth, controllers.LogoutU)

	//---------------Google
	r.GET("/google/login", controllers.GoogleLogin)
	r.GET("/google/callback", controllers.GoogleCallback)

	//---------------Home & Product
	r.GET("/home", controllers.UserHome)
	r.GET("/product", controllers.UserShowP)
	r.GET("/home/sort", controllers.SortProduct)
	r.GET("/home/search", controllers.UserSearchP)
	r.GET("/home/filter", controllers.FilterProduct)

	//---------------Profile
	r.GET("/profile", middleware.Auth, controllers.UserProfile)
	r.PATCH("/password", middleware.Auth, controllers.UpdatePass)
	r.PATCH("/profile", middleware.Auth, controllers.EditProfile)

	//---------------Address
	r.POST("/address", middleware.Auth, controllers.AddAddress)
	r.PUT("/address", middleware.Auth, controllers.EditAddress)
	r.DELETE("/address", middleware.Auth, controllers.DeleteAddress)

	//---------------Rating
	r.POST("/rating", middleware.Auth, controllers.AddRating)
	r.PUT("/rating", middleware.Auth, controllers.EditRating)

	//---------------Cart
	r.POST("/cart", middleware.Auth, controllers.AddCart)
	r.GET("/cart", middleware.Auth, controllers.ShowCart)
	r.PATCH("/cart", middleware.Auth, controllers.LessCart)
	r.DELETE("/cart", middleware.Auth, controllers.DeleteCart)

	//---------------Order
	r.POST("/cart/checkout", middleware.Auth, controllers.CheckoutCart)
	r.PATCH("/order", middleware.Auth, controllers.CancelOrder)
	r.GET("/order", middleware.Auth, controllers.ShowOrder)

	//---------------Payment
	r.GET("/payment", middleware.Auth, controllers.RazorPay)
	r.POST("/payment", middleware.Auth, controllers.RazorPayVerify)

	//---------------Wishlist
	r.GET("/wishlist", middleware.Auth, controllers.ShowWishlist)
	r.POST("/wishlist", middleware.Auth, controllers.AddWishlist)
	r.DELETE("/wishlist", middleware.Auth, controllers.RemoveWishlist)
}

package routers

import (
	controllers "project/controllers/admin"
	"project/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouters(r *gin.RouterGroup) {

	//---------------Login
	r.POST("/login", controllers.PostLoginA)
	r.DELETE("/logout", middleware.Auth, controllers.LogoutA)

	//---------------Product
	r.GET("/product", middleware.Auth, controllers.ShowProduct)
	r.POST("/product", middleware.Auth, controllers.AddProduct)
	r.PUT("/product/:Name", middleware.Auth, controllers.EditProduct)
	r.DELETE("/product/:Name", middleware.Auth, controllers.DeleteProduct)

	//---------------Category
	r.GET("/category", middleware.Auth, controllers.ShowCategory)
	r.POST("/category", middleware.Auth, controllers.AddCtgry)
	r.PUT("/category/:Name", middleware.Auth, controllers.EditCategory)
	r.DELETE("/category/:Name", middleware.Auth, controllers.DeleteCategory)
	r.PATCH("/category/:Name", middleware.Auth, controllers.BlockingCategory)

	//---------------User
	r.GET("/user", middleware.Auth, controllers.ShowUser)
	r.PATCH("/user/:Id", middleware.Auth, controllers.BlockingUser)

	//---------------Coupon
	r.GET("/coupon", middleware.Auth, controllers.ShowCoupon)
	r.POST("/coupon", middleware.Auth, controllers.AddCoupon)
	r.PUT("/coupon/:Id", middleware.Auth, controllers.EditCoupon)
	r.DELETE("/coupon/:Id", middleware.Auth, controllers.DeleteCoupon)
}

package routers

import (
	controllers "github.com/rishad004/My-Ecommerce/controllers/admin"
	"github.com/rishad004/My-Ecommerce/middleware"

	"github.com/gin-gonic/gin"
)

func AdminRouters(r *gin.RouterGroup) {

	//---------------Login
	r.POST("/login", controllers.PostLoginA)
	r.DELETE("/logout", middleware.Auth, controllers.LogoutA)

	//---------------Product
	r.GET("/product", middleware.Auth, controllers.ShowProduct)
	r.POST("/product", middleware.Auth, controllers.AddProduct)
	r.PUT("/product", middleware.Auth, controllers.EditProduct)
	r.DELETE("/product", middleware.Auth, controllers.DeleteProduct)

	//---------------Category
	r.GET("/category", middleware.Auth, controllers.ShowCategory)
	r.POST("/category", middleware.Auth, controllers.AddCtgry)
	r.PUT("/category", middleware.Auth, controllers.EditCategory)
	r.DELETE("/category", middleware.Auth, controllers.DeleteCategory)
	r.PATCH("/category", middleware.Auth, controllers.BlockingCategory)

	//---------------User
	r.GET("/user", middleware.Auth, controllers.ShowUser)
	r.PATCH("/user", middleware.Auth, controllers.BlockingUser)

	//---------------Coupon
	r.GET("/coupon", middleware.Auth, controllers.ShowCoupon)
	r.POST("/coupon", middleware.Auth, controllers.AddCoupon)
	r.PUT("/coupon", middleware.Auth, controllers.EditCoupon)
	r.DELETE("/coupon", middleware.Auth, controllers.DeleteCoupon)

	//---------------Order
	r.GET("/order", middleware.Auth, controllers.ShowOrders)
	r.PATCH("/order", middleware.Auth, controllers.OrdersStatusChange)

	//---------------Sale Report
	r.GET("/report", middleware.Auth, controllers.GetReportData)

	//---------------Admin Dashboard
	r.GET("/dashboard", middleware.Auth, controllers.Dashboard)
}

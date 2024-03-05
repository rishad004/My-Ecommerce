package routers

import (
	controllers "project/controllers/admin"

	"github.com/gin-gonic/gin"
)

func AdminRouters(r *gin.RouterGroup) {

	//---------------Login
	r.POST("/login", controllers.PostLoginA)

	//---------------Product
	r.GET("/product", controllers.ShowProduct)
	r.POST("/product", controllers.AddProduct)
	r.PUT("/product/:Name", controllers.EditProduct)
	r.DELETE("/product/:Name", controllers.DeleteProduct)

	//---------------Category
	r.GET("/category", controllers.ShowCategory)
	r.POST("/category", controllers.AddCtgry)
	r.PUT("/category/:Name", controllers.EditCategory)
	r.DELETE("/category/:Name", controllers.DeleteCategory)
	r.PATCH("/category/:Name", controllers.BlockingCategory)

	//---------------User
	r.GET("/user", controllers.ShowUser)
	r.PATCH("/user/:Id", controllers.BlockingUser)

	//---------------Coupon
	r.GET("/coupon", controllers.ShowCoupon)
	r.POST("/coupon", controllers.AddCoupon)
	r.PUT("/coupon/:Id", controllers.EditCoupon)
	r.DELETE("/coupon/:Id", controllers.DeleteCoupon)
}

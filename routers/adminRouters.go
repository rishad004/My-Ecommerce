package routers

import (
	"project/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRouters(r *gin.RouterGroup) {

	//---------------Login
	r.POST("/login", controllers.PostLoginA)

	//---------------Product
	r.POST("/product", controllers.AddProduct)
	r.PUT("/product/:Name", controllers.EditProduct)
	r.DELETE("/product/:Name", controllers.DeleteProduct)

	//---------------Category
	r.POST("/category", controllers.AddCtgry)
	r.PUT("/category/:Name", controllers.EditCategory)
	r.DELETE("/category/:Name", controllers.DeleteCategory)
	r.PATCH("/category/:Name", controllers.BlockingCategory)

	//---------------User
	r.PATCH("/user/:Id",controllers.BlockingUser)
}

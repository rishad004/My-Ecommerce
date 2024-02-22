package routers

import (
	"project/controllers"

	"github.com/gin-gonic/gin"
)

func AdminRouters(r *gin.RouterGroup) {

	//---------------Login
	r.POST("/login", controllers.PostLoginA)

	//---------------Product
	r.POST("/addproduct", controllers.AddProduct)
	r.PUT("/editproduct/:Name", controllers.EditProduct)
	r.DELETE("/deleteproduct/:Name", controllers.DeleteProduct)

	//---------------Category
	r.POST("/addcategory", controllers.AddCtgry)
	r.PUT("/editcategory/:Name", controllers.EditCategory)
	r.DELETE("/deletecategory/:Name", controllers.DeleteCategory)
	r.PATCH("/blockingcategory/:Name", controllers.BlockingCategory)

	//---------------User
	r.PATCH("/blockinguser/:Id",controllers.BlockingUser)
}

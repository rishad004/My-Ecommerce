package routers

import (
	"project/controllers"

	"github.com/gin-gonic/gin"
)

func UserRouters(r *gin.RouterGroup) {

	//---------------Login&SignUp
	r.POST("/signup", controllers.PostSignupU)
	r.POST("/otp", controllers.PostOtpU)
	r.POST("/login", controllers.PostLoginU)

}

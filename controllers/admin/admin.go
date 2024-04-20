package controllers

import (
	"fmt"
	"time"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/middleware"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// PostLoginA godoc
// @Summary Admin Login
// @Description Admin Login with email and password
// @Tags Admin
// @Accept  multipart/form-data
// @Produce  json
// @Param email formData string true "admin email"
// @Param pass formData string true "admin password"
// @Router /admin/login [post]
func PostLoginA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING IN----------------------")

	var adminlog models.Admin
	var check models.Admin

	adminlog.Email = c.Request.FormValue("email")
	adminlog.Pass = c.Request.FormValue("pass")

	database.Db.First(&check, "Email=?", adminlog.Email)
	err := bcrypt.CompareHashAndPassword([]byte(check.Pass), []byte(adminlog.Pass))
	if err != nil {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Invalid Email or Password!",
			"Data":    gin.H{},
		})
	} else {
		token, err := middleware.JwtCreate(c, check.Id, check.Email, "Admin")
		if err != nil {
			fmt.Println("=======Error JWT Create", err)
			c.JSON(403, gin.H{
				"Status":  "Error!",
				"Code":    400,
				"Error":   err.Error(),
				"Message": "Failed to create Token!",
				"Data":    gin.H{},
			})
			return
		}
		c.SetCookie("Jwt-Admin", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Successfully Logged in!",
			"Data": gin.H{
				"Token": token,
			},
		})
	}
}

// LogoutA godoc
// @Summary Admin Logout
// @Description Admin Logout by clearing cookie
// @Tags Admin
// @Produce  json
// @Router /admin/logout [delete]
func LogoutA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING OUT----------------------")

	c.SetCookie("Jwt-Admin", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Logged out successfully!",
		"Data":    gin.H{},
	})
}

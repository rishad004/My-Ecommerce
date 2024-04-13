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
// @Accept  json
// @Produce  json
// @Param adminlog body models.Login true "Admin Login"
// @Router /admin/login [post]
func PostLoginA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING IN----------------------")

	var adminlog models.Login
	var check models.Admin

	c.BindJSON(&adminlog)

	database.Db.First(&check, "Email=?", adminlog.Email)
	err := bcrypt.CompareHashAndPassword([]byte(check.Pass), []byte(adminlog.Password))
	if err != nil {
		c.JSON(401, gin.H{"Message": "Invalid Email or Password", "Error": err.Error()})
	} else {
		token, err := middleware.JwtCreate(c, check.Id, check.Email, "Admin")
		if err != nil {
			fmt.Println("=======Error JWT Create", err)
			c.JSON(403, gin.H{"Error": "Failed to create Token"})
			return
		}
		c.SetCookie("Jwt-Admin", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
		c.JSON(200, gin.H{"Message": "Admin Login Successfull", "token": token})
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
	c.JSON(200, gin.H{"Message": "Logged out successfully."})
}

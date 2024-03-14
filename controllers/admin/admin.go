package controllers

import (
	"fmt"
	"project/database"
	"project/middleware"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PostLoginA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING IN----------------------")

	var adminlog models.Admin
	var check models.Admin

	c.BindJSON(&adminlog)

	database.Db.First(&check, "Email=?", adminlog.Email)
	err := bcrypt.CompareHashAndPassword([]byte(check.Pass), []byte(adminlog.Pass))
	if err != nil {
		c.JSON(401, gin.H{"message": "Invalid Email or Password"})
	} else {
		middleware.JwtCreate(c, check.Id, check.Email, "Admin")
		c.JSON(200, gin.H{"message": "Admin Login Successfull"})
	}
}

func LogoutA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING OUT----------------------")

	tokenString := c.MustGet("token").(string)
	middleware.BlacklistedTokens[tokenString] = true
	c.JSON(200, gin.H{"message": "Logged out successfully."})
}

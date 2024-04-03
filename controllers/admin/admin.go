package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/middleware"
	"github.com/rishad004/My-Ecommerce/models"
	"time"

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
		c.JSON(401, gin.H{"Message": "Invalid Email or Password", "Error": err.Error()})
	} else {
		token, err := middleware.JwtCreate(c, check.Id, check.Email, "Admin")
		if err != nil {
			fmt.Println("=======Error JWT Create", err)
			c.JSON(403, gin.H{
				"Error": "Failed to create Token",
			})
			return
		}
		c.SetCookie("Jwt-Admin", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
		c.JSON(200, gin.H{"Message": "Admin Login Successfull", "token": token})
	}
}

func LogoutA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING OUT----------------------")

	c.SetCookie("Jwt-Admin", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"Message": "Logged out successfully."})
}

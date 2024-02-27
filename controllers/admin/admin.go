package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PostLoginA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING IN----------------------")

	var adminlog models.Users
	var check models.Users

	c.BindJSON(&adminlog)

	database.Db.First(&check, "Email=?", adminlog.Email)

	err := bcrypt.CompareHashAndPassword([]byte(check.Pass), []byte(adminlog.Pass))

	if err != nil {
		c.JSON(401, "Invalid Email or Password")
	} else {
		if !check.Admin {
			c.JSON(401, "Not an Admin")
		} else {
			c.JSON(200, "Admin Login Successfull")
		}
	}
}

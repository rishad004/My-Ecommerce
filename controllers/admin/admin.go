package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func PostLoginA(c *gin.Context) {
	fmt.Println("")
	fmt.Println("------------------ADMIN LOGGING IN----------------------")

	var adminlog models.Admin
	var check models.Admin

	c.BindJSON(&adminlog)

	database.Db.First(&check, "Email=?", adminlog.Email)

	if adminlog.Pass != check.Pass {
		c.JSON(401, "Invalid Email or Password")
	} else {
		c.JSON(200, "Admin Login Successfull")
	}
}

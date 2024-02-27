package controllers

import (
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func UserProfile(c *gin.Context) {

	var user models.Users
	var address []models.Address

	database.Db.First(&user, Logged)
	database.Db.Find(&address, "User_Id=?", user.ID)

	c.JSON(200, gin.H{
		"Name":   user.Name,
		"Email":  user.Email,
		"Phone":  user.Phone,
		"Gender": user.Gender,
	})
	for _, k := range address {
		c.JSON(200, gin.H{
			"Address": k,
		})
	}
}

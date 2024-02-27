package controllers

import (
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {
	var address models.Address

	c.BindJSON(&address)
	if Logged != 0 {
		address.User_Id = Logged

		database.Db.Create(&address)
		c.JSON(201, "Address added  successfully")
	} else {
		c.JSON(401, "You must be logged in to add an address!")
	}
}

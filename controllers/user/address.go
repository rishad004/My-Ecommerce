package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS ADDING------------------------")

	Logged := c.MustGet("Id").(uint)

	var address models.Address

	c.BindJSON(&address)
	if Logged != 0 {
		address.User_Id = Logged

		database.Db.Create(&address)
		c.JSON(201, gin.H{"message": "Address added  successfully"})
	} else {
		c.JSON(401, gin.H{"error": "You must be logged in to add an address!"})
	}
}

func EditAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS EDITING------------------------")

	var address, ad models.Address

	Id, _ := strconv.Atoi(c.Param("Id"))
	c.BindJSON(&ad)

	database.Db.First(&address, "Id=?", uint(Id))

	database.Db.Model(&address).Update("Name", ad.Name)
	database.Db.Model(&address).Update("Phone", ad.Phone)
	database.Db.Model(&address).Update("PinCode", ad.PinCode)
	database.Db.Model(&address).Update("City", ad.City)
	database.Db.Model(&address).Update("State", ad.State)
	database.Db.Model(&address).Update("Landmark", ad.Landmark)
	database.Db.Model(&address).Update("Address", ad.Address)

	c.JSON(200, gin.H{"message": "The Address has been updated."})
}

func DeleteAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS DELETING------------------------")

	Id, _ := strconv.Atoi(c.Param("Id"))

	var ad models.Address
	database.Db.First(&ad, "Id=?", uint(Id))

	database.Db.Delete(&ad)

	c.JSON(200, gin.H{"message": "Address deleted  successfully"})
}

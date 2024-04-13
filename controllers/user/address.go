package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

func AddAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS ADDING------------------------")

	Logged := c.MustGet("Id").(uint)

	var address models.Address

	c.BindJSON(&address)

	address.User_Id = Logged

	if err := database.Db.Create(&address).Error; err != nil {
		c.JSON(500, gin.H{
			"Status":  "Fail!",
			"Code":    500,
			"Message": "Address coudln't add!",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(201, gin.H{
		"Status":  "Success!",
		"Code":    201,
		"Message": "Address added successfully",
		"Data":    address,
	})
}

func EditAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS EDITING------------------------")

	var address, ad models.Address

	Id, _ := strconv.Atoi(c.Param("Id"))
	Logged := c.MustGet("Id").(uint)
	c.BindJSON(&ad)

	if err := database.Db.First(&address, "Id=?", uint(Id)).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Address found!",
			"Data":    gin.H{},
		})
		return
	}

	if address.User_Id != Logged {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Address found!",
			"Data":    gin.H{},
		})
		return
	}

	database.Db.Model(&address).Update("Name", ad.Name)
	database.Db.Model(&address).Update("Phone", ad.Phone)
	database.Db.Model(&address).Update("PinCode", ad.PinCode)
	database.Db.Model(&address).Update("City", ad.City)
	database.Db.Model(&address).Update("State", ad.State)
	database.Db.Model(&address).Update("Landmark", ad.Landmark)
	database.Db.Model(&address).Update("Address", ad.Address)

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Address  Updated Successfully.",
		"Data":    ad,
	})
}

func DeleteAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS DELETING------------------------")

	Id, _ := strconv.Atoi(c.Param("Id"))
	Logged := c.MustGet("Id").(uint)

	var ad models.Address
	err := database.Db.Where("Id=?", uint(Id)).First(&ad)

	if err.Error != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Address found!",
			"Data":    gin.H{},
		})
		return
	}
	if ad.User_Id != Logged {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Address found!",
			"Data":    gin.H{},
		})
		return
	}

	database.Db.Delete(&ad)

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Address Deleted Successfully.",
		"Data":    gin.H{},
	})
}

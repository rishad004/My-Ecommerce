package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// AddAddress godoc
// @Summary Address Add
// @Description Adding Address with it's details
// @Tags User Address
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "Name"
// @Param phone formData string true "Phone"
// @Param pincode formData string true "Pincode"
// @Param city formData string true "City"
// @Param state formData string true "State"
// @Param landmark formData string true "Landmark"
// @Param address formData string true "Address"
// @Router /user/address [post]
func AddAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS ADDING------------------------")

	Logged := c.MustGet("Id").(uint)

	var address models.Address

	address.Name = c.Request.FormValue("name")
	address.Phone = c.Request.FormValue("phone")
	address.PinCode = c.Request.FormValue("pincode")
	address.City = c.Request.FormValue("city")
	address.State = c.Request.FormValue("state")
	address.Landmark = c.Request.FormValue("landmark")
	address.Address = c.Request.FormValue("address")

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

// EditAddress godoc
// @Summary Address Edit
// @Description Editing Address with it's details
// @Tags User Address
// @Accept  multipart/form-data
// @Produce  json
// @Param id query string true "address id"
// @Param name formData string true "Name"
// @Param phone formData string true "Phone"
// @Param pincode formData string true "Pincode"
// @Param city formData string true "City"
// @Param state formData string true "State"
// @Param landmark formData string true "Landmark"
// @Param address formData string true "Address"
// @Router /user/address [put]
func EditAddress(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ADDRESS EDITING------------------------")

	var address models.Address

	Id, _ := strconv.Atoi(c.Query("Id"))
	Logged := c.MustGet("Id").(uint)

	if err := database.Db.First(&address, "Id=?", uint(Id)).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Address found!",
			"Error":   err.Error(),
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
	address.Name = c.Request.FormValue("name")
	address.Phone = c.Request.FormValue("phone")
	address.PinCode = c.Request.FormValue("pincode")
	address.City = c.Request.FormValue("city")
	address.State = c.Request.FormValue("state")
	address.Landmark = c.Request.FormValue("landmark")
	address.Address = c.Request.FormValue("address")
	if err := database.Db.Save(&address).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Message": "Address  Not Updated!",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Address  Updated Successfully.",
		"Data":    address,
	})
}

// DeleteAddress godoc
// @Summary Address Delete
// @Description Deleting address completely
// @Tags User Address
// @Produce  json
// @Param id query string true "name search by id"
// @Router /user/address [delete]
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

package controllers

import (
	"fmt"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/helper"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// ProfileShow godoc
// @Summary Show Profile
// @Description Showing User Profile
// @Tags User Profile
// @Router /user/profile [get]
func UserProfile(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------USER PROFILE------------------------")

	var use models.Users
	var address []models.Address
	var addresShow []gin.H

	Logged := c.MustGet("Id").(uint)

	database.Db.First(&use, Logged)
	database.Db.Find(&address, "User_Id=?", Logged)

	personShow := gin.H{
		"Name":   use.Name,
		"Email":  use.Email,
		"Phone":  use.Phone,
		"Gender": use.Gender,
	}
	for _, k := range address {
		addresShow = append(addresShow, gin.H{
			"Id":       k.Id,
			"Name":     k.Name,
			"Phone":    k.Phone,
			"PinCode":  k.PinCode,
			"City":     k.City,
			"State":    k.State,
			"Landmark": k.Landmark,
			"Address":  k.Address,
		})
	}
	c.JSON(200, gin.H{
		"User":    personShow,
		"Address": addresShow,
	})
}

// UpdatePass godoc
// @Summary Updating Password
// @Description Updating user password
// @Tags User Profile
// @Param address formData string true "address id"
// @Param coupon formData string true "coupon code"
// @Param method formData string true "payment method"
// @Produce  json
// @Router /user/password [patch]
func UpdatePass(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CHANGING PASSWORD------------------------")

	CurrentPass := c.Request.FormValue("currentpass")
	NewPass := c.Request.FormValue("newpass")
	RepeatPass := c.Request.FormValue("repeatpass")

	var use models.Users
	Logged := c.MustGet("Id").(uint)

	database.Db.First(&use, Logged)
	err := bcrypt.CompareHashAndPassword([]byte(use.Pass), []byte(CurrentPass))

	if err != nil {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Current password is not correct!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}
	if NewPass != RepeatPass {
		c.JSON(406, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Both new and repeat pass are not the same!",
			"Data":    gin.H{},
		})
		return
	}
	use.Pass = helper.HashPass(NewPass)
	database.Db.Save(&use)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Successfully updated your password!",
		"Data":    gin.H{},
	})
}

// EditProfile godoc
// @Summary Editing Profile
// @Description Editing user profile
// @Tags User Profile
// @Param name formData string true "user's name"
// @Param phone formData string true "user's phone"
// @Param gender formData string true "user's gender"
// @Produce  json
// @Router /user/profile [patch]
func EditProfile(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------EDITING PROFILE------------------------")

	var use models.Users

	Logged := c.MustGet("Id").(uint)
	database.Db.First(&use, Logged)
	use.Name = c.Request.FormValue("name")
	use.Phone = c.Request.FormValue("phone")
	use.Gender = c.Request.FormValue("gender")
	err := database.Db.Save(&use)

	if err.Error != nil {
		c.JSON(409, gin.H{
			"Status":  "Fail!",
			"Code":    409,
			"Message": "User already exist with this number!",
			"Error":   err.Error,
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Profile udpated successfully!",
		"Data":    gin.H{},
	})
}

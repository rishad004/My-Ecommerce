package controllers

import (
	"fmt"
	"project/database"
	"project/helper"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

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

func UpdatePass(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CHANGING PASSWORD------------------------")

	type Pass struct {
		CurrentPass string `json:"currentpass"`
		NewPass     string `json:"newpass"`
		RepeatPass  string `json:"repeatpass"`
	}
	var data Pass
	var use models.Users
	Logged := c.MustGet("Id").(uint)
	c.BindJSON(&data)

	database.Db.First(&use, Logged)
	err := bcrypt.CompareHashAndPassword([]byte(use.Pass), []byte(data.CurrentPass))

	if err != nil {
		c.JSON(401, gin.H{"Message": "Current password is not correct"})
		return
	}
	if data.NewPass != data.RepeatPass {
		c.JSON(406, gin.H{"Message": "Both new and repeat pass are not the same"})
		return
	}
	use.Pass = helper.HashPass(data.NewPass)
	database.Db.Save(&use)
	c.JSON(200, gin.H{"Message": "Successfully updated your password"})
}

func EditProfile(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------EDITING PROFILE------------------------")

	var us, use models.Users

	c.BindJSON(&us)

	Logged := c.MustGet("Id").(uint)
	database.Db.First(&use, Logged)
	use.Name = us.Name
	use.Phone = us.Phone
	use.Gender = us.Gender
	err := database.Db.Save(&use)

	if err.Error != nil {
		c.JSON(409, gin.H{"Message": "User already exist with this number."})
		return
	}
	c.JSON(200, gin.H{"Message": "Profile udpated successfully"})
}

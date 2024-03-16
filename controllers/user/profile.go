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

	var user models.Users
	var address []models.Address

	Logged := c.MustGet("Id").(uint)

	database.Db.First(&user, Logged)
	database.Db.Find(&address, "User_Id=?", Logged)

	c.JSON(200, gin.H{
		"Name":   user.Name,
		"Email":  user.Email,
		"Phone":  user.Phone,
		"Gender": user.Gender,
	})
	for _, k := range address {
		c.JSON(200, gin.H{
			"Name":     k.Name,
			"Phone":    k.Phone,
			"Pincode":  k.PinCode,
			"City":     k.City,
			"State":    k.State,
			"Landmark": k.Landmark,
			"Address":  k.Address,
		})
	}
}

func UpdatePass(c *gin.Context) {
	type Pass struct {
		CurrentPass string `json:"currentpass"`
		NewPass     string `json:"newpass"`
		RepeatPass  string `json:"repeatpass"`
	}
	var data Pass
	var user models.Users
	Logged := c.MustGet("Id").(uint)
	c.BindJSON(&data)

	database.Db.First(&user, Logged)
	err := bcrypt.CompareHashAndPassword([]byte(user.Pass), []byte(data.CurrentPass))

	if err != nil {
		c.JSON(401, gin.H{"message": "Current password is not correct"})
		return
	}
	if data.NewPass != data.RepeatPass {
		c.JSON(406, gin.H{"message": "Both new and repeat pass are not the same"})
		return
	}
	user.Pass = helper.HashPass(data.NewPass)
	database.Db.Save(&user)
	c.JSON(200, gin.H{"message": "Successfully updated your password"})
}

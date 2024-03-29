package controllers

import (
	"fmt"
	"project/database"
	"project/helper"
	"project/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type addres struct {
	Id       uint
	Name     string
	Phone    uint
	PinCode  uint
	City     string
	State    string
	Landmark string
	Address  string
}

type person struct {
	Name   string
	Email  string
	Phone  string
	Gender string
}

func UserProfile(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------USER PROFILE------------------------")

	var use models.Users
	var address []models.Address
	var addresShow []addres

	Logged := c.MustGet("Id").(uint)

	database.Db.First(&use, Logged)
	database.Db.Find(&address, "User_Id=?", Logged)

	personShow := person{
		Name:   use.Name,
		Email:  use.Email,
		Phone:  use.Phone,
		Gender: use.Gender,
	}
	for _, k := range address {
		addresShow = append(addresShow, addres{k.Id, k.Name, k.Phone, k.PinCode, k.City, k.State, k.Landmark, k.Address})
	}
	c.JSON(200, gin.H{
		"1user":    personShow,
		"2address": addresShow,
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
		c.JSON(401, gin.H{"message": "Current password is not correct"})
		return
	}
	if data.NewPass != data.RepeatPass {
		c.JSON(406, gin.H{"message": "Both new and repeat pass are not the same"})
		return
	}
	use.Pass = helper.HashPass(data.NewPass)
	database.Db.Save(&use)
	c.JSON(200, gin.H{"message": "Successfully updated your password"})
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
		c.JSON(409, gin.H{"message": "User already exist with this number."})
		return
	}
	c.JSON(200, gin.H{"message": "Profile udpated successfully"})
}

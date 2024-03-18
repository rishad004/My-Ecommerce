package controllers

import (
	"fmt"
	"math/rand"
	"project/database"
	"project/helper"
	"project/middleware"
	"project/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var user models.Users

func PostSignupU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER SIGNING UP----------------------")

	var otp models.Otp
	var oottpp models.Otp

	c.ShouldBindJSON(&user)

	user.Blocking = true

	user.Pass = helper.HashPass(user.Pass)

	randomNum := strconv.Itoa(rand.Intn(900000) + 100000)

	database.Db.First(&oottpp, "User_Mail=?", user.Email)
	database.Db.Delete(&oottpp)
	otp = models.Otp{
		Id:       oottpp.Id,
		UserMail: user.Email,
		Otp:      randomNum,
		Expr:     time.Now().Add(2 * time.Minute),
	}
	database.Db.Save(&otp)

	helper.SendMail(c, user.Email, "Otp", "Your verification code is "+otp.Otp)

	c.JSON(200, "Verify Otp, Please check your mail. "+otp.Otp)
}

func PostOtpU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------OTP VERIFYING----------------------")

	var check models.Otp
	var rc models.Otp

	c.BindJSON(&rc)

	database.Db.Where("User_Mail=? AND Expr > ?", user.Email, time.Now()).First(&check)

	if check.Otp == rc.Otp {
		database.Db.Delete(&check)

		err := database.Db.Create(&user)

		if err.Error != nil {
			c.JSON(409, "User already exist")
		} else {
			c.JSON(200, gin.H{
				"message": "Successfully signed up",
				"userId":  user.ID,
			})
		}
	} else {
		c.JSON(401, "Otp expired or invalid, Please try again")
	}

}

func PostLoginU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING IN----------------------")

	var userlog models.Users
	var check models.Users

	c.BindJSON(&userlog)

	database.Db.First(&check, "Email=?", userlog.Email)

	err := bcrypt.CompareHashAndPassword([]byte(check.Pass), []byte(userlog.Pass))

	if !check.Blocking {
		c.JSON(401, gin.H{"message": "User blocked by admin"})
	} else {
		if err != nil {
			c.JSON(401, gin.H{"message": "Inavlid Email or Password"})
		} else {
			token, erro := middleware.JwtCreate(c, check.ID, check.Email, "User")
			if erro != nil {
				fmt.Println("=======Error JWT Create", err)
				c.JSON(403, gin.H{"Error": "Failed to create Token"})
				return
			}
			c.JSON(200, gin.H{"message": "Successfully Logged in", "token": token})
		}
	}
}

func LogoutU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING OUT----------------------")

	tokenString := c.MustGet("token").(string)
	middleware.BlacklistedTokens[tokenString] = true
	c.JSON(200, gin.H{"message": "Logged out successfully."})
}

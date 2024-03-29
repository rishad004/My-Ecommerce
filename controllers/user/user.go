package controllers

import (
	"fmt"
	"project/database"
	"project/helper"
	"project/middleware"
	"project/models"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func PostSignupU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER SIGNING UP----------------------")

	var user models.Users
	var otp models.Otp
	var oottpp models.Otp

	c.ShouldBindJSON(&user)
	session := sessions.Default(c)

	user.Pass = helper.HashPass(user.Pass)

	randomNum := helper.GenerateInt()

	database.Db.First(&oottpp, "User_Mail=?", user.Email)
	database.Db.Delete(&oottpp)
	otp = models.Otp{
		Id:       oottpp.Id,
		UserMail: user.Email,
		Otp:      randomNum,
		Expr:     time.Now().Add(2 * time.Minute),
	}
	database.Db.Save(&otp)

	err := helper.SendMail(c, user.Email, "Otp", "Your verification code is "+otp.Otp)
	if err != nil {
		c.JSON(503, gin.H{
			"Message": "We couldn't send the mail, Please check email address-----------------",
			"Error":   err,
		})
		return
	}
	session.Set("signupEmail", user.Email)
	session.Set("signupName", user.Name)
	session.Set("signupPhone", user.Phone)
	session.Set("signupPass", user.Pass)
	session.Set("signupGender", user.Gender)
	session.Save()

	c.JSON(200, "Verify Otp, Please check your mail. "+otp.Otp)
}

func PostOtpU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------OTP VERIFYING----------------------")

	var check models.Otp
	var rc models.Otp
	var user models.Users

	c.BindJSON(&rc)
	session := sessions.Default(c)
	user.Email = session.Get("signupEmail").(string)
	user.Name = session.Get("signupName").(string)
	user.Phone = session.Get("signupPhone").(string)
	user.Gender = session.Get("signupGender").(string)
	user.Pass = session.Get("signupPass").(string)
	user.Blocking = true
	fmt.Println(user)

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
			session.Delete("signupEmail")
			session.Delete("signupName")
			session.Delete("signupPhone")
			session.Delete("signupGender")
			session.Delete("signupPass")
			session.Save()
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
			c.SetCookie("Jwt-User", token, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
			c.JSON(200, gin.H{"Message": "Successfully Logged in", "Token": token, "Id": check.ID})
		}
	}
}

func LogoutU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING OUT----------------------")

	c.SetCookie("Jwt-User", "", -1, "/", "localhost", false, true)
	c.JSON(200, gin.H{"message": "Logged out successfully."})
}

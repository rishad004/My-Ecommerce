package controllers

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

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
		"Status":  "Success!",
		"Code":    200,
		"Message": "User profile showing",
		"Data": gin.H{
			"User":    personShow,
			"Address": addresShow,
		},
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

// FogotPassInit godoc
// @Summary Forgot Pass Init
// @Description Initialilzing forgot password without login
// @Tags User ForgetPass
// @Accept  multipart/form-data
// @Produce  json
// @Param email formData string true "User email"
// @Router /user/reset/password [post]
func ForgotPassword(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------FORGOT PASSWORD----------------------")

	var user models.Users
	var otp models.Otp

	email := c.Request.FormValue("email")

	if err := database.Db.First(&user, "Email=?", email).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "User not found with this email!",
			"Data":    gin.H{},
		})
		return
	}

	rand := helper.GenerateInt()

	otp = models.Otp{
		UserMail: email,
		Otp:      rand,
		Expr:     time.Now().Add(2 * time.Minute),
	}

	if err := database.Db.Create(&otp).Error; err != nil {
		c.JSON(409, gin.H{
			"Status":  "Error!",
			"Code":    409,
			"Error":   err.Error(),
			"Message": "Couldn't create otp!",
			"Data":    gin.H{},
		})
		return
	}

	code := base64.StdEncoding.EncodeToString([]byte(email + "_" + rand))
	fmt.Println("")
	fmt.Println(code)

	if err := helper.SendMail(email, "Reset your Password", "https://byecom.shop/user/reset?code="+code); err != nil {
		c.JSON(503, gin.H{
			"Status":  "Fail!",
			"Code":    503,
			"Error":   err.Error(),
			"Message": "We couldn't send the mail, Please check email address.",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "We have sent you a link to reset, Please check your email!",
		"Data": gin.H{
			"Link": "https://byecom.shop/user/reset?code=" + code,
		},
	})

}

// ChangePass godoc
// @Summary Changing Pass
// @Description Changing user password without login
// @Tags User ForgetPass
// @Accept  multipart/form-data
// @Produce  json
// @Param code query string true "Reset code"
// @Param newpass formData string true "New Pass"
// @Param repeatpass formData string true "Repeat Pass"
// @Router /user/reset [post]
func ResetPassword(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------FORGOT PASSWORD----------------------")

	var user models.Users
	var otp models.Otp

	code := c.Query("code")

	decodedBytes, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Decode error!",
			"Data":    gin.H{},
		})
	}
	fmt.Println(string(decodedBytes))
	detail := strings.Split(string(decodedBytes), "_")
	email := detail[0]
	rand := detail[1]
	fmt.Println(email, "  ", rand)

	if err := database.Db.First(&user, "Email=?", email).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Invalid User or Code!",
			"Data":    gin.H{},
		})
		return
	}

	if err := database.Db.Where("Otp=? AND Expr > ?", rand, time.Now()).First(&otp).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Invalid User or Code!",
			"Data":    gin.H{},
		})
		return
	}

	if otp.UserMail != email {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Invalid User or Code!",
			"Data":    gin.H{},
		})
		return
	}

	new := c.Request.FormValue("newpass")
	repeat := c.Request.FormValue("repeatpass")

	if new != repeat || new == "" || repeat == "" {
		c.JSON(406, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Both new and repeat pass are not the same or password isn't in the format!",
			"Data":    gin.H{},
		})
		return
	}
	user.Pass = helper.HashPass(new)
	if err := database.Db.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Fail!",
			"Code":    400,
			"Message": "Couldn't update new password!",
			"Data":    gin.H{},
		})
	}

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Password updated successfully!",
		"Data":    gin.H{},
	})
}

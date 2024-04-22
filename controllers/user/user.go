package controllers

import (
	"fmt"
	"time"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/helper"
	"github.com/rishad004/My-Ecommerce/middleware"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignupUser godoc
// @Summary Signup User
// @Description Sigining up user
// @Tags User Login&Signup
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "User name"
// @Param email formData string true "User email"
// @Param pass formData string true "User pass"
// @Param phone formData string true "User phone"
// @Param gender formData string true "User gender"
// @Router /user/signup [post]
func PostSignupU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER SIGNING UP----------------------")

	var user models.Users
	var otp models.Otp
	var oottpp models.Otp

	user.Email = c.Request.FormValue("email")
	user.Name = c.Request.FormValue("name")
	user.Phone = c.Request.FormValue("phone")
	user.Gender = c.Request.FormValue("gender")
	user.Pass = c.Request.FormValue("pass")

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

	err := helper.SendMail(user.Email, "Otp", "Your verification code is "+otp.Otp)
	if err != nil {
		c.JSON(503, gin.H{
			"Status":  "Fail!",
			"Code":    503,
			"Error":   err.Error(),
			"Message": "We couldn't send the mail, Please check email address.",
			"Data":    gin.H{},
		})
		return
	}
	session.Set("signupEmail", user.Email)
	session.Set("signupName", user.Name)
	session.Set("signupPhone", user.Phone)
	session.Set("signupPass", user.Pass)
	session.Set("signupGender", user.Gender)
	session.Save()

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Verify Email, Please check your mail!",
		"Data": gin.H{
			"Otp": otp.Otp,
		},
	})
}

// OtpVerify godoc
// @Summary Otp verify
// @Description Verifying the user otp and creating user
// @Tags User Login&Signup
// @Accept  multipart/form-data
// @Produce  json
// @Param otp formData string true "User Otp"
// @Router /user/signup/otp [post]
func PostOtpU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------OTP VERIFYING----------------------")

	var check models.Otp
	var user models.Users
	var wallet models.Wallet

	otp := c.Request.FormValue("otp")

	session := sessions.Default(c)
	user.Email = session.Get("signupEmail").(string)
	user.Name = session.Get("signupName").(string)
	user.Phone = session.Get("signupPhone").(string)
	user.Gender = session.Get("signupGender").(string)
	user.Pass = session.Get("signupPass").(string)
	user.Blocking = true
	fmt.Println(user)

	database.Db.Where("User_Mail=? AND Expr > ?", user.Email, time.Now()).First(&check)

	if check.Otp == otp {
		database.Db.Delete(&check)

		if err := database.Db.Create(&user).Error; err != nil {
			c.JSON(409, gin.H{
				"Status":  "Fail!",
				"Code":    409,
				"Error":   err.Error(),
				"Message": "User already exist!",
				"Data":    gin.H{},
			})
			return
		}

		wallet.UserId = user.ID
		if er := database.Db.Create(&wallet).Error; er != nil {
			c.JSON(500, gin.H{"Error": "Couldn't create the wallet!"})
			return
		}
		session.Delete("signupEmail")
		session.Delete("signupName")
		session.Delete("signupPhone")
		session.Delete("signupGender")
		session.Delete("signupPass")
		session.Save()

		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Successfully signed up!",
			"Data": gin.H{
				"UserId": user,
			},
		})
	} else {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Otp expired or invalid, Please try again!",
			"Data":    gin.H{},
		})
	}

}

// LoginUser godoc
// @Summary Login User
// @Description Logging in user with email and pass
// @Tags User Login&Signup
// @Accept  multipart/form-data
// @Produce  json
// @Param email formData string true "User email"
// @Param pass formData string true "User pass"
// @Router /user/login [post]
func PostLoginU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING IN----------------------")

	var userlog models.Users
	var check models.Users

	userlog.Email = c.Request.FormValue("email")
	userlog.Pass = c.Request.FormValue("pass")

	if err := database.Db.First(&check, "Email=?", userlog.Email).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "User not found!",
			"Data":    gin.H{},
		})
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(check.Pass), []byte(userlog.Pass))

	if !check.Blocking {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "User blocked by admin!",
			"Data":    gin.H{},
		})
	} else {
		if err != nil {
			c.JSON(401, gin.H{
				"Status":  "Fail!",
				"Code":    401,
				"Error":   err.Error(),
				"Message": "Inavlid Email or Password!",
				"Data":    gin.H{},
			})
		} else {
			token, erro := middleware.JwtCreate(c, check.ID, check.Email, "User")
			if erro != nil {
				fmt.Println("=======Error JWT Create", err)
				c.JSON(403, gin.H{
					"Status":  "Error!",
					"Code":    400,
					"Error":   erro.Error(),
					"Message": "Failed to create Token!",
					"Data":    gin.H{},
				})
				return
			}
			c.SetCookie("Jwt-User", token, int((time.Hour * 1).Seconds()), "/", "byecom.shop", false, false)
			c.JSON(200, gin.H{
				"Status":  "Success!",
				"Code":    200,
				"Message": "Successfully Logged in!",
				"Data": gin.H{
					"Token": token,
					"Id":    check.ID,
				},
			})
		}
	}
}

// LogoutUser godoc
// @Summary Logout User
// @Description Logging out user
// @Tags User Login&Signup
// @Produce  json
// @Router /user/logout [delete]
func LogoutU(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------USER LOGGING OUT----------------------")

	c.SetCookie("Jwt-User", "", -1, "/", "byecom.shop", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Logged out successfully!",
		"Data":    gin.H{},
	})
}

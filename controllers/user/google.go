package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/helper"
	"github.com/rishad004/My-Ecommerce/middleware"
	"github.com/rishad004/My-Ecommerce/models"
)

// GoogleLogin godoc
// @Summary Google Login
// @Description Logging in/Signing up with google auth
// @Tags User Google
// @Produce  json
// @Router /user/google/login [get]
func GoogleLogin(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------GOOGLE LOGIN STARTED---------------------------")

	conf := helper.Google()

	url := conf.AuthCodeURL("randomState")
	// c.JSON(200, gin.H{
	// 	"Status":  "Success!",
	// 	"Code":    200,
	// 	"Message": "Redirecting to Google",
	// 	"Data": gin.H{
	// 		"Url": url,
	// 	},
	// })

	c.Redirect(302, url)
}

// GoogleCallback godoc
// @Summary Google Callback
// @Description Callback function after getting details
// @Tags User Google
// @Produce  json
// @Router /user/google/callback [get]
func GoogleCallback(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------GOOGLE LOGIN CHECKING------------------------")

	var User models.Users

	conf := helper.Google()

	if c.Request.URL.Query().Get("state") != "randomState" {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "State isn't valid!",
			"Data":    gin.H{},
		})
		return
	}

	token, err := conf.Exchange(context.Background(), c.Request.URL.Query().Get("code"))
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Token not found!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}
	client := conf.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Get request error!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}

	defer resp.Body.Close()

	var user map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Couldn't read body!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}

	if !user["email_verified"].(bool) {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Email not verified!",
			"Data":    gin.H{},
		})
	}

	if er := database.Db.First(&User, "Email=?", user["email"].(string)).Error; er != nil {
		User.Name = user["name"].(string)
		User.Email = user["email"].(string)
		User.Pass = user["sub"].(string)
		User.Blocking = true
		if er := database.Db.Create(&User).Error; er != nil {
			c.JSON(400, gin.H{
				"Status":  "Error!",
				"Code":    400,
				"Message": "Couldn't create user!",
				"Error":   er.Error(),
				"Data":    gin.H{},
			})
			return
		}
	}
	jwtTok, erro := middleware.JwtCreate(c, User.ID, User.Email, "User")
	if erro != nil {
		c.JSON(403, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   erro.Error(),
			"Message": "Failed to create Token!",
			"Data":    gin.H{},
		})
		return
	}
	c.SetCookie("Jwt-User", jwtTok, int((time.Hour * 1).Seconds()), "/", "localhost", false, true)
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Successfully Logged in!",
		"Data": gin.H{
			"Token": jwtTok,
			"Id":    User.ID,
		},
	})
}

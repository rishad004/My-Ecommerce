package controllers

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rishad004/My-Ecommerce/helper"
	"golang.org/x/oauth2"
)

func GoogleLogin(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------GOOGLE LOGIN STARTED---------------------------")

	conf := helper.Google()

	url := conf.AuthCodeURL("randomState")

	c.Redirect(302, url)
}

func GoogleCallback(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------GOOGLE LOGIN CHECKING------------------------")

	conf := helper.Google()

	if c.Request.FormValue("state") != "randomState" {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "State isn't valid!",
			"Data":    gin.H{},
		})
		return
	}

	token, err := conf.Exchange(oauth2.NoContext, c.Request.FormValue("code"))
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
	resp, err := http.Get("https://www.gooogleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
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

	details, err := ioutil.ReadAll(resp.Body)
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
	fmt.Println("details:==== ", details)
	c.Redirect(302, "/user/home")
}

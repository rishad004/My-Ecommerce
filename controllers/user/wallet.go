package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
)

// WalletShow godoc
// @Summary Show Wallet
// @Description Showing User Wallet
// @Tags User Wallet
// @Router /user/wallet [get]
func ShowWallet(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------WALLET SHOWING----------------------")

	Logged := c.MustGet("Id").(uint)

	var wallet models.Wallet

	if err := database.Db.First(&wallet, "User_Id=?", Logged).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Wallet not found!",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Wallet found!",
		"Data":    gin.H{"Balance": wallet.Balance},
	})
}

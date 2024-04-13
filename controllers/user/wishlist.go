package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

func AddWishlist(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------ADDING WISHLIST----------------------")

	var Wishlist models.Wishlist

	Logged := c.MustGet("Id").(uint)
	ProductID, _ := strconv.Atoi(c.Query("Id"))

	if err := database.Db.Where("Product_Id=? AND User_Id=?", uint(ProductID), Logged).First(&Wishlist).Error; err == nil {
		c.JSON(409, gin.H{
			"Status":  "Fail!",
			"Code":    409,
			"Error":   err,
			"Message": "Product already exist in wishlist!",
			"Data":    gin.H{},
		})
		return
	}
	Wishlist = models.Wishlist{
		UserId:    Logged,
		ProductId: uint(ProductID),
	}
	if err := database.Db.Create(&Wishlist).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Couldn't create the wishlist!",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Added to wishlist!",
		"Data":    gin.H{},
	})
}

func RemoveWishlist(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------REMOVING WISHLIST----------------------")

	var Wishlist models.Wishlist

	Logged := c.MustGet("Id").(uint)
	ProductID, _ := strconv.Atoi(c.Query("Id"))

	if err := database.Db.Where("Product_Id=? AND User_Id=?", uint(ProductID), Logged).First(&Wishlist).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Product not found in wishlist!",
			"Data":    gin.H{},
		})
		return
	}

	if err := database.Db.Delete(&Wishlist).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Couldn't delete the wishlist!",
			"Data":    gin.H{},
		})
		return
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Deleted from wishlist!",
		"Data":    gin.H{},
	})
}

func ShowWishlist(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------SHOWING WISHLIST----------------------")

	var wishlist []models.Wishlist
	var show []gin.H
	var img string

	Logged := c.MustGet("Id").(uint)

	if err := database.Db.Preload("Product").Where("User_Id=?", Logged).Find(&wishlist).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No products found in wishlist!",
			"Data":    gin.H{},
		})
		return
	}

	for _, v := range wishlist {
		if v.Product.ImageURLs == nil {
			img = ""
		} else {
			img = v.Product.ImageURLs[0]
		}
		show = append(show, gin.H{
			"Id":    v.Id,
			"Name":  v.Product.Name,
			"Price": v.Product.Price,
			"Image": img,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved whishlist data!",
		"Data": gin.H{
			"Products": show,
		},
	})
}

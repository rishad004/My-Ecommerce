package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddWishlist(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------ADDING WISHLIST----------------------")

	var Wishlist models.Wishlist

	Logged := c.MustGet("Id").(uint)
	ProductID, _ := strconv.Atoi(c.Query("Id"))

	if err := database.Db.Where("Product_Id=? AND User_Id=?", uint(ProductID), Logged).First(&Wishlist).Error; err == nil {
		c.JSON(409, gin.H{"Message": "Product already exist in wishlist!"})
		return
	}
	Wishlist = models.Wishlist{
		UserId:    Logged,
		ProductId: uint(ProductID),
	}
	if err := database.Db.Create(&Wishlist).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't create the wishlist!"})
		return
	}
	c.JSON(200, gin.H{"Message": "Added to wishlist!"})
}

func RemoveWishlist(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------REMOVING WISHLIST----------------------")

	var Wishlist models.Wishlist

	Logged := c.MustGet("Id").(uint)
	ProductID, _ := strconv.Atoi(c.Query("Id"))

	if err := database.Db.Where("Product_Id=? AND User_Id=?", uint(ProductID), Logged).First(&Wishlist).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Product not found in wishlist!"})
		return
	}

	if err := database.Db.Delete(&Wishlist).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't delete the wishlist!"})
		return
	}
	c.JSON(200, gin.H{"Message": "Deleted from wishlist!"})
}

func ShowWishlist(c *gin.Context) {

	fmt.Println("")
	fmt.Println("------------------SHOWING WISHLIST----------------------")

	var wishlist []models.Wishlist
	var show []gin.H
	var img string

	Logged := c.MustGet("Id").(uint)

	if err := database.Db.Preload("Product").Where("User_Id=?", Logged).Find(&wishlist).Error; err != nil {
		c.JSON(404, gin.H{"Message": "No products found in wishlist!"})
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
	c.JSON(200, gin.H{"Products": show})
}

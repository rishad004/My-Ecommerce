package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type details struct {
	Images      pq.StringArray
	Name        string
	Price       int
	Color       string
	Category    string
	Description string
	Status      string
}

type home struct {
	Image pq.StringArray
	Name  string
	Price int
	Rating int
}

func UserShowP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product models.Products
	var category models.Category
	var show details
	var s string
	var p []models.Products

	Id := c.Param("Id")

	database.Db.First(&product, Id)

	database.Db.First(&category, product.CtgryId)
	if category.Blocking {
		if product.Quantity != 0 {
			s = "On stock"
		} else {
			s = "Out of stock"
		}
		show = details{
			Images:      product.ImageURLs,
			Name:        product.Name,
			Price:       product.Price,
			Color:       product.Color,
			Category:    category.Name,
			Description: product.Dscptn,
			Status:      s,
		}
		database.Db.Where("Ctgry_Id=?", category.Id).Find(&p)
	}
	c.JSON(200, show)
	c.JSON(200, "Related Products")
	for i := 0; i < len(p); i++ {
		if p[i].Id != product.Id {
			c.JSON(200, gin.H{
				"Image": p[i].ImageURLs[0],
				"Name":  p[i].Name,
				"Price": p[i].Price,
			})
		}
	}
	product = models.Products{}
	category = models.Category{}
}

func UserSearchP(c *gin.Context) {

}

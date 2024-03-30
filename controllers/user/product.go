package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func UserShowP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product models.Products
	var category models.Category
	var show gin.H
	var s string
	var p []models.Products
	var r []models.Rating
	var relatedShow, ratingShow []gin.H

	Id, _ := strconv.Atoi(c.Param("Id"))

	err := database.Db.Where("id=?", uint(Id)).First(&product).Error
	if err != nil {
		c.JSON(404, gin.H{"Error": "Product not found!"})
		return
	}
	database.Db.Find(&r, "Prdct_Id=?", product.ID)

	database.Db.First(&category, product.CtgryId)
	if category.Blocking {
		if product.Quantity != 0 {
			s = "On stock"
		} else {
			s = "Out of stock"
		}
		for _, k := range r {
			ratingShow = append(ratingShow, gin.H{"Rating": k.Rating, "Review": k.Review})
		}
		var Avg string
		if product.AvrgRating != 0 {
			Avg = fmt.Sprint(product.AvrgRating, "/5")
		} else {
			Avg = "0 Rating"
		}
		show = gin.H{
			"Images":      product.ImageURLs,
			"Name":        product.Name,
			"Price":       product.Price,
			"Color":       product.Color,
			"Category":    category.Name,
			"Description": product.Dscptn,
			"Rating":      Avg,
			"Status":      s,
		}
		database.Db.Where("Ctgry_Id=?", category.Id).Find(&p)
	} else {
		c.JSON(404, gin.H{"Error": "Product not found!"})
		return
	}
	for i := 0; i < len(p); i++ {
		if p[i].ID != product.ID {
			relatedShow = append(relatedShow, gin.H{"Image": p[i].ImageURLs[0], "Name": p[i].Name, "Price": p[i].Price})
		}
	}
	c.JSON(200, gin.H{
		"Product":         show,
		"RatingReview":    ratingShow,
		"RelatedProducts": relatedShow,
	})
	product = models.Products{}
	category = models.Category{}
}

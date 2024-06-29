package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowProduct godoc
// @Summary Show Product
// @Description Showing Product details
// @Tags User Home&Product
// @Param id query string true "product id"
// @Produce  json
// @Router /user/product [get]
func UserShowP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product models.Products
	var category models.Category
	var show gin.H
	var s, img string
	var p []models.Products
	var r []models.Rating
	var relatedShow, ratingShow []gin.H

	Id, _ := strconv.Atoi(c.Query("id"))

	err := database.Db.Where("id=?", uint(Id)).First(&product).Error
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Message": "Product not found!",
			"Data":    gin.H{}})
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
			"Offer":       product.Offer,
			"Color":       product.Color,
			"Category":    category.Name,
			"Description": product.Dscptn,
			"Rating":      Avg,
			"Status":      s,
		}
		database.Db.Where("Ctgry_Id=?", category.Id).Find(&p)
	} else {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Message": "Product not found!",
			"Data":    gin.H{},
		})
		return
	}
	for i := 0; i < len(p); i++ {
		if p[i].ID != product.ID {
			if p[i].ImageURLs == nil {
				img = "No Photos"
			} else {
				img = p[i].ImageURLs[0]
			}
			relatedShow = append(relatedShow, gin.H{"Image": img, "Name": p[i].Name, "Price": p[i].Offer})
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved product detail with it's related products!",
		"Data": gin.H{
			"Product":         show,
			"RatingReview":    ratingShow,
			"RelatedProducts": relatedShow,
		},
	})
}

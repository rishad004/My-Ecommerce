package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type rating struct {
	Rating float32
	Review string
}

type related struct {
	Image string
	Name  string
	Price int
}

type details struct {
	Images      pq.StringArray
	Name        string
	Price       int
	Color       pq.StringArray
	Category    string
	Description string
	Rating      string
	Status      string
}

type home struct {
	Image  pq.StringArray
	Name   string
	Price  int
	Rating float32
}

func UserShowP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product models.Products
	var category models.Category
	var show details
	var s string
	var p []models.Products
	var r []models.Rating
	var ratingShow []rating
	var relatedShow []related

	Id, _ := strconv.Atoi(c.Param("Id"))

	err := database.Db.Where("id=?", uint(Id)).First(&product).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "Product not found!"})
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
			ratingShow = append(ratingShow, rating{k.Rating, k.Review})
		}
		var Avg string
		if product.AvrgRating != 0 {
			Avg = fmt.Sprint(product.AvrgRating, "/5")
		} else {
			Avg = "0 Rating"
		}
		show = details{
			Images:      product.ImageURLs,
			Name:        product.Name,
			Price:       product.Price,
			Color:       product.Color,
			Category:    category.Name,
			Description: product.Dscptn,
			Rating:      Avg,
			Status:      s,
		}
		database.Db.Where("Ctgry_Id=?", category.Id).Find(&p)
	} else {
		c.JSON(404, gin.H{"error": "Product not found!"})
		return
	}
	for i := 0; i < len(p); i++ {
		if p[i].ID != product.ID {
			relatedShow = append(relatedShow, related{p[i].ImageURLs[0], p[i].Name, p[i].Price})
		}
	}
	c.JSON(200, gin.H{
		"product":          show,
		"rating_review":    ratingShow,
		"related_products": relatedShow,
	})
	product = models.Products{}
	category = models.Category{}
}

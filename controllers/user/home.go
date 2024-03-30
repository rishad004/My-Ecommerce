package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func UserHome(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------HOME SHOWING------------------------")

	var product []models.Products
	var show []gin.H
	var rate float32

	database.Db.Find(&product)
	for i := 0; i < len(product); i++ {
		if product[i].CtgryBlock {
			var r []models.Rating
			database.Db.Find(&r, "Prdct_Id=?", product[i].ID)
			for _, k := range r {
				rate += k.Rating
			}
			if len(r) == 0 {
				rate = 0
			} else {
				rate = rate / float32(len(r))
			}
			show = append(show, gin.H{
				"Image":  product[i].ImageURLs[0:1],
				"Name":   product[i].Name,
				"Price":  product[i].Price,
				"Rating": rate,
			})
			rate = 0
		}
	}
	c.JSON(200, gin.H{"products": show})
}

func SortProduct(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SORTING PRODUCT------------------------")

	var p []models.Products
	var show []gin.H

	sortType := c.Query("type")
	fmt.Println(sortType)

	switch sortType {
	case "aA-zZ":
		database.Db.Order("name asc").Find(&p)
	case "zZ-aA":
		database.Db.Order("name desc").Find(&p)
	case "low to high":
		database.Db.Order("price asc").Find(&p)
	case "high to low":
		database.Db.Order("price desc").Find(&p)
	case "recent":
		database.Db.Order("id desc").Find(&p)
	case "most rated":
		database.Db.Order("Avrg_Rating desc").Find(&p)
	default:
		c.JSON(404, gin.H{"error": "Products not found"})
		return
	}
	for _, product := range p {
		var rate float32
		var r []models.Rating
		database.Db.Find(&r, "Prdct_Id=?", product.ID)
		for _, k := range r {
			rate += k.Rating
		}
		if len(r) == 0 {
			rate = 0
		} else {
			rate = rate / float32(len(r))
		}
		if product.CtgryBlock {
			show = append(show, gin.H{
				"Image":  product.ImageURLs[0:1],
				"Name":   product.Name,
				"Price":  product.Price,
				"Rating": rate,
			})
		}
	}
	c.JSON(200, gin.H{"products": show})
}

func UserSearchP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SEARCHING PRODUCT------------------------")

	var products []models.Products
	var show []gin.H

	searchQuery := c.Query("search")
	fmt.Println(searchQuery)

	database.Db.Where("name ILIKE ?", "%"+searchQuery+"%").Find(&products)
	if len(products) == 0 {
		c.JSON(404, gin.H{"message": "Products not found"})
		return
	}
	for _, product := range products {
		var rate float32
		var r []models.Rating
		database.Db.Find(&r, "Prdct_Id=?", product.ID)
		for _, k := range r {
			rate += k.Rating
		}
		if len(r) == 0 {
			rate = 0
		} else {
			rate = rate / float32(len(r))
		}
		if product.CtgryBlock {
			show = append(show, gin.H{
				"Image":  product.ImageURLs[0:1],
				"Name":   product.Name,
				"Price":  product.Price,
				"Rating": rate,
			})
		}
	}
	c.JSON(200, gin.H{"products": show})
}

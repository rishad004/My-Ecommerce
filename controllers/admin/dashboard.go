package controllers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
)

func Dashboard(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------ADMIN DASHBOARD----------------------")

	var products []models.Products
	var category []models.Category
	var show, Show []gin.H

	if err := database.Db.Order("Sold desc").Find(&products).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Message": "No products found!",
			"Data":    gin.H{},
		})
		return
	}

	if err := database.Db.Order("Sold desc").Find(&category).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Fail!",
			"Code":    404,
			"Message": "No category found!",
			"Data":    gin.H{},
		})
		return
	}

	for i := 0; i < len(category); i++ {
		if i < 10 {
			Show = append(Show, gin.H{
				"Name":        category[i].Name,
				"Description": category[i].Dscptn,
				"Sold":        products[i].Sold,
			})
		} else {
			break
		}
	}

	for i := 0; i < len(products); i++ {
		if i < 10 {
			show = append(show, gin.H{
				"Image":  products[i].ImageURLs,
				"Name":   products[i].Name,
				"Price":  products[i].Offer,
				"Rating": products[i].AvrgRating,
				"Sold":   products[i].Sold,
			})
		} else {
			break
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Products found!",
		"Data": gin.H{
			"TopCategories": Show,
			"TopProducts":   show,
		},
	})
}

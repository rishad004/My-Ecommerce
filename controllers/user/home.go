package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

func UserHome(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product []models.Products
	var show []home

	database.Db.Find(&product)
	for i := 0; i < len(product); i++ {
		if product[i].CtgryBlock {
			l := home{
				Image: product[i].ImageURLs[0:1],
				Name:  product[i].Name,
				Price: product[i].Price,
			}
			show = append(show, l)
		}
	}
	c.JSON(200, show)
}

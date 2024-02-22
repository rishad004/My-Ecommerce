package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

type details struct {
	Name        string
	Price       int
	Color       string
	Category    string
	Description string
	Status      string
}

func UserShowP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product []models.Products
	var category models.Category
	var show []details
	var s string

	database.Db.Find(&product)
	for i := 0; i < len(product); i++ {
		database.Db.First(&category).Update("Id", product[i].CtgryId)
		if category.Blocking {
			if product[i].Quantity != 0 {
				s = "On stock"
			} else {
				s = "Out of stock"
			}
			l := details{
				Name:        product[i].Name,
				Price:       product[i].Price,
				Color:       product[i].Color,
				Category:    category.Name,
				Description: product[i].Dscptn,
				Status:      s,
			}
			show = append(show, l)
		}
	}
	c.JSON(200, show)
}

func UserSearchP(c *gin.Context) {

}

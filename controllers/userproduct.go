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

type home struct {
	Name  string
	Price int
}

func UserHome(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product []models.Products
	var category models.Category
	var show []home

	database.Db.Find(&product)
	for i := 0; i < len(product); i++ {
		database.Db.First(&category).Update("Id", product[i].CtgryId)
		if category.Blocking {
			l := home{
				Name:  product[i].Name,
				Price: product[i].Price,
			}
			show = append(show, l)
			category = models.Category{}
		}
	}
	c.JSON(200, show)
}

func UserSearchP(c *gin.Context) {

}

func UserShowP(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PRODUCT SHOWING------------------------")

	var product models.Products
	var category models.Category
	var show details
	var s string

	name := c.Param("Name")

	database.Db.First(&product).Update("Name", name)

	database.Db.First(&category).Update("Id", product.CtgryId)
	if category.Blocking {
		if product.Quantity != 0 {
			s = "On stock"
		} else {
			s = "Out of stock"
		}
		show = details{
			Name:        product.Name,
			Price:       product.Price,
			Color:       product.Color,
			Category:    category.Name,
			Description: product.Dscptn,
			Status:      s,
		}

	}

	c.JSON(200, show)
}

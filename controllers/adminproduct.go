package controllers

import (
	"fmt"
	"project/database"
	"project/models"

	"github.com/gin-gonic/gin"
)

type pp struct {
	models.Products
	Ctgry string `json:"category"`
}

func AddProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT ADDING--------------------")

	var product pp
	var Product models.Products
	var check models.Category

	c.BindJSON(&product)

	database.Db.First(&check, "Name=?", product.Ctgry)

	Product = models.Products{
		Name:       product.Name,
		Price:      product.Price,
		Color:      product.Color,
		Quantity:   product.Quantity,
		Dscptn:     product.Dscptn,
		CtgryId:    check.Id,
		CtgryBlock: true,
	}

	if Product.CtgryId == 0 {
		c.JSON(404, "The category not found, Please add the category first.")
	} else {
		e := database.Db.Create(&Product)
		if e.Error != nil {
			c.JSON(409, "Product already exist, Please try to edit.")
		} else {
			c.JSON(200, "Product added successfully")
		}
	}
}

func EditProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT EDITING--------------------")

	name := c.Param("Name")

	var product pp
	var Product models.Products
	var p models.Products
	var check models.Category

	c.BindJSON(&product)

	database.Db.First(&check, "Name=?", product.Ctgry)

	Product = models.Products{
		Name:     product.Name,
		Price:    product.Price,
		Color:    product.Color,
		Quantity: product.Quantity,
		Dscptn:   product.Dscptn,
		CtgryId:  check.Id,
	}

	if Product.CtgryId == 0 {
		c.JSON(404, "The category not found, Please add the category first.")
	} else {
		database.Db.First(&p, "Name=?", name)
		database.Db.Model(&models.Products{}).Where("Name=?", name).Updates(Product)
		if p.Id == 0 {
			c.JSON(404, "Product not found")
		} else {
			c.JSON(200, "Product edited successfully")
		}
	}

}

func DeleteProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT DELETING--------------------")

	var product models.Products

	name := c.Param("Name")

	database.Db.First(&product, "Name=?", name)

	e := database.Db.Delete(&product)
	if e.Error != nil {
		c.JSON(422, "Couldn't delete the product, Please try again.")
	} else {
		c.JSON(200, "Product deleted successfully")
	}

}

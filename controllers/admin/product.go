package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type pp struct {
	models.Products
	Ctgry string `json:"category"`
}

func ShowProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT SHOWING--------------------")

	var product []models.Products
	var l []gin.H

	database.Db.Order("Id asc").Find(&product)

	for i := 0; i < len(product); i++ {
		var cat models.Category
		database.Db.First(&cat, product[i].CtgryId)
		l = append(l, gin.H{
			"ImageURLs":   product[i].ImageURLs,
			"Id":          product[i].ID,
			"Name":        product[i].Name,
			"Price":       product[i].Price,
			"Color":       product[i].Color,
			"Quantity":    product[i].Quantity,
			"Category":    cat.Name,
			"Description": product[i].Dscptn,
		})
	}
	c.JSON(200, gin.H{"Products": l})
}

func AddProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT ADDING--------------------")

	var product pp
	var Product models.Products
	var check models.Category

	file, _ := c.MultipartForm()
	product.Name = file.Value["name"][0]
	product.Price, _ = strconv.Atoi(file.Value["price"][0])
	product.Color = file.Value["color"]
	product.Quantity, _ = strconv.Atoi(file.Value["quantity"][0])
	product.Dscptn = file.Value["description"][0]
	product.Ctgry = file.Value["category"][0]
	image := file.File["image"]

	for _, k := range image {
		product.ImageURLs = append(product.ImageURLs, "./image/"+k.Filename)
		if err := c.SaveUploadedFile(k, "./assets/images/"+k.Filename); err != nil {
			c.JSON(400, gin.H{"Error": "Failed to save"})
		}
	}

	database.Db.First(&check, "Name=?", product.Ctgry)

	Product = models.Products{
		Name:       product.Name,
		Price:      product.Price,
		Color:      product.Color,
		Quantity:   product.Quantity,
		Dscptn:     product.Dscptn,
		CtgryId:    check.Id,
		CtgryBlock: true,
		ImageURLs:  product.ImageURLs,
	}

	if Product.CtgryId == 0 {
		c.JSON(404, gin.H{"Message": "The category not found, Please add the category first."})
	} else {
		e := database.Db.Create(&Product)
		if e.Error != nil {
			c.JSON(409, gin.H{"Message": "Product already exist, Please try to edit."})
		} else {
			c.JSON(200, gin.H{"Message": "Product added successfully"})
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
		c.JSON(404, gin.H{"Message": "The category not found, Please add the category first."})
	} else {
		database.Db.First(&p, "Name=?", name)
		database.Db.Model(&models.Products{}).Where("Name=?", name).Updates(Product)
		if p.ID == 0 {
			c.JSON(404, gin.H{"Error": "Product not found"})
		} else {
			c.JSON(200, gin.H{"Message": "Product edited successfully"})
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
		c.JSON(422, gin.H{"Error": "Couldn't delete the product, Please try again."})
	} else {
		c.JSON(200, gin.H{"Message": "Product deleted successfully"})
	}

}

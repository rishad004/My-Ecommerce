package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"

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

	file, _ := c.MultipartForm()
	product.Name = file.Value["name"][0]
	product.Price, _ = strconv.Atoi(file.Value["price"][0])
	product.Color = file.Value["color"][0]
	product.Quantity, _ = strconv.Atoi(file.Value["quantity"][0])
	product.Dscptn = file.Value["description"][0]
	product.Ctgry = file.Value["category"][0]
	image := file.File["image"]

	for _, k := range image {
		product.ImageURLs = append(product.ImageURLs, "./image/"+k.Filename)
		if err:=c.SaveUploadedFile(k, "./assets/images/"+k.Filename);err!=nil{
			c.JSON(400, "Failed to save")
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

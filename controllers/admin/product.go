package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowProduct godoc
// @Summary Products Show
// @Description Showing Products details in admin side
// @Tags Admin Product
// @Produce  json
// @Router /admin/product [get]
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
			"Offer":       product[i].Offer,
			"Color":       product[i].Color,
			"Quantity":    product[i].Quantity,
			"Category":    cat.Name,
			"Description": product[i].Dscptn,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved product details!",
		"Data": gin.H{
			"Products": l,
		},
	})
}

// AddProduct godoc
// @Summary Product Add
// @Description Adding Product with it's details
// @Tags Admin Product
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "Product Name"
// @Param price formData integer true "Product Price"
// @Param offer formData integer true "Product Offer"
// @Param color formData []string true "Product Color"
// @Param quantity formData integer true "Product Quantity"
// @Param description formData string true "Product Description"
// @Param category formData string true "Product Category"
// @Param image formData []file true "Product Image"
// @Router /admin/product [post]
func AddProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT ADDING--------------------")

	var Product models.Products
	var check models.Category

	file, _ := c.MultipartForm()
	Product.Name = c.Request.FormValue("name")
	Product.Price, _ = strconv.Atoi(c.Request.FormValue("price"))
	Product.Offer, _ = strconv.Atoi(c.Request.FormValue("offer"))
	Product.Color = file.Value["color"]
	Product.Quantity, _ = strconv.Atoi(c.Request.FormValue("quantity"))
	Product.Dscptn = c.Request.FormValue("description")
	Product.CtgryBlock = true
	Category := c.Request.FormValue("category")
	image := file.File["image"]

	if Product.Offer > Product.Price {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Offer price can't be more than the original price.!",
			"Data":    gin.H{},
		})
		return
	}

	for _, k := range image {
		Product.ImageURLs = append(Product.ImageURLs, "./assets/products/"+k.Filename)
		if err := c.SaveUploadedFile(k, "./assets/products/"+k.Filename); err != nil {
			c.JSON(400, gin.H{
				"Status":  "Error!",
				"Code":    400,
				"Error":   err.Error(),
				"Message": "Failed to save!",
				"Data":    gin.H{},
			})
		}
	}

	database.Db.First(&check, "Name=?", Category)

	Product.CtgryId = check.Id
	if Product.Offer == 0 {
		Product.Offer = Product.Price
	}

	if Product.CtgryId == 0 {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "The category not found, Please add the category first!",
			"Data":    gin.H{},
		})
	} else {
		e := database.Db.Create(&Product)
		if e.Error != nil {
			c.JSON(409, gin.H{
				"Status":  "Error!",
				"Code":    409,
				"Error":   e.Error.Error(),
				"Message": "Product already exist, Please try to edit!",
				"Data":    gin.H{},
			})
		} else {
			c.JSON(200, gin.H{
				"Status":  "Success!",
				"Code":    200,
				"Message": "Product added successfully!",
				"Data":    gin.H{},
			})
		}
	}
}

// EditProduct godoc
// @Summary Product Edit
// @Description Editing Product with it's details
// @Tags Admin Product
// @Accept  multipart/form-data
// @Produce  json
// @Param id query string true "name search by id"
// @Param name formData string true "Product Name"
// @Param price formData integer true "Product Price"
// @Param color formData []string true "Product Color"
// @Param quantity formData integer true "Product Quantity"
// @Param description formData string true "Product Description"
// @Param category formData string true "Product Category"
// @Param image formData []file true "Product Image"
// @Router /admin/product [put]
func EditProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT EDITING--------------------")

	name := c.Query("id")

	var Product models.Products
	var p models.Products
	var check models.Category

	file, _ := c.MultipartForm()
	Product.Name = c.Request.FormValue("name")
	Product.Price, _ = strconv.Atoi(c.Request.FormValue("price"))
	Product.Offer, _ = strconv.Atoi(c.Request.FormValue("offer"))
	Product.Color = file.Value["color"]
	Product.Quantity, _ = strconv.Atoi(c.Request.FormValue("quantity"))
	Product.Dscptn = c.Request.FormValue("description")
	Product.CtgryBlock = true
	Category := c.Request.FormValue("category")
	image := file.File["image"]

	if Product.Offer > Product.Price {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Offer price can't be more than the original price.!",
			"Data":    gin.H{},
		})
		return
	}
	if Product.Offer == 0 {
		Product.Offer = Product.Price
	}

	for _, k := range image {
		Product.ImageURLs = append(Product.ImageURLs, "./assets/products/"+k.Filename)
		if err := c.SaveUploadedFile(k, "./assets/products/"+k.Filename); err != nil {
			c.JSON(400, gin.H{
				"Status":  "Error!",
				"Code":    400,
				"Error":   err.Error(),
				"Message": "Failed to save!",
				"Data":    gin.H{},
			})
		}
	}

	if err := database.Db.First(&check, "Name=?", Category).Error; err != nil {
		c.JSON(404, gin.H{"Message": "The category not found, Please add the category first."})
	} else {
		Product.CtgryId = check.Id

		if err := database.Db.First(&p, "Id=?", name).Error; err != nil {
			c.JSON(404, gin.H{"Error": "Product not found"})
		} else {
			database.Db.Model(&models.Products{}).Where("Id=?", name).Updates(Product)
			c.JSON(200, gin.H{"Message": "Product edited successfully"})
		}
	}

}

// DeleteProduct godoc
// @Summary Product Delete
// @Description Deleting product completely
// @Tags Admin Product
// @Produce  json
// @Param id query string true "name search by id"
// @Router /admin/product [delete]
func DeleteProduct(c *gin.Context) {
	fmt.Println("")
	fmt.Println("---------------------PRODUCT DELETING--------------------")

	var product models.Products

	Id := c.Query("id")

	database.Db.First(&product, "ID=?", Id)

	e := database.Db.Delete(&product)
	if e.Error != nil {
		c.JSON(422, gin.H{
			"Status":  "Error!",
			"Code":    422,
			"Error":   e.Error.Error(),
			"Message": "Couldn't delete the product, Please try again!",
			"Data":    gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Product deleted successfully!",
			"Data":    gin.H{},
		})
	}

}

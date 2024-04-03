package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Scart struct {
	Product     string
	Color       string
	Quantity    uint
	Description string
	Price       int
}

func AddCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART ADDING------------------------")

	Logged := c.MustGet("Id").(uint)

	Id, _ := strconv.Atoi(c.Param("Id"))
	Color, _ := strconv.Atoi(c.Param("Color"))

	var product models.Products
	var color string
	var cc models.Cart

	err := database.Db.Where("ID=?", uint(Id)).First(&product).Error
	if err != nil {
		c.JSON(404, gin.H{"Error": "Product not found."})
	} else {
		eror := database.Db.First(&cc, "Product_Id=?", Id)

		if eror.Error == nil {
			if cc.Quantity < 10 && cc.Quantity < uint(product.Quantity) {
				cc.Quantity++
				database.Db.Save(&cc)
				c.JSON(200, gin.H{"Message": "Quantity increased successfully"})
			} else {
				c.JSON(409, gin.H{"Message": "This product can't be added to cart anymore"})
			}
		} else {
			if product.Quantity <= 0 {
				c.JSON(404, gin.H{"Error": "This product is out of stock!"})
				return
			}
			for i := 0; i < len(product.Color); i++ {
				if i == Color {
					color = product.Color[i]
				}
			}
			cart := models.Cart{
				UserId:    Logged,
				ProductId: uint(Id),
				Color:     color,
				Quantity:  1,
			}
			database.Db.Create(&cart)
			c.JSON(200, gin.H{"Message": "Product added successfully"})
		}
	}
}

func ShowCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART SHOWING------------------------")

	Logged := c.MustGet("Id").(uint)

	var cart []models.Cart
	var products []models.Products
	var l []gin.H
	var SubTotal int

	database.Db.Find(&cart, "User_Id=?", Logged)
	for i := 0; i < len(cart); i++ {
		var product models.Products
		database.Db.First(&product, cart[i].ProductId)
		products = append(products, product)
	}
	for i := 0; i < len(cart); i++ {
		l = append(l, gin.H{
			"Product":     products[i].Name,
			"Color":       cart[i].Color,
			"Quantity":    cart[i].Quantity,
			"Description": products[i].Dscptn,
			"Price":       products[i].Price,
		})
		SubTotal += int(cart[i].Quantity) * products[i].Price
	}
	c.JSON(200, gin.H{
		"Products": l,
		"SubTotal": SubTotal,
	})
}

func LessCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART LESSING------------------------")

	Logged := c.MustGet("Id").(uint)

	Id, _ := strconv.Atoi(c.Param("Id"))

	var cc models.Cart

	database.Db.First(&cc, "Product_Id=? AND User_Id=?", Id, Logged)
	if cc.ProductId == uint(Id) && cc.UserId == Logged {
		if cc.Quantity <= 1 {
			database.Db.Delete(&cc)
			c.JSON(200, gin.H{"Message": "Removed product from cart"})
		} else {
			cc.Quantity--
			database.Db.Save(&cc)
			c.JSON(200, gin.H{"Message": "Quantity decreased successfully"})
		}
	} else {
		c.JSON(404, gin.H{"Error": "Product not found in your cart"})
	}
}

func DeleteCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART REMOVING------------------------")

	var cc models.Cart

	Logged := c.MustGet("Id").(uint)
	Id, _ := strconv.Atoi(c.Param("Id"))

	database.Db.First(&cc, "Product_Id=? AND User_Id=?", Id, Logged)
	err := database.Db.Delete(&cc)

	if err.Error != nil {
		c.JSON(400, gin.H{"Error": "Couldn't delete data"})
		return
	}
	c.JSON(200, gin.H{"Message": "Product removed from cart."})
}

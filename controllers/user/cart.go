package controllers

import (
	"fmt"
	"project/database"
	"project/models"
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

	database.Db.First(&product, Id)
	if product.ID == 0 {
		c.JSON(404, gin.H{"error": "Product not found."})
	} else {
		eror := database.Db.First(&cc, "Product_Id=?", Id)

		if eror.Error == nil {
			if cc.Quantity < 10 && cc.Quantity < uint(product.Quantity) {
				cc.Quantity++
				database.Db.Save(&cc)
				c.JSON(200, gin.H{"message": "Quantity increased successfully"})
			} else {
				c.JSON(409, gin.H{"message": "This product can't be added to cart anymore"})
			}
		} else {
			if product.Quantity <= 0 {
				c.JSON(404, gin.H{"error": "This product is out of stock!"})
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
			c.JSON(200, gin.H{"message": "Product added successfully"})
		}
	}
}

func ShowCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART SHOWING------------------------")

	Logged := c.MustGet("Id").(uint)

	var cart []models.Cart
	var products []models.Products
	var show []Scart
	var SubTotal int

	database.Db.Find(&cart, "User_Id=?", Logged)
	for i := 0; i < len(cart); i++ {
		var product models.Products
		database.Db.First(&product, cart[i].ProductId)
		products = append(products, product)
	}
	for i := 0; i < len(cart); i++ {
		l := Scart{
			Product:     products[i].Name,
			Color:       cart[i].Color,
			Quantity:    cart[i].Quantity,
			Description: products[i].Dscptn,
			Price:       products[i].Price,
		}
		SubTotal += int(l.Quantity) * l.Price
		show = append(show, l)
	}
	c.JSON(200, gin.H{
		"Products": show,
		"SubTotal": SubTotal,
	})
}

func LessCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART LESSING------------------------")

	Logged := c.MustGet("Id").(uint)

	Id, _ := strconv.Atoi(c.Param("Id"))

	var cc models.Cart

	if Logged == 0 {
		c.JSON(401, "Please login first")
	} else {
		database.Db.First(&cc, "Product_Id=? AND User_Id=?", Id, Logged)
		if cc.ProductId == uint(Id) && cc.UserId == Logged {
			if cc.Quantity <= 1 {
				database.Db.Delete(&cc)
				c.JSON(200, gin.H{"message": "Removed product from cart"})
			} else {
				cc.Quantity--
				database.Db.Save(&cc)
				c.JSON(200, gin.H{"message": "Quantity decreased successfully"})
			}
		} else {
			c.JSON(404, gin.H{"error": "Product not found in your cart"})
		}
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
		c.JSON(400, gin.H{"error": "Couldn't delete data"})
		return
	}
	c.JSON(200, gin.H{"message": "Product removed from cart."})
}

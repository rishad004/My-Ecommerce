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
	Quantity    int
	Description string
	Price       int
}

func AddCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART ADDING------------------------")

	Id, _ := strconv.Atoi(c.Param("Id"))
	Color, _ := strconv.Atoi(c.Param("Color"))

	var product models.Products
	var color string
	var cc models.Cart

	if Logged == 0 {
		c.JSON(401, "Please login first")
	} else {
		database.Db.First(&product, Id)
		if product.Id == 0 {
			c.JSON(404, "Product not found.")
		} else {
			database.Db.First(&cc, "Productid=?", Id)

			if cc.Productid == uint(Id) && cc.UserId == Logged {
				cc.Quantity++
				database.Db.Save(&cc)
				c.JSON(200, "Product added successfully")
			} else {

				for i := 0; i < len(product.Color); i++ {
					if i == Color {
						color = product.Color[i]
					}
				}

				cart := models.Cart{
					UserId:    Logged,
					Productid: uint(Id),
					Color:     color,
					Quantity:  1,
				}
				database.Db.Create(&cart)
				c.JSON(200, "Product added successfully")
			}
		}
	}
}

func ShowCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CART SHOWING------------------------")

	var cart []models.Cart
	var products []models.Products
	var show []Scart
	var SubTotal int

	database.Db.Find(&cart, "User_Id=?", Logged)
	for i := 0; i < len(cart); i++ {
		var product models.Products
		database.Db.First(&product, cart[i].Productid)
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
		SubTotal += l.Quantity * l.Price
		show = append(show, l)
	}
	c.JSON(200, gin.H{
		"Products": show,
		"SubTotal": SubTotal,
	})
}

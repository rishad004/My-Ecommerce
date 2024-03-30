package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func ShowOrders(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SHOW ORDER------------------------")

	var order []models.Orderitem
	var show []gin.H

	err := database.Db.Preload("Order").Preload("Prdct").Preload("Order.User").Find(&order).Error

	if err != nil {
		c.JSON(404, gin.H{"message": "No orders found!"})
		return
	}
	for _, v := range order {
		show = append(show, gin.H{
			"Id":           v.Id,
			"OrderId":      v.OrderId,
			"Username":     v.Order.User.Name,
			"User_Email":   v.Order.User.Email,
			"Product_Name": v.Prdct.Name,
			"Image":        v.Prdct.ImageURLs[0],
			"Color":        v.Color,
			"Quantity":     v.Quantity,
			"Status":       v.Status,
		})
	}
	c.JSON(200, gin.H{"Orders": show})
}

func OrdersStatusChange(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ORDER STATUS CHANGING------------------------")

	ord, _ := strconv.Atoi(c.Query("order"))
	status := c.PostForm("status")

	var order models.Orderitem
	err := database.Db.Preload("Order").Preload("Prdct").First(&order, uint(ord)).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Order found!"})
		return
	}
	order.Status = status
	if status == "cancelled" {
		fmt.Println("Cancelling......................")
		order.Prdct.Quantity += order.Quantity
	}
	er := database.Db.Save(&order).Error
	if er != nil {
		c.JSON(401, gin.H{"error": "Couldn't cancel this order!"})
		return
	}
	if status == "cancelled" {
		c.JSON(200, gin.H{"message": "Order cancelled succesfully!"})
		return
	}
	c.JSON(200, gin.H{"message": "Order status updated successfully!"})
}

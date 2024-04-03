package controllers

import (
	"fmt"
	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"
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
		c.JSON(404, gin.H{"Message": "No orders found!"})
		return
	}
	for _, v := range order {
		var img string
		if v.Prdct.ImageURLs != nil {
			img = v.Prdct.ImageURLs[0]
		}
		show = append(show, gin.H{
			"Id":           v.Id,
			"OrderId":      v.OrderId,
			"Username":     v.Order.User.Name,
			"User_Email":   v.Order.User.Email,
			"Product_Name": v.Prdct.Name,
			"Image":        img,
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
	var payment models.Payment
	var wallet models.Wallet

	if err := database.Db.Preload("Order").Preload("Order.Coupon").Preload("Prdct").First(&order, uint(ord)).Error; err != nil {
		c.JSON(404, gin.H{"Error": "No such Order found!"})
		return
	}
	if err := database.Db.First(&payment, "Order_Id=?", order.OrderId).Error; err != nil {
		c.JSON(500, gin.H{"Error": "No such payment!"})
		return
	}
	if err := database.Db.First(&wallet, "User_Id=?", order.Order.UserId).Error; err != nil {
		c.JSON(501, gin.H{"Error": "Failed to find the user wallet!"})
		return
	}
	if status == "cancelled" {
		if order.Status == "cancelled" {
			c.JSON(409, gin.H{"Error": "This order is already cancelled"})
			return
		}

		order.Status = status

		fmt.Println("Cancelling......................")
		order.Order.SubTotal = order.Order.SubTotal - float32(order.Prdct.Price*order.Quantity)
		if order.Order.SubTotal < float32(order.Order.Coupon.Condition) {
			order.Order.Amount = order.Order.SubTotal
			order.Order.CouponId = 1
		} else {
			order.Order.Amount = order.Order.SubTotal - (order.Order.SubTotal * float32(order.Order.Coupon.Value) / 100)
		}
		if er := database.Db.Save(&order.Order).Error; er != nil {
			c.JSON(500, gin.H{"Error": "Can't decrease the order amount!"})
			return
		}
		order.Prdct.Quantity += order.Quantity
		if er := database.Db.Save(&order.Prdct).Error; er != nil {
			c.JSON(500, gin.H{"Error": "Can't increase product quantity!"})
			return
		}
		if payment.Status == "recieved" {
			wallet.Balance += (float32(order.Prdct.Price) * float32(order.Quantity)) - (float32(order.Prdct.Price) * float32(order.Quantity) * float32(order.Order.Coupon.Value) / 100)
			if err := database.Db.Save(&wallet).Error; err != nil {
				c.JSON(500, gin.H{"Error": "Couldn't update wallet!"})
				return
			}
			payment.Status = "partially refunded"
			if err := database.Db.Model(&payment).Update("Status", payment.Status).Error; err != nil {
				c.JSON(500, gin.H{"Error": "Failed to set payment as refunded"})
				return
			}
		}
	} else if status == "shipped" {
		order.Status = status
	} else if status == "delivered" {
		order.Status = status
	} else {
		c.JSON(400, gin.H{"Error": "This status can't be assigned!"})
		return
	}
	er := database.Db.Save(&order).Error
	if er != nil {
		c.JSON(401, gin.H{"Error": "Couldn't change the order status!"})
		return
	}
	if status == "cancelled" {
		c.JSON(200, gin.H{"Message": "Order cancelled succesfully!"})
		return
	}
	c.JSON(200, gin.H{"Message": "Order status updated successfully!"})
}

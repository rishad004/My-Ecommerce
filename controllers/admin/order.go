package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowOrder godoc
// @Summary Orders Show
// @Description Showing orders details in admin side
// @Tags Admin Order
// @Produce  json
// @Router /admin/order [get]
func ShowOrders(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SHOW ORDER------------------------")

	var order []models.Orderitem
	var show []gin.H

	err := database.Db.Preload("Order").Preload("Prdct").Preload("Order.User").Find(&order).Error

	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No orders found!!",
			"Data":    gin.H{},
		})
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
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved order details!",
		"Data": gin.H{
			"Orders": show,
		},
	})
}

// ChangerOrderStatus godoc
// @Summary Order status change
// @Description Changing order status on any order by admin
// @Tags Admin Order
// @Accept  multipart/form-data
// @Produce  json
// @Param status formData string true "Status of order"
// @Param order query string true "name search by order"
// @Router /admin/order [patch]
func OrdersStatusChange(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------ORDER STATUS CHANGING------------------------")

	status := c.Request.FormValue("status")

	ord, _ := strconv.Atoi(c.Query("order"))

	fmt.Println(status)

	var order models.Orderitem

	if err := database.Db.Preload("Order").Preload("Order.Coupon").Preload("Prdct").First(&order, uint(ord)).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "No such Order found!",
			"Data":    gin.H{},
		})
		return
	}
	if status == "shipped" {
		order.Status = status
	} else if status == "delivered" {
		order.Status = status
	} else {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Message": "This status can't be assigned!",
			"Data":    gin.H{},
		})
		return
	}
	er := database.Db.Save(&order).Error
	if er != nil {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   er.Error(),
			"Message": "Couldn't change the order status!",
			"Data":    gin.H{},
		})
		return
	}

	c.JSON(200, gin.H{
		"Status":  "Error!",
		"Code":    200,
		"Message": "Order status updated successfully!",
		"Data":    gin.H{},
	})
}

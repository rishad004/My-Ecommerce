package controllers

import (
	"fmt"
	"strconv"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/helper"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// CheckoutCart godoc
// @Summary Cart Chekout Only
// @Description Buying cart products
// @Tags User Order
// @Param address formData string true "address id"
// @Param coupon formData string true "coupon code"
// @Param method formData string true "payment method"
// @Produce  json
// @Router /user/cart/checkout [get]
func CheckoutCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CHECKOUT------------------------")

	var ca []models.Cart
	var coupon models.Coupons
	var order models.Orders
	var a models.Address
	var SubTotal int

	Logged := c.MustGet("Id").(uint)
	// Email := c.MustGet("Email").(string)
	Coupon := c.PostForm("coupon")
	method := c.PostForm("method")
	address, _ := strconv.Atoi(c.PostForm("address"))

	database.Db.First(&a, "Id=?", uint(address)).Where("User_Id = ? ", Logged)
	database.Db.Preload("Product").Find(&ca, "User_Id=?", Logged)

	if a.Id == 0 || a.User_Id != Logged {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No address found!",
			"Data":    gin.H{},
		})
		return
	}
	order.AddressId = uint(address)
	if len(ca) == 0 {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Your cart is empty!",
			"Data":    gin.H{},
		})
		return
	}
	errorr := database.Db.First(&coupon, "Code=?", Coupon)

	num := helper.GenerateInt()
	order.Ordernum, _ = strconv.Atoi(num)
	for _, v := range ca {
		SubTotal += int(v.Quantity) * v.Product.Price
	}

	order.SubTotal = float32(SubTotal)
	order.UserId = Logged
	order.Amount = order.SubTotal - (order.SubTotal * float32(coupon.Value) / 100)
	if coupon.Condition < SubTotal && errorr.Error == nil {
		order.CouponId = coupon.Id
	} else if coupon.Condition > SubTotal {
		c.JSON(401, gin.H{
			"Status":  "Fail!",
			"Code":    401,
			"Message": "Amount doesn't meet the min of coupon!",
			"Data":    gin.H{},
		})
		return
	} else if errorr.Error != nil {
		if Coupon == "" {
			order.CouponId = 1
		} else {
			c.JSON(404, gin.H{
				"Status":  "Error!",
				"Code":    404,
				"Message": "Not a valid coupon code!",
				"Data":    gin.H{},
			})
			return
		}
	}
	database.Db.Create(&order)
	for _, v := range ca {
		orderitem := models.Orderitem{
			OrderId:  order.Ordernum,
			PrdctId:  v.ProductId,
			Color:    v.Color,
			Quantity: int(v.Quantity),
			Status:   "Pending",
		}
		if er := database.Db.Create(&orderitem); er.Error != nil {
			c.JSON(403, gin.H{
				"Status":  "Error!",
				"Code":    403,
				"Message": "Couldn't place the order. Please try again later.",
				"Error":   er.Error,
				"Data":    gin.H{},
			})
			return
		}
		v.Product.Quantity -= int(v.Quantity)
		database.Db.Model(v.Product).Update("Quantity", v.Product.Quantity)
	}

	payment := models.Payment{
		OrderId: order.Ordernum,
		UserId:  Logged,
		Amount:  int(order.Amount),
		Status:  "pending",
	}
	if method == "COD" {

		payment.PMethod = "COD"

		for _, v := range ca {
			database.Db.Delete(&v)
		}

		if errr := database.Db.Create(&payment); errr.Error != nil {
			c.JSON(403, gin.H{
				"Status":  "Error!",
				"Code":    403,
				"Message": "Payment creation failed! Try again later!",
				"Error":   errr.Error,
				"Data":    gin.H{},
			})
			return
		}

		if err := Invoice(c, order.Ordernum); err != nil {
			c.JSON(500, gin.H{
				"Status":  "Error!",
				"Code":    500,
				"Message": "Error on invoice create!",
				"Error":   err.Error(),
				"Data":    gin.H{},
			})
			fmt.Println("Error on invoice create!", err.Error())
			return
		}

		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Order placed on COD!",
			"Data":    gin.H{},
		})
	} else if method == "PAY NOW" {

		payment.PMethod = "RAZOR PAY"

		razorId, err := helper.Executerazorpay(num, payment.Amount)

		if err != nil {
			c.JSON(406, gin.H{
				"Status":  "Error!",
				"Code":    406,
				"Message": "Payment gateway not initiated!",
				"Data":    gin.H{},
			})
			return
		}

		payment.PaymentId = razorId
		if errr := database.Db.Create(&payment); errr.Error != nil {
			c.JSON(403, gin.H{
				"Status":  "Error!",
				"Code":    403,
				"Message": "Payment creation failed! Try again later.!",
				"Data":    gin.H{},
			})
			return
		}

		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Complete the payment to place the order!",
			"Data": gin.H{
				"Payment": razorId,
			},
		})
	}
}

// CancelOrder godoc
// @Summary Cancelling Order
// @Description Cancelling ordered products individually
// @Tags User Order
// @Param order query string true "order id"
// @Produce  json
// @Router /user/order [patch]
func CancelOrder(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CANCEL ORDER------------------------")

	ord, _ := strconv.Atoi(c.Query("order"))
	Logged := c.MustGet("Id").(uint)

	var order models.Orderitem
	var wallet models.Wallet
	var payment models.Payment

	if err := database.Db.Preload("Order").Preload("Order.Coupon").Preload("Prdct").First(&order, uint(ord)).Error; err != nil || order.Order.UserId != Logged {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Order found!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}
	if order.Status == "cancelled" {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No such Order found!",
			"Data":    gin.H{},
		})
		return
	}

	if err := database.Db.First(&payment, "Order_Id=?", order.OrderId).Error; err != nil {
		c.JSON(500, gin.H{
			"Status":  "Error!",
			"Code":    500,
			"Message": "Couldn't fetch the payment!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}

	if order.Status == "delivered" {
		order.Status = "returned"
	} else {
		order.Status = "cancelled"
	}

	if order.Order.SubTotal == (float32(order.Prdct.Price) * float32(order.Quantity)) {
		payment.Status = "refunded"
		wall := wallet.Balance + order.Order.Amount
		database.Db.Model(&order.Order).Update("coupon_id", 1)
		if payment.Status == "recieved" {
			if err := database.Db.First(&wallet, "User_Id=?", Logged).Error; err != nil {
				c.JSON(404, gin.H{
					"Status":  "Error!",
					"Code":    404,
					"Message": "No such wallet found!",
					"Error":   err.Error(),
					"Data":    gin.H{},
				})
				return
			}
			database.Db.Model(&wallet).Update("balance", wall)
		}
		database.Db.Model(&order.Order).Update("sub_total", 0)
		database.Db.Model(&order.Order).Update("amount", 0)
		database.Db.Save(&order)

		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Order cancelled succesfully!",
			"Data":    gin.H{},
		})
		return
	}

	order.Order.SubTotal = order.Order.SubTotal - float32(order.Prdct.Price*order.Quantity)

	if order.Order.SubTotal < float32(order.Order.Coupon.Condition) {

		order.Order.Amount = order.Order.SubTotal
		order.Order.CouponId = 1

	} else {

		order.Order.Amount = order.Order.SubTotal - (order.Order.SubTotal * float32(order.Order.Coupon.Value) / 100)
	}

	order.Prdct.Quantity = order.Prdct.Quantity + order.Quantity
	database.Db.Model(&order.Prdct).Updates(&order.Prdct)
	database.Db.Model(&order.Order).Updates(&order.Order)
	if err := database.Db.Save(&order).Where("Order.User_Id AND Id", Logged, uint(ord)).Error; err != nil {
		c.JSON(401, gin.H{
			"Status":  "Error!",
			"Code":    401,
			"Message": "Couldn't cancel this order!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
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
			c.JSON(500, gin.H{
				"Status":  "Error!",
				"Code":    500,
				"Message": "Failed to set payment as refunded!",
				"Error":   err.Error(),
				"Data":    gin.H{},
			})
			return
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Order cancelled succesfully!",
		"Data":    gin.H{},
	})
}

// ShowOrders godoc
// @Summary Show Orders
// @Description Showing Orders with its details
// @Tags User Order
// @Produce  json
// @Router /user/order [get]
func ShowOrder(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SHOW ORDER------------------------")

	Logged := c.MustGet("Id").(uint)

	var orderitem []models.Orderitem
	var show []gin.H

	err := database.Db.Preload("Order").Preload("Prdct").Find(&orderitem).Where("Order.User_Id=?", Logged).Error
	if err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "No Orders found!",
			"Error":   err.Error(),
			"Data":    gin.H{},
		})
		return
	}
	for _, v := range orderitem {
		if v.Status != "cancelled" {
			if len(v.Prdct.ImageURLs) > 0 {
				show = append(show, gin.H{
					"Id":       v.Id,
					"OrderId":  v.OrderId,
					"Name":     v.Prdct.Name,
					"Image":    v.Prdct.ImageURLs[0],
					"Color":    v.Color,
					"Quantity": v.Quantity,
				})
			} else {
				show = append(show, gin.H{
					"Id":       v.Id,
					"OrderId":  v.OrderId,
					"Name":     v.Prdct.Name,
					"Image":    "",
					"Color":    v.Color,
					"Quantity": v.Quantity,
				})
			}
		}
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved orders and showing!",
		"Data": gin.H{
			"Orders": show,
		},
	})
}

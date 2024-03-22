package controllers

import (
	"fmt"
	"project/database"
	"project/helper"
	"project/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type OrderType struct {
	Id       uint
	OrderId  int
	Name     string
	Image    string
	Color    string
	Quantity int
}

func CheckoutCart(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CHECKOUT------------------------")

	var ca []models.Cart
	var coupon models.Coupons
	var order models.Orders
	var a models.Address
	var SubTotal int

	Logged := c.MustGet("Id").(uint)
	Email := c.MustGet("Email").(string)
	Coupon := c.PostForm("coupon")
	method := c.PostForm("method")
	address, _ := strconv.Atoi(c.PostForm("address"))

	errr := database.Db.First(&a, "Id=?", uint(address)).Where("User_Id = ? ", Logged).Error
	database.Db.Preload("Product").Find(&ca, "User_Id=?", Logged)

	if errr != nil {
		c.JSON(404, gin.H{"error": "No Address Found!"})
		return
	}
	order.AddressId = uint(address)
	if len(ca) == 0 {
		c.JSON(404, gin.H{"error": "Your cart is empty!"})
		return
	}
	errorr := database.Db.First(&coupon, "Code=?", Coupon)

	num := helper.GenerateInt()
	order.Ordernum, _ = strconv.Atoi(num)
	for _, v := range ca {
		SubTotal += int(v.Quantity) * v.Product.Price
	}

	order.SubTotal = SubTotal
	order.UserId = Logged
	order.Amount = order.SubTotal - (order.SubTotal * coupon.Value / 100)
	if coupon.Condition < SubTotal && errorr.Error == nil {
		order.CouponId = coupon.Id
	} else if coupon.Condition > SubTotal {
		c.JSON(401, gin.H{"message": "This coupon is not valid for this amount!"})
		return
	} else if errorr.Error != nil {
		if Coupon == "" {
			order.CouponId = 1
		} else {
			c.JSON(404, gin.H{"message": "Coupon code is not valid."})
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
			c.JSON(403, gin.H{"error": "Couldn't place the order the. Please try again later."})
			return
		}
		v.Product.Quantity -= int(v.Quantity)
		database.Db.Model(v.Product).Update("Quantity", v.Product.Quantity)
	}

	payment := models.Payment{
		OrderId: num,
		UserId:  Logged,
		Amount:  order.Amount,
		Status:  false,
	}
	if method == "COD" {
		payment.PMethod = "COD"
		for _, v := range ca {
			database.Db.Delete(&v)
		}
		body := fmt.Sprintln("\n*Order Id: ", num, "Amount: ", order.Amount, "Order will be delivered within 7 days.....*")
		helper.SendMail(c, Email, "Order placed succesfully", body)
		c.JSON(200, gin.H{"message": "Order placed on COD"})
	} else if method == "PAY NOW" {
		payment.PMethod = "RAZOR PAY"
		c.JSON(200, "Complete the payment to place the order.")
	}
	if errr := database.Db.Create(&payment); errr.Error != nil {
		c.JSON(403, gin.H{"error": "Payment  gateway failed! Try again later."})
	}

}

func CancelOrder(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------CANCEL ORDER------------------------")

	ord, _ := strconv.Atoi(c.Query("order"))
	Logged := c.MustGet("Id").(uint)

	var order models.Orderitem
	err := database.Db.Preload("Order").Preload("Prdct").First(&order, uint(ord)).Where("Order.User_Id=?", Logged).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "No such Order found!"})
		return
	}
	order.Status = "cancelled"
	order.Prdct.Quantity += order.Quantity
	er := database.Db.Save(&order).Error
	if er != nil {
		c.JSON(401, gin.H{"error": "Couldn't cancel this order!"})
		return
	}
	c.JSON(200, gin.H{"message": "Order cancelled succesfully!"})
}

func ShowOrder(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------SHOW ORDER------------------------")

	Logged := c.MustGet("Id").(uint)

	var orderitem []models.Orderitem
	var show []OrderType

	err := database.Db.Preload("Order").Preload("Prdct").Find(&orderitem).Where("Order.User_Id=?", Logged).Error
	if err != nil {
		c.JSON(404, gin.H{"message": "No orders found!"})
		return
	}
	for _, v := range orderitem {
		if v.Status != "cancelled" {
			if len(v.Prdct.ImageURLs) > 0 {
				show = append(show, OrderType{v.Id, v.OrderId, v.Prdct.Name, v.Prdct.ImageURLs[0], v.Color, v.Quantity})
			} else {
				show = append(show, OrderType{v.Id, v.OrderId, v.Prdct.Name, "", v.Color, v.Quantity})
			}
		}
	}
	c.JSON(200, gin.H{"Orders": show})
}

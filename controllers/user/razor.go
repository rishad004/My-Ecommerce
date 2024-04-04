package controllers

import (
	"fmt"
	"os"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/helper"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

type Razor struct {
	Order     string `json:"OrderID"`
	Payment   string `json:"PaymentID"`
	Signature string `json:"Signature"`
}

func RazorPay(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------RAZOR PAY------------------------")

	var payment models.Payment
	var orderitems []models.Orderitem
	var detail string

	orderId := c.Param("payment")
	Logged := c.MustGet("Id").(uint)

	if err := database.Db.Preload("User").First(&payment, "Payment_Id=?", orderId).Error; err != nil {
		c.JSON(404, gin.H{"Error": "Order not found!"})
		return
	}
	if payment.UserId != Logged {
		c.JSON(404, gin.H{"Error": "Order not found!"})
		return
	}
	if err := database.Db.Preload("Prdct").Find(&orderitems, "Order_Id=?", payment.OrderId).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Cannot fetch Order Items!"})
		return
	}
	for i := 0; i < len(orderitems); i++ {
		if i == (len(orderitems) - 1) {
			detail += "and " + orderitems[i].Prdct.Name
		} else {
			detail += orderitems[i].Prdct.Name + ", "
		}
	}

	c.HTML(200, "razor.html", gin.H{
		"Order":   orderId,
		"Amounr":  payment.Amount,
		"Key":     os.Getenv("RAZOR_KEY"),
		"Name":    payment.User.Name,
		"Eamil":   payment.User.Email,
		"Phone":   payment.User.Phone,
		"Product": "Your products " + detail + ". Pay for them now!",
	})
}

func RazorPayVerify(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------PAYMENT VERIFY------------------------")

	var verify Razor
	var order []models.Orderitem
	var payment models.Payment
	var ca []models.Cart

	Logged := c.MustGet("Id").(uint)

	err := c.ShouldBindJSON(&verify)
	if err != nil {
		c.JSON(404, gin.H{"Error": "Couldn't find any data bind!"})
		fmt.Println("Couldn't find any data bind!")
		return
	}

	er := database.Db.First(&payment, "Payment_Id=?", verify.Order).Error
	if er != nil {
		c.JSON(404, gin.H{"Error": "No such order found!"})
		fmt.Println("No such order found!")
		return
	}

	if err := database.Db.Preload("Order").Preload("Prdct").Find(&order, "Order_Id=?", payment.OrderId).Error; err != nil {
		c.JSON(500, gin.H{"Error": "Couldn't find order items from databse!"})
		fmt.Println("Couldn't find order items from databse!")
		return
	}
	eror := helper.RazorPaymentVerification(verify.Signature, verify.Order, verify.Payment)
	if eror != nil {
		c.JSON(402, gin.H{"Message": "Payment failed!"})
		fmt.Println("Payment failed!")
		return
	}
	payment.TransactionId = verify.Payment
	payment.Status = "recieved"
	erorr := database.Db.Save(&payment).Error

	errors := database.Db.Preload("Product").Find(&ca, "User_Id=?", Logged).Error
	if errors != nil {
		c.JSON(404, gin.H{"Error": "Nothing in cart!"})
		return
	}
	for _, v := range ca {
		v.Product.Quantity -= int(v.Quantity)
		erro := database.Db.Model(&v.Product).Updates(&v.Product).Error
		if erro != nil {
			c.JSON(402, gin.H{"Error": "Error while updating product"})
			return
		}
	}

	for _, v := range ca {
		database.Db.Delete(&v)
	}

	if erorr != nil {
		c.JSON(500, gin.H{"Error": "Couldn't update payment success in databse!"})
		fmt.Println("Couldn't update payment success in databse!")
		return
	}
	if err := Invoice(c, payment.OrderId); err != nil {
		c.JSON(500, gin.H{"Error": "Error on invoice create!"})
		fmt.Println("Error on invoice create!  ",err)
		return
	}
	c.JSON(200, gin.H{"Message": "Payment Succesfull!"})
	fmt.Println("Payment Succesfull!")
}

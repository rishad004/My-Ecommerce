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

// RazorPay godoc
// @Summary Razor Pay
// @Description Razor Pay payment gateway
// @Param id query string true "Payment id"
// @Tags User Payment
// @Router /user/payment [get]
func RazorPay(c *gin.Context) {

	fmt.Println("")
	fmt.Println("-----------------------------RAZOR PAY------------------------")

	var payment models.Payment
	var orderitems []models.Orderitem
	var detail string

	orderId := c.Query("id")
	Logged := c.MustGet("Id").(uint)

	if err := database.Db.Preload("User").First(&payment, "Payment_Id=?", orderId).Error; err != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Order not found!",
			"Data":    gin.H{},
		})
		return
	}
	if payment.UserId != Logged {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Order not found!",
			"Data":    gin.H{},
		})
		return
	}
	if err := database.Db.Preload("Prdct").Find(&orderitems, "Order_Id=?", payment.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Message": "Cannot fetch Order Items!",
			"Data":    gin.H{},
		})
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

// PaymentVerify godoc
// @Summary Verify Payment
// @Description Verifying Payment and updating payment status
// @Tags User Payment
// @Accept  json
// @Produce  json
// @Param verify body Razor true "Payment details"
// @Router /user/payment [post]
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
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   err.Error(),
			"Message": "Couldn't find any data bind!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't find any data bind!")
		return
	}

	er := database.Db.First(&payment, "Payment_Id=?", verify.Order).Error
	if er != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   er.Error(),
			"Message": "No such order found!",
			"Data":    gin.H{},
		})
		fmt.Println("No such order found!")
		return
	}

	if err := database.Db.Preload("Order").Preload("Prdct").Find(&order, "Order_Id=?", payment.OrderId).Error; err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Couldn't find order items from databse!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't find order items from databse!")
		return
	}
	eror := helper.RazorPaymentVerification(verify.Signature, verify.Order, verify.Payment)
	if eror != nil {
		c.JSON(402, gin.H{
			"Status":  "Error!",
			"Code":    402,
			"Error":   eror.Error(),
			"Message": "Payment failed!",
			"Data":    gin.H{},
		})
		fmt.Println("Payment failed!")
		return
	}
	payment.TransactionId = verify.Payment
	payment.Status = "recieved"
	erorr := database.Db.Save(&payment).Error

	errors := database.Db.Preload("Product").Find(&ca, "User_Id=?", Logged).Error
	if errors != nil {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Error":   errors.Error(),
			"Message": "Nothing in cart!",
			"Data":    gin.H{},
		})
		return
	}
	for _, v := range ca {
		v.Product.Quantity -= int(v.Quantity)
		erro := database.Db.Model(&v.Product).Updates(&v.Product).Error
		if erro != nil {
			c.JSON(402, gin.H{
				"Status":  "Error!",
				"Code":    402,
				"Error":   erro.Error(),
				"Message": "Error while updating product!",
				"Data":    gin.H{},
			})
			return
		}
	}

	for _, v := range ca {
		database.Db.Delete(&v)
	}

	if erorr != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   erorr.Error(),
			"Message": "Couldn't update payment success in databse!",
			"Data":    gin.H{},
		})
		fmt.Println("Couldn't update payment success in databse!")
		return
	}
	if err := Invoice(c, payment.OrderId); err != nil {
		c.JSON(400, gin.H{
			"Status":  "Error!",
			"Code":    400,
			"Error":   err.Error(),
			"Message": "Error on invoice create!",
			"Data":    gin.H{},
		})
		fmt.Println("Error on invoice create!  ", err)
		return
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    400,
		"Message": "Payment Succesfull!",
		"Data":    gin.H{},
	})
	fmt.Println("Payment Succesfull!")
}

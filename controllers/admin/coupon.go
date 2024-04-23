package controllers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rishad004/My-Ecommerce/database"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/gin-gonic/gin"
)

// ShowCoupon godoc
// @Summary Coupon Show
// @Description Showing Coupon details in admin side
// @Tags Admin Coupon
// @Produce  json
// @Router /admin/coupon [get]
func ShowCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON SHOWING----------------------")

	var coupon []models.Coupons
	var l []gin.H

	database.Db.Order("Id asc").Find(&coupon)

	for i := 0; i < len(coupon); i++ {
		diff := time.Until(coupon[i].Expr)
		days := int(diff.Hours() / 24)

		l = append(l, gin.H{
			"Id":           coupon[i].Id,
			"Name":         coupon[i].Name,
			"Description":  coupon[i].Dscptn,
			"Code":         coupon[i].Code,
			"Value":        coupon[i].Value,
			"Expires_Days": days,
		})
	}
	c.JSON(200, gin.H{
		"Status":  "Success!",
		"Code":    200,
		"Message": "Retrieved coupon details!",
		"Data":    l,
	})
}

// AddCoupon godoc
// @Summary Coupon Add
// @Description Adding Coupon with it's details
// @Tags Admin Coupon
// @Accept  multipart/form-data
// @Produce  json
// @Param name formData string true "coupon name"
// @Param description formData string true "coupon description"
// @Param code formData string true "coupon code"
// @Param condition formData string true "coupon condition"
// @Param value formData string true "coupon value"
// @Param expires formData string true "coupon expires"
// @Router /admin/coupon [post]
func AddCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON ADDING----------------------")

	var coupon models.Coupons

	coupon.Name = c.Request.FormValue("name")
	coupon.Dscptn = c.Request.FormValue("description")
	coupon.Code = c.Request.FormValue("code")
	coupon.Condition, _ = strconv.Atoi(c.Request.FormValue("condition"))
	coupon.Value, _ = strconv.Atoi(c.Request.FormValue("value"))
	exp, _ := strconv.Atoi(c.Request.FormValue("expires"))
	coupon.Expr = time.Now().AddDate(0, 0, exp)

	err := database.Db.Create(&coupon)
	if err.Error != nil {
		c.JSON(409, gin.H{
			"Status":  "Fail!",
			"Code":    409,
			"Message": "Coupon name or code already exist, please try to edit!",
			"Data":    gin.H{},
		})
	} else {
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Coupon added successfully!",
			"Data":    gin.H{},
		})
	}
}

// EditCoupon godoc
// @Summary Coupon Edit
// @Description Editing Coupon with it's details
// @Tags Admin Coupon
// @Accept  json
// @Produce  json
// @Param id query string true "name search by id"
// @Param name formData string true "coupon name"
// @Param description formData string true "coupon description"
// @Param code formData string true "coupon code"
// @Param condition formData string true "coupon condition"
// @Param value formData string true "coupon value""
// @Router /admin/coupon [put]
func EditCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON EDITING----------------------")

	Id, _ := strconv.Atoi(c.Query("id"))

	var cpp models.Coupons

	if err := database.Db.First(&cpp, Id).Error; err != nil {

		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Coupon not found!",
			"Data":    gin.H{},
		})
	} else {
		cpp.Name = c.Request.FormValue("name")
		cpp.Dscptn = c.Request.FormValue("description")
		cpp.Code = c.Request.FormValue("code")
		cpp.Condition, _ = strconv.Atoi(c.Request.FormValue("condition"))
		cpp.Value, _ = strconv.Atoi(c.Request.FormValue("value"))

		ERR := database.Db.Save(&cpp).Error
		if ERR != nil {
			c.JSON(400, gin.H{
				"Status":  "Error!",
				"Code":    404,
				"Message": "Coupon not updated!",
				"Error":   ERR.Error(),
				"Data":    gin.H{},
			})
			return
		}
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Coupon edited succesfully!",
			"Data":    gin.H{},
		})
	}

}

// DeleteCoupon godoc
// @Summary Coupon Delete
// @Description Deleting coupon completely
// @Tags Admin Coupon
// @Produce  json
// @Param id query string true "name search by id"
// @Router /admin/coupon [delete]
func DeleteCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON DELETING----------------------")

	Id, _ := strconv.Atoi(c.Query("id"))

	var coupon models.Coupons

	database.Db.First(&coupon, Id)

	if coupon.Id == 0 {
		c.JSON(404, gin.H{
			"Status":  "Error!",
			"Code":    404,
			"Message": "Coupon not found!",
			"Data":    gin.H{},
		})
	} else {
		database.Db.Delete(&coupon)
		c.JSON(200, gin.H{
			"Status":  "Success!",
			"Code":    200,
			"Message": "Coupon deleted successfully!",
			"Data":    gin.H{},
		})
	}
}

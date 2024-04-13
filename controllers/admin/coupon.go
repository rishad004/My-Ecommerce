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
	c.JSON(200, gin.H{"Coupons": l})
}

// AddCoupon godoc
// @Summary Coupon Add
// @Description Adding Coupon with it's details
// @Tags Admin Coupon
// @Accept  json
// @Produce  json
// @Param rcc body models.Coup true "Add Coupon"
// @Router /admin/coupon [post]
func AddCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON ADDING----------------------")

	var coupon models.Coupons
	var rcc models.Coup

	c.BindJSON(&rcc)
	coupon.Name = rcc.Name
	coupon.Dscptn = rcc.Description
	coupon.Code = rcc.Code
	coupon.Condition = rcc.Condition
	coupon.Value = rcc.Value
	coupon.Expr = time.Now().AddDate(0, 0, rcc.Duration)

	err := database.Db.Create(&coupon)
	if err.Error != nil {
		c.JSON(409, gin.H{"Message": "Coupon name or code already exist, please try to edit"})
	} else {
		c.JSON(200, gin.H{"Message": "Coupon added successfully"})
	}
}

// EditCoupon godoc
// @Summary Coupon Edit
// @Description Editing Coupon with it's details
// @Tags Admin Coupon
// @Accept  json
// @Produce  json
// @Param id query string false "name search by id"
// @Param rc body models.Coup true "Edit Coupon"
// @Router /admin/coupon [put]
func EditCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON EDITING----------------------")

	Id, _ := strconv.Atoi(c.Query("id"))

	var rc models.Coup
	var cpp models.Coupons

	c.BindJSON(&rc)

	database.Db.First(&cpp, Id)
	database.Db.Model(&models.Coupons{}).Where("Id=?", Id).Updates(rc)

	if cpp.Id == 0 {
		c.JSON(404, gin.H{"Error": "Coupon not found."})
	} else {
		c.JSON(200, gin.H{"Message": "Coupon edited succesfully."})
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
		c.JSON(404, gin.H{"Error": "Coupon not found."})
	} else {
		database.Db.Delete(&coupon)
		c.JSON(200, gin.H{"Message": "Coupon deleted successfully"})
	}
}

package controllers

import (
	"fmt"
	"project/database"
	"project/models"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type cp struct {
	models.Coupons
	Duration int `json:"day"`
}

type showcp struct {
	Id           uint
	Name         string
	Description  string
	Code         string
	Value        int
	Expires_Days int
}

func ShowCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON SHOWING----------------------")

	var coupon []models.Coupons
	var show []showcp

	database.Db.Order("Id asc").Find(&coupon)

	for i := 0; i < len(coupon); i++ {
		diff := time.Until(coupon[i].Expr)
		days := int(diff.Hours() / 24)

		l := showcp{
			Id:           coupon[i].Id,
			Name:         coupon[i].Name,
			Description:  coupon[i].Dscptn,
			Code:         coupon[i].Code,
			Value:        coupon[i].Value,
			Expires_Days: days,
		}
		show = append(show, l)
	}
	c.JSON(200, gin.H{"coupons": show})
}

func AddCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON ADDING----------------------")

	var coupon models.Coupons
	var rcc cp

	c.BindJSON(&rcc)
	coupon.Name = rcc.Name
	coupon.Dscptn = rcc.Dscptn
	coupon.Code = rcc.Code
	coupon.Value = rcc.Value
	coupon.Expr = time.Now().AddDate(0, 0, rcc.Duration)

	err := database.Db.Create(&coupon)
	if err.Error != nil {
		c.JSON(409, gin.H{"message": "Coupon name or code already exist, please try to edit"})
	} else {
		c.JSON(200, gin.H{"message": "Coupon added successfully"})
	}
}

func EditCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON EDITING----------------------")

	Id, _ := strconv.Atoi(c.Param("Id"))

	var rc models.Coupons
	var cpp models.Coupons

	c.BindJSON(&rc)

	database.Db.First(&cpp, Id)
	database.Db.Model(&models.Coupons{}).Where("Id=?", Id).Updates(rc)

	if cpp.Id == 0 {
		c.JSON(404, gin.H{"error": "Coupon not found."})
	} else {
		c.JSON(200, gin.H{"message": "Coupon edited succesfully."})
	}

}

func DeleteCoupon(c *gin.Context) {

	fmt.Println("")
	fmt.Println("---------------------------COUPON DELETING----------------------")

	Id, _ := strconv.Atoi(c.Param("Id"))

	var coupon models.Coupons

	database.Db.First(&coupon, Id)

	if coupon.Id == 0 {
		c.JSON(404, gin.H{"error": "Coupon not found."})
	} else {
		database.Db.Delete(&coupon)
		c.JSON(200, gin.H{"message": "Coupon deleted successfully"})
	}
}

package database

import (
	"fmt"
	"os"

	"github.com/rishad004/My-Ecommerce/helper"
	"github.com/rishad004/My-Ecommerce/models"

	"github.com/joho/godotenv"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func EnvLoad() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("------------------Error to Load ENV-------------------------")
	}
}

func DbConnect() {
	var err error
	var ad []models.Admin
	var coupon models.Coupons
	DSN := "host=" + os.Getenv("HOST") + " user=" + os.Getenv("USER_ID") + " password=" + os.Getenv("PASS") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("PORT") + " sslmode=disable TimeZone=Asia/Taipei"
	Db, err = gorm.Open(postgres.Open(DSN))
	if err != nil {
		fmt.Println("!!!!!!!!!!!!!!!!! Db connection failed !!!!!!!!!!!!!!!!!!")
	}
	Db.AutoMigrate(&models.Users{}, &models.Address{}, &models.Banner{}, &models.Cart{}, &models.Category{}, &models.Coupons{}, &models.Orders{}, &models.Otp{}, &models.Payment{}, &models.Products{}, &models.Wishlist{}, &models.Rating{}, &models.Admin{}, &models.Orderitem{}, &models.Wallet{}, &models.Referral{})

	Db.Find(&ad)

	for i := 0; i < len(ad); i++ {
		if len(ad[i].Pass) < 10 {
			ad[i].Pass = helper.HashPass(ad[i].Pass)
			Db.Save(&ad[i])
		}
	}
	coupon.Value = 0
	coupon.Name = "NO COUPON"
	coupon.Dscptn = "NO COUPON"
	coupon.Condition = 0
	coupon.Code = "NO COUPON"
	if ERROR := Db.Create(&coupon); ERROR.Error != nil {
		fmt.Println("Error this coupon is already exist")
	}
}

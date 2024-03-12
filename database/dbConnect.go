package database

import (
	"fmt"
	"os"
	"project/helper"
	"project/models"

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
	DSN := "host=" + os.Getenv("HOST") + " user=" + os.Getenv("USER") + " password=" + os.Getenv("PASS") + " dbname=" + os.Getenv("DB_NAME") + " port=" + os.Getenv("PORT") + " sslmode=disable TimeZone=Asia/Taipei"
	Db, err = gorm.Open(postgres.Open(DSN))
	if err != nil {
		fmt.Println("!!!!!!!!!!!!!!!!! Db connection failed !!!!!!!!!!!!!!!!!!")
	}
	Db.AutoMigrate(&models.Users{}, &models.Address{}, &models.Banner{}, &models.Cart{}, &models.Category{}, &models.Coupons{}, &models.Orders{}, &models.Otp{}, &models.Payment{}, &models.Products{}, &models.Wishlist{}, &models.Rating{}, &models.Admin{})

	Db.Find(&ad)

	for i := 0; i < len(ad); i++ {
		if len(ad[i].Pass) < 17 {
			ad[i].Pass = helper.HashPass(ad[i].Pass)
			Db.Save(&ad[i])
		}
	}
}

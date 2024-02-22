package database

import (
	"fmt"
	"os"
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
	Db, err = gorm.Open(postgres.Open(os.Getenv("DSN")))
	if err != nil {
		fmt.Println("!!!!!!!!!!!!!!!!! Db connection failed !!!!!!!!!!!!!!!!!!")
	}
	Db.AutoMigrate(&models.Users{}, &models.Address{}, &models.Banner{}, &models.Cart{}, &models.Category{}, &models.Coupons{}, &models.Orders{}, &models.Otp{}, &models.Payment{}, &models.Products{}, &models.Wishlist{})
}

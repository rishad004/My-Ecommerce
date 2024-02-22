package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	*gorm.Model
	Name       string `gorm:"not null" json:"username"`
	Email      string `gorm:"not null;unique" json:"usermail"`
	Pass       string `gorm:"not null" json:"userpass"`
	Phone      string `gorm:"not null;unique" json:"userphone"`
	Gender     string `gorm:"not null" json:"usergender"`
	Address_id uint
	Blocking   bool
	Admin      bool
}
type Otp struct {
	Id       uint   `gorm:"primaryKey"`
	Otp      string `gorm:"not null" json:"otp"`
	UserMail string
	Expr     time.Time
}
type Address struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null"`
	Phone    uint   `gorm:"not null"`
	PinCode  uint   `gorm:"not null"`
	City     string `gorm:"not null"`
	State    string `gorm:"not null"`
	Landmark string `gorm:"not null"`
	Address  string `gorm:"not null"`
}
type Wishlist struct {
	Id        uint `gorm:"primaryKey"`
	Productid uint `gorm:"not null"`
	UserId    uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
}
type Payment struct {
	Id            uint   `gorm:"primaryKey"`
	OrderId       uint   `gorm:"not null"`
	UserId        uint   `gorm:"not null"`
	Amount        int    `gorm:"not null"`
	Status        bool   `gorm:"not null"`
	PMethod       string `gorm:"not null"`
	TransactionId string
	Created_at    time.Time `gorm:"not null"`
}
type Cart struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint `gorm:"not null"`
	Productid uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	SubTotal  int  `gorm:"not null"`
}
type Orders struct {
	Id        uint   `gorm:"primaryKey"`
	UserId    uint   `gorm:"not null"`
	Productid uint   `gorm:"not null"`
	CouponId  uint   `gorm:"not null"`
	Amount    int    `gorm:"not null"`
	Status    string `gorm:"not null"`
}

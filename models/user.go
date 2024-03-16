package models

import (
	"time"

	"gorm.io/gorm"
)

type Users struct {
	*gorm.Model
	Name     string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null;unique" json:"usermail"`
	Pass     string `gorm:"not null" json:"userpass"`
	Phone    string `gorm:"not null;unique" json:"userphone"`
	Gender   string `gorm:"not null" json:"usergender"`
	Blocking bool
}
type Otp struct {
	Id       uint   `gorm:"primaryKey"`
	Otp      string `gorm:"not null" json:"otp"`
	UserMail string
	Expr     time.Time
}
type Address struct {
	Id       uint   `gorm:"primaryKey"`
	User_Id  uint   `gorm:"not null"`
	Name     string `gorm:"not null" json:"name"`
	Phone    uint   `gorm:"not null" json:"phone"`
	PinCode  uint   `gorm:"not null" json:"pincode"`
	City     string `gorm:"not null" json:"city"`
	State    string `gorm:"not null" json:"state"`
	Landmark string `gorm:"not null" json:"landmark"`
	Address  string `gorm:"not null" json:"address"`
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
	Id        uint   `gorm:"primaryKey"`
	UserId    uint   `gorm:"not null"`
	Productid uint   `gorm:"not null"`
	Color     string `gorm:"not null"`
	Quantity  uint   `gorm:"not null"`
}
type Orders struct {
	Id       uint `gorm:"primaryKey"`
	UserId   uint `gorm:"not null"`
	CouponId uint
	Amount   int    `gorm:"not null"`
	Status   string `gorm:"not null"`
}
type Rating struct {
	Id       uint    `gorm:"primaryKey"`
	Rating   float32 `json:"rating"`
	Review   string  `json:"review"`
	User_Id  uint
	Prdct_Id uint
}
type Orderitem struct {
	Id       uint `gorm:"primaryKey"`
	OrderId  uint
	PrdctId  uint
	Quantity int
}

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
	ProductId uint `gorm:"not null"`
	Product   Products
	UserId    uint `gorm:"not null"`
}
type Payment struct {
	*gorm.Model
	OrderId       int `gorm:"not null"`
	Order         Orders
	UserId        uint `gorm:"not null"`
	User          Users
	Amount        int    `gorm:"not null"`
	Status        string `gorm:"not null"`
	PMethod       string `gorm:"not null"`
	PaymentId     string
	TransactionId string
}
type Cart struct {
	Id        uint `gorm:"primaryKey"`
	UserId    uint `gorm:"not null"`
	ProductId uint `gorm:"not null"`
	Product   Products
	Color     string `gorm:"not null"`
	Quantity  uint   `gorm:"not null"`
}
type Orders struct {
	Ordernum  int  `gorm:"primaryKey"`
	UserId    uint `gorm:"not null"`
	User      Users
	SubTotal  float32
	AddressId uint
	Address   Address
	CouponId  uint
	Coupon    Coupons
	Amount    float32
	CreatedAt time.Time
	UpdatedAt time.Time
}
type Rating struct {
	Id      uint    `gorm:"primaryKey"`
	Rating  float32 `json:"rating"`
	Review  string  `json:"review"`
	User_Id uint
	PrdctId uint
	Prdct   Products
}
type Orderitem struct {
	Id       uint `gorm:"primaryKey"`
	OrderId  int
	PrdctId  uint
	Color    string
	Quantity int
	Status   string
	Order    Orders
	Prdct    Products
}
type Wallet struct {
	Id      uint `gorm:"primaryKey"`
	UserId  uint
	Balance float32
}
type Referral struct {
	Id     uint `gorm:"primaryKey"`
	UserId uint
	User   Users
	Code   string
	ById   uint
	By     Users
	Count  int
}

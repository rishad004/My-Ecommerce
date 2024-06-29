package models

import (
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type Admin struct {
	Id    uint `gorm:"primaryKey"`
	Name  string
	Email string `json:"adminmail"`
	Pass  string `json:"adminpass"`
}
type Products struct {
	*gorm.Model
	Name       string `gorm:"not null; unique" json:"name"`
	Price      int    `gorm:"not null" json:"price"`
	Offer      int
	Color      pq.StringArray `gorm:"not null;type:text[]" json:"color"`
	Quantity   int            `gorm:"not null" json:"quantity"`
	Dscptn     string         `gorm:"not null" json:"description"`
	AvrgRating float32
	ImageURLs  pq.StringArray `gorm:"type:text[]"`
	CtgryId    uint
	Sold       int
	Ctgry      Category
	CtgryBlock bool
}
type Category struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null; unique" json:"name"`
	Dscptn   string `gorm:"not null" json:"description"`
	Sold     int
	Blocking bool `gorm:"not null"`
}
type Coupons struct {
	Id        uint   `gorm:"primaryKey"`
	Name      string `gorm:"not null; unique" json:"name"`
	Dscptn    string `gorm:"not null" json:"description"`
	Code      string `gorm:"not null; unique" json:"code"`
	Condition int
	Value     int `gorm:"not null" json:"off"`
	Expr      time.Time
}
type Banner struct {
	Id         uint   `gorm:"primaryKey"`
	MainHeader string `gorm:"not null"`
	SubHeader  string `gorm:"not null"`
	Product_Id uint
	URL        string `gorm:"not null"`
}

package models

import "time"

type Products struct {
	Id         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null; unique" json:"name"`
	Price      int    `gorm:"not null" json:"price"`
	Color      string `gorm:"not null" json:"color"`
	Quantity   int    `gorm:"not null" json:"quantity"`
	Dscptn     string `gorm:"not null" json:"description"`
	CtgryId    uint
	CtgryBlock bool
}
type Category struct {
	Id       uint   `gorm:"primaryKey"`
	Name     string `gorm:"not null; unique" json:"name"`
	Dscptn   string `gorm:"not null" json:"description"`
	Blocking bool   `gorm:"not null"`
}
type Coupons struct {
	Id     uint      `gorm:"primaryKey"`
	Name   string    `gorm:"not null; unique"`
	Dscptn string    `gorm:"not null"`
	Code   string    `gorm:"not null"`
	Value  string    `gorm:"not null"`
	Expr   time.Time `gorm:"not null"`
}
type Banner struct {
	Id         uint   `gorm:"primaryKey"`
	MainHeader string `gorm:"not null"`
	SubHeader  string `gorm:"not null"`
	Product_Id uint
	URL        string `gorm:"not null"`
}

package models

import (
	"time"

	"github.com/lib/pq"
)

type Products struct {
	Id         uint           `gorm:"primaryKey"`
	Name       string         `gorm:"not null; unique" json:"name"`
	Price      int            `gorm:"not null" json:"price"`
	Color      pq.StringArray `gorm:"not null;type:text[]" json:"color"`
	Quantity   int            `gorm:"not null" json:"quantity"`
	Dscptn     string         `gorm:"not null" json:"description"`
	ImageURLs  pq.StringArray `gorm:"type:text[]"`
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
	Name   string    `gorm:"not null; unique" json:"name"`
	Dscptn string    `gorm:"not null" json:"description"`
	Code   string    `gorm:"not null; unique" json:"code"`
	Value  int       `gorm:"not null" json:"off"`
	Expr   time.Time `gorm:"not null"`
}
type Banner struct {
	Id         uint   `gorm:"primaryKey"`
	MainHeader string `gorm:"not null"`
	SubHeader  string `gorm:"not null"`
	Product_Id uint
	URL        string `gorm:"not null"`
}

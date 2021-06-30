package entity

import (
	"github.com/jackc/pgtype"
	"gorm.io/gorm"
)
type User struct {
	gorm.Model
	UserName string  `validate:"required,min=3,max=30"`
	UserEmail string `validate:"required,email,max=30" gorm:"unique"`
	UserPassword string `validate:"required,alphanum"`
}
type Category struct {
	gorm.Model
	Products []Product //`gorm:"foreignKey:category_id"`
	ParentID uint
	CategoryName string
	CategoryKeywords string
	CategoryDescription string
	CategorySlug string
	CategoryStatus bool
}

type Product struct {
	gorm.Model
	CategoryID uint
	ProductName string
	ProductKeywords string
	ProductDescription string
	ProductPrice float64
	ProductQuantity int
	ProductImage pgtype.Bytea
	ProductDiscount int // indirim tutarı
	ProductStatus bool
}
type Faq struct {
	gorm.Model
	FaqTitle string
	FaqDescription string
}

type ErrorResponse struct {
	FailedField string
	Field         interface{}
	Value       string
}
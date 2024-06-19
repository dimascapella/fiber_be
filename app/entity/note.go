package entity

import "gorm.io/gorm"

type Note struct {
	gorm.Model
	ProductID uint
	Note      string
	Qty       int
}

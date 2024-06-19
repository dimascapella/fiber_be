package entity

import "gorm.io/gorm"

type Product struct {
	gorm.Model
	Code          string `gorm:"size:255;unique;not null" validate:"required,max=10"`
	Nama          string `gorm:"size:255;not null" validate:"required,max=20"`
	Jumlah        int
	Deskripsi     string `gorm:"text"`
	Status_active bool
	Notes         []Note
}

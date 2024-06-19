package request

type ProductRequest struct {
	Code          string `form:"code" validate:"required,max=10"`
	Nama          string `form:"name" validate:"required,max=20"`
	Jumlah        int    `form:"qty" validate:"numeric"`
	Deskripsi     string `form:"description"`
	Status_active bool   `form:"status"`
}

type ProductUpdateRequest struct {
	Nama          string `form:"name" validate:"required,max=20"`
	Jumlah        int    `form:"qty" validate:"numeric"`
	Deskripsi     string `form:"description"`
	Status_active bool   `form:"status"`
}

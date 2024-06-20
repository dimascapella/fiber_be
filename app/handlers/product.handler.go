package handlers

import (
	"fiber_be/app/entity"
	"fiber_be/app/request"
	"fiber_be/database"
	"fiber_be/helper"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetProducts(ctx *fiber.Ctx) error {
	var products []entity.Product

	db := database.DB.Db

	db.Preload("Notes").Order("id desc").Find(&products)

	return ctx.JSON(fiber.Map{
		"status":   "success",
		"message":  "Products Found",
		"products": products,
	})
}

func GetProductById(ctx *fiber.Ctx) error {
	code := ctx.Params("code")

	db := database.DB.Db

	var check entity.Product
	result := db.Preload("Notes").Where("code = ?", code).First(&check)
	if result.Error != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Product not found",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Product Found",
		"product": check,
	})
}

func CreateProduct(ctx *fiber.Ctx) error {
	var productRequest request.ProductRequest
	db := database.DB.Db

	err := ctx.BodyParser(&productRequest)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	Validator := validator.New()
	validation := Validator.Struct(productRequest)
	if validation != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": validation.Error(),
		})
	}

	if productRequest.Jumlah <= 0 {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Jumlah Tidak Boleh 0",
		})
	}

	generateNumber := helper.RandomNumber(10)

	var check entity.Product
	db.Where("code = ?", generateNumber).First(&check)

	if check.Code == generateNumber {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Duplicate Product",
		})
	}

	product := entity.Product{
		Code:          generateNumber,
		Nama:          productRequest.Nama,
		Jumlah:        productRequest.Jumlah,
		Deskripsi:     productRequest.Deskripsi,
		Status_active: productRequest.Status_active,
	}

	db.Create(&product)

	note := entity.Note{
		ProductID: product.ID,
		Note:      "Masuk",
		Qty:       product.Jumlah,
	}

	db.Create(&note)

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Product Created",
		"product": productRequest,
	})
}

func UpdateProduct(ctx *fiber.Ctx) error {
	code := ctx.Params("code")
	var productRequest request.ProductUpdateRequest

	db := database.DB.Db

	err := ctx.BodyParser(&productRequest)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	Validator := validator.New()
	validation := Validator.Struct(productRequest)
	if validation != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": validation.Error(),
		})
	}

	var findProduct entity.Product
	result := db.Where("code = ?", code).First(&findProduct)
	if result.Error != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Product not found",
		})
	}

	var check entity.Product
	db.Where("nama = ?", productRequest.Nama).First(&check)
	if check.Nama == productRequest.Nama && check.Code != code {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Duplicate Product",
		})
	}

	result.Updates(map[string]interface{}{
		"Nama":          productRequest.Nama,
		"Jumlah":        productRequest.Jumlah,
		"Deskripsi":     productRequest.Deskripsi,
		"Status_active": productRequest.Status_active,
	})
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Product Updated",
		"product": findProduct,
	})
}

func DeleteProduct(ctx *fiber.Ctx) error {
	code := ctx.Params("code")

	db := database.DB.Db

	var findProduct entity.Product
	result := db.Where("code = ?", code).First(&findProduct)
	if result.Error != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Product not found",
		})
	}

	db.Unscoped().Where("product_id = ?", findProduct.ID).Delete(&entity.Note{})

	result.Unscoped().Delete(&findProduct)
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Product Deleted",
		"product": findProduct,
	})
}

func PostNote(ctx *fiber.Ctx) error {
	var NoteRequest request.NoteRequest

	db := database.DB.Db

	err := ctx.BodyParser(&NoteRequest)
	if err != nil {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": err.Error(),
		})
	}

	if NoteRequest.Qty <= 0 {
		return ctx.JSON(fiber.Map{
			"status":  "failed",
			"message": "Jumlah Tidak Boleh 0",
		})
	}

	var findProduct entity.Product
	result := db.Where("code = ?", NoteRequest.Code).First(&findProduct)

	if NoteRequest.Note == "Masuk" {
		result.Update("jumlah", findProduct.Jumlah+NoteRequest.Qty)
	} else {
		if NoteRequest.Qty > findProduct.Jumlah {
			return ctx.JSON(fiber.Map{
				"status":  "failed",
				"message": "Gagal Jumlah Barang Tidak Mencukupi",
			})
		}

		result.Update("jumlah", findProduct.Jumlah-NoteRequest.Qty)
	}

	note := entity.Note{
		ProductID: findProduct.ID,
		Note:      NoteRequest.Note,
		Qty:       NoteRequest.Qty,
	}

	db.Create(&note)

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Pencatatan Berhasil",
	})
}

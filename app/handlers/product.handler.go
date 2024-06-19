package handlers

import (
	"fiber_be/app/entity"
	"fiber_be/app/request"
	"fiber_be/database"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func GetProducts(ctx *fiber.Ctx) error {
	var products []entity.Product

	db := database.DB.Db

	db.Preload("Notes").Find(&products)

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message":  "Products Found",
		"products": products,
	})
}

func GetProductById(ctx *fiber.Ctx) error {
	code := ctx.Params("code")

	db := database.DB.Db

	var check entity.Product
	result := db.Where("code = ?", code).First(&check)
	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Product Found",
		"product": check,
	})
}

func CreateProduct(ctx *fiber.Ctx) error {
	var productRequest request.ProductRequest
	db := database.DB.Db

	err := ctx.BodyParser(&productRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	Validator := validator.New()
	validation := Validator.Struct(productRequest)
	if validation != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": validation.Error(),
		})
	}

	var check entity.Product
	db.Where("code = ?", productRequest.Code).First(&check)

	if check.Code == productRequest.Code {
		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Duplicate Product",
		})
	}

	product := entity.Product{
		Code:          productRequest.Code,
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

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	Validator := validator.New()
	validation := Validator.Struct(productRequest)
	if validation != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": validation.Error(),
		})
	}

	var findProduct entity.Product
	result := db.Where("code = ?", code).First(&findProduct)
	if result.Error != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	if findProduct.Nama == productRequest.Nama {
		return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
			"message": "Duplicate Product",
		})
	}

	updateProduct := entity.Product{
		Nama:          productRequest.Nama,
		Jumlah:        productRequest.Jumlah,
		Deskripsi:     productRequest.Deskripsi,
		Status_active: productRequest.Status_active,
	}

	result.Updates(updateProduct)
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Product not found",
		})
	}

	db.Where("product_id = ?", findProduct.ID).Delete(&entity.Note{})

	result.Delete(&findProduct)
	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Product Deleted",
		"product": findProduct,
	})
}

func PostNote(ctx *fiber.Ctx) error {
	var NoteRequest request.NoteRequest

	db := database.DB.Db

	err := ctx.BodyParser(&NoteRequest)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	var findProduct entity.Product
	result := db.Where("code = ?", NoteRequest.Code).First(&findProduct)

	if NoteRequest.Note == "Masuk" {
		result.Update("jumlah", findProduct.Jumlah+NoteRequest.Qty)
	} else {
		if NoteRequest.Qty > findProduct.Jumlah {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Order Over Product QTY",
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

	return ctx.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"message": "Product QTY Updated",
	})
}

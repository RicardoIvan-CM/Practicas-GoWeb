package handler

import (
	"log"
	"regexp"
	"strconv"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/product"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service product.Service
}

func ValidatePostProductRequest(req *CreateProductRequest) error {
	if req.Name == "" {
		return ErrProductNameRequired
	}
	if req.Quantity < 0 {
		return ErrProductQuantityInvalid
	}
	if req.CodeValue == "" {
		return ErrProductCodeValueRequired
	}
	if req.Expiration == "" {
		return ErrProductExpirationRequired
	}
	match, _ := regexp.MatchString("\\d{2}/\\d{2}/\\d{4}", req.Expiration)
	if !match {
		return ErrProductExpirationInvalid
	}
	if req.Price < 0 {
		return ErrProductPriceInvalid
	}
	return nil
}

func (handler ProductHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Obtener peticion y validarla
		var req CreateProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{
				"message": "Invalid request",
			})
			log.Println("Error :", err.Error())
			return
		}
		err := ValidatePostProductRequest(&req)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		newProduct := req.ToDomain()
		if err = handler.Service.Create(&newProduct); err != nil {
			ctx.JSON(400, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, newProduct)
	}
}

func (handler ProductHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := handler.Service.GetAll()
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, products)
	}
}

func (handler ProductHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "The product ID must be an integer",
			})
		}
		product, err := handler.Service.GetByID(id)
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, product)
	}
}

func (handler ProductHandler) GetBySearch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		priceGt, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(400, gin.H{
				"message": "The PriceGt parameter must be a number",
			})
		}
		product, err := handler.Service.GetByID(int(priceGt))
		if err != nil {
			ctx.JSON(500, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(200, product)
	}
}

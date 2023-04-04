package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/product"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/pkg/rest"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service product.Service
}

func ValidatePostProductRequest(req *ProductRequest) error {
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
		var req ProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			log.Println("Error :", err.Error())
			return
		}
		err := ValidatePostProductRequest(&req)
		if err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		newProduct := req.ToDomain()
		if err = handler.Service.Create(&newProduct); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		ctx.JSON(200, rest.SuccessfulResponse{
			Data: newProduct,
		})
	}
}

func (handler ProductHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := handler.Service.GetAll()
		if err != nil {
			ctx.JSON(500, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		ctx.JSON(200, rest.SuccessfulResponse{
			Data: products,
		})
	}
}

func (handler ProductHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
		}
		product, err := handler.Service.GetByID(id)
		if err != nil {
			ctx.JSON(500, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		ctx.JSON(200, rest.SuccessfulResponse{
			Data: product,
		})
	}
}

func (handler ProductHandler) GetBySearch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		priceGt, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
		}
		product, err := handler.Service.GetByID(int(priceGt))
		if err != nil {
			ctx.JSON(500, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		ctx.JSON(200, rest.SuccessfulResponse{
			Data: product,
		})
	}
}

func (handler ProductHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ProductRequest
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			log.Println("Error :", err.Error())
			return
		}

		if err := ValidatePostProductRequest(&req); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		newProduct := req.ToDomain()
		newProduct.ID = id

		if err := handler.Service.Update(&newProduct); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx.JSON(200, rest.SuccessfulResponse{
			Data: newProduct,
		})
	}
}

func (handler ProductHandler) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
		}
		product, err := handler.Service.GetByID(id)
		if err != nil {
			ctx.JSON(500, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		if err := json.NewDecoder(ctx.Request.Body).Decode(&product); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
		}
		product.ID = id

		if err := handler.Service.Update(product); err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}
		ctx.JSON(200, rest.SuccessfulResponse{
			Data: product,
		})
	}
}

func (handler ProductHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, rest.ErrorResponse{
				Error: err.Error(),
			})
		}

		if err := handler.Service.Delete(id); err != nil {
			ctx.JSON(500, rest.ErrorResponse{
				Error: err.Error(),
			})
			return
		}

		ctx.Header("Location", fmt.Sprintf("/products/%d", id))
		ctx.JSON(204, nil)
	}
}

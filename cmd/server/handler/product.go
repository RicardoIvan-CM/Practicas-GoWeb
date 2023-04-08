package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/domain"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/product"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/pkg/web"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	Service product.Service
}

func NewProductHandler(service product.Service) *ProductHandler {
	return &ProductHandler{
		Service: service,
	}
}

func ValidateProductRequest(req *ProductRequest) error {
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

var ErrInvalidToken = errors.New("The user token is not valid")

// @Summary Create Product
// @Tags Products
// @Description Create a Product
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Param body body ProductRequest true "Product"
// @Success 201 {object} web.SuccessfulResponse
// @Failure 400 {object} web.ErrorResponse
// @Router /products [post]
func (handler ProductHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Obtener peticion y validarla
		var req ProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, web.InvalidRequestResponse)
			log.Println("Error :", err.Error())
			return
		}
		err := ValidateProductRequest(&req)
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}
		newProduct := req.ToDomain()
		created, err := handler.Service.Create(newProduct)
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(201, web.SuccessfulResponse{
			Data: created,
		})
	}
}

// @Summary List Products
// @Tags Products
// @Description Get All Products
// @Accept json
// @Produce json
// @Param token header string true "token"
// @Success 200 {object} web.SuccessfulResponse
// @Router /products [get]
func (handler ProductHandler) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		products, err := handler.Service.GetAll()
		if err != nil {
			ctx.JSON(500, web.ErrorResponse{
				Status:  500,
				Code:    "InternalError",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, web.SuccessfulResponse{
			Data: products,
		})
	}
}

type ConsumerPriceData struct {
	Products   []domain.Product `json:"products"`
	TotalPrice float64          `json:"total_price"`
}

func (handler ProductHandler) GetConsumerPrice() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		listStr := ctx.Query("list")
		valuesStr := listStr[1 : len(listStr)-1]
		valuesStrArr := strings.Split(valuesStr, ",")

		var ids []int

		for _, idStr := range valuesStrArr {
			idStr = strings.TrimSpace(idStr)
			id, err := strconv.Atoi(idStr)
			if err != nil {
				ctx.JSON(400, web.ErrorResponse{
					Status:  400,
					Code:    "RequestError",
					Message: "An entered id is not valid",
				})
				return
			}
			ids = append(ids, id)
		}

		consumerPrice, products, err := handler.Service.GetConsumerPrice(ids)
		if err != nil {
			ctx.JSON(500, web.ErrorResponse{
				Status:  500,
				Code:    "InternalError",
				Message: err.Error(),
			})
			return
		}

		ctx.JSON(200, web.SuccessfulResponse{
			Data: ConsumerPriceData{
				products,
				consumerPrice,
			},
		})
	}
}

func (handler ProductHandler) GetByID() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: "The ID is not valid",
			})
			return
		}
		producto, err := handler.Service.GetByID(id)
		if err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.NotFoundResponse)
			} else {
				ctx.JSON(500, web.ErrorResponse{
					Status:  500,
					Code:    "InternalError",
					Message: err.Error(),
				})
			}
			return
		}
		ctx.JSON(200, web.SuccessfulResponse{
			Data: producto,
		})
	}
}

func (handler ProductHandler) GetBySearch() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		priceGt, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: "The value of the priceGt parameter is not valid",
			})
			return
		}
		product, err := handler.Service.GetByID(int(priceGt))
		if err != nil {
			ctx.JSON(500, web.ErrorResponse{
				Status:  500,
				Code:    "InternalError",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, web.SuccessfulResponse{
			Data: product,
		})
	}
}

func (handler ProductHandler) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req ProductRequest
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: "The ID is not valid",
			})
			return
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, web.InvalidRequestResponse)
			log.Println("Error :", err.Error())
			return
		}

		if err := ValidateProductRequest(&req); err != nil {
			ctx.JSON(400, web.InvalidRequestResponse)
			return
		}

		newProduct := req.ToDomain()
		newProduct.ID = id

		createdProduct, err := handler.Service.Update(newProduct)

		if err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.NotFoundResponse)
			} else {
				ctx.JSON(500, web.ErrorResponse{
					Status:  500,
					Code:    "InternalError",
					Message: err.Error(),
				})
			}
			return
		}

		ctx.JSON(200, web.SuccessfulResponse{
			Data: createdProduct,
		})
	}
}

func (handler ProductHandler) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: "The ID is not valid",
			})
			return
		}
		producto, err := handler.Service.GetByID(id)
		if err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.NotFoundResponse)
			} else {
				ctx.JSON(500, web.ErrorResponse{
					Status:  500,
					Code:    "InternalError",
					Message: err.Error(),
				})
			}
			return
		}
		if err := json.NewDecoder(ctx.Request.Body).Decode(&producto); err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: "The request is not valid",
			})
			return
		}
		producto.ID = id

		product, err := handler.Service.Update(producto)
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, web.SuccessfulResponse{
			Data: product,
		})
	}
}

func (handler ProductHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: "The ID is not valid",
			})
			return
		}

		if err := handler.Service.Delete(id); err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.NotFoundResponse)
			} else {
				ctx.JSON(500, web.ErrorResponse{
					Status:  500,
					Code:    "InternalError",
					Message: err.Error(),
				})
			}
			return
		}

		ctx.Header("Location", fmt.Sprintf("/products/%d", id))
		ctx.JSON(204, nil)
	}
}

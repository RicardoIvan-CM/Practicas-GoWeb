package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
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

func verifyToken(userToken string) error {
	token := os.Getenv("TOKEN")
	if userToken != token {
		return ErrInvalidToken
	}
	return nil
}

func (handler ProductHandler) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Verificacion de Token
		userToken := ctx.GetHeader("TOKEN")
		if err := verifyToken(userToken); err != nil {
			ctx.JSON(401, web.ErrorResponse{
				Status:  401,
				Code:    "InvalidTokenError",
				Message: err.Error(),
			})
			return
		}

		//Obtener peticion y validarla
		var req ProductRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
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
		if err = handler.Service.Create(&newProduct); err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(201, web.SuccessfulResponse{
			Data: newProduct,
		})
	}
}

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
					Message: err.Error(),
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
				Message: err.Error(),
			})
		}
		producto, err := handler.Service.GetByID(id)
		if err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.ErrorResponse{
					Status:  404,
					Code:    "NotFoundError",
					Message: err.Error(),
				})
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
				Message: err.Error(),
			})
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
		//Verificacion de Token
		userToken := ctx.GetHeader("TOKEN")
		if err := verifyToken(userToken); err != nil {
			ctx.JSON(401, web.ErrorResponse{
				Status:  401,
				Code:    "InvalidTokenError",
				Message: err.Error(),
			})
			return
		}

		var req ProductRequest
		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			log.Println("Error :", err.Error())
			return
		}

		if err := ValidateProductRequest(&req); err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}

		newProduct := req.ToDomain()
		newProduct.ID = id

		if err := handler.Service.Update(&newProduct); err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.ErrorResponse{
					Status:  404,
					Code:    "NotFoundError",
					Message: err.Error(),
				})
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
			Data: newProduct,
		})
	}
}

func (handler ProductHandler) UpdatePartial() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Verificacion de Token
		userToken := ctx.GetHeader("TOKEN")
		if err := verifyToken(userToken); err != nil {
			ctx.JSON(401, web.ErrorResponse{
				Status:  401,
				Code:    "InvalidTokenError",
				Message: err.Error(),
			})
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}
		producto, err := handler.Service.GetByID(id)
		if err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.ErrorResponse{
					Status:  404,
					Code:    "NotFoundError",
					Message: err.Error(),
				})
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
				Message: err.Error(),
			})
			return
		}
		producto.ID = id

		if err := handler.Service.Update(producto); err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}
		ctx.JSON(200, web.SuccessfulResponse{
			Data: producto,
		})
	}
}

func (handler ProductHandler) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		//Verificacion de Token
		userToken := ctx.GetHeader("TOKEN")
		if err := verifyToken(userToken); err != nil {
			ctx.JSON(401, web.ErrorResponse{
				Status:  401,
				Code:    "InvalidTokenError",
				Message: err.Error(),
			})
			return
		}

		id, err := strconv.Atoi(ctx.Param("id"))
		if err != nil {
			ctx.JSON(400, web.ErrorResponse{
				Status:  400,
				Code:    "RequestError",
				Message: err.Error(),
			})
			return
		}

		if err := handler.Service.Delete(id); err != nil {
			if errors.Is(err, product.ErrProductNotFound) {
				ctx.JSON(404, web.ErrorResponse{
					Status:  404,
					Code:    "NotFoundError",
					Message: err.Error(),
				})
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

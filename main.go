package main

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostProductRequest struct {
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var Products []Product
var lastProductID int

var (
	ErrMovieNameRequired       = errors.New("The name is required")
	ErrMovieQuantityInvalid    = errors.New("The quantity is not valid")
	ErrMovieCodeValueRequired  = errors.New("The codevalue is required")
	ErrMovieCodeValueExists    = errors.New("The codevalue already exists")
	ErrMovieExpirationRequired = errors.New("The expiration date is required")
	ErrMovieExpirationInvalid  = errors.New("The expiration date is not valid")
	ErrMoviePriceInvalid       = errors.New("The price is invalid")
)

func ValidatePostProductRequest(req *PostProductRequest) error {
	if req.Name == "" {
		return ErrMovieNameRequired
	}
	if req.Quantity < 0 {
		return ErrMovieQuantityInvalid
	}
	if req.CodeValue == "" {
		return ErrMovieCodeValueRequired
	}
	for _, product := range Products {
		if product.CodeValue == req.CodeValue {
			return ErrMovieCodeValueExists
		}
	}
	if req.Expiration == "" {
		return ErrMovieExpirationRequired
	}
	match, _ := regexp.MatchString("\\d{2}/\\d{2}/\\d{4}", req.Expiration)
	if !match {
		return ErrMovieExpirationInvalid
	}
	if req.Price < 0 {
		return ErrMoviePriceInvalid
	}
	return nil
}

func PostProduct(ctx *gin.Context) {
	var req PostProductRequest
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
	lastProductID++
	var newProduct = Product{
		ID:          lastProductID,
		Name:        req.Name,
		Quantity:    req.Quantity,
		CodeValue:   req.CodeValue,
		IsPublished: req.IsPublished,
		Expiration:  req.Expiration,
		Price:       req.Price,
	}
	Products = append(Products, newProduct)
	ctx.JSON(200, newProduct)
}

func GetAllProducts(ctx *gin.Context) {
	ctx.JSON(200, Products)
}

func GetProductByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "El parámetro id debe ser un número entero",
		})
		return
	}
	for _, product := range Products {
		if product.ID == id {
			ctx.JSON(200, product)
			return
		}
	}
	ctx.JSON(404, gin.H{
		"message": "No se encontró el producto",
	})
}

func GetSearchProduct(ctx *gin.Context) {
	priceGt, err := strconv.ParseFloat(ctx.Query("priceGt"), 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "El parámetro PriceGt debe ser un número",
		})
		return
	}
	foundProducts := []Product{}
	for _, product := range Products {
		if product.Price > priceGt {
			foundProducts = append(foundProducts, product)
		}
	}
	ctx.JSON(200, foundProducts)
}

func GetJSONProducts() {
	file, err := os.Open("products.json")
	if err != nil {
		panic("El archivo no pudo ser abierto")
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic("El archivo no pudo ser leido")
	}

	json.Unmarshal(bytes, &Products)
}

func main() {
	//GetJSONProducts()
	Products = []Product{}

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	gopher := router.Group("/products")
	gopher.GET("/", GetAllProducts)
	gopher.POST("/", PostProduct)
	gopher.GET("/:id", GetProductByID)
	gopher.GET("/search", GetSearchProduct)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

package main

import (
	"encoding/json"
	"io"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Product struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	CodeValue   string  `json:"code_value"`
	IsPublished bool    `json:"is_published"`
	Expiration  string  `json:"expiration"`
	Price       float64 `json:"price"`
}

var Products []Product = []Product{}

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

func main() {
	file, err := os.Open("products.json")
	if err != nil {
		panic("El archivo no pudo ser abierto")
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		panic("El archivo no pudo ser leido")
	}

	json.Unmarshal(bytes, &Products)
	router := gin.Default()

	gopher := router.Group("/products")
	gopher.GET("/", GetAllProducts)
	gopher.GET("/:id", GetProductByID)
	gopher.GET("/search", GetSearchProduct)

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
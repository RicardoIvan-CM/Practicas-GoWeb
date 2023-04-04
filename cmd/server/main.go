package main

import (
	"github.com/RicardoIvan-CM/Practicas-GoWeb/cmd/server/handler"
	"github.com/gin-gonic/gin"
)

/*
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
}*/

func main() {
	//GetJSONProducts()

	router := gin.Default()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.String(200, "pong")
	})

	var handler handler.ProductHandler

	gopher := router.Group("/products")
	{
		gopher.GET("/", handler.GetAll())
		gopher.POST("/", handler.Create())
		gopher.GET("/:id", handler.GetByID())
		gopher.GET("/search", handler.GetBySearch())
	}

	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}

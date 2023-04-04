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

	//Generar nuevo router en Gin
	server := gin.New()

	//condigurar el router

	router := handler.Router{
		Engine: server,
	}
	router.Setup()

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}

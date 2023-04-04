package main

import (
	"github.com/RicardoIvan-CM/Practicas-GoWeb/cmd/server/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
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

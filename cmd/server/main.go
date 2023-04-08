package main

import (
	"os"

	"github.com/RicardoIvan-CM/Practicas-GoWeb/cmd/server/docs"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/cmd/server/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)

// @title           MELI Bootcamp API
// @version         1.0
// @description     This API simulates handling of MELI Products.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Ricardo Cruz
// @contact.url    http://www.swagger.io/support
// @contact.email  ricardoivan.cruz@mercadolibre.com.mx

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		panic(err)
	}
	//Generar nuevo router en Gin
	server := gin.New()
	docs.SwaggerInfo.Host = os.Getenv(("HOST"))
	server.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//Configurar el router
	router := handler.Router{
		Engine: server,
	}
	router.Setup()

	if err := server.Run(":8080"); err != nil {
		panic(err)
	}
}

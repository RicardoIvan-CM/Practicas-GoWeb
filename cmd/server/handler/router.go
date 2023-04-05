package handler

import (
	"github.com/RicardoIvan-CM/Practicas-GoWeb/internal/product"
	"github.com/RicardoIvan-CM/Practicas-GoWeb/pkg/store"
	"github.com/gin-gonic/gin"
)

type Router struct {
	Engine *gin.Engine
}

func (router *Router) Setup() {
	//Set default middlewares
	router.Engine.Use(gin.Logger())
	router.Engine.Use(gin.Recovery())

	//Set routes
	router.SetProductRoutes()
}

func (router *Router) SetProductRoutes() {
	repository, err := store.NewJSONRepository("../../products.json")
	if err != nil {
		panic(err)
	}
	//repository := product.NewSliceRepository()
	service := product.NewDefaultService(repository)
	handler := NewProductHandler(service)

	group := router.Engine.Group("/products")
	{
		group.GET("/", handler.GetAll())
		group.POST("/", handler.Create())
		group.GET("/:id", handler.GetByID())
		group.GET("/search", handler.GetBySearch())
		group.PUT("/:id", handler.Update())
		group.PATCH("/:id", handler.UpdatePartial())
		group.DELETE("/:id", handler.Delete())
		group.GET("/consumer_price", handler.GetConsumerPrice())
	}
}

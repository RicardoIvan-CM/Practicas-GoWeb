package handler

import (
	"github.com/RicardoIvan-CM/Practicas-GoWeb/cmd/server/middleware"
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
	//router.Engine.Use(gin.Recovery())
	router.Engine.Use(middleware.Recovery())

	//Set routes
	router.SetProductRoutes()
}

func (router *Router) SetProductRoutes() {
	storage := store.NewJSONstore("../../products.json")
	repository := product.NewDefaultRepository(storage)
	service := product.NewDefaultService(repository)
	handler := NewProductHandler(service)

	group := router.Engine.Group("/products")
	{
		group.GET("/", handler.GetAll())
		group.GET("/search", handler.GetBySearch())
		group.GET("/consumer_price", handler.GetConsumerPrice())
		group.Use(middleware.ValidarToken())
		group.POST("/", handler.Create())
		group.GET("/:id", handler.GetByID())

		group.PUT("/:id", handler.Update())
		group.PATCH("/:id", handler.UpdatePartial())
		group.DELETE("/:id", handler.Delete())
	}

	router.Engine.GET("/panic", func(ctx *gin.Context) {
		panic("Nani!?")
	})
}

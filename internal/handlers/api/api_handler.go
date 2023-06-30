package api

import (
	"github.com/cost_control/internal/handlers/api/product"
	productRepos "github.com/cost_control/internal/repository/product"
	"github.com/cost_control/internal/service"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type ApiHandler struct {
	productHandler product.ProductApiHandler
}

func New(db *mongo.Collection) *ApiHandler {
	repos := productRepos.New(db)
	return &ApiHandler{productHandler: *product.New(service.New(repos))}
}

func (h *ApiHandler) InitRoutes() *gin.Engine {
	router := gin.New()

	productsEndpoints := router.Group("/products")
	{
		productsEndpoints.GET("/", h.productHandler.Get)
		productsEndpoints.GET("/report", h.productHandler.Report)
		productsEndpoints.POST("/", h.productHandler.Create)
	}
	productEndpoints := router.Group("/product")
	{
		productEndpoints.GET("/:id", h.productHandler.GetById)
		productEndpoints.PUT("/:id", h.productHandler.Update)
		productEndpoints.DELETE("/:id", h.productHandler.Delete)
	}

	return router
}

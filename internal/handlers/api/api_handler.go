package api

import (
	"github.com/cost_control/config"
	"github.com/cost_control/internal/handlers/api/auth"
	"github.com/cost_control/internal/handlers/api/middleware"
	"github.com/cost_control/internal/handlers/api/product"
	productRepos "github.com/cost_control/internal/repository/product"
	userRepos "github.com/cost_control/internal/repository/user"
	product2 "github.com/cost_control/internal/service/product"
	"github.com/cost_control/internal/service/user"
	"github.com/cost_control/pkg/jwt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	productCollection = "product"
	userCollection    = "user"
)

type ApiHandler struct {
	productHandler product.ProductApiHandler
	auth           auth.Auth
	jwtManager     *jwt.Token
}

func New(db *mongo.Database, config *config.Config) *ApiHandler {
	pRepo := productRepos.New(db.Collection(productCollection))
	uRepo := userRepos.New(db.Collection(userCollection))
	return &ApiHandler{
		productHandler: *product.New(product2.New(pRepo)),
		auth:           *auth.New(user.New(uRepo), config)}
}

func (h *ApiHandler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())
	authMiddleware := middleware.New()

	router.POST("/login", h.auth.Login)
	productsEndpoints := router.Group("/products", authMiddleware.AuthMiddleware)
	{
		productsEndpoints.GET("/", h.productHandler.Get)
		productsEndpoints.GET("/report", h.productHandler.Report)
		productsEndpoints.POST("/", h.productHandler.Create)
	}
	productEndpoints := router.Group("/product", authMiddleware.AuthMiddleware)
	{
		productEndpoints.GET("/:id", h.productHandler.GetById)
		productEndpoints.PUT("/:id", h.productHandler.Update)
		productEndpoints.DELETE("/:id", h.productHandler.Delete)
	}

	return router
}

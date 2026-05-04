package http

import (
	"github.com/gin-gonic/gin"
	v1 "ss-catalog-service/internal/delivery/http/v1"
	"ss-catalog-service/internal/domain"
)

// RouterConfig holds all pre-built usecases injected from main.go.
// The router knows NOTHING about the database or implementations.
type RouterConfig struct {
	ProductCommandUsecase domain.ProductCommandUsecase
	ProductQueryUsecase   domain.ProductQueryUsecase
	VariantCommandUsecase domain.VariantCommandUsecase
}

// SetupRouter wires the HTTP routes using already-constructed usecases.
// Dependency Injection is performed by the caller (main.go), not here.
func SetupRouter(r *gin.Engine, cfg RouterConfig) {
	productHandler := v1.NewProductHandler(cfg.ProductCommandUsecase, cfg.ProductQueryUsecase)

	api := r.Group("/api/v1")
	{
		products := api.Group("/products")
		{
			products.POST("", productHandler.CreateProduct)
			products.GET("", productHandler.GetProducts)
			products.GET("/:id", productHandler.GetProduct)
		}

		variantHandler := v1.NewVariantHandler(cfg.VariantCommandUsecase)
		variants := api.Group("/variants")
		{
			variants.POST("", variantHandler.CreateVariant)
		}
	}
}

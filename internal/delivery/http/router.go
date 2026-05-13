package http

import (
	"github.com/gin-gonic/gin"
	v1 "ss-catalog-service/internal/delivery/http/v1"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/internal/delivery/http/middleware"
	"ss-catalog-service/config"
)

// RouterConfig holds all pre-built usecases injected from main.go.
// The router knows NOTHING about the database or implementations.
type RouterConfig struct {
	ProductCommandUsecase domain.ProductCommandUsecase
	ProductQueryUsecase   domain.ProductQueryUsecase
	VariantCommandUsecase domain.VariantCommandUsecase
	InventoryCommandUsecase domain.InventoryCommandUsecase
	ReviewUsecase           domain.ReviewUsecase
	BundleUsecase           domain.BundleUsecase
	PriceHistoryRepository  domain.PriceHistoryRepository
	ImportUsecase           domain.ImportUsecase
	CategoryUsecase         domain.CategoryUsecase
	JWT                   config.JWTConfig
}

// SetupRouter wires the HTTP routes using already-constructed usecases.
// Dependency Injection is performed by the caller (main.go), not here.
func SetupRouter(r *gin.Engine, cfg RouterConfig) {
	productHandler := v1.NewProductHandler(cfg.ProductCommandUsecase, cfg.ProductQueryUsecase)
	variantHandler := v1.NewVariantHandler(cfg.VariantCommandUsecase)
	inventoryHandler := v1.NewInventoryHandler(cfg.InventoryCommandUsecase)
	sellerHandler := v1.NewSellerHandler()
	auditHandler := v1.NewAuditHandler()
	reviewHandler := v1.NewReviewHandler(cfg.ReviewUsecase, cfg.ProductQueryUsecase)
	bundleHandler := v1.NewBundleHandler(cfg.BundleUsecase)
	priceHandler := v1.NewPriceHandler(cfg.PriceHistoryRepository)
	importHandler := v1.NewImportHandler(cfg.ImportUsecase)
	categoryHandler := v1.NewCategoryHandler(cfg.CategoryUsecase)

	// Global Middlewares
	r.Use(middleware.CorrelationIDMiddleware())

	// Health Check for Gateway
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Catalog API Group
	api := r.Group("/api/catalog/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWT))
	{
		products := api.Group("/products")
		{
			products.POST("", middleware.RequireAuth(), productHandler.CreateProduct)
			products.PUT("/:id", middleware.RequireAuth(), productHandler.UpdateProduct)
			products.GET("", productHandler.GetProducts)
			// NOTE: /search must be registered BEFORE /:id to avoid static route collision in Gin
			products.GET("/search", productHandler.SearchProducts)
			products.GET("/faceted-search", productHandler.FacetedSearch)
			products.GET("/:id", productHandler.GetProduct)
		}

		variants := api.Group("/variants")
		{
			variants.POST("", variantHandler.CreateVariant)
		}

		inventory := api.Group("/inventory")
		{
			inventory.POST("/adjust", inventoryHandler.AdjustStock)
		}

		sellers := api.Group("/sellers")
		{
			sellers.POST("", sellerHandler.RegisterSeller)
			sellers.GET("/:code", sellerHandler.GetSeller)
		}

		audit := api.Group("/audit-logs")
		{
			audit.GET("", auditHandler.GetAuditLogs)
		}

		reviews := api.Group("/reviews")
		{
			reviews.POST("", middleware.RequireAuth(), reviewHandler.SubmitReview)
			reviews.GET("/product/:id", reviewHandler.GetProductReviews)
			reviews.GET("/product/:id/summary", reviewHandler.GetRatingSummary)
		}

		bundles := api.Group("/bundles")
		{
			bundles.POST("", middleware.RequireAuth(), bundleHandler.CreateBundle)
			bundles.GET("", bundleHandler.GetBundles)
			bundles.GET("/:id", bundleHandler.GetBundle)
		}

		api.GET("/products/:id/price-history", priceHandler.GetPriceHistory)

		imports := api.Group("/imports")
		{
			imports.POST("", middleware.RequireAuth(), importHandler.TriggerImport)
			imports.GET("/:id", importHandler.GetJobStatus)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
		}
	}
}

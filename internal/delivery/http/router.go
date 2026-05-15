package http

import (
	"github.com/gin-gonic/gin"
	v1 "ss-catalog-service/internal/delivery/http/v1"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/internal/delivery/http/middleware"
	"ss-catalog-service/config"
)

type AppUsecases struct {
	ProductCommand   domain.ProductCommandUsecase
	ProductQuery     domain.ProductQueryUsecase
	VariantCommand   domain.VariantCommandUsecase
	InventoryCommand domain.InventoryCommandUsecase
	InventoryQuery   domain.InventoryQueryUsecase
	Review           domain.ReviewUsecase
	Bundle           domain.BundleUsecase
	Import           domain.ImportUsecase
	Category         domain.CategoryUsecase
	Brand            domain.BrandUsecase
	Attribute        domain.AttributeUsecase
	Tag              domain.TagUsecase
	Warehouse        domain.WarehouseUsecase
	Digital          domain.DigitalUsecase
}

type AppRepositories struct {
	PriceHistory     domain.PriceHistoryRepository
	Seller           domain.SellerRepository
}

// RouterConfig holds all pre-built dependencies injected from main.go.
type RouterConfig struct {
	Usecases     AppUsecases
	Repositories AppRepositories
	JWT          config.JWTConfig
}

// SetupRouter wires the HTTP routes using already-constructed usecases.
// Dependency Injection is performed by the caller (main.go), not here.
func SetupRouter(r *gin.Engine, cfg RouterConfig) {
	productHandler := v1.NewProductHandler(cfg.Usecases.ProductCommand, cfg.Usecases.ProductQuery)
	variantHandler := v1.NewVariantHandler(cfg.Usecases.VariantCommand)
	inventoryHandler := v1.NewInventoryHandler(cfg.Usecases.InventoryCommand, cfg.Usecases.InventoryQuery)
	sellerHandler := v1.NewSellerHandler()
	auditHandler := v1.NewAuditHandler()
	reviewHandler := v1.NewReviewHandler(cfg.Usecases.Review, cfg.Usecases.ProductQuery)
	bundleHandler := v1.NewBundleHandler(cfg.Usecases.Bundle)
	priceHandler := v1.NewPriceHandler(cfg.Repositories.PriceHistory)
	importHandler := v1.NewImportHandler(cfg.Usecases.Import)
	categoryHandler := v1.NewCategoryHandler(cfg.Usecases.Category)
	brandHandler := v1.NewBrandHandler(cfg.Usecases.Brand)
	attributeHandler := v1.NewAttributeHandler(cfg.Usecases.Attribute)
	tagHandler := v1.NewTagHandler(cfg.Usecases.Tag)
	warehouseHandler := v1.NewWarehouseHandler(cfg.Usecases.Warehouse)
	digitalHandler := v1.NewDigitalHandler(cfg.Usecases.Digital)

	// Global Middlewares
	r.Use(middleware.CorrelationIDMiddleware())

	// Health Check for Gateway
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "UP"})
	})

	// Catalog API Group
	api := r.Group("/api/catalog/v1")
	api.Use(middleware.AuthMiddleware(cfg.JWT, cfg.Repositories.Seller))
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
			inventory.GET("", inventoryHandler.GetInventory)
			inventory.POST("/adjust", inventoryHandler.AdjustStock)
		}

		sellers := api.Group("/sellers")
		{
			sellers.POST("", sellerHandler.RegisterSeller)
			sellers.GET("/:code", sellerHandler.GetSeller)
		}

		brands := api.Group("/brands")
		{
			brands.GET("", brandHandler.GetBrands)
			brands.GET("/:id", brandHandler.GetBrand)
			brands.POST("", middleware.RequireAuth(), brandHandler.CreateBrand)
			brands.PUT("/:id", middleware.RequireAuth(), brandHandler.UpdateBrand)
			brands.DELETE("/:id", middleware.RequireAuth(), brandHandler.DeleteBrand)
		}

		attributes := api.Group("/attributes")
		{
			attributes.GET("", attributeHandler.GetAttributes)
			attributes.POST("", middleware.RequireAuth(), attributeHandler.CreateAttribute)
			attributes.DELETE("/:id", middleware.RequireAuth(), attributeHandler.DeleteAttribute)
		}

		tags := api.Group("/tags")
		{
			tags.GET("", tagHandler.GetTags)
			tags.POST("", middleware.RequireAuth(), tagHandler.CreateTag)
			tags.DELETE("/:id", middleware.RequireAuth(), tagHandler.DeleteTag)
		}

		warehouses := api.Group("/warehouses")
		{
			warehouses.GET("", warehouseHandler.GetWarehouses)
			warehouses.POST("", middleware.RequireAuth(), warehouseHandler.CreateWarehouse)
			warehouses.PUT("/:id", middleware.RequireAuth(), warehouseHandler.UpdateWarehouse)
			warehouses.DELETE("/:id", middleware.RequireAuth(), warehouseHandler.DeleteWarehouse)
		}

		audit := api.Group("/audit-logs")
		{
			audit.GET("", auditHandler.GetAuditLogs)
		}

		reviews := api.Group("/reviews")
		{
			reviews.GET("", middleware.RequireAuth(), reviewHandler.GetAllReviews)
			reviews.POST("", middleware.RequireAuth(), reviewHandler.SubmitReview)
			reviews.GET("/product/:id", reviewHandler.GetProductReviews)
			reviews.GET("/product/:id/summary", reviewHandler.GetRatingSummary)
			reviews.PATCH("/:id/status", middleware.RequireAuth(), reviewHandler.UpdateReviewStatus)
		}

		bundles := api.Group("/bundles")
		{
			bundles.POST("", middleware.RequireAuth(), bundleHandler.CreateBundle)
			bundles.GET("", bundleHandler.GetBundles)
			bundles.GET("/:id", bundleHandler.GetBundle)
			bundles.PUT("/:id", middleware.RequireAuth(), bundleHandler.UpdateBundle)
			bundles.DELETE("/:id", middleware.RequireAuth(), bundleHandler.DeleteBundle)
		}

		api.GET("/products/:id/price-history", priceHandler.GetPriceHistory)

		imports := api.Group("/imports")
		{
			imports.GET("", middleware.RequireAuth(), importHandler.GetImportJobs)
			imports.POST("", middleware.RequireAuth(), importHandler.TriggerImport)
			imports.GET("/:id", importHandler.GetJobStatus)
		}

		digital := api.Group("/digital")
		{
			digital.POST("/licenses", middleware.RequireAuth(), digitalHandler.AddLicenses)
			digital.GET("/product/:product_id", digitalHandler.GetDigitalDetails)
			digital.POST("/upload", middleware.RequireAuth(), digitalHandler.UploadFile)
		}

		categories := api.Group("/categories")
		{
			categories.GET("", categoryHandler.GetCategories)
			categories.GET("/:id", categoryHandler.GetCategory)
			categories.POST("", middleware.RequireAuth(), categoryHandler.CreateCategory)
			categories.PUT("/:id", middleware.RequireAuth(), categoryHandler.UpdateCategory)
			categories.DELETE("/:id", middleware.RequireAuth(), categoryHandler.DeleteCategory)
		}
	}
}

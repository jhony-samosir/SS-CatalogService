package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"ss-catalog-service/config"
	"ss-catalog-service/internal/domain"
	apphttp "ss-catalog-service/internal/delivery/http"
	"ss-catalog-service/internal/infrastructure/cache"
	"ss-catalog-service/internal/infrastructure/database"
	"ss-catalog-service/internal/infrastructure/messaging"
	pgmodel "ss-catalog-service/internal/repository/postgres"
	inventoryusecase "ss-catalog-service/internal/usecase/inventory"
	productusecase "ss-catalog-service/internal/usecase/product"
	variantusecase "ss-catalog-service/internal/usecase/variant"
	reviewusecase "ss-catalog-service/internal/usecase/review"
	bundleusecase "ss-catalog-service/internal/usecase/bundle"
	importusecase "ss-catalog-service/internal/usecase/import_job"
	categoryusecase "ss-catalog-service/internal/usecase/category"
	msrepo "ss-catalog-service/internal/repository/meilisearch"
	"ss-catalog-service/internal/worker"
)

func main() {
	// --- Load Centralized Config ---
	cfg := config.Load()

	// --- Infrastructure: Database ---
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}

	// Run SQL Migrations (Source of Truth)
	if err := database.RunMigrations(db, "db/migrations"); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// --- Infrastructure: Messaging ---
	logBroker := messaging.NewLogBroker()

	// --- Dependency Injection (Composition Root) ---
	txManager := pgmodel.NewTransactionManager(db)
	outboxRepo := pgmodel.NewOutboxRepository(db)

	productRepo := pgmodel.NewProductRepository(db)
	
	// Define active languages for cache invalidation (Opsi A)
	activeLangs := []string{"id-ID", "en-US"}
	
	baseCacheRepo, err := cache.NewProductCacheRepository(10*time.Minute, activeLangs)
	if err != nil {
		log.Fatalf("Cache initialization failed: %v", err)
	}
	
	// Wrap with Prometheus metrics decorator
	productCacheRepo := cache.NewProductCacheMetricsDecorator(baseCacheRepo)

	meiliURL := os.Getenv("MEILI_URL")
	meiliKey := os.Getenv("MEILI_MASTER_KEY")
	var searchRepo domain.SearchRepository
	if meiliURL != "" {
		var err error
		searchRepo, err = msrepo.NewProductSearchRepository(meiliURL, meiliKey)
		if err != nil {
			log.Printf("Failed to connect to Meilisearch: %v", err)
		}
	}

	productCmd := productusecase.NewProductCommandUsecase(productRepo, productCacheRepo, outboxRepo, txManager)
	productQry := productusecase.NewProductQueryUsecase(productRepo, searchRepo, productCacheRepo, cfg.App.DefaultLang)

	variantRepo := pgmodel.NewVariantRepository(db)
	variantCmd := variantusecase.NewVariantCommandUsecase(variantRepo, productRepo, txManager)

	inventoryRepo := pgmodel.NewInventoryRepository(db)
	inventoryCmd := inventoryusecase.NewInventoryCommandUsecase(inventoryRepo, txManager)

	reviewRepo := pgmodel.NewReviewRepository(db)
	reviewUsecase := reviewusecase.NewReviewUsecase(reviewRepo)

	bundleRepo := pgmodel.NewBundleRepository(db)
	bundleUsecase := bundleusecase.NewBundleUsecase(bundleRepo)

	priceRepo := pgmodel.NewPriceHistoryRepository(db)

	importRepo := pgmodel.NewImportRepository(db)
	importUsecase := importusecase.NewImportUsecase(importRepo)

	categoryRepo := pgmodel.NewCategoryRepository(db)
	categoryUsecase := categoryusecase.NewCategoryUsecase(categoryRepo)

	sellerRepo := pgmodel.NewSellerRepository(db)

	// --- Background Workers ---
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	outboxWorker := worker.NewOutboxWorker(outboxRepo, logBroker, 5*time.Second)
	go outboxWorker.Start(ctx)

	// --- HTTP Router ---
	r := gin.Default()
	apphttp.SetupRouter(r, apphttp.RouterConfig{
		ProductCommandUsecase: productCmd,
		ProductQueryUsecase:   productQry,
		VariantCommandUsecase: variantCmd,
		InventoryCommandUsecase: inventoryCmd,
		ReviewUsecase:           reviewUsecase,
		BundleUsecase:           bundleUsecase,
		PriceHistoryRepository:  priceRepo,
		ImportUsecase:           importUsecase,
		CategoryUsecase:         categoryUsecase,
		SellerRepository:        sellerRepo,
		JWT:                   cfg.JWT,
	})

	// --- Start Server ---
	port := cfg.App.Port

	log.Printf("🚀 Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

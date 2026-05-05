package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	apphttp "ss-catalog-service/internal/delivery/http"
	"ss-catalog-service/internal/infrastructure/database"
	"ss-catalog-service/internal/infrastructure/messaging"
	pgmodel "ss-catalog-service/internal/repository/postgres"
	inventoryusecase "ss-catalog-service/internal/usecase/inventory"
	productusecase "ss-catalog-service/internal/usecase/product"
	variantusecase "ss-catalog-service/internal/usecase/variant"
	"ss-catalog-service/internal/worker"
)

func main() {
	// Load .env (non-fatal, env vars may come from OS in production)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, reading from system environment")
	}

	// --- Infrastructure: Database ---
	dbCfg := database.NewConfig()
	db, err := database.NewPostgresDB(dbCfg)
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
	productCmd := productusecase.NewProductCommandUsecase(productRepo, outboxRepo, txManager)

	defaultLang := os.Getenv("DEFAULT_LANG")
	if defaultLang == "" {
		defaultLang = "id-ID"
	}
	productQry := productusecase.NewProductQueryUsecase(productRepo, defaultLang)

	variantRepo := pgmodel.NewVariantRepository(db)
	variantCmd := variantusecase.NewVariantCommandUsecase(variantRepo, productRepo, txManager)

	inventoryRepo := pgmodel.NewInventoryRepository(db)
	inventoryCmd := inventoryusecase.NewInventoryCommandUsecase(inventoryRepo, txManager)

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
	})

	// --- Start Server ---
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("🚀 Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

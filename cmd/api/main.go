package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"

	apphttp "ss-catalog-service/internal/delivery/http"
	"ss-catalog-service/internal/infrastructure/database"
	pgmodel "ss-catalog-service/internal/repository/postgres"
	productusecase "ss-catalog-service/internal/usecase/product"
	variantusecase "ss-catalog-service/internal/usecase/variant"
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

	// --- Dependency Injection (Composition Root) ---
	productRepo := pgmodel.NewProductRepository(db)
	productCmd := productusecase.NewProductCommandUsecase(productRepo)
	productQry := productusecase.NewProductQueryUsecase(productRepo)

	txManager := pgmodel.NewTransactionManager(db)
	variantRepo := pgmodel.NewVariantRepository(db)
	variantCmd := variantusecase.NewVariantCommandUsecase(variantRepo, productRepo, txManager)

	// --- HTTP Router ---
	r := gin.Default()
	apphttp.SetupRouter(r, apphttp.RouterConfig{
		ProductCommandUsecase: productCmd,
		ProductQueryUsecase:   productQry,
		VariantCommandUsecase: variantCmd,
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

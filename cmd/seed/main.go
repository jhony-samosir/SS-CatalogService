package main

import (
	"context"
	"log"
	"ss-catalog-service/config"
	"ss-catalog-service/internal/domain"
	"ss-catalog-service/internal/infrastructure/database"
	"ss-catalog-service/internal/repository/postgres"
	msrepo "ss-catalog-service/internal/repository/meilisearch"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	msRepo, err := msrepo.NewProductSearchRepository("http://localhost:7700", "masterKey")
	if err != nil {
		log.Fatalf("Failed to connect to Meilisearch: %v", err)
	}

	seedData(db, msRepo)
}

func seedData(db *gorm.DB, msRepo domain.SearchRepository) {
	// ... (categories code same)
	categories := []postgres.CategoryModel{
		{Name: "Snacks", Slug: "snacks", Description: "Various snacks"},
		{Name: "Drinks", Slug: "drinks", Description: "Fresh drinks"},
	}
	for i := range categories {
		db.FirstOrCreate(&categories[i], postgres.CategoryModel{Slug: categories[i].Slug})
	}

	// 2. Products
	products := []postgres.ProductModel{
		{
			Name: "Kripik Balado",
			Slug: "kripik-balado",
			Description: "Spicy and crispy cassava chips from Padang.",
			ShortDesc: "Authentic Spicy Cassava Chips",
			Status: "active",
			IsFeatured: true,
			BaseModel: postgres.BaseModel{PublicID: uuid.New()},
		},
		{
			Name: "Kue Nastar Premium",
			Slug: "kue-nastar-premium",
			Description: "Pineapple tart cookies with premium butter.",
			ShortDesc: "Premium Pineapple Tarts",
			Status: "active",
			IsFeatured: true,
			BaseModel: postgres.BaseModel{PublicID: uuid.New()},
		},
		{
			Name: "Emping Melinjo",
			Slug: "emping-melinjo",
			Description: "Authentic bitter-sweet crackers from Java.",
			ShortDesc: "Traditional Melinjo Crackers",
			Status: "active",
			IsFeatured: false,
			BaseModel: postgres.BaseModel{PublicID: uuid.New()},
		},
	}
	
	// Add Image URLs (using local public assets for now)
	imageUrls := []string{
		"/images/cat-keripik.png",
		"/images/cat-kue-kering.png",
		"/images/cat-kerupuk.png",
	}

	for i := range products {
		products[i].ImageURL = imageUrls[i]
		if err := db.Where(postgres.ProductModel{Slug: products[i].Slug}).FirstOrCreate(&products[i]).Error; err == nil {
			// Update image URL if already exists
			db.Model(&products[i]).Update("ImageURL", imageUrls[i])
			
			// Add to category
			db.Model(&products[i]).Association("Categories").Append(&categories[0])
			
			// Index in Meilisearch
			pDomain := domain.Product{
				BaseEntity: domain.BaseEntity{
					ID:       products[i].ID,
					PublicID: products[i].PublicID,
				},
				Name:        products[i].Name,
				Description: products[i].Description,
				ImageURL:    products[i].ImageURL,
				Categories: []domain.Category{
					{Name: categories[0].Name},
				},
			}
			msRepo.IndexProduct(context.Background(), pDomain)
		}
	}

	log.Println("Seeding completed successfully with Meilisearch sync!")
}

package msrepo

import (
	"context"
	"encoding/json"
	"fmt"
	"ss-catalog-service/internal/domain"

	"github.com/meilisearch/meilisearch-go"
)

type productSearchRepository struct {
	service meilisearch.ServiceManager
	index   meilisearch.IndexManager
}

func NewProductSearchRepository(url, masterKey string) (domain.SearchRepository, error) {
	service := meilisearch.New(url, meilisearch.WithAPIKey(masterKey))
	index := service.Index("products")
	
	// Configure index settings
	filterable := []interface{}{"categories", "price", "brand"}
	_, err := index.UpdateFilterableAttributes(&filterable)
	if err != nil {
		return nil, err
	}
	
	searchable := []string{"name", "description", "tags"}
	_, err = index.UpdateSearchableAttributes(&searchable)
	if err != nil {
		return nil, err
	}

	return &productSearchRepository{
		service: service,
		index:   index,
	}, nil
}

func (r *productSearchRepository) IndexProduct(ctx context.Context, product domain.Product) error {
	doc := map[string]interface{}{
		"id":          product.ID,
		"public_id":   product.PublicID.String(),
		"name":        product.Name,
		"description": product.Description,
		"price":       0, // Should come from base variant
		"categories":  []string{},
		"brand":       "SamStore",
	}

	for _, c := range product.Categories {
		doc["categories"] = append(doc["categories"].([]string), c.Name)
	}

	_, err := r.index.AddDocuments(doc, nil)
	return err
}

func (r *productSearchRepository) Search(ctx context.Context, q domain.GetProductSearchQuery) (*domain.FacetedSearchResult, error) {
	searchReq := &meilisearch.SearchRequest{
		Limit:  int64(q.Limit),
		Facets: []string{"categories", "brand"},
	}

	keyword := ""
	if q.Keyword != nil {
		keyword = *q.Keyword
	}

	// Build filter string
	filter := ""
	if q.MinPrice != nil {
		filter = fmt.Sprintf("price >= %f", *q.MinPrice)
	}
	if q.MaxPrice != nil {
		if filter != "" {
			filter += " AND "
		}
		filter += fmt.Sprintf("price <= %f", *q.MaxPrice)
	}
	searchReq.Filter = filter

	resp, err := r.index.Search(keyword, searchReq)
	if err != nil {
		return nil, err
	}

	items := make([]domain.Product, len(resp.Hits))
	for i, hit := range resp.Hits {
		var p domain.Product
		// Meilisearch Hit is a map[string]json.RawMessage
		if err := hit.DecodeInto(&p); err == nil {
			items[i] = p
		}
	}

	facets := []domain.SearchFacet{}
	if len(resp.FacetDistribution) > 0 {
		var distribution map[string]map[string]int
		if err := json.Unmarshal(resp.FacetDistribution, &distribution); err == nil {
			for name, values := range distribution {
				facet := domain.SearchFacet{Name: name}
				for val, count := range values {
					facet.Values = append(facet.Values, struct {
						Value string `json:"value"`
						Count int    `json:"count"`
					}{Value: val, Count: count})
				}
				facets = append(facets, facet)
			}
		}
	}

	return &domain.FacetedSearchResult{
		Items:     items,
		Facets:    facets,
		TotalHint: resp.TotalHits,
	}, nil
}

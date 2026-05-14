package dto

type CreateWarehouseRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Code        string `json:"code" validate:"required,uppercase"`
	City        string `json:"city" validate:"required"`
	Province    string `json:"province"`
	CountryCode string `json:"country_code" validate:"required,len=2"`
	PostalCode  string `json:"postal_code"`
	Address     string `json:"address" validate:"required"`
	IsActive    bool   `json:"is_active"`
}

type UpdateWarehouseRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Code        string `json:"code" validate:"required,uppercase"`
	City        string `json:"city" validate:"required"`
	Province    string `json:"province"`
	CountryCode string `json:"country_code" validate:"required,len=2"`
	PostalCode  string `json:"postal_code"`
	Address     string `json:"address" validate:"required"`
	IsActive    bool   `json:"is_active"`
}

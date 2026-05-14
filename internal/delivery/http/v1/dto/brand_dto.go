package dto

type CreateBrandRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Slug        string `json:"slug" validate:"required,lowercase"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	WebsiteURL  string `json:"website_url" validate:"omitempty,url"`
	Description string `json:"description" validate:"max=1000"`
	IsActive    bool   `json:"is_active"`
}

type UpdateBrandRequest struct {
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Slug        string `json:"slug" validate:"required,lowercase"`
	LogoURL     string `json:"logo_url" validate:"omitempty,url"`
	WebsiteURL  string `json:"website_url" validate:"omitempty,url"`
	Description string `json:"description" validate:"max=1000"`
	IsActive    bool   `json:"is_active"`
}

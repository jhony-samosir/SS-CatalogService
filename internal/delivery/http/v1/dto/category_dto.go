package dto

type CreateCategoryRequest struct {
	ParentID    *int   `json:"parent_id"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Slug        string `json:"slug" validate:"required,lowercase,alphanumhyphen"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	Description string `json:"description" validate:"max=1000"`
	SortOrder   int    `json:"sort_order" validate:"min=0"`
	IsActive    bool   `json:"is_active"`
}

type UpdateCategoryRequest struct {
	ParentID    *int   `json:"parent_id"`
	Name        string `json:"name" validate:"required,min=2,max=255"`
	Slug        string `json:"slug" validate:"required,lowercase,alphanumhyphen"`
	IconURL     string `json:"icon_url" validate:"omitempty,url"`
	Description string `json:"description" validate:"max=1000"`
	SortOrder   int    `json:"sort_order" validate:"min=0"`
	IsActive    bool   `json:"is_active"`
}

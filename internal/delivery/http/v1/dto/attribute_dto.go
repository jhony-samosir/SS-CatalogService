package dto

type CreateAttributeRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=255"`
	Code      string `json:"code" binding:"required,uppercase"`
	InputType string `json:"input_type" binding:"required,oneof=text select multiselect boolean number"`
	IsVariant bool   `json:"is_variant"`
	SortOrder int    `json:"sort_order" binding:"min=0"`
}

type UpdateAttributeRequest struct {
	Name      string `json:"name" binding:"required,min=2,max=255"`
	Code      string `json:"code" binding:"required,uppercase"`
	InputType string `json:"input_type" binding:"required,oneof=text select multiselect boolean number"`
	IsVariant bool   `json:"is_variant"`
	SortOrder int    `json:"sort_order" binding:"min=0"`
}

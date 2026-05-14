package dto

type CreateAttributeRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=255"`
	Code      string `json:"code" validate:"required,uppercase"`
	InputType string `json:"input_type" validate:"required,oneof=text select multiselect boolean number"`
	IsVariant bool   `json:"is_variant"`
	SortOrder int    `json:"sort_order" validate:"min=0"`
}

type UpdateAttributeRequest struct {
	Name      string `json:"name" validate:"required,min=2,max=255"`
	Code      string `json:"code" validate:"required,uppercase"`
	InputType string `json:"input_type" validate:"required,oneof=text select multiselect boolean number"`
	IsVariant bool   `json:"is_variant"`
	SortOrder int    `json:"sort_order" validate:"min=0"`
}

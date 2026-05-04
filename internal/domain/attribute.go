package domain

// AttributeInputType defines how the attribute value is collected.
type AttributeInputType string

const (
	InputTypeText        AttributeInputType = "text"
	InputTypeSelect      AttributeInputType = "select"
	InputTypeMultiSelect AttributeInputType = "multiselect"
	InputTypeBoolean     AttributeInputType = "boolean"
	InputTypeNumber      AttributeInputType = "number"
)

// ProductAttribute represents attribute definitions (e.g., Color, Size).
type ProductAttribute struct {
	BaseEntity
	Name      string
	Code      string
	InputType AttributeInputType
	IsVariant bool
	SortOrder int
}

// AttributeValue represents possible values for each attribute.
type AttributeValue struct {
	BaseEntity
	AttributeID int
	Value       string
	ColorHex    string
	SortOrder   int
}

// ProductVariantAttribute maps a variant to its specific attribute values.
type ProductVariantAttribute struct {
	BaseEntity
	VariantID        int
	AttributeID      int
	AttributeValueID int
}

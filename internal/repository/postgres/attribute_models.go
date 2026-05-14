package postgres

// ProductAttributeModel represents the database schema for product attributes.
type ProductAttributeModel struct {
	BaseModel
	Name      string `gorm:"type:varchar(255);not null"`
	Code      string `gorm:"type:varchar(100);not null;uniqueIndex"`
	InputType string `gorm:"type:varchar(50);not null;default:'text'"`
	IsVariant bool   `gorm:"not null;default:false"`
	SortOrder int    `gorm:"not null;default:0"`

	// Associations
	Values []AttributeValueModel `gorm:"foreignKey:AttributeID"`
}

func (ProductAttributeModel) TableName() string { return "product_attributes" }

// AttributeValueModel represents the database schema for attribute values.
type AttributeValueModel struct {
	BaseModel
	AttributeID int    `gorm:"not null;index"`
	Value       string `gorm:"type:varchar(255);not null"`
	ColorHex    string `gorm:"type:varchar(50)"`
	SortOrder   int    `gorm:"not null;default:0"`
}

func (AttributeValueModel) TableName() string { return "attribute_values" }

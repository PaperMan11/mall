package model

type Cart struct {
	Model
	UserID   uint       `gorm:"unique" json:"user_id"`
	ItemInfo []ItemInfo `gorm:"foreignKey:CartId" json:"item_info,omitempty"`
}

// -------------------------------------------------------------------

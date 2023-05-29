package model

type Category struct {
	Model
	CategoryName string `gorm:"unique" json:"category_name"`
}

// -------------------------------------------------------------------

type CategoryList struct {
	Total         int         `json:"total"`
	CategoryInfos []*Category `json:"category_infos"`
}

type CategoryAddReq struct {
	CategoryName string `json:"category_name"`
}

type CategoryRemoveReq struct {
	CategoryId uint `json:"category_id"`
}

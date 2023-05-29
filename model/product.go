package model

type Product struct {
	Model
	Name          string       `gorm:"size:255;index" json:"name"`
	CategoryId    uint         `gorm:"not null" json:"category_id"`
	Category      Category     `gorm:"foreignKey:CategoryId" json:"category,omitempty"`
	Title         string       `json:"title"`
	Info          string       `gorm:"size:1000" json:"info"`
	Price         string       `json:"price"`
	DiscountPrice string       `json:"discount_price"`
	OnSale        bool         `gorm:"default:false" json:"on_sale"` // 是否出售
	Num           int          `json:"num"`
	BossId        uint         `json:"boss_id"`
	ProductImgs   []ProductImg `gorm:"foreignKey:ProductId" json:"product_imgs,omitempty"`
}

// -------------------------------------------------------------------

type ProductCreateReq struct {
	Name          string `form:"name" json:"name"`
	CategoryID    uint   `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`
}

type ProductList struct {
	Products []*Product `json:"products"`
	Total    int64      `json:"total"`
}

type ProductUpdateReq struct {
	Name          string `form:"name" json:"name"`
	CategoryID    uint   `form:"category_id" json:"category_id"`
	Title         string `form:"title" json:"title" `
	Info          string `form:"info" json:"info" `
	Price         string `form:"price" json:"price"`
	DiscountPrice string `form:"discount_price" json:"discount_price"`
	OnSale        bool   `form:"on_sale" json:"on_sale"`
	Num           int    `form:"num" json:"num"`
}

type ProductInfo struct {
	Product
	ViewCount uint64 `json:"view_count"`
}

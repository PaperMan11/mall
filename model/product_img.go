package model

type ProductImg struct {
	Model
	ProductId uint   `json:"product_id"`
	ImgPath   string `json:"img_path"`
}

// -------------------------------------------------------------------

type ProductImgList struct {
	Total           int64         `json:"total"`
	ProductImgInfos []*ProductImg `json:"product_img_infos"`
}

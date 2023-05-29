package service

import (
	"fmt"
	"io"
	"mall/middleware"
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/utils"
	"mall/repositry/mysqldb"
	"mall/svc"
	"mime/multipart"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewProductLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *ProductLogic {
	return &ProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ProductLogic) ProductCreateLogic(req *model.ProductCreateReq, files []*multipart.FileHeader) (resp *model.ProductInfo, code int) {

	// 1 创建商品
	userIdV, _ := l.ctx.Get(middleware.CtxUserIdKey)
	userId, ok := userIdV.(uint)
	if !ok {
		l.svcCtx.Log.Errorf("key [%s] get value faild", middleware.CtxUserIdKey)
		return nil, e.ERROR
	}

	product := model.Product{
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		Name:          req.Name,
		CategoryId:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		Num:           req.Num,
		OnSale:        true,
		BossId:        userId,
	}
	if err := l.svcCtx.ProductModel.CreateProduct(&product); err != nil {
		l.svcCtx.Log.Errorf("ProductModel.CreateProduct err: %s", err)
		return nil, e.ErrorDatabase
	}

	// 2 创建商品图片
	for _, file := range files {
		fd, err := file.Open()
		if err != nil {
			l.svcCtx.Log.Errorf("file.Open err: %s", err)
			return nil, e.ErrorUploadFile
		}
		filePath, err := l.productUpload2Local(fd, userId, file.Filename)
		if err != nil {
			l.svcCtx.Log.Errorf("productUpload2Local err: %s", err)
			return nil, e.ErrorUploadFile
		}
		productImg := &model.ProductImg{
			Model: model.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			ProductId: product.ID,
			ImgPath:   filePath,
		}
		if err = l.svcCtx.ProductImgModel.CreateProductImg(productImg); err != nil {
			l.svcCtx.Log.Errorf("ProductImgModel.CreateProductImg err: %s", err)
			return nil, e.ErrorDatabase
		}
	}

	return &model.ProductInfo{
		Product: product,
	}, e.SUCCESS
}

func (l *ProductLogic) productUpload2Local(file multipart.File, bossId uint, productName string) (filePath string, err error) {
	// base path
	bId := strconv.Itoa(int(bossId))
	basePath := path.Join(".", l.svcCtx.Config.LocalUploadConf.ProductPhotoPath, "boss"+bId)
	if !utils.DirExist(basePath) {
		utils.CreateDir(basePath)
	}
	// read file and save
	productPath := path.Join(basePath, productName)
	fd, err := os.OpenFile(productPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		l.svcCtx.Log.Errorf("os.OpenFile err: %s", err)
		return "", err
	}
	_, err = io.Copy(fd, file)
	if err != nil {
		l.svcCtx.Log.Errorf("io.Copy err: %s", err)
		return "", err
	}
	return fmt.Sprintf("boss%s/%s", bId, productName), err
}

func (l *ProductLogic) ProductUpdateLogic(req *model.ProductUpdateReq, product_id uint) (nProductInfo *model.ProductInfo, code int) {
	// update
	newProduct := model.Product{
		Model: model.Model{
			UpdatedAt: time.Now(),
		},
		Name:          req.Name,
		CategoryId:    req.CategoryID,
		Title:         req.Title,
		Info:          req.Info,
		Price:         req.Price,
		DiscountPrice: req.DiscountPrice,
		OnSale:        req.OnSale,
	}

	var product *model.Product
	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		productDAO := mysqldb.NewProductModel(tx, "product")
		err := productDAO.UpdateProduct(product_id, &newProduct)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction UpdateProduct err: %s", err)
			return err
		}

		product, err = productDAO.GetProductById(product_id)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction GetProductById err: %s", err)
			return err
		}
		return nil
	}); err != nil {
		l.svcCtx.Log.Errorf("Transaction ProductUpdateSelect err:%s", err)
		return nil, e.ErrorDatabase
	}

	// 查 redis
	viewCount, err := l.svcCtx.ProductRedisCache.ViewClick(product_id)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductRedisCache.ViewClick err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &model.ProductInfo{
		Product:   *product,
		ViewCount: viewCount,
	}, e.SUCCESS
}

func (l *ProductLogic) ProductDeleteLogic(product_id uint) (code int) {
	// 1 product exist?
	count, err := l.svcCtx.ProductModel.CountProductByCondition(map[string]interface{}{
		"id": product_id,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.CountProductByCondition err: %s", err)
		return e.ErrorDatabase
	}
	if count == 0 {
		l.svcCtx.Log.Errorf("product [%d] not exist", product_id)
		return e.ErrorNotExistProduct
	}
	// 2 delete local_imgs
	productImgs, err := l.svcCtx.ProductImgModel.ListProductImgByProductId(product_id)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.ListProductImgByProductId err: %s", err)
		return e.ErrorDatabase
	}
	basePath := path.Join(".", l.svcCtx.Config.LocalUploadConf.ProductPhotoPath)
	for _, pImg := range productImgs {
		if err := utils.DeletFiles(fmt.Sprintf("%s/%s", basePath, pImg.ImgPath)); err != nil {
			l.svcCtx.Log.Errorf("delete productImgs err: %s", err)
			return e.ERROR
		}
	}
	// 3 delete product and imgs
	if err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		pImgDAO := mysqldb.NewProductImgModel(tx, "product_img")
		err := pImgDAO.DeleteProductImgByPid(product_id)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction DeleteProductImgByPid err: %s", err)
			return err
		}

		productDAO := mysqldb.NewProductModel(tx, "product")
		if err = productDAO.DeleteProduct(product_id); err != nil {
			l.svcCtx.Log.Errorf("Transaction DeleteProductById err: %s", err)
			return err
		}
		return nil
	}); err != nil {
		l.svcCtx.Log.Errorf("Transaction DeleteProductAndImgs err:%s", err)
		return e.ErrorDatabase
	}

	return e.SUCCESS
}

func (l *ProductLogic) ProductListLogic(req *utils.BasePage, catgory_id uint) (resp *model.ProductList, code int) {
	var (
		total     int64
		condition = make(map[string]interface{})
	)

	condition["category_id"] = catgory_id

	// total
	total, err := l.svcCtx.ProductModel.CountProductByCondition(condition)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.CountProductByCondition err: %s", err)
		return nil, e.ErrorDatabase
	}
	if total == 0 {
		return &model.ProductList{Products: nil, Total: total}, e.SUCCESS
	}

	// products
	// condition["category_id"] = catgory_id
	products, err := l.svcCtx.ProductModel.ListProductByCondition(condition, req)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.ListProductByCondition err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &model.ProductList{
		Products: products,
		Total:    total,
	}, e.SUCCESS
}

func (l *ProductLogic) ProductShowLogic(product_id uint) (resp *model.ProductInfo, code int) {
	// search
	count, err := l.svcCtx.ProductModel.CountProductByCondition(map[string]interface{}{
		"id": product_id,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.CountProductByCondition err: %s", err)
		return nil, e.ErrorDatabase
	}
	if count == 0 {
		l.svcCtx.Log.Errorf("product [%d] not exist", product_id)
		return nil, e.ErrorNotExistProduct
	}

	p, err := l.svcCtx.ProductModel.GetProductInfoById(product_id)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.GetProductById err: %s", err)
		return nil, e.ErrorDatabase
	}

	// 点击
	l.svcCtx.ProductRedisCache.AddClick(p.ID)
	viewCount, err := l.svcCtx.ProductRedisCache.ViewClick(p.ID)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductRedisCache.ViewClick err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &model.ProductInfo{
		Product:   *p,
		ViewCount: viewCount,
	}, e.SUCCESS
}

func (l *ProductLogic) ProductSearchLogic(req *utils.BasePage, searchInfo string) (resp *model.ProductList, code int) {

	products, count, err := l.svcCtx.ProductModel.SearchProduct(searchInfo, req)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.SearchProduct err: %s", err)
		return nil, e.ErrorDatabase
	}
	if count == 0 {
		return &model.ProductList{
			Products: nil,
			Total:    int64(count),
		}, e.SUCCESS
	}

	return &model.ProductList{
		Products: products,
		Total:    int64(count),
	}, e.SUCCESS
}

func (l *ProductLogic) ProductImgListLogic(product_id uint) (resp *model.ProductImgList, code int) {
	productImgInfos, err := l.svcCtx.ProductImgModel.ListProductImgByProductId(product_id)
	if err != nil {
		l.svcCtx.Log.Errorf("ProductRedisCache.ViewClick err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &model.ProductImgList{
		Total:           int64(len(productImgInfos)),
		ProductImgInfos: productImgInfos,
	}, e.SUCCESS
}

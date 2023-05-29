package service

import (
	"mall/model"
	"mall/pkg/e"
	"mall/repositry/mysqldb"
	"mall/svc"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ItemInfoLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewItemInfoLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *ItemInfoLogic {
	return &ItemInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ItemInfoLogic) ItemInfoAddLogic(cartID, productID uint, req *model.ItemInfoAddReq) (itemResp *model.ItemInfo, code int) {
	// 1 是否有该商品
	count, err := l.svcCtx.ProductModel.CountProductByCondition(map[string]interface{}{
		"id": productID,
	})
	if err != nil {
		l.svcCtx.Log.Errorf("ProductModel.CountProductByCondition err: %s", err)
		return nil, e.ErrorDatabase
	}
	if count == 0 {
		l.svcCtx.Log.Infof("product [%d] not exist", productID)
		return nil, e.ErrorNotExistProduct
	}

	// 2 购物车中是否有
	item, exist, err := l.svcCtx.ItemInfoModel.ExistByCartID(cartID, productID)
	if err != nil {
		l.svcCtx.Log.Errorf("ItemInfoModel.ExistByCartID err: %s", err)
		return nil, e.ErrorDatabase
	}
	// 3 添加到购物车
	if exist {
		item.Num += req.Num
		item.UpdatedAt = time.Now()
		if err = l.svcCtx.ItemInfoModel.UpdateByID(item.ID, item); err != nil {
			l.svcCtx.Log.Errorf("ItemInfoModel.UpdateByID err: %s", err)
			return nil, e.ErrorDatabase
		}
	} else {
		itemInfo := model.ItemInfo{
			Model: model.Model{
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			CartId:    cartID,
			ProductId: productID,
			Num:       req.Num,
		}
		if err = l.svcCtx.ItemInfoModel.Create(&itemInfo); err != nil {
			l.svcCtx.Log.Errorf("ItemInfoModel.Create err: %s", err)
			return nil, e.ErrorDatabase
		}
	}

	return item, e.SUCCESS
}

func (l *ItemInfoLogic) ItemInfoUpdateLogic(itemInfoID uint, req *model.ItemInfoUpdateReq) (itemInfo *model.ItemInfo, code int) {
	// 1 exist?
	item, exist, err := l.svcCtx.ItemInfoModel.ExistByID(itemInfoID)
	if err != nil {
		l.svcCtx.Log.Errorf("ItemInfoModel.ExistByID err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("item_info [%d] not exist", itemInfoID)
		return nil, e.ErrorNotExistItem
	}

	// update select
	item.Num += req.Num
	item.UpdatedAt = time.Now()
	if item.Num < 0 {
		item.Num = 0
	}

	if err := l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		ItemInfoDAO := mysqldb.NewItemInfoModel(tx, "item_info")
		err := ItemInfoDAO.UpdateByID(itemInfoID, item)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction UpdateByID err: %s", err)
			return err
		}

		itemInfo, err = ItemInfoDAO.ShowById(itemInfoID)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction ShowById err: %s", err)
			return err
		}
		return nil
	}); err != nil {
		l.svcCtx.Log.Errorf("Transaction ItemInfoUpdateSelect err:%s", err)
		return nil, e.ErrorDatabase
	}

	return itemInfo, e.SUCCESS
}

func (l *ItemInfoLogic) ItemInfoDeleteLogic(itemInfoID uint) (code int) {

	_, exist, err := l.svcCtx.ItemInfoModel.ExistByID(itemInfoID)
	if err != nil {
		l.svcCtx.Log.Errorf("ItemInfoModel.ExistByID err: %s", err)
		return e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("item_info [%d] not exist", itemInfoID)
		return e.ErrorNotExistItem
	}

	err = l.svcCtx.ItemInfoModel.Delete(itemInfoID)
	if err != nil {
		l.svcCtx.Log.Errorf("ItemInfoModel.Delete err: %s", err)
		return e.ErrorDatabase
	}
	return e.SUCCESS
}

func (l *ItemInfoLogic) ItemInfoShowLogic(itemInfoID uint) (item *model.ItemInfo, code int) {
	_, exist, err := l.svcCtx.ItemInfoModel.ExistByID(itemInfoID)
	if err != nil {
		l.svcCtx.Log.Errorf("ItemInfoModel.ExistByID err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("item_info [%d] not exist", itemInfoID)
		return nil, e.ErrorNotExistItem
	}

	i, err := l.svcCtx.ItemInfoModel.ShowById(itemInfoID)
	if err != nil {
		l.svcCtx.Log.Errorf("ItemInfoModel.ShowById err: %s", err)
		return nil, e.ErrorDatabase
	}
	return i, e.SUCCESS
}

package service

import (
	"mall/model"
	"mall/pkg/e"
	"mall/svc"

	"github.com/gin-gonic/gin"
)

type CartLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCartLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *CartLogic {
	return &CartLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CartLogic) CartInfoLogic(cartId uint) (*model.Cart, int) {

	ci, err := l.svcCtx.CartModel.ShowCartById(cartId)
	if err != nil {
		l.svcCtx.Log.Errorf("CartModel.ShowCartById err: %s", err)
		return nil, e.ErrorDatabase
	}
	return ci, e.SUCCESS
}

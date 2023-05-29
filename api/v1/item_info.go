package v1

import (
	"mall/model"
	"mall/pkg/e"
	"mall/service"
	"mall/svc"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 	添加购物车商品
// @Description	添加购物车商品
// @Tags		item_info
// @Accept		json
// @Produce		json
// @Param 		req body model.ItemInfoAddReq true "商品数量"
// @Param		cart_id path uint true "cart id"
// @Param		product_id path uint true "product id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/carts/{cart_id}/item-infos/products/{product_id} [post]
func ItemInfoAddHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		cartIDStr := ctx.Param("cart_id")
		cartID, err := strconv.ParseUint(cartIDStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", cartIDStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		productIDStr := ctx.Param("product_id")
		productID, err := strconv.ParseUint(productIDStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", productIDStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		var req model.ItemInfoAddReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewItemInfoLogic(ctx, svc)
		resp, code := l.ItemInfoAddLogic(uint(cartID), uint(productID), &req)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary		更新购物车商品
// @Description	更新购物车商品
// @Tags		item_info
// @Accept		json
// @Produce		json
// @Param 		req body model.ItemInfoUpdateReq true "商品数量"
// @Param		id path uint true "iteminfo id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/carts/item-infos/{id} [put]
func ItemInfoUpdateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("id")

		itemInfoID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		var req model.ItemInfoUpdateReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewItemInfoLogic(ctx, svc)
		resp, code := l.ItemInfoUpdateLogic(uint(itemInfoID), &req)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary		删除购物车商品
// @Description	删除购物车商品
// @Tags		item_info
// @Accept		json
// @Produce		json
// @Param		id path uint true "iteminfo id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/carts/item-infos/{id} [delete]
func ItemInfoDeleteHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		itemInfoID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewItemInfoLogic(ctx, svc)
		code := l.ItemInfoDeleteLogic(uint(itemInfoID))
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary		获取购物车商品
// @Description	获取购物车商品
// @Tags		item_info
// @Accept		json
// @Produce		json
// @Param		id path uint true "iteminfo id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/carts/item-infos/{id} [get]
func ItemInfoShowHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("id")
		itemInfoID, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewItemInfoLogic(ctx, svc)
		resp, code := l.ItemInfoShowLogic(uint(itemInfoID))
		model.RespWithCode(ctx, resp, code)
	}
}

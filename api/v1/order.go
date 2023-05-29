package v1

import (
	"mall/middleware"
	"mall/model"
	"mall/pkg/e"
	"mall/service"
	"mall/svc"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 	订单创建
// @Description	订单创建
// @Tags		order
// @Accept		json
// @Produce		json
// @Param 		req body model.OrderCreateReq true "订单创建信息"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/orders [post]
func OrderCreateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.OrderCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}

		userIdV, _ := ctx.Get(middleware.CtxUserIdKey)
		userId, ok := userIdV.(uint)
		if !ok {
			svc.Log.Errorf("key [%s] get value faild", middleware.CtxUserIdKey)
			model.RespError(ctx, nil, e.ErrorAdminFindUser)
			return
		}

		l := service.NewOrderLogic(ctx, svc)
		resp, code := l.OrderCreateLogic(userId, &req)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	订单删除
// @Description	订单删除
// @Tags		order
// @Produce		json
// @Param 		order_id path int true "订单号"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/orders/{order_id} [delete]
func OrderDeleteHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("order_id")

		oId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewOrderLogic(ctx, svc)
		code := l.OrderDeleteLogic(oId)
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary 	订单查询
// @Description	订单查询
// @Tags		order
// @Produce		json
// @Param 		order_id path int true "订单号"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/orders/{order_id} [get]
func OrderGetHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("order_id")

		oId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewOrderLogic(ctx, svc)
		resp, code := l.OrderGetLogic(oId)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	订单列表
// @Description	订单列表（单个用户）
// @Tags		order
// @Produce		json
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/orders [get]
func OrderListHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		l := service.NewOrderLogic(ctx, svc)
		resp, code := l.OrderListLogic()
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	支付
// @Description	支付
// @Tags		order
// @Produce		json
// @Param 		order_id path int true "订单号"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/orders/{order_id}/payment [post]
func OrderPaymentHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("order_id")

		oId, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewOrderLogic(ctx, svc)
		code := l.OrderPaymentLogic(oId)
		model.RespWithCode(ctx, nil, code)
	}
}

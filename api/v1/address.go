package v1

import (
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/utils"
	"mall/service"
	"mall/svc"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 	地址创建
// @Description	地址创建
// @Tags		address
// @Accept		json
// @Produce		json
// @Param 		req body model.AddressCreateReq true "地址信息"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/addresses [post]
func AddressCreateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.AddressCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewAddressLogic(ctx, svc)
		resp, code := l.AddressCreateLogic(&req)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	地址获取
// @Description	地址获取
// @Tags		address
// @Accept		json
// @Produce		json
// @Param 		id path uint true "address id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/addresses/{id} [get]
func AddressGetHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		addresssId := ctx.Param("id")

		id, err := strconv.ParseUint(addresssId, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", addresssId)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewAddressLogic(ctx, svc)
		resp, code := l.AddressGetLogic(uint(id))
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	地址删除
// @Description	地址删除
// @Tags		address
// @Produce		json
// @Param 		id path uint true "address id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/addresses/{id} [delete]
func AddressDeleteHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		addresssId := ctx.Param("id")

		id, err := strconv.ParseUint(addresssId, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", addresssId)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewAddressLogic(ctx, svc)
		code := l.AddressDeleteLogic(uint(id))
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary 	地址更新
// @Description	地址更新
// @Tags		address
// @Accept		json
// @Produce		json
// @Param 		req body model.AddressUpdateReq true "地址信息"
// @Param 		id path uint true "address id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/addresses/{id} [put]
func AddressUpdateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		addressIdStr := ctx.Param("id")
		if addressIdStr == "" {
			svc.Log.Errorf("id not exist")
			model.RespError(ctx, e.InvalidParams)
			return
		}

		addressId, err := strconv.ParseUint(addressIdStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", addressIdStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		var req model.AddressUpdateReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewAddressLogic(ctx, svc)
		resp, code := l.AddressUpdateLogic(&req, uint(addressId))
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	地址获取
// @Description	地址获取
// @Tags		address
// @Produce		json
// @Param		offset query int false "page number"
// @Param		limit query int false "page size"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/addresses [get]
func AddressListHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var (
			offset, limit int
			err           error
		)

		offsetStr := ctx.Query("offset")
		limitStr := ctx.Query("limit")

		if offset, err = strconv.Atoi(offsetStr); err != nil {
			offset = 1
		}

		if limit, err = strconv.Atoi(limitStr); err != nil {
			limit = 10
		}

		req := utils.BasePage{
			PageNum:  offset,
			PageSize: limit,
		}
		l := service.NewAddressLogic(ctx, svc)
		resp, code := l.AddressListLogic(&req)
		model.RespWithCode(ctx, resp, code)
	}
}

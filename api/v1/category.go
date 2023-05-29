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

// @Summary 	商品种类
// @Description	商品种类
// @Tags		category
// @Produce		json
// @Param		offset query int false "page number"
// @Param		limit query int false "page size"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/categories [get]
func CategoryListHandler(svc *svc.ServiceContext) gin.HandlerFunc {
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

		page := utils.BasePage{
			PageNum:  offset,
			PageSize: limit,
		}

		l := service.NewCategoryLogic(ctx, svc)
		resp, code := l.CategoryListLogic(&page)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品种类创建
// @Description	商品种类创建
// @Tags		category
// @Accept		json
// @Produce		json
// @Param 		req body model.CategoryAddReq true "商品种类信息"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/categories [post]
func CategoryCreateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req model.CategoryAddReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewCategoryLogic(ctx, svc)
		resp, code := l.CategoryCreateLogic(&req)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品种类删除
// @Description	商品种类删除
// @Tags		category
// @Accept		json
// @Produce		json
// @Param 		id path uint true "category id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/categories/{id} [delete]
func CategoryDeleteHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("id")

		categroy_id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewCategoryLogic(ctx, svc)
		code := l.CategoryRemoveLogic(uint(categroy_id))
		model.RespWithCode(ctx, nil, code)
	}
}

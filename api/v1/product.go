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

// @Summary 	商品创建
// @Description	商品创建
// @Tags		product
// @Accept		json,mpfd
// @Produce		json
// @Param 		req body model.ProductCreateReq true "商品创建信息"
// @Param		file formData file true "file"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products [get]
func ProductCreateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1 bind
		var req model.ProductCreateReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}

		// 2 获取表单文件
		form, err := ctx.MultipartForm()
		if err != nil {
			svc.Log.Errorf("MultipartForm err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		files := form.File["file"]
		l := service.NewProductLogic(ctx, svc)
		resp, code := l.ProductCreateLogic(&req, files)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品更新
// @Description	商品更新
// @Tags		product
// @Accept		json
// @Produce		json
// @Param		id path uint true "product id"
// @Param 		req body model.ProductUpdateReq true "商品更新信息"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products/{id} [put]
func ProductUpdateHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("id")

		product_id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		var req model.ProductUpdateReq
		if err := ctx.ShouldBind(&req); err != nil {
			svc.Log.Errorf("ShouldBind err: %s", err)
			model.RespError(ctx, nil, e.InvalidParams)
			return
		}
		l := service.NewProductLogic(ctx, svc)
		resp, code := l.ProductUpdateLogic(&req, uint(product_id))
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品删除
// @Description	商品删除
// @Tags		product
// @Produce		json
// @Param		id path uint true "product id"
// @Param		access_token header string true "access_token"
// @Param		refresh_token header string true "refresh_token"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products/{id} [delete]
func ProductDeleteHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("id")
		product_id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewProductLogic(ctx, svc)
		code := l.ProductDeleteLogic(uint(product_id))
		model.RespWithCode(ctx, nil, code)
	}
}

// @Summary 	商品列表
// @Description	商品列表（按种类显示）
// @Tags		product
// @Produce		json
// @Param		category_id path uint true "category id"
// @Param		offset query int false "page number"
// @Param		limit query int false "page size"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products-categories/{category_id} [get]
func ProductListHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		var (
			offset, limit int
			category_id   int
			err           error
		)

		offsetStr := ctx.Query("offset")
		limitStr := ctx.Query("limit")
		categoryIdStr := ctx.Param("category_id")

		if offset, err = strconv.Atoi(offsetStr); err != nil {
			offset = 1
		}

		if limit, err = strconv.Atoi(limitStr); err != nil {
			limit = 10
		}

		if category_id, err = strconv.Atoi(categoryIdStr); err != nil {
			category_id = 1
		}

		req := utils.BasePage{
			PageNum:  offset,
			PageSize: limit,
		}

		l := service.NewProductLogic(ctx, svc)
		resp, code := l.ProductListLogic(&req, uint(category_id))
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品显示
// @Description	商品显示
// @Tags		product
// @Produce		json
// @Param		id path uint true "product id"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products/{id} [get]
func ProductShowHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {

		idStr := ctx.Param("id")
		product_id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewProductLogic(ctx, svc)
		resp, code := l.ProductShowLogic(uint(product_id))
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品查找
// @Description	商品查找
// @Tags		product
// @Produce		json
// @Param		offset query int false "page number"
// @Param		limit query int false "page size"
// @Param		search_info query string false "page size"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products [get]
func ProductSearchHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			offset, limit int
			searchInfo    string
			err           error
		)

		offsetStr := ctx.Query("offset")
		limitStr := ctx.Query("limit")
		searchInfo = ctx.Query("search_info")

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
		l := service.NewProductLogic(ctx, svc)
		resp, code := l.ProductSearchLogic(&req, searchInfo)
		model.RespWithCode(ctx, resp, code)
	}
}

// @Summary 	商品图片
// @Description	商品图片
// @Tags		product
// @Produce		json
// @Param		id path uint true "product id"
// @Success     200  {object} model.Response
// @Failure     400  {object} model.Response
// @Router      /api/v1/products/imgs/{id} [get]
func ProductImgListHandler(svc *svc.ServiceContext) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idStr := ctx.Param("id")
		product_id, err := strconv.ParseUint(idStr, 10, 64)
		if err != nil {
			svc.Log.Errorf("[%s] InvalidParams", idStr)
			model.RespError(ctx, e.InvalidParams)
			return
		}

		l := service.NewProductLogic(ctx, svc)
		resp, code := l.ProductImgListLogic(uint(product_id))
		model.RespWithCode(ctx, resp, code)
	}
}

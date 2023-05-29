package service

import (
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/utils"
	"mall/svc"
	"time"

	"github.com/gin-gonic/gin"
)

type CategoryLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewCategoryLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *CategoryLogic {
	return &CategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CategoryLogic) CategoryListLogic(page *utils.BasePage) (resp *model.CategoryList, code int) {
	categories, err := l.svcCtx.CategoryModel.ListCategory(page)
	if err != nil {
		l.svcCtx.Log.Errorf("CategoryModel.ListCategory err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &model.CategoryList{
		Total:         len(categories),
		CategoryInfos: categories,
	}, e.SUCCESS
}

func (l *CategoryLogic) CategoryCreateLogic(req *model.CategoryAddReq) (resp *model.Category, code int) {

	// 1 判断是否存在
	exist, err := l.svcCtx.CategoryModel.ExistCategoryByName(req.CategoryName)
	if err != nil {
		l.svcCtx.Log.Errorf("CategoryModel.ExistCategoryByName err: %s", err)
		return nil, e.ErrorDatabase
	}
	if exist {
		l.svcCtx.Log.Infof("[%s] category exist: %s", req.CategoryName)
		return nil, e.ErrorExistCategory
	}
	// 2 创建
	category := model.Category{
		Model: model.Model{
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		CategoryName: req.CategoryName,
	}
	err = l.svcCtx.CategoryModel.AddCategory(&category)
	if err != nil {
		l.svcCtx.Log.Errorf("CategoryModel.AddCategory err: %s", err)
		return nil, e.ErrorDatabase
	}

	return &category, e.SUCCESS
}

func (l *CategoryLogic) CategoryRemoveLogic(category_id uint) (code int) {
	err := l.svcCtx.CategoryModel.RemoveCategoryById(category_id)
	if err != nil {
		l.svcCtx.Log.Errorf("CategoryModel.RemoveCategoryById err: %s", err)
		return e.ErrorDatabase
	}

	return e.SUCCESS
}

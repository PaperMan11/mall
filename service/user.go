package service

import (
	"fmt"
	"mall/model"
	"mall/pkg/e"
	"mall/pkg/utils/jwt"
	"mall/repositry/mysqldb"
	"mall/svc"
	"path"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserLogic struct {
	ctx    *gin.Context
	svcCtx *svc.ServiceContext
}

func NewUserLogic(ctx *gin.Context, svcCtx *svc.ServiceContext) *UserLogic {
	return &UserLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UserLogic) UserRegisterLogic(req *model.UserRegisterReq) (code int) {
	// vaildator 判断
	// // 1 密钥
	// if req.Key == "" || len(req.Key) != 16 {
	// 	return e.ErrorKey
	// }
	user := &model.User{
		UserName: req.UserName,
		NickName: req.NickName,
		Status:   model.Active,
		Email:    req.Email,
		Money:    l.svcCtx.Encrypt.AesEncoding("10000"), // 初始金额,
		Avatar:   "avatar.JPG",
	}
	// 2 是否存在
	_, exist, err := l.svcCtx.UserModel.ExistByUserName(req.UserName)
	if err != nil {
		l.svcCtx.Log.Errorf("UserModel.ExistByUserName err: %s", err)
		return e.ErrorDatabase
	}
	if exist {
		l.svcCtx.Log.Infof("User [%s] exist", req.UserName)
		return e.ErrorExistUser
	}

	// 3 create
	if err := user.SetPassword(req.Password); err != nil {
		l.svcCtx.Log.Errorf("user.SetPassword err: %s", err)
		return e.ERROR
	}

	if err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		userDAO := mysqldb.NewUserModel(tx, "user")
		err := userDAO.CreateUser(user)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction CreateUser err: %s", err)
			return err
		}

		cartDAO := mysqldb.NewCartModel(tx, "cart")
		cart := model.Cart{
			UserID: user.ID,
		}
		err = cartDAO.Create(&cart)
		if err != nil {
			l.svcCtx.Log.Errorf("Transaction CreateCart err: %s", err)
			return err
		}

		return nil
	}); err != nil {
		l.svcCtx.Log.Errorf("Transaction CreateUserAndCart err:%s", err)
		return e.ErrorDatabase
	}
	return e.SUCCESS
}

func (l *UserLogic) UserLoginLogic(req *model.UserLoginReq) (resp *model.UserLoginResp, code int) {
	// 1 user exist
	u, exist, err := l.svcCtx.UserModel.ExistByUserName(req.UserName)
	if err != nil {
		l.svcCtx.Log.Errorf("UserModel.ExistByUserName err: %s", err)
		return nil, e.ErrorDatabase
	}
	if !exist {
		l.svcCtx.Log.Infof("User [%s] not exist", req.UserName)
		return nil, e.ErrorNotExistUser
	}

	// 2 password
	if ok := u.CheckPassword(req.Password); !ok {
		l.svcCtx.Log.Infof("User [%s] password different", req.UserName)
		return nil, e.ErrorNotComparePassword
	}

	// 3 gen token
	accessToken, refreshToken, err := jwt.GenerateToken(u.ID, u.UserName)
	if err != nil {
		l.svcCtx.Log.Errorf("jwt.GenerateToken err: %s", err)
		return nil, e.ErrorAuthToken
	}

	userInfo := &model.UserInfo{
		Id:       u.ID,
		UserName: u.UserName,
		NickName: u.NickName,
		Email:    u.Email,
		Status:   u.Status,
		CreateAt: u.CreatedAt.Unix(),
	}
	userInfo.Avatar = path.Join(
		fmt.Sprintf("%s:%d", l.svcCtx.Config.LocalUploadConf.PhotoHost,
			l.svcCtx.Config.HttpServerConf.Port),
		l.svcCtx.Config.LocalUploadConf.AvatarPath,
		u.Avatar)
	return &model.UserLoginResp{
		TokenInfo: &model.TokenInfo{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
		UserInfo: userInfo,
	}, e.SUCCESS
}

func (l *UserLogic) UserShowLogic(userId uint) (user *model.User, code int) {
	user, err := l.svcCtx.UserModel.GetUserById(userId)
	if err != nil {
		l.svcCtx.Log.Errorf("UserModel.GetUserById err: %s", err)
		return nil, e.ErrorDatabase
	}
	user.Money = l.svcCtx.Encrypt.AesDecoding(user.Money)
	return user, e.SUCCESS
}

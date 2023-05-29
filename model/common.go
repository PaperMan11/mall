package model

import (
	"mall/pkg/e"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// common model

type Model struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
}

// --------------------------------------------------

type Response struct {
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
	Msg    string      `json:"msg"`
}

func RespWithCode(ctx *gin.Context, data interface{}, code int) {
	ctx.JSON(http.StatusOK, &Response{
		Status: code,
		Data:   data,
		Msg:    e.GetMsg(code),
	})
}

func RespSuccess(ctx *gin.Context, data interface{}, code ...int) {
	status := e.SUCCESS
	if code != nil {
		status = code[0]
	}
	ctx.JSON(http.StatusOK, &Response{
		Status: status,
		Data:   data,
		Msg:    e.GetMsg(status),
	})
}

func RespError(ctx *gin.Context, data interface{}, code ...int) {
	status := e.ERROR
	if code != nil {
		status = code[0]
	}
	ctx.JSON(http.StatusBadRequest, &Response{
		Status: status,
		Data:   data,
		Msg:    e.GetMsg(status),
	})
}

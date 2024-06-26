/*
HTTP请求工具

BetaX Harbor
Copyright © 2024 SkyeZhang <skai-zhang@hotmail.com>
*/

package util

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/skye-z/olympus/model"
)

func ReturnError(ctx *gin.Context, err CustomError) {
	ctx.JSON(200, err)
	ctx.Abort()
}

type commonResponse struct {
	State   bool   `json:"state"`
	Message string `json:"message"`
	Data    any    `json:"data"`
	Time    int64  `json:"time"`
}

func ReturnMessage(ctx *gin.Context, state bool, message string) {
	ctx.JSON(200, commonResponse{
		State:   state,
		Message: message,
		Time:    time.Now().Unix(),
	})
	ctx.Abort()
}

func ReturnData(ctx *gin.Context, state bool, obj any) {
	ctx.JSON(200, commonResponse{
		State: state,
		Data:  obj,
		Time:  time.Now().Unix(),
	})
	ctx.Abort()
}

func ReturnMessageData(ctx *gin.Context, state bool, message string, obj any) {
	ctx.JSON(200, commonResponse{
		State:   state,
		Message: message,
		Data:    obj,
		Time:    time.Now().Unix(),
	})
	ctx.Abort()
}

// 校验权限
func CheckAuth(ctx *gin.Context) bool {
	param, exists := ctx.Get("user")
	if exists {
		user := param.(model.User)
		return user.Admin
	} else {
		return false
	}
}

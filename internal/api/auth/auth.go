package auth

import (
	"Excel-Props/internal/api"
	"Excel-Props/pkg/config"
	"Excel-Props/pkg/constant"
	"Excel-Props/pkg/db"
	"Excel-Props/pkg/log"
	"Excel-Props/pkg/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	req := api.LoginReq{}
	resp := api.LoginResp{}

	if api.IsInterruptBindJson(&req, &resp.CommResp, c) {
		return
	}
	user, err := db.DB.MysqlDB.GetAccountInfo(req.UserID)
	if err != nil {
		log.NewError(operationID, "not user info", err.Error(), req)
		resp.ErrCode = constant.NotUserInfo
		resp.ErrMsg = "not user info"
		return
	}
	if user.Password != req.Password {
		log.NewError(operationID, "password err", req)
		resp.ErrCode = constant.PasswordError
		resp.ErrMsg = "password err"
		return
	}
	tokenString, _ := token.CreateToken(req.UserID, config.Config.TokenPolicy.AccessExpire)
	resp.Data.UserName = user.UserName
	resp.Data.Token = tokenString
	c.JSON(http.StatusOK, resp)
}
func ParseToken(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	resp := api.PareTokenResp{}
	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error())
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		return
	}
	_, err = db.DB.MysqlDB.GetAccountInfo(userID)
	if err != nil {
		log.NewError(operationID, "not user info", err.Error())
		resp.ErrCode = constant.NotUserInfo
		resp.ErrMsg = "not user info"
		return
	}
	c.JSON(http.StatusOK, resp)
}

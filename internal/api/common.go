package api

import (
	"Excel-Props/pkg/constant"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsInterruptBindJson(req interface{}, commResp *CommResp, c *gin.Context) bool {
	if err := c.BindJSON(req); err != nil {
		commResp.ErrCode = constant.ArgsErr
		commResp.ErrMsg = err.Error()
		c.JSON(http.StatusOK, commResp)
		return true
	}
	return false
}

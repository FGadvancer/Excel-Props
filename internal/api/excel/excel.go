package excel

import (
	"Excel-Props/internal/api"
	"github.com/gin-gonic/gin"
	"net/http"
)

func FileUpload(c *gin.Context) {
	//operationID := c.Request.Header.Get("operationID")
	//token := c.Request.Header.Get("token")
	//req := api.ExcelFileUploadReq{}
	resp := api.ExcelFileUploadResp{}
	//
	//
	//if api.IsInterruptBindJson(&req, &resp.CommResp, c) {
	//	return
	//}
	//
	//var ok bool
	//var errInfo string
	//ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	//if !ok {
	//	errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//
	//log.NewInfo(req.OperationID, "ForceLogout args ", req.String())
	//etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	//if etcdConn == nil {
	//	errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//client := rpc.NewAuthClient(etcdConn)
	//reply, err := client.ForceLogout(context.Background(), req)
	//if err != nil {
	//	errMsg := req.OperationID + " UserToken failed " + err.Error() + req.String()
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//resp := api.ForceLogoutResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	//log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " return ", resp)
	c.JSON(http.StatusOK, resp)
}
func GetAllExcelFiles(c *gin.Context) {
	//operationID := c.Request.Header.Get("operationID")
	//token := c.Request.Header.Get("token")
	//req := api.GetAllExcelFilesReq{}
	resp := api.GetAllExcelFilesResp{}
	//if api.IsInterruptBindJson(&req, &resp.CommResp, c) {
	//	return
	//}
	//var ok bool
	//var errInfo string
	//ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	//if !ok {
	//	errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//
	//log.NewInfo(req.OperationID, "ForceLogout args ", req.String())
	//etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	//if etcdConn == nil {
	//	errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//client := rpc.NewAuthClient(etcdConn)
	//reply, err := client.ForceLogout(context.Background(), req)
	//if err != nil {
	//	errMsg := req.OperationID + " UserToken failed " + err.Error() + req.String()
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//resp := api.ForceLogoutResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	//log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " return ", resp)
	c.JSON(http.StatusOK, resp)
}
func GetOneExcelDetail(c *gin.Context) {
	//operationID := c.Request.Header.Get("operationID")
	//token := c.Request.Header.Get("token")
	//req := api.GetOneExcelDetailReq{}
	resp := api.GetOneExcelDetailResp{}
	//if api.IsInterruptBindJson(&req, &resp.CommResp, c) {
	//	return
	//}
	//
	//var ok bool
	//var errInfo string
	//ok, req.OpUserID, errInfo = token_verify.GetUserIDFromToken(c.Request.Header.Get("token"), req.OperationID)
	//if !ok {
	//	errMsg := req.OperationID + " " + "GetUserIDFromToken failed " + errInfo + " token:" + c.Request.Header.Get("token")
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//
	//log.NewInfo(req.OperationID, "ForceLogout args ", req.String())
	//etcdConn := getcdv3.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImAuthName, req.OperationID)
	//if etcdConn == nil {
	//	errMsg := req.OperationID + " getcdv3.GetDefaultConn == nil"
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//client := rpc.NewAuthClient(etcdConn)
	//reply, err := client.ForceLogout(context.Background(), req)
	//if err != nil {
	//	errMsg := req.OperationID + " UserToken failed " + err.Error() + req.String()
	//	log.NewError(req.OperationID, errMsg)
	//	c.JSON(http.StatusInternalServerError, gin.H{"errCode": 500, "errMsg": errMsg})
	//	return
	//}
	//resp := api.ForceLogoutResp{CommResp: api.CommResp{ErrCode: reply.CommonResp.ErrCode, ErrMsg: reply.CommonResp.ErrMsg}}
	//log.NewInfo(params.OperationID, utils.GetSelfFuncName(), " return ", resp)
	c.JSON(http.StatusOK, resp)
}

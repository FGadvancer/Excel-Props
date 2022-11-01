package excel

import (
	"Excel-Props/internal/api"
	"Excel-Props/pkg/constant"
	"Excel-Props/pkg/db"
	"Excel-Props/pkg/log"
	"Excel-Props/pkg/token"
	"Excel-Props/pkg/utils"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func FileUpload(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.ExcelFileUploadReq{}
	resp := api.ExcelFileUploadResp{}
	//
	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error())
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		c.JSON(http.StatusOK, resp)
		return
	}
	if api.IsInterruptBindJson(&req, &resp.CommResp, c) {
		return
	}
	log.NewDebug(operationID, "input args is:", req)
	temp, err := db.DB.MysqlDB.GetTemplateInfo(req.SheetID)
	if err != nil {
		log.NewError(operationID, "not template info", err.Error(), req)
		resp.ErrCode = constant.NotTemplateInfo
		resp.ErrMsg = "not template info"
		c.JSON(http.StatusOK, resp)
		return
	}
	var tempMaterialList []*db.SheetAndMaterial
	for _, v := range req.SheetList {
		temp, err := db.DB.MysqlDB.GetMaterialInfo(v.MaterialKey, v.MaterialStandard)
		if err != nil {
			log.NewError(operationID, "not material info", err.Error(), req)
			resp.ErrCode = constant.NotMaterialInfo
			resp.ErrMsg = "not material info"
			c.JSON(http.StatusOK, resp)
			return
		}
		material := new(db.SheetAndMaterial)
		material.SheetID = req.SheetID
		material.MaterialKey = v.MaterialKey
		material.MaterialStandard = v.MaterialStandard
		material.MaterialCategory = temp.MaterialCategory
		material.MaterialName = temp.MaterialName
		material.MaterialSubstance = temp.MaterialSubstance
		material.Quantity = v.Quantity
		material.MaterialUnit = temp.MaterialUnit
		material.ProcessingCategory = temp.ProcessingCategory
		material.RemarkOne = temp.RemarkOne
		material.RemarkTwo = temp.RemarkTwo
		material.IsPurchase = temp.IsPurchase
		material.StandardCraft = temp.StandardCraft
		material.LastModifyTime = time.Now()
		material.LastModifierUserID = userID
		material.LastModifyCount = v.Quantity
		material.SubMaterialKey = utils.StructToJsonString(v.SubMaterialKey)
		tempMaterialList = append(tempMaterialList, material)
	}
	//抢占分布式锁
	err = db.DB.Redis.LockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "this sheetID locked by others ", err.Error(), req)
		resp.ErrCode = constant.SheetBusy
		resp.ErrMsg = "this sheetID locked by others"
		c.JSON(http.StatusOK, resp)
		return
	}

	//分布式锁抢占成功
	sheet, err := db.DB.MysqlDB.GetSheetInfo(req.SheetID)
	if err != nil {
		tx := db.DB.MysqlDB.Db().Begin()
		var sheet db.Sheet
		sheet.SheetID = req.SheetID
		sheet.CommodityName = temp.MachineKind + "_" + temp.ProductName
		sheet.Code = temp.Code
		sheet.CreatorUserID = userID
		sheet.Version = 1
		sheet.LastModifierIP = c.Request.RemoteAddr
		sheet.CreateTime = time.Now()
		sheet.LastModifyTime = time.Now()
		sheet.LastModifierUserID = userID
		err := db.DB.MysqlDB.InsertSheet(&sheet)
		if err != nil {
			log.NewError(operationID, "this sheet db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "this sheet db operation error"
			c.JSON(http.StatusOK, resp)
			return
		}
		err = db.DB.MysqlDB.BatchInsertSheetAndMaterialList(tempMaterialList)
		if err != nil {
			tx.Rollback()
			log.NewError(operationID, "BatchInsertSheetAndMaterialList db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "BatchInsertSheetAndMaterialList err"
			c.JSON(http.StatusOK, resp)
			return
		}
		tx.Commit()

	} else {
		tx := db.DB.MysqlDB.Db().Begin()
		var newSheet db.Sheet
		newSheet.SheetID = sheet.SheetID
		newSheet.Version = sheet.Version + 1
		newSheet.LastModifierUserID = userID
		newSheet.LastModifyTime = time.Now()
		sheet.LastModifierIP = c.Request.RemoteAddr
		err := db.DB.MysqlDB.UpdateSheet(&newSheet)
		if err != nil {
			log.NewError(operationID, "UpdateSheet db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "UpdateSheet err"
			c.JSON(http.StatusOK, resp)
			return
		}
		for _, material := range tempMaterialList {
			oldMaterialInfo, err := db.DB.MysqlDB.GetSheetAndMaterialInfo(material.SheetID, material.MaterialKey, material.MaterialStandard)
			if err != nil {
				newErr := db.DB.MysqlDB.InsertSheetAndMaterial(material)
				if newErr != nil {
					tx.Rollback()
					log.NewError(operationID, "InsertSheetAndMaterial db operation error", newErr.Error(), req)
					resp.ErrCode = constant.SheetDBError
					resp.ErrMsg = "InsertSheetAndMaterial err"
					c.JSON(http.StatusOK, resp)
					return
				}
			} else {
				var a []string
				utils.JsonStringToStruct(oldMaterialInfo.SubMaterialKey, &a)
				var b []string
				utils.JsonStringToStruct(material.SubMaterialKey, &b)
				a = append(a, b...)
				oldMaterialInfo.Quantity = oldMaterialInfo.Quantity + material.Quantity
				oldMaterialInfo.SubMaterialKey = utils.StructToJsonString(utils.RemoveRepeatedStringInList(a))
				oldMaterialInfo.LastModifyCount = material.LastModifyCount
				oldMaterialInfo.LastModifierUserID = material.LastModifierUserID
				oldMaterialInfo.LastModifyTime = material.LastModifyTime
				newErr := db.DB.MysqlDB.UpdateSheetAndMaterial(oldMaterialInfo)
				if newErr != nil {
					tx.Rollback()
					log.NewError(operationID, "UpdateSheetAndMaterial db operation error", newErr.Error(), req)
					resp.ErrCode = constant.SheetDBError
					resp.ErrMsg = "InsertSheetAndMaterial err"
					c.JSON(http.StatusOK, resp)
					return
				}

			}

		}
		tx.Commit()
	}
	//解开分布式锁
	err = db.DB.Redis.UnLockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "unLockSheetID err:", err.Error(), req)
	}
	c.JSON(http.StatusOK, resp)
}
func GetAllExcelFiles(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	resp := api.GetAllExcelFilesResp{}
	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error())
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		c.JSON(http.StatusOK, resp)
		return
	}
	sheetList, err := db.DB.MysqlDB.GetAllSheetsInfo()
	if err != nil {
		log.NewError(operationID, "GetAllSheetsInfo db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "GetAllSheetsInfo err"
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data.SheetList = sheetList
	c.JSON(http.StatusOK, resp)
}
func GetOneExcelDetail(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.GetOneExcelDetailReq{}
	resp := api.GetOneExcelDetailResp{}
	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error(), userID)
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		c.JSON(http.StatusOK, resp)
		return
	}
	if api.IsInterruptBindJson(&req, &resp.CommResp, c) {
		return
	}

	sheet, err := db.DB.MysqlDB.GetSheetInfo(req.SheetID)
	if err != nil {
		log.NewError(operationID, "sheet info not exist", err.Error())
		resp.ErrCode = constant.NotSheetInfo
		resp.ErrMsg = "sheet info not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	sheetAndMaterialList, err := db.DB.MysqlDB.GetSheetAndMaterialInfoBySheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "material infos not exist", err.Error())
		resp.ErrCode = constant.NotSheetInfo
		resp.ErrMsg = "material infos not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data.Sheet = sheet
	resp.Data.SheetMaterialList = sheetAndMaterialList
	c.JSON(http.StatusOK, resp)
}

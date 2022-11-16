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
	log.NewDebug(operationID, "req", req)
	user, err := db.DB.MysqlDB.GetAccountInfo(userID)
	if err != nil {
		log.NewError(operationID, "not user info", err.Error(), req)
		resp.ErrCode = constant.NotUserInfo
		resp.ErrMsg = "not user info"
		c.JSON(http.StatusOK, resp)
		return
	}
	log.NewDebug(operationID, "input args is:", req)
	temp, err := db.DB.MysqlDB.GetTemplateSheetInfo(req.SheetID)
	if err != nil {
		log.NewError(operationID, "not template info", err.Error(), req)
		resp.ErrCode = constant.NotTemplateInfo
		resp.ErrMsg = "not template info"
		c.JSON(http.StatusOK, resp)
		return
	}
	var tempMaterialList []*db.SheetAndMaterial
	var recordList []*db.VersionUpLoadRecord
	for _, v := range req.SheetList {
		temp, err := db.DB.MysqlDB.GetMaterialInfo(v.MaterialKey)
		if err != nil {
			log.NewError(operationID, "not material info", err.Error(), req)
			resp.ErrCode = constant.NotTemPlateMaterialInfo
			resp.ErrMsg = "not material info"
			c.JSON(http.StatusOK, resp)
			return
		}
		record := new(db.VersionUpLoadRecord)
		record.SheetID = req.SheetID
		record.MaterialKey = v.MaterialKey
		record.MaterialStandard = v.MaterialStandard
		record.Version = 1
		record.SubVersion = 1
		record.MaterialCategory = temp.MaterialCategory
		record.MaterialName = temp.MaterialName
		record.MaterialSubstance = temp.MaterialSubstance
		record.Quantity = v.Quantity
		record.MaterialUnit = temp.MaterialUnit
		record.ProcessingCategory = temp.ProcessingCategory
		record.RemarkOne = temp.RemarkOne
		record.RemarkTwo = temp.RemarkTwo
		record.IsPurchase = temp.IsPurchase
		record.StandardCraft = temp.StandardCraft
		record.SubMaterialKey = utils.StructToJsonString(v.SubMaterialKey)
		record.CommitTime = time.Now()
		record.ModifierUserID = userID
		record.ModifierName = user.UserName
		recordList = append(recordList, record)

		material := new(db.SheetAndMaterial)
		material.SheetID = req.SheetID
		material.MaterialKey = v.MaterialKey
		material.MaterialStandard = v.MaterialStandard
		material.Version = 1
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
		material.LastModifierName = user.UserName
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
	//摸具信息存储
	sheet, err := db.DB.MysqlDB.GetSheetInfo(req.SheetID)
	if err != nil {
		tx := db.DB.MysqlDB.Db().Begin()
		var sheet db.Sheet
		sheet.SheetID = req.SheetID
		sheet.CommodityName = temp.MachineKind + "_" + temp.ProductName
		sheet.Code = temp.Code
		sheet.CreatorUserID = userID
		sheet.Version = 1
		sheet.SubVersion = 1
		sheet.IsCompleteVersion = false
		sheet.LastModifierIP = c.Request.RemoteAddr
		sheet.CreateTime = time.Now()
		sheet.LastModifyTime = time.Now()
		sheet.LastModifierUserID = userID
		sheet.LastModifierName = user.UserName
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
		err = db.DB.MysqlDB.BatchInsertVersionUpLoadRecordList(recordList)
		if err != nil {
			tx.Rollback()
			log.NewError(operationID, "BatchInsertVersionUpLoadRecordList db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "BatchInsertVersionUpLoadRecordList err"
			c.JSON(http.StatusOK, resp)
			return
		}
		tx.Commit()

	} else {
		var isUpdateCompleteVersion bool
		tx := db.DB.MysqlDB.Db().Begin()
		var newSheet db.Sheet
		newSheet.SheetID = sheet.SheetID
		if sheet.IsCompleteVersion {
			newSheet.Version = sheet.Version + 1
			isUpdateCompleteVersion = true
		} else {
			newSheet.Version = sheet.Version
		}
		newSheet.SubVersion = sheet.SubVersion + 1
		newSheet.LastModifierUserID = userID
		newSheet.LastModifyTime = time.Now()
		newSheet.LastModifierIP = c.Request.RemoteAddr
		newSheet.LastModifierName = user.UserName
		err := db.DB.MysqlDB.UpdateSheet(&newSheet, isUpdateCompleteVersion)
		if err != nil {
			log.NewError(operationID, "UpdateSheet db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "UpdateSheet err"
			c.JSON(http.StatusOK, resp)
			return
		}
		for _, material := range tempMaterialList {
			oldMaterialInfo, err := db.DB.MysqlDB.GetSheetAndMaterialInfo(material.SheetID, material.MaterialKey, material.MaterialStandard, newSheet.Version)
			if err != nil {
				material.Version = newSheet.Version
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
				oldMaterialInfo.LastModifierName = material.LastModifierName
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
		for _, record := range recordList {
			record.Version = newSheet.Version
			record.SubVersion = newSheet.SubVersion
		}
		err = db.DB.MysqlDB.BatchInsertVersionUpLoadRecordList(recordList)
		if err != nil {
			tx.Rollback()
			log.NewError(operationID, "BatchInsertVersionUpLoadRecordList db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "BatchInsertVersionUpLoadRecordList err"
			c.JSON(http.StatusOK, resp)
			return
		}
		tx.Commit()
	}
	//解开分布式锁
	err = db.DB.Redis.UnLockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "unLockSheetID err:", err.Error(), req)
	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}
func GetAllExcelFiles(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	resp := api.GetAllExcelFilesResp{}
	log.NewDebug(operationID, "req", tokenString)
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
	if err := utils.CopyStructFields(&resp.Data.SheetList, sheetList); err != nil {
		log.NewDebug(operationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewDebug(operationID, "resp", resp)
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
	log.NewDebug(operationID, "req", req)
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
	if err := utils.CopyStructFields(&resp.Data.SheetMaterialList, sheetAndMaterialList); err != nil {
		log.NewDebug(operationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	resp.Data.Sheet = *sheet
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}

func CompleteSheetVersion(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.CompleteSheetVersionReq{}
	resp := api.CompleteSheetVersionResp{}
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
	log.NewDebug(operationID, "req", req)
	//抢占分布式锁
	err = db.DB.Redis.LockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "this sheetID locked by others ", err.Error(), req)
		resp.ErrCode = constant.SheetBusy
		resp.ErrMsg = "this sheetID locked by others"
		c.JSON(http.StatusOK, resp)
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
	if sheet.IsCompleteVersion {
		log.NewError(operationID, "sheet has been completed")
		resp.ErrCode = constant.HasCompleteVersion
		resp.ErrMsg = "sheet has been completed"
		c.JSON(http.StatusOK, resp)
		return
	}
	err = db.DB.MysqlDB.UpdateSheetColumns(sheet.SheetID, map[string]interface{}{"is_complete_version": true, "sub_version": 0})
	if err != nil {
		log.NewError(operationID, "UpdateSheet db operation error", err.Error(), req)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "UpdateSheet err"
		c.JSON(http.StatusOK, resp)
		return
	}
	//解开分布式锁
	err = db.DB.Redis.UnLockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "unLockSheetID err:", err.Error(), req)
	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}

func RevokeSheetVersion(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.RevokeSheetVersionReq{}
	resp := api.RevokeSheetVersionResp{}

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
	log.NewDebug(operationID, "req", req)
	//抢占分布式锁
	err = db.DB.Redis.LockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "this sheetID locked by others ", err.Error(), req)
		resp.ErrCode = constant.SheetBusy
		resp.ErrMsg = "this sheetID locked by others"
		c.JSON(http.StatusOK, resp)
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
	if sheet.IsCompleteVersion {
		log.NewError(operationID, "sheet has been completed")
		resp.ErrCode = constant.HasCompleteVersion
		resp.ErrMsg = "sheet has been completed"
		c.JSON(http.StatusOK, resp)
		return
	}
	if sheet.Version-1 < 0 {
		log.NewError(operationID, "sheet can not be revoked")
		resp.ErrCode = constant.SheetVersionZero
		resp.ErrMsg = "sheet can not be revoked"
		c.JSON(http.StatusOK, resp)
		return
	}
	oldVersion := sheet.Version
	tx := db.DB.MysqlDB.Db().Begin()
	if sheet.Version-1 == 0 {
		err = db.DB.MysqlDB.DeleteSheet(sheet.SheetID)
		if err != nil {
			log.NewError(operationID, "DeleteSheet db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "DeleteSheet err"
			c.JSON(http.StatusOK, resp)
			return
		}
	} else {
		err = db.DB.MysqlDB.UpdateSheetColumns(sheet.SheetID, map[string]interface{}{"is_complete_version": true, "sub_version": 0, "version": sheet.Version - 1})
		if err != nil {
			log.NewError(operationID, "UpdateSheet db operation error", err.Error(), req)
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "UpdateSheet err"
			c.JSON(http.StatusOK, resp)
			return
		}
	}

	err = db.DB.MysqlDB.DeleteSheetAndMaterialInfoBySheetIDAndVersion(req.SheetID, oldVersion)
	if err != nil {
		tx.Rollback()
		log.NewError(operationID, "DeleteSheetAndMaterialInfoBySheetIDAndVersion db operation error", err.Error(), req)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "DeleteSheetAndMaterialInfoBySheetIDAndVersion err"
		c.JSON(http.StatusOK, resp)
		return
	}
	err = db.DB.MysqlDB.DeleteVersionRecordListBySheetIDAndVersion(req.SheetID, oldVersion)
	if err != nil {
		tx.Rollback()
		log.NewError(operationID, "DeleteVersionRecordListBySheetIDAndVersion db operation error", err.Error(), req)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "DeleteVersionRecordListBySheetIDAndVersion err"
		c.JSON(http.StatusOK, resp)
		return
	}
	tx.Commit()
	//解开分布式锁
	err = db.DB.Redis.UnLockSheetID(req.SheetID)
	if err != nil {
		log.NewError(operationID, "unLockSheetID err:", err.Error(), req)
	}
	log.NewDebug(operationID, "resp", resp)

	c.JSON(http.StatusOK, resp)
}
func GetRecordSheetVersion(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.GetRecordSheetVersionReq{}
	resp := api.GetRecordSheetVersionResp{}

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
	log.NewDebug(operationID, "req", req)
	sheet, err := db.DB.MysqlDB.GetSheetInfo(req.SheetID)
	if err != nil {
		log.NewError(operationID, "sheet info not exist", err.Error())
		resp.ErrCode = constant.NotSheetInfo
		resp.ErrMsg = "sheet info not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	recordList, err := db.DB.MysqlDB.GetVersionRecordList(req.SheetID, sheet.Version)
	if err != nil {
		log.NewError(operationID, "material infos not exist", err.Error())
		resp.ErrCode = constant.NotSheetInfo
		resp.ErrMsg = "material infos not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	var result []*api.AllRecordList

	var temp api.AllRecordList
	for i := 0; i < len(recordList); i++ {

		if i == 0 || recordList[i].SubVersion == recordList[i-1].SubVersion {
			temp.CommitTime = recordList[i].CommitTime
			temp.SubVersion = recordList[i].SubVersion
			temp.ModifierUserID = recordList[i].ModifierUserID
			temp.ModifierName = recordList[i].ModifierName
			temp.RecordList = append(temp.RecordList, recordList[i])
		} else {
			args := new(api.AllRecordList)
			args.RecordList = temp.RecordList
			args.CommitTime = temp.CommitTime
			args.SubVersion = temp.SubVersion
			args.ModifierUserID = temp.ModifierUserID
			args.ModifierName = temp.ModifierName
			result = append(result, args)
			temp.CommitTime = recordList[i].CommitTime
			temp.SubVersion = recordList[i].SubVersion
			temp.ModifierUserID = recordList[i].ModifierUserID
			temp.ModifierName = recordList[i].ModifierName
			temp.RecordList = nil
			temp.RecordList = append(temp.RecordList, recordList[i])
		}
	}

	result = append(result, &temp)
	if err := utils.CopyStructFields(&resp.Data.VersionUpLoadRecordList, result); err != nil {
		log.NewDebug(operationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	resp.Data.Sheet = *sheet
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}
func RevokeRecordSheetVersion(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.RevokeRecordSheetVersionReq{}
	resp := api.RevokeRecordSheetVersionResp{}

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
	log.NewDebug(operationID, "req", req)
	sheet, err := db.DB.MysqlDB.GetSheetInfo(req.SheetID)
	if err != nil {
		log.NewError(operationID, "sheet info not exist", err.Error())
		resp.ErrCode = constant.NotSheetInfo
		resp.ErrMsg = "sheet info not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	recordList, err := db.DB.MysqlDB.GetSubVersionRecordList(req.SheetID, sheet.Version, req.SubVersion)
	if err != nil {
		log.NewError(operationID, "record infos not exist", err.Error())
		resp.ErrCode = constant.NotRecordInfo
		resp.ErrMsg = "record infos not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	var result []*api.AllRecordList

	var temp api.AllRecordList
	for i := 0; i < len(recordList); i++ {

		if i == 0 || recordList[i].SubVersion == recordList[i-1].SubVersion {
			temp.CommitTime = recordList[i].CommitTime
			temp.SubVersion = recordList[i].SubVersion
			temp.ModifierUserID = recordList[i].ModifierUserID
			temp.ModifierName = recordList[i].ModifierName
			temp.RecordList = append(temp.RecordList, recordList[i])
		} else {
			args := new(api.AllRecordList)
			args.RecordList = temp.RecordList
			args.CommitTime = temp.CommitTime
			args.SubVersion = temp.SubVersion
			args.ModifierUserID = temp.ModifierUserID
			args.ModifierName = temp.ModifierName
			result = append(result, args)
			temp.CommitTime = recordList[i].CommitTime
			temp.SubVersion = recordList[i].SubVersion
			temp.ModifierUserID = recordList[i].ModifierUserID
			temp.ModifierName = recordList[i].ModifierName
			temp.RecordList = nil
			temp.RecordList = append(temp.RecordList, recordList[i])
		}
	}

	result = append(result, &temp)
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}
func GetSubSheetList(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")

	resp := api.GetSubSheetListResp{}

	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error(), userID)
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		c.JSON(http.StatusOK, resp)
		return
	}

	list, err := db.DB.MysqlDB.GetSheetSubList()
	if err != nil {
		log.NewError(operationID, "GetSheetSubList not exist", err.Error())
		resp.ErrCode = constant.NotRecordInfo
		resp.ErrMsg = "record infos not exist"
		c.JSON(http.StatusOK, resp)
		return
	}
	resp.Data.SubSheetIDList = list
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}
func ModifySubSheetList(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req := api.AddSubSheetListReq{}
	resp := api.AddSubSheetListResp{}

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
	log.NewDebug(operationID, "req", req)
	switch req.OperationType {
	case constant.AddSubSheet:
		var list []*db.SheetSub
		for _, v := range req.SubSheetIDList {
			_, newErr := db.DB.MysqlDB.GetSheetSubInfo(v)
			if newErr != nil {
				temp := new(db.SheetSub)
				temp.SubSheetID = v
				list = append(list, temp)
			}
		}
		err = db.DB.MysqlDB.BatchInsertSheetSubList(list)
		if err != nil {
			log.NewError(operationID, "BatchInsertSheetSubList", err.Error())
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "BatchInsertSheetSubList err" + err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}
	case constant.DeleteSubSheet:
		err = db.DB.MysqlDB.BatchDeleteSheetSubList(req.SubSheetIDList)
		if err != nil {
			log.NewError(operationID, "BatchDeleteSheetSubList", err.Error())
			resp.ErrCode = constant.SheetDBError
			resp.ErrMsg = "BatchDeleteSheetSubList err" + err.Error()
			c.JSON(http.StatusOK, resp)
			return
		}

	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}

func GetTemplateSheetList(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	resp := api.GetTemplateSheetListResp{}
	log.NewDebug(operationID, "req", tokenString)
	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error())
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		c.JSON(http.StatusOK, resp)
		return
	}
	templateSheetList, err := db.DB.MysqlDB.GetAllSheetTemplates()
	if err != nil {
		log.NewError(operationID, "GetAllSheetTemplates db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "GetAllSheetTemplates err"
		c.JSON(http.StatusOK, resp)
		return
	}
	if err := utils.CopyStructFields(&resp.Data.TemplateSheetList, templateSheetList); err != nil {
		log.NewDebug(operationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}

func UpdateTemplateSheetList(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req:=api.UpdateTemplateSheetListReq{}
	resp := api.UpdateTemplateSheetListResp{}
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
	log.NewDebug(operationID, "req", req)
	err = db.DB.MysqlDB.DeleteAllTemplateSheet()
	if err != nil {
		log.NewError(operationID, "DeleteAllTemplateSheet db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "DeleteAllTemplateSheet err"
		c.JSON(http.StatusOK, resp)
		return
	}
	err = db.DB.MysqlDB.ImportDataToTemplateSheet(req.TemplateSheetList)
	if err != nil {
		log.NewError(operationID, "ImportDataToTemplateSheet db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "ImportDataToTemplateSheet err"
		c.JSON(http.StatusOK, resp)
		return
	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}

func GetTemplateMaterialList(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	resp := api.GetTemplateMaterialListResp{}
	log.NewDebug(operationID, "req", tokenString)
	userID, err := token.GetUserIDFromToken(tokenString)
	if err != nil {
		log.NewError(operationID, "token parse failed", err.Error())
		resp.ErrCode = constant.ParseTokenFailed
		resp.ErrMsg = "token parse failed"
		c.JSON(http.StatusOK, resp)
		return
	}
	templateMaterialList, err := db.DB.MysqlDB.GetAllMaterialTemplates()
	if err != nil {
		log.NewError(operationID, "GetAllMaterialTemplates db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "GetAllMaterialTemplates err"
		c.JSON(http.StatusOK, resp)
		return
	}
	if err := utils.CopyStructFields(&resp.Data.TemplateMaterialList, templateMaterialList); err != nil {
		log.NewDebug(operationID, utils.GetSelfFuncName(), "CopyStructFields failed", err.Error())
	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}

func UpdateTemplateMaterialList(c *gin.Context) {
	operationID := c.Request.Header.Get("operationID")
	tokenString := c.Request.Header.Get("token")
	req:=api.UpdateTemplateMaterialListReq{}
	resp := api.UpdateTemplateMaterialListResp{}
	log.NewDebug(operationID, "req", tokenString)
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
	log.NewDebug(operationID, "req", req)
	err = db.DB.MysqlDB.DeleteAllTemplateMaterial()
	if err != nil {
		log.NewError(operationID, "DeleteAllTemplateMaterial db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "DeleteAllTemplateMaterial err"
		c.JSON(http.StatusOK, resp)
		return
	}
	err = db.DB.MysqlDB.ImportDataToTemplateMaterial(req.TemplateMaterialList)
	if err != nil {
		log.NewError(operationID, "ImportDataToTemplateMaterial db operation error", err.Error(), userID)
		resp.ErrCode = constant.SheetDBError
		resp.ErrMsg = "ImportDataToTemplateMaterial err"
		c.JSON(http.StatusOK, resp)
		return
	}
	log.NewDebug(operationID, "resp", resp)
	c.JSON(http.StatusOK, resp)
}
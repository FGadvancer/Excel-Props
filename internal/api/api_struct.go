package api

import "Excel-Props/pkg/db"

type (
	LoginReq struct {
		UserID   string `json:"userID" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResp struct {
		CommResp
		Data struct {
			Token    string `json:"token,omitempty"`
			UserName string `json:"userName,omitempty"`
		} `json:"data,omitempty"`
	}
	PareTokenResp struct {
		CommResp
		Data struct {
			UserName string `json:"userName,omitempty"`
		} `json:"data,omitempty"`
	}
	ExcelFileUploadReq struct {
		SheetID   string        `json:"sheetID" binding:"required"`
		SheetList []SheetObject `json:"sheetList"`
	}
	ExcelFileUploadResp struct {
		CommResp
	}
	GetAllExcelFilesReq struct {
	}
	GetAllExcelFilesResp struct {
		CommResp
		Data struct {
			SheetList []*db.Sheet `json:"sheetList"`
		} `json:"data,omitempty"`
	}
	GetOneExcelDetailReq struct {
		SheetID string `json:"sheetID" binding:"required"`
	}
	GetOneExcelDetailResp struct {
		CommResp
		Data struct {
			*db.Sheet
			SheetMaterialList []*db.SheetAndMaterial `json:"sheetMaterialList"`
		} `json:"data,omitempty"`
	}
	CompleteSheetVersionReq struct {
		SheetID string `json:"sheetID" binding:"required"`
	}
	CompleteSheetVersionResp struct {
		CommResp
	}
	RevokeSheetVersionReq struct {
		SheetID string `json:"sheetID" binding:"required"`
	}
	RevokeSheetVersionResp struct {
		CommResp
	}
	GetRecordSheetVersionReq struct {
		SheetID string `json:"sheetID" binding:"required"`
	}
	GetRecordSheetVersionResp struct {
		CommResp
		Data struct {
			*db.Sheet
			VersionUpLoadRecordList []*db.VersionUpLoadRecord `json:"versionUpLoadRecord"`
		} `json:"data,omitempty"`
	}
)
type CommResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}
type SheetObject struct {
	MaterialKey      string   `json:"materialKey" `
	SubMaterialKey   []string `json:"subMaterialKey" `
	MaterialStandard string   `json:"materialStandard" `
	Quantity         int32    `json:"quantity" `
}

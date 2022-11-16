package api

import (
	"Excel-Props/pkg/db"
	"time"
)

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
			ManagerLevel int `json:"managerLevel"`
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
			SheetList []db.Sheet `json:"sheetList"`
		} `json:"data,omitempty"`
	}
	GetOneExcelDetailReq struct {
		SheetID string `json:"sheetID" binding:"required"`
	}
	GetOneExcelDetailResp struct {
		CommResp
		Data struct {
			db.Sheet
			SheetMaterialList []db.SheetAndMaterial `json:"sheetMaterialList"`
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
			db.Sheet
			VersionUpLoadRecordList []AllRecordList `json:"versionUpLoadRecord"`
		} `json:"data,omitempty"`
	}
	RevokeRecordSheetVersionReq struct {
		SheetID    string `json:"sheetID" binding:"required"`
		SubVersion int32  ` json:"subVersion" binding:"required"`
	}
	RevokeRecordSheetVersionResp struct {
		CommResp
	}

	GetSubSheetListResp struct {
		CommResp
		Data struct {
			SubSheetIDList []string `json:"subSheetIDList"`
		} `json:"data,omitempty"`
	}
	AddSubSheetListReq struct {
		SubSheetIDList []string `json:"subSheetIDList" binding:"required"`
		OperationType  int32    `json:"operationType" binding:"required"`
	}
	AddSubSheetListResp struct {
		CommResp
	}
	GetTemplateSheetListResp struct {
		CommResp
		Data struct {
			TemplateSheetList []db.TemplateSheet `json:"templateSheetList"`
		} `json:"data,omitempty"`
	}
	GetTemplateMaterialListResp struct {
		CommResp
		Data struct {
			TemplateMaterialList []db.TemplateMaterial `json:"templateMaterialList"`
		} `json:"data,omitempty"`
	}
	UpdateTemplateSheetListReq struct {
		TemplateSheetList []*db.TemplateSheet `json:"templateSheetList"`
	}
	UpdateTemplateSheetListResp struct {
		CommResp
	}
	UpdateTemplateMaterialListReq struct {
		TemplateMaterialList []*db.TemplateMaterial `json:"templateMaterialList"`
	}
	UpdateTemplateMaterialListResp struct {
		CommResp
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
type AllRecordList struct {
	SubVersion     int32                     ` json:"subVersion"`
	CommitTime     time.Time                 ` json:"commitTime"`
	ModifierUserID string                    `json:"modifierUserID"`
	ModifierName   string                    ` json:"modifierName"`
	RecordList     []*db.VersionUpLoadRecord `json:"recordList"`
}

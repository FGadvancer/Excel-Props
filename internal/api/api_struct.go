package api

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
	}
	ExcelFileUploadReq struct {
	}
	ExcelFileUploadResp struct {
		CommResp
	}
	GetAllExcelFilesReq struct {
	}
	GetAllExcelFilesResp struct {
		CommResp
	}
	GetOneExcelDetailReq struct {
	}
	GetOneExcelDetailResp struct {
		CommResp
	}
)
type CommResp struct {
	ErrCode int32  `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
}

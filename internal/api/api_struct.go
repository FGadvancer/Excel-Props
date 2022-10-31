package api

type (
	LoginReq struct {
		UserID   string `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResp struct {
		CommResp
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

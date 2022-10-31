package api

type (
	LoginReq struct {
		UserID   string `json:"userId" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	LoginResp struct {
	}
	PareTokenResp struct {
	}
	ExcelFileUploadReq struct {
	}
	ExcelFileUploadResp struct {
	}
	GetAllExcelFilesReq struct {
	}
	GetAllExcelFilesResp struct {
	}
	GetOneExcelDetailReq struct {
	}
	GetOneExcelDetailResp struct {
	}
)

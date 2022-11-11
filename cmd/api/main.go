package main

import (
	"Excel-Props/internal/api/auth"
	"Excel-Props/internal/api/excel"
	"Excel-Props/pkg/config"
	"Excel-Props/pkg/constant"
	"Excel-Props/pkg/log"
	"Excel-Props/pkg/utils"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	"io"
	"os"
	"strconv"
)

// @title Excel-Props API
// @version 1.0
// @description  Excel-Props 的API服务器文档, 文档中所有请求都有一个operationID字段用于链路追踪

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @BasePath /
func main() {
	log.NewPrivateLog(constant.LogFile)
	gin.SetMode(gin.ReleaseMode)
	f, _ := os.Create("../logs/api.log")
	gin.DefaultWriter = io.MultiWriter(f)
	//	gin.SetMode(gin.DebugMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(utils.CorsHandler())
	log.Info("load config: ", config.Config)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// auth routing group, which handles user registration and login services
	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/login", auth.Login)
		authRouterGroup.POST("/parse_token", auth.ParseToken)
	}
	fileRouterGroup := r.Group("/file")
	{
		fileRouterGroup.POST("/excel_files_upload", excel.FileUpload) //1
		//fileRouterGroup.POST("/get_generated_excel_files", excel.SetGlobalRecvMessageOpt)
		fileRouterGroup.POST("/get_all_generated_excel_list", excel.GetAllExcelFiles)
		fileRouterGroup.POST("/get_excel_detail", excel.GetOneExcelDetail)
		fileRouterGroup.POST("/complete_sheet_version", excel.CompleteSheetVersion)
		fileRouterGroup.POST("/revoke_sheet_version", excel.RevokeSheetVersion)
		fileRouterGroup.POST("/get_record_sheet_version", excel.GetRecordSheetVersion)
		fileRouterGroup.POST("/revoke_record_sheet_version", excel.RevokeRecordSheetVersion)
		fileRouterGroup.POST("/get_sub_sheet_list", excel.GetSubSheetList)
		fileRouterGroup.POST("/add_sub_sheet_list", excel.AddSubSheetList)

	}

	defaultPorts := config.Config.Api.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	fmt.Println("start api server, address: ", address, "Excel-Props version: ", constant.CurrentVersion, "\n")
	err := r.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		panic("api start failed " + err.Error())
	}
}

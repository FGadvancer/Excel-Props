package main

import (
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

	// user routing group, which handles user registration and login services
	userRouterGroup := r.Group("/auth")
	{
		userRouterGroup.POST("/login", user.UpdateUserInfo) //1
		userRouterGroup.POST("/parse_token", user.SetGlobalRecvMessageOpt)
	}
	fileRouterGroup := r.Group("/file")
	{
		userRouterGroup.POST("/excel_files_upload", user.UpdateUserInfo) //1
		userRouterGroup.POST("/get_generated_excel_files", user.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_all_generated_excel_list", user.SetGlobalRecvMessageOpt)
		userRouterGroup.POST("/get_excel_detail", user.SetGlobalRecvMessageOpt)
	}

	defaultPorts := config.Config.Api.GinPort
	ginPort := flag.Int("port", defaultPorts[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	fmt.Println("start api server, address: ", address, "OpenIM version: ", constant.CurrentVersion, "\n")
	err := r.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		panic("api start failed " + err.Error())
	}
}

/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"bank-ocr/controller"
	"fmt"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var logger *log.Logger

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// default allow all origins
	r.Use(cors.Default())

	//
	r.LoadHTMLGlob("./app/views")
	r.Static("/assets", "./app/assets")

	r.GET("/", controller.IndexV2)
	r.GET("/status", controller.StatusV2)
	r.POST("/base64", controller.Base64V2)
	r.POST("/file", controller.FileUploadV2)

	logger = log.New(os.Stdout, fmt.Sprintf("[%s] ", "ocrserver"), 0)
	//r.Use(&filters.LogFilter{Logger: logger})

	return r
}

/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:44
 * @Desc: application main router
 */
package router

import (
	"bank-ocr/controller"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// default allow all origins
	r.Use(cors.Default())

	// static
	r.LoadHTMLGlob("./app/views/*")
	r.Static("/assets", "./app/assets")

	// APIs
	r.GET("/", controller.IndexV2)
	r.GET("/status", controller.StatusV2)
	r.POST("/base64", controller.Base64V2)
	r.POST("/file", controller.FileUploadV2)

	return r
}

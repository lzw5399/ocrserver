/**
 * @Author: lzw5399
 * @Date: 2020/9/30 14:25
 * @Desc: home page controller
 */
package controller

import (
	"net/http"

	"bank-ocr/global"
	"bank-ocr/global/response"
	"bank-ocr/service"

	"github.com/gin-gonic/gin"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"AppName": global.BANK_CONFIG.App.Name,
	})
}

func Info(c *gin.Context) {
	info, err := service.GetTesseractInfo()
	if err != nil {
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	response.OkWithData(c, info)
}

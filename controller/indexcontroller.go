/**
 * @Author: lzw5399
 * @Date: 2020/9/30 14:25
 * @Desc: index page controller
 */
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
)

func Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"AppName": "bank-ocr",
	})
}

func Info(c *gin.Context) {
	langs, err := gosseract.GetAvailableLanguages()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err,
		})
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": "Hello!",
		"version": version,
		"tesseract": gin.H{
			"version":   client.Version(),
			"languages": langs,
		},
	})
}

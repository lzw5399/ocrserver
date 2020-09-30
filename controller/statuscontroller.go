/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:49
 * @Desc: status controller
 */
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
	"github.com/otiai10/marmoset"
)

func StatusV2(c *gin.Context) {
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
		"tesseract": marmoset.P{
			"version":   client.Version(),
			"languages": langs,
		},
	})
}

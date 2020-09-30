/**
 * @Author: lzw5399
 * @Date: 2020/9/30 14:25
 * @Desc: index page controller
 */
package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func IndexV2(c *gin.Context){
	c.HTML(http.StatusOK, "index.html", gin.H{
		"AppName": "bank-ocr",
	})
}

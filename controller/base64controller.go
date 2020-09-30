/**
 * @Author: lzw5399
 * @Date: 2020/9/30 13:54
 * @Desc: base64 controller
 */
package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
)

func Base64V2(c *gin.Context) {
	var body = new(struct {
		Base64    string `json:"base64"`
		Trim      string `json:"trim"`
		Languages string `json:"languages"`
		Whitelist string `json:"whitelist"`
	})

	err := json.NewDecoder(c.Request.Body).Decode(body)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	if len(body.Base64) == 0 {
		c.JSON(http.StatusBadRequest, fmt.Errorf("base64 string required"))
		return
	}
	body.Base64 = regexp.MustCompile("data:image\\/png;base64,").ReplaceAllString(body.Base64, "")
	b, err := base64.StdEncoding.DecodeString(body.Base64)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	tempfile.Write(b)

	client := gosseract.NewClient()
	defer client.Close()

	client.Languages = []string{"eng"}
	if body.Languages != "" {
		client.Languages = strings.Split(body.Languages, ",")
	}
	client.SetImage(tempfile.Name())
	if body.Whitelist != "" {
		client.SetWhitelist(body.Whitelist)
	}

	text, err := client.Text()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result":  strings.Trim(text, body.Trim),
		"version": version,
	})
}

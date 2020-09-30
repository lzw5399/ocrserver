/**
 * @Author: lzw5399
 * @Date: 2020/9/30 23:24
 * @Desc: ocr related functionality
 */
package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
)

var (
	imgExp  = regexp.MustCompile("^image")
	version = 0.2
)

func FileUpload(c *gin.Context) {
	// Get uploaded file
	_ = c.Request.ParseMultipartForm(32 << 20)
	upload, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer upload.Close()

	// Create physical file
	tempFile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
	}()

	// Make uploaded physical
	if _, err = io.Copy(tempFile, upload); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempFile.Name())
	client.Languages = []string{"eng"}

	if langs := c.Request.FormValue("languages"); langs != "" {
		client.Languages = strings.Split(langs, ",")
	}

	if whitelist := c.Request.FormValue("whitelist"); whitelist != "" {
		client.SetWhitelist(whitelist)
	}

	var out string
	escapeHtml := true

	switch c.Request.FormValue("format") {
	case "hocr":
		out, err = client.HOCRText()
		escapeHtml = false

	default:
		out, err = client.Text()
	}

	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	if escapeHtml {
		c.JSON(http.StatusOK, gin.H{
			"result":  strings.Trim(out, c.Request.FormValue("trim")),
			"version": version,
		})
	} else {
		c.PureJSON(http.StatusOK, gin.H{
			"result":  strings.Trim(out, c.Request.FormValue("trim")),
			"version": version,
		})
	}
}

func Base64(c *gin.Context) {
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

	tempFile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	defer func() {
		tempFile.Close()
		os.Remove(tempFile.Name())
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
	tempFile.Write(b)

	client := gosseract.NewClient()
	defer client.Close()

	client.Languages = []string{"eng"}
	if body.Languages != "" {
		client.Languages = strings.Split(body.Languages, ",")
	}
	client.SetImage(tempFile.Name())
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

/**
 * @Author: lzw5399
 * @Date: 2020/9/30 14:08
 * @Desc: file controller
 */
package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/otiai10/gosseract/v2"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
)

var (
	imgexp = regexp.MustCompile("^image")
)

func FileUploadV2(c *gin.Context) {
	// Get uploaded file
	c.Request.ParseMultipartForm(32 << 20)
	upload, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	defer upload.Close()

	// Create physical file
	tempfile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	defer func() {
		tempfile.Close()
		os.Remove(tempfile.Name())
	}()

	// Make uploaded physical
	if _, err = io.Copy(tempfile, upload); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.SetImage(tempfile.Name())
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

/**
 * @Author: lzw5399
 * @Date: 2020/10/25 19:51
 * @Desc:
 */
package service

import (
	"encoding/base64"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"bank-ocr/global"
	"bank-ocr/util"

	"github.com/satori/go.uuid"
)

var supportBase64Type = [5]string{
	"data:image/png;base64,",
	"data:image/jpeg;base64,",
	"data:image/gif;base64,",
	"data:image/tiff;base64,",
	"data:application/pdf;base64,",
}

func EnsureContentType(str string) (base64 string, isPdf bool, err error) {
	for i, v := range supportBase64Type {
		if strings.HasPrefix(str, v) {
			base64 = str[len(v):]
			if i == 4 {
				isPdf = true
			}
			return
		}
	}
	err = errors.New("invalid or unsupported content type")
	return
}

// 处理pdf
func PdfToImgsThenGetBytes(base64 string) ([][]byte, error) {
	// pdf先保存到本地
	filePath, err := save(base64)
	if err != nil {
		return nil, err
	}
	defer func() {
		os.Remove(filePath)
	}()

	// pdf分页转成png
	dirToSave, _ := os.Getwd()
	imgs, err := util.PdfToImgs(filePath, dirToSave)
	global.BANK_LOGGER.Info("util.PdfToImgs=", imgs)
	defer func() {
		for _, path := range imgs {
			os.Remove(path)
		}
	}()

	// 读取png成[]byte
	var finalArray [][]byte
	for _, imgPath := range imgs {
		byteArray, err := ioutil.ReadFile(imgPath)
		global.BANK_LOGGER.Info("ioutil.ReadFile=", byteArray)
		if err != nil {
			return nil, err
		}
		finalArray = append(finalArray, byteArray)
	}

	return finalArray, nil
}

func save(base64Str string) (filePath string, err error) {
	buf, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		return
	}

	filePath = uuid.NewV4().String() + ".pdf"
	file, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer file.Close()

	if _, err = file.Write(buf); err != nil {
		return
	}

	err = file.Sync()
	return
}

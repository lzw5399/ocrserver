/**
 * @Author: lzw5399
 * @Date: 2020/10/25 22:01
 * @Desc:
 */
package util

import (
	"bank-ocr/global"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"gopkg.in/gographics/imagick.v2/imagick"
)

func PdfToImgs(filePath, dirToSave string) (imgs []string, err error) {
	if !PathExists(filePath){
		global.BANK_LOGGER.Info("不存在filePath")
		err = errors.New("filePath doesn't exist")
		return
	}

	if !PathExists(dirToSave){
		global.BANK_LOGGER.Info("不存在dirToSave")
		err = errors.New("dirToSave doesn't exist")
		return
	}

	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	// count page of pdf
	if err = mw.PingImage(filePath); err != nil {
		return
	}
	pdfPages := mw.GetNumberImages()

	fileName := strings.TrimSuffix(filepath.Base(filePath), filepath.Ext(filePath))

	// pdf按页转换成png
	for i := uint(0); i < pdfPages; i++ {
		pageName := fmt.Sprintf("%s[%v]", filePath, i)
		imgName := fmt.Sprintf("%s-%v.png", fileName, i)
		imgPath := filepath.Join(dirToSave, imgName)

		// clear resources associated with the wand
		mw.Clear()

		err = mw.ReadImage(pageName)
		if err != nil {
			return
		}

		mwc := mw.Clone()
		err = mwc.WriteImage(imgPath)
		if err != nil {
			return
		}

		imgs = append(imgs, imgName)
	}

	return
}

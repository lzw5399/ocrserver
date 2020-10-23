/**
 * @Author: lzw5399
 * @Date: 2020/9/30 23:42
 * @Desc: ocr core service
 */
package service

import (
	"bytes"
	"errors"
	"golang.org/x/image/tiff"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"strings"

	"bank-ocr/global"
	"bank-ocr/model/request"
	"bank-ocr/model/response"

	"github.com/otiai10/gosseract/v2"
)

// 从 image.Image 获取 string
func GetTextFromImage(img image.Image, contentType string, re request.OcrBase) (string, error) {
	buf, err := imageToBytes(img, contentType)
	if err != nil {
		return "", err
	}

	return OcrTextFromBytes(re, buf)
}

// 从 []image.Image 获取 []string
func OcrTextFromImages(imgs []image.Image, contentType string, re request.OcrBase) ([]string, error) {
	var results []string

	global.BANK_LOGGER.Debug("start ocr  image")
	for _, img := range imgs {
		buf, err := imageToBytes(img, contentType)

		if err != nil {
			return nil, err
		}
		global.BANK_LOGGER.Debug(" ocr one image")
		result, err := OcrTextFromBytes(re, buf)
		if err != nil {
			return nil, err
		}
		global.BANK_LOGGER.Debug(" ocr one image result=" + result)

		results = append(results, result)
	}

	return results, nil
}

func GetTesseractInfo() (*response.InfoResponse, error) {
	langs, err := gosseract.GetAvailableLanguages()
	if err != nil {
		return nil, err
	}

	client := gosseract.NewClient()
	defer client.Close()

	info := response.TesseractInfo{
		Version:   client.Version(),
		Languages: langs,
	}

	return &response.InfoResponse{
		Tesseract: info,
	}, nil
}

func OcrTextFromBytes(req request.OcrBase, bytes []byte) (string, error) {
	client := gosseract.NewClient()
	defer client.Close()

	var err error
	if err = client.SetImageFromBytes(bytes); err != nil {
		return "", err
	}

	client.Languages = global.BANK_CONFIG.Ocr.DefaultLangs

	if req.Languages != "" {
		client.Languages = strings.Split(req.Languages, ",")
	}

	if req.Whitelist != "" {
		if err := client.SetWhitelist(req.Whitelist); err != nil {
			return "", err
		}
	}

	var text string
	if req.HOCRMode {
		text, err = client.HOCRText()
	} else {
		text, err = client.Text()
	}
	if err != nil {
		return "", err
	}

	if req.TrimLineFeed {
		text = strings.Replace(text, "\n", "", -1)
	}

	return text, nil
}

func imageToBytes(img image.Image, contentType string) ([]byte, error) {
	var err error
	buf := bytes.NewBuffer(nil)

	switch contentType {
	case "image/png":
		err = png.Encode(buf, img)
	case "image/jpeg":
		err = jpeg.Encode(buf, img, nil)
	case "image/gif":
		err = gif.Encode(buf, img, nil)
	case "image/tiff":
		err = tiff.Encode(buf, img, nil)
	default:
		err = errors.New("cannot support current file type")
	}

	return buf.Bytes(), err
}

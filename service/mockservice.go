/**
 * @Author: lzw5399
 * @Date: 2020/9/30 15:12
 * @Desc:
 */
package service

import (
	"bank-ocr/model/request"
	"bank-ocr/model/response"
	"bytes"
	"errors"
	"github.com/disintegration/imaging"
	"image"
	"io"
	"mime/multipart"
)

func GetTextFromImageV2(img image.Image, contentType string, re request.FileFormRequest) (string, error) {
	return "123", nil
}

func GetHOCRTextFromImageV2(img image.Image, contentType string, re request.FileFormRequest) (string, error) {
	return "456", nil
}

func GetTesseractInfoV2() (*response.InfoResponse, error) {
	info := response.TesseractInfo{
		Version:   "mock",
		Languages: nil,
	}

	return &response.InfoResponse{
		Tesseract: info,
	}, nil
}

func OcrTextFromImagesV2(imgs []image.Image, contentType string, re request.FileFormRequest) ([]string, error) {
	return nil, errors.New("mock")
}

func GrayImageV2(f multipart.File, re request.FileFormRequest) (image.Image, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	return imaging.Decode(reader)
}

/**
 * @Author: lzw5399
 * @Date: 2020/10/2 15:12
 * @Desc:
 */
package service

import (
	"bank-ocr/model/request"
	"bank-ocr/model/response"
	"errors"
	"image"
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
/**
 * @Author: lzw5399
 * @Date: 2020/9/30 23:45
 * @Desc: image processing service
 */
package service

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"bank-ocr/global"
	req "bank-ocr/model/request"

	"github.com/disintegration/imaging"
	uuid "github.com/satori/go.uuid"
)

const (
	JPEG = iota
	PNG
	GIF
)

// 通过像素点裁剪
func CropImageByPixelPosition() {
	img, err := imaging.Open("images/xuezhixia.png")
	if err != nil {
		global.BANK_LOGGER.Fatalf("failed to open image: %v\n", err)
	}

	// 根据两个坐标点裁剪
	img = imaging.Crop(img, image.Rect(0, 0, 200, 200))

	// 灰度化
	img = imaging.Grayscale(img)

	id := uuid.NewV4().String()
	fmt.Printf("current file name is %s \n", id)

	err = imaging.Save(img, fmt.Sprintf("images/%s.png", id))
	if err != nil {
		global.BANK_LOGGER.Fatal(err)
	}
}

func CopyToTempFile(upload multipart.File) (*os.File, error) {
	// Create temp physical file
	tempFile, err := ioutil.TempFile("", "ocrserver"+"-")
	if err != nil {
		return nil, err
	}

	// Make uploaded physical
	_, err = io.Copy(tempFile, upload)

	return tempFile, err
}

// 图片灰度化
func GrayscaleImage() {

}

var supportImgType = [4]string{
	"image/png",
	"image/jpeg",
	"image/gif",
	"image/tiff",
}

func EnsureFileType(f multipart.File) (bool, string, error) {
	buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	if _, err := f.Read(buff); err != nil {
		return false, "", err
	}

	contentType := http.DetectContentType(buff)

	// 把偏移量移回0
	f.Seek(0, 0)

	for _, v := range supportImgType {
		if contentType == v {
			return true, contentType, nil
		}
	}

	return false, contentType, nil
}

// 像素点切割和灰度化, 返回image切片
func CropAndGrayImage(f multipart.File, re req.FileWithPixelPointRequest) ([]image.Image, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	imgs := make([]image.Image, len(re.MatrixPixels))

	for i, v := range re.MatrixPixels {
		var tempImg image.Image = imaging.Crop(img, image.Rect(v.PointA.X, v.PointA.Y, v.PointB.X, v.PointB.Y))
		tempImg = imaging.Grayscale(tempImg)
		imaging.Save(tempImg, "testdata/666.png")
		imgs[i] = tempImg
	}

	return imgs, nil
}

// 图像灰度化
func GrayImage(f multipart.File, re req.FileFormRequest) (image.Image, error) {
	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, f); err != nil {
		return nil, err
	}

	reader := bytes.NewReader(buf.Bytes())

	img, err := imaging.Decode(reader)
	if err != nil {
		return nil, err
	}

	img = imaging.Grayscale(img)

	return img, nil
}

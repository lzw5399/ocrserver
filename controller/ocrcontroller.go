/**
 * @Author: lzw5399
 * @Date: 2020/9/30 23:24
 * @Desc: ocr related functionality
 */
package controller

import (
	"encoding/json"
	"net/http"

	. "bank-ocr/global"
	"bank-ocr/global/response"
	"bank-ocr/model/request"
	"bank-ocr/service"

	"github.com/gin-gonic/gin"
)

const (
	INVALID_IMG_TYPE_MSG = "invalid file or unsupported file type. Only support .jpg .jpeg .png .gif .tiff, please double check!"
)

// @Tags ocr
// @Summary OCR识别上传的整张图片
// @Accept  x-www-form-urlencoded
// @Produce json
// @param file formData file true "图片文件"
// @param languages formData string false "可选项: 指定要识别的语言种类，如eng(英文) chi_sim(简体中文)，可以用逗号隔开指定多个, 不指定默认是eng"
// @param whitelist formData string false "可选项: 为空检测全部字符。如果填写，仅会检测白名单之内的字符"
// @param hocrMode formData bool false "可选项: 是否开始HOCR，一般默认为false"
// @Success 200 {object} response.HttpResponse
// @Router /api/ocr/file [post]
func ScanFile(c *gin.Context) {
	var r request.FileFormRequest
	if err := c.ShouldBind(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	upload, err := r.File.Open()
	if err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, INVALID_IMG_TYPE_MSG)
		return
	}
	upload.Close()

	// 确保file类型是支持的image类型
	valid, contentType, err := service.EnsureFileType(upload)
	if err != nil || !valid {
		response.FailWithMsg(c, http.StatusBadRequest, INVALID_IMG_TYPE_MSG)
		return
	}

	// 针对像素坐标点进行裁剪并灰度化
	img, err := service.GrayImageV2(upload, r)
	if err != nil {
		BANK_LOGGER.Error(err)
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	// 根据format类型返回ocr最终的值
	text, err := service.GetTextFromImage(img, contentType, r)

	if err != nil {
		BANK_LOGGER.Error(err)
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	response.OkWithData(c, text)
}

// @Tags ocr
// @Summary OCR识别像素点裁剪之后的图片文字
// @Accept  x-www-form-urlencoded
// @Produce json
// @param file formData file true "图片文件"
// @param matrixPixels formData string true "裁剪像素点。必须是下面格式的合法json字符串： [{ pointA: {x: 127, y: 249}, pointB: {x: 983, y: 309}}]"
// @param languages formData string false "可选项: 指定要识别的语言种类，如eng(英文) chi_sim(简体中文)，可以用逗号隔开指定多个, 不指定默认是eng"
// @param whitelist formData string false "可选项: 为空检测全部字符。如果填写，仅会检测白名单之内的字符"
// @param hocrMode formData bool false "可选项: 是否开始HOCR，一般默认为false"
// @Success 200 {object} response.HttpResponse
// @Router /api/ocr/scan-crop-file [post]
func ScanCropFile(c *gin.Context) {
	var r request.FileWithPixelPointRequest
	if err := c.ShouldBind(&r); err != nil {
		response.Failed(c, http.StatusBadRequest)
		return
	}

	// 绑定像素点 (gin的bind不能绑定formdata的对象数组)
	b := c.PostFormArray("matrixPixels")
	if len(b) == 0 {
		response.FailWithMsg(c, http.StatusBadRequest, "matrixPixels invalid.")
		return
	}
	var matrixPixels []request.MatrixPixel
	err := json.Unmarshal([]byte(b[0]), &matrixPixels)
	if err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, "matrixPixels invalid.")
		return
	}
	r.MatrixPixels = matrixPixels

	// 获取file
	upload, err := r.File.Open()
	if err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, INVALID_IMG_TYPE_MSG)
		return
	}
	defer upload.Close()

	// 确保file类型是支持的image类型
	valid, contentType, err := service.EnsureFileType(upload)
	if err != nil || !valid {
		response.FailWithMsg(c, http.StatusBadRequest, INVALID_IMG_TYPE_MSG)
		return
	}

	// 针对像素坐标点进行裁剪并灰度化
	imgs, err := service.CropAndGrayImage(upload, r)
	if err != nil {
		BANK_LOGGER.Error(err)
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	// 裁剪之后的图片进行ocr识别
	texts, err := service.OcrTextFromImages(imgs, contentType, r.ToFileFormRequest())
	if err != nil {
		BANK_LOGGER.Error(err)
		response.Failed(c, http.StatusInternalServerError)
		return
	}

	response.OkWithData(c, texts)
}

func Base64(c *gin.Context) {
	//var body = new(struct {
	//	Base64    string `json:"base64"`
	//	Trim      string `json:"trim"`
	//	Languages string `json:"languages"`
	//	Whitelist string `json:"whitelist"`
	//})
	//
	//err := json.NewDecoder(c.Request.Body).Decode(body)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, err)
	//	return
	//}
	//
	//tempFile, err := ioutil.TempFile("", "ocrserver"+"-")
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, err)
	//	return
	//}
	//defer func() {
	//	tempFile.Close()
	//	os.Remove(tempFile.Name())
	//}()
	//
	//if len(body.Base64) == 0 {
	//	c.JSON(http.StatusBadRequest, fmt.Errorf("base64 string required"))
	//	return
	//}
	//body.Base64 = regexp.MustCompile("data:image\\/png;base64,").ReplaceAllString(body.Base64, "")
	//b, err := base64.StdEncoding.DecodeString(body.Base64)
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, err)
	//	return
	//}
	//tempFile.Write(b)
	//
	//client := gosseract.NewClient()
	//defer client.Close()
	//
	//client.Languages = []string{"eng"}
	//if body.Languages != "" {
	//	client.Languages = strings.Split(body.Languages, ",")
	//}
	//client.SetImage(tempFile.Name())
	//if body.Whitelist != "" {
	//	client.SetWhitelist(body.Whitelist)
	//}
	//
	//text, err := client.Text()
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, err)
	//	return
	//}
	//
	//c.JSON(http.StatusOK, gin.H{
	//	"result":  strings.Trim(text, body.Trim),
	//	"version": version,
	//})
}

/**
 * @Author: lzw5399
 * @Date: 2020/10/2 15:18
 * @Desc:
 */
package response

type InfoResponse struct {
	Tesseract TesseractInfo `json:"tesseract"`
}

type TesseractInfo struct {
	Version   string   `json:"version"`
	Languages []string `json:"languages"`
}

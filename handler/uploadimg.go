package handler

import (
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/disposition"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// 获取token 已废弃
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error"`
}

/////////////////////////////////////////////////////////////////////////路由2/////////////////////////////////////////////

// 路由2的处理器: 上传图片并发送给AI服务器
func UploadHandler(c *gin.Context) {
	// 标记2_1: 获取上传的文件
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		// 标记2_1_1: 文件解析失败
		respondWithJSON(c, http.StatusBadRequest, "文件解析失败", nil)
		return
	}
	defer file.Close()

	// 标记2_2: 构造保存文件路径
	ext := filepath.Ext(header.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(disposition.UploadDir, filename)

	// 标记2_3: 保存文件到本地
	savedFile, err := os.Create(filePath)
	if err != nil {
		// 标记2_3_1: 文件保存失败
		respondWithJSON(c, http.StatusInternalServerError, "文件保存失败", nil)
		return
	}
	defer savedFile.Close()
	if _, err := io.Copy(savedFile, file); err != nil {
		// 标记2_3_2: 文件复制失败
		respondWithJSON(c, http.StatusInternalServerError, "文件保存失败", nil)
		return
	}

	// 标记2_4: 发送图片给AI服务器
	aiResult, err := sendImageToAIServer(filePath)
	if err != nil {
		// 标记2_4_1: AI处理失败
		respondWithJSON(c, http.StatusInternalServerError, "AI处理失败", nil)
		return
	}

	// 标记2_5: 保存AI处理结果
	if err := saveAIResult(aiResult, disposition.AiResultsDir, filename+"_results.json"); err != nil {
		// 标记2_5_1: 保存失败
		respondWithJSON(c, http.StatusInternalServerError, "保存AI结果失败", nil)
		return
	}

	// 标记2_6: 返回图片的URL
	imageURL := fmt.Sprintf("%s/images/%s", disposition.ServerHost, filename)
	respondWithJSON(c, http.StatusOK, "文件上传成功", map[string]string{"image_url": imageURL})
}

/////////////////////////////////////////////////////辅助函数//////////////////////////////////////////////////////////////////

// 辅助函数: 发送图片到AI服务器并返回结果
func sendImageToAIServer(imagePath string) (string, error) {
	// 标记2_4_2_1: 读取图片文件
	imgData, err := ioutil.ReadFile(imagePath)
	if err != nil {
		return "", err
	}

	// 标记2_4_2_2: 编码图片为Base64
	imgStr := base64.StdEncoding.EncodeToString(imgData)
	imgParam := url.QueryEscape(imgStr)

	// 标记2_4_2_3: 获取access_token
	accessToken := getAccessToken()
	if accessToken == "" {
		return "", fmt.Errorf("获取access_token失败")
	}

	// 标记2_4_2_4: 构造请求并发送
	requesturl := fmt.Sprintf("%s?access_token=%s", disposition.OcrURL, accessToken)
	payload := strings.NewReader("image=" + imgParam + "&probability=false&poly_location=false")
	client := &http.Client{}
	req, err := http.NewRequest("POST", requesturl, payload)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// 标记2_4_2_5: 读取响应体
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// 辅助函数: 保存AI处理结果到指定文件
func saveAIResult(aiResult, dir, filename string) error {
	// 标记2_5_2: 创建结果保存路径
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 标记2_5_3: 写入AI结果
	_, err = file.WriteString(aiResult)
	return err
}

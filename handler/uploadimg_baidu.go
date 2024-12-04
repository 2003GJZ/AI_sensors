package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/disposition"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// 获取token 已废弃
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Error       string `json:"error"`
}

/////////////////////////////////////////////////////////////////////////路由2/////////////////////////////////////////////

// 路由2的处理器: 上传图片并发送给AI服务器
//func UploadHandler(c *gin.Context) {
//	// 标记2_1: 获取上传的文件
//	id := c.PostForm("id")
//	file, header, err := c.Request.FormFile("image")
//	if err != nil {
//		// 标记2_1_1: 文件解析失败
//		respondWithJSON(c, http.StatusBadRequest, "文件解析失败", nil)
//		return
//	}
//	defer file.Close()
//
//	// 标记2_2: 构造保存文件路径
//	ext := filepath.Ext(header.Filename)
//	filename := fmt.Sprintf("%d%s%s", time.Now().UnixNano(), ext)
//	filePath := filepath.Join(disposition.UploadDir, filename)
//
//	// 标记2_3: 保存文件到本地
//	savedFile, err := os.Create(filePath)
//	if err != nil {
//		// 标记2_3_1: 文件保存失败
//		respondWithJSON(c, http.StatusInternalServerError, "文件保存失败", nil)
//		return
//	}
//	defer savedFile.Close()
//	if _, err := io.Copy(savedFile, file); err != nil {
//		// 标记2_3_2: 文件复制失败
//		respondWithJSON(c, http.StatusInternalServerError, "文件保存失败", nil)
//		return
//	}
//
//	// 标记2_4: 发送图片给AI服务器
//	aiResult, err := sendImageToAIServer(filePath)
//	if err != nil {
//		// 标记2_4_1: AI处理失败
//		respondWithJSON(c, http.StatusInternalServerError, "AI处理失败", nil)
//		return
//	}
//
//	// 标记2_5: 保存AI处理结果
//	if err := saveAIResult(aiResult, disposition.AiResultsDir, filename+"_results.json"); err != nil {
//		// 标记2_5_1: 保存失败
//		respondWithJSON(c, http.StatusInternalServerError, "保存AI结果失败", nil)
//		return
//	}
//
//	// 标记2_6: 返回图片的URL
//	imageURL := fmt.Sprintf("%s/images/%s", disposition.ServerHost, filename)
//	respondWithJSON(c, http.StatusOK, "文件上传成功", map[string]string{"image_url": imageURL})
//}

func UploadHandler(c *gin.Context) {
	// 标记2_1: 获取上传的文件
	id := c.PostForm("id")
	file, header, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	// 创建保存文件的目录
	//savePath := "/var/ftp/ftpuser/" + id
	savePath := "D:\\var\\ftp\\ftpuser\\" + id
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		if err := os.MkdirAll(savePath, 0755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create directory: " + err.Error()})
			return
		}
	}

	// 构建完整的文件路径
	outPath := filepath.Join(savePath, header.Filename)

	// 将文件保存到指定路径
	if err := c.SaveUploadedFile(header, outPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file: " + err.Error()})
	}
	//通知图片处理
	go func() {
		irl := "http://127.0.0.1:4399/upload_success"
		// 创建一个buffer用于存储表单数据
		var buffer bytes.Buffer

		// 创建一个multipart writer
		writer := multipart.NewWriter(&buffer)

		// 添加字段
		err := writer.WriteField("id", id)
		if err != nil {
			panic(err)
		}
		//将文件名传过去
		err = writer.WriteField("imgname", header.Filename)
		if err != nil {
			panic(err)
		}

		// 关闭writer以确保所有数据都被写入
		err = writer.Close()
		if err != nil {
			panic(err)
		}

		// 创建请求
		http.Post(irl, writer.FormDataContentType(), &buffer)
		//fmt.Println(&buffer)

	}()
	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "id": id, "filename": header.Filename})
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

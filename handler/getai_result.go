package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/disposition"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type MeterResponse struct {
	LogID     interface{} `json:"log_id"`
	WordsNum  uint32      `json:"words_result_num"`
	WordsList []struct {
		Words string `json:"words"`
	} `json:"words_result"`
}

// 路由3的处理器: 获取AI结果
func AiResultHandler(c *gin.Context) {
	// 标记3_1: 获取文件名
	filename := c.Param("filename")
	if filename == "" {
		// 标记3_1_1: 未传递文件名
		respondWithJSON(c, http.StatusBadRequest, "缺少文件名", nil)
		return
	}

	// 标记3_2: 构造结果文件路径
	filePath := filepath.Join(disposition.AiResultsDir, filename+"_results.json")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// 标记3_2_1: 文件不存在
		respondWithJSON(c, http.StatusNotFound, "结果文件不存在", nil)
		return
	}

	// 标记3_3: 读取文件内容
	body, err := ioutil.ReadFile(filePath)
	if err != nil {
		// 标记3_3_1: 文件读取失败
		respondWithJSON(c, http.StatusInternalServerError, "读取结果文件失败", nil)
		return
	}
	//3_4 解析JSON
	var meterResp MeterResponse
	err = json.Unmarshal(body, &meterResp)
	if err != nil {
		respondWithJSON(c, 201, "解析JSON失败", nil)
		return
	}

	// 3_5检查是否有"words"字段
	if meterResp.WordsNum == 0 {
		respondWithJSON(c, 201, "AI处理失败", nil)
		return
	}

	// 3_6提取"words"字段并用逗号分隔
	var words []string
	for _, item := range meterResp.WordsList {
		words = append(words, item.Words)
	}
	resultText := strings.Join(words, ",")

	// 标记3_4: 返回结果
	respondWithJSON(c, 200, "成功", map[string]string{"text": resultText})
}

// 辅助函数: 统一响应封装
func respondWithJSON(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, gin.H{
		"code":    statusCode,
		"message": message,
		"data":    data,
	})
}

// 辅助函数: 获取Access Token
func getAccessToken() string {
	// 标记: 构造URL和参数
	values := url.Values{}
	values.Set("grant_type", "client_credentials")
	values.Set("client_id", disposition.API_KEY)
	values.Set("client_secret", disposition.SECRET_KEY)

	// 标记: 发送请求
	resp, err := http.PostForm(disposition.TokenURL, values)
	if err != nil {
		log.Printf("获取Access Token失败: %v", err)
		return ""
	}
	defer resp.Body.Close()

	// 标记: 解析响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("读取Access Token响应体失败: %v", err)
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		log.Printf("解析Access Token响应失败: %v", err)
		return ""
	}

	// 标记: 提取Access Token
	if token, ok := result["access_token"].(string); ok {
		return token
	}

	log.Printf("Access Token不存在于响应中")
	return ""
}

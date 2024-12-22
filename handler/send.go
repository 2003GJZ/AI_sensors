package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
)

// 发送报文到网关
func SendReportHandler(c *gin.Context) {
	var requestBody struct {
		MeterID string `json:"meter_id"`
		Data    struct {
			ACurrent string `json:"ACurrent"`
			BCurrent string `json:"BCurrent"`
			CCurrent string `json:"CCurrent"`
			AVoltage string `json:"AVoltage"`
			BVoltage string `json:"BVoltage"`
			CVoltage string `json:"CVoltage"`
			Power    string `json:"Power"`
		} `json:"data"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("请求参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数绑定失败",
		})
		return
	}

	meterID := requestBody.MeterID
	if meterID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "缺少 meter_id 参数",
		})
		return
	}

	//// 直接使用请求中的数据
	//attributes := map[string]string{
	//	"ACurrent": requestBody.Data.ACurrent,
	//	"BCurrent": requestBody.Data.BCurrent,
	//	"CCurrent": requestBody.Data.CCurrent,
	//	"AVoltage": requestBody.Data.AVoltage,
	//	"BVoltage": requestBody.Data.BVoltage,
	//	"CVoltage": requestBody.Data.CVoltage,
	//	"Power":    requestBody.Data.Power,
	//}
	//
	//// 将 attributes 转换为 JSON 字符串
	//attributesJson, err := json.Marshal(attributes)
	//if err != nil {
	//	log.Printf("构建 JSON 属性失败: %v", err)
	//	c.JSON(http.StatusInternalServerError, gin.H{
	//		"status":  "error",
	//		"message": "构建 JSON 属性失败",
	//	})
	//	return
	//}

	attributesJson, err := json.Marshal(requestBody.Data)
	fmt.Println(attributesJson)

	// 对 JSON 字符串进行 Base64 编码
	base64Message := base64.StdEncoding.EncodeToString(attributesJson)

	// 发送报文到网关
	err = sendToGateway(meterID, base64Message)
	if err != nil {
		log.Printf("发送报文失败: %v", err)
		fmt.Printf("发送报文失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "发送报文失败",
		})
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "报文已成功发送",
	})
}

func sendToGateway(topic string, message string) error {
	payload := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("构建 JSON Payload 失败: %v", err)
		fmt.Printf("构建 JSON Payload 失败: %v\n", err)
		return err
	}

	resp, err := http.Post("http://localhost:4366/publish", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("发送 HTTP POST 请求失败: %v", err)
		fmt.Printf("发送 HTTP POST 请求失败: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("网关返回错误状态码: %d, 响应体: %s", resp.StatusCode, string(body))
		fmt.Printf("网关返回错误状态码: %d, 响应体: %s\n", resp.StatusCode, string(body))
		return fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	log.Printf("成功发送报文到网关")
	return nil
}

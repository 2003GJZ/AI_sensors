package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/protocol_stack"
	_ "imgginaimqtt/protocol_stack"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
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
	err1 := subscribeToGateway(meterID)
	if err1 != nil {
		log.Fatalf("订阅主题失败: %v", err1)
	}
	//---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------
	// 定义字段与命令码的映射
	fieldCommandMap := map[string]struct {
		Command byte
		Data    []byte
	}{
		"ACurrent": {0x11, []byte{0x00, 0x00, 0x01, 0x02}},
		"BCurrent": {0x12, []byte{0x00, 0x00, 0x01, 0x03}},
		"CCurrent": {0x13, []byte{0x00, 0x00, 0x01, 0x04}},
		"AVoltage": {0x14, []byte{0x00, 0x00, 0x01, 0x05}},
		"BVoltage": {0x15, []byte{0x00, 0x00, 0x01, 0x06}},
		"CVoltage": {0x16, []byte{0x00, 0x00, 0x01, 0x07}},
		"Power":    {0x17, []byte{0x00, 0x00, 0x01, 0x08}},
	}

	var newFrame []byte
	var err error

	// 循环遍历 Data 结构体中的字段
	for fieldName, fieldData := range fieldCommandMap {
		fieldValue := reflect.ValueOf(requestBody.Data).FieldByName(fieldName).String()
		if fieldValue != "" {
			newFrame, err = protocol_stack.BuildDLT645Frame(meterID, fieldData.Command, fieldData.Data)
			if err != nil {
				log.Printf("生成失败: %v\n", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "生成帧失败",
				})
				return
			}
			// 对 JSON 字符串进行 Base64 编码
			base64Message := base64.StdEncoding.EncodeToString(newFrame)

			// 发送报文到网关
			err = sendToGateway(meterID, base64Message)
			if err != nil {
				log.Printf("发送报文失败: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"status":  "error",
					"message": "发送报文失败",
				})
				return
			}
		}
	}
	//---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------

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
		return err
	}

	resp, err := http.Post("http://localhost:4366/publish", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("发送 HTTP POST 请求失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("网关返回错误状态码: %d, 响应体: %s", resp.StatusCode, string(body))
		return fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	log.Printf("成功发送报文到网关")
	return nil
}

// 订阅网关
func subscribeToGateway(topic string) error {
	payload := map[string]interface{}{
		"topic": topic,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("构建 JSON Payload 失败: %v", err)
		return err
	}

	resp, err := http.Post("http://localhost:4366/subscribe", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Printf("发送 HTTP POST 请求失败: %v", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		log.Printf("网关返回错误状态码: %d, 响应体: %s", resp.StatusCode, string(body))
		return fmt.Errorf("网关返回错误状态码: %d", resp.StatusCode)
	}

	log.Printf("成功订阅主题: %s", topic)
	return nil
}

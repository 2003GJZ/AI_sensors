package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/protocol_stack"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// 发送报文到网关
func SendReportHandler(c *gin.Context) {

	var requestBody dao.Message
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("请求参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "请求参数绑定失败",
		})
		return
	}

	meterID := requestBody.DeviceID
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

	hexKey, exists := protocol_stack.GetKeyByDescription(requestBody.Data)
	if exists != nil {
		log.Printf("错误: %v\n", exists)
	} else {
		log.Printf("找到的键: %s\n", hexKey)
	}

	byteArray, err := HexKeyToByteArray(hexKey)
	if err != nil {
		log.Printf("转换失败: %v\n", err)
	} else {
		log.Printf("转换成功: %v\n", byteArray)
	}

	// 地址
	address := meterID
	// 控制码
	control := byte(0x11)
	// 数据域
	data := byteArray

	var newFrame []byte

	newFrame, err = protocol_stack.BuildDLT645Frame(address, control, data)
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

// HexKeyToByteArray 将字符串键转换为字节数组
func HexKeyToByteArray(hexKey string) ([]byte, error) {
	// 分割字符串
	parts := strings.Split(hexKey, "-")
	if len(parts) != 4 {
		return nil, fmt.Errorf("无效的 hexKey 格式: %s", hexKey)
	}

	// 初始化字节数组
	byteArray := make([]byte, 4)

	// 转换每个部分为字节并加上偏移量 0x33
	for i, part := range parts {
		value, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析部分 %s 失败: %v", part, err)
		}
		byteArray[i] = byte(value) + 0x33
	}

	return byteArray, nil
}

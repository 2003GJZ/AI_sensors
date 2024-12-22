package handler

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"imgginaimqtt/protocol_stack"
	"io/ioutil"
	"log"
	"net/http"
)

// MQTT -> HTTP 消息结构
type MQTTMessage struct {
	Topic   string `json:"topic"`
	Message string `json:"message"` // Base64 编码后的消息
}

// 用来保存电表结果
var AmmeterMap = make(map[string]*dao.Ammeter)

// MQTT 处理器
func MqttBaes64Handler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			respond(c, 500, "服务器内部错误", nil)
		}
	}()

	// 读取请求体
	//body, err := ioutil.ReadAll(c.Request.Body)
	//接收请求体,用结构体
	var req MQTTMessage
	if err := c.ShouldBindJSON(&req); err != nil || req.Topic == "" || req.Message == "" {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		return
	}

	defer c.Request.Body.Close()

	//bese64解码
	bytes, err1 := protocol_stack.MyBase64ToBytes(string(req.Message))
	if err1 != nil {
		respond(c, 400, "base64解码失败", nil)
	}
	//dlt协议解析
	err, i := protocol_stack.ElectricityAnswer(bytes)
	if err1 != nil {

		log.Printf("应答解析失败: %v\n", err1)
		respond(c, 400, "DLT解码失败", nil)
	}
	frame, err := protocol_stack.ParseDLT645Frame(i)
	// 数据域去偏移
	decodedData := protocol_stack.OffsetData(frame.Data, false)
	// 调用解析函数
	dataType, value, phase, err := protocol_stack.ParseDataSegment(decodedData)

	if err != nil {
		log.Printf("解析失败: %v\n", err)
		respond(c, 400, "DLT解码失败", nil)
		return
	}

	log.Printf("解析结果: 类型 = %s, 值 = %s, 相位 = %s\n", dataType, value, phase)
	//查询Map
	ammeter, ok := AmmeterMap[frame.Address]
	if !ok {
		// 如果不存在，则创建一个新的电表对象并添加到Map中
		AmmeterMap[frame.Address] = &dao.Ammeter{DeviceID: frame.Address}
		ammeter = AmmeterMap[frame.Address]
	}
	// 更新电表对象
	switch dataType {
	case "I":
		switch phase {
		case "A":
			ammeter.ACurrent = value
		case "B":
			ammeter.BCurrent = value
		case "C":
			ammeter.CCurrent = value
		case "O":
			ammeter.Current = value
		}

	case "V":
		switch phase {
		case "A":
			ammeter.AVoltage = value
		case "B":
			ammeter.BVoltage = value
		case "C":
			ammeter.CVoltage = value
		case "O":
			ammeter.Voltage = value
		}

	case "P":
		ammeter.Power = value
	}

	// 保存到Redis-------->
	// 将结构体转换为JSON
	jsonData, _ := json.Marshal(ammeter)

	// 将JSON数据转换为字符串
	jsonString := string(jsonData)
	//TODO 将字符串保存到Redis
	ammeter.DeviceID = "123456"

	link, _ := mylink.GetredisLink()
	link.Client.HSet(link.Ctx, "ai_value", "123456", jsonString)
	// 将字符串保存到Redis
	respond(c, 200, "数据处理成功并保存到 Redis！", nil)
	//通知前端
	//TODO 通知前端
	dao.NoticeUpdate(req.Topic)
}

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
		fmt.Printf("请求参数绑定失败: %v\n", err)
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

	// 直接使用请求中的数据
	attributes := map[string]string{
		"ACurrent": requestBody.Data.ACurrent,
		"BCurrent": requestBody.Data.BCurrent,
		"CCurrent": requestBody.Data.CCurrent,
		"AVoltage": requestBody.Data.AVoltage,
		"BVoltage": requestBody.Data.BVoltage,
		"CVoltage": requestBody.Data.CVoltage,
		"Power":    requestBody.Data.Power,
	}

	// 将 attributes 转换为 JSON 字符串
	attributesJson, err := json.Marshal(attributes)
	if err != nil {
		log.Printf("构建 JSON 属性失败: %v", err)
		fmt.Printf("构建 JSON 属性失败: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "构建 JSON 属性失败",
		})
		return
	}

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
		"data":    requestBody.Data,
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

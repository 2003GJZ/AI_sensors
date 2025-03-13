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
	"strconv"
	"strings"
	"time"
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
	link, _ := mylink.GetredisLink()
	//地址

	var metermac string
	link.Client.HGet(link.Ctx, "topic", meterID).Scan(&metermac)

	if metermac == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "未找到 meter_id 对应的 mac 地址",
		})
		return
	}

	// 控制码
	control := byte(0x11)
	// 数据域
	data := byteArray

	var newFrame []byte

	newFrame, err = protocol_stack.BuildDLT645Frame(metermac, control, data)
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

	//meterID = meterID

	fmt.Println("发送报文到网关:", meterID, "///", base64Message)

	//id_V	:	报文
	key := meterID + "_" + requestBody.Data
	dao.BaowenMap.Store(key, base64Message)

	// 发送报文到网关
	err = SendToGateway(meterID, base64Message)
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

func SendToGateway(topic string, message string) error {
	topic = "/down/" + topic
	payload := map[string]interface{}{
		"topic":   topic,
		"message": message,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		log.Printf("构建 JSON Payload 失败: %v", err)
		return err
	}
	//fmt.Println("发送报文到网关:", topic, "///", message)
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

// 网关订阅
func subscribeToGateway(topic string) error {
	topic = "/up/" + topic
	// 定义要发送的订阅请求
	subscribeRequest := map[string]string{
		"topic": topic,
	}

	// 将请求体编码为 JSON
	jsonData, err := json.Marshal(subscribeRequest)
	if err != nil {
		fmt.Printf("编码 JSON 失败: %v\n", err)
		return nil
	}

	// 创建 HTTP POST 请求
	resp, err := http.Post("http://localhost:4366/subscribe", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("发送 POST 请求失败: %v\n", err)
		return nil
	}
	defer resp.Body.Close()

	// 打印响应状态码
	fmt.Printf("响应状态码: %d\n", resp.StatusCode)
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
	// 反转字节数组
	for i, j := 0, len(byteArray)-1; i < j; i, j = i+1, j-1 {
		byteArray[i], byteArray[j] = byteArray[j], byteArray[i]
	}

	return byteArray, nil
}

func BaowenMapFor() {
	for {
		dao.BaowenMap.Range(func(key, value interface{}) bool {
			k := key.(string)
			v := value.(string)
			// 使用Split函数按照下划线分割
			parts := strings.Split(k, "_")
			// 检查分割结果
			if len(parts) == 2 {
				topic := parts[0]
				// 发送报文到网关
				err := SendToGateway(topic, v)
				if err != nil {
					log.Println("发送报文失败: %v", err)
					return false
				}
			} else {
				log.Println("字符串格式不正确")
			}
			//停止1秒
			time.Sleep(1 * time.Second)
			return true
		})
		// 停止指定的时间间隔
		// 将 Interval 转换为 int64 并乘以 time.Minute
		sleepDuration := time.Duration(Interval * float64(time.Minute))
		time.Sleep(sleepDuration)
		fmt.Println("停止", Interval, "分钟")
	}
}

package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"io/ioutil"
	"log"
	"net/http"
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

	// 发送报文到网关
	err := SendToGateway(meterID, requestBody.Data)
	if err != nil {
		log.Printf("发送报文失败: %v", err)
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

func SendToGateway(topic string, message string) error {
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

//func BaowenMapFor() {
//	for {
//		dao.BaowenMap.Range(func(key, value interface{}) bool {
//			k := key.(string)
//			v := value.(string)
//			parts := strings.Split(k, "_")
//			if len(parts) == 2 {
//				topic := parts[0]
//				err := SendToGateway(topic, v)
//				if err != nil {
//					log.Println("发送报文失败: %v", err)
//					return false
//				}
//			} else {
//				log.Println("字符串格式不正确")
//			}
//			time.Sleep(1 * time.Second)
//			return true
//		})
//		sleepDuration := time.Duration(Interval * float64(time.Minute))
//		time.Sleep(sleepDuration)
//		fmt.Println("停止", Interval, "分钟")
//	}
//}

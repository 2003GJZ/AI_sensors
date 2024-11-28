package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"imgginaimqtt/protocol_stack"
	"io/ioutil"
	"log"
)

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
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, dao.ResponseEER_400("err"))
		return
	}
	defer c.Request.Body.Close()

	//bese64解码
	bytes, err1 := protocol_stack.MyBase64ToBytes(string(body))
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

	link, _ := mylink.GetredisLink()
	link.Client.HSet(link.Ctx, "ai_value", ammeter.DeviceID, jsonString)
	// 将字符串保存到Redis
	respond(c, 200, "数据处理成功并保存到 Redis！", nil)
}

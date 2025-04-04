package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	_ "imgginaimqtt/protocol_stack"
	"log"
	"reflect"
)

// 统一响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// 统一响应方法
func respond(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

// 根据 key 从 hset 中获取结构体类型
func getStructTypeFromRedis(hsetKey, key string) (string, error) {
	link, _ := mylink.GetredisLink()
	var structType string
	// 从 Redis 的 hset 中查询结构体类型
	err := link.Client.HGet(link.Ctx, hsetKey, key).Scan(&structType)
	if err != nil || structType == "" {
		return "", fmt.Errorf("Redis HGet 返回空值或错误: %v", err)
	}
	return structType, nil
}

// 动态初始化结构体的新实例
func initStruct(structType string) (interface{}, int) {
	if structTemplateValue, exists := dao.StructRegistry.Load(structType); exists {
		structTemplate, ok := structTemplateValue.(interface{})
		if !ok {
			log.Printf("结构体模板 %s 类型断言失败", structType)
			return make(map[string]interface{}), 0
		}
		// 使用反射创建新实例
		newInstance := reflect.New(reflect.TypeOf(structTemplate)).Interface()
		return newInstance, 1
	}
	// 如果不存在，返回 map[string]interface{}
	return make(map[string]interface{}), 0
}

// 将数据序列化为 JSON 并存入 Redis
func serializeAndSaveToRedis(hsetKey, key string, data interface{}) error {
	link, _ := mylink.GetredisLink()
	// 将数据序列化为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("序列化失败: %v", err)
	}
	// 将 JSON 数据存入 Redis 的 hset
	err = link.Client.HSet(link.Ctx, hsetKey, key, jsonData).Err()
	if err != nil {
		return fmt.Errorf("保存到 Redis 失败: %v", err)
	}
	return nil
}

// 主逻辑方法
func processAndSaveData(structHSetKey, dataHSetKey, key, inputData string) error {
	// 第一步：获取结构体类型
	structType, err := getStructTypeFromRedis(structHSetKey, key)
	if err != nil {
		return fmt.Errorf("获取结构体类型失败: %v", err)
	}

	// 第二步：初始化对应结构体
	dataStruct, typen := initStruct(structType)

	// 第三步：填充数据
	switch typen {
	case 1: // 使用结构体过滤属性
		err := json.Unmarshal([]byte(inputData), &dataStruct)
		if err != nil {
			log.Printf("反序列化失败: %v\t设备: %s\t内容: %s", err, key, inputData)
			return fmt.Errorf("反序列化失败: %v", err)
		}
	case 0: // 如果是动态 map，直接赋值,存储所有属性
		dataStruct = inputData
	default:
		return fmt.Errorf("未知的结构体类型: %v", structType)
	}

	// 第四步：序列化并存储
	err = serializeAndSaveToRedis(dataHSetKey, key, dataStruct)
	if err != nil {
		return fmt.Errorf("保存数据失败: %v", err)
	}

	return nil
}

//// MQTT 处理器
//func MqttHandler(c *gin.Context) {
//	defer func() {
//		if err := recover(); err != nil {
//			respond(c, 500, "服务器内部错误", nil)
//		}
//	}()
//
//	structHSetKey := "type"   // 结构体查询表名
//	dataHSetKey := "ai_value" // 数据存储表名
//
//	// 接收数据
//	message := dao.Message{}
//	if err := c.ShouldBindJSON(&message); err != nil {
//		respond(c, 400, "请求参数错误", nil)
//		return
//	}
//
//	// 执行主逻辑
//	if err := processAndSaveData(structHSetKey, dataHSetKey, message.DeviceID, message.Data); err != nil {
//		respond(c, 500, fmt.Sprintf("数据处理失败: %v", err), nil)
//		return
//	}
//
//	respond(c, 200, "数据处理成功并保存到 Redis！", nil)
//}

// MQTT 处理器
//func MqttHandler(c *gin.Context) {
//	defer func() {
//		if err := recover(); err != nil {
//			respond(c, 500, "服务器内部错误", nil)
//		}
//	}()
//
//	// 读取请求体
//	body, err := ioutil.ReadAll(c.Request.Body)
//	if err != nil {
//		c.JSON(400, dao.ResponseEER_400("err"))
//		return
//	}
//	defer c.Request.Body.Close()
//	//bese64解码
//	body, _ = protocol_stack.MyBase64ToBytes(string(body))
//	respond(c, 200, "数据处理成功并保存到 Redis！", nil)
//}

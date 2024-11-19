package handler

import (
	"encoding/json"
	"testing"

	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"reflect"
)

func TestProcessAndSaveData(t *testing.T) {
	// 初始化 Redis 连接
	link, err := mylink.GetredisLink()
	if err != nil {
		t.Fatalf("Redis 连接失败: %v", err)
	}

	// 注册 Ammeter 结构体
	dao.StructRegistry = map[string]interface{}{
		"Ammeter": dao.Ammeter{},
	}

	// Redis 键
	structHSetKey := "type"
	dataHSetKey := "ai_value"
	key := "device_1"
	inputData := `{"device_id": "123", "current": "10", "voltage": "220", "power": "2000"}`

	// 在 Redis 设置结构体类型
	err = link.Client.HSet(link.Ctx, structHSetKey, key, "Ammeter").Err()
	if err != nil {
		t.Fatalf("设置结构体类型失败: %v", err)
	}

	// 调用主逻辑
	err = processAndSaveData(structHSetKey, dataHSetKey, key, inputData)
	if err != nil {
		t.Errorf("处理数据时出现错误: %v", err)
	}

	// 验证 Redis 数据
	result, err := link.Client.HGet(link.Ctx, dataHSetKey, key).Result()
	if err != nil {
		t.Fatalf("获取 Redis 数据失败: %v", err)
	}

	// 检查存储的 JSON 数据
	expectedData := dao.Ammeter{
		DeviceID: "123",
		Current:  "10",
		Voltage:  "220",
		Power:    "2000",
	}
	expectedJSON, _ := json.Marshal(expectedData)
	if result != string(expectedJSON) {
		t.Errorf("存储数据与预期不符，预期: %s，实际: %s", expectedJSON, result)
	}
}

func TestInitStruct(t *testing.T) {
	// 注册 Ammeter 结构体
	dao.StructRegistry = map[string]interface{}{
		"Ammeter": dao.Ammeter{},
	}

	// 测试已注册结构体
	instance, typ := initStruct("Ammeter")
	if typ != 1 {
		t.Errorf("预期返回类型 1，但收到: %d", typ)
	}
	if _, ok := instance.(*dao.Ammeter); !ok {
		t.Errorf("预期返回 *Ammeter 实例，但收到: %T", instance)
	}

	// 测试未注册结构体
	instance, typ = initStruct("UnknownStruct")
	if typ != 0 {
		t.Errorf("预期返回类型 0，但收到: %d", typ)
	}
	if reflect.TypeOf(instance).Kind() != reflect.Map {
		t.Errorf("预期返回 map 类型，但收到: %T", instance)
	}
}

func TestGetStructTypeFromRedis(t *testing.T) {
	// 初始化 Redis 连接
	link, err := mylink.GetredisLink()
	if err != nil {
		t.Fatalf("Redis 连接失败: %v", err)
	}

	// Redis 键
	hsetKey := "type"
	key := "device_1"
	expectedType := "Ammeter"

	// 设置 Redis 数据
	err = link.Client.HSet(link.Ctx, hsetKey, key, expectedType).Err()
	if err != nil {
		t.Fatalf("设置 Redis 数据失败: %v", err)
	}

	// 调用测试函数
	structType, err := getStructTypeFromRedis(hsetKey, key)
	if err != nil {
		t.Errorf("获取结构体类型时出现错误: %v", err)
	}
	if structType != expectedType {
		t.Errorf("预期返回 %s，但收到: %s", expectedType, structType)
	}
}

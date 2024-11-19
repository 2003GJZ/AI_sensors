package handler

import (
	"encoding/json"
	"fmt"
	"testing"

	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"reflect"
)

func TestProcessAndSaveData2(t *testing.T) {
	// 初始化 Redis 连接
	link, err := mylink.GetredisLink()
	if err != nil {
		t.Fatalf("Redis 连接失败: %v", err)
	}
	fmt.Println("[TestProcessAndSaveData] Redis 连接成功")

	// 注册 Ammeter 结构体
	dao.StructRegistry = map[string]interface{}{
		"Ammeter": dao.Ammeter{},
	}
	fmt.Println("[TestProcessAndSaveData] 已注册结构体: Ammeter")

	// Redis 键和输入数据
	structHSetKey := "type"
	dataHSetKey := "ai_value"
	key := "device_1"
	inputData := `{"device_id": "123", "current": "10", "voltage": "220", "power": "2000"}`

	// 设置 Redis 中的结构体类型
	fmt.Printf("[TestProcessAndSaveData] 向 Redis 设置类型: %s -> %s\n", key, "Ammeter")
	err = link.Client.HSet(link.Ctx, structHSetKey, key, "Ammeter").Err()
	if err != nil {
		t.Fatalf("设置结构体类型失败: %v", err)
	}

	// 调用主逻辑
	fmt.Printf("[TestProcessAndSaveData] 开始处理数据: key=%s, data=%s\n", key, inputData)
	err = processAndSaveData(structHSetKey, dataHSetKey, key, inputData)
	if err != nil {
		t.Errorf("处理数据时出现错误: %v", err)
	} else {
		fmt.Println("[TestProcessAndSaveData] 数据处理成功")
	}

	// 验证 Redis 数据
	fmt.Printf("[TestProcessAndSaveData] 从 Redis 获取数据: key=%s\n", key)
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
	fmt.Printf("[TestProcessAndSaveData] 预期数据: %s\n", expectedJSON)
	fmt.Printf("[TestProcessAndSaveData] 实际数据: %s\n", result)

	if result != string(expectedJSON) {
		t.Errorf("存储数据与预期不符，预期: %s，实际: %s", expectedJSON, result)
	} else {
		fmt.Println("[TestProcessAndSaveData] 数据验证成功")
	}
}

func TestInitStruct2(t *testing.T) {
	// 注册 Ammeter 结构体
	dao.StructRegistry = map[string]interface{}{
		"Ammeter": dao.Ammeter{},
	}
	fmt.Println("[TestInitStruct] 已注册结构体: Ammeter")

	// 测试已注册结构体
	fmt.Println("[TestInitStruct] 测试已注册结构体: Ammeter")
	instance, typ := initStruct("Ammeter")
	if typ != 1 {
		t.Errorf("预期返回类型 1，但收到: %d", typ)
	} else {
		fmt.Println("[TestInitStruct] 类型验证成功: Ammeter")
	}
	if _, ok := instance.(*dao.Ammeter); !ok {
		t.Errorf("预期返回 *Ammeter 实例，但收到: %T", instance)
	} else {
		fmt.Println("[TestInitStruct] 实例验证成功: Ammeter")
	}

	// 测试未注册结构体
	fmt.Println("[TestInitStruct] 测试未注册结构体: UnknownStruct")
	instance, typ = initStruct("UnknownStruct")
	if typ != 0 {
		t.Errorf("预期返回类型 0，但收到: %d", typ)
	} else {
		fmt.Println("[TestInitStruct] 类型验证成功: UnknownStruct")
	}
	if reflect.TypeOf(instance).Kind() != reflect.Map {
		t.Errorf("预期返回 map 类型，但收到: %T", instance)
	} else {
		fmt.Println("[TestInitStruct] 实例验证成功: UnknownStruct")
	}
}

func TestGetStructTypeFromRedis2(t *testing.T) {
	// 初始化 Redis 连接
	link, err := mylink.GetredisLink()
	if err != nil {
		t.Fatalf("Redis 连接失败: %v", err)
	}
	fmt.Println("[TestGetStructTypeFromRedis] Redis 连接成功")

	// Redis 键
	hsetKey := "type"
	key := "device_1"
	expectedType := "Ammeter"

	// 设置 Redis 数据
	fmt.Printf("[TestGetStructTypeFromRedis] 设置 Redis 数据: %s -> %s\n", key, expectedType)
	err = link.Client.HSet(link.Ctx, hsetKey, key, expectedType).Err()
	if err != nil {
		t.Fatalf("设置 Redis 数据失败: %v", err)
	}

	// 调用测试函数
	fmt.Printf("[TestGetStructTypeFromRedis] 获取 Redis 数据: key=%s\n", key)
	structType, err := getStructTypeFromRedis(hsetKey, key)
	if err != nil {
		t.Errorf("获取结构体类型时出现错误: %v", err)
	}
	fmt.Printf("[TestGetStructTypeFromRedis] 获取到的结构体类型: %s\n", structType)

	if structType != expectedType {
		t.Errorf("预期返回 %s，但收到: %s", expectedType, structType)
	} else {
		fmt.Println("[TestGetStructTypeFromRedis] 类型验证成功")
	}
}

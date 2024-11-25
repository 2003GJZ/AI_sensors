package dao

// 定义结构体注册表
var StructRegistry = map[string]interface{}{}

type UpdataMac struct {
	NeedsImage bool
	NewMAC     string
}

// 模拟MAC地址和状态的存储
var MacAddressStatus = map[string]UpdataMac{}

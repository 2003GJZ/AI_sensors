package dao

// 定义结构体注册表
var StructRegistry = map[string]interface{}{}

// 模拟MAC地址和状态的存储
var MacAddressStatus = map[string]*UpdataMacImg{}

// aimodel表
var AimodelTable = map[string]Aimodel{}

var BaowenMap = map[string]string{}

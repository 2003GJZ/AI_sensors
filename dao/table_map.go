package dao

import (
	"sync"
)

// 定义结构体注册表
var StructRegistry sync.Map

// 模拟MAC地址和状态的存储
var MacAddressStatus sync.Map

// aimodel表
var AimodelTable sync.Map

var BaowenMap sync.Map

//
//// 读取StructRegistry中的值
//value, ok := StructRegistry.Load("key")
//if ok {
//    // 进行类型断言
//    structValue, ok := value.(YourStructType)
//    if ok {
//        // 使用structValue
//    }
//}
//
//// 读取MacAddressStatus中的值
//value, ok = MacAddressStatus.Load("key")
//if ok {
//    // 进行类型断言
//    macStatus, ok := value.(*UpdataMacImg)
//    if ok {
//        // 使用macStatus
//    }
//}
//
//// 读取AimodelTable中的值
//value, ok = AimodelTable.Load("key")
//if ok {
//    // 进行类型断言
//    aimodel, ok := value.(Aimodel)
//    if ok {
//        // 使用aimodel
//    }
//}
//
//// 读取BaowenMap中的值
//value, ok = BaowenMap.Load("key")
//if ok {
//    // 进行类型断言
//    baowenValue, ok := value.(string)
//    if ok {
//        // 使用baowenValue
//    }
//}

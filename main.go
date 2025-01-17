package main

import (
	"imgginaimqtt/dao"
	"imgginaimqtt/disposition"
	"imgginaimqtt/routing"
	"log"
	"os"
)

var logFile *os.File

/*
    'Ammeter': 'AI_Model_1',         // 电表类型对应 AI_Model_1
    'WaterMeter': 'AI_Model_2',      // 水表类型对应 AI_Model_2
    'PressureMeter': 'AI_Model_3',   // 压力表类型对应 AI_Model_3
    'LevelMeter': 'AI_Model_4',      // 液位表类型对应 AI_Model_4
    'SewageAlarm': 'AI_Model_5',     // 污水报警类型对应 AI_Model_5
    'ControlPanel': 'AI_Model_6',    // 控制灯板类型对应 AI_Model_6
    'TemperatureMeter': 'AI_Model_7' // 温度表类型对应 AI_Model_7
};
*/

func init() {
	/*后加入nds维护*/
	// TODO 加添ai模型对应路径到表中
	dao.AimodelTable["AI_Model_1"] = dao.Aimodel{ //电表
		AimodelUrl:  "http://127.0.0.1:5000/recognize/indicator",
		AimodelName: "Ammeter_ai",
	}
	dao.AimodelTable["AI_Model_2"] = dao.Aimodel{ //水表
		AimodelUrl:  "http://127.0.0.1:5000/recognize/water-meter",
		AimodelName: "WaterMeter_ai",
	}
	dao.AimodelTable["AI_Model_3"] = dao.Aimodel{ //压力表
		AimodelUrl:  "http://127.0.0.1:5000/recognize/pressure",
		AimodelName: "PressureMeter_ai",
	}
	dao.AimodelTable["AI_Model_4"] = dao.Aimodel{ //液位表
		AimodelUrl:  "http://127.0.0.1:5000/recognize/levelmeter",
		AimodelName: "LevelMeter_ai",
	}
	dao.AimodelTable["AI_Model_5"] = dao.Aimodel{ //污水报警
		AimodelUrl:  "http://127.0.0.1:5000/recognize/indicator",
		AimodelName: "SewageAlarm_ai",
	}
	dao.AimodelTable["AI_Model_6"] = dao.Aimodel{ //控制灯板
		AimodelUrl:  "http://127.0.0.1:5000/recognize/indicator",
		AimodelName: "ControlPanel_ai",
	}
	dao.AimodelTable["AI_Model_7"] = dao.Aimodel{ //温度表
		AimodelUrl:  "http://127.0.0.1:5000/recognize/temperature",
		AimodelName: "TemperatureMeter_ai",
	}

	// 注册结构体到DAO

	dao.StructRegistry["Ammeter"] = dao.Ammeter{} //电表

	dao.StructRegistry["temphum"] = dao.TempHum{} //温湿度

	// 创建存储文件夹
	//if err := os.MkdirAll(disposition.UploadDir, os.ModePerm); err != nil {
	//	log.Fatalf("无法创建目录: %v", err)
	//}
	//if err := os.MkdirAll(disposition.AiResultsDir, os.ModePerm); err != nil {
	//	log.Fatalf("无法创建AI结果目录: %v", err)
	//}
	var err error
	// 打开或创建日志文件
	logFile, err = os.OpenFile(disposition.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}

	//写一个hashmap存字符串
	dao.BaowenMap = make(map[string]string)
}

// 主函数入口
func main() {

	defer logFile.Close()
	log.SetOutput(logFile)

	log.Println("启动服务，监听端口4399...")
	router := routing.Router()
	if err := router.Run(":4399"); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

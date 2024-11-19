package main

import (
	"imgginaimqtt/dao"
	"imgginaimqtt/disposition"
	"imgginaimqtt/routing"
	"log"
	"os"
)

var logFile *os.File

func init() {
	// 注册结构体到DAO

	dao.StructRegistry["Ammeter"] = dao.Ammeter{} //电表

	dao.StructRegistry["temphum"] = dao.TempHum{} //温湿度

	dao.StructRegistry["device"] = dao.Device{} //设备

	// 创建存储文件夹
	if err := os.MkdirAll(disposition.UploadDir, os.ModePerm); err != nil {
		log.Fatalf("无法创建目录: %v", err)
	}
	if err := os.MkdirAll(disposition.AiResultsDir, os.ModePerm); err != nil {
		log.Fatalf("无法创建AI结果目录: %v", err)
	}
	var err error
	// 打开或创建日志文件
	logFile, err = os.OpenFile(disposition.LogFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("无法打开日志文件: %v", err)
	}
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

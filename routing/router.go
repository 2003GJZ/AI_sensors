package routing

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/handler"
)

func Router() *gin.Engine {
	// 初始化路由
	router := gin.Default()
	router.Use(cors.Default()) // 允许跨域

	// 路由标记1: 获取图片列表（已废弃）
	//router.GET("/images", handler.ImagesHandler)

	// 路由标记2: 上传图片并处理AI结果（已废弃）
	//router.POST("/upload", handler.UploadHandler)

	// 路由标记3: 获取AI结果（已废弃）
	//router.GET("/ai_result/:filename", handler.AiResultHandler)

	// 路由标记4: mqtt协议支持     需要搭配mqtt转http网关使用（暂未启用）
	router.POST("/mqtt", handler.MqttHandler)

	// 路由标记5: 图片上传成功触发AI识别(ftp图片上传)                       frmo ftp to 百度ai or YOLOU
	// TODO YOLOU 识别 上传nfs path 路径
	router.POST("/upload_success", handler.UploadFtpHandler)

	//路由标记6:iot端询问是否需要图片，回包带上是否需要更新mac地址
	router.POST("/need_imagemac", handler.DeviceRequestHandle)

	//路由标记7:客户端发起mac地址更新请求,或者是否需要图片请求，触发被动下行      MAC更新完了后会触发文件删除，和redis重置
	router.POST("/update_imagemac", handler.UpdataMacImgHandler)

	//DltHttp(测试用)
	router.POST("/dlthttp", handler.DltHttp)

	//路由标记8:接收mqqt协议数据，解析为DLT645协议，数据存储在redis中，使用电表地址作为key
	router.POST("/mqttdlt645base64", handler.MqttBaes64Handler)

	// TODO 路由标记8:
	//1.DLT645-2007协议解析栈 (ok)
	// TODO 2.处理逻辑实现

	return router

}

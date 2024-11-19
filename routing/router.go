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

	// 路由标记1: 获取图片列表
	router.GET("/images", handler.ImagesHandler)

	// 路由标记2: 上传图片并处理AI结果
	router.POST("/upload", handler.UploadHandler)

	// 路由标记3: 获取AI结果
	router.GET("/ai_result/:filename", handler.AiResultHandler)

	// 路由标记4: mqtt协议支持
	router.POST("/mqtt", handler.MqttHandler)

	return router

}

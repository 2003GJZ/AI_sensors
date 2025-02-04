package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Ai_Indicator(c *gin.Context) {
	// 通过POST请求接收矩形数据
	var receivedAnnotations []Rectangle

	if err := c.ShouldBindJSON(&receivedAnnotations); err != nil {
		c.JSON(400, gin.H{"error": "无效的矩形数据"})
		return
	}

	// 将接收到的矩形数据存储到 annotations 映射中
	for _, rect := range receivedAnnotations {
		annotations[rect.Label] = rect
	}

	// 打印接收到的矩形数据（调试用）
	fmt.Println("接收到的矩形数据:", annotations)

	// 返回结果
	c.JSON(200, "存储成功")
}

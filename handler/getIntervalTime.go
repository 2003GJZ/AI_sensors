package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 定义一个包含时间字段的结构体
type TimeConfig struct {
	IntervalTime float64 `json:"IntervalTime"` // 添加 json 标签以便解析 JSON 数据
}

// 创建一个全局的 TimeConfig 实例
var timeConfig = TimeConfig{
	IntervalTime: 30.00,
}

// 使用结构体中的时间间隔
var Interval = timeConfig.IntervalTime

func GetIntervalTime(c *gin.Context) {
	var newTimeConfig TimeConfig

	// 解析请求体中的 JSON 数据到 newTimeConfig
	if err := c.ShouldBindJSON(&newTimeConfig); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	// 更新全局的 IntervalTime 变量
	timeConfig.IntervalTime = newTimeConfig.IntervalTime
	Interval = timeConfig.IntervalTime

	fmt.Println("获取到时间间隔:", Interval)

	// 示例：将时间间隔作为响应返回
	c.JSON(200, gin.H{
		"interval": Interval,
	})
}

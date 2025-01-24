package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// 定义一个包含时间字段的结构体
type TimeConfig struct {
	IntervalTime float64 // 时间间隔，单位为分钟
}

// 创建一个全局的 TimeConfig 实例
var timeConfig = TimeConfig{
	IntervalTime: 30.00,
}

// 使用结构体中的时间间隔
var Interval = timeConfig.IntervalTime

func GetIntervalTime(c *gin.Context) {

	fmt.Println("获取到时间间隔:", Interval)

	// 示例：将时间间隔作为响应返回
	c.JSON(200, gin.H{
		"interval": Interval,
	})
}

package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
)

// 定义全局变量 annotations 为 map 类型
var annotations = make(map[string]Rectangle)

type Rectangle struct {
	Label  string
	Points [4][2]float64
}

func Indicator(c *gin.Context) {
	// 通过POST请求接收数据
	var receivedData []struct {
		Coordinate []float64 `json:"coordinate"`
		Status     string    `json:"status"`
	}

	if err := c.ShouldBindJSON(&receivedData); err != nil {
		c.JSON(400, gin.H{"error": "无效的请求数据"})
		return
	}

	fmt.Println("接收到的数据:", receivedData)

	fmt.Println("接收到的矩形数据:", annotations)
	// 计算矩形的中心点并覆盖 receivedData
	result := replacePointsWithCenter(receivedData, annotations)

	// 格式化 JSON 输出
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		c.JSON(500, gin.H{"error": "无法格式化JSON数据"})
		return
	}

	// 打印格式化后的 JSON 数据
	fmt.Println(string(jsonData))

	// 返回结果
	c.JSON(200, result)
}

// 辅助函数：计算矩形的中心点
func calculateCenter(rect [4][2]float64) (float64, float64) {
	var sumX, sumY float64
	for _, point := range rect {
		sumX += point[0]
		sumY += point[1]
	}
	return sumX / 4, sumY / 4
}

// 辅助函数：替换 pointsData 中的点为矩形的中心点
func replacePointsWithCenter(pointsData []struct {
	Coordinate []float64 `json:"coordinate"`
	Status     string    `json:"status"`
}, annotations map[string]Rectangle) []map[string]string {
	var result []map[string]string

	for _, pointData := range pointsData {
		x, y := pointData.Coordinate[0], pointData.Coordinate[1]
		for _, rect := range annotations {
			if isPointInRectangle(x, y, rect.Points) {
				centerX, centerY := calculateCenter(rect.Points)
				newCoordStr := fmt.Sprintf("(%.2f, %.2f)", centerX, centerY)
				result = append(result, map[string]string{newCoordStr: pointData.Status})
				break
			}
		}
	}

	return result
}

// 辅助函数：检查点是否在矩形内
func isPointInRectangle(x, y float64, rect [4][2]float64) bool {
	// 检查点是否在矩形内的逻辑
	// 假设矩形的四个点是按顺序排列的
	minX, maxX := rect[0][0], rect[0][0]
	minY, maxY := rect[0][1], rect[0][1]
	for _, point := range rect {
		if point[0] < minX {
			minX = point[0]
		}
		if point[0] > maxX {
			maxX = point[0]
		}
		if point[1] < minY {
			minY = point[1]
		}
		if point[1] > maxY {
			maxY = point[1]
		}
	}
	if x < minX || x > maxX || y < minY || y > maxY {
		return false
	}
	// 使用射线法检查点是否在多边形内
	return isPointInPolygon(x, y, rect)
}

// 辅助函数：使用射线法检查点是否在多边形内
func isPointInPolygon(x, y float64, polygon [4][2]float64) bool {
	intersect := false
	j := len(polygon) - 1
	for i := 0; i < len(polygon); i++ {
		xi, yi := polygon[i][0], polygon[i][1]
		xj, yj := polygon[j][0], polygon[j][1]
		if ((yi > y) != (yj > y)) && (x < (xj-xi)*(y-yi)/(yj-yi)+xi) {
			intersect = !intersect
		}
		j = i
	}
	return intersect
}

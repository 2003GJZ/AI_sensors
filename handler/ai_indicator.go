package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"regexp"
	"strconv"
	"strings"
)

type Rectangle struct {
	Label  string
	Points [4][2]float64
}

type DataPoint struct {
	Coordinates [2]int
	Status      string
}

func Ai_Indicator(msg dao.Message) (string, error) {

	// 连接Redis
	link, err := mylink.GetredisLink()
	if err != nil {
		return "", fmt.Errorf("无法连接到Redis: %v", err)
	}
	defer link.Client.Close()

	//var msg dao.Message
	//if err := c.ShouldBindJSON(&msg); err != nil {
	//	respond(c, 400, "无效的请求数据", nil)
	//	return
	//}

	// 从 Redis 中获取 annotations
	var annotationsData string
	err = link.Client.HGet(link.Ctx, "indicator", msg.DeviceID).Scan(&annotationsData)
	if err != nil {
		return "", fmt.Errorf("无法获取annotations数据: %v", err)
	}

	//// 调试信息：打印 annotationsData
	//fmt.Println("annotationsData:", annotationsData)

	// 预处理 annotationsData
	processedAnnotationsData, err := preprocessAnnotationsData(annotationsData)
	fmt.Println("processedAnnotationsData:", processedAnnotationsData)
	if err != nil {
		return "", fmt.Errorf("无法预处理annotations数据: %v", err)
	}

	//// 调试信息：打印 processedAnnotationsData
	//fmt.Println("processedAnnotationsData:", processedAnnotationsData)

	// 定义全变量 annotations 为 map 类型
	var annotations = make(map[string]Rectangle)
	// 解析 annotations 数据
	err = json.Unmarshal([]byte(processedAnnotationsData), &annotations)
	if err != nil {
		return "", fmt.Errorf("无法解析annotations数据: %v", err)
	}
	fmt.Println("annotations:", annotations)

	// 解析请求中的 Data 字段
	//fmt.Println("Data:", msg.Data)
	dataPoints, err := parseDataPoints(msg.Data)
	if err != nil {
		return "", fmt.Errorf("无法解析Data数据: %v", err)
	}

	// 初始化 results
	results := make(map[string]string)
	for label := range annotations {
		results[label] = "off"
	}

	// 遍历解析后的数据点
	for _, dataPoint := range dataPoints {
		x, y := float64(dataPoint.Coordinates[0]), float64(dataPoint.Coordinates[1])
		found := false
		for label, rect := range annotations {
			if isPointInRectangle(x, y, rect.Points) {
				results[label] = dataPoint.Status
				found = true
				break
			}
		}
		if !found {
			// 如果没有找到匹配的矩形，跳过该点
			continue
		}
	}

	// 格式化 JSON 输出
	jsonData, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		return "", errors.New("无法格式化JSON数据")
	}

	// 打印格式化后的 JSON 数据
	fmt.Println(string(jsonData))

	//// 存储结果到 Redis
	//err = link.Client.HSet(link.Ctx, "ai_value", msg.DeviceID, jsonData).Err()
	//if err != nil {
	//	return "", errors.New("无法存储结果到Redis")
	//}

	// 返回结果
	return string(jsonData), nil
}

// 辅助函数：预处理 annotationsData
func preprocessAnnotationsData(data string) (string, error) {
	var annotations []map[string][][2]float64
	if err := json.Unmarshal([]byte(data), &annotations); err != nil {
		return "", fmt.Errorf("解析 annotations 数据失败: %v", err)
	}

	// 构建新的 annotationsData
	var newAnnotationsData strings.Builder
	newAnnotationsData.WriteString("{")
	firstEntry := true

	for _, annotation := range annotations {
		for label, points := range annotation {
			if len(points) != 4 {
				return "", fmt.Errorf("矩形 %s 的坐标点数量不正确: %d", label, len(points))
			}

			if !firstEntry {
				newAnnotationsData.WriteString(",")
			}
			firstEntry = false

			newAnnotationsData.WriteString(fmt.Sprintf(`"%s":{"Points":[`, label))
			for i, point := range points {
				if i > 0 {
					newAnnotationsData.WriteString(",")
				}
				newAnnotationsData.WriteString(fmt.Sprintf(`[%f,%f]`, point[0], point[1]))
			}
			newAnnotationsData.WriteString("]}")
		}
	}

	newAnnotationsData.WriteString("}")

	return newAnnotationsData.String(), nil
}

// 辅助函数：解析 Data 字段
func parseDataPoints(data string) ([]DataPoint, error) {
	fmt.Println("Data:", data)
	// 使用正则表达式匹配数据点
	re := regexp.MustCompile(`\{\((\d+),\s*(\d+)\):\s*'([A-Z]+)'`)
	matches := re.FindAllStringSubmatch(data, -1)

	if matches == nil || len(matches) == 0 {
		return nil, errors.New("no data points found")
	}

	var dataPoints []DataPoint
	for _, match := range matches {
		// 将字符串转换为整数
		x, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
		y, err := strconv.Atoi(match[2])
		if err != nil {
			return nil, err
		}

		// 创建DataPoint实例并添加到切片中
		dataPoint := DataPoint{
			Coordinates: [2]int{x, y},
			Status:      match[3],
		}
		dataPoints = append(dataPoints, dataPoint)
	}

	return dataPoints, nil
}

// 辅助函数：检查点是否在矩形内
func isPointInRectangle(x, y float64, rect [4][2]float64) bool {
	// 假设矩形的点是按顺时针或逆时针顺序给出的
	// 计算矩形的边界
	var minX, maxX, minY, maxY float64
	for _, point := range rect {
		if point[0] < minX || minX == 0 {
			minX = point[0]
		}
		if point[0] > maxX {
			maxX = point[0]
		}
		if point[1] < minY || minY == 0 {
			minY = point[1]
		}
		if point[1] > maxY {
			maxY = point[1]
		}
	}

	// 检查点是否在矩形边界内
	return x >= minX && x <= maxX && y >= minY && y <= maxY
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

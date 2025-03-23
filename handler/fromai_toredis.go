package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/disposition"
	"imgginaimqtt/mylink"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func GetAitoRedis(c *gin.Context) {
	link, _ := mylink.GetredisLink()
	//获取请求体转json

	//body, err := ioutil.ReadAll(c.Request.Body)
	//if err != nil {
	//	return
	//}

	var aiRespones dao.Message
	err := c.ShouldBindJSON(&aiRespones)
	if err != nil {
		log.Println("解析JSON失败-------------------------------->>ERR>>>>", err)
		return
	}

	//////////////////////////////////////////////临时造假
	if aiRespones.Data == "" || aiRespones.Data == "null" || aiRespones.Data == "0" {

		link.Client.HGet(link.Ctx, "ai_value", aiRespones.DeviceID).Scan(&aiRespones.Data)

		fmt.Fprintln(os.Stdout, "AI处理结果为空或null或0-------------------------------->>ERR>>>>")
	}

	//fmt.Println(aiRespones.DeviceID)
	var tableType string
	var id_value string
	link.Client.HGet(link.Ctx, "type", aiRespones.DeviceID).Scan(&tableType)
	if tableType == "" {
		log.Println("没有找到该设备对应的表-------------------------------->>ERR>>>>")
		return
	}
	//fmt.Println(tableType)
	id_value = aiRespones.DeviceID + ":" + aiRespones.Data

	logName := tableType + "_" + aiRespones.DeviceID + ".log"

	if tableType == "Indicator" {
		//fmt.Println("1233")
		result, err := Ai_Indicator(aiRespones)
		if err != nil {
			log.Println("AI处理失败-------------------------------->>ERR>>>>", err)
			return
		}

		link.Client.HSet(link.Ctx, "ai_value", aiRespones.DeviceID, result)

		// 替换所有的换行符和空格
		dataOneline := strings.ReplaceAll(result, "\n", "")
		dataOneline = strings.ReplaceAll(dataOneline, "  ", "")

		result = aiRespones.DeviceID + ":" + dataOneline
		//把ai处理结果写到文件里
		err = saveResult(dataOneline, disposition.AiResultsDir, logName)
		if err != nil {
			log.Println("保存AI处理结果到文件失败-------------------------------->>ERR>>>>", err)
			return
		}

		respond(c, 200, "redis保存成功", nil)
	} else {
		//把ai处理结果写到文件里
		err = saveResult(id_value, disposition.AiResultsDir, logName)
		if err != nil {
			log.Println("保存AI处理结果到文件失败-------------------------------->>ERR>>>>", err)
			return
		}

		// 标记5: 保存AI处理结果 到redis
		link.Client.HSet(link.Ctx, "ai_value", aiRespones.DeviceID, aiRespones.Data)

		respond(c, 200, "redis保存成功", nil)
	}

	////直接递增价格Billing
	//err = Billing(aiRespones.DeviceID)
	//if err != nil {
	//	log.Println(aiRespones.DeviceID, "直接递增价格失败-------------------------------->>ERR>>>>", err)
	//	return
	//}

	//_, err = link.Client.HIncrBy(link.Ctx, "increment_results", "ai_num", 1).Result()
	//if err != nil {
	//	log.Printf("在 Redis 中无法提高价格: %v\n", err)
	//	return
	//}
} // AI结果存储目录

// 辅助函数: 保存AI处理结果到指定文件
func saveResult(aiResult, dir, filename string) error {
	// 获取当前时间戳并格式化
	timestamp := time.Now().Format("2006-01-02 15:04:05")

	// 将时间戳添加到AI结果前面
	aiResultWithTimestamp := fmt.Sprintf("%s %s", timestamp, aiResult)

	// 标记2_5_2: 创建结果保存路径
	filePath := filepath.Join(dir, filename)
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 标记2_5_3: 写入结果
	_, err = file.WriteString(aiResultWithTimestamp + "\n") // 添加换行符以便每次写入后换行
	return err
}

package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"log"
)

func GetAitoRedis(c *gin.Context) {
	link1, _ := mylink.GetredisLink()
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
	// 标记5: 保存AI处理结果 到redis
	link1.Client.HSet(link1.Ctx, "ai_value", aiRespones.DeviceID, aiRespones.Data)

	//直接递增价格
	_, err = link1.Client.HIncrBy(link1.Ctx, "increment_results", "ai_all_num", 1).Result()
	if err != nil {
		log.Printf("在 Redis 中无法提高价格: %v\n", err)
		return
	}
}

package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/mylink"
)

func Billing(c *gin.Context) {
	monitor_id := c.Query("monitor_id")
	if monitor_id == "" {
		respond(c, 400, "monitor_id 参数缺失", nil)
		return
	}

	// 连接Redis
	link, err := mylink.GetredisLink()
	if err != nil {
		c.Error(fmt.Errorf("无法连接到Redis: %v", err))
		respond(c, 500, "无法连接到Redis", nil)
		return
	}
	defer link.Client.Close()

	// 使用加密id查询用户名
	var name string
	err = link.Client.HGet(link.Ctx, "monitor_user", monitor_id).Scan(&name)
	if err != nil {
		c.Error(fmt.Errorf("查询用户名失败: %v", err))
		respond(c, 500, "查询用户名失败", nil)
		return
	}
	if name == "" {
		respond(c, 404, "未找到用户名", nil)
		return
	}
	username := name + "_Increment"

	// 执行计费逻辑
	err = chargeUser(username, link)
	if err != nil {
		c.Error(fmt.Errorf("计费失败: %v", err))
		respond(c, 500, "计费失败", nil)
		return
	}

	respond(c, 200, "计费成功", nil)
}

func chargeUser(username string, link *mylink.RedisLink) error {
	// 更新计费信息
	_, err := link.Client.HIncrBy(link.Ctx, username, "03", 1).Result()
	if err != nil {
		return err
	}
	return nil
}

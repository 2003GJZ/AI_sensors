package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
)

// 询问是否需要上传图片
func DeviceRequestHandle(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(400, dao.ResponseEER_400("id not null"))
	}

	// 检查MAC地址是否存在

	status := dao.MacAddressStatus[id]
	if status != nil {
		//当前时间,用来做间隔时间上传
		//UpdataTime := time.Now().UnixNano()
		//|| (UpdataTime-status.LastUpdata) > int64(disposition.Interval)
		if status.NeedsImage == "1" { // 检查是否需要图片
			c.JSON(601, dao.ResponseSuccess_601())
			status.NeedsImage = "0"
		} else {
			//不需要
			c.JSON(600, dao.ResponseSuccess_600())
		}
	} else {
		//首次上传
		dao.MacAddressStatus[id] = &dao.UpdataMacImg{}
		c.JSON(601, dao.ResponseSuccess_601())

	}

}

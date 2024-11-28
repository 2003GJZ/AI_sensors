package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"net/http"
)

// 客户端通知图片上传成功
func DeviceRequestHandle(c *gin.Context) {
	var req dao.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		return
	}

	// 检查MAC地址是否存在

	status := dao.MacAddressStatus[req.MACAddress]
	if status != nil {
		//当前时间,用来做间隔时间上传
		//UpdataTime := time.Now().UnixNano()
		//|| (UpdataTime-status.LastUpdata) > int64(disposition.Interval)
		if status.NeedsImage == "1" { // 检查是否需要图片
			c.JSON(http.StatusOK, dao.ResponseSuccess_600())
			status.NeedsImage = "0"
		} else {
			//不需要
			c.JSON(http.StatusOK, dao.ResponseSuccess_601())
		}
	} else {
		//首次上传
		dao.MacAddressStatus[req.MACAddress] = &dao.UpdataMacImg{}
		c.JSON(http.StatusOK, dao.ResponseSuccess_600())

	}

}

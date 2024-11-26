package handler

import (
	"imgginaimqtt/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 通知图片上传成功和新的MAC地址
func DeviceRequestHandle(c *gin.Context) {
	var req dao.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		return
	}

	// 检查MAC地址是否存在

	status := dao.MacAddressStatus[req.MACAddress]
	if status != nil {
		// 检查是否需要更新MAC地址
		if status.NewMAC != "" && status.NewMAC != req.MACAddress {
			c.JSON(http.StatusOK, dao.ResponseSuccess_610(status.NewMAC))
			//清空需求
			status.NewMAC = ""
			dao.MacAddressStatus[req.MACAddress] = status
			return
		} else if status.NeedsImage == "1" { // 检查是否需要图片
			c.JSON(http.StatusOK, dao.ResponseSuccess_600())
			status.NeedsImage = "0"
		} else {
			c.JSON(http.StatusOK, dao.ResponseSuccess_601())
		}
	} else {
		//MAC地址不存在，即当前mac无需下行
		//使用mac地址hash后加上时间戳取模进行，概率定时上传
		// TODO 加上上次上传的时间戳

	}

}

package handler

import (
	"imgginaimqtt/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

func handleDeviceRequest(c *gin.Context) {
	var req dao.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		return
	}

	// 检查MAC地址是否存在

	status, exists := dao.MacAddressStatus[req.MACAddress]
	if exists {
		// 检查是否需要更新MAC地址
		if status.NewMAC != "" {
			c.JSON(http.StatusOK, dao.ResponseSuccess_610(status.NewMAC))
			return
		}
	} else {
		// 检查是否需要图片
		if status.NeedsImage {
			c.JSON(http.StatusOK, dao.ResponseSuccess_600())
		} else {
			c.JSON(http.StatusOK, dao.ResponseSuccess_601())
		}
	}

}

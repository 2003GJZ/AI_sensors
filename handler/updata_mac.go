package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
)

func UpdataMacImgHandler(c *gin.Context) { //被动下行接口
	var req dao.Request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		return
	}

	//首先查询
	status := dao.MacAddressStatus[req.MACAddress]
	if status == nil {
		//没有查询到
		dao.MacAddressStatus[req.MACAddress] = &req.UpdataMacImg
	} else {
		// TODO: 更新mac地址
		status.NewMAC = req.UpdataMacImg.NewMAC
		status.NeedsImage = req.UpdataMacImg.NeedsImage
	}
	c.JSON(200, dao.ResponseSuccess("Success"))

}

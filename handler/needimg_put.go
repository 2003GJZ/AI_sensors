package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
)

// 是否需要上传图片
func UpdataImgHandler(c *gin.Context) { //被动下行接口
	id := c.PostForm("id")
	imgneedimg := c.PostForm("needimg")
	if imgneedimg != "0" && imgneedimg != "1" {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		return
	}
	req := dao.Request{
		MACAddress: id,
		UpdataImg: dao.UpdataMacImg{
			NeedsImage: imgneedimg,
		},
	}

	//var req dao.Request
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	c.JSON(400, dao.ResponseEER_400("Invalid request format"))
	//	return
	//}
	//首先查询
	status := dao.MacAddressStatus[req.MACAddress]
	if status == nil {
		//没有查询到
		dao.MacAddressStatus[req.MACAddress] = &req.UpdataImg
	} else {
		status.NeedsImage = req.UpdataImg.NeedsImage
	}
	c.JSON(200, dao.ResponseSuccess("Success"))

}

//"\"indicator center list\": [[89, 131],[195,133],[301,1331, [411,133],[84,238], [192,238], [300, 241], [410, 242]],\"indicator_namelist\":[\"0\",\"1\",\"2\",\"3\",\"4\",\"5\",\"6\",\"7\"]"

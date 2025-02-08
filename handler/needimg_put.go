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

	// 获取或创建 MacAddressStatus 条目
	ptrValue, exists := dao.MacAddressStatus.Load(id)
	var ptr *dao.UpdataMacImg
	if exists {
		var ok bool
		ptr, ok = ptrValue.(*dao.UpdataMacImg)
		if !ok {
			c.JSON(400, dao.ResponseEER_400("类型断言失败"))
			return
		}
		ptr.NeedsImage = imgneedimg
	} else {
		ptr = &dao.UpdataMacImg{NeedsImage: imgneedimg}
		dao.MacAddressStatus.Store(id, ptr)
	}

	c.JSON(200, dao.ResponseSuccess("Success"))
}

package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"log"
	"time"
)

// 询问是否需要上传图片
func DeviceRequestHandle(c *gin.Context) {
	id := c.PostForm("id")
	if id == "" {
		c.JSON(400, dao.ResponseEER_400("id not null"))
	}

	// 检查MAC地址是否存在
	ptrValue, exists := dao.MacAddressStatus.Load(id)
	if exists {
		ptr, ok := ptrValue.(*dao.UpdataMacImg)
		if !ok {
			log.Printf("类型断言失败: %v 不是 *UpdataMacImg 类型", ptrValue)
			c.JSON(400, dao.ResponseEER_400("类型断言失败"))
			return
		}
		checkAndModifyNeedImage(ptr)

		if ptr.NeedsImage == "1" { // 检查是否需要图片
			c.JSON(211, dao.ResponseSuccess_211())
		} else {
			//不需要
			c.JSON(210, dao.ResponseSuccess_210())
		}
	} else {
		//首次上传询问
		newPtr := &dao.UpdataMacImg{}
		dao.MacAddressStatus.Store(id, newPtr)
		c.JSON(211, dao.ResponseSuccess_211())
	}
}

func checkAndModifyNeedImage(ptr *dao.UpdataMacImg) {
	currentTime := time.Now().UnixNano()
	timeDifference := currentTime - ptr.LastUpdata // 获取时间差（纳秒）
	// 检查时间差是否大于5秒
	if timeDifference > int64(5*time.Second) {
		ptr.NeedsImage = "1"
	} else {
		ptr.NeedsImage = "0"
	}
}

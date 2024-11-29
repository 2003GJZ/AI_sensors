package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"imgginaimqtt/tool"
	"io/ioutil"
)

// 获取水位，和是否漏水
func WaterLevelHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, dao.ResponseEER_400("err"))
		return
	}
	defer c.Request.Body.Close()
	bodystring := string(body)

	//声明一个byte数组，长度为4
	arr := make([]string, 4)
	var splitString string
	var err0 error
	for i := 0; i < 4; i++ {
		splitString, bodystring, err0 = tool.SplitString(bodystring, ",")
		if err0 != nil {
			c.JSON(400, dao.ResponseEER_400("err"))
		}
		arr[i] = splitString
	}
	//字符串拼接
	waterlevel := "waterlevel:\"" + arr[0] + arr[1] + arr[2] + "\"" + "leaking:\"" + arr[3] + "\""
	fmt.Println(waterlevel)
	link, _ := mylink.GetredisLink()
	link.Client.HSet(link.Ctx, "ai_value", "10086", waterlevel)
	c.JSON(200, dao.ResponseSuccess("ok"))

}

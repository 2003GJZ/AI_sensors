package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"imgginaimqtt/tool"
	"io/ioutil"
)

type WaterLevel struct {
	Leaking    string `json:"leaking"`
	Waterlevel string `json:"waterlevel"`
}

// 获取水位，和是否漏水
func WaterLevelHandler(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, dao.ResponseEER_400("err"))
		return
	}
	defer c.Request.Body.Close()

	bodystring := string(body)
	fmt.Println("gbk:+++++++++" + bodystring)
	if len(body) != 7 {
		c.JSON(401, dao.ResponseEER_400("err"))
		return
	}
	//fmt.Println("gbk:+++++++++" + bodystring)

	//声明一个byte数组，长度为4
	arr := make([]string, 4)
	var splitString string
	for i := 0; i < 4; i++ {
		splitString, bodystring, _ = tool.SplitString(bodystring, ",")
		arr[i] = splitString
	}
	//字符串拼接
	waterlevel := arr[0] + arr[1] + arr[2]
	fmt.Println(waterlevel)
	leaking := arr[3]
	if leaking == "0" {
		leaking = "漏水"
	} else {
		leaking = "正常"
	}
	data := WaterLevel{
		Waterlevel: waterlevel,
		Leaking:    leaking,
	}

	//转json
	json, _ := json.Marshal(data)
	fmt.Println(string(json))
	link, _ := mylink.GetredisLink()
	link.Client.HSet(link.Ctx, "ai_value", "10086", string(json))
	c.JSON(200, dao.ResponseSuccess("ok"))
	dao.NoticeUpdate("10086")

}

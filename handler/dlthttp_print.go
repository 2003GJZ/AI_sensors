package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/protocol_stack"
	"io/ioutil"
	"net/http"
)

func DltHttp(c *gin.Context) {

	// 读取请求体
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(400, dao.ResponseEER_400("err"))
		return
	}
	defer c.Request.Body.Close()
	fmt.Println("数据数组:", body)

	// 将请求体转换为16进制数组
	err1, i := protocol_stack.ElectricityAnswer(body)
	//rawFrame := []byte{0x68, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10, 0x68, 0x11, 0x04, 0x33, 0x33, 0x34, 0x35, 0xC5, 0x16}
	if err1 != nil {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		fmt.Printf("应答解析失败: %v\n", err1)
	}
	frame, err := protocol_stack.ParseDLT645Frame(i)
	if err != nil {
		c.JSON(400, dao.ResponseEER_400("Invalid request format"))
		fmt.Printf("解析失败: %v\n", err)
		return
	}
	//frame.Address//主键
	fmt.Printf("解析成功: %+v\n", frame)
	// 数据域去偏移
	decodedData := protocol_stack.OffsetData(frame.Data, false)
	fmt.Printf("去偏移后的数据域: %X\n", decodedData)
	respondWithJSON(c, http.StatusOK, "ok", nil)

}

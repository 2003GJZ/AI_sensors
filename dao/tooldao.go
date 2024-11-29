package dao

import (
	"bytes"
	"encoding/json"
	"net/http"
)

//{
/*
{
"mac":"10010",
"body":"\"needsImage\":\"0\",\"newMAC\":\"10061\"}"
}
*/

// 统一错误返回

// 一个通知函数，我去发一个post请求给http://127.0.0.1:9000/
func Notice_post(message Message) {
	//发起一个http请求到http://127.0.0.1:9000/
	//把message的body传过去
	// 将结构体编码为JSON
	jsonData, err := json.Marshal(message)
	if err != nil {
		panic(err)
	}
	http.Post("http://127.0.0.1:9000/", "application/json", bytes.NewBuffer(jsonData))

}

func ResponseEER_400(err string) *Response {
	return &Response{
		Code:    400,
		Message: err,
		Data:    nil,
	}
}

// 统一成功返回
func ResponseSuccess(data interface{}) *Response {
	return &Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

/*问我是否要图片*/
//通知改名字
func ResponseSuccess_610(newname string) *Response {
	return &Response{
		Code:    610,
		Message: newname,
		Data:    nil,
	}
}

// 需要图片
func ResponseSuccess_600() *Response {
	return &Response{
		Code:    600,
		Message: "no",
		Data:    nil,
	}
}

// 不需要图片
func ResponseSuccess_601() *Response {
	return &Response{
		Code:    601,
		Message: "yes",
		Data:    nil,
	}
}

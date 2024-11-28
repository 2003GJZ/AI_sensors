package dao

type UpdataMacImg struct {
	NeedsImage string `json:"needsImage"` // 是否需要更新图片   "1"需要   “0”不需要
	LastUpdata int64  `json:"listUpdata"` //上次更新时间
}

// 改名601
// 返回统一格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type Request struct {
	MACAddress   string       `json:"mac" binding:"required"`
	UpdataMacImg UpdataMacImg `json:"updataMacImg" binding:"required"`
}

//{
/*
{
"mac":"10010",
"body":"\"needsImage\":\"0\",\"newMAC\":\"10061\"}"
}
*/

// 统一错误返回
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
		Message: "yes",
		Data:    nil,
	}
}

// 不需要图片
func ResponseSuccess_601() *Response {
	return &Response{
		Code:    601,
		Message: "no",
		Data:    nil,
	}
}

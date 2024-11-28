package dao

// 表公共接口
type Table interface {
	GetDeviceID() string
}

// Ammeter 电表
type Ammeter struct { // 电表
	DeviceID string `json:"device_id"`
	// 设备ID
	Current string `json:"current"`
	// 电流
	Voltage string `json:"voltage"`
	// 电压
	Power string `json:"power"`
	// 功率
	ACurrent string `json:"a_current"`
	BCurrent string `json:"b_current"`
	CCurrent string `json:"c_current"`
	//ABC	相电流
	AVoltage string `json:"a_voltage"`
	BVoltage string `json:"b_voltage"`
	CVoltage string `json:"c_voltage"`
	//ABC	相电压
}

// TempHum 温湿度
type TempHum struct { // 温湿度
	DeviceID string `json:"device_id"`
	// 设备ID
	Temp string `json:"temp"`
	// 温度
	Humidity string `json:"humidity"`
	// 湿度
}

type Message struct { // 消息
	DeviceID string `json:"device_id"`
	// 设备ID
	Data string `json:"data"`
	// 消息内容
}

// 实现接口
func (a *Ammeter) GetDeviceID() string {
	return a.DeviceID
}

func (a *TempHum) GetDeviceID() string {
	return a.DeviceID
}

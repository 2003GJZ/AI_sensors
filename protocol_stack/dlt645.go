package protocol_stack

import (
	"bytes"
	"errors"
	"fmt"
)

// DLT645Frame 表示DL/T645协议的数据帧
type DLT645Frame struct {
	Address    string // 地址
	Control    byte   // 控制字节
	DataLength byte   // 数据域长度
	Data       []byte // 数据域
	Checksum   byte   // 校验码
}

const Offset byte = 0x33 // 偏移量

// ParseDLT645Frame 解析DLT645协议数据帧
func ParseDLT645Frame(rawData []byte) (*DLT645Frame, error) {
	if len(rawData) < 12 {
		return nil, errors.New("数据帧长度不足")
	}
	// 检查帧头和帧尾
	if rawData[0] != 0x68 || rawData[7] != 0x68 || rawData[len(rawData)-1] != 0x16 {
		return nil, errors.New("帧头或帧尾标识无效")
	}
	// 解析地址
	addressBytes := rawData[1:7]
	//地址反转
	address := fmt.Sprintf("%02X%02X%02X%02X%02X%02X",
		addressBytes[5], addressBytes[4], addressBytes[3],
		addressBytes[2], addressBytes[1], addressBytes[0])

	// 解析控制字节和数据长度
	control := rawData[8]
	dataLength := rawData[9]
	// 校验长度
	if len(rawData) != 10+int(dataLength)+2 {
		return nil, fmt.Errorf("数据帧长度与数据域不匹配：期望长度 %d, 实际长度 %d", 10+int(dataLength)+2, len(rawData))
	}
	// 提取数据域
	data := rawData[10 : 10+dataLength]
	// 校验码验证
	calculatedChecksum := calculateChecksum(rawData[:len(rawData)-2])
	if calculatedChecksum != rawData[len(rawData)-2] {
		return nil, fmt.Errorf("校验码错误：期望 %02X, 实际 %02X", calculatedChecksum, rawData[len(rawData)-2])
	}
	return &DLT645Frame{
		Address:    address,
		Control:    control,
		DataLength: dataLength,
		Data:       data,
		Checksum:   rawData[len(rawData)-2],
	}, nil
}

// 电表应答帧
func ElectricityAnswer(rawData []byte) (error, []byte) {
	if len(rawData) < 16 {
		return errors.New("数据帧长度不足FE"), nil
	}

	if rawData[0] != 0xFE || rawData[1] != 0xFE || rawData[2] != 0xFE || rawData[3] != 0xFE {
		return errors.New("电表响应标识错误"), nil
	}

	return nil, rawData[4:]
}

// calculateChecksum 计算校验码
func calculateChecksum(data []byte) byte {
	var sum byte
	for _, b := range data {
		sum += b
	}
	return sum
}

// OffsetData 数据加偏移或去偏移
func OffsetData(data []byte, add bool) []byte {
	result := make([]byte, len(data))
	for i, b := range data {
		if add {
			result[i] = b + Offset
		} else {
			result[i] = b - Offset
		}
	}
	return result
}

// DataIdentifier 定义数据标识结构
type DataIdentifier struct {
	Format      string // 数据格式（如 XXX.X 或 XX.XXXX）
	Length      int    // 数据长度（字节数）
	Unit        string // 单位（如 V、A）
	Description string // 数据项描述
	Type        string // 数据类型（电压、电流、有功功率等）
	Phase       string // 相位（A、B、C 或 O）
}

// 数据标识表，根据图片内容完全构建
var dataIdentifierTable = map[string]DataIdentifier{
	"02-01-01-00": {"XXX.X", 2, "V", "A相电压", "V", "A"},
	"02-01-02-00": {"XXX.X", 2, "V", "B相电压", "V", "B"},
	"02-01-03-00": {"XXX.X", 2, "V", "C相电压", "V", "C"},
	"02-01-FF-00": {"XXX.X", 2, "V", "电压数据块", "V", "O"},
	"02-02-01-00": {"XXX.XXX", 3, "A", "A相电流", "I", "A"},
	"02-02-02-00": {"XXX.XXX", 3, "A", "B相电流", "I", "B"},
	"02-02-03-00": {"XXX.XXX", 3, "A", "C相电流", "I", "C"},
	"02-02-FF-00": {"XXX.XXX", 3, "A", "电流数据块", "I", "O"},
	"02-03-00-00": {"XX.XXXX", 3, "kW", "瞬时总有功功率", "P", "O"},
	"02-03-01-00": {"XX.XXXX", 3, "kW", "瞬时A有功功率", "P", "A"},
	"02-03-02-00": {"XX.XXXX", 3, "kW", "瞬时B有功功率", "P", "B"},
	"02-03-03-00": {"XX.XXXX", 3, "kW", "瞬时C有功功率", "P", "C"},
}

// ParseDataSegment 直接解析完整的数据段
func ParseDataSegment(data []byte) (string, string, string, error) {
	// 校验数据段长度
	if len(data) < 4 {
		return "", "0", "", errors.New("数据段长度不足，无法提取数据标识")
	}

	// 提取数据标识
	diKey := fmt.Sprintf("%02X-%02X-%02X-%02X", data[3], data[2], data[1], data[0])

	// 查找数据标识表
	diInfo, exists := dataIdentifierTable[diKey]
	if !exists {
		return "", "0", "", fmt.Errorf("未知的数据标识: %s", diKey)
	}

	// 提取数据值部分
	dataValue := data[4:] // 数据标识之后的部分
	if len(dataValue) < diInfo.Length {
		return "", "0", "", fmt.Errorf("数据段长度不足，期望 %d 字节数据值，实际 %d 字节", diInfo.Length, len(dataValue))
	}

	// 解析数据值
	var value string
	switch diInfo.Length {
	case 2: // 2 字节数据
		value = fmt.Sprintf("%d%d", dataValue[0], dataValue[1])
	case 3: // 3 字节数据
		value = fmt.Sprintf("%d%d%d", dataValue[0], dataValue[1], dataValue[2])
	default:
		return "", "0", "", fmt.Errorf("不支持的数据长度: %d 字节", diInfo.Length)
	}

	//加小数点
	switch diInfo.Format {
	case "XXX.X":
		value = InsertDot(value, 3)
	case "XX.XXXX":
		value = InsertDot(value, 2)
	case "XXX.XXX":
		value = InsertDot(value, 3)
	default:

	}

	// 返回解析结果
	return diInfo.Type, value, diInfo.Phase, nil
}

// InsertDot 函数在字符串的第n个字符后面插入一个小数点
func InsertDot(s string, n int) string {
	// 检查n是否在字符串长度范围内
	if n < 1 || n > len(s) {
		return s // 如果n超出范围，返回原始字符串
	}
	// 将字符串分割为两部分，一部分是第n个字符之前，另一部分是之后
	return s[:n] + "." + s[n:]
}

/*----------------------------------------------------暂未启用------------------------------------------------------------------------*/

// BuildDLT645Frame 构建DLT645协议数据帧
func BuildDLT645Frame(address string, control byte, data []byte) ([]byte, error) {
	if len(address) != 12 {
		return nil, fmt.Errorf("地址长度必须为12位，但收到的是: %d", len(address))
	}
	// 地址高低字节翻转
	addressBytes := make([]byte, 6)
	for i := 0; i < 6; i++ {
		_, err := fmt.Sscanf(address[2*(5-i):2*(6-i)], "%02X", &addressBytes[i])
		if err != nil {
			return nil, fmt.Errorf("地址格式无效: %v", err)
		}
	}
	// 构建帧
	dataLength := byte(len(data))
	frame := &bytes.Buffer{}
	frame.WriteByte(0x68)
	frame.Write(addressBytes)
	frame.WriteByte(0x68)
	frame.WriteByte(control)
	frame.WriteByte(dataLength)
	frame.Write(data)
	checksum := calculateChecksum(frame.Bytes())
	frame.WriteByte(checksum)
	frame.WriteByte(0x16)
	return frame.Bytes(), nil
}

package main

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
		return errors.New("数据帧长度不足EF"), nil
	}

	if rawData[0] != 0xFE || rawData[1] != 0xFE || rawData[2] != 0xFE || rawData[3] != 0xFE {
		return errors.New("电表响应标识错误"), nil
	}

	return nil, rawData[4:]
}

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

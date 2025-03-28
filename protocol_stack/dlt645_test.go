package protocol_stack

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func TestParseDLT645Frame(t *testing.T) {

	//fefefefe 68 60 07 92 04 00 91 68 d1 01 35 65 16
	//fefefefe6860079204009168d101356516
	rawFrame := []byte{0xFE, 0xFE, 0xFE, 0xFE, 0x68, 0x60, 0x07, 0x92, 0x04, 0x00, 0x91, 0x68, 0xd1, 0x01, 0x35, 0x65, 0x16}
	//fefefefe6860079204009168d101356516

	//fefefefe686007920400916891083335333395335a331a16
	//rawFrame := []byte{0xFE, 0xFE, 0xFE, 0xFE, 0x68, 0x60, 0x07, 0x92, 0x04, 0x00, 0x91, 0x68, 0x91, 0x08, 0x33, 0x35, 0x33, 0x33, 0x95, 0x33, 0x5a, 0x33, 0x1a, 0x16}

	err1, i := ElectricityAnswer(rawFrame)
	//rawFrame := []byte{0x68, 0x01, 0x00, 0x00, 0x00, 0x00, 0x10, 0x68, 0x11, 0x04, 0x33, 0x33, 0x34, 0x35, 0xC5, 0x16}
	// FE FE FE FE 68 60 07 92 04 00 91 68 91 08 33 33 33 33 97 35 A8 33 6A 16
	if err1 != nil {
		fmt.Printf("应答解析失败: %v\n", err1)
	}
	frame, err := ParseDLT645Frame(i)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
		return
	}
	fmt.Printf("解析成功: %+v\n", frame)
	// 数据域去偏移
	decodedData := OffsetData(frame.Data, false)
	fmt.Printf("去偏移后的数据域: %X\n", decodedData)

	// 示例完整数据段：数据标识 + 数据值
	//dataSegment := []byte{0x02, 0x01, 0x00, 0x00, 0x12, 0x34} // DI 为 02-01-00-00，数据值为 0x3412

	// 调用解析函数
	dataType, value, phase, err := ParseDataSegment(decodedData)
	if err != nil {
		fmt.Printf("解析失败: %v\n", err)
		return
	}

	fmt.Printf("解析结果: 类型 = %s, 值 = %s, 相位 = %s\n", dataType, value, phase)

}

func TestBuildDLT645Frame(t *testing.T) {
	// 示例生成
	newFrame, err := BuildDLT645Frame("000000000010", 0x11, OffsetData([]byte{0x00, 0x00, 0x01, 0x02}, true))
	if err != nil {
		fmt.Printf("生成失败: %v\n", err)
		return
	}
	fmt.Printf("生成的新数据帧: %X\n", newFrame)
}
func TestParseDataSegment(t *testing.T) {

	hexKey := "02-01-01-00"
	byteArray, err := HexKeyToByteArray(hexKey)
	if err != nil {
		fmt.Printf("转换失败: %v\n", err)
	} else {
		// 使用 %02X 格式化每个字节为两位十六进制数，并用空格分隔
		fmt.Printf("转换成功: %02X %02X %02X %02X\n", byteArray[0], byteArray[1], byteArray[2], byteArray[3])
	}

	// 地址
	address := "000000000001"
	// 控制码
	control := byte(0x11)
	// 数据域
	data := byteArray

	// 构建标准DLT645帧
	frame, err := BuildDLT645Frame(address, control, data)
	if err != nil {
		fmt.Printf("生成失败: %v\n", err)
		return
	}

	// 打印生成的帧
	fmt.Printf("生成的新数据帧: %X\n", frame)

}

// HexKeyToByteArrayWithOffset 将字符串键转换为字节数组，并对每个字节加上偏移量 0x33
func HexKeyToByteArray(hexKey string) ([]byte, error) {
	// 分割字符串
	parts := strings.Split(hexKey, "-")
	if len(parts) != 4 {
		return nil, fmt.Errorf("无效的 hexKey 格式: %s", hexKey)
	}

	// 初始化字节数组
	byteArray := make([]byte, 4)

	// 转换每个部分为字节并加上偏移量 0x33
	for i, part := range parts {
		value, err := strconv.ParseUint(part, 16, 8)
		if err != nil {
			return nil, fmt.Errorf("解析部分 %s 失败: %v", part, err)
		}
		byteArray[i] = byte(value) + 0x33
	}

	return byteArray, nil
}

func TestBuildDLT645Frame1(t *testing.T) {
	// 示例生成
	newFrame, err := BuildDLT645Frame("910004920760", 0x11, OffsetData([]byte{0x00, 0x02, 0x00, 0x00}, true))
	if err != nil {
		fmt.Printf("生成失败: %v\n", err)
		return
	}
	fmt.Printf("生成的新数据帧: %X\n", newFrame)
}

func TestBuildDLT645Frame2(t *testing.T) {
	//生成bytes数组
	bytes := []byte{0x68, 0x60, 0x07, 0x92, 0x04, 0x00, 0x91, 0x68, 0x11, 0x04, 0x33, 0x33, 0x33, 0x33, 0x33, 0x3F, 0x16}
	//bytes := "68600792040091681104333333333F16"
	dlt645Frame, err := ParseDLT645Frame(bytes)
	fmt.Println(dlt645Frame)

	if err != nil {
		fmt.Println(err)
		return
	}

	//fmt.Printf("解析结果: %X\n", i)

}

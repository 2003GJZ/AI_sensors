package protocol_stack

import (
	"bytes"
	"fmt"
	"testing"
)

func TestBase64(t *testing.T) {
	//rawFrame := []byte{0x68, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68, 0x11, 0x04, 0x35, 0x34, 0x34, 0x33, 0xAF, 0x16}
	//FEFEFEFE6802000000000068110435343433AF16
	//rawFrame := []byte{0xfe, 0xfe, 0xfe, 0xfe, 0x68, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68, 0x11, 0x04, 0x35, 0x34, 0x34, 0x33, 0xAF, 0x16}
	//FEFEFEFE6802000000000068110433333333AB16
	//rawFrame := []byte{0xfe, 0xfe, 0xfe, 0xfe, 0x68, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68, 0x11, 0x04, 0x33, 0x33, 0x33, 0x33, 0xAB, 0x16}
	//FEFEFEFE6801000000000068110435343433AE16
	//FEFEFEFE68600792040091681104333333333F16
	rawFrame := []byte{0xfe, 0xfe, 0xfe, 0xfe, 0x68, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x68, 0x11, 0x04, 0x35, 0x34, 0x34, 0x33, 0xAE, 0x16}
	base64 := MyBytesToBase64(rawFrame)
	fmt.Println("原始数据:", rawFrame)
	fmt.Println("Base64编码:", base64)
	fmt.Println("///////////////")

	// 解码
	decodedBody, err := MyBase64ToBytes(base64)
	if err != nil {
		t.Errorf("解码失败: %v", err)
		return
	}
	fmt.Printf("解码成功，原始Base64: %s, 解码后内容: %v\n", base64, decodedBody)

	// 验证解码后的数据是否与原始数据一致
	if !bytes.Equal(decodedBody, rawFrame) {
		t.Errorf("解码后的数据与原始数据不一致。期望: %v, 实际: %v", rawFrame, decodedBody)
	}
}

// TestMyBase64ToBytes 测试 MyBase64ToBytes 函数
func TestMyBase64ToBytes(t *testing.T) {
	body := "SGVsbG8gV29ybGQh" // Base64 编码的 "Hello World!"
	expected := []byte("Hello World!")
	decodedBody, err := MyBase64ToBytes(body)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !bytes.Equal(decodedBody, expected) {
		t.Errorf("解码结果不正确。期望: %v, 实际: %v", expected, decodedBody)
	}
	fmt.Println("数据数组 b:", decodedBody)
}

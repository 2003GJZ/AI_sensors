package mylink

import (
	"testing"
)

// 测试链接
func TestLink(t *testing.T) {
	link, _ := NewRedisLink(0)
	var data string
	link.Client.HGet(link.Ctx, "test", "test").Scan(&data)
	t.Log(data)
}

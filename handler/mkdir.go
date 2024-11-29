package handler

import (
	"github.com/gin-gonic/gin"
	"os"
)

// CreateDir 创建指定路径的文件夹，如果文件夹已存在则不进行任何操作
func CreateDir(dirName string) error {
	// 检查文件夹是否已经存在
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		// 尝试创建文件夹
		return os.MkdirAll(dirName, 0755)
	}
	// 如果文件夹已经存在，返回nil表示没有错误
	return nil
}

func MkdirHandler(c *gin.Context) {

}

package handler

import (
	"github.com/gin-gonic/gin"
	"imgginaimqtt/disposition"
	"net/http"
	"os"
	"path/filepath"
)

// 路由1的处理器: 获取图片列表
// 获取上传目录中的所有图片文件路径并返回
func ImagesHandler(c *gin.Context) {
	// 标记1_1: 获取图片文件列表
	files, err := getAllImageFiles(disposition.UploadDir)
	if err != nil {
		// 标记1_1_1: 获取文件失败时返回错误
		respondWithJSON(c, http.StatusInternalServerError, "获取文件列表失败", nil)
		return
	}
	// 标记1_1_2: 获取成功时返回文件列表
	respondWithJSON(c, http.StatusOK, "获取成功", files)

}

// 辅助函数: 遍历目录获取所有文件
func getAllImageFiles(dir string) ([]string, error) {
	var files []string
	// 标记1_1_2_1: 遍历目录中的文件
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			// 仅返回文件名
			files = append(files, info.Name())
		}
		return nil
	})
	return files, err
}

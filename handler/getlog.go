package handler

import (
	"bufio"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/disposition"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// GetLogHandler 处理获取日志文件最后10行的请求
func GetLogHandler(c *gin.Context) {
	// 获取文件名参数
	var aiRespones dao.Message
	err := c.ShouldBindJSON(&aiRespones)
	if err != nil {
		log.Println("解析JSON失败-------------------------------->>ERR>>>>", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "解析JSON失败"})
		return
	}
	deviceID := aiRespones.DeviceID
	devicetype := aiRespones.Data
	var logFilePath string
	// devicetype, "_", deviceID, ".log"
	logFileName := devicetype + "_" + strings.ReplaceAll(deviceID, "/", "_") + ".log"
	logFilePath = filepath.Join(disposition.AiResultsDir, logFileName)

	fmt.Println(logFilePath)

	// 读取文件的最后10行
	lastLines, err := readLastNLines(logFilePath, 20)
	if err != nil {
		log.Println("读取文件失败-------------------------------->>ERR>>>>", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "读取文件失败"})
		return
	}

	// 将最后10行内容写入HTTP响应
	for _, line := range lastLines {
		fmt.Fprintln(c.Writer, line)
	}
}

// readLastNLines 读取文件的最后n行
func readLastNLines(filePath string, n int) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > n {
			lines = lines[1:]
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/disposition"
	"imgginaimqtt/mylink"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// 路由5的处理器: 读取最新图片并发送给AI服务器
func UploadFtpHandler(c *gin.Context) {
	// 标记1: 获取mac字段
	mac := c.PostForm("mac")
	if mac == "" {
		c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("缺少mac字段"))
		return
	}

	// 标记2: 构造文件夹路径
	dirPath := filepath.Join("/var/ftp/ftpuser", mac)
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("目录不存在"))
		return
	}

	// 标记3: 获取最新图片文件
	latestImage, err := getLatestFile(dirPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("获取最新图片失败"))
		return
	}

	link, _ := mylink.GetredisLink()
	//查询redis选择表类型
	var tableType string
	var aimodel string
	link.Client.HGet(link.Ctx, "type", mac).Scan(&tableType)

	if tableType != "" {
		link.Client.HGet(link.Ctx, "aiModel", tableType).Scan(&aimodel)
	} else {
		respondWithJSON(c, http.StatusInternalServerError, "未查询到表类型", nil)
		return
	}
	//选择ai摸型
	switch aimodel {
	case "aimodel1":
		// 标记4: 发送图片给AI服务器
		aiResult, err := sendImageToAIServer(latestImage)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("AI处理失败"))
			return
		}
		err, s := ParseJson(aiResult)
		if err != nil {
			c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("AI处理失败"))
			return
		}
		// 标记5: 保存AI处理结果 到redis
		link.Client.HSet(link.Ctx, "ai_value", mac, s)
		respondWithJSON(c, http.StatusOK, "AI处理成功", nil)

	case "aimodel2":

		//TODO YOLOU AI
	default:
		fmt.Println("未查询到ai模型")
		respondWithJSON(c, http.StatusInternalServerError, "未查询到ai模型", nil)

	}

}

// FileInfoWithTime 包含文件信息和修改时间
type FileInfoWithTime struct {
	FileInfo os.FileInfo
	ModTime  time.Time
}

// getLatestFile 获取指定目录下最新修改的文件
// getLatestFile 获取指定目录下最新修改的文件，同时控制文件数量
func getLatestFile(dirPath string) (string, error) {
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return "", err
	}

	// 筛选非目录文件
	var filesWithTime []os.FileInfo
	for _, file := range files {
		if !file.IsDir() {
			filesWithTime = append(filesWithTime, file)
		}
	}

	if len(filesWithTime) == 0 {
		return "", fmt.Errorf("no file found in directory")
	}

	// 按修改时间升序排序
	sort.Slice(filesWithTime, func(i, j int) bool {
		return filesWithTime[i].ModTime().Before(filesWithTime[j].ModTime())
	})

	// 最新文件是排序后的最后一个
	latestFile := filesWithTime[len(filesWithTime)-1]

	// 如果文件数量超过限制，删除最旧的文件
	if len(filesWithTime) > disposition.MaxfileImg {
		oldestFile := filesWithTime[0]
		oldestFilePath := filepath.Join(dirPath, oldestFile.Name())
		if err := os.Remove(oldestFilePath); err != nil {
			return "", fmt.Errorf("failed to delete file %s: %v", oldestFilePath, err)
		}
		fmt.Printf("Deleted oldest file: %s\n", oldestFile.Name())
	}

	return latestFile.Name(), nil
}

func ParseJson(body string) (error, string) { //提取百度ai返回的json数据
	//3_4 解析JSON
	var meterResp MeterResponse
	err := json.Unmarshal([]byte(body), &meterResp)
	if err != nil {
		return err, ""
	}

	// 3_5检查是否有"words"字段
	if meterResp.WordsNum == 0 {
		return errors.New("AI识别结果为空"), ""
	}

	// 3_6提取"words"字段并用逗号分隔
	var words []string
	for _, item := range meterResp.WordsList {
		words = append(words, item.Words)
	}
	resultText := strings.Join(words, ",")
	return nil, resultText
}

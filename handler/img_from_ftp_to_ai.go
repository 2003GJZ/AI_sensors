package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/disposition"
	"imgginaimqtt/mylink"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type imgai struct {
	Img_path string `json:"img_path"`
	Data     string `json:"data"`
}

// 上传完图片通知
// 路由5的处理器: 读取最新图片并发送给AI服务器
func UploadFtpHandler(c *gin.Context) {

	// 标记1: 获
	//取mac字段
	//用结构体接收
	id := c.PostForm("id")
	imgname := c.PostForm("imgname")

	if id == "" || imgname == "" {
		c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("not mac error"))
		return
	}
	//var req dao.Request
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	c.JSON(400, dao.ResponseEER_400("Invalid request format"))
	//	return
	//}
	// 标记2: 构造文件夹路径
	ptr := dao.MacAddressStatus[id]
	lastUpdataTime := time.Now().UnixNano()

	//填充当前时间
	if ptr == nil {
		dao.MacAddressStatus[id] = &dao.UpdataMacImg{LastUpdata: lastUpdataTime}
	} else {
		ptr.LastUpdata = lastUpdataTime
	}

	//检查文件夹是否存在
	//dirPath := filepath.Join(disposition.FtpPathex, id)
	//if _, err := os.Stat(dirPath); os.IsNotExist(err) {
	//	c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("path not found"))
	//	return
	//}

	// 标记3: 获取最新图片文件
	//latestImage, err := getLatestFile(dirPath)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("img not found"))
	//	return
	//}

	//根据发过来的图片path获取图片

	link, _ := mylink.GetredisLink()
	//查询redis选择表类型
	var tableType string
	var aimodelName string
	link.Client.HGet(link.Ctx, "type", id).Scan(&tableType)

	//查询redis选择ai模型
	if tableType != "" {
		link.Client.HGet(link.Ctx, "aiModel", tableType).Scan(&aimodelName)
	} else {
		respondWithJSON(c, http.StatusInternalServerError, "not tableType", nil)
		return
	}

	//ch := make(chan int)
	//启动匿名函数发送给ai
	go func() {
		link1, _ := mylink.GetredisLink()
		//选择ai摸型
		//查询AI地址
		status := dao.MacAddressStatus[id]
		aimodel, ok := dao.AimodelTable[aimodelName]

		if !ok {
			//没有此模型
			log.Println(aimodelName, "没有此模型-------------------------------->>ERR>>>>")
			return
		}
		// 标记4: 发送图片给AI服务器,

		//参数”img“
		imgmaseg := imgai{
			Img_path: "",
			Data:     "",
		}
		imgpath := disposition.FtpPathex + "/" + id + "/" + imgname

		// imgpath := "D:\\\\PythonProject\\\\Project\\\\lige\\\\PythonProject2\\\\static\\\\indicator\\\\H1L12A.jpg"
		// imgmaseg.Data = "{\"indicator_center_list\": [[89, 131], [195, 133], [301, 133], [411, 133], [84, 238], [192, 238], [300, 241], [410, 242]], \"indicator_namelist\": [\"0\", \"1\", \"2\", \"3\", \"4\", \"5\", \"6\", \"7\"]}"

		//参数body
		imgmaseg.Img_path = imgpath
		var tabe string
		link1.Client.HGet(link.Ctx, "imgRes", aimodelName).Scan(&tabe)
		imgmaseg.Data = tabe

		//response, err2 := http.Post(aimodel.AimodelUrl, "application/json", strings.NewReader(imgpath))
		//打印请求
		fmt.Println("请求地址：", aimodel.AimodelUrl)
		fmt.Println("请求参数：", imgpath)
		//转json
		jsonData, _ := json.Marshal(imgmaseg)
		fmt.Println("请求参数：", string(jsonData))

		response, err2 := http.Post(aimodel.AimodelUrl, "application/json", bytes.NewBuffer(jsonData))
		if err2 != nil {
			//ch <- 600
			//ch <- 600
			log.Println(aimodel.AimodelName, "请求地址错误-------------------------------->>ERR>>>>", err2)
			return
		}
		//判断code
		if response.StatusCode != http.StatusOK {
			//识别失败继续上传
			if status != nil {
				status.NeedsImage = "1"
			}
			//ch <- 601
			//ch <- 601
		} else {
			//读出返回的数据
			body, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return
			}
			// 标记5: 保存AI处理结果 到redis
			link1.Client.HSet(link1.Ctx, "ai_value", id, string(body))
			//ch <- 600
			//ch <- 600
			// 标记6: 通知前端
			dao.NoticeUpdate(id)
		}

	}()
	// 标记7: 返回成功响应 ，不需要图片
	c.JSON(210, dao.ResponseSuccess_210())
	return
	//code := 600
	////等待20ns
	//select {
	//case <-ch:
	//	code = <-ch
	//	if code == 601 {
	//		c.JSON(601, dao.ResponseSuccess_601())
	//	} else {
	//		c.JSON(600, dao.ResponseSuccess_600())
	//	}
	//	return
	//case <-time.After(200 * time.Second):
	//	c.JSON(600, dao.ResponseSuccess_600())
	//}

}

/*********************************************废弃**************************************************************************/
// FileInfoWithTime 包含文件信息和修改时间
type FileInfoWithTime struct {
	FileInfo os.FileInfo
	ModTime  time.Time
}

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
		log.Printf("Deleted oldest file: %s\n", oldestFile.Name())
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

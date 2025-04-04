package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"imgginaimqtt/dao"
	"imgginaimqtt/mylink"
	"io/ioutil"
	"log"
	"net/http"
)

type imgai struct {
	Id       string `json:"id"`
	Img_path string `json:"img_path"`
}

// 结构体来接收请求参数
type ImageRequest struct {
	ID      string `form:"id"`
	ImgName string `form:"imgname"`
}

//var ptr *dao.UpdataMacImg

// 上传完图片通知
// 路由5的处理器: 读取最新图片并发送给AI服务器
// 使用map[string]interface{}接收请求参数
func UploadFtpHandler(c *gin.Context) {
	link, err := mylink.GetredisLink()
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("Redis connection failed"))
		return
	}
	defer link.Client.Close()

	// 使用map[string]interface{}接收请求参数
	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 现在你可以通过 req["id"] 和 req["imgname"] 来访问这些数据
	id, ok1 := req["id"].(string)
	imgname, ok2 := req["imgname"].(string)

	if !ok1 || !ok2 || id == "" || imgname == "" {
		c.JSON(http.StatusInternalServerError, dao.ResponseEER_400("not mac error"))
		return
	}

	// 标记5: 保存图片路径
	link.Client.HSet(link.Ctx, "Image", id, imgname)

	// 存储url
	var tableType string
	var aimodelName string
	// 从哈希表查询模型url
	if value, ok := dao.AimodelTable.Load(id); ok {
		aimodelName = value.(string)
	} else {
		// 处理未找到的情况
		log.Println("未找到模型URL，ID为:", id)

		// 查询redis选择url
		link.Client.HGet(link.Ctx, "type", id).Scan(&tableType)

		if tableType == "" {
			log.Println("没有找到该设备对应的表-------------------------------->>ERR>>>>，ID为：", id)
			respond(c, 400, "not tableType", nil)
			return
		}

		// 查询redis选择ai的url
		if tableType != "" {
			link.Client.HGet(link.Ctx, "aiModel", tableType).Scan(&aimodelName)
		} else {
			respond(c, http.StatusInternalServerError, "not tableType", nil)
			return
		}
		// 存进哈希表
		dao.AimodelTable.Store(id, aimodelName)
	}

	go func() {
		// 参数”img“
		imgmaseg := imgai{
			Id:       id,
			Img_path: imgname,
		}

		// 打印请求
		fmt.Println("请求地址：", aimodelName)
		fmt.Println("请求参数：", imgname)

		// 转json
		jsonData, err := json.Marshal(imgmaseg)
		if err != nil {
			log.Println("JSON序列化错误:", err)
			return
		}
		fmt.Println("请求参数：", string(jsonData))

		response, err := http.Post(aimodelName, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Println(aimodelName, "请求地址错误-------------------------------->>ERR>>>>", err)
			return
		}
		defer response.Body.Close()

		// 判断code
		if response.StatusCode != http.StatusOK {
			log.Println(aimodelName, "识别失败-------------------------------->>ERR>>>>", err)
		} else {
			// 读出返回的数据
			_, err := ioutil.ReadAll(response.Body)
			if err != nil {
				log.Println("读取响应错误:", err)
				return
			}
		}
	}()

	// 标记7: 返回成功响应 ，不需要图片
	c.JSON(210, dao.ResponseSuccess_210())
	return
}

/*********************************************废弃**************************************************************************/
//// FileInfoWithTime 包含文件信息和修改时间
//type FileInfoWithTime struct {
//	FileInfo os.FileInfo
//	ModTime  time.Time
//}
//
//// getLatestFile 获取指定目录下最新修改的文件，同时控制文件数量
//func getLatestFile(dirPath string) (string, error) {
//	files, err := ioutil.ReadDir(dirPath)
//	if err != nil {
//		return "", err
//	}
//
//	// 筛选非目录文件
//	var filesWithTime []os.FileInfo
//	for _, file := range files {
//		if !file.IsDir() {
//			filesWithTime = append(filesWithTime, file)
//		}
//	}
//
//	if len(filesWithTime) == 0 {
//		return "", fmt.Errorf("no file found in directory")
//	}
//
//	// 按修改时间升序排序
//	sort.Slice(filesWithTime, func(i, j int) bool {
//		return filesWithTime[i].ModTime().Before(filesWithTime[j].ModTime())
//	})
//
//	// 最新文件是排序后的最后一个
//	latestFile := filesWithTime[len(filesWithTime)-1]
//
//	// 如果文件数量超过限制，删除最旧的文件
//	if len(filesWithTime) > disposition.MaxfileImg {
//		oldestFile := filesWithTime[0]
//		oldestFilePath := filepath.Join(dirPath, oldestFile.Name())
//		if err := os.Remove(oldestFilePath); err != nil {
//			return "", fmt.Errorf("failed to delete file %s: %v", oldestFilePath, err)
//		}
//		log.Printf("Deleted oldest file: %s\n", oldestFile.Name())
//	}
//
//	return latestFile.Name(), nil
//}

//func ParseJson(body string) (error, string) { //提取百度ai返回的json数据
//	//3_4 解析JSON
//	var meterResp MeterResponse
//	err := json.Unmarshal([]byte(body), &meterResp)
//	if err != nil {
//		return err, ""
//	}
//
//	// 3_5检查是否有"words"字段
//	if meterResp.WordsNum == 0 {
//		return errors.New("AI识别结果为空"), ""
//	}
//
//	// 3_6提取"words"字段并用逗号分隔
//	var words []string
//	for _, item := range meterResp.WordsList {
//		words = append(words, item.Words)
//	}
//	resultText := strings.Join(words, ",")
//	return nil, resultText
//}

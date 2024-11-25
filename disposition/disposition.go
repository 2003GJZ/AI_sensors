package disposition

// 配置文件参数

var UploadDir = "/opt/var/www/images"               // 图片存储目录
var AiResultsDir = "/opt/var/www/ai_results"        // AI结果存储目录
var ServerHost = "http://120.46.0.42:4398"          // 图片服务器地址
var LogFilePath = "/opt/images_logfile.log"         // 日志文件路径
var API_KEY = "w80tNJGJfeIMFM07Rv4tG1Be"            // 百度API_KEY
var SECRET_KEY = "OK0aqjJ5dRIOeoPFInsrmx8EHTpDUUid" // SECRET_KEY
var TokenURL = "https://aip.baidubce.com/oauth/2.0/token"
var OcrURL = "https://aip.baidubce.com/rest/2.0/ocr/v1/meter"
var MaxfileImg = 5

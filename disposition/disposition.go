package disposition

/*后通过配置文件写入*/
// 配置文件参数

var LogFilePath = "./images_logfile.log" // 日志文件路径

/*百度-----弃用*/
var UploadDir = "/var/ftp/ftpuser"                  // 图片存储目录
var AiResultsDir = "/opt/var/www/ai_results"        // AI结果存储目录
var ServerHost = "http://120.46.0.42:4398"          // 图片服务器地址
var API_KEY = "w80tNJGJfeIMFM07Rv4tG1Be"            // 百度API_KEY
var SECRET_KEY = "OK0aqjJ5dRIOeoPFInsrmx8EHTpDUUid" // SECRET_KEY
var TokenURL = "https://aip.baidubce.com/oauth/2.0/token"
var OcrURL = "https://aip.baidubce.com/rest/2.0/ocr/v1/meter"

/*文件控制-----弃用*/
var MaxfileImg = 5      //文件数目
var Interval = 60000000 //间隔时间us

var FtpPathex = "D:\\var\\ftp\\ftpuser" //ftp目录
// /var/ftp/ftpuser/9981
var NoticeUpdataUrl = "http://127.0.0.1:3000/api/monitor/updataAiRes" //ftp地址 通知前端哪个IP处理结果完成了去redis拿

var RedisAddr = "127.0.0.1:6379"
var RedisDB = 0
var RedisPassword = ""

package mockservice

/// 配置文件
const ConfigFile = "config.yml"

/// 默认配置的项目值
const defaultPort = "9527"         /// 默认端口
const defaultLogFile = "mock.log"  /// 默认日志文件
const defaultHostFile = "host.yml" /// 默认目标主机文件
const defaultInfoFile = "info.yml" /// 默认模拟服务信息文件

/// 默认日志等级
/// -1: DEBUG
/// 0: INFO
/// 1: WARN
/// 2: ERROR
/// 3: FATAL
const defaultLogLevel = 0 // 默认为INFO

/// 内部配置项
const pathRequest = "request"   // 请求信息存放相对根目录(程序运行目录下的目录)
const pathResponse = "response" // 响应信息存放相对根目录(程序运行目录下的目录)
const pathBackup = "backup"     // 备份目录

const pathResponseHeader = "response_header"                   // 响应信息存放相对根目录(程序运行目录下的目录)
const fileResponseHeader = "response_header.json"              // 响应头信息文件
const fileCommonResponseHeader = "common_response_header.json" // 模拟服务通用响应头信息文件

/// HTTP请求方法
const httpMethodGet = "GET"
const httpMethodPost = "POST"
const httpMethodPut = "PUT"
const httpMethodDelete = "DELETE"
const httpMethodHead = "HEAD"
const httpMethodCopy = "COPY"
const httpMethodView = "VIEW"

/// HTTP头部键
const httpHeaderContentEncoding = "Content-Encoding"
const httpHeaderContentType = "Content-Type"
const httpHeaderContentLength = "Content-Length"

/// HTTP头部值
const encodingGzip = "gzip"
const contentTypeJSON = "application/json"
const contentTypeJSONUTF8 = "application/json;charset=UTF-8"

/// 日志
const logFieldError = "错误信息"
const logFieldFile = "文件"
const logFieldPath = "目录"
const logFieldURL = "URL"
const logFieldHTTPMethod = "HTTP请求方法"
const logFieldHost = "主机"
const logFieldRequestHeader = "请求头"
const logFieldRequestBody = "请求体"
const logFieldResponseHeader = "响应头"
const logFieldResponseBody = "响应体"
const logFieldMockStatusCode = "模拟服务响应状态码"
const logFieldMockResponseFile = "模拟服务响应文件"
const logFieldMockResponseBody = "模拟服务响应体"

/// 程序所需文件(Mac开发用)
// const ConfigFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/mockservice/command/config.yml"
// const HostFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/mockservice/command/host.yml"
// const InfoFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/mockservice/command/input.yml"

/// 程序所需文件(Win开发用)
// const ConfigFile = "E:/develop/github.com/bettersun/mockservice/command/config.yml"
// const HostFile = "E:/develop/github.com/bettersun/mockservice/command/host.yml"
// const InfoFile = "E:/develop/github.com/bettersun/mockservice/command/input.yml"

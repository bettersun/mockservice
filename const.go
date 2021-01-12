package mockservice

/// 配置文件
const ConfigFile = "config.yml"

/// 默认端口
const defaultPort = "9527"

/// 目标主机文件
const HostFile = "host.yml"

/// 输入文件
const InfoFile = "input.yml"

/// URL响应头部信息文件
const ResponseHeaderFile = "response_header.yml"

/// 请求目录
const requestPath = "/request"

/// 响应目录
const responsePath = "/response"

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

/// 日志文件
const logFile = "run.log"

/// 日志
const logFieldError = "错误信息"
const logFieldFile = "文件"
const logFieldURI = "URI"
const logFieldHTTPMethod = "HTTP请求方法"
const logFieldHost = "主机"
const logFieldRequestHeader = "请求头"
const logFieldRequestBody = "请求体"
const logFieldResponseHeader = "响应头"
const logFieldResponseBody = "响应体"

/// 配置文件(开发用)
// const ConfigFile = "E:/develop/github.com/bettersun/mockservice/command/config.yml"
// const ConfigFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/mockservice/command/config.yml"

/// 目标主机文件(开发用)
// const HostFile = "E:/develop/github.com/bettersun/mockservice/command/host.yml"
// const HostFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/mockservice/command/host.yml"

/// 输入文件(开发用)
// const InfoFile = "E:/develop/github.com/bettersun/mockservice/command/input.yml"
// const InfoFile = "/Users/sunjiashu/Documents/Develop/github.com/bettersun/mockservice/command/input.yml"

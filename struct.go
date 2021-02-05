package mockservice

/// 配置
type Config struct {
	Port     string `yaml:"port"`     // 端口
	LogFile  string `yaml:"logFile"`  // 日志文件
	LogLevel int    `yaml:"logLevel"` // 日志等级
	HostFile string `yaml:"hostFile"` // 目标主机文件
	InfoFile string `yaml:"infoFile"` // 模拟服务信息文件

	// 使用模拟服务通用响应头
	// 仅当URL对应的响应头不存在时使用
	UseMockCommonResponseHeader bool `yaml:"useMockCommonResponseHeader"`
}

/// 模拟服务信息
type MockServiceInfo struct {
	URL                  string `yaml:"url" json:"url"`                                   // URL
	Method               string `yaml:"method" json:"method"`                             // HTTP请求方法
	TargetHost           string `yaml:"targetHost" json:"targetHost"`                     // 目标主机
	UseDefaultTargetHost bool   `yaml:"useDefaultTargetHost" json:"useDefaultTargetHost"` // 使用默认目标主机
	UseMockService       bool   `yaml:"useMockService" json:"useMockService"`             // 使用模拟服务
	StatusCode           int    `yaml:"statusCode" json:"statusCode"`                     // 响应状态码
	ResponseFile         string `yaml:"responseFile" json:"responseFile"`                 // 响应文件
	Description          string `yaml:"description" json:"description"`                   //  说明
}

/// 响应头部
type ResponseHeader struct {
	URL    string              `yaml:"url" json:"url"`       // URL
	Method string              `yaml:"method" json:"method"` // HTTP请求方法
	Header map[string][]string `yaml:"header" json:"header"` // 响应头部
}

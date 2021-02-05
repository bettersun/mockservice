package mockservice

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/bettersun/moist"
	yml "github.com/bettersun/moist/yaml"
	"github.com/sirupsen/logrus"
)

/// 全局变量：日志记录器
var logger = logrus.New()

/// 初始化
func init() {
	// 配置文件
	cfgFile := fmt.Sprintf("%v/%v", moist.CurrentDir(), ConfigFile)

	// 配置文件不存在
	if !moist.IsExist(cfgFile) {
		// 默认配置
		cfg := defaultConfig()
		// 地址设给全局变量
		config = &cfg

		// 初始化日志配置(默认)
		initLogConfig(cfg.LogFile, cfg.LogLevel)

		// 输出日志
		logger.WithFields(logrus.Fields{
			"端口":       cfg.Port,
			"日志文件":     cfg.LogFile,
			"日志等级":     cfg.LogLevel,
			"目标主机文件":   cfg.LogFile,
			"模拟服务信息文件": cfg.LogLevel,
		}).Warn("配置文件不存在，使用默认配置")

		// Go的init()方法里向Flutter端发送消息，Flutter端接收不到。
		// // 通知Flutter
		// msg := fmt.Sprintf("配置文件[%v]不存在，使用默认端口[%v]", cfgFile, defaultPort)
		// go Notify(msg)
	}

	// 配置文件存在
	if moist.IsExist(cfgFile) {
		// 读取配置
		var cfg Config
		err := yml.YamlFileToStruct(cfgFile, &cfg)

		// 配置文件读取成功
		if err == nil {
			// 地址设给全局变量
			config = &cfg

			// 初始化日志配置
			initLogConfig(cfg.LogFile, cfg.LogLevel)

			// 输出日志
			logger.WithFields(logrus.Fields{"配置": config}).Info()
		}

		// 配置文件读取失败
		if err != nil {
			// 默认配置
			cfg := defaultConfig()
			// 地址设给全局变量
			config = &cfg

			// 初始化日志配置(默认)
			initLogConfig(cfg.LogFile, cfg.LogLevel)

			// 输出日志
			logger.WithFields(logrus.Fields{
				"端口":       cfg.Port,
				"日志文件":     cfg.LogFile,
				"日志等级":     cfg.LogLevel,
				"目标主机文件":   cfg.LogFile,
				"模拟服务信息文件": cfg.LogLevel,
			}).Warn("读取配置文件发生错误，使用默认配置")

			// Go的init()方法里向Flutter端发送消息，Flutter端接收不到。
			// // 通知Flutter
			// msg := fmt.Sprintf("读取配置文件[%v]发生错误，使用默认端口[%v]", cfgFile, defaultPort)
			// go Notify(msg)
		}
	}
}

/// 初始化日志配置
func initLogConfig(file string, level int) {

	file = fmt.Sprintf("%v/%v", moist.CurrentDir(), file)

	// JSON格式
	// logger.SetFormatter(&logrus.JSONFormatter{})

	// 日志输出到文件
	logFile, err := os.OpenFile(file, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	// 日志同时输出到控制台和文件
	if err == nil {
		mw := io.MultiWriter(logFile, os.Stdout)
		logger.SetOutput(mw)
	} else {
		logger.SetOutput(os.Stdout)
		log.Println("只输出到控制台")
	}

	// 日志级别
	switch level {
	case -1:
		logger.SetLevel(logrus.DebugLevel)
	case 0:
		logger.SetLevel(logrus.InfoLevel)
	case 1:
		logger.SetLevel(logrus.WarnLevel)
	case 2:
		logger.SetLevel(logrus.ErrorLevel)
	case 3:
		logger.SetLevel(logrus.FatalLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}
}

/// 默认配置
func defaultConfig() Config {
	var cfg Config

	cfg.Port = defaultPort
	cfg.LogFile = defaultLogFile
	cfg.LogLevel = defaultLogLevel
	cfg.HostFile = defaultHostFile
	cfg.InfoFile = defaultInfoFile

	return cfg
}

/// 加载配置和输入
func Load() error {

	// 读取目标主机
	LoadHost(config)
	// 读取模拟服务信息
	LoadMockServiceInfo(config)
	// 读取模拟服务通用响应头
	LoadCommonResponseHeader(config)
	// 读取响应头信息
	LoadResponseHeader(config)

	return nil
}

/// 读取目标主机
func LoadHost(config *Config) {
	// 目标主机文件
	// hFile := config.HostFile
	hFile := fmt.Sprintf("%v/%v", moist.CurrentDir(), config.HostFile)

	// 目标主机文件不存在
	if !moist.IsExist(hFile) {
		msg := "目标主机文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: hFile}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, hFile)
		go Notify(msg)
	}

	// 目标主机文件存在
	if moist.IsExist(hFile) {
		// 读取目标主机
		var hosts []string
		err := yml.YamlFileToStruct(hFile, &hosts)

		if err == nil {
			hostSlice = hosts
			// 调试模式输出
			logger.WithFields(logrus.Fields{"目标主机": hostSlice}).Debug()
		}

		if err != nil {
			msg := "读取目标主机文件发生错误"
			logger.WithFields(logrus.Fields{logFieldFile: hFile, logFieldError: err}).Warn(msg)

			// 通知Flutter
			msg = fmt.Sprintf("%v[%v]", msg, hFile)
			go Notify(msg)
		}
	}

	// 默认目标主机
	if len(hostSlice) > 0 {
		defaultTargetHost = hostSlice[0]
		// 调试模式输出
		logger.WithFields(logrus.Fields{"默认目标主机": defaultTargetHost}).Debug()
	}
}

/// 读取模拟服务信息
func LoadMockServiceInfo(config *Config) {

	// 模拟服务信息文件
	// infoFile := config.InfoFile
	infoFile := fmt.Sprintf("%v/%v", moist.CurrentDir(), config.InfoFile)

	// 模拟服务信息文件不存在
	if !moist.IsExist(infoFile) {
		msg := "模拟服务信息文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: infoFile}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, infoFile)
		go Notify(msg)
	}

	// 模拟服务信息文件存在
	if moist.IsExist(infoFile) {
		// 读取模拟服务信息
		var info []MockServiceInfo
		err := yml.YamlFileToStruct(infoFile, &info)

		if err == nil {
			mockServiceInfoSlice = info
			// 调试模式输出
			logger.WithFields(logrus.Fields{"模拟服务信息": mockServiceInfoSlice}).Debug()
		}

		if err != nil {
			msg := "读取模拟服务信息文件发生错误"
			logger.WithFields(logrus.Fields{logFieldFile: infoFile, logFieldError: err}).Warn(msg)

			// 通知Flutter
			msg = fmt.Sprintf("%v[%v]", msg, infoFile)
			go Notify(msg)
		}
	}

}

/// 读取模拟服务通用响应头
func LoadCommonResponseHeader(config *Config) {
	// 模拟服务通用响应头文件
	fCommonResponesHeader := filePathCommonResponseHeader()

	// 模拟服务通用响应头文件不存在
	if !moist.IsExist(fCommonResponesHeader) {
		msg := "模拟服务通用响应头文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: fCommonResponesHeader}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, fCommonResponesHeader)
		go Notify(msg)
	}

	// 模拟服务通用响应头文件存在
	if moist.IsExist(fCommonResponesHeader) {
		// 读取模拟服务通用响应头
		var cmnRespHdr http.Header
		err := moist.JsonFileToStruct(fCommonResponesHeader, &cmnRespHdr)

		if err == nil {
			commonResponseHeader = cmnRespHdr
			// 调试模式输出
			logger.WithFields(logrus.Fields{"模拟服务通用响应头": commonResponseHeader}).Debug()
		}

		if err != nil {
			msg := "模拟服务通用响应头文件发生错误"
			logger.WithFields(logrus.Fields{logFieldFile: fCommonResponesHeader, logFieldError: err}).Warn(msg)

			// 通知Flutter
			msg = fmt.Sprintf("%v[%v]", msg, fCommonResponesHeader)
			go Notify(msg)
		}
	}
}

/// 读取响应头信息
/// 使用代理转发请求时，会自动生成URI对应的响应头信息
func LoadResponseHeader(config *Config) {
	// 响应头信息文件
	fResponseHeader := filePathResponseHeader()

	// 响应头信息文件不存在
	if !moist.IsExist(fResponseHeader) {
		msg := "响应头信息文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: fResponseHeader}).Info(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, fResponseHeader)
		go Notify(msg)

		// 空的响应头信息
		mapResponseHeader = make(map[string]http.Header)
	}

	// 响应头信息文件存在
	if moist.IsExist(fResponseHeader) {
		mHeader := make(map[string]http.Header)

		err := yml.YamlFileToStruct(fResponseHeader, &mHeader)
		if err != nil {
			msg := "读取响应头信息文件发生错误"
			logger.WithFields(logrus.Fields{
				logFieldFile:  fResponseHeader,
				logFieldError: err,
			}).Warn(msg)
		}
		mapResponseHeader = mHeader

		// 调试模式输出
		logger.WithFields(logrus.Fields{"响应头信息": mapResponseHeader}).Debug()
	}
}

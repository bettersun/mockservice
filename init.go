package mockservice

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/bettersun/moist"
	yml "github.com/bettersun/moist/yaml"
	"github.com/sirupsen/logrus"
)

// 全局变量：配置
var config *Config

// 全局变量：目标主机切片
var hostSlice []string

// 全局变量： 模拟服务信息切片
var mockServiceInfoSlice []MockServiceInfo

// 全局变量： URL对应模拟服务信息Map
var mapURIMockServiceInfo map[string]([]MockServiceInfo)

// 全局变量：HTTP服务
var server http.Server

// 全局变量：默认目标主机
var defaultTargetHost string

// 全局变量：默认目标主机
var targetHostSlice []string

// 全局变量：响应头
var mapResponseHeader map[string]http.Header

// 全局变量：运行中标志
var isRunning bool

// 全局变量：已注册URI
var registedURI []string

// 全局变量：日志记录器
var logger = logrus.New()

/// 初始化
func init() {

	// 日志设置
	// JSON格式
	// logger.SetFormatter(&logrus.JSONFormatter{})

	// 日志输出到文件
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		logger.SetOutput(file)
	} else {
		logger.Info("Failed to log to file, using default stderr")
	}

	// 日志级别：INFO
	logger.SetLevel(logrus.InfoLevel)
}

/// 加载配置和输入
func Load() error {

	// 配置文件
	cfgFile := ConfigFile

	// 模拟服务信息文件
	infoFile := InfoFile

	// 目标主机文件
	hFile := HostFile

	// 配置文件不存在
	if !moist.IsExist(cfgFile) {
		logger.WithFields(logrus.Fields{logFieldFile: cfgFile}).Warn("配置文件不存在")
		log.Println(fmt.Sprintf("配置文件不存在[%v]", cfgFile))

		logger.WithFields(logrus.Fields{"默认端口": defaultPort}).Warn("使用默认端口")
		log.Println(fmt.Sprintf("使用默认端口[%v]", defaultPort))

		var cfg Config
		cfg.Port = defaultPort
		config = &cfg

		// 通知Flutter
		msg := fmt.Sprintf("配置文件[%v]不存在，使用默认端口[%v]", cfgFile, defaultPort)
		Notify(msg)
	}

	// 配置文件存在
	if moist.IsExist(cfgFile) {
		// 读取配置
		var cfg Config
		result, err := yml.YamlFileToStruct(cfgFile, &cfg)

		if err == nil {
			config = result.(*Config)
			logger.WithFields(logrus.Fields{"配置": config}).Info()
		}

		if err != nil {
			logger.WithFields(logrus.Fields{logFieldFile: cfgFile, logFieldError: err}).Warn("读取配置文件发生错误")
			log.Println(fmt.Sprintf("读取配置文件发生错误[%v] %v", cfgFile, err))

			logger.WithFields(logrus.Fields{"默认端口": defaultPort}).Warn("使用默认端口")
			log.Println(fmt.Sprintf("使用默认端口[%v]", defaultPort))

			var cfg Config
			cfg.Port = defaultPort
			config = &cfg

			// 通知Flutter
			msg := fmt.Sprintf("读取配置文件[%v]发生错误，使用默认端口[%v]", cfgFile, defaultPort)
			Notify(msg)
		}
	}

	//=============================================================

	// 目标主机文件不存在
	if !moist.IsExist(hFile) {
		msg := "目标主机文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: hFile}).Warn(msg)
		log.Println(fmt.Sprintf("目标主机文件不存在[%v]", hFile))

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, hFile)
		Notify(msg)
	}

	// 目标主机文件存在
	if moist.IsExist(hFile) {
		// 读取目标主机
		var hosts []string
		result, err := yml.YamlFileToStruct(hFile, &hosts)

		if err == nil {
			hostSlice = *result.(*[]string)
			// 调试模式输出
			logger.WithFields(logrus.Fields{"目标主机": hostSlice}).Debug()
		}

		if err != nil {
			msg := "读取目标主机文件发生错误"
			logger.WithFields(logrus.Fields{logFieldFile: hFile, logFieldError: err}).Warn(msg)
			log.Println(fmt.Sprintf("读取目标主机文件发生错误[%v] %v", hFile, err))

			// 通知Flutter
			msg = fmt.Sprintf("%v[%v]", msg, hFile)
			Notify(msg)
		}
	}

	// 默认目标主机
	if len(hostSlice) > 0 {
		defaultTargetHost = hostSlice[0]
		// 调试模式输出
		logger.WithFields(logrus.Fields{"默认目标主机": defaultTargetHost}).Debug()
	}

	//=============================================================

	// 模拟服务信息文件不存在
	if !moist.IsExist(infoFile) {
		msg := "模拟服务信息文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: infoFile}).Warn(msg)
		log.Println(fmt.Sprintf("模拟服务信息文件不存在[%v]", infoFile))

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, infoFile)
		Notify(msg)
	}

	// 模拟服务信息文件存在
	if moist.IsExist(infoFile) {

		// 读取模拟服务信息
		var info []MockServiceInfo
		result, err := yml.YamlFileToStruct(infoFile, &info)

		if err == nil {
			mockServiceInfoSlice = *(result.(*[]MockServiceInfo))
			// 调试模式输出
			logger.WithFields(logrus.Fields{"模拟服务信息": mockServiceInfoSlice}).Debug()
			// log.Println(mockServiceInfoSlice)
		}

		if err != nil {
			msg := "读取模拟服务信息文件发生错误"
			logger.WithFields(logrus.Fields{logFieldFile: infoFile, logFieldError: err}).Warn("读取模拟服务信息文件发生错误")
			log.Println(fmt.Sprintf("读取模拟服务信息文件发生错误[%v] %v", infoFile, err))

			// 通知Flutter
			msg = fmt.Sprintf("%v[%v]", msg, infoFile)
			Notify(msg)
		}
	}

	// 基于模拟服务信息切片生成URL对应模拟服务信息Map
	// 相同的URI不同的HTTP请求方法
	mURIInfo := make(map[string]([]MockServiceInfo))
	for _, v := range mockServiceInfoSlice {
		vURIInfo, ok := mURIInfo[v.URI]
		if ok {
			vURIInfo = append(vURIInfo, v)
			mURIInfo[v.URI] = vURIInfo
		} else {
			tmp := []MockServiceInfo{}
			tmp = append(tmp, v)
			mURIInfo[v.URI] = tmp
		}
	}
	mapURIMockServiceInfo = mURIInfo
	// 调试模式输出
	logger.WithFields(logrus.Fields{"URL对应模拟服务信息": mapURIMockServiceInfo}).Debug()
	// log.Println(mapUrlMockServiceInfoSlice)

	//=============================================================

	// 读取响应头信息
	// 使用代理转发请求时，会自动生成URI对应的响应头信息

	// 响应头信息文件
	responseHeaderFile := fmt.Sprintf("%v%v/%v", moist.CurrentDir(), responsePath, ResponseHeaderFile)

	// 响应头信息文件不存在
	if !moist.IsExist(responseHeaderFile) {
		msg := "响应头信息文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: responseHeaderFile}).Warn(msg)
		log.Println(fmt.Sprintf("响应头信息文件不存在[%v]", responseHeaderFile))

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, responseHeaderFile)
		Notify(msg)

		mapResponseHeader = make(map[string]http.Header)
	}

	// 响应头信息文件存在
	if moist.IsExist(responseHeaderFile) {
		mHeader := make(map[string]http.Header)
		result, err := yml.YamlFileToStruct(responseHeaderFile, &mHeader)
		if err != nil {
			msg := "读取响应头信息文件发生错误"
			logger.WithFields(logrus.Fields{
				"响应头信息文件":     responseHeaderFile,
				logFieldError: err,
			}).Warn(msg)

			log.Println(fmt.Sprintf("%v[%v] %v", msg, responseHeaderFile, err))
		}
		mapResponseHeader = *result.(*map[string]http.Header)

		// 调试模式输出
		logger.WithFields(logrus.Fields{"已保存的响应头信息": mapResponseHeader}).Debug()
	}

	return nil
}

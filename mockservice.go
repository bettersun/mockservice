package mockservice

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/bettersun/moist"
	"github.com/sirupsen/logrus"
)

/// 全局变量：配置
var config *Config

/// 全局变量：目标主机切片
var hostSlice []string

/// 全局变量： 模拟服务信息切片
var mockServiceInfoSlice []MockServiceInfo

/// 全局变量：HTTP服务
var server http.Server

/// 全局变量：默认目标主机
var defaultTargetHost string

/// 全局变量：响应头Map
var mapResponseHeader map[string]http.Header

/// 全局变量：模拟服务通用响应头
var commonResponseHeader http.Header

/// 模拟服务(命令行调用)
func MockServiceCommand() {
	MockService()
	log.Println(moist.CurrentDir())

	// 命令行执行时，需要下面这段代码
	ch := make(chan string)
	<-ch
}

/// 模拟服务
func MockService() {
	// 监听端口
	port := ":" + config.Port
	server = http.Server{
		Addr: port,
	}

	// TODO 需要确认端口是否被占用
	go http.ListenAndServe(port, http.HandlerFunc(DoHandle))

	msg := fmt.Sprintf("服务运行中... 端口[%v]", config.Port)
	logger.Info(msg)
	go Notify(msg)
}

/// 关闭服务
func CloseServer() error {
	err := server.Close()
	if err != nil {
		msg := "关闭服务发生错误"
		logger.Warn(msg)
	}

	msg := "服务已关闭"
	logger.Info(msg)
	go Notify(msg)

	return err
}

/// 响应函数
func DoHandle(w http.ResponseWriter, r *http.Request) {

	url := r.URL.String()
	// URL(有参数时，问号之前的部分)
	bURL := baseURL(url)

	msg := "请求URL"
	logger.WithFields(logrus.Fields{logFieldURL: url, logFieldHTTPMethod: r.Method}).Info(msg)
	msg = fmt.Sprintf("%v[%v: %v]", msg, r.Method, url)
	go Notify(msg)

	// ======================================================================================

	isConfigedURL := false
	// 存在模拟服务信息
	for _, vInfo := range mockServiceInfoSlice {
		// 根据请求的URL和请求方法确定响应函数
		if vInfo.URL == bURL && vInfo.Method == r.Method {
			isConfigedURL = true

			// 使用模拟服务
			if vInfo.UseMockService {
				doMockService(w, r, vInfo, config)
				return
			}

			// 不使用模拟服务
			if !vInfo.UseMockService {
				// 目标主机
				var host string
				if vInfo.UseDefaultTargetHost {
					host = defaultTargetHost
				} else {
					host = vInfo.TargetHost
				}

				doProxyService(w, r, host)
				return
			}
		}
	}

	// 请求的URL和请求方法对应的模拟服务信息
	if !isConfigedURL {
		var info MockServiceInfo
		info.URL = bURL                     // URL(有参数时，问号之前的部分)
		info.Method = r.Method              // 请求方法
		info.TargetHost = defaultTargetHost // 目标主机: 默认目标主机
		info.UseDefaultTargetHost = true    // 使用默认目标主机
		info.UseMockService = false         // 不使用虚拟服务
		info.StatusCode = http.StatusOK     // 默认返回200
		info.ResponseFile = ""              // 响应文件默认为空
		info.Description = bURL             // 说明默认使用URL(有参数时，问号之前的部分)

		// 添加到模拟服务信息切片
		mockServiceInfoSlice = append(mockServiceInfoSlice, info)

		msg = "模拟服务信息未配置，使用默认目标主机"
		logger.Info(msg)
		go Notify(msg)

		// 响应
		doProxyService(w, r, defaultTargetHost)

		// 保存响应信息
		SaveMockServiceInfo()

		// 向Flutter发送消息触发Flutter事件
		go NotifyAddMockServiceInfo(info)
	}
}

/// 转发请求
// 参考：https://www.cnblogs.com/boxker/p/11046342.html
func doProxyService(w http.ResponseWriter, r *http.Request, host string) {
	go Notify("访问目标主机")

	url := r.URL.String()
	// 创建一个HttpClient用于转发请求
	cli := &http.Client{}

	// 读取请求的Body
	// 读取后 r.Body 即关闭，无法再次读取
	// 若需要再次读取，需要用读取到的内容再次构建Reader
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "读取请求体发生错误"
		logger.WithFields(logrus.Fields{logFieldURL: url, logFieldHTTPMethod: r.Method}).Info(msg)
	}

	// 日志
	logger.WithFields(logrus.Fields{logFieldHost: host}).Info("访问目标主机")

	// 转发的URL
	reqURL := host + url

	// 输出请求头
	OutRequest(url, r.Method, r.Header, string(body))

	// 创建转发用的请求
	reqProxy, err := http.NewRequest(r.Method, reqURL, strings.NewReader(string(body)))
	if err != nil {
		msg := "创建转发请求发生错误"
		logger.WithFields(logrus.Fields{
			logFieldHost:       host,
			logFieldURL:        url,
			logFieldHTTPMethod: r.Method,
			logFieldError:      err,
		}).Error(msg)
		go Notify(msg)
		return
	}

	// 转发请求的表头
	for k, v := range r.Header {
		reqProxy.Header.Set(k, v[0])
	}

	// 发起请求
	responseProxy, err := cli.Do(reqProxy)
	if err != nil {
		msg := "转发请求发生错误"
		logger.WithFields(logrus.Fields{
			logFieldHost:       host,
			logFieldURL:        url,
			logFieldHTTPMethod: r.Method,
			logFieldError:      err,
		}).Error(msg)
		go Notify(msg)

		// StatusCode
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	defer responseProxy.Body.Close()

	// 响应头
	for k, v := range responseProxy.Header {
		w.Header().Set(k, v[0])
	}

	// 响应头键
	k := keyResponseHeader(url, r.Method)
	_, ok := mapResponseHeader[k]
	// 响应头Map中不存在，则添加到响应头Map中
	if !ok {
		mapResponseHeader[k] = responseProxy.Header
	}
	// 输出响应头(每次都输出用于保存最新的响应头)
	OutResponseHeader(mapResponseHeader)

	// 响应为JSON判断
	isResponseJSON := isResponseJSON(responseProxy.Header)
	// gzip压缩判断
	isGzipped := isGzipped(responseProxy.Header)

	// 响应体数据
	var data []byte
	// 读取响应体
	data, err = ioutil.ReadAll(responseProxy.Body)
	if err != nil {
		msg := "读取响应体发生错误"
		logger.WithFields(logrus.Fields{logFieldError: err}).Info(msg)
		log.Println(fmt.Sprintf("%v %v", msg, err))
	}

	// 输出数据
	var dataOutput []byte
	// gzip压缩编码数据
	if isGzipped {
		resProxyGzippedBody := ioutil.NopCloser(bytes.NewBuffer(data))
		defer resProxyGzippedBody.Close() // 延时关闭

		// gzip Reader
		gr, err := gzip.NewReader(resProxyGzippedBody)
		if err != nil {
			msg := "创建gzip读取器发生错误"
			logger.WithFields(logrus.Fields{logFieldError: err}).Info(msg)
			log.Println(fmt.Sprintf("%v %v", msg, err))
		}
		defer gr.Close()

		// 读取gzip对象内容
		dataOutput, err = ioutil.ReadAll(gr)
		if err != nil {
			msg := "读取gzip对象内容发生错误"
			logger.WithFields(logrus.Fields{logFieldError: err}).Info(msg)
			log.Println(fmt.Sprintf("%v %v", msg, err))
		}
	} else { // 非gzip压缩
		dataOutput = data
	}

	// 输出响应体到文件
	OutResponseBody(r.Method, url, isResponseJSON, dataOutput)

	// response的Body不能多次读取，需要重新生成可读取的Body
	resProxyBody := ioutil.NopCloser(bytes.NewBuffer(data))
	defer resProxyBody.Close() // 延时关闭

	// StatusCode
	w.WriteHeader(responseProxy.StatusCode)
	// 复制转发的响应Body到响应Body
	io.Copy(w, resProxyBody)
}

/// 模拟服务
func doMockService(w http.ResponseWriter, r *http.Request, info MockServiceInfo, config *Config) {
	msg := "使用模拟服务"
	logger.Info(msg)
	go Notify(msg)

	url := r.URL.String()

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := "读取请求体发生错误"
		logger.WithFields(logrus.Fields{logFieldURL: url, logFieldHTTPMethod: r.Method}).Info(msg)
		log.Println(fmt.Sprintf("%v %v", msg, err))
	}
	// 输出请求信息
	OutRequest(url, r.Method, r.Header, string(body))

	// =====================================================================
	var header http.Header

	// 响应头: URL_Method
	k := keyResponseHeader(url, r.Method)
	// 获取URL对应的响应头
	header, ok := mapResponseHeader[k]

	// 无法获取URL对应的响应头，且不使用模拟服务通用响应头时
	if !ok && !config.UseMockCommonResponseHeader {
		msg := "无法获取URL对应的响应头信息,可先访问目标主机以保存相关信息或使用模拟服务通用响应头"
		logger.WithFields(logrus.Fields{logFieldURL: url, logFieldHTTPMethod: r.Method}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, k)
		go Notify(msg)

		m := make(map[string]interface{}, 0)
		m["message"] = msg

		msgStream, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(msgStream)
		return
	}
	// 无法获取URL对应的响应头，且使用模拟服务通用响应头时
	if !ok && config.UseMockCommonResponseHeader {
		// 模拟服务通用响应头
		header = commonResponseHeader

		msg := "URL的响应头信息不存在，使用模拟服务通用响应头"
		logger.Warn(msg)
		go Notify(msg)
	}
	// 能够获取URL对应的响应头
	if ok {
		// 仅输出日志
		msg := "使用已保存的响应头信息"
		logger.Info(msg)
		go Notify(msg)
	}

	// 调试模式输出
	logger.WithFields(logrus.Fields{"响应头信息": header}).Debug()

	// 响应是否为JSON
	isResponseJSON := false
	// 响应是否gzip压缩
	isGzipped := false
	for k, v := range header {
		// Content-Length 不添加到响应头
		// http.ResponseWriter会自动计算
		if k == httpHeaderContentLength {
			continue
		}

		w.Header().Set(k, v[0])

		// 响应是否为JSON
		if k == httpHeaderContentType && moist.Contains(v, contentTypeJSON) {
			isResponseJSON = true
		}

		// 响应是否Gzip压缩
		if k == httpHeaderContentEncoding && moist.Contains(v, encodingGzip) {
			isGzipped = true
		}
	}

	// ==============================================================================================

	// 响应文件
	responseFile := info.ResponseFile

	var stream []byte
	// 响应体文件未指定
	if responseFile == "" {
		msg := "模拟服务响应文件未指定，响应体返回空内容。"
		logger.WithFields(logrus.Fields{logFieldFile: responseFile}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, responseFile)
		go Notify(msg)
	}
	// 响应体文件指定时进行处理
	if responseFile != "" {
		// 日志：模拟服务响应文件
		logger.WithFields(logrus.Fields{logFieldMockResponseFile: responseFile}).Info()
		msg = fmt.Sprintf("%v[%v]", logFieldMockResponseFile, responseFile)
		go Notify(msg)

		// 完整路径
		fResponse := fmt.Sprintf("%v/%v", moist.CurrentDir(), responseFile)

		// 响应文件不存在
		if !moist.IsExist(fResponse) {
			msg = "模拟服务响应文件不存在"
			logger.WithFields(logrus.Fields{logFieldFile: fResponse}).Warn(msg)

			// 通知Flutter
			msg = fmt.Sprintf("%v[%v]", msg, fResponse)
			go Notify(msg)

			m := make(map[string]interface{}, 0)
			m["message"] = msg

			msgStream, err := json.Marshal(m)
			if err != nil {
				log.Println(err)
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write(msgStream)
			return
		}

		// 响应文件存在
		// 响应体为JSON
		if isResponseJSON {
			// 响应文件转换成Map
			data, err := moist.JsonFileToMap(fResponse)
			if err != nil {
				log.Println(err)
			}

			// 转换成字节
			stream, err = json.Marshal(data)
			if err != nil {
				log.Println(err)
			}
		} else {
			// 响应体非JSON
			data, err := ioutil.ReadFile(fResponse)
			if err != nil {
				log.Println(err)
			}

			stream = data
		}
	}

	// 响应状态码，必须放在w.Header().Set(k, v)之后
	statusCode := info.StatusCode
	if statusCode == 0 {
		// 响应状态码为0时，设为200
		statusCode = http.StatusOK
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(info.StatusCode)
	}
	logger.WithFields(logrus.Fields{logFieldMockStatusCode: statusCode}).Warn()
	// 通知Flutter
	msg = fmt.Sprintf("%v[%v]", logFieldMockStatusCode, statusCode)
	go Notify(msg)

	// 调试模式输出
	logger.WithFields(logrus.Fields{
		logFieldURL:              url,
		logFieldHTTPMethod:       r.Method,
		logFieldResponseHeader:   w.Header(),
		logFieldMockStatusCode:   statusCode,
		logFieldMockResponseFile: responseFile,
		logFieldMockResponseBody: string(stream),
	}).Debug()

	// 响应
	if isGzipped {
		// gzip压缩
		buffer := new(bytes.Buffer)
		gw := gzip.NewWriter(buffer)
		gw.Write(stream)
		gw.Flush()

		w.Write(buffer.Bytes())
	} else {
		w.Write(stream)
	}
}

/// gzip压缩判断
func isGzipped(header http.Header) bool {

	if header == nil {
		return false
	}

	contentEncoding := header.Get(httpHeaderContentEncoding)
	isGzipped := false
	if strings.Contains(contentEncoding, encodingGzip) {
		isGzipped = true
	}

	return isGzipped
}

/// 判断响应类型是否是JSON
func isResponseJSON(header http.Header) bool {

	if header == nil {
		return false
	}

	v := header.Get(httpHeaderContentType)
	result := false
	if strings.Contains(v, contentTypeJSON) {
		result = true
	}

	return result
}

/// 响应头的键值
func keyResponseHeader(url string, method string) string {
	// URL(有参数时，问号之前的部分)
	bURL := baseURL(url)
	return fmt.Sprintf("%v_%v", bURL, method)
}

/// 获取目标主机列表
func ListHost() []string {
	return hostSlice
}

/// 设置默认目标主机
func SetDefaultTargetHost(targetHost string) bool {
	defaultTargetHost = targetHost

	logger.WithFields(logrus.Fields{"默认目标主机": targetHost}).Info("默认目标主机已改变")
	return true
}

/// 获取模拟服务信息
func ListMockServiceInfo() []MockServiceInfo {
	return mockServiceInfoSlice
}

/// 更新单个模拟服务信息和URL对应模拟服务信息Map
func UpdateMockServiceInfo(info MockServiceInfo) {
	// 更新内存中的模拟服务信息
	for i, v := range mockServiceInfoSlice {
		if v.URL == info.URL && v.Method == info.Method {
			mockServiceInfoSlice[i] = info
		}
	}

	logger.WithFields(logrus.Fields{"模拟服务信息": info}).Info("模拟服务信息已更新")
}

/// 更新单个模拟服务信息和URL对应模拟服务信息Map
func UpdateAllMockServiceInfo(newInfoSlice []MockServiceInfo) {
	// 更新内存中的模拟服务信息
	for i, v := range mockServiceInfoSlice {
		for _, v2 := range newInfoSlice {
			if v.URL == v2.URL && v.Method == v2.Method {
				mockServiceInfoSlice[i] = v2
				break
			}
		}
	}

	logger.WithFields(logrus.Fields{"全部的模拟服务信息": newInfoSlice}).Info("全部的模拟服务信息已更新")
}

/// 保存模拟服务信息
func SaveMockServiceInfo() error {
	// 保存模拟服务信息
	err := OutputMockServiceInfo(*config, mockServiceInfoSlice)
	if err != nil {
		log.Println(err)
		logger.WithFields(logrus.Fields{
			logFieldFile: config.InfoFile,
		}).Warn("模拟服务信息保存失败")
	}
	logger.WithFields(logrus.Fields{
		logFieldFile: config.InfoFile,
	}).Info("模拟服务信息已保存")

	return err
}

/// 保存主机列表
func SaveHost() error {
	err := OutputHost(*config, hostSlice)
	if err != nil {
		log.Println(err)
		logger.WithFields(logrus.Fields{
			logFieldFile: config.InfoFile,
		}).Warn("主机列表保存失败")
	}
	logger.WithFields(logrus.Fields{
		logFieldFile: config.InfoFile,
	}).Info("主机列表已保存")

	return err
}

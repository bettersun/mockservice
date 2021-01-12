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

// 模拟服务(命令行调用)
func MockServiceCommand() {

	MockService()

	log.Println(moist.CurrentDir())

	// 命令行执行时，需要下面这段代码
	// 从UI执行时，需要注释掉下面这段代码
	ch := make(chan string)
	<-ch
}

/// 模拟服务
func MockService() error {

	if isRunning {
		log.Println("The Server is running")
	}

	// 监听
	port := ":" + config.Port
	server = http.Server{
		Addr: port,
	}

	// 注册
	for url := range mapURIMockServiceInfo {
		// 对未注册的URL进行注册
		if !moist.IsInSlice(registedURI, url) {
			http.HandleFunc(url, DoHandle)
			registedURI = append(registedURI, url)
		}
	}

	// TODO 需要确认端口是否被占用
	go server.ListenAndServe()
	isRunning = true
	log.Printf("Listen and serve [%v]", config.Port)

	logger.Info("服务运行中")
	Notify("服务运行中")

	return nil
}

/// 关闭服务
func CloseServer() error {
	isRunning = false
	err := server.Close()
	if err == nil {
		log.Println("Server closed")
	}

	logger.Info("服务已关闭")
	Notify("服务已关闭")
	return err
}

/// 请求处理函数
func DoHandle(w http.ResponseWriter, r *http.Request) {

	uri := r.URL.String()
	logger.WithFields(logrus.Fields{logFieldURI: uri, logFieldHTTPMethod: r.Method}).Info()

	message := fmt.Sprintf("URL[%v: %v]", r.Method, uri)
	Notify(message)

	// 输出请求到文件
	OutRequest(r)

	// URL对应的模拟服务信息
	for url, v := range mapURIMockServiceInfo {
		// 模拟服务信息切片
		if url == r.URL.String() {
			for _, vInfo := range v {
				// 根据请求的URL和请求方法确定响应函数
				if vInfo.URI == r.URL.String() && vInfo.Method == r.Method {
					if vInfo.UseMockService {
						doMockService(w, r, vInfo)
						return
					} else {
						doProxyService(w, r, &vInfo)
						return
					}
				}
			}
		}
	}

	msg := "目标主机无法访问或模拟服务信息未设置"
	logger.WithFields(logrus.Fields{logFieldURI: uri, logFieldHTTPMethod: r.Method}).Info(msg)
	Notify(msg)

	log.Println(msg)
}

/// 转发请求
// 参考：https://www.cnblogs.com/boxker/p/11046342.html
func doProxyService(w http.ResponseWriter, r *http.Request, info *MockServiceInfo) {

	// 通知Flutter
	Notify("访问目标主机")

	uri := r.URL.String()
	// 创建一个HttpClient用于转发请求
	cli := &http.Client{}

	// 读取请求的Body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Print("ioutil.ReadAll(r.Body) ", err.Error())
	}

	var host string
	if info.UseDefaultTargetHost {
		host = defaultTargetHost
	} else {
		host = info.TargetHost
	}

	// 转发的URL
	reqURL := host + uri

	logger.WithFields(logrus.Fields{
		logFieldHost: host, logFieldURI: uri,
		logFieldHTTPMethod:    r.Method,
		logFieldRequestHeader: r.Header,
		// logFieldRequestBody:   string(body),	// 请求体不输出到日志，可参考保存的请求体文件
	},
	).Info("访问目标主机")

	// 创建转发用的请求
	reqProxy, err := http.NewRequest(r.Method, reqURL, strings.NewReader(string(body)))
	if err != nil {
		log.Print("http.NewRequest ", err.Error())
		return
	}

	// 转发请求的表头
	for k, v := range r.Header {
		reqProxy.Header.Set(k, v[0])
	}

	// 发起请求
	responseProxy, err := cli.Do(reqProxy)
	if err != nil {
		fmt.Print("cli.Do() ", err.Error())
		return
	}
	defer responseProxy.Body.Close()

	// 响应头
	for k, v := range responseProxy.Header {
		w.Header().Set(k, v[0])
	}
	// log.Println(responseProxy.Header)

	// 输出响应头
	k := fmt.Sprintf("%v_%v", uri, r.Method)
	_, ok := mapResponseHeader[k]
	if !ok {
		mapResponseHeader[k] = responseProxy.Header
	}
	OutResponseHeader(ResponseHeaderFile, mapResponseHeader)

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
	OutResponseBody(r.Method, uri, dataOutput)
	// log.Println(string(dataOutput))

	// response的Body不能多次读取，需要重新生成可读取的Body
	resProxyBody := ioutil.NopCloser(bytes.NewBuffer(data))
	defer resProxyBody.Close() // 延时关闭

	// StatusCode
	w.WriteHeader(responseProxy.StatusCode)
	// 复制转发的响应Body到响应Body
	io.Copy(w, resProxyBody)
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

/// 模拟服务
func doMockService(w http.ResponseWriter, r *http.Request, info MockServiceInfo) {

	// 通知Flutter
	Notify("使用模拟服务")

	uri := r.URL.String()

	// 响应头: URI_Method
	k := fmt.Sprintf("%v_%v", uri, r.Method)
	// 获取URI对应的响应头
	header, ok := mapResponseHeader[k]
	// 无法获取URI对应的响应头时
	if !ok {
		msg := "无法获取URI对应的响应头信息,可先访问目标主机以保存相关信息"
		logger.WithFields(logrus.Fields{
			logFieldURI:        uri,
			logFieldHTTPMethod: r.Method,
		}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, k)
		Notify(msg)
	}

	// 调试模式输出
	logger.WithFields(logrus.Fields{
		logFieldURI:        uri,
		logFieldHTTPMethod: r.Method,
		"响应头信息":            header,
	}).Debug()

	// 响应是否为JSON
	isResponseJSON := false
	// 响应是否gzip压缩
	isGzipped := false
	if ok {
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
	}

	// 调试模式输出
	logger.WithFields(logrus.Fields{
		logFieldURI:        uri,
		logFieldHTTPMethod: r.Method,
		"响应体JSON":          isResponseJSON,
	}).Debug()

	// ==============================================================================================

	// 响应文件
	responseFile := moist.CurrentDir() + info.ResponseFile

	// 响应文件不存在
	if !moist.IsExist(responseFile) {
		msg := "模拟服务响应文件不存在"
		logger.WithFields(logrus.Fields{logFieldFile: responseFile}).Warn(msg)

		// 通知Flutter
		msg = fmt.Sprintf("%v[%v]", msg, responseFile)
		Notify(msg)

		m := make(map[string]interface{}, 0)
		m["message"] = msg

		msgStream, err := json.Marshal(m)
		if err != nil {
			log.Println(err)
		}
		w.Write(msgStream)
		return
	}

	// // 通知Flutter
	// message := fmt.Sprintf("模拟服务响应文件[%v]", responseFile)
	// Notify(message)

	// 响应体
	var stream []byte
	if isResponseJSON {
		// 响应文件转换成Map
		data, err := moist.JsonFileToMap(responseFile)
		if err != nil {
			log.Println(err)
		}

		// 转换成字节
		stream, err = json.Marshal(data)
		if err != nil {
			log.Println(err)
		}

	} else {
		data, err := ioutil.ReadFile(responseFile)
		if err != nil {
			log.Println(err)
		}

		stream = data
	}
	// log.Println(string(stream))

	// 调试模式输出
	logger.WithFields(logrus.Fields{
		logFieldURI:        uri,
		logFieldHTTPMethod: r.Method,
		"模拟服务响应文件":         responseFile,
		"模拟服务响应体":          string(stream),
	}).Debug()

	// 响应状态码，必须放在w.Header().Set(k, v)之后
	if info.StatusCode == 0 {
		w.WriteHeader(200)
	} else {
		w.WriteHeader(info.StatusCode)
	}

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

/// 是否运行中
func IsRunning() bool {
	return isRunning
}

/// 获取目标主机列表
func ListHost() []string {
	return hostSlice
}

/// 获取模拟服务信息
func ListMockServiceInfo() []MockServiceInfo {
	return mockServiceInfoSlice
}

/// 更新单个模拟服务信息和URL对应模拟服务信息Map
func UpdateMockServiceInfo(info MockServiceInfo) {

	// 更新内存中的模拟服务信息
	for i, v := range mockServiceInfoSlice {
		if v.URI == info.URI && v.Method == info.Method {
			mockServiceInfoSlice[i] = info
		}
	}

	// 更新内存中的 URL对应模拟服务信息Map
	for url, v := range mapURIMockServiceInfo {
		if url == info.URI {
			var infoSlice []MockServiceInfo
			for _, vInfo := range v {
				// 只更新对应方法
				if info.Method == vInfo.Method {
					infoSlice = append(infoSlice, info)
				} else {
					infoSlice = append(infoSlice, vInfo)
				}
			}
			mapURIMockServiceInfo[url] = infoSlice
		}
	}
}

/// 更新单个模拟服务信息和URL对应模拟服务信息Map
func UpdateAllMockServiceInfo(newInfoSlice []MockServiceInfo) {

	// 更新内存中的模拟服务信息
	for i, v := range mockServiceInfoSlice {
		for _, v2 := range newInfoSlice {
			if v.URI == v2.URI && v.Method == v2.Method {
				mockServiceInfoSlice[i] = v2
				break
			}
		}
	}

	// 更新内存中的 URL对应模拟服务信息Map
	for url, v := range mapURIMockServiceInfo {

		// 新的模拟服务信息
		for _, vNew := range newInfoSlice {

			// URL对应模拟服务信息Map 的键 对应 新的模拟服务信息的URI
			if url == vNew.URI {
				var infoSlice []MockServiceInfo

				// URL对应模拟服务信息Map 的 值
				for _, vInfo := range v {
					if vNew.Method == vInfo.Method {
						infoSlice = append(infoSlice, vNew)
					}
				}

				mapURIMockServiceInfo[url] = infoSlice
				break
			}
		}
	}
}

// 保存模拟服务信息
func SaveMockServiceInfo() error {

	err := OutputMockServiceInfo(InfoFile, mockServiceInfoSlice)
	if err != nil {
		log.Println(err)
	}

	logger.WithFields(logrus.Fields{
		logFieldFile: InfoFile,
	}).Info("模拟服务信息已保存")

	return err
}

// 保存主机列表
func SaveHost() error {
	err := OutputHost(HostFile, hostSlice)
	return err
}

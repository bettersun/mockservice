package ch

import (
	"log"

	"github.com/go-flutter-desktop/go-flutter"
	"github.com/go-flutter-desktop/go-flutter/plugin"

	"github.com/bettersun/mockservice"
	"github.com/bettersun/moist"
)

// go-flutter插件需要声明包名和函数名
// Flutter代码中调用时需要指定相应的包名和函数名
const (
	channelName = "bettersun.go-flutter.plugin.mockservice"

	funcNameHelo      = "hello"
	funcNameRun       = "run"
	funcNameClose     = "close"
	funcNameLoad      = "load"
	funcNameIsRunning = "IsRunning"

	funcNameUpdateInfo    = "updateInfo"
	funcNameUpdateAllInfo = "updateAllInfo"
	funcNameSaveInfo      = "saveInfo"

	funcNameHostlist     = "hostlist"
	funcNameInfolist     = "infolist"
	funcNameResponselist = "responselist"

	// funcNameNotify = "notify"
)

/// 声明插件结构体
type MockServicePlugin struct{}

/// 指定为go-flutter插件
var _ flutter.Plugin = &MockServicePlugin{}

/// 初始化插件
func (MockServicePlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	ch := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})

	ch.HandleFunc(funcNameHelo, helloFunc)
	ch.HandleFunc(funcNameRun, runFunc)
	ch.HandleFunc(funcNameClose, closeFunc)
	ch.HandleFunc(funcNameLoad, loadFunc)
	ch.HandleFunc(funcNameIsRunning, isRunningFunc)

	ch.HandleFunc(funcNameUpdateInfo, updateInfoFunc)
	ch.HandleFunc(funcNameUpdateAllInfo, updateAllInfoFunc)
	ch.HandleFunc(funcNameSaveInfo, saveInfoFunc)

	ch.HandleFunc(funcNameHostlist, hostListFunc)
	ch.HandleFunc(funcNameInfolist, infoListFunc)
	ch.HandleFunc(funcNameResponselist, responseListFunc)

	// 用于向flutter端发送消息
	mockservice.Channel = ch

	return nil
}

/// Hello
func helloFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("helloFunc()")
	return "message from mockservice", nil
}

/// 启动服务
func runFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("runFunc()")
	mockservice.MockService()
	return "running", nil
}

/// 运行中
func isRunningFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("isRunningFunc()")
	isRunning := mockservice.IsRunning()
	return isRunning, nil
}

/// 关闭服务
func closeFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("closeFunc()")
	err = mockservice.CloseServer()
	if err != nil {
		return false, nil
	}
	return true, nil
}

/// 加载(配置及输入文件)
func loadFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("loadFunc()")
	err = mockservice.Load()
	if err != nil {
		return false, nil
	}
	return true, nil
}

/// 更新模拟服务信息
func updateInfoFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("updateInfoFunc()")

	// 检查参数
	m := arguments.(map[interface{}]interface{})

	useDefaultTargetHost, ok := m["useDefaultTargetHost"].(bool)
	if !ok {
		log.Println("ERROR: [useDefaultTargetHost]")
	}
	useMockService, ok := m["useMockService"].(bool)
	if !ok {
		log.Println("ERROR: [useMockService]")
	}
	targetHost, ok := m["targetHost"].(string)
	if !ok {
		log.Println("ERROR: [targetHost]")
	}
	uri, ok := m["uri"].(string)
	if !ok {
		log.Println("ERROR: [uri]")
	}
	method, ok := m["method"].(string)
	if !ok {
		log.Println("ERROR: [method]")
	}
	// Flutter端的int，在Go端需要使用int32来接收，然后再转换成Int
	statusCode, ok := m["statusCode"].(int32)
	if !ok {
		log.Println("ERROR: [statusCode]")
	}
	responseFile, ok := m["responseFile"].(string)
	if !ok {
		log.Println("ERROR: [responseFile]")
	}
	description, ok := m["description"].(string)
	if !ok {
		log.Println("ERROR: [description]")
	}

	var info mockservice.MockServiceInfo

	info.UseDefaultTargetHost = useDefaultTargetHost
	info.UseMockService = useMockService
	info.TargetHost = targetHost
	info.URI = uri
	info.Method = method
	info.StatusCode = int(statusCode)
	info.ResponseFile = responseFile
	info.Description = description

	mockservice.UpdateMockServiceInfo(info)
	return true, nil
}

/// 更新所有模拟服务信息
func updateAllInfoFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("updateAllInfoFunc()")

	// 检查参数
	m := arguments.([]interface{})

	var infoSlice []mockservice.MockServiceInfo

	for _, v := range m {
		// 检查参数
		vm := v.(map[interface{}]interface{})

		useDefaultTargetHost, ok := vm["useDefaultTargetHost"].(bool)
		if !ok {
			log.Println("ERROR: [useDefaultTargetHost]")
		}
		useMockService, ok := vm["useMockService"].(bool)
		if !ok {
			log.Println("ERROR: [useMockService]")
		}
		targetHost, ok := vm["targetHost"].(string)
		if !ok {
			log.Println("ERROR: [targetHost]")
		}
		uri, ok := vm["uri"].(string)
		if !ok {
			log.Println("ERROR: [uri]")
		}
		method, ok := vm["method"].(string)
		if !ok {
			log.Println("ERROR: [method]")
		}
		// Flutter端的int，在Go端需要使用int32来接收，然后再转换成Int
		statusCode, ok := vm["statusCode"].(int32)
		if !ok {
			log.Println("ERROR: [statusCode]")
		}
		responseFile, ok := vm["responseFile"].(string)
		if !ok {
			log.Println("ERROR: [responseFile]")
		}
		description, ok := vm["description"].(string)
		if !ok {
			log.Println("ERROR: [description]")
		}

		var info mockservice.MockServiceInfo

		info.UseDefaultTargetHost = useDefaultTargetHost
		info.UseMockService = useMockService
		info.TargetHost = targetHost
		info.URI = uri
		info.Method = method
		info.StatusCode = int(statusCode)
		info.ResponseFile = responseFile
		info.Description = description

		infoSlice = append(infoSlice, info)
	}

	mockservice.UpdateAllMockServiceInfo(infoSlice)
	return true, nil
}

/// 保存信息
func saveInfoFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("saveInfoFunc()")

	err = mockservice.SaveMockServiceInfo()
	if err != nil {
		return false, nil
	}

	return true, nil
}

/// 获取目标主机列表
func hostListFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("hostListFunc()")

	m := make(map[interface{}]interface{})

	var slice []interface{}
	list := mockservice.ListHost()
	for _, v := range list {
		slice = append(slice, v)
	}

	m["HostList"] = slice

	return m, nil
}

/// 获取模拟服务信息列表
func infoListFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("infoListFunc()")

	m := make(map[interface{}]interface{})

	var slice []interface{}
	list := mockservice.ListMockServiceInfo()
	for _, v := range list {
		entity, err := moist.StructToIfKeyMap(v)
		if err != nil {
			log.Println(err)
		}

		slice = append(slice, entity)
	}

	m["InfoList"] = slice

	return m, nil
}

/// 获取响应文件列表
func responseListFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("responseListFunc()")

	// 参数
	m := arguments.(map[interface{}]interface{})
	uri, ok := m["uri"].(string)
	if !ok {
		log.Println("ERROR: [uri]")
	}
	method, ok := m["method"].(string)
	if !ok {
		log.Println("ERROR: [method]")
	}

	var slice []interface{}
	list, err := mockservice.LoadResponseFile(uri, method)
	if err != nil {
		log.Println(err)
		return slice, err
	}

	for _, v := range list {
		slice = append(slice, v)
	}

	mResult := make(map[interface{}]interface{})
	mResult["ResponseList"] = slice
	return mResult, nil
}

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

	funcNameRun   = "run"
	funcNameClose = "close"
	funcNameLoad  = "load"

	funcNameUpdateInfo    = "updateInfo"
	funcNameUpdateAllInfo = "updateAllInfo"
	funcNameSaveInfo      = "saveInfo"

	funcNameHostlist             = "hostlist"
	funcNameSetDefaultTargetHost = "setDefaultTargetHost"
	funcNameInfolist             = "infolist"
	funcNameResponselist         = "responselist"

	funcNameRenameResponseFile = "renameResponseFile"
)

/// 声明插件结构体
type MockServicePlugin struct{}

/// 指定为go-flutter插件
var _ flutter.Plugin = &MockServicePlugin{}

/// 初始化插件
func (MockServicePlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	ch := plugin.NewMethodChannel(messenger, channelName, plugin.StandardMethodCodec{})

	ch.HandleFunc(funcNameRun, runFunc)
	ch.HandleFunc(funcNameClose, closeFunc)
	ch.HandleFunc(funcNameLoad, loadFunc)

	ch.HandleFunc(funcNameUpdateInfo, updateInfoFunc)
	ch.HandleFunc(funcNameUpdateAllInfo, updateAllInfoFunc)
	ch.HandleFunc(funcNameSaveInfo, saveInfoFunc)

	ch.HandleFunc(funcNameHostlist, hostListFunc)
	ch.HandleFunc(funcNameSetDefaultTargetHost, setDefaultTargetHostFunc)
	ch.HandleFunc(funcNameInfolist, infoListFunc)
	ch.HandleFunc(funcNameResponselist, responseListFunc)

	ch.HandleFunc(funcNameRenameResponseFile, renameResponseFileFunc)

	// 用于向flutter端发送消息
	mockservice.Channel = ch

	return nil
}

/// 启动服务
func runFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("runFunc()")
	mockservice.MockService()

	return "start", nil
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

	// 参数转interface{}类型Key的map
	mArgs := arguments.(map[interface{}]interface{})

	// interface{}类型Key的map转为string类型Key的map
	m, err := moist.ToStringKeyMap(mArgs)
	if err != nil {
		log.Println(err)
	}

	// map转为struct
	var info mockservice.MockServiceInfo
	err = moist.MapToStruct(m, &info)
	if err != nil {
		log.Println(err)
	}

	// 更新模拟服务信息
	mockservice.UpdateMockServiceInfo(info)
	return true, nil
}

/// 更新所有模拟服务信息
func updateAllInfoFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("updateAllInfoFunc()")

	// 参数
	mArgs := arguments.([]interface{})

	var infoSlice []mockservice.MockServiceInfo
	for _, v := range mArgs {

		// 单个切片元素转换为interface{}类型Key的map
		mTmp := v.(map[interface{}]interface{})

		// interface{}类型Key的map转为string类型Key的map
		m, err := moist.ToStringKeyMap(mTmp)
		if err != nil {
			log.Println(err)
		}

		// map转为struct
		var info mockservice.MockServiceInfo
		err = moist.MapToStruct(m, &info)
		if err != nil {
			log.Println(err)
		}

		infoSlice = append(infoSlice, info)
	}

	// 更新全部模拟服务信息
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

/// 设置默认目标主机
func setDefaultTargetHostFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("setDefaultTargetHostFunc()")

	// 参数
	m := arguments.(map[interface{}]interface{})
	// 目标主机
	targetHost, ok := m["targetHost"].(string)
	if !ok {
		log.Println("ERROR: [targetHost]")
	}

	v := mockservice.SetDefaultTargetHost(targetHost)

	return v, nil
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
	url, ok := m["url"].(string)
	if !ok {
		log.Println("ERROR: [url]")
	}
	method, ok := m["method"].(string)
	if !ok {
		log.Println("ERROR: [method]")
	}

	var slice []interface{}
	list, err := mockservice.LoadResponseFile(url, method)
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

/// 重命名响应文件
func renameResponseFileFunc(arguments interface{}) (reply interface{}, err error) {
	log.Println("renameResponseFileFunc()")

	// 参数
	m := arguments.(map[interface{}]interface{})
	// 响应文件
	responseFile, ok := m["responseFile"].(string)
	if !ok {
		log.Println("ERROR: [responseFile]")
	}
	// 新文件名
	fileName, ok := m["fileName"].(string)
	if !ok {
		log.Println("ERROR: [fileName]")
	}

	err = mockservice.RenameResponseFile(responseFile, fileName)
	if err != nil {
		log.Println(err)
		return false, err
	}

	return true, nil
}

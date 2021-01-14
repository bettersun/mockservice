package mockservice

import (
	"fmt"
	"log"
	"time"

	"github.com/bettersun/moist"
	"github.com/go-flutter-desktop/go-flutter/plugin"
)

// 指向go-flutter的Plugin
// 用于向Flutter端发送消息
var Channel *plugin.MethodChannel

// 函数名: 向Flutter发送通知表示信息
var funcNameNotify = "notify"

// 函数名: 向Flutter发送通知，添加新的模拟服务信息
var funcNameNotifyAddMockServiceInfo = "notifyAddMockServiceInfo"

//向Flutter发送通知表示信息
func Notify(message string) {

	// 指向go-flutter的Plugin后才可使用
	if Channel != nil {
		// 消息格式
		notification := fmt.Sprintf("[Go]%v: %v", moist.NowYmdHmsSlash(), message)

		// 向通道发送消息
		err := Channel.InvokeMethod(funcNameNotify, notification)
		if err != nil {
			log.Println(fmt.Sprintf("向Flutter端发送消息失败[%v]", funcNameNotify))
		}
	}
}

// 向Flutter发送通知，添加新的模拟服务信息
func NotifyAddMockServiceInfo(info interface{}) {

	// 指向go-flutter的Plugin后才可使用
	if Channel != nil {
		entity, err := moist.StructToIfKeyMap(info)
		if err != nil {
			log.Println(err)
		}

		// 休眠0.2秒
		time.Sleep(time.Millisecond * 200)

		// 向通道发送消息
		err = Channel.InvokeMethod(funcNameNotifyAddMockServiceInfo, entity)
		if err != nil {
			log.Println(fmt.Sprintf("向Flutter端发送消息失败[%v]", funcNameNotifyAddMockServiceInfo))
		}
	}
}

package mockservice

import (
	"fmt"
	"log"

	"github.com/bettersun/moist"
	"github.com/go-flutter-desktop/go-flutter/plugin"
)

// 指向go-flutter的Plugin
// 用于向Flutter端发送消息
var Channel *plugin.MethodChannel

// 向Flutter发送消息的方法名
var funcNameNotify = "notify"

// 向Flutter发送消息
func Notify(message string) {

	// 指向go-flutter的Plugin后才可使用
	if Channel != nil {
		// 消息格式
		notification := fmt.Sprintf("[Go]%v: %v", moist.NowYmdHmsSlash(), message)

		// 向通道发送消息
		err := Channel.InvokeMethod(funcNameNotify, notification)
		if err != nil {
			log.Println("向Flutter端发送消息失败")
		}
	}
}

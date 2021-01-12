package mockservice

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/bettersun/moist"
	yml "github.com/bettersun/moist/yaml"
)

/// 请求文件路径
func pathRequest() string {
	pathRequest := fmt.Sprintf("%v%v", moist.CurrentDir(), requestPath)
	return pathRequest
}

/// 反斜线转下划线
/// 用于 URI对应的目录
func EscapseSlash(s string) string {
	result := strings.Replace(strings.ReplaceAll(s, "/", "_"), "_", "", 1)
	return result
}

/// URL对应的响应文件路径
func pathURLResponse(uri string, method string) string {
	// 输出子目录
	pURI := EscapseSlash(uri)

	//递归创建目录
	p := fmt.Sprintf("%v/%v/%v", responsePath, pURI, method)

	return p
}

/// URL对应的响应文件路径
func fullPathURLResponse(uri string, method string) string {

	//递归创建目录
	p := fmt.Sprintf("%v%v", moist.CurrentDir(), pathURLResponse(uri, method))

	return p
}

/// 输出请求到文件
func OutRequest(r *http.Request) {

	uri := r.URL.String()
	// 输出子目录
	fileName := EscapseSlash(uri)

	//递归创建目录
	filePath := fmt.Sprintf("%v/%v", pathRequest(), fileName)

	// 完整路径
	fileFullPath := fmt.Sprintf("%v/request_%v.json", filePath, moist.NowYmdHms())
	// log.Println(fileFullPath)

	if !moist.IsExist(fileFullPath) {
		// 创建目录
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			log.Printf("创建目录失败。[%v][%v]", filePath, err)
		}
	}

	// 读取请求体
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Print("ioutil.ReadAll() ", err.Error())
	}
	body := string(data)

	m := make(map[string]interface{})
	m["header"] = r.Header
	m["body"] = body

	// 输出请求到文件
	err = moist.OutJson(fileFullPath, m)
	if err != nil {
		log.Print(err)
	}
}

/// 输出响应体到文件
func OutResponseBody(method string, uri string, body []byte) {

	path := fullPathURLResponse(uri, method)
	// log.Println(path)
	if !moist.IsExist(path) {
		// 创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			log.Printf("创建目录失败。[%v][%v]", path, err)
		}
	}

	// 完整路径
	fileFullPath := fmt.Sprintf("%v/body_%v_%v.json", path, method, moist.NowYmdHms())
	// log.Println(fileFullPath)

	content := string(body)
	m, err := moist.JsonToMap(content)
	// 字符串能转换成JSON时，输出Map(或Map切片)到文件
	if err == nil {
		errOut := moist.OutJson(fileFullPath, m)
		if errOut != nil {
			log.Print(errOut)
		}
		return
	}

	// 字符串不能转换成JSON时，直接输出响应体到文件
	if err != nil {
		errOut := moist.WriteFile(fileFullPath, []string{content})
		if errOut != nil {
			log.Print(errOut)
		}
		return
	}
}

/// 输出响应到文件
func OutResponseHeader(file string, mHeader map[string]http.Header) {

	filePath := fmt.Sprintf("%v%v/%v", moist.CurrentDir(), responsePath, file)

	// 输出响应头
	err := moist.OutJson(filePath, mHeader)
	if err != nil {
		log.Print(err)
	}

	return
}

/// 获取URL的响应文件列表
func LoadResponseFile(uri string, method string) ([]string, error) {

	// URL对应响应目录下的URI对应的目录
	path := fullPathURLResponse(uri, method)

	var file []string
	sub, err := ioutil.ReadDir(path)
	if err != nil {
		log.Printf("目录不存在，或打开错误。[%v]", path)
		// 不返回error
		return file, nil
	}

	// 响应文件所在的目录
	filePath := pathURLResponse(uri, method)

	for _, f := range sub {
		if !f.IsDir() {
			fname := fmt.Sprintf("%v/%v", filePath, f.Name())
			file = append(file, fname)
		}
	}

	return file, nil
}

/// 保存模拟服务信息
func OutputMockServiceInfo(inputFile string, infoSlice []MockServiceInfo) error {

	bkFile := fmt.Sprintf("%v/backup/input_%v.yml", moist.CurrentDir(), moist.NowYmdHms())
	// log.Println(bkFile)
	err := moist.CopyFile(inputFile, bkFile)
	if err != nil {
		log.Println(err)
		return err
	}

	// 备份成功后覆盖当前yml文件
	if moist.IsExist(bkFile) {
		err = yml.OutYaml(inputFile, infoSlice)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

/// 保存目标主机
func OutputHost(file string, hostSlice []string) error {
	bkFile := fmt.Sprintf("%v/backup/host_%v.yml", moist.CurrentDir(), moist.NowYmdHms())
	// log.Println(bkFile)
	err := moist.CopyFile(file, bkFile)
	if err != nil {
		log.Println(err)
		return err
	}

	// 备份成功后覆盖当前yml文件
	if moist.IsExist(bkFile) {
		err = yml.OutYaml(file, hostSlice)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

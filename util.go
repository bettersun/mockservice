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
	"github.com/sirupsen/logrus"
)

/// 获取URL里的地址部分
///  去掉URL里?后面的内容
func baseURL(s string) string {
	tmp := ""

	// 第一个问号的位置
	qIndex := strings.Index(s, "?")

	if qIndex > 0 {
		tmp = s[0:qIndex]
	} else {
		tmp = s
	}

	return tmp
}

/// 反斜线转下划线
/// 去掉第一个下划线
func escapseURL(s string) string {
	// 去掉问号后面的部分
	s = baseURL(s)

	// 反斜线转下划线，去掉第一个下划线
	result := strings.Replace(strings.ReplaceAll(s, "/", "_"), "_", "", 1)
	return result
}

/// 请求信息文件存放目录
func pathURLRequest(url string, method string) string {
	pURL := escapseURL(url)
	p := fmt.Sprintf("%v/%v/%v_%v", moist.CurrentDir(), pathRequest, pURL, method)

	return p
}

/// 响应信息文件存放目录(绝对)
func pathURLResponse(url string, method string) string {
	p := fmt.Sprintf("%v/%v", moist.CurrentDir(), pathRelativeURLResponse(url, method))

	return p
}

/// 响应信息文件存放目录(相对)
func pathRelativeURLResponse(url string, method string) string {
	pURL := escapseURL(url)
	p := fmt.Sprintf("%v/%v_%v", pathResponse, pURL, method)

	return p
}

/// 请求信息文件名
func fileRequest() string {
	f := fmt.Sprintf("req_%v.json", moist.NowMdHms())
	return f
}

/// 响应信息文件名
func fileResponse(isJSON bool) string {
	var f string

	if isJSON {
		f = fmt.Sprintf("body_%v.json", moist.NowMdHms())
	} else {
		f = fmt.Sprintf("body_%v.txt", moist.NowMdHms())
	}

	return f
}

/// 响应头信息文件
func filePathResponseHeader() string {
	f := fmt.Sprintf("%v/%v/%v", moist.CurrentDir(), pathResponseHeader, fileResponseHeader)
	return f
}

/// 模拟服务通用响应头信息文件
func filePathCommonResponseHeader() string {
	f := fmt.Sprintf("%v/%v/%v", moist.CurrentDir(), pathResponseHeader, fileCommonResponseHeader)
	return f
}

/// 输出请求到文件
func OutRequest(url string, method string, header http.Header, body string) {
	path := pathURLRequest(url, method)
	// log.Println(path)
	if !moist.IsExist(path) {
		//递归创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logger.WithFields(logrus.Fields{logFieldPath: path, logFieldError: err}).Warn("创建目录失败")
		}
	}

	isBodyJSON := false
	mBody, err := moist.JsonToMap(body)
	if err != nil {
		logger.Debug("请求体非JSON")
	}
	if err == nil {
		isBodyJSON = true
	}

	// 输出内容
	m := make(map[string]interface{})
	m["url"] = url
	m["header"] = header
	m["method"] = method
	// 请求体为JSON时
	if isBodyJSON {
		m["body"] = mBody
	} else {
		m["body"] = body
	}

	// 文件完整路径
	fileFullPath := fmt.Sprintf("%v/%v", path, fileRequest())
	// log.Println(fileFullPath)

	// 输出请求到文件
	err = moist.OutJson(fileFullPath, m)
	if err != nil {
		log.Print(err)
	}
}

/// 输出响应体到文件
func OutResponseBody(method string, url string, isJSON bool, body []byte) {
	path := pathURLResponse(url, method)
	// log.Println(path)
	if !moist.IsExist(path) {
		// 创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logger.WithFields(logrus.Fields{logFieldPath: path, logFieldError: err}).Warn("创建目录失败")
		}
	}

	// 文件完整路径
	fileFullPath := fmt.Sprintf("%v/%v", path, fileResponse(isJSON))
	// log.Println(fileFullPath)

	content := string(body)

	isSaved := false
	if isJSON {
		m, err := moist.JsonToMap(content)
		if err != nil {
			logger.WithFields(logrus.Fields{logFieldFile: fileFullPath, logFieldError: err}).Warn("响应信息转换JSON失败")
		}

		// 字符串能转换成JSON时，输出Map(或Map切片)到文件
		if err == nil {
			errOut := moist.OutJson(fileFullPath, m)
			if errOut != nil {
				logger.WithFields(logrus.Fields{logFieldFile: fileFullPath, logFieldError: errOut}).Warn("保存响应信息失败")
			}
			return
		}

		isSaved = true
	}

	// 上面的处理未保存成功时
	if !isSaved {
		// 保存为普通文件内容
		err := moist.WriteFile(fileFullPath, []string{content})
		if err != nil {
			logger.WithFields(logrus.Fields{logFieldFile: fileFullPath, logFieldError: err}).Warn("保存响应信息失败")
		}
	}
}

/// 输出响应到文件
func OutResponseHeader(mHeader map[string]http.Header) {
	path := pathResponseHeader
	// log.Println(path)
	if !moist.IsExist(path) {
		// 创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			logger.WithFields(logrus.Fields{logFieldPath: path, logFieldError: err}).Warn("创建目录失败")
		}
	}

	// 文件完整路径
	fileFullPath := filePathResponseHeader()

	// 输出响应头
	err := moist.OutJson(fileFullPath, mHeader)
	if err != nil {
		log.Print(err)
	}

	return
}

/// 获取URL的响应文件列表
func LoadResponseFile(url string, method string) ([]string, error) {
	// URL对应响应目录下的URL对应的目录
	path := pathURLResponse(url, method)

	var file []string
	sub, err := ioutil.ReadDir(path)
	if err != nil {
		logger.WithFields(logrus.Fields{logFieldPath: path, logFieldError: err}).Warn("目录不存在，或打开错误")

		// 不返回error
		return file, nil
	}

	p := pathRelativeURLResponse(url, method)
	// 文件列表
	for _, f := range sub {
		if !f.IsDir() {
			fname := fmt.Sprintf("%v/%v", p, f.Name())
			file = append(file, fname)
		}
	}

	return file, nil
}

/// 保存模拟服务信息
func OutputMockServiceInfo(config Config, infoSlice []MockServiceInfo) error {
	// 备份
	bkFileName := strings.Replace(config.InfoFile, ".", fmt.Sprintf("_%v.", moist.NowYmdHms()), 1)
	bkFile := fmt.Sprintf("%v/%v", pathBackup, bkFileName)
	// log.Println(bkFile)

	// 复制
	infoFile := fmt.Sprintf("%v/%v", moist.CurrentDir(), config.InfoFile)
	err := moist.CopyFile(infoFile, bkFile)
	if err != nil {
		log.Println(err)
		return err
	}

	// 备份成功后覆盖当前yml文件
	if moist.IsExist(bkFile) {
		err = yml.OutYaml(infoFile, infoSlice)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

/// 保存目标主机
func OutputHost(config Config, hostSlice []string) error {
	// 备份
	bkFileName := strings.Replace(config.HostFile, ".", fmt.Sprintf("_%v.", moist.NowYmdHms()), 1)
	bkFile := fmt.Sprintf("%v/%v", pathBackup, bkFileName)
	// log.Println(bkFile)

	// 复制
	hostFile := fmt.Sprintf("%v/%v", moist.CurrentDir(), config.HostFile)
	err := moist.CopyFile(hostFile, bkFile)
	if err != nil {
		log.Println(err)
		return err
	}

	// 备份成功后覆盖当前yml文件
	if moist.IsExist(bkFile) {
		err = yml.OutYaml(hostFile, hostSlice)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

/// 保存目标主机
func RenameResponseFile(responseFile string, fileName string) error {

	file := fmt.Sprintf("%v/%v", moist.CurrentDir(), responseFile)
	path := moist.FilePath(file, "/")
	newFile := fmt.Sprintf("%v/%v", path, fileName)

	log.Println(file)
	log.Println(path)
	log.Println(newFile)

	//重命名文件
	err := os.Rename(file, newFile)
	if err != nil {
		logger.WithFields(logrus.Fields{"响应文件": responseFile, "新文件名": fileName, logFieldError: err}).Error("重命名响应文件失败")
		return err
	}

	return nil
}

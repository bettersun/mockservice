# mockservice
一个简单的模拟服务工具

依赖Github仓库:  
https://github.com/bettersun/moist

## 简介

前端连接到该工具，该工具实现请求转发、请求和响应信息的保存。  
该工具也可提供模拟服务，即作为一个虚拟的服务器提供所需的响应数据。

## 运行工具

进入到 command 目录，运行 go build 命令后生成可执行文件，运行可执行文件即可。

## 配置文件

程序运行所需的配置，修改该文件的内容后，需要重新启动该工具以生效。  
格式为 yaml 格式。

config.yml
``` yml
# 此文件内容修改后需重启
# 端口
port: 9999
# 日志文件
logFile: mock.log
# 日志等级
# -1: DEBUG
# 0: INFO
# 1: WARN
# 2: ERROR
# 3: FATAL
logLevel: 0
# 主机文件
hostFile: host.yml
# 模拟服务信息文件
infoFile: info.yml
# 使用模拟服务通用响应头
useMockCommonResponseHeader: true
```

1. 端口

    端口为该工具的监听接口，前端连接到此工具所在主机的该端口。
    本机连接的话，手机的模拟器不能连接 127.0.0.1 或者 localhost。

    若不配置，默认使用 8527。

2. 日志文件

    日志文件为该工具运行时日志保存的文件名。

    若不配置，默认使用 mock.log。

3. 日志等级

    日志保存的等级，高于或等于该等级的日志会被保存。

    数字对应等级为：  
    -1: DEBUG
    0: INFO
    1: WARN
    2: ERROR
    3: FATAL

    若不配置，默认使用 0 (信息)。

4. 主机文件

    前端通过该工具连接的对象主机列表的文件名。

    若不配置，默认使用 host.yml。

5. 模拟服务信息文件

    前端请求的各个 URL 对应的配置信息列表的文件名。

    若不配置，默认使用 info.yml。

6. 使用模拟服务通用响应头

    该工具转发请求时，真实主机返回的响应头信息会保存。
    位置为 response_header/response_header.json。
    
    不使用该工具转发，使用该工具提供模拟服务时，对于各个 URL 的响应，需要响应头信息。
    默认提供了一个通用响应头，response_header/common_response_header.json。

    对于 URL 的请求，首先会查找是否存在真实主机的响应头。
    如果不存在，并且配置文件里使用模拟服务通用响应头为 true 时，则会使用 response_header/common_response_header.json 的内容作为模拟服务的响应头。

    对于 URL 的请求，首先会查找是否存在真实主机的响应头，如果不存在，并且配置文件里使用模拟服务通用响应头为 false 时，则不作特殊处理。

    若不配置，默认使用 false。

## 主机文件

前端需要连接的主机列表，包括主机 IP 和服务监听的端口，也可使用服务映射的根 URL。

配置文件中的配置项为： hostFile

格式为 yaml 格式。

``` yml
- http://localhost:8012
- http://127.0.0.1:8015
- http://192.168.9.12:8016
- http://www.helloworld.cn
```

## 模拟服务信息

前端请求的各个 URL（包括请求方法）对应的配置信息列表，该工具通过该配置信息来转发请求或提供模拟服务。

若不配置该文件，前端通过该工具转发请求时，会自动记录对应的的模拟服务信息并保存到文件。
有个前提是工具运行时必须指定默认目标主机，这样工具才能知道请求要转发的目标主机。

配置文件中的配置项为： infoFile

格式为 yaml 格式。

``` yml
- url: /bettersun/hello
  method: POST
  targetHost: http://localhost:8012
  useDefaultTargetHost: true
  useMockService: true
  statusCode: 200
  responseFile: json/hello.json
  description: hello_POST
- url: /bettersun/hello
  method: GET
  targetHost: http://localhost:8012
  useDefaultTargetHost: true
  useMockService: false
  statusCode: 200
  responseFile: ""
  description: hello_GET
- url: /bettersun/hello
  method: PUT
  targetHost: http://localhost:8012
  useDefaultTargetHost: true
  useMockService: false
  statusCode: 200
  responseFile: ""
  description: Hello_PUT
```

1. url

    前端请求的 URL。

2. method

    前端请求URL 时的请求方法。
    GET/POST/PUT/DELETE 等。

3. targetHost

    目标主机，为前端请求需要连接的真实主机，即该工具转发请求需要连接的真实主机。
    当使用模拟服务（useMockService 为 true）时，则由该工具提供模拟服务，不再连接目标主机。

4. useDefaultTargetHost

    使用默认目标主机标志。

    程序启动时，默认目标主机为主机文件里配置的第一个。
    可在画面中选择默认目标主机，适用于前端请求的大部分URL为同一台主机提供真实服务的情况。

    当 使用默认目标主机标志 为 true 时，则使用当前选择的默认目标主机。
    若 使用默认目标主机标志 为 false 时，则使用配置信息中的目标主机（targetHost）。

5. useMockService

    使用模拟服务标志。

    当 使用模拟服务标志 为 true 时，该工具不会转发前端请求的 URL，会对前端请求的 URL提供模拟服务。
    当 使用模拟服务标志 为 false 时，该工具会转发前端请求的 URL。

6. statusCode

    模拟服务响应状态码。

    当使用模拟服务时，该工具对于前端请求的 URL 返回的响应状态码。

7. responseFile

    模拟服务响应文件名。

    当使用模拟服务时，该工具对于前端请求的 URL 返回的响应体内容所在的文件。
    对应文件里保存模拟服务的响应体内容。

    **文件的路径规则：**

    响应体内容为 JSON 时，
    response/URL(转换)_请求方法/body_MdHms.json
    响应体内容非 JSON 时，
    response/URL(转换)_请求方法/body_MdHms.txt

    - URL(转换) 为将请求 URL 去掉首个斜线，并将斜线转为下划线后的值。
    例:
    /bettersun/hello
    转换后为
    bettersun_hello

    - 请求方法为大写。

    - MdHms 为 月日时分秒各两位的数字。

    文件例：

    response/bettersun_hello_GET/body_0112091212.txt
    response/bettersun_hello_GET/body_0112091216.json

8. description

    该 URL（包括请求方法）对应的说明或描述。

## 模拟服务通用响应头

位置是 response_header/common_response_header.json

可自定义内容，默认的三个选项为支持跨域、gzip 压缩和响应格式为json(UTF8)。

格式为 json 格式。

``` json
{
    "Access-Control-Allow-Origin": [
        "*"
    ],
    "Content-Encoding": [
        "gzip"
    ],
    "Content-Type": [
        "application/json;charset=UTF-8"
    ]
}
```

该工具转发请求时，真实主机返回的响应头信息会保存。
位置为 response_header/response_header.json。

不使用该工具转发，使用该工具提供模拟服务时，对于各个 URL 的响应，需要响应头信息。
默认提供了一个通用响应头，response_header/common_response_header.json。

对于 URL 的请求，首先会查找是否存在真实主机的响应头。
如果不存在，并且配置文件里使用模拟服务通用响应头为 true 时，则会使用 response_header/common_response_header.json 的内容作为模拟服务的响应头。

对于 URL 的请求，首先会查找是否存在真实主机的响应头，如果不存在，并且配置文件里使用模拟服务通用响应头为 false 时，则不作特殊处理。

## 请求响应记录

通过该工具转发请求到真实的目标主机时，会自动记录请求信息和响应信息。

1. 记录的请求信息

    记录的请求信息保存在下面目录。

    requeset/URL(转换)_请求方法/req_MdHms.json

    - URL(转换) 为将请求 URL 去掉首个斜线，并将斜线转为下划线后的值。  
    例:
    /bettersun/hello
    转换后为
    bettersun_hello

    - 请求方法为大写。

    - MdHms 为 月日时分秒各两位的数字。

    记录的请求信息内容为 请求的 URL、请求头和请求体。

2. 记录的响应信息

    记录的响应头信息保存在下面文件。  
    response_header/response_header.json

    文件内容是 URL(转换)_请求方法 作为键，响应头信息作为值的 Map 转换后的 JSON 文本。

    当使用模拟服务时，会首先查找该文件并从该文件中查找 URL（包括请求方法）对应的响应头信息。若该文件不存在或查找不到对应的响应头信息，并且配置文件里的使用模拟服务通用响应头为 true 时，则会读取模拟服务通用响应头文件的内容作为响应头信息返回。
    
    响应体信息保存在下面目录。  

    响应体内容为 JSON 时，
    response/URL(转换)_请求方法/body_MdHms.json
    响应体内容非 JSON 时，
    response/URL(转换)_请求方法/body_MdHms.txt

    - URL(转换) 为将请求 URL 去掉首个斜线，并将斜线转为下划线后的值。
    例:
    /bettersun/hello
    转换后为
    bettersun_hello

    - 请求方法为大写。

    - MdHms 为 月日时分秒各两位的数字。

    文件例：

    response/bettersun_hello_GET/body_0112091212.txt
    response/bettersun_hello_GET/body_0112091216.json

## 坑
Go的init()方法里向Flutter端发送消息，Flutter端接收不到。
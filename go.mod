module github.com/bettersun/mockservice

go 1.15

require (
	github.com/bettersun/moist v0.0.6
	github.com/bettersun/moist/yaml v0.0.0-20201229122944-f48691f8f589
	github.com/go-flutter-desktop/go-flutter v0.42.0
	github.com/sirupsen/logrus v1.7.0
)

// 使用本地目录改写依赖
replace github.com/bettersun/moist => ../moist

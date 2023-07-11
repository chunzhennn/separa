package plugin

import (
	"regexp"
	"testing"
)

func TestMain(t *testing.T) {
	// 定义一个 HTTP 响应头
	header := "HTTP/1.1 302 Found\r\nLocation: /123\r\nContent-Length: 0\r\n\r\n"

	// 获取 Location 字段的值
	location := getLocation(header)
	t.Logf("Location: %s", location)
	// 检查 Location 字段的值是否正确
	if location != "/123" {
		t.Errorf("Location should be [/123], but got [%s]", location)
	}
}

// 从 HTTP 响应头中获取 Location 字段的值
func getLocation(header string) string {
	// 定义正则表达式
	re, _ := regexp.Compile("Location: (.*)\r\n")

	// 查找 Location 字段的值
	match := re.FindStringSubmatch(header)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

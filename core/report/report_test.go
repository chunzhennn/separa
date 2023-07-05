package report

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestJson(t *testing.T) {
	ResultKV.KV = make(map[string]*ResultUnit)

	PushResult("13.23.9.62", NewResultUnit())

	// resultUnit := Get("13.23.9.62")
	// resultUnit.Honeypot = make([]string, 0)
	// resultUnit.Honeypot = append(resultUnit.Honeypot, "test")

	// 将 map 序列化成 JSON 格式
	data, err := json.Marshal(ResultKV.KV)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 打印序列化后的 JSON 字符串
	fmt.Println(string(data))
}

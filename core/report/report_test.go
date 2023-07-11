package report

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"gopkg.in/yaml.v2"
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

type config struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Author  string `yaml:"author"`
}

func TestYaml(t *testing.T) {

	// 将 YAML 格式的文本写入文件
	content, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}
	var cfg config
	err = yaml.Unmarshal(content, &cfg)
	fmt.Printf("%s\n", cfg.Name)
	fmt.Printf("%s\n", cfg.Author)
	fmt.Printf("%s\n", cfg.Version)
	// fmt.Println("Config file written successfully.")

}

package main

import (
	"fmt"
	"strings"
)

func main() {
	name := "nginx"
	name = strings.ToLower(strings.ReplaceAll(name, "\t", ""))
	// 可能有 version 信息
	index := strings.LastIndex(name, "/")
	if index == -1 {
		fmt.Println(name)
		return
	}
	prod := name[:index]
	version := name[index+1:]
	fmt.Println(prod, version)
}

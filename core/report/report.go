package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"separa/common"
	"separa/common/log"
	"time"
)

var ResultKV struct {
	KV map[string]*ResultUnit // Change the map value type to a pointer to ResultUnit
}

type ResultUnit struct {
	Services   []ServiceUnit
	Deviceinfo string
	Honeypot   []string
	Timestamp  string
}

type ServiceUnit struct {
	Port       int
	Protocol   string
	ServiceApp []string
}

func NewServiceUnit(port int, protocol string, serviceApp []string) *ServiceUnit {
	ServiceUnit := &ServiceUnit{
		Port:       port,
		Protocol:   protocol,
		ServiceApp: serviceApp,
	}
	return ServiceUnit
}

func NewResultUnit() *ResultUnit {
	currentTime := time.Now()
	ResultUnit := &ResultUnit{
		Timestamp: currentTime.Format("2006-01-02 15:04:05"),
	}
	return ResultUnit
}

func Init() {
	ResultKV.KV = map[string]*ResultUnit{}
}

func PushIP(ip string) {
	ResultKV.KV[ip] = NewResultUnit()
}

func PushResult(ip string, result *ResultUnit) {
	ResultKV.KV[ip] = result
}

func Get(ip string) *ResultUnit {
	return ResultKV.KV[ip]
}

func AppendService(ip string, service *ServiceUnit) {
	if ResultKV.KV[ip].Services == nil {
		ResultKV.KV[ip].Services = make([]ServiceUnit, 0)
	}
	ResultKV.KV[ip].Services = append(ResultKV.KV[ip].Services, *service)
}

func Save() {
	data, err := json.MarshalIndent(ResultKV.KV, "", "    ")
	if err != nil {
		log.Log.Printf("Error: %s", err)
		return
	}
	file, err := os.OpenFile(common.Setting.Output, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Log.Printf("Error: %s", err)
		return
	}
	defer file.Close()
	_, err = io.WriteString(file, string(data))
	if err != nil {
		fmt.Println(err)
		return
	}
}

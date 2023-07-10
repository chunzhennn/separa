package report

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"separa/common"
	"separa/common/log"
	"strings"
	"time"
)

var ResultKV struct {
	KV map[string]*ResultUnit // Change the map value type to a pointer to ResultUnit
}

type ResultUnit struct {
	Services   []ServiceUnit `json:"services"`
	Deviceinfo string        `json:"deviceinfo"`
	Honeypot   []string      `json:"honeypot"`
	Timestamp  string        `json:"timestamp"`
}

type ServiceUnit struct {
	Port       int      `json:"port"`
	Protocol   string   `json:"protocol"`
	ServiceApp []string `json:"service_app"`
}

func AttachVersion(app string) (string, string) {
	mayBeVersion := strings.Split(app, "/")
	if len(mayBeVersion) > 1 {
		return mayBeVersion[0], mayBeVersion[1]
	} else {
		return mayBeVersion[0], "N"
	}
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

func AppendHonypot(ip string, honeypot string) {
	if ResultKV.KV[ip].Honeypot == nil {
		ResultKV.KV[ip].Honeypot = make([]string, 0)
	}
	ResultKV.KV[ip].Honeypot = append(ResultKV.KV[ip].Honeypot, honeypot)
}

func UpdateDeviceinfo(ip string, deviceinfo string) {
	ResultKV.KV[ip].Deviceinfo = deviceinfo
}

func Save() {
	data, err := json.MarshalIndent(ResultKV.KV, "", "    ")
	if err != nil {
		log.Log.Printf("Error: %s", err)
		return
	}
	path := filepath.FromSlash(common.Setting.Output)
	dir := filepath.Dir(path)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Log.Printf("Error: %s", err)
			return
		}
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
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

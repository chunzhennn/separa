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
	"sync"
	"time"
)

var ResultKV struct {
	sync.RWMutex
	KV map[string]*ResultUnit // Change the map value type to a pointer to ResultUnit
}

type ResultUnit struct {
	Services   []ServiceUnit `json:"services"`
	Deviceinfo []string      `json:"deviceinfo"`
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
	if len(serviceApp) > 0 {
		ServiceUnit := &ServiceUnit{
			Port:       port,
			Protocol:   protocol,
			ServiceApp: serviceApp,
		}
		return ServiceUnit
	} else {
		ServiceUnit := &ServiceUnit{
			Port:     port,
			Protocol: protocol,
		}
		return ServiceUnit
	}
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
	ResultKV.Lock()
	ResultKV.KV[ip] = NewResultUnit()
	ResultKV.Unlock()
}

func PushResult(ip string, result *ResultUnit) {
	ResultKV.Lock()
	ResultKV.KV[ip] = result
	ResultKV.Unlock()
}

func Get(ip string) *ResultUnit {
	return ResultKV.KV[ip]
}

func AppendService(ip string, service *ServiceUnit) {
	ResultKV.Lock()
	if ResultKV.KV[ip] == nil {
		ResultKV.KV[ip] = NewResultUnit()
	}

	if ResultKV.KV[ip].Services == nil {
		ResultKV.KV[ip].Services = make([]ServiceUnit, 0)
	}
	ResultKV.KV[ip].Services = append(ResultKV.KV[ip].Services, *service)
	ResultKV.Unlock()
}

func AppendHonypot(ip string, honeypot string) {
	ResultKV.Lock()
	if ResultKV.KV[ip] == nil {
		ResultKV.KV[ip] = NewResultUnit()
	}

	if ResultKV.KV[ip].Honeypot == nil {
		ResultKV.KV[ip].Honeypot = make([]string, 0)
	}
	ResultKV.KV[ip].Honeypot = append(ResultKV.KV[ip].Honeypot, honeypot)
	ResultKV.Unlock()
}

func AppendDeviceinfo(ip string, deviceinfo string) {
	ResultKV.Lock()
	if ResultKV.KV[ip] == nil {
		ResultKV.KV[ip] = NewResultUnit()
	}

	if ResultKV.KV[ip].Deviceinfo == nil {
		ResultKV.KV[ip].Deviceinfo = make([]string, 0)
	}
	for _, v := range ResultKV.KV[ip].Deviceinfo {
		if v == deviceinfo {
			return
		}
	}
	ResultKV.KV[ip].Deviceinfo = append(ResultKV.KV[ip].Deviceinfo, deviceinfo)
	ResultKV.Unlock()
}

func Save() {
	data, err := json.MarshalIndent(ResultKV.KV, "", "    ")
	if err != nil {
		log.Err("Error: %s", err)
		return
	}
	path := filepath.FromSlash(common.Setting.Output)
	dir := filepath.Dir(path)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			log.Err("Error: %s", err)
			return
		}
	}
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Err("Error: %s", err)
		return
	}
	defer file.Close()
	_, err = io.WriteString(file, string(data))
	if err != nil {
		fmt.Println(err)
		return
	}
}

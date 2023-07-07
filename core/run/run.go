package run

import (
	"fmt"
	"net"
	"net/url"
	"separa/common"
	"separa/common/log"
	"separa/common/uri"
	"separa/core/report"
	"separa/core/scanner"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/lcvvvv/appfinger"

	"github.com/lcvvvv/gonmap"
)

var (
	IPScanner    *scanner.IPScanner
	ProtoScanner *scanner.ProtoScanner
	URLScanner   *scanner.URLScanner
)

const (
	scannerNum = 3
)

// main entry point
func Start(targets *[]string) {
	var wg = &sync.WaitGroup{}
	wg.Add(scannerNum)

	// initialize the scanner
	initialize(wg)

	// run the scanner
	log.Log.Printf("IPScanner start")
	go IPScanner.Run()
	log.Log.Printf("ProtoScanner start")
	go ProtoScanner.Run()
	log.Log.Printf("URLScanner start")
	go URLScanner.Run()
	time.Sleep(time.Second * 1)

	// distribute the target
	distributeTraget(&common.Setting.Target)

	// check if the scanner is done
	go checkStop()
	wg.Wait()

	defer report.Save()
}

func IPScannerInit(wg *sync.WaitGroup) {
	config := scanner.DefaultConfig()
	config.Timeout = 200 * time.Millisecond
	IPScanner = scanner.NewIPScanner(config, 255)
	IPScanner.Defer(func() {
		wg.Done()
	})
	IPScanner.HandlerActive = func(addr net.IP) {
		log.Log.Printf("IPScanner active: %s", addr.String())
		report.PushIP(addr.String())
		for _, port := range common.Setting.Port {
			ProtoScanner.Push(addr, port)
		}
	}
}

func getTimeout(i int) time.Duration {
	switch {
	case i > 10000:
		return time.Millisecond * 200
	case i > 5000:
		return time.Millisecond * 300
	case i > 1000:
		return time.Millisecond * 400
	default:
		return time.Millisecond * 500
	}
}

func ProtoScannerInit(wg *sync.WaitGroup) {
	config := scanner.DefaultConfig()
	config.Timeout = getTimeout(len(common.Setting.Port))
	ProtoScanner = scanner.NewProtoScanner(config, 800)
	ProtoScanner.Defer(func() {
		wg.Done()
	})
	ProtoScanner.HandlerOpen = func(addr net.IP, port int) {
		log.Log.Printf("ProtoScanner open: %s:%d", addr.String(), port)
		protocol := gonmap.GuessProtocol(port)
		report.AppendService(addr.String(), report.NewServiceUnit(port, protocol, nil))
	}

	ProtoScanner.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		// log.Log.Printf("ProtoScanner matched: %s:%d", addr.String(), port)
		// log.Log.Printf("ProtoScanner matched:%+v", response.FingerPrint)
		// log.Log.Printf("ProtoScanner matched:%+v", response)
		var protocol string
		if response.FingerPrint.Service != "" {
			protocol = response.FingerPrint.Service
		} else {
			protocol = gonmap.GuessProtocol(port)
		}
		URLRaw := fmt.Sprintf("%s://%s:%d", response.FingerPrint.Service, addr.String(), port)
		URL, _ := url.Parse(URLRaw)
		if appfinger.SupportCheck(URL.Scheme) == true {
			URLScanner.Push(URL, response, nil, nil)
			return
		}
		report.AppendService(addr.String(), report.NewServiceUnit(port, protocol, nil))
	}
}
func URLScannerInit(wg *sync.WaitGroup) {
	config := scanner.DefaultConfig()
	config.Threads = config.Threads/2 + 1
	URLScanner = scanner.NewURLScanner(config)
	URLScanner.HandlerMatched = func(URL *url.URL, banner *appfinger.Banner, finger *appfinger.FingerPrint) {
		host := URL.Hostname()
		port := URL.Port()
		scheme := URL.Scheme
		if port == "" {
			port = "80" // 默认端口号
		}
		iPort, _ := strconv.Atoi(port)

		// 去重处理
		productMap := make(map[string]string)

		var productName []string
		for _, name := range finger.ProductName {
			// 去除finger.ProductName里的 '\t' 并小写
			name = strings.ToLower(strings.ReplaceAll(name, "\t", ""))
			// 可能有 version 信息
			index := strings.LastIndex(name, "/")
			// 如果 version 为空，则默认为 N
			version := "N"
			prod := name
			if index != -1 {
				prod = name[:index]
				version = name[index+1:]
			}

			// 如果 productMap 中已经存在 prod，则比较 version 是否为 N，为 N 且新的不为 N 则替换
			_, ok := productMap[prod]
			if ok && productMap[prod] != "N" {
				continue
			}
			productMap[prod] = version
		}

		for k, v := range productMap {
			name := k + "/" + v
			productName = append(productName, name)
		}

		report.AppendService(host, report.NewServiceUnit(iPort, scheme, productName))
	}
	URLScanner.HandlerError = func(url *url.URL, err error) {
		log.Log.Printf("URLScanner %s Error: %s", url.String(), err)
	}
	URLScanner.Defer(func() {
		wg.Done()
	})
}

func initialize(wg *sync.WaitGroup) {
	report.Init()
	common.ConfigInit()
	IPScannerInit(wg)
	ProtoScannerInit(wg)
	URLScannerInit(wg)
}

func distributeTraget(targets *[]string) {
	for _, target := range *targets {
		PushTarget(target)
	}
}

func checkStop() {
	for {
		time.Sleep(time.Second * 10)
		if IPScanner.RunningThreads() == 0 && !IPScanner.IsDone() {
			IPScanner.Stop()
			log.Log.Printf("IPScanner finish")
		}
		if !IPScanner.IsDone() {
			continue
		}
		if ProtoScanner.RunningThreads() == 0 && !ProtoScanner.IsDone() {
			ProtoScanner.Stop()
			log.Log.Printf("ProtoScanner finish")
		}
		if !ProtoScanner.IsDone() {
			continue
		}
		if URLScanner.RunningThreads() == 0 && !URLScanner.IsDone() {
			URLScanner.Stop()
			log.Log.Printf("URLScanner finish")
		}
	}
}

// push target to the scanner， string to net.IP
func PushTarget(expr string) {
	if expr == "" {
		return
	}
	if uri.IsIPv4(expr) {
		IPScanner.Push(net.ParseIP(expr))
		return
	}
	if uri.IsCIDR(expr) {
		for _, ip := range uri.CIDRToIP(expr) {
			PushTarget(ip.String())
		}
		return
	}
	if uri.IsIPRanger(expr) {
		for _, ip := range uri.RangerToIP(expr) {
			PushTarget(ip.String())
		}
		return
	}
}

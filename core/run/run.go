package run

import (
	"net"
	"separa/common"
	"separa/common/log"
	"separa/common/uri"
	"separa/core/report"
	"separa/core/scanner"
	"sync"
	"time"

	"github.com/lcvvvv/gonmap"
)

var (
	IPScanner    *scanner.IPScanner
	ProtoScanner *scanner.ProtoScanner
)

const (
	scannerNum = 2
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

	// distribute the target
	distributeTraget(&common.Setting.Target)

	// check if the scanner is done
	go checkStop()
	wg.Wait()

	defer report.Save()
}

func IPScannerInit(wg *sync.WaitGroup) {
	IPScanner = scanner.NewIPScanner(scanner.DefaultConfig(), 16)
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

func ProtoScannerInit(wg *sync.WaitGroup) {
	ProtoScanner = scanner.NewProtoScanner(scanner.DefaultConfig(), 16)
	ProtoScanner.Defer(func() {
		wg.Done()
	})
	ProtoScanner.HandlerOpen = func(addr net.IP, port int) {
		log.Log.Printf("ProtoScanner open: %s:%d", addr.String(), port)
		protocol := gonmap.GuessProtocol(port)
		report.AppendService(addr.String(), report.NewServiceUnit(port, protocol, nil))
	}

	ProtoScanner.HandlerMatched = func(addr net.IP, port int, response *gonmap.Response) {
		log.Log.Printf("ProtoScanner matched: %s:%d", addr.String(), port)
		protocol := gonmap.GuessProtocol(port)
		report.AppendService(addr.String(), report.NewServiceUnit(port, protocol, nil))
	}
}

func initialize(wg *sync.WaitGroup) {
	report.Init()
	common.ConfigInit()
	IPScannerInit(wg)
	ProtoScannerInit(wg)
}

func distributeTraget(targets *[]string) {
	for _, target := range *targets {
		PushTarget(target)
	}
}

func checkStop() {
	for {
		time.Sleep(time.Second * 2)
		if IPScanner.RunningThreads() == 0 && !IPScanner.IsDone() {
			IPScanner.Stop()
			log.Log.Printf("IPScanner finish")
		}
		if ProtoScanner.RunningThreads() == 0 && !ProtoScanner.IsDone() {
			ProtoScanner.Stop()
			log.Log.Printf("ProtoScanner finish")
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

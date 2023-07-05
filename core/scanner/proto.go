package scanner

import (
	"net"
	"separa/common/log"

	"github.com/lcvvvv/gonmap"
)

type ProtoScanner struct {
	*client

	HandlerClosed     func(addr net.IP, port int)
	HandlerOpen       func(addr net.IP, port int)
	HandlerNotMatched func(addr net.IP, port int, response string)
	HandlerMatched    func(addr net.IP, port int, response *gonmap.Response)
	HandlerError      func(addr net.IP, port int, err error)
}

type ipPort struct {
	addr net.IP
	port int
}

func NewProtoScanner(config *Config, threads int) (ps *ProtoScanner) {
	ps = &ProtoScanner{
		client:            newConfig(config, threads),
		HandlerClosed:     func(addr net.IP, port int) {},
		HandlerOpen:       func(addr net.IP, port int) {},
		HandlerNotMatched: func(addr net.IP, port int, response string) {},
		HandlerMatched:    func(addr net.IP, port int, response *gonmap.Response) {},
		HandlerError:      func(addr net.IP, port int, err error) {},
	}
	ps.pool.Function = func(in interface{}) {
		nmap := gonmap.New()
		nmap.SetTimeout(config.Timeout)
		if config.DeepInspection == true {
			nmap.OpenDeepIdentify()
		}
		value := in.(ipPort)
		log.Log.Printf("scan %s:%d", value.addr.String(), value.port)
		status, response := nmap.ScanTimeout(value.addr.String(), value.port, 10*config.Timeout)
		switch status {
		case gonmap.Closed:
			ps.HandlerClosed(value.addr, value.port)
		case gonmap.Open:
			ps.HandlerOpen(value.addr, value.port)
		case gonmap.NotMatched:
			ps.HandlerNotMatched(value.addr, value.port, response.Raw)
		case gonmap.Matched:
			ps.HandlerMatched(value.addr, value.port, response)
		}
	}
	return
}

func (c *ProtoScanner) Push(ip net.IP, port int) {
	c.pool.Push(ipPort{ip, port})
}

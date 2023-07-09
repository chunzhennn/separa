package scanner

import (
	"net"
	"separa/core/plugin"
	"separa/pkg"
	"strconv"

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
		value := in.(ipPort)
		result := pkg.NewResult(value.addr.String(), strconv.Itoa(value.port))
		plugin.Dispatch(result)
		println(result.FullOutput())
	}
	return
}

func (c *ProtoScanner) Push(ip net.IP, port int) {
	c.pool.Push(ipPort{ip, port})
}

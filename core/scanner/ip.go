package scanner

import (
	"net"
	"separa/common/osping"
	"separa/common/tcpping"
)

type IPScanner struct {
	*client
	HandlerActive func(addr net.IP)
	HandlerDie    func(addr net.IP)
}

func NewIPScanner(config *Config, threads int) (ips *IPScanner) {
	ips = &IPScanner{
		client:        newConfig(config, threads),
		HandlerActive: func(addr net.IP) {},
		HandlerDie:    func(addr net.IP) {},
	}
	ips.pool.Function = func(in interface{}) {
		ip := in.(net.IP)
		if ips.config.HostDiscoverClosed {
			ips.HandlerActive(ip)
			return

		}
		if osping.Ping(ip.String()) {
			ips.HandlerActive(ip)
			return
		}
		if err := tcpping.PingPorts(ip.String(), config.Timeout); err == nil {
			ips.HandlerActive(ip)
			return
		}
		ips.HandlerDie(ip)
	}
	return
}

func (c *IPScanner) Push(ips ...net.IP) {
	for _, ip := range ips {
		c.pool.Push(ip)
	}
}

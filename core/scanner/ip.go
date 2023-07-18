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
		// 设置中开启不进行 ip 检测，默认所有 ip 均开放
		if ips.config.HostDiscoverClosed {
			ips.HandlerActive(ip)
			return

		}
		// 调用系统 ping 进行探测
		if osping.Ping(ip.String()) {
			ips.HandlerActive(ip)
			return
		}
		// 调用 tcp ping 进行探测，使用端口：22, 23, 80, 139, 512, 443, 445, 3389
		if err := tcpping.PingPorts(ip.String(), config.Timeout); err == nil {
			ips.HandlerActive(ip)
			return
		}
		ips.HandlerDie(ip)
	}
	return ips
}

func (c *IPScanner) Push(ips ...net.IP) {
	for _, ip := range ips {
		c.pool.Push(ip)
	}
}

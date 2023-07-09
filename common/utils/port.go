package utils

import (
	"strconv"
	"strings"
)

func NewPorts(s string) Ports {
	return ParsePort(s)
}

type Ports []string

func (ps Ports) String() string {
	return strings.Join(ps, ",")
}

var (
	NameMap *PortMapper = &PortMapper{}
	PortMap *PortMapper = &PortMapper{}
	TagMap  *PortMapper = &PortMapper{}
)

type PortMapper map[string][]string

func (p PortMapper) Get(name string) []string {
	return p[name]
}

func (p *PortMapper) Set(name string, ports []string) {
	(*p)[name] = ports
}

func (p *PortMapper) Append(name string, ports ...string) {
	(*p)[name] = append((*p)[name], ports...)
}

func ParsePort(portstring string) []string {
	portstring = strings.TrimSpace(portstring)
	portstring = strings.Replace(portstring, "\r", "", -1)
	return ParsePorts(strings.Split(portstring, ","))
}

func ParsePorts(ports []string) []string {
	var portSlice []string
	for _, portname := range ports {
		portSlice = append(portSlice, choicePorts(portname)...)
	}
	portSlice = expandPorts(portSlice)
	portSlice = sliceUnique(portSlice)
	return portSlice
}

// 将string格式的port range 转为单个port组成的slice
func expandPorts(ports []string) []string {
	var tmpports []string
	for _, pr := range ports {
		if len(pr) == 0 {
			continue
		}
		pr = strings.TrimSpace(pr)
		if pr[0] == 45 {
			pr = "1" + pr
		}
		if pr[len(pr)-1] == 45 {
			pr = pr + "65535"
		}
		tmpports = append(tmpports, expandPort(pr)...)
	}
	return tmpports
}

// 将string格式的port range 转为单个port组成的slice
func expandPort(port string) []string {
	var tmpports []string
	if strings.Contains(port, "-") {
		sf := strings.Split(port, "-")
		start, _ := strconv.Atoi(sf[0])
		fin, _ := strconv.Atoi(sf[1])
		for port := start; port <= fin; port++ {
			tmpports = append(tmpports, strconv.Itoa(port))
		}
	} else {
		tmpports = append(tmpports, port)
	}
	return tmpports
}

// 端口预设
func choicePorts(portname string) []string {
	var ports []string
	if portname == "all" {
		for p := range *PortMap {
			ports = append(ports, p)
		}
		return ports
	}

	if NameMap.Get(portname) != nil {
		ports = append(ports, NameMap.Get(portname)...)
		return ports
	} else if TagMap.Get(portname) != nil {
		ports = append(ports, TagMap.Get(portname)...)
		return ports
	} else {
		return []string{portname}
	}
}

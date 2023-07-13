package uri

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"strings"
)

// IsIPv4 IsIP checks if a string is either IP version 4 Alias for `net.ParseIP`
func IsIPv4(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == '.' {
			return net.ParseIP(str) != nil
		}
	}
	return false
}

// IsCIDR checks if the string is an valid CIDR notation (IPV4)
func IsCIDR(str string) bool {
	_, _, err := net.ParseCIDR(str)
	return err == nil
}

// IsIPRanger checks if the string is an valid CIDR notation (IPV4)
func IsIPRanger(str string) bool {
	ip1, ip2 := parseIPPairs(str)
	if ip1 != nil && ip2 != nil {
		return true
	}
	return false
}

func RangerToIP(ranger string) (IPs []net.IP) {
	first, last := parseIPPairs(ranger)
	return pairsToIP(first, last)
}

func CIDRToIP(cidr string) (IPs []net.IP) {
	_, network, _ := net.ParseCIDR(cidr)
	first := FirstIP(network)
	last := LastIP(network)
	return pairsToIP(first, last)
}

func LastIP(network *net.IPNet) net.IP {
	firstIP := FirstIP(network)
	mask, _ := network.Mask.Size()
	size := math.Pow(2, float64(32-mask))
	lastIP := toIP(toInt(firstIP) + uint32(size) - 1)
	return net.ParseIP(lastIP)
}

func FirstIP(network *net.IPNet) net.IP {
	return network.IP
}

// IsIPRanger parse the string is an ip pairs
// 192.168.0.1-192.168.2.255
// 192.168.0.1-255
// 192.168.0.1-2.255
func parseIPPairs(str string) (ip1 net.IP, ip2 net.IP) {
	if strings.Count(str, "-") != 1 {
		return nil, nil
	}
	r := strings.Split(str, "-")
	s1 := r[0]
	s2 := r[1]
	if ip1 = net.ParseIP(s1); ip1 == nil {
		return nil, nil
	}
	i := strings.Count(s2, ".")
	if i > 3 {
		return ip1, nil
	}
	rs1 := strings.Split(s1, ".")
	rs2 := strings.Join(append(rs1[:3-i], s2), ".")
	ip2 = net.ParseIP(rs2)
	if ip2 == nil {
		return ip1, nil
	}
	if toInt(ip1) >= toInt(ip2) {
		return nil, nil
	}
	return ip1, ip2
}

func pairsToIP(ip1, ip2 net.IP) (IPs []net.IP) {
	start := toInt(ip1)
	end := toInt(ip2)
	for i := start; i <= end; i++ {
		IPs = append(IPs, net.ParseIP(toIP(i)))
	}
	return IPs
}

// IPToInteger converts an IP address to its integer representation.
// It supports both IPv4
func toInt(ip net.IP) uint32 {
	var buf = []byte(ip)
	if len(buf) > 12 {
		buf = buf[12:]
	}
	buffer := bytes.NewBuffer(buf)
	var i uint32
	_ = binary.Read(buffer, binary.BigEndian, &i)
	return i
}

func toIP(i uint32) string {
	buf := bytes.NewBuffer([]byte{})
	_ = binary.Write(buf, binary.BigEndian, i)
	b := buf.Bytes()
	return fmt.Sprintf("%v.%v.%v.%v", b[0], b[1], b[2], b[3])
}

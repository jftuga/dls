package main

import (
	"fmt"
	"net"
)

func IsPrivateIPv4(s string) bool {
	// RFC 1918 CIDRS
	var rfc1918_cidrs []string = []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"}

	var ip net.IP = net.ParseIP(s)

	if ip == nil {
		return false
	}

	for _, cidr := range rfc1918_cidrs {
		_, net, _ := net.ParseCIDR(cidr)
		if net.Contains(ip) {
			return true
		}
	}

	return false
}

func main() {
	var addrs []string = []string{"10.1.1.1", "192.168.10.2", "172.32.0.0"}
	for _, addr := range addrs {
		fmt.Println("IsPrivateIPv4", addr, IsPrivateIPv4(addr))
	}
}

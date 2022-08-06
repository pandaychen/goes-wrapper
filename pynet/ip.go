package pynet

import (
	"net"
	"regexp"
	"strings"
)

// 获取本机IP
func GetLocalIP() (ret string) {
	localIP := "localhost"
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return localIP
	}
	for _, addr := range addrs {
		items := strings.Split(addr.String(), "/")
		if len(items) < 2 || items[0] == "127.0.0.1" {
			continue
		}
		if match, err := regexp.MatchString(`\d+\.\d+\.\d+\.\d+`, items[0]); err == nil && match {
			localIP = items[0]
		}
	}
	return localIP
}

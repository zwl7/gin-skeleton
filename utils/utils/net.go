package utils

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetMacAddrs() (macAddrs []string) {
	netInterfaces, err := net.Interfaces()
	if err != nil {
		fmt.Printf("fail to get net interfaces: %v", err)
		return macAddrs
	}

	for _, netInterface := range netInterfaces {
		macAddr := netInterface.HardwareAddr.String()
		if len(macAddr) == 0 {
			continue
		}

		macAddrs = append(macAddrs, macAddr)
	}
	return macAddrs
}

// 获取客户端 ip
// X-Forwarded-For : client ip, proxy ip, proxy ip ..., 理论上应该取逗号分隔后的第一个 ip
// X-Real-Ip : 最接近 nginx 的 ip, 如果请求经过多个代理转发，那么获取到的最后一个代理服务器的 ip
// remote addr : go net 包取的 ip, 如果经过 nginx upstream 反向代理到本机的 go 程序，则会取到 127.0.0.1
func GetRealIp(r *http.Request) string {
	ip := ""
	// X-Forwarded-For
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		proxyIps := strings.Split(xForwardedFor, ",")
		if len(proxyIps) > 0 {
			ip = proxyIps[0]
			return ip
		}
	}
	// X-Real-IP, X-Real-Ip
	ip = r.Header.Get("X-Real-IP")
	if ip == "" {
		ip = r.Header.Get("X-Real-Ip")
	}
	if ip != "" {
		return ip
	}
	// Remote Addr
	s := strings.Split(r.RemoteAddr, ":")
	if len(s) > 0 {
		ip = s[0]
	}
	return ip
}

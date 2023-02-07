package reverseproxy

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func GetIpAddress(r *http.Request) (net.IP, error) {
	addr := r.RemoteAddr
	if xReal := r.Header.Get("X-Real-Ip"); xReal != "" {
		addr = xReal
	} else if xFowarded := r.Header.Get("X-Forwarded-For"); xFowarded != "" {
		addr = xFowarded
	}

	ip := addr
	if strings.Contains(addr, ":") {
		var err error
		ip, _, err = net.SplitHostPort(addr)
		if err != nil {
			return nil, fmt.Errorf("addr: %q is not ip:port %w", addr, err)
		}
	}

	userIP := net.ParseIP(ip)
	if userIP == nil {
		return nil, fmt.Errorf("ip: %q is not a valid IP address", ip)
	}

	return userIP, nil
}

package v1

import (
	"net"
	"net/http"
)

func getClientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		// в случае ошибки вернем полный RemoteAddr
		return r.RemoteAddr
	}
	return ip
}

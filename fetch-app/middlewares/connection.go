package middlewares

import (
	"net"
	"strings"

	"github.com/beego/beego/v2/server/web/context"
)

var privateIPBlocks []*net.IPNet

func IsPrivateIP(ctx *context.Context) {
	ip := ""
	ip = strings.Split(ctx.Request.Header.Get("X-FORWARDED-FOR"), ",")[0]

	if ip == "" {
		ip = strings.Split(ctx.Request.Header.Get("X-REAL-IP"), ",")[0]
	}

	if ip == "" {
		ip = ctx.Input.IP()
	}

	ipNet := net.ParseIP(ip)
	if ipNet.IsLoopback() || ipNet.IsLinkLocalUnicast() || ipNet.IsLinkLocalMulticast() {
		return
	}

	for _, block := range privateIPBlocks {
		if block.Contains(ipNet) {
			return
		}
	}

	errAuth(ctx)

}

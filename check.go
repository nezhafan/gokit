package toolkit

import (
	"math"
	"net"
	"unicode"
)

// 是否为ASCII
func IsASCII(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] > unicode.MaxASCII {
			return false
		}
	}
	return true
}

// 是否有效数字（存在NaN和Inf的可能）
func IsNumber[T Float](f T) bool {
	return !math.IsNaN(float64(f)) && !math.IsInf(float64(f), 0)
}

// 是否国内手机号
func IsPhoneNumber(s string) bool {
	if len(s) != 11 {
		return false
	}
	if s[0] != '1' {
		return false
	}
	for i := 1; i < 11; i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

// 判断是否为内网或者本机IP
func IsLocalIP(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return false
	}
	return ip.IsPrivate() || ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast()
}

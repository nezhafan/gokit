package gokit

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"runtime"
	"strings"
)

func MustJSON(v any) string {
	raw := MustRaw(v)
	return string(raw)
}

func MustRaw(v any) json.RawMessage {
	bs, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return bs
}

// 除法，截断保留x位小数 (非四舍五入)
func DivTrunc[E Number](a, b E, precision uint8) float64 {
	if a == 0 || b == 0 {
		return 0
	}
	return Trunc(float64(a)/float64(b), precision)
}

// 截断保留x位小数 (非四舍五入)
func Trunc(n float64, precision uint8) float64 {
	pow := math.Pow(10, float64(precision))
	n = math.Trunc(n * pow)
	if !IsNumber(n) {
		return 0
	}
	return n / pow
}

// 加密
func Sha256(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

// 加密
func HmacSha256(s string, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func GetCaller(skip int) string {
	pc, file, line, ok := runtime.Caller(skip + 2)
	if !ok {
		return ""
	}

	if idx := strings.LastIndex(file, "/"); idx != -1 {
		if prev := strings.LastIndex(file[:idx], "/"); prev != -1 {
			file = file[prev+1:]
		}
	}

	fn := runtime.FuncForPC(pc)
	funcName := ""
	if fn != nil {
		funcName = fn.Name()
		if idx := strings.LastIndex(funcName, "/"); idx != -1 {
			funcName = funcName[idx+1:]
		}
	}

	return fmt.Sprintf("%s:%d (%s)", file, line, funcName)
}

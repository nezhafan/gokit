package gokit

import (
	"crypto/rand"
	"errors"
	"math/big"
)

// 随机数 [0, max) ，适用于高安全性的场景。
// 对于大量且非安全性高的使用 math/rand/v2 的 rand.IntN() 即可，速度相差几十倍
func RandomIntN(max int) int {
	if max <= 0 {
		return 0
	}

	m := big.NewInt(int64(max))
	// 在现代 Go 版本中，使用 rand.Reader 时几乎不会出错
	n, err := rand.Int(rand.Reader, m)
	if err != nil {
		panic(err)
	}
	return int(n.Int64())
}

var (
	defaultCharset = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

// 随机字符串
// 随机数方法使用本包的 RandomIntN 或 math/rand/v2 的 rand.IntN
func RandomString(intN func(int) int, size int, charset []byte) (string, error) {
	if intN == nil {
		return "", errors.New("intN can not be nil")
	}

	if size <= 0 {
		return "", errors.New("size must be greater than 0")
	}

	if len(charset) == 0 {
		charset = defaultCharset
	}

	b := make([]byte, size)
	charsetLen := len(charset)
	for i := range b {
		n := intN(charsetLen)
		b[i] = charset[n]
	}

	return string(b), nil
}

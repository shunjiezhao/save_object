package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"io"
	"net/http"
	path2 "path"
	"strconv"
	"strings"
)

func GetOffsetFromHeader(h http.Header) int64 {
	byteRange := h.Get("range")
	if len(byteRange) < 7 {
		return 0
	}
	if byteRange[:6] != "bytes=" {
		return 0
	}
	bytePos := strings.Split(byteRange[6:], "-")
	offset, _ := strconv.ParseInt(bytePos[0], 0, 64)
	return offset
}

// digest: SHA-256=
func GetHashFromHeader(h http.Header) string {
	digest := h.Get("digest")
	if len(digest) < 9 {
		return ""
	}
	if digest[:8] != "SHA-256=" {
		return ""
	}
	return digest[8:]
}

func GetSizeFromHeader(h http.Header) int64 {
	i, _ := strconv.ParseInt(h.Get("Content-Length"), 0, 64)
	return i
}

func CalculateHash(r io.Reader) string {
	h := sha256.New()
	readAll, _ := io.ReadAll(r)
	h.Write(readAll)
	// 注意需要先将byte转换为16进制表示
	return base64.StdEncoding.EncodeToString([]byte(hex.EncodeToString(h.Sum(nil))))
}
func SetAddr(addr string) (string, string) {
	return path2.Join(addr, "temp"), path2.Join(addr, "objects")
}

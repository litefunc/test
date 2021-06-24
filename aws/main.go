package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"
)

func hmacString(data, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func sha256String(s string) string {
	sum := sha256.Sum256([]byte(s))
	return fmt.Sprintf("%x", sum)
}

func sha256Hash(s string) []byte {
	h := sha256.New()
	h.Write([]byte(s))
	return h.Sum(nil)
}

func Hex(src []byte) string {

	maxEnLen := hex.EncodedLen(len(src)) // 最大编码长度
	dst := make([]byte, maxEnLen)
	n := hex.Encode(dst, src)
	return hex.EncodeToString(src)[:n]
}

func hexString(src string) string {
	by := []byte(src)
	maxEnLen := hex.EncodedLen(len(by)) // 最大编码长度
	dst := make([]byte, maxEnLen)
	n := hex.Encode(dst, by)
	return hex.EncodeToString(by)[:n]
}

func getHMAC(key []byte, data []byte) []byte {
	hash := hmac.New(sha256.New, key)
	hash.Write(data)
	return hash.Sum(nil)
}

// var t = time.Now().UTC()

var t = time.Date(2013, 5, 24, 0, 0, 0, 0, time.UTC)

var HTTPMethod = "GET\n"
var CanonicalURI = "/test.txt\n"
var CanonicalQueryString = "\n"
var CanonicalHeaders = map[string]string{
	"host":                 "examplebucket.s3.amazonaws.com",
	"range":                "bytes=0-9",
	"x-amz-content-sha256": HashPayload,
	"x-amz-date":           t.Format("20060102T150405Z"),
}

// var SignedHeaders = "host;range;x-amz-content-sha256;x-amz-date\n"
var SignedHeaders = []string{}

var HashPayload = Hex(sha256Hash(""))

func main() {
	fmt.Println("Canonical Request:")

	s := HTTPMethod + CanonicalURI + CanonicalQueryString
	for k := range CanonicalHeaders {
		SignedHeaders = append(SignedHeaders, k)
	}
	sort.Strings(SignedHeaders)
	for _, k := range SignedHeaders {
		s += fmt.Sprintf("%s:%s\n", k, CanonicalHeaders[k])
	}

	s += "\n"

	s += fmt.Sprintf("%s\n", strings.Join(SignedHeaders, ";"))
	s += HashPayload

	fmt.Println(s)

	s1 := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/us-east-1/s3/aws4_request\n%s", CanonicalHeaders["x-amz-date"], t.Format("20060102"), Hex(sha256Hash(s)))
	fmt.Println("StringToSign:")
	fmt.Println(s1)

	fmt.Println("SigningKey:")

	hash := getHMAC([]byte("AWS4"+"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"), []byte("20130524"))
	hash = getHMAC(hash, []byte("us-east-1"))
	hash = getHMAC(hash, []byte("s3"))
	hash = getHMAC(hash, []byte("aws4_request"))

	fmt.Println("signature:")

	sum := getHMAC(hash, []byte(s1))
	signature := Hex(sum)
	fmt.Println(signature)
}

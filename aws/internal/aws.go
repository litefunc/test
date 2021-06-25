package internal

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"test/logger"
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

type CanonicalRequest struct {
	HTTPMethod           string
	CanonicalURI         string
	CanonicalQueryString string
	CanonicalHeaders     map[string]string
	HashPayload          string
}

func (rec CanonicalRequest) String() string {

	s := fmt.Sprintf("%s\n%s\n%s\n", rec.HTTPMethod, rec.CanonicalURI, rec.CanonicalQueryString)
	SignedHeaders := []string{}
	for k := range rec.CanonicalHeaders {
		SignedHeaders = append(SignedHeaders, k)
	}
	sort.Strings(SignedHeaders)
	for _, k := range SignedHeaders {
		s += fmt.Sprintf("%s:%s\n", k, rec.CanonicalHeaders[k])
	}

	s += fmt.Sprintf("\n%s\n%s", strings.Join(SignedHeaders, ";"), rec.HashPayload)

	fmt.Printf("Canonical Request:\n%s\n\n", s)
	return s
}

type StringToSign struct {
	t      time.Time
	region string
	cr     string
}

func (rec StringToSign) String() string {

	s := fmt.Sprintf("AWS4-HMAC-SHA256\n%s\n%s/%s/s3/aws4_request\n%s", rec.t.Format("20060102T150405Z"), rec.t.Format("20060102"), rec.region, Hex(sha256Hash(rec.cr)))

	fmt.Printf("StringToSign:\n%s\n\n", s)
	return s
}

type Signature struct {
	t               time.Time
	SecretAccessKey string
	region          string
	sts             string
}

func (rec Signature) String() string {

	hash := getHMAC([]byte("AWS4"+rec.SecretAccessKey), []byte(rec.t.Format("20060102")))
	hash = getHMAC(hash, []byte(rec.region))
	hash = getHMAC(hash, []byte("s3"))
	hash = getHMAC(hash, []byte("aws4_request"))

	s := Hex(getHMAC(hash, []byte(rec.sts)))
	fmt.Printf("Signature:\n%s\n\n", s)
	return s
}

func genAuthorization(cr CanonicalRequest, sts StringToSign, sg Signature, accessKeyId string) string {
	s := fmt.Sprintf("%s %s HTTP/1.1\n", cr.HTTPMethod, cr.CanonicalURI)

	SignedHeaders := []string{}
	for k := range cr.CanonicalHeaders {
		SignedHeaders = append(SignedHeaders, k)
	}
	sort.Strings(SignedHeaders)
	for _, k := range SignedHeaders {
		v := cr.CanonicalHeaders[k]
		if !strings.HasPrefix(k, "x-") {
			k = strings.Title(k)
		}
		s += fmt.Sprintf("%s: %s\n", k, v)
	}

	s1 := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/s3/aws4_request", accessKeyId, sts.t.Format("20060102"), sts.region)

	return fmt.Sprintf("Authorization: %s,SignedHeaders=%s,Signature=%s", s1, strings.Join(SignedHeaders, ";"), sg.String())
}

func genRequest(cr CanonicalRequest, sts StringToSign, sg Signature, accessKeyId string) string {
	s := fmt.Sprintf("%s %s HTTP/1.1\n", cr.HTTPMethod, cr.CanonicalURI)

	SignedHeaders := []string{}
	for k := range cr.CanonicalHeaders {
		SignedHeaders = append(SignedHeaders, k)
	}
	sort.Strings(SignedHeaders)
	for _, k := range SignedHeaders {
		v := cr.CanonicalHeaders[k]
		if !strings.HasPrefix(k, "x-") {
			k = strings.Title(k)
		}
		s += fmt.Sprintf("%s: %s\n", k, v)
	}

	s1 := fmt.Sprintf("AWS4-HMAC-SHA256 Credential=%s/%s/%s/s3/aws4_request", accessKeyId, sts.t.Format("20060102"), sts.region)

	s += fmt.Sprintf("Authorization: %s,SignedHeaders=%s,Signature=%s", s1, strings.Join(SignedHeaders, ";"), sg.String())

	fmt.Printf("Request:\n%s\n\n", s)

	return s
}

func AwsExample() {

	HashPayload := Hex(sha256Hash(""))

	// t := time.Now().UTC()

	t := time.Date(2013, 5, 24, 0, 0, 0, 0, time.UTC)

	cr := CanonicalRequest{
		HTTPMethod:           "GET",
		CanonicalURI:         "/test.txt",
		CanonicalQueryString: "",
		CanonicalHeaders: map[string]string{
			"host":                 "examplebucket.s3.amazonaws.com",
			"range":                "bytes=0-9",
			"x-amz-content-sha256": HashPayload,
			"x-amz-date":           t.Format("20060102T150405Z"),
		},
		HashPayload: HashPayload,
	}

	sts := StringToSign{
		t:      t,
		region: "us-east-1",
		cr:     cr.String(),
	}

	sig := Signature{
		t:               t,
		SecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		region:          "us-east-1",
		sts:             sts.String(),
	}

	genRequest(cr, sts, sig, "AKIAIOSFODNN7EXAMPLE")

}

func Example1() {

	HashPayload := Hex(sha256Hash(""))

	// t := time.Now().UTC()

	t := time.Date(2021, 6, 25, 9, 39, 48, 0, time.UTC)

	cr := CanonicalRequest{
		HTTPMethod:           "GET",
		CanonicalURI:         "/thumb/images/vod/cover/VOD-01F4EKV7MBD91JZ9A7GQDNJ72M-P.png",
		CanonicalQueryString: "",
		CanonicalHeaders: map[string]string{
			"host":                 "noovo-files.s3.eu-west-3.amazonaws.com",
			"x-amz-content-sha256": HashPayload,
			"x-amz-date":           t.Format("20060102T150405Z"),
			"x-amz-security-token": "FwoGZXIvYXdzEGMaDHlmmxU4Nl2mP1Y+OCLEAfUexRp6z0/fZGMDCAW62ADUAiIQbHsMW6/4oRWPEh5xBzjoZcwCPSRMmMWNKZBGA54qZyqKEazykQPhs1XwTabkYepoJkyp9mijnFUa1DsGSblj6baJc49lHK7aXkldf5pUJLu/PYr3AzDhrOn/GvAAvcnBymB4lcR71a+gbH5FNG4s77MRNxLuM1I4/tK+Ureh6BNcXSOJWUzBIdB0oMx2+Hy7XDOx7widIZpUmzcsIgk/LWNoA5lgfSnUY5xPmFQJgeEo+MfWhgYyLZvDFAoV5fbO7HeF0wHwGEHAJ8StP9vRSIHS9M+nTFZb2hulQ5xWBG4RybAKFw==",
		},
		HashPayload: HashPayload,
	}

	sts := StringToSign{
		t:      t,
		region: "eu-west-3",
		cr:     cr.String(),
	}

	sig := Signature{
		t:               t,
		SecretAccessKey: "RnDrIWwS5l8F+Pqcbvfy9LCZ1DdHVUps2Qsm+CwO",
		region:          "eu-west-3",
		sts:             sts.String(),
	}

	want := "Authorization: AWS4-HMAC-SHA256 Credential=ASIARLGJNELUGR53OOU2/20210625/eu-west-3/s3/aws4_request,SignedHeaders=host;x-amz-content-sha256;x-amz-date;x-amz-security-token,Signature=903f4158520724fd2d1264f5f794f14905c790cf38fdb9b6cb6da455f4b19046"
	if got := genAuthorization(cr, sts, sig, "ASIARLGJNELUGR53OOU2"); got != want {
		logger.Debug(want)
		logger.Debug(got)
	}

	genRequest(cr, sts, sig, "ASIARLGJNELUGR53OOU2")
}

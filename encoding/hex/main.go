package main

import (
	"encoding/hex"
	"fmt"
)

func main() {
	src := []byte("hello")
	maxEnLen := hex.EncodedLen(len(src)) // 最大编码长度
	dst1 := make([]byte, maxEnLen)
	n := hex.Encode(dst1, src)
	dst2 := hex.EncodeToString(src)
	fmt.Println(src)
	fmt.Println(dst1)
	fmt.Printf("%s", hex.Dump(src))
	fmt.Println("编码后的结果:", string(dst1[:n]))
	fmt.Println("编码后的结果:", dst2)

	fmt.Println(hexString("hello"))
}

func hexString(src string) string {
	by := []byte(src)
	maxEnLen := hex.EncodedLen(len(by)) // 最大编码长度
	dst := make([]byte, maxEnLen)
	n := hex.Encode(dst, by)
	return hex.EncodeToString(by)[:n]
}

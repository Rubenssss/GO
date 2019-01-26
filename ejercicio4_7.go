package main

import (
	"unicode/utf8"
	"fmt"
)

func reverseBytes(a []byte)[]byte {
	if utf8.RuneCount(a) == 1 {
		return a
	}
	_,s := utf8.DecodeRune(a)
	return append(reverseBytes(a[s:]), a[:s]...)
}

func main() {
	a := []byte("abb")
	fmt.Println(string(reverseBytes(a)))
}


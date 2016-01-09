package main

import (
	"io"
	"os"
	"strings"
)

var dict1 = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
var dict2 = []byte("NOPQRSTUVWXYZABCDEFGHIJKLMnopqrstuvwxyzabcdefghijklm")

type rot13Reader struct {
	r io.Reader
}

func (rr rot13Reader) Read(b []byte) (int, error) {
	n, err := rr.r.Read(b)
	for i, v := range b {
		for j, d := range dict1 {
			if d == v {
				b[i] = dict2[j]
			}
		}
	}
	return n, err
}

func main() {
	s := strings.NewReader("Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	io.Copy(os.Stdout, &r)
}

package main

import (
	"bytes"
	"strconv"
)

// 1366783e726174667d22537e727a7f3d6166

func cf(a string) (s string) {
	var e bytes.Buffer
	r, _ := strconv.ParseInt(a[0:2], 16, 0)
	for n := 4; n < len(a)+2; n += 2 {
		i, _ := strconv.ParseInt(a[n-2:n], 16, 0)
		e.WriteString(string(rune(i ^ r)))
	}
	return e.String()
}

func decryptEmail(encodedString string) string {
	return cf(encodedString)
}

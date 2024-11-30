package main

import (
	"fmt"
	"testing"
)

func TestEmailDecript(t *testing.T) {
	criptoStr := "1366783e726174667d22537e727a7f3d6166"
	decrypt := decryptEmail(criptoStr)

	fmt.Println(decrypt)
	fmt.Println(cf(criptoStr))

}

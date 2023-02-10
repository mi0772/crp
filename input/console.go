package input

import (
	"bytes"
	"fmt"
	"syscall"

	"golang.org/x/term"
)

func ReadPassword() (password string, err error) {
	fmt.Printf("specify key:")
	var p1 = make([]byte, 32)
	p1, err = term.ReadPassword(syscall.Stdin)
	if err != nil {
		return
	}
	password = string(p1)
	return
}

func ReadPasswordAndConfirm() (password string, err error) {
	fmt.Printf("specify key:")
	var p1 = make([]byte, 32)
	var p2 = make([]byte, 32)
	p1, err = term.ReadPassword(syscall.Stdin)
	if err != nil {
		return
	}

	fmt.Printf("\nenter again:")
	p2, err = term.ReadPassword(syscall.Stdin)
	if err != nil {
		return
	}

	if bytes.Compare(p1, p2) != 0 {
		fmt.Println("\npasswords do not match")
		return
	}
	fmt.Println()
	password = string(p1)
	return
}

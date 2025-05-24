package io

import (
	"bufio"
	"fmt"
	"io"
)

func ReadUserInput(r io.Reader) string {
	reader := bufio.NewReader(r)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("broski shit went sideways :\\ - %v\n", err)
	}
	return input
}

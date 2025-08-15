package util

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

var inputReader *bufio.Reader = bufio.NewReader(os.Stdin)

func SetInputReader(r *os.File) {
	inputReader = bufio.NewReader(r)
}

func GetInput(prompt string) (string, error) {
	fmt.Printf("%v", prompt)
	input, err := inputReader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

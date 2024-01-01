package cli

import (
	"bufio"
	"fmt"
	"os"
)

func WriteSystemMessage(msg string) {
	fmt.Printf("🤖 %s\n\n", msg)
	return
}

func GetUserInput() (string, error) {
	fmt.Print("> ")
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Print("\n")
	input := s[:len(s)-1]
	if input == "" {
		return "", fmt.Errorf("no user input")
	}
	return input, nil
}

func ExitWithError(err error) {
	fmt.Println("Error: ", err)
	os.Exit(1)
}

package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Exit() {
	fmt.Println("bye!")
	os.Exit(0)
}

func ExitWithError(err error) {
	fmt.Println("Error: ", err)
	os.Exit(1)
}

var Yellow = color.New(color.FgYellow).SprintFunc()
var Red = color.New(color.FgRed).SprintFunc()

func WriteAssistantMessage(msg string) {
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
	} else if input == "exit" {
		Exit()
	}
	return input, nil
}

// Package cli provides utilities for working with the CLI
package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/fatih/color"
)

// Exit exits the CLI with a success code
func Exit() {
	fmt.Println("bye!")
	os.Exit(0)
}

// ExitWithError exits the CLI with an error code
func ExitWithError(err error) {
	fmt.Println("Error: ", err)
	os.Exit(1)
}

// Yellow is a function which makes a string yellow, ie. Yellow("warning")
var Yellow = color.New(color.FgYellow).SprintFunc()

// Red is a function which makes a string red, ie. Red("danger")
var Red = color.New(color.FgRed).SprintFunc()

// WriteAssistantMessage prints a formatted message from the assistant
func WriteAssistantMessage(msg string) {
	fmt.Printf("🤖 %s\n\n", msg)
	return
}

// GetUserInput prompts the user to enter input. If the user types "exit", the
// CLI will exit with a success code. If the user doesn't type anything, this
// function returns an error.
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

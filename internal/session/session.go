package session

import (
	"bufio"
	"fmt"
	"os"
)

type Role string

const (
	RoleSystem Role = "system"
	RoleUser   Role = "user"
	RoleTool   Role = "tool"
)

type Message struct {
	Role    Role
	Content string
}

type Session struct {
	Messages []*Message
}

func promptUser() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("> ")
	s, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	fmt.Println()
	return s[:len(s)-1], nil
}

func (s *Session) Advance(msg *Message) error {
	s.Messages = append(s.Messages, msg)
	switch msg.Role {
	case RoleSystem:
		fmt.Printf("🤖 %s\n\n", msg.Content)

		u, err := promptUser()
		if err != nil {
			return err
		}
		if u == "" {
			return fmt.Errorf("no user input")
		}

		err = s.Advance(&Message{
			Role:    RoleUser,
			Content: u,
		})
		if err != nil {
			return err
		}
	case RoleUser:
		fmt.Print("Querying LLM...\n\n")
		// TODO: Get a chat completion from the LLM
		fmt.Print("I suggest you run: `ls`\n\n")
		fmt.Print("Want to run it? (y)es / (n)o\n\n")
		// TODO: Add a user prompt for yes/no
		fmt.Print("> y\n\n")
		err := s.Advance(&Message{
			Role:    RoleTool,
			Content: "ls",
		})
		if err != nil {
			return err
		}
	case RoleTool:
		fmt.Print("Running `ls`...\n\n")
		fmt.Print("Output:\n\nDesktop\nDownloads\nDocuments\n\n")
		// TODO: Call Advance with RoleSystem and see if user wants anything else
	}
	return nil
}

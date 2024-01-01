package session

import (
	"fmt"

	"github.com/dgurns/clix/internal/cli"
	"github.com/dgurns/clix/internal/llm"
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
	LLM      *llm.LLM
	Messages []*Message
}

func New(llm *llm.LLM) (*Session, error) {
	if llm == nil {
		return nil, fmt.Errorf("no LLM passed to session")
	}
	return &Session{
		LLM:      llm,
		Messages: []*Message{},
	}, nil
}

func (s *Session) Advance(msg *Message) error {
	s.Messages = append(s.Messages, msg)
	switch msg.Role {
	case RoleSystem:
		cli.WriteSystemMessage(msg.Content)

		u, err := cli.GetUserInput()
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

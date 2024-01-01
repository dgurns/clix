package session

import (
	"fmt"

	"github.com/dgurns/clix/internal/cli"
	"github.com/dgurns/clix/internal/llm"
)

type Session struct {
	LLM      llm.LLM
	Messages []*llm.Message
}

func New(l llm.LLM) (*Session, error) {
	if l == nil {
		return nil, fmt.Errorf("no LLM passed to session")
	}
	return &Session{
		LLM:      l,
		Messages: []*llm.Message{},
	}, nil
}

func (s *Session) Advance(msg *llm.Message) error {
	s.Messages = append(s.Messages, msg)

	switch msg.Role {
	case llm.RoleSystem:
		cli.WriteSystemMessage(msg.Content)

		u, err := cli.GetUserInput()
		if err != nil {
			return err
		}
		if u == "" {
			return fmt.Errorf("no user input")
		}

		err = s.Advance(&llm.Message{
			Role:    llm.RoleUser,
			Content: u,
		})
		if err != nil {
			return err
		}
	case llm.RoleUser:
		cli.WriteSystemMessage("Querying LLM...")

		c, err := s.LLM.CreateChatCompletion(s.Messages)
		if err != nil {
			return err
		}

		cli.WriteSystemMessage(fmt.Sprintf("I suggest you run: %s", c.Content))

		cli.WriteSystemMessage("Want to run it? (y)es / (n)o")

		u, err := cli.GetUserInput()
		if err != nil {
			return err
		}
		if u != "y" {
			s.Advance(&llm.Message{
				Role:    llm.RoleSystem,
				Content: "Okay, what would you like to do instead?",
			})
		}

		err = s.Advance(&llm.Message{
			Role:    llm.RoleTool,
			Content: c.Content,
		})
		if err != nil {
			return err
		}
	case llm.RoleTool:
		cli.WriteSystemMessage(fmt.Sprintf("Running command: %s", msg.Content))

		// TODO: run the commmand

		fmt.Print("Output:\n\nDesktop\nDownloads\nDocuments\n\n")
		// TODO: Call Advance with RoleSystem and see if user wants anything else
	}
	return nil
}

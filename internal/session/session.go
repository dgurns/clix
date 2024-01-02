package session

import (
	"encoding/json"
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
		if err := s.Advance(&llm.Message{
			Role: llm.RoleAssistant,
			Content: `Welcome to Clix! I can help you run commands on your computer. What would you like to do? 
For example, "Reorganize my desktop" or "Initialize a new git repository"`,
		}); err != nil {
			return err
		}

	case llm.RoleAssistant:
		if len(msg.ToolCalls) > 0 {
			args := map[string]string{}
			err := json.Unmarshal([]byte(msg.ToolCalls[0].Function.Arguments), &args)
			if err != nil {
				return err
			}

			cmd, ok := args["command"]
			if !ok {
				return fmt.Errorf("no command provided by tool")
			}
			rationale, ok := args["rationale"]
			if !ok {
				return fmt.Errorf("no rationale provided by tool")
			}

			cli.WriteAssistantMessage(fmt.Sprintf(
				"%s\n\nI suggest you run: %s\n\nWant to run it? (y)es / (n)o",
				rationale,
				cmd,
			))

			u, err := cli.GetUserInput()
			if err != nil {
				return err
			}
			if u != "y" {
				if err := s.Advance(&llm.Message{
					Role:       llm.RoleTool,
					Name:       llm.FunctionNameRunTerminalCommand,
					ToolCallID: msg.ToolCalls[0].ID,
					Content:    "User chose not to call this function",
				}); err != nil {
					return err
				}
			}

			cli.WriteAssistantMessage(fmt.Sprintf("Running command: %s", cmd))

			// TODO: run the commmand

			mockOutput := "Desktop\nDownloads\nDocuments"

			fmt.Printf("Output:\n\n%s\n\n", mockOutput)

			if err = s.Advance(&llm.Message{
				Role:       llm.RoleTool,
				Name:       llm.FunctionNameRunTerminalCommand,
				Content:    mockOutput,
				ToolCallID: msg.ToolCalls[0].ID,
			}); err != nil {
				return err
			}
		}

		// the assistant isn't trying to run any tool calls, so handle it like
		// a normal assistant message

		cli.WriteAssistantMessage(msg.Content)

		u, err := cli.GetUserInput()
		if err != nil {
			return err
		}
		if u == "" {
			return fmt.Errorf("no user input")
		}

		if err = s.Advance(&llm.Message{
			Role:    llm.RoleUser,
			Content: u,
		}); err != nil {
			return err
		}

	case llm.RoleUser:
		cli.WriteAssistantMessage("Querying LLM...")

		c, err := s.LLM.CreateChatCompletion(s.Messages)
		if err != nil {
			return err
		}
		if err = s.Advance(c); err != nil {
			return err
		}

	case llm.RoleTool:
		cli.WriteAssistantMessage("Sending output to LLM...")

		c, err := s.LLM.CreateChatCompletion(s.Messages)
		if err != nil {
			return err
		}
		if err = s.Advance(c); err != nil {
			return err
		}
	}

	return nil
}

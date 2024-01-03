// Package main provides the CLI entrypoint to the app
package main

import (
	"github.com/dgurns/clix/internal/cli"
	"github.com/dgurns/clix/internal/config"
	"github.com/dgurns/clix/internal/llm"
	"github.com/dgurns/clix/internal/session"
)

func main() {
	c, err := config.Init()
	if err != nil {
		cli.ExitWithError(err)
	}

	l, err := llm.NewOpenAi(c.OpenAiAPIKey)
	if err != nil {
		cli.ExitWithError(err)
	}

	s, err := session.New(l)
	if err != nil {
		cli.ExitWithError(err)
	}

	err = s.Advance(&llm.Message{
		Role:    llm.RoleSystem,
		Content: llm.SystemPrompt,
	})
	if err != nil {
		cli.ExitWithError(err)
	}
}

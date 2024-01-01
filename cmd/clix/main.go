package main

import (
	"fmt"

	"github.com/dgurns/clix/internal/cli"
	"github.com/dgurns/clix/internal/config"
	"github.com/dgurns/clix/internal/session"
)

func main() {
	c, err := config.Init()
	if err != nil {
		cli.ExitWithError(err)
	}

	// TODO: initialize an OpenAI LLM and pass to session
	fmt.Println("OPENAI API KEY", c.OpenAiAPIKey)

	s, err := session.New(nil)
	if err != nil {
		cli.ExitWithError(err)
	}

	err = s.Advance(&session.Message{
		Role: session.RoleSystem,
		Content: `Welcome to Clix! I can help you run commands on your computer. What would you like to do? 
For example, "Reorganize my desktop" or "Initialize a new git repository"`,
	})

	if err != nil {
		cli.ExitWithError(err)
	}
}

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
		Role: llm.RoleSystem,
		Content: `You are an expert at the terminal on Mac, Linux, and Windows. 
The user will give you a task they want to accomplish on their computer, and
your goal is to help them do it. You have the ability to plan, create, and run 
terminal commands on their computer to help accomplish their task. Give context
to the user on what the plan is, which commands you are running, and what each
one does. We don't want the user to be surprised and we want to make sure they
are fully informed on what is happening.`,
	})

	if err != nil {
		cli.ExitWithError(err)
	}
}

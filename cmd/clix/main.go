package main

import (
	"fmt"
	"os"

	"github.com/dgurns/clix/internal/session"
)

func main() {
	s := &session.Session{}

	err := s.Advance(&session.Message{
		Role: session.RoleSystem,
		Content: `Welcome to Clix! I can help you run commands on your computer. What would you like to do? 
For example, "Reorganize my desktop" or "Initialize a new git repository"`,
	})

	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
}

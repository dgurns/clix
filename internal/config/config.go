package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgurns/clix/internal/cli"
)

type Config struct {
	OpenAiAPIKey string
}

var openAiAPIKeyName = "OPENAI_API_KEY"

// Init returns config for the app like API keys. Config is stored on the user's
// system at ~/.clix/.env and will be created if it doesn't exist.
func Init() (*Config, error) {
	configDir := os.Getenv("HOME") + "/.clix"
	envFile := configDir + "/.env"

	if _, err := os.Stat(configDir); err != nil {
		if err = os.Mkdir(configDir, 0755); err != nil {
			return nil, err
		}
	}
	if _, err := os.Stat(envFile); err != nil {
		if _, err = os.Create(envFile); err != nil {
			return nil, err
		}
	}

	raw, err := os.ReadFile(envFile)
	if err != nil {
		return nil, err
	}

	envs := strings.Split(string(raw), "\n")
	for _, e := range envs {
		pre := fmt.Sprintf("%s=", openAiAPIKeyName)
		if strings.HasPrefix(e, pre) {
			return &Config{
				OpenAiAPIKey: strings.TrimPrefix(e, pre),
			}, nil
		}
	}

	cli.WriteAssistantMessage("Please enter your OpenAI API key")

	openAiAPIKey, err := cli.GetUserInput()
	if err != nil {
		return nil, err
	}

	cli.WriteAssistantMessage("Saving OpenAI API key locally to ~/.clix/.env")

	f, err := os.OpenFile(envFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	e := fmt.Sprintf("%s=%s\n", openAiAPIKeyName, openAiAPIKey)
	if _, err := f.WriteString(e); err != nil {
		return nil, err
	}

	cli.WriteAssistantMessage("OpenAI API key saved")

	return &Config{
		OpenAiAPIKey: openAiAPIKey,
	}, nil
}

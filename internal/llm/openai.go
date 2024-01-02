package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

type OpenAiLLM struct {
	client *openai.Client
}

var _ LLM = OpenAiLLM{}

func NewOpenAi(apiKey string) (*OpenAiLLM, error) {
	c := openai.NewClient(apiKey)
	return &OpenAiLLM{
		client: c,
	}, nil
}

func formatForOpenAI(msgs []*Message) []openai.ChatCompletionMessage {
	var formatted []openai.ChatCompletionMessage

	for _, m := range msgs {
		role := openai.ChatMessageRoleUser
		if m.Role == RoleSystem {
			role = openai.ChatMessageRoleSystem
		} else if m.Role == RoleAssistant {
			role = openai.ChatMessageRoleAssistant
		} else if m.Role == RoleTool {
			role = openai.ChatMessageRoleTool
		}
		tc := []openai.ToolCall{}
		for _, t := range m.ToolCalls {
			functionName := string(t.Function.Name)
			tc = append(tc, openai.ToolCall{
				ID:   t.ID,
				Type: openai.ToolTypeFunction,
				Function: openai.FunctionCall{
					Name:      functionName,
					Arguments: t.Function.Arguments,
				},
			})
		}
		formatted = append(formatted, openai.ChatCompletionMessage{
			Role:       role,
			Content:    m.Content,
			ToolCalls:  tc,
			Name:       string(m.Name),
			ToolCallID: m.ToolCallID,
		})
	}

	return formatted
}

func formatFromOpenAI(msg openai.ChatCompletionMessage) *Message {
	role := RoleUser
	if msg.Role == openai.ChatMessageRoleSystem {
		role = RoleSystem
	} else if msg.Role == openai.ChatMessageRoleAssistant {
		role = RoleAssistant
	} else if msg.Role == openai.ChatMessageRoleTool {
		role = RoleTool
	}

	tc := []ToolCall{}
	for _, t := range msg.ToolCalls {
		tc = append(tc, ToolCall{
			ID:   t.ID,
			Type: ToolTypeFunction,
			Function: FunctionCall{
				Name:      FunctionName(t.Function.Name),
				Arguments: t.Function.Arguments,
			},
		})
	}

	return &Message{
		Role:       role,
		Content:    msg.Content,
		ToolCalls:  tc,
		Name:       FunctionName(msg.Name),
		ToolCallID: msg.ToolCallID,
	}
}

var tools = []openai.Tool{
	{
		Type: openai.ToolTypeFunction,
		Function: openai.FunctionDefinition{
			Name:        string(FunctionNameRunTerminalCommand),
			Description: "Run a terminal command on the user's computer",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"command": {
						Type:        jsonschema.String,
						Description: "The terminal command to run, e.g. 'cat ./dir/somefile.txt'",
					},
					"rationale": {
						Type:        jsonschema.String,
						Description: "Explain to the user why we're running this command, e.g. 'In order to see what is in the file, we will run the `cat` command to print its contents to the terminal'",
					},
				},
				Required: []string{"command", "rationale"},
			},
		},
	},
}

func (l OpenAiLLM) CreateChatCompletion(msgs []*Message) (*Message, error) {
	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    openai.GPT3Dot5Turbo,
			Tools:    tools,
			Messages: formatForOpenAI(msgs),
		},
	)
	if err != nil {
		fmt.Printf("OpenAI CreateChatCompletion error: %v\n", err)
		return nil, err
	}
	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response choices returned by OpenAI")
	}

	msg := formatFromOpenAI(resp.Choices[0].Message)
	return msg, nil
}

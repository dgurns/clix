package llm

import (
	"context"
	"errors"
	"fmt"
	"io"

	openai "github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

// OpenAiLLM is an LLM that uses the OpenAI API
type OpenAiLLM struct {
	client *openai.Client
}

var _ LLM = OpenAiLLM{}

// GPT4o is an alias to the latest GPT-4o model ID
var GPT4o = openai.GPT4o20240513

// NewOpenAi initializes a new OpenAI LLM client
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
		Function: &openai.FunctionDefinition{
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

// CreateChatCompletion calls the OpenAI API and returns the chat completion
func (l OpenAiLLM) CreateChatCompletion(msgs []*Message) (*Message, error) {
	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    GPT4o,
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

// CreateChatCompletionWithStreaming calls the OpenAI API and prints the
// completion chunks as they come in. It returns the final completion as a
// single message.
func (l OpenAiLLM) CreateChatCompletionWithStreaming(msgs []*Message) (*Message, error) {
	stream, err := l.client.CreateChatCompletionStream(
		context.Background(),
		openai.ChatCompletionRequest{
			Model:    GPT4o,
			Tools:    tools,
			Messages: formatForOpenAI(msgs),
			Stream:   true,
		},
	)
	if err != nil {
		fmt.Printf("OpenAI CreateChatCompletionWithStreaming error: %v\n", err)
		return nil, err
	}
	defer stream.Close()

	msgText := ""
	toolID := ""
	toolName := ""
	toolArgs := ""

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Print("\n")
			break
		}
		if err != nil {
			return nil, err
		}
		if response.Choices[0].FinishReason == openai.FinishReasonStop {
			fmt.Print("\n")
			break
		}
		delta := response.Choices[0].Delta
		if delta.ToolCalls != nil {
			// only use the first tool call to keep things simple
			tc := delta.ToolCalls[0]
			if tc.ID != "" {
				toolID = tc.ID
			}
			if tc.Function.Name != "" {
				toolName = string(tc.Function.Name)
			}
			if tc.Function.Arguments != "" {
				toolArgs += tc.Function.Arguments
			}
		}
		fmt.Printf(delta.Content)
		msgText += delta.Content
	}

	if toolID != "" {
		return &Message{
			Role: RoleAssistant,
			ToolCalls: []ToolCall{
				{
					ID:       toolID,
					Type:     ToolTypeFunction,
					Function: FunctionCall{Name: FunctionName(toolName), Arguments: toolArgs},
				},
			},
		}, nil
	}

	return &Message{
		Role:    RoleAssistant,
		Content: msgText,
	}, nil
}

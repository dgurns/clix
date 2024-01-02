package llm

import (
	"context"
	"fmt"

	openai "github.com/sashabaranov/go-openai"
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

func (l OpenAiLLM) CreateChatCompletion(msgs []*Message) (*Message, error) {
	resp, err := l.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Hello!",
				},
			},
		},
	)

	if err != nil {
		fmt.Printf("ChatCompletion error: %v\n", err)
		return nil, err
	}

	fmt.Println("MESSAGE", resp.Choices[0].Message.Content)

	return &Message{
		Role:    RoleTool,
		Content: "ls",
	}, nil
}

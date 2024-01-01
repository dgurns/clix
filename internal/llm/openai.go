package llm

type OpenAiLLM struct {
	OpenAiAPIKey string
}

var _ LLM = OpenAiLLM{}

func NewOpenAi(apiKey string) (*OpenAiLLM, error) {
	return &OpenAiLLM{
		OpenAiAPIKey: apiKey,
	}, nil
}

func (llm OpenAiLLM) CreateChatCompletion(msgs []*Message) (*Message, error) {
	// TODO: actually call the OpenAI API
	return &Message{
		Role:    RoleTool,
		Content: "ls",
	}, nil
}

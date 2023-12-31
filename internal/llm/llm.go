package llm

type LLM interface {
	CreateChatCompletion(msgs []string) (string, error)
}

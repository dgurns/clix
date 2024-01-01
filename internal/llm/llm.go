package llm

type Role string

const (
	RoleSystem Role = "system"
	RoleUser   Role = "user"
	RoleTool   Role = "tool"
)

type Message struct {
	Role    Role
	Content string
}

type LLM interface {
	CreateChatCompletion(msgs []*Message) (*Message, error)
}

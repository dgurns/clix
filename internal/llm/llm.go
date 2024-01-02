package llm

type Role string

const (
	RoleSystem    Role = "system"
	RoleAssistant Role = "assistant"
	RoleUser      Role = "user"
	RoleTool      Role = "tool"
)

type ToolType string

const (
	ToolTypeFunction ToolType = "function"
)

type FunctionName string

const (
	FunctionNameRunTerminalCommand FunctionName = "run_terminal_command"
)

type FunctionCall struct {
	Name      FunctionName `json:"name,omitempty"`
	Arguments string       `json:"arguments,omitempty"`
}

type ToolCall struct {
	ID       string
	Type     ToolType
	Function FunctionCall
}

type Message struct {
	Role Role
	// the message content, which could be the function response if a
	// function was called
	Content   string
	ToolCalls []ToolCall
	// if message is a tool response, include the function name that was called
	// and the tool call id
	Name       FunctionName
	ToolCallID string
}

type LLM interface {
	CreateChatCompletion(msgs []*Message) (*Message, error)
}

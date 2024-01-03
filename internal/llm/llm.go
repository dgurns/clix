// Package llm is used to interface with LLMs like GPT3.5, Claude, etc.
package llm

// SystemPrompt is passed as context to the LLM
var SystemPrompt = `You are an expert at the terminal on Mac, Linux, and Windows. 
The user will give you a task they want to accomplish on their computer, and
your goal is to help them do it. You have the ability to plan, create, and run 
terminal commands on their computer to help accomplish their task. When it 
makes sense, err on the side of writing commands instead of giving written 
explanations. You can make a plan that involves running multiple commands in 
a row. Give context to the user on what the plan is, which commands you are 
running, and what each one does. We don't want the user to be surprised and we 
want to make sure they are fully informed on what is happening.`

// Role is the role of the message
type Role string

const (
	// RoleSystem is used to set the system prompt
	RoleSystem Role = "system"
	// RoleAssistant is the assistant role
	RoleAssistant Role = "assistant"
	// RoleUser is the user role
	RoleUser Role = "user"
	// RoleTool is used to send back tool responses to the LLM
	RoleTool Role = "tool"
)

// ToolType is the type of tool being used
type ToolType string

const (
	// ToolTypeFunction is a function tool
	ToolTypeFunction ToolType = "function"
)

// FunctionName is the name of the function, ie. "run_terminal_command"
type FunctionName string

const (
	// FunctionNameRunTerminalCommand is a function name for running a terminal
	// command
	FunctionNameRunTerminalCommand FunctionName = "run_terminal_command"
)

// FunctionCall is a function that the assistant decides to call
type FunctionCall struct {
	Name      FunctionName `json:"name,omitempty"`
	Arguments string       `json:"arguments,omitempty"`
}

// ToolCall is a tool the assistant decides to call, which might be a function
type ToolCall struct {
	ID       string
	Type     ToolType
	Function FunctionCall
}

// Message is a message sent to and from the LLM
type Message struct {
	// Role should be tool when responding to a tool call
	Role Role
	// Content can be text or a tool response
	Content string
	// ToolCalls are the tool calls that the assistant decides to call
	ToolCalls []ToolCall
	// Name is included when the message is a tool response
	Name FunctionName
	// ToolCallID must be included when the message is a tool response
	ToolCallID string
}

// LLM is an interface layer on top of LLMs like GPT3.5, Claude, etc.
type LLM interface {
	CreateChatCompletion(msgs []*Message) (*Message, error)
}

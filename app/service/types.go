package service

type CompletionRequest struct {
	Prompt string `json:"prompt"`
}

type CompletionResponse struct {
	Response string `json:"response"`
	Provider string `json:"provider"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ChatGptHandlerType string = "chat-gpt"
	AmazonBedrockType  string = "aws-bedrock"
)

type GenAIHandler interface {
	Init(apiKey string)
	Type() string
	ChatCompletion(message string) (string, error)
}

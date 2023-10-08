package service

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ChatGptHandlerType string = "chat-gpt"
	AmazonBedrockType  string = "aws-bedrock"
)

type GenAIHandler interface {
	Type() string
	ChatCompletion(req *CompletionRequest) (string, error)
}

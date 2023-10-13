package service

type ErrorResponse struct {
	Error string `json:"error"`
}

const (
	ChatGptHandlerType string = "chat-gpt"
	AmazonBedrockType  string = "aws-bedrock"
)

const (
	Completions     string = "Completions"
	ImageGeneration string = "Image generation"
	TextToSpeech    string = "Text to speech"
	SpeechToText    string = "Speech to text"
)

type GenAIHandler interface {
	Type() string
	GetCapabilities() []string
	ChatCompletion(req *CompletionRequest) (string, error)
	GenerateImages(req *GenerateImagesRequest) (*GenerateImagesResponse, error)
}

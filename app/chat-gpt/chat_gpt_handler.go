package chatgpt

import (
	"context"
	"gen-ai-service/app/service"

	"github.com/sashabaranov/go-openai"
)

type ChatGptHandler struct {
	Client *openai.Client
}

func (c *ChatGptHandler) Init(apiKey string) {
	c.Client = openai.NewClient(apiKey)
}

func (c *ChatGptHandler) Type() string {
	return service.ChatGptHandlerType
}

func (c *ChatGptHandler) ChatCompletion(message string) (string, error) {
	resp, err := c.Client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: message,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

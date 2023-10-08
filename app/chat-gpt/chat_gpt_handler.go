package chatgpt

import (
	"context"
	"gen-ai-service/app/service"

	"github.com/sashabaranov/go-openai"
)

type ChatGptHandler struct {
	Client *openai.Client
}

func NewChatGptHandler(apiKey string) *ChatGptHandler {
	return &ChatGptHandler{
		Client: openai.NewClient(apiKey),
	}
}

func (c *ChatGptHandler) Type() string {
	return service.ChatGptHandlerType
}

func (c *ChatGptHandler) ChatCompletion(req *service.CompletionRequest) (string, error) {
	resp, err := c.Client.CreateChatCompletion(
		context.Background(),
		marshallChatGptCompletionRequest(req),
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *ChatGptHandler) GenerateImages(req *service.GenerateImagesRequest) (*service.GenerateImagesResponse, error) {
	gptResp, err := c.Client.CreateImage(
		context.Background(),
		marshallChatGptGenerateImageRequest(req),
	)

	if err != nil {
		return nil, err
	}

	images := make([]string, 0)
	for _, img := range gptResp.Data {
		switch *req.ResponseFormat {
		case service.ImageFormatBase64:
			images = append(images, img.B64JSON)
		default:
			images = append(images, img.URL)
		}
	}

	return &service.GenerateImagesResponse{
		Images:   images,
		Provider: c.Type(),
	}, nil
}

func marshallChatGptCompletionRequest(r *service.CompletionRequest) openai.ChatCompletionRequest {
	req := openai.ChatCompletionRequest{
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: r.Prompt,
			},
		},
	}
	if r.Model != nil {
		req.Model = *r.Model
	} else {
		req.Model = openai.GPT4
	}
	if r.MaxTokens != nil {
		req.MaxTokens = *r.MaxTokens
	}
	if r.Temperature != nil {
		req.Temperature = *r.Temperature
	}
	if r.TopP != nil {
		req.TopP = *r.TopP
	}
	return req
}

func marshallChatGptGenerateImageRequest(r *service.GenerateImagesRequest) openai.ImageRequest {
	req := openai.ImageRequest{
		Prompt: r.Prompt,
	}
	if r.Size != nil {
		req.Size = *r.Size
	}
	if r.ResponseFormat != nil {
		req.ResponseFormat = string(*r.ResponseFormat)
	}
	if r.Number != nil {
		req.N = *r.Number
	}
	return req
}

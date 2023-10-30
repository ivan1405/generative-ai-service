package chatgpt

import (
	"context"
	"gen-ai-service/app/service"

	"github.com/sashabaranov/go-openai"
)

type Handler struct {
	Client *openai.Client
}

func NewChatGptHandler(apiKey string) *Handler {
	return &Handler{
		Client: openai.NewClient(apiKey),
	}
}

func (h *Handler) Type() string {
	return service.ChatGptHandlerType
}

func (h *Handler) GetCapabilities() []string {
	return []string{
		service.Completions,
		service.ImageGeneration,
	}
}

func (h *Handler) ChatCompletion(req *service.CompletionRequest) (string, error) {
	resp, err := h.Client.CreateChatCompletion(
		context.Background(),
		marshallChatGptCompletionRequest(req),
	)

	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}

func (h *Handler) GenerateImages(req *service.GenerateImagesRequest) (*service.GenerateImagesResponse, error) {
	gptResp, err := h.Client.CreateImage(
		context.Background(),
		marshallChatGptGenerateImageRequest(req),
	)

	if err != nil {
		return nil, err
	}

	images := make([]string, 0)
	for _, img := range gptResp.Data {
		images = append(images, img.B64JSON)
	}

	return &service.GenerateImagesResponse{
		Image:    images[0],
		Provider: h.Type(),
	}, nil
}

func (h *Handler) TextToSpeech(req *service.TextToSpeechRequest) (*service.TextToSpeechResponse, error) {
	return nil, nil
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
	return openai.ImageRequest{
		Prompt:         r.Prompt,
		Size:           "512x512",
		N:              1,
		ResponseFormat: "b64_json",
	}
}

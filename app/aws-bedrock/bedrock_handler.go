package awsbedrock

import (
	"encoding/json"
	"gen-ai-service/app/service"
	"math/rand"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/bedrockruntime"
	"golang.org/x/exp/slog"
)

type Handler struct {
	Client *bedrockruntime.BedrockRuntime
}

type BedrockInferenceRequest struct {
	Prompt      string  `json:"prompt"`
	MaxTokens   int     `json:"max_tokens"`
	Temperature float32 `json:"temperature"`
	TopP        float32 `json:"p"`
}

type BedrockInferenceResponse struct {
	Generations []BedrockGeneration `json:"generations"`
}

type BedrockGeneration struct {
	Id   string `json:"id"`
	Text string `json:"text"`
}

type BedrockImageGenerationRequest struct {
	TextPrompts []TextPrompt `json:"text_prompts"`
	//  Determines how much the final image portrays the prompt.
	// Use a lower number to increase randomness in the generation.
	ConfigScale float32 `json:"cfg_scale"`
	// Generation step determines how many times the image is sampled.
	// More steps can result in a more accurate result.
	Steps int `json:"steps"`
	// The seed determines the initial noise setting.
	// Use the same seed and the same settings as a previous run to allow inference to create a similar image.
	// If you don't set this value, it is set as a random number.
	Seed int `json:"seed"`
}

type TextPrompt struct {
	Text string `json:"text"`
}

type BedrockImageGenerationResponse struct {
	Result    string                 `json:"result"`
	Artifacts []BedrockImageArtifact `json:"artifacts"`
}

type BedrockImageArtifact struct {
	Seed   int    `json:"seed"`
	Base64 string `json:"base64"`
}

func NewAWSBedrockHandler() *Handler {
	s := session.Must(session.NewSession())
	return &Handler{
		Client: bedrockruntime.New(s, aws.NewConfig().WithRegion("us-east-1")),
	}
}

func (h *Handler) Type() string {
	return service.AmazonBedrockType
}

func (h *Handler) GetCapabilities() []string {
	return []string{
		service.Completions,
		service.ImageGeneration,
	}
}

func (h *Handler) ChatCompletion(req *service.CompletionRequest) (string, error) {
	bedrockReq, err := marshallAWSBedrockCompletionRequest(req)
	if err != nil {
		slog.Error("Error parsing request")
		return "", err
	}

	r, err := h.Client.InvokeModel(bedrockReq)
	if err != nil {
		return "", err
	}

	var resp BedrockInferenceResponse
	if err := json.Unmarshal(r.Body, &resp); err != nil {
		slog.Error("Error unmarshalling response from Bedrock")
		return "", err
	}
	return strings.TrimSpace(resp.Generations[0].Text), nil
}

func (h *Handler) GenerateImages(req *service.GenerateImagesRequest) (*service.GenerateImagesResponse, error) {
	bedrockReq, err := marshallAWSBedrockImageGenerationRequest(req)
	if err != nil {
		slog.Error("Error parsing request")
		return nil, err
	}

	r, err := h.Client.InvokeModel(bedrockReq)
	if err != nil {
		return nil, err
	}

	var resp BedrockImageGenerationResponse
	if err := json.Unmarshal(r.Body, &resp); err != nil {
		slog.Error("Error unmarshalling response from Bedrock")
		return nil, err
	}
	return &service.GenerateImagesResponse{
		Image:    resp.Artifacts[0].Base64,
		Provider: h.Type(),
	}, nil
}

func (h *Handler) TextToSpeech(req *service.TextToSpeechRequest) (*service.TextToSpeechResponse, error) {
	return nil, nil
}

func marshallAWSBedrockCompletionRequest(r *service.CompletionRequest) (*bedrockruntime.InvokeModelInput, error) {
	req := &BedrockInferenceRequest{
		Prompt: r.Prompt,
	}
	if r.MaxTokens != nil {
		req.MaxTokens = *r.MaxTokens
	} else {
		req.MaxTokens = 300
	}
	if r.Temperature != nil {
		req.Temperature = *r.Temperature
	} else {
		req.Temperature = 0.1
	}
	if r.TopP != nil {
		req.TopP = *r.TopP
	} else {
		req.TopP = 0.2
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	modelId := "cohere.command-text-v14"
	if r.Model != nil {
		modelId = *r.Model
	}
	accept := "application/json"
	contentType := "application/json"
	return &bedrockruntime.InvokeModelInput{
		Body:        b,
		ModelId:     &modelId,
		Accept:      &accept,
		ContentType: &contentType,
	}, nil
}

func marshallAWSBedrockImageGenerationRequest(r *service.GenerateImagesRequest) (*bedrockruntime.InvokeModelInput, error) {
	// Config for inference request on their official documentation
	// https://docs.aws.amazon.com/bedrock/latest/userguide/model-parameters.html#model-parameters-diffusion
	req := &BedrockImageGenerationRequest{
		TextPrompts: []TextPrompt{
			{Text: r.Prompt},
		},
		ConfigScale: 15,
		Steps:       80,
		Seed:        rand.Intn(40) + 1,
	}

	b, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	modelId := "stability.stable-diffusion-xl-v0"
	accept := "application/json"
	contentType := "application/json"
	return &bedrockruntime.InvokeModelInput{
		Body:        b,
		ModelId:     &modelId,
		Accept:      &accept,
		ContentType: &contentType,
	}, nil
}

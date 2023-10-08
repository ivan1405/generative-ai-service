package awsbedrock

import (
	"encoding/json"
	"gen-ai-service/app/service"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/bedrockruntime"
	"golang.org/x/exp/slog"
)

type AWSBedrockHandler struct {
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

func NewAWSBedrockHandler() *AWSBedrockHandler {
	s := session.Must(session.NewSession())
	return &AWSBedrockHandler{
		Client: bedrockruntime.New(s, aws.NewConfig().WithRegion("us-east-1")),
	}
}

func (c *AWSBedrockHandler) Type() string {
	return service.AmazonBedrockType
}

func (c *AWSBedrockHandler) ChatCompletion(prompt string) (string, error) {
	req := &BedrockInferenceRequest{
		Prompt:      prompt,
		MaxTokens:   300,
		Temperature: 0.1,
		TopP:        0.9,
	}
	b, err := json.Marshal(req)
	if err != nil {
		slog.Error("Error parsing request")
		return "", err
	}

	modelId := "cohere.command-text-v14"
	accept := "application/json"
	contentType := "application/json"
	r, err := c.Client.InvokeModel(&bedrockruntime.InvokeModelInput{
		Body:        b,
		ModelId:     &modelId,
		Accept:      &accept,
		ContentType: &contentType,
	})
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

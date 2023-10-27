package elevenlabs

import (
	"bytes"
	"encoding/json"
	"gen-ai-service/app/service"
	"io"
	"net/http"
)

type ElevenLabsHandler struct {
	apiKey string
	client *http.Client
}

type ElevenLabsErrorResponse struct {
	Detail struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	} `json:"detail"`
}

// CustomError is a custom error type that implements the error interface.
type CustomError struct {
	Message string
}

// Error returns the error message. Implements the Error interface to be
// able to return it as a normal error in the response
func (e *CustomError) Error() string {
	return e.Message
}

func NewElevenLabsHandler(apiKey string) *ElevenLabsHandler {
	return &ElevenLabsHandler{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (c *ElevenLabsHandler) Type() string {
	return service.ElevenLabsType
}

func (c *ElevenLabsHandler) GetCapabilities() []string {
	return []string{
		service.TextToSpeech,
		service.SpeechToText,
	}
}

func (c *ElevenLabsHandler) ChatCompletion(req *service.CompletionRequest) (string, error) {
	return "not supported", nil
}

func (c *ElevenLabsHandler) GenerateImages(req *service.GenerateImagesRequest) (*service.GenerateImagesResponse, error) {
	return nil, nil
}

func (c *ElevenLabsHandler) TextToSpeech(req *service.TextToSpeechRequest) (*service.TextToSpeechResponse, error) {

	url := "https://api.elevenlabs.io/v1/text-to-speech/" + req.VoiceId

	data := map[string]interface{}{
		"text":     req.Prompt,
		"model_id": req.Model,
		"voice_settings": map[string]float64{
			"stability":        0.8,
			"similarity_boost": 0.5,
		},
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	elevenLabsReq, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	elevenLabsReq.Header.Set("Accept", "audio/mpeg")
	elevenLabsReq.Header.Set("Content-Type", "application/json")
	elevenLabsReq.Header.Set("xi-api-key", c.apiKey)

	resp, err := c.client.Do(elevenLabsReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Elevent labs does not return error, but a valid response
	// with status codes. So we need a custom error handling
	if resp.StatusCode != 200 {
		var errorResp ElevenLabsErrorResponse
		err := json.Unmarshal(respBody, &errorResp)
		if err != nil {
			return nil, &CustomError{
				Message: "Error unmarshalling response error",
			}
		}
		return nil, &CustomError{
			Message: errorResp.Detail.Message,
		}
	}

	return &service.TextToSpeechResponse{
		AudioData: respBody,
	}, nil
}

package elevenlabs

import (
	"bytes"
	"encoding/json"
	"gen-ai-service/app/service"
	"io"
	"net/http"
)

type Handler struct {
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

func NewElevenLabsHandler(apiKey string) *Handler {
	return &Handler{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

func (h *Handler) Type() string {
	return service.ElevenLabsType
}

func (h *Handler) GetCapabilities() []string {
	return []string{
		service.TextToSpeech,
		service.SpeechToText,
	}
}

func (h *Handler) ChatCompletion(req *service.CompletionRequest) (string, error) {
	return "not supported", nil
}

func (h *Handler) GenerateImages(req *service.GenerateImagesRequest) (*service.GenerateImagesResponse, error) {
	return nil, nil
}

func (h *Handler) TextToSpeech(req *service.TextToSpeechRequest) (*service.TextToSpeechResponse, error) {

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
	elevenLabsReq.Header.Set("xi-api-key", h.apiKey)

	resp, err := h.client.Do(elevenLabsReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// Eleven labs does not return error, but a valid response
	// with status codes. So we need a custom error handling
	if resp.StatusCode != http.StatusOK {
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

package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type TextToSpeechRequest struct {
	Prompt  string  `json:"prompt"`
	VoiceId string  `json:"voiceId"`
	Model   *string `json:"model,omitempty"`
}

type TextToSpeechResponse struct {
	AudioData []byte `json:"audio"`
}

func textToSpeech(c *gin.Context) {
	var req *TextToSpeechRequest
	if err := c.BindJSON(&req); err != nil {
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	handler := GetHandler(c.Request.Header.Get("Provider"))
	resp, err := handler.TextToSpeech(req)
	if err != nil {
		slog.Error("Error processing text to speech %s", err)
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	// Set headers for an MP3 file
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Content-Disposition", "attachment; filename=audio_file.mp3")

	// Write audio data to response
	c.Writer.Write(resp.AudioData)
}

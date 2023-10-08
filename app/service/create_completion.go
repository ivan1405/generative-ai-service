package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type CompletionRequest struct {
	Prompt      string   `json:"prompt"`
	Model       *string  `json:"model,omitempty"`
	MaxTokens   *int     `json:"max_tokens,omitempty"`
	Temperature *float32 `json:"temperature,omitempty"`
	TopP        *float32 `json:"top_p,omitempty"`
}

type CompletionResponse struct {
	Response string `json:"response"`
	Provider string `json:"provider"`
}

func createCompletion(c *gin.Context) {
	var compReq *CompletionRequest
	if err := c.BindJSON(&compReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	handler := getHandler(c)
	r, err := handler.ChatCompletion(compReq)
	if err != nil {
		slog.Error("Error creating chat completion %s", err)
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, &CompletionResponse{
		Response: r,
		Provider: handler.Type(),
	})
}

func getHandler(c *gin.Context) GenAIHandler {
	handlerType := c.Request.Header.Get("Provider")
	handler, exists := svc.Handlers[handlerType]
	// default handler chat-gpt
	if !exists {
		return svc.Handlers[ChatGptHandlerType]
	}
	return handler
}

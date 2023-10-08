package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

var svc *GenAIService

type GenAIService struct {
	Router   *gin.Engine
	Handlers map[string]GenAIHandler
}

func InitRouter() *GenAIService {
	svc = &GenAIService{
		Router: gin.Default(),
	}
	svc.Router.POST("/completion", createCompletion)
	return svc
}

func (genAIApi *GenAIService) ConfigureHandlers(handlers ...GenAIHandler) {
	handlersMap := make(map[string]GenAIHandler)
	for _, h := range handlers {
		handlersMap[h.Type()] = h
	}
	svc.Handlers = handlersMap
}

func (genAIApi *GenAIService) Start() {
	svc.Router.Run()
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
	r, err := handler.ChatCompletion(compReq.Prompt)
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

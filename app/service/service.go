package service

import (
	"github.com/gin-gonic/gin"
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
	svc.Router.GET("/ai-capabilities", getCapabilities)
	svc.Router.POST("/completion", createCompletion)
	svc.Router.POST("/images/generation", generateImages)
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

func GetHandler(handlerType string) GenAIHandler {
	handler, exists := svc.Handlers[handlerType]
	// default handler chat-gpt
	if !exists {
		return svc.Handlers[ChatGptHandlerType]
	}
	return handler
}

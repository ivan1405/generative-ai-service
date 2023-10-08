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

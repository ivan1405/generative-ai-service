package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func getCapabilities(c *gin.Context) {
	handlerTypes := []string{
		ChatGptHandlerType,
		AmazonBedrockType,
		ElevenLabsType,
	}

	var getCapabilitiesResponse = make(map[string][]string)

	for _, handlerType := range handlerTypes {
		handler := GetHandler(handlerType)
		resp := handler.GetCapabilities()
		getCapabilitiesResponse[handlerType] = resp
	}
	c.JSON(http.StatusOK, getCapabilitiesResponse)
}

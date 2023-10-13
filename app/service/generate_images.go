package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type GenerateImagesRequest struct {
	Prompt string `json:"prompt"`
}

type GenerateImagesResponse struct {
	Image    string `json:"image"`
	Provider string `json:"provider"`
}

func generateImages(c *gin.Context) {
	var genImgReq *GenerateImagesRequest
	if err := c.BindJSON(&genImgReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	handler := GetHandler(c.Request.Header.Get("Provider"))
	resp, err := handler.GenerateImages(genImgReq)
	if err != nil {
		slog.Error("Error generating image %s", err)
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, &resp)
}

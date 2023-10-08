package service

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
)

type GenerateImagesRequest struct {
	Prompt         string       `json:"prompt"`
	Number         *int         `json:"n,omitempty"`
	ResponseFormat *ImageFormat `json:"response_format,omitempty"`
	Size           *string      `json:"size,omitempty"`
}

type GenerateImagesResponse struct {
	Images   []string `json:"images"`
	Provider string   `json:"provider"`
}

type ImageFormat string

const (
	ImageFormatBase64 ImageFormat = "b64_json"
	ImageFormatUrl    ImageFormat = "url"
)

func generateImages(c *gin.Context) {
	var genImgReq *GenerateImagesRequest
	if err := c.BindJSON(&genImgReq); err != nil {
		c.IndentedJSON(http.StatusBadRequest, &ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	handler := GetHandler(c)
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

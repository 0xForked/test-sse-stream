package rest

import (
	"github.com/aasumitro/tts/internal/domain"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TTSRESTHandler struct {
	service domain.ITTSService
}

func (handler *TTSRESTHandler) Publish(ctx *gin.Context) {
	type respond struct {
		Code    int
		Status  string
		Message string
	}

	if err := handler.service.PublishEvent(ctx.Request.Context()); err != nil {
		ctx.JSON(http.StatusBadRequest, respond{
			Code:    http.StatusBadRequest,
			Status:  http.StatusText(http.StatusBadRequest),
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, respond{
		Code:    http.StatusCreated,
		Status:  http.StatusText(http.StatusCreated),
		Message: "Event Published",
	})
}

func NewTTSRESTHandler(
	router *gin.RouterGroup,
	service domain.ITTSService,
) {
	handler := &TTSRESTHandler{service: service}
	router.GET("/publish", handler.Publish)
}

package sse

import (
	"context"
	"fmt"
	"github.com/aasumitro/tts/internal/domain"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
)

type TTSSSEHandler struct {
	service domain.ITTSService
}

func (handler *TTSSSEHandler) SSEHeader() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Writer.Header().Set("Content-Type", "text/event-stream")
		ctx.Writer.Header().Set("Cache-Control", "no-cache")
		ctx.Writer.Header().Set("Connection", "keep-alive")
		ctx.Writer.Header().Set("Transfer-Encoding", "chunked")
		ctx.Next()
	}
}

func (handler *TTSSSEHandler) SSEServer(eventStream *domain.TTSEvent) gin.HandlerFunc {
	return func(c *gin.Context) {
		clientChan := make(domain.ClientChan)
		eventStream.NewClients <- clientChan
		defer func() {
			eventStream.ClosedClients <- clientChan
		}()
		c.Set("clientChan", clientChan)
		c.Next()
	}
}

func (handler *TTSSSEHandler) SSEWriter(ctx *gin.Context) {
	v, ok := ctx.Get("clientChan")
	if !ok {
		log.Println("here")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	clientChan, ok := v.(domain.ClientChan)
	if !ok {
		log.Println("here2")
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ctx.Stream(func(w io.Writer) bool {
		if msg, ok := <-clientChan; ok {
			go func() {
				cmd := exec.Command("python3", "-c", fmt.Sprintf(domain.PyScript, msg))
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Println(err)
				}
			}()
			ctx.SSEvent("message", msg)
			return true
		}
		return false
	})
}

func NewSSERESTHandler(
	router *gin.RouterGroup,
	service domain.ITTSService,
) {
	handler := &TTSSSEHandler{service: service}

	eventStream := &domain.TTSEvent{
		Message:       make(chan *domain.TTS),
		NewClients:    make(chan chan string),
		ClosedClients: make(chan chan string),
		TotalClients:  make(map[chan string]bool),
	}

	go eventStream.Listen()

	ttsChan := make(chan *domain.TTS)
	go service.StreamEvent(context.Background(), ttsChan)
	go func() {
		for tts := range ttsChan {
			eventStream.Message <- tts
		}
	}()

	router.GET("/stream",
		handler.SSEHeader(),
		handler.SSEServer(eventStream),
		handler.SSEWriter)
}

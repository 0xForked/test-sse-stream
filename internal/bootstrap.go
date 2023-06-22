package internal

import (
	"github.com/aasumitro/tts/internal/delivery/rest"
	"github.com/aasumitro/tts/internal/delivery/sse"
	"github.com/aasumitro/tts/internal/service"
	"github.com/aasumitro/tts/web"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
)

type boostrap struct {
	ginEngine   *gin.Engine
	redisClient *redis.Client
}

type BoostrapOption func(*boostrap)

func WithGinEngine(ginEngine *gin.Engine) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.ginEngine = ginEngine
	}
}

func WithRedisClient(redisClient *redis.Client) BoostrapOption {
	return func(boostrap *boostrap) {
		boostrap.redisClient = redisClient
	}
}

func RunApp(options ...BoostrapOption) {
	boot := &boostrap{}
	for _, option := range options {
		option(boot)
	}
	boot.newPublicAPIProvider()
	boot.newTTSAPIProvider()
}

func (bootstrap *boostrap) newPublicAPIProvider() {
	router := bootstrap.ginEngine
	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(http.StatusTemporaryRedirect, "/stream")
	})
	router.StaticFS("/stream", http.FS(web.UIResource))
}

func (bootstrap *boostrap) newTTSAPIProvider() {
	routerGroupV1 := bootstrap.ginEngine.Group("api/v1/events")
	ttsService := service.NewTTSService(bootstrap.redisClient)
	rest.NewTTSRESTHandler(routerGroupV1, ttsService)
	sse.NewSSERESTHandler(routerGroupV1, ttsService)
}

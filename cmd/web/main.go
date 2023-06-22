package main

import (
	"github.com/aasumitro/tts/config"
	"github.com/aasumitro/tts/internal"
	"github.com/spf13/viper"
	"log"
)

func main() {
	viper.SetConfigFile(".env")
	config.LoadEnv()
	config.Instance.InitRedisConn()
	config.Instance.InitGinEngine()

	log.Printf("Run %s (%s) . . . .",
		config.Instance.AppName,
		config.Instance.AppVersion)

	internal.RunApp(
		internal.WithGinEngine(config.GinEngine),
		internal.WithRedisClient(config.RedisClient))

	log.Fatalln(config.GinEngine.Run(config.Instance.AppPort))
}

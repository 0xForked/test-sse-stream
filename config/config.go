package config

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"log"
	"sync"
)

var (
	cfgSingleton       sync.Once
	redisSingleton     sync.Once
	ginEngineSingleton sync.Once

	Instance    *Config
	RedisClient *redis.Client
	GinEngine   *gin.Engine
)

type Config struct {
	AppName       string `mapstructure:"APP_NAME"`
	AppDebug      bool   `mapstructure:"APP_DEBUG"`
	AppVersion    string `mapstructure:"APP_VERSION"`
	AppPort       string `mapstructure:"APP_PORT"`
	RedisDsnURL   string `mapstructure:"REDIS_DSN_URL"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
}

func LoadEnv() {
	// notify that app try to load config file
	log.Println("Load configuration file . . . .")
	cfgSingleton.Do(func() {
		// find environment file
		viper.AutomaticEnv()
		// error handling for specific case
		if err := viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				// Config file not found; ignore error if desired
				panic(".env file not found!, please copy .env.example and paste as .env")
			}
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
		// notify that config file is ready
		log.Println("configuration file: ready")
		// extract config to struct
		if err := viper.Unmarshal(&Instance); err != nil {
			panic(fmt.Sprintf("ENV_ERROR: %s", err.Error()))
		}
	})
}

func (cfg *Config) InitGinEngine() {
	log.Println("Trying to init engine . . . .")
	ginEngineSingleton.Do(func() {
		gin.SetMode(func() string {
			if cfg.AppDebug {
				return gin.DebugMode
			}
			return gin.ReleaseMode
		}())
		GinEngine = gin.Default()
		log.Printf("Gin Engine (%s) created  . . . .", gin.Version)
	})
}

func (cfg *Config) InitRedisConn() {
	log.Println("Trying to open redis connection pool . . . .")
	redisSingleton.Do(func() {
		RedisClient = redis.NewClient(&redis.Options{
			Addr:     cfg.RedisDsnURL,
			Password: cfg.RedisPassword,
			DB:       0,
		})
		if err := RedisClient.Ping(context.Background()).Err(); err != nil {
			panic(fmt.Sprintf("REDIS_ERROR: %s", err.Error()))
		}
		log.Println("Redis connection pool created . . . .")
	})
}

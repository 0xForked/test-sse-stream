package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aasumitro/tts/internal/domain"
	"github.com/redis/go-redis/v9"
	"log"
	"math/rand"
	"strconv"
	"time"
)

type ttsService struct {
	redisClient *redis.Client
}

func (service ttsService) PublishEvent(ctx context.Context) error {
	sequence := func() int {
		seq := 1
		data, err := service.redisClient.Get(ctx, domain.SequenceKey).Result()
		if err == nil {
			if cacheSeq, err := strconv.Atoi(data); err == nil {
				seq = cacheSeq
			}
		}
		service.redisClient.Set(ctx, domain.SequenceKey, strconv.Itoa(seq+1), 24*time.Hour)
		return seq
	}()
	locket := func() int {
		lockets := []int{1, 2}
		source := rand.NewSource(time.Now().UnixNano())
		random := rand.New(source)
		randomIndex := random.Intn(len(lockets))
		return lockets[randomIndex]
	}()
	payload, err := json.Marshal(map[string]any{
		"sequence": sequence,
		"locket":   locket,
		"message": fmt.Sprintf(
			"Nomor antrian %d, silahkan menuju ke loket %d",
			sequence, locket),
	})
	if err != nil {
		return err
	}
	service.redisClient.Publish(ctx, domain.TTSChannelName, payload)
	return nil
}

func (service ttsService) StreamEvent(ctx context.Context, ttsChan chan<- *domain.TTS) {
	subscriber := service.redisClient.Subscribe(ctx, domain.TTSChannelName)
	go func() {
		for {
			msg, err := subscriber.ReceiveMessage(context.Background())
			if err != nil {
				ptn := "[%d] - SYNC_EVENT_ERR (QUEUE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				log.Println(msg)
				continue
			}

			eventData := &domain.TTS{}
			if err := json.Unmarshal([]byte(msg.Payload), &eventData); err != nil {
				ptn := "[%d] - SYNC_EVENT_ERR (DECODE): %s"
				msg := fmt.Sprintf(ptn, time.Now().Unix(), err.Error())
				log.Println(msg)
				continue
			}

			ttsChan <- eventData
		}
	}()
}

func NewTTSService(
	redisClient *redis.Client,
) domain.ITTSService {
	return &ttsService{
		redisClient: redisClient,
	}
}

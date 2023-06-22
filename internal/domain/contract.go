package domain

import "context"

type ITTSService interface {
	PublishEvent(ctx context.Context) error
	StreamEvent(ctx context.Context, ttsChan chan<- *TTS)
}

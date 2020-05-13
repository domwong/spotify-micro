package subscriber

import (
	"context"

	log "github.com/micro/go-micro/v2/logger"

	spotify "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
)

type Spotify struct{}

func (e *Spotify) Handle(ctx context.Context, msg *spotify.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *spotify.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}

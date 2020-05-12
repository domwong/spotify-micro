package main

import (
	"spotify/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	spotify "spotify/proto/spotify"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.api.spotify"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	spotify.RegisterSpotifyHandler(service.Server(), new(handler.Spotify))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"github.com/domwong/spotify-micro/api-spotify/handler"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	sps "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
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
	service.Server().Handle(
		service.Server().NewHandler(
			&handler.Spotify{
				Client: sps.NewSpotifyService("go.micro.service.spotify", service.Client()),
			},
		),
	)

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

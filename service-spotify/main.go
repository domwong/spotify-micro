package main

import (
	"github.com/domwong/spotify-micro/service-spotify/dao"
	"github.com/domwong/spotify-micro/service-spotify/handler"
	"github.com/domwong/spotify-micro/service-spotify/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	spotify "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.spotify"),
		micro.Version("latest"),
	)
	// Initialise service
	service.Init()
	dao.Init(service)
	handler.Init()
	// Register Handler
	spotify.RegisterSpotifyHandler(service.Server(), new(handler.Spotify))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.spotify", service.Server(), new(subscriber.Spotify))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"os"

	"github.com/domwong/spotify-micro/service-spotify/dao"
	"github.com/domwong/spotify-micro/service-spotify/handler"
	"github.com/domwong/spotify-micro/service-spotify/subscriber"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	spotify "github.com/zmb3/spotify"

	sproto "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.spotify"),
		micro.Version("latest"),
	)
	auth := spotify.NewAuthenticator(os.Getenv("SPOTIFY_REDIRECT"), spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistReadPrivate)
	// Initialise service
	service.Init()
	dao.Init(service)

	// Register Handler
	sproto.RegisterSpotifyHandler(service.Server(), &handler.Spotify{
		Client: service.Client(),
		Auth:   &auth,
	})

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.service.spotify", service.Server(), &subscriber.Spotify{Auth: &auth})

	// Run it
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}

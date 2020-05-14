package handler

import (
	"context"
	"time"

	sps "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
	api "github.com/micro/go-micro/v2/api/proto"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
)

// Spotify is the api
type Spotify struct {
	Client sps.SpotifyService
}

// RootRedirect is used to start the auth dance
func (s *Spotify) RootRedirect(ctx context.Context, req *api.Request, rsp *api.Response) error {
	rrsp, err := s.Client.RootRedirect(ctx, &sps.RedirectRequest{})
	if err != nil {
		return err
	}
	rsp.StatusCode = 301
	rsp.Header = map[string]*api.Pair{}
	rsp.Header["Location"] = &api.Pair{
		Key:    "Location",
		Values: []string{rrsp.GetRedirectUrl()},
	}
	return nil
}

// Callback does the end of the auth dance
func (s *Spotify) Callback(ctx context.Context, req *api.Request, rsp *api.Response) error {
	var code, state string
	if errPair, ok := req.Get["error"]; ok {
		log.Errorf("Error on auth %s", errPair.Values[0])
		return nil
	}
	if statePair, ok := req.Get["state"]; ok {
		state = statePair.Values[0]
	}
	if codePair, ok := req.Get["code"]; ok {
		code = codePair.Values[0]
	}
	_, err := s.Client.Callback(ctx, &sps.CallbackRequest{
		Code:  code,
		State: state,
	})

	rsp.Body = "Success"
	return err
}

// Save does the hard work
func (s *Spotify) Save(ctx context.Context, req *api.Request, rsp *api.Response) error {

	userPair, ok := req.Get["user"]
	if !ok {
		return errors.BadRequest("go.micro.api.spotify.save", "missing user")
	}

	_, err := s.Client.Save(ctx, &sps.SaveRequest{
		UserName: userPair.Values[0],
	}, client.WithRequestTimeout(5*time.Second)) // don't actually need to change this but this is how you would do it

	return err
}

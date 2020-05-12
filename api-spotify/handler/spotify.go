package handler

import (
	"context"

	sps "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
)

type Spotify struct {
	Client sps.SpotifyService
}

func (s *Spotify) RootRedirect(context.Context, *RedirectRequest, *RedirectResponse) error {

}
func (s *Spotify) Callback(context.Context, *CallbackRequest, *CallbackResponse) error {

}
func (s *Spotify) Save(context.Context, *SaveRequest, *SaveResponse) error {

}

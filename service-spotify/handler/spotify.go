package handler

import (
	"context"
	"encoding/json"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"golang.org/x/oauth2"

	"github.com/domwong/spotify-micro/service-spotify/dao"
	sp "github.com/domwong/spotify-micro/service-spotify/proto/spotify"

	spotify "github.com/zmb3/spotify"
)

var (
	state = "abc123"
)

// Spotify struct
type Spotify struct {
	Client client.Client
	Auth   *spotify.Authenticator
}

// RootRedirect called to kick off the auth dance by returning the URL to redirect the user to
func (e *Spotify) RootRedirect(ctx context.Context, req *sp.RedirectRequest, rsp *sp.RedirectResponse) error {
	log.Info("Received Spotify.RootRedirect request")
	// URL will be auth.AuthURL(state)
	rsp.RedirectUrl = e.Auth.AuthURL(state)
	return nil
}

// Callback invoked by spotify after the auth dance
func (e *Spotify) Callback(ctx context.Context, req *sp.CallbackRequest, rsp *sp.CallbackResponse) error {
	if err := req.GetError(); err != "" {
		return errors.BadRequest("go.micro.service.spotify.callback", "Auth failed "+err)
	}
	code := req.GetCode()
	if code == "" {
		return errors.BadRequest("go.micro.service.spotify.callback", "Didn't get access code")
	}
	actualState := req.GetState()
	if actualState != state {
		return errors.BadRequest("go.micro.service.spotify.callback", "Redirect state parameter doesn't match")
	}
	tok, err := e.Auth.Exchange(code)
	if err != nil {
		log.Errorf("Error exchanging code %s", err)
		return errors.InternalServerError("go.mirco.service.spotify.callback", "Error exchanging code"+err.Error())
	}

	client := e.Auth.NewClient(tok)
	user, err := client.CurrentUser()
	if err != nil {
		log.Errorf("Error retrieving user %s", err)
		return err
	}

	// store tok
	_, err = storeToken(user.ID, tok)
	if err != nil {
		return errors.InternalServerError("go.micro.service.spotify.callback", err.Error())
	}
	return nil
}

// Save the desired playlists
func (e *Spotify) Save(ctx context.Context, in *sp.SaveRequest, out *sp.SaveResponse) error {
	username := in.GetUserName()
	if username == "" {
		return errors.BadRequest("go.micro.service.spotify.save", "Missing username")
	}

	_, err := dao.ReadUserEntry(username) // check this user exists
	if err != nil {
		return errors.InternalServerError("go.micro.service.spotify.save", "Error retrieving entry for user")
	}

	ev := micro.NewEvent("go.micro.service.spotify", e.Client)
	payload := map[string]interface{}{"user": username} // TODO make this a proto defined payload
	b, err := json.Marshal(payload)
	if err != nil {
		return errors.InternalServerError("go.micro.service.spotify.save", "Error marshaling save request")
	}
	if err := ev.Publish(ctx, &sp.Event{Type: "save.request", Payload: b}); err != nil {
		return errors.InternalServerError("go.micro.service.spotify.save", "Error publishing save request")
	}

	return nil
}

func storeToken(username string, tok *oauth2.Token) (*dao.UserEntry, error) {
	ue := &dao.UserEntry{
		Username:  username,
		Token:     *tok,
		Playlists: []string{"Discover Weekly", "Release Radar"}, // TODO make this configurable
	}

	if err := dao.CreateUserEntry(ue); err != nil {
		return nil, err
	}
	return ue, nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Spotify) Stream(ctx context.Context, req *sp.StreamingRequest, stream sp.Spotify_StreamStream) error {
	log.Infof("Received Spotify.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&sp.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Spotify) PingPong(ctx context.Context, stream sp.Spotify_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&sp.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}

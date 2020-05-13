package handler

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/micro/go-micro/v2/errors"
	log "github.com/micro/go-micro/v2/logger"
	"golang.org/x/oauth2"

	"github.com/domwong/spotify-micro/service-spotify/dao"
	sp "github.com/domwong/spotify-micro/service-spotify/proto/spotify"

	spotify "github.com/zmb3/spotify"
)

var (
	redirectURI string // TODO get from env
	auth        spotify.Authenticator
	state       = "abc123"
)

// Init the package
func Init() error {
	redirectURI = os.Getenv("SPOTIFY_REDIRECT")
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistReadPrivate)
	log.Infof("Auth obj created %s", redirectURI)
	return nil
}

// Spotify struct
type Spotify struct{}

// RootRedirect called to kick off the auth dance by returning the URL to redirect the user to
func (e *Spotify) RootRedirect(ctx context.Context, req *sp.RedirectRequest, rsp *sp.RedirectResponse) error {
	log.Info("Received Spotify.RootRedirect request")
	// URL will be auth.AuthURL(state)
	rsp.RedirectUrl = auth.AuthURL(state)
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
	tok, err := auth.Exchange(code)
	if err != nil {
		log.Errorf("Error exchanging code %s", err)
		return errors.InternalServerError("go.mirco.service.spotify.callback", "Error exchanging code"+err.Error())
	}

	client := auth.NewClient(tok)
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

	ue, err := dao.ReadUserEntry(username)
	if err != nil {
		return errors.InternalServerError("go.micro.service.spotify.save", "Error retrieving entry for user")
	}
	auth = spotify.NewAuthenticator(redirectURI, spotify.ScopePlaylistModifyPrivate, spotify.ScopePlaylistReadPrivate)
	client := auth.NewClient(&ue.Token)
	if err := savePlaylists(&client, ue); err != nil {
		return errors.InternalServerError("go.micro.service.spotify.save", err.Error())
	}

	return nil
}

func savePlaylists(client *spotify.Client, userEntry *dao.UserEntry) error {
	limit := 50
	off := 0
	for {
		pls, err := client.CurrentUsersPlaylistsOpt(&spotify.Options{
			Limit:  &limit,
			Offset: &off,
		})

		if err != nil {
			return err
		}
		for _, v := range pls.Playlists {
			found := false
			for _, p := range userEntry.Playlists {
				if v.Name == p {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			pl, err := client.GetPlaylist(v.ID)
			if err != nil {
				return err
			}
			plTracks := make([]spotify.ID, len(pl.Tracks.Tracks))
			for i, tr := range pl.Tracks.Tracks {
				plTracks[i] = tr.Track.ID
			}

			cpl, err := client.CreatePlaylistForUser(userEntry.Username, fmt.Sprintf("%s %s", v.Name, time.Now().Format("2006-01-02")), "Autosaved snapshot", false)
			if err != nil {
				return err
			}

			if _, err = client.AddTracksToPlaylist(cpl.ID, plTracks...); err != nil {
				return err
			}

		}
		if len(pls.Playlists) < limit {
			break
		}
		off++

	}
	return nil
}

func storeToken(username string, tok *oauth2.Token) (*dao.UserEntry, error) {
	ue := &dao.UserEntry{
		Username:  username,
		Token:     *tok,
		Playlists: []string{"Discover Weekly"}, // TODO make this configurable
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

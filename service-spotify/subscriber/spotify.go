package subscriber

import (
	"context"
	"fmt"
	"time"

	"github.com/micro/go-micro/errors"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/zmb3/spotify"
	"gopkg.in/square/go-jose.v2/json"

	"github.com/domwong/spotify-micro/service-spotify/dao"
	sproto "github.com/domwong/spotify-micro/service-spotify/proto/spotify"
)

type Spotify struct {
	Auth *spotify.Authenticator
}

func (e *Spotify) Handle(ctx context.Context, msg *sproto.Event) error {
	log.Info("Handler Received message: ", msg.Type)
	switch msg.Type {
	case "save.request":
		out := map[string]string{}
		if err := json.Unmarshal(msg.Payload, &out); err != nil {
			log.Errorf("Error processing save request %s", err)
			return nil
		}
		if err := e.savePlaylists(out["user"]); err != nil {
			log.Errorf("Error processing save request %s", err)
			return nil
		}
	default:
		log.Warnf("Received unrecognised message type %s", msg.Type)
	}
	return nil
}

func (e *Spotify) savePlaylists(username string) error {
	ue, err := dao.ReadUserEntry(username)
	if err != nil {
		return errors.InternalServerError("go.micro.service.spotify.save", "Error retrieving entry for user")
	}

	client := e.Auth.NewClient(&ue.Token)
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
			for _, p := range ue.Playlists {
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

			cpl, err := client.CreatePlaylistForUser(ue.Username, fmt.Sprintf("%s %s", v.Name, time.Now().Format("2006-01-02")), "Autosaved snapshot", false)
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
		off = off + limit

	}
	log.Infof("Playlist created for user %s", ue.Username)
	return nil
}

package dao

import (
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/store"
	"golang.org/x/oauth2"
)

// UserEntry represents the user and their auth
type UserEntry struct {
	Username  string       `json:"username"`
	Token     oauth2.Token `json:"token"`
	Playlists []string     `json:"playlists"`
}

var (
	s micro.Service
)

// Init the package
func Init(srv micro.Service) {
	s = srv
}

// CreateUserEntry creates in the store
func CreateUserEntry(ue *UserEntry) error {
	b, err := json.Marshal(ue)
	if err != nil {
		return err
	}
	return s.Options().Store.Write(&store.Record{
		Key:   ue.Username,
		Value: b,
	})
}

// ReadUserEntry returns record for this user
func ReadUserEntry(userName string) (*UserEntry, error) {
	records, err := s.Options().Store.Read(userName)
	if err != nil {
		return nil, err
	}
	if len(records) != 1 {
		return nil, fmt.Errorf("Number of records is incorrect %d", len(records))
	}
	retVal := &UserEntry{}
	if err := json.Unmarshal(records[0].Value, retVal); err != nil {
		return nil, err
	}
	return retVal, nil
}

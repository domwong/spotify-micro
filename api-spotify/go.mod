module github.com/domwong/spotify-micro/api-spotify

go 1.14

require (
	github.com/domwong/spotify-micro/service-spotify v0.0.0-20200514114908-9480832df0cb
	github.com/micro/go-micro/v2 v2.6.0
)

replace github.com/domwong/spotify-micro/service-spotify => ../service-spotify

module github.com/domwong/spotify-micro/api-spotify

go 1.14

require (
	github.com/domwong/spotify-micro/service-spotify v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.4.1
	github.com/micro/go-micro/v2 v2.6.0
	google.golang.org/protobuf v1.22.0
)

replace github.com/domwong/spotify-micro/service-spotify => /Users/domwong/development/spotify-project/service-spotify

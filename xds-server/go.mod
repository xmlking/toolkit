module github.com/xmlking/toolkit/xds-server

go 1.16

replace github.com/xmlking/toolkit/logger => ./

require (
	github.com/cockroachdb/errors v1.8.4
	github.com/rs/zerolog v1.22.0
	google.golang.org/grpc v1.38.0
)

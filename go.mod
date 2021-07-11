module github.com/xmlking/toolkit

go 1.16

require (
	cloud.google.com/go/pubsub v1.12.2
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/cockroachdb/errors v1.8.4
	github.com/cockroachdb/redact v1.0.9
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/rs/zerolog v1.22.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/api v0.50.0
	google.golang.org/genproto v0.0.0-20210707164411-8c882eb9abba
	google.golang.org/grpc v1.39.0
)

//replace github.com/xmlking/toolkit => ./
//
//replace github.com/xmlking/toolkit/confy => ./confy

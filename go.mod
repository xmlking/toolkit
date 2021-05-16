module github.com/xmlking/toolkit

go 1.16

require (
	cloud.google.com/go/pubsub v1.10.1
	github.com/cloudevents/sdk-go/v2 v2.4.1
	github.com/cockroachdb/errors v1.8.4
	github.com/cockroachdb/redact v1.0.9
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/rs/zerolog v1.21.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	google.golang.org/api v0.46.0
	google.golang.org/genproto v0.0.0-20210429181445-86c259c2b4ab
	google.golang.org/grpc v1.37.0
)

//replace github.com/xmlking/toolkit => ./
//
//replace github.com/xmlking/toolkit/confy => ./confy

module github.com/xmlking/toolkit

go 1.16

require (
	cloud.google.com/go/pubsub v1.9.1
	github.com/cloudevents/sdk-go/v2 v2.3.1
	github.com/cockroachdb/errors v1.8.2
	github.com/cockroachdb/redact v1.0.9
	github.com/gogo/protobuf v1.3.2
	github.com/google/uuid v1.2.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2
	github.com/hashicorp/golang-lru v0.5.4
	github.com/rs/zerolog v1.20.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20201207232520-09787c993a3a
	google.golang.org/api v0.36.0
	google.golang.org/genproto v0.0.0-20201209185603-f92720507ed4
	google.golang.org/grpc v1.35.0
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

//replace github.com/xmlking/toolkit => ./
//
//replace github.com/xmlking/toolkit/confy => ./confy

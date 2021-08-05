module github.com/xmlking/toolkit

go 1.16

require (
	github.com/cockroachdb/errors v1.8.6
	github.com/cockroachdb/redact v1.1.3
	github.com/gogo/protobuf v1.3.2
	github.com/google/go-cmp v0.5.6 // indirect
	github.com/google/uuid v1.3.0
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/hashicorp/golang-lru v0.5.4
	github.com/rs/zerolog v1.23.0
	github.com/stretchr/testify v1.7.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	google.golang.org/genproto v0.0.0-20210722135532-667f2b7c528f // indirect
	google.golang.org/grpc v1.39.0
)

//replace github.com/xmlking/toolkit => ./
//
//replace github.com/xmlking/toolkit/confy => ./confy

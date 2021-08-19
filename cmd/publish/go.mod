module github.com/xmlking/toolkit/cmd/publish

go 1.16

replace github.com/xmlking/toolkit/broker/pubsub => ../../broker/pubsub

require (
	cloud.google.com/go/pubsub v1.13.0
	github.com/google/uuid v1.3.0
	github.com/rs/zerolog v1.23.0
	github.com/xmlking/toolkit/broker/pubsub v0.2.3
)

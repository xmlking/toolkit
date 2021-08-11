module github.com/xmlking/toolkit/cmd/subscribe

go 1.16

//replace github.com/xmlking/toolkit/broker/pubsub => ./broker/pubsub

require (
	cloud.google.com/go/pubsub v1.13.0
	github.com/rs/zerolog v1.23.0
	github.com/xmlking/toolkit/broker/pubsub v0.2.1
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

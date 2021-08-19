module github.com/xmlking/toolkit/broker/cloudevents

go 1.16

replace github.com/xmlking/toolkit => ../..

require (
	github.com/cloudevents/sdk-go/v2 v2.5.0
	github.com/rs/zerolog v1.23.0
	github.com/xmlking/toolkit v0.2.3
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
)

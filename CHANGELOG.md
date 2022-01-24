# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]


<a name="v0.2.4"></a>
## [v0.2.4] - 2022-01-24
### Build
- **cleanup:** remove unused file
- **deps:** updated deps
- **deps:** updated to golang 1.17 and deps
- **deps:** updated deps
- **deps:** updated deps
- **deps:** updated deps
- **deps:** updated deps
- **deps:** deps updated
- **makefile:** polish makefile
- **mod:** updated command to go mod tidy -compat=1.17
- **mod:** fix go mod
- **release:** update changelog

### Feat
- **xds:** adding SDK to bootstrap xDS server
- **xds:** initial commit for xds-server

### Fix
- **config:** getEnvironment() now check if DefaultConfy is nil


<a name="logger/v0.2.3"></a>
## [logger/v0.2.3] - 2021-08-14

<a name="broker/cloudevents/v0.2.3"></a>
## [broker/cloudevents/v0.2.3] - 2021-08-14

<a name="cmd/publish/v0.2.3"></a>
## [cmd/publish/v0.2.3] - 2021-08-14

<a name="broker/pubsub/v0.2.3"></a>
## [broker/pubsub/v0.2.3] - 2021-08-14

<a name="telemetry/v0.2.3"></a>
## [telemetry/v0.2.3] - 2021-08-14

<a name="cmd/subscribe/v0.2.3"></a>
## [cmd/subscribe/v0.2.3] - 2021-08-14

<a name="cmd/emulator/v0.2.3"></a>
## [cmd/emulator/v0.2.3] - 2021-08-14

<a name="confy/v0.2.3"></a>
## [confy/v0.2.3] - 2021-08-14

<a name="v0.2.3"></a>
## [v0.2.3] - 2021-08-14
### Build
- **deps:** updated deps
- **release:** update changelog


<a name="broker/cloudevents/v0.2.2"></a>
## [broker/cloudevents/v0.2.2] - 2021-08-11

<a name="cmd/subscribe/v0.2.2"></a>
## [cmd/subscribe/v0.2.2] - 2021-08-11

<a name="logger/v0.2.2"></a>
## [logger/v0.2.2] - 2021-08-11

<a name="telemetry/v0.2.2"></a>
## [telemetry/v0.2.2] - 2021-08-11

<a name="cmd/publish/v0.2.2"></a>
## [cmd/publish/v0.2.2] - 2021-08-11

<a name="broker/pubsub/v0.2.2"></a>
## [broker/pubsub/v0.2.2] - 2021-08-11

<a name="confy/v0.2.2"></a>
## [confy/v0.2.2] - 2021-08-11

<a name="cmd/emulator/v0.2.2"></a>
## [cmd/emulator/v0.2.2] - 2021-08-11

<a name="v0.2.2"></a>
## [v0.2.2] - 2021-08-11
### Build
- **deps:** updated deps
- **release:** update changelog

### Ci
- **release:** updated release command for Taskfile

### Fix
- **telemetry:** now auto Register NewGoCollector


<a name="cmd/subscribe/v0.2.1"></a>
## [cmd/subscribe/v0.2.1] - 2021-08-05

<a name="cmd/publish/v0.2.1"></a>
## [cmd/publish/v0.2.1] - 2021-08-05

<a name="logger/v0.2.1"></a>
## [logger/v0.2.1] - 2021-08-05

<a name="broker/cloudevents/v0.2.1"></a>
## [broker/cloudevents/v0.2.1] - 2021-08-05

<a name="broker/pubsub/v0.2.1"></a>
## [broker/pubsub/v0.2.1] - 2021-08-05

<a name="telemetry/v0.2.1"></a>
## [telemetry/v0.2.1] - 2021-08-05

<a name="cmd/emulator/v0.2.1"></a>
## [cmd/emulator/v0.2.1] - 2021-08-05

<a name="confy/v0.2.1"></a>
## [confy/v0.2.1] - 2021-08-05

<a name="v0.2.1"></a>
## [v0.2.1] - 2021-08-05
### Build
- **deps:** fix deps
- **deps:** updated deps
- **release:** update changelog

### Feat
- **modules:** adding go.mod for each module
- **telemetry:** adding openTelemetry initialization helpers

### Fix
- **deps:** updated deps
- **telemetry:** adding const for GCP/PROMETHEUS/STDOUT


<a name="logger/v0.2.0"></a>
## [logger/v0.2.0] - 2021-06-05

<a name="confy/v0.2.0"></a>
## [confy/v0.2.0] - 2021-06-05

<a name="v0.2.0"></a>
## [v0.2.0] - 2021-06-05
### Build
- **deps:** updated deps
- **release:** update changelog

### Docs
- **logger:** updated logger docs and GetListener now support xds

### Refactor
- **broker:** refactor/polich broker and grpc server components
- **server:** now using healthServer ref for updates
- **server:** now using healthServer ref for updates
- **server:** now using healthServer ref for updates
- **server:** added API to updateing Health Status anytime from anywhere


<a name="logger/v0.1.6"></a>
## [logger/v0.1.6] - 2021-04-03

<a name="confy/v0.1.6"></a>
## [confy/v0.1.6] - 2021-04-03

<a name="v0.1.6"></a>
## [v0.1.6] - 2021-04-03
### Build
- **deps:** bump cloud.google.com/go/pubsub from 1.9.1 to 1.10.1
- **release:** update changelog

### Fix
- **xfs:** change log level to debug


<a name="confy/v0.1.5"></a>
## [confy/v0.1.5] - 2021-04-01

<a name="logger/v0.1.5"></a>
## [logger/v0.1.5] - 2021-04-01

<a name="v0.1.5"></a>
## [v0.1.5] - 2021-04-01
### Build
- **release:** update changelog
- **release:** update changelog

### Docs
- **logger:** updated docs on auto and manual init modes

### Fix
- **signals:** adding build tags for windows and linux

### Style
- **code:** fix code lint


<a name="confy/v0.1.4"></a>
## [confy/v0.1.4] - 2021-03-26

<a name="logger/v0.1.4"></a>
## [logger/v0.1.4] - 2021-03-26

<a name="v0.1.4"></a>
## [v0.1.4] - 2021-03-26
### Build
- **go:** fix go mod
- **release:** updated changelogs

### Refactor
- **xfs:** adding debug logs

### Test
- **xfs:** adding more test cases


<a name="logger/v0.1.3"></a>
## [logger/v0.1.3] - 2021-03-25

<a name="confy/v0.1.3"></a>
## [confy/v0.1.3] - 2021-03-25

<a name="v0.1.3"></a>
## [v0.1.3] - 2021-03-25
### Build
- **deps:** updated deps
- **docker:** updated docker
- **release:** updated changelog

### Feat
- **logger:** seperate auto logger config file
- **taskfile:** updated release task

### Fix
- **logger:** fix space before CONFY_LOG_FORMAT
- **xfs:** now support loading files with absolute path


<a name="v0.1.2"></a>
## [v0.1.2] - 2021-02-19

<a name="confy/v0.1.2"></a>
## [confy/v0.1.2] - 2021-02-19
### Build
- **deps:** updated deps

### Chore
- **deps:** update actions/setup-go action to v2
- **deps:** update actions/checkout action to v2
- **deps:** update github/super-linter docker tag to v2.2.2
- **deps:** update actions/upload-release-asset action to v1.0.2

### Feat
- **broker:** adding RecoveryHandler option
- **broker:** removed WithSubscriptionID
- **configurator:** moved https://github.com/xmlking/configor to toolkit
- **confy:** updated confy
- **confy:** upgraded confy to use golang 1.16 FileSystem

### Fix
- **broker:** handle error with client.Close() due to pubsub v1.8.0
- **configurator:** fix tests

### Refactor
- **broker:** rename Subscriber to NewSubscriber

### Test
- **broker:** adding mocks for Ack() Nack() methods on pubsub.Message
- **configurator:** fix tests


<a name="v0.1.1"></a>
## [v0.1.1] - 2020-09-30
### Build
- **deps:** updated deps
- **release:** updating changelog

### Chore
- **deps:** update golang docker tag to v1.15

### Docs
- **util:** updated tls util
- **util:** updated tls util

### Feat
- **broker:** adding pubsub broker, errors, logger
- **cache:** adding cache package
- **crypto:** adding crypto , errors modules
- **util:** polish ioutil, adding cmd packages for broker testing

### Fix
- **errors:** makring errors.Code as interface

### Refactor
- **broker:** moving original broker to cloudevents
- **cleanup:** prune

### Test
- **ioutils:** fix tests


<a name="v0.1.0"></a>
## v0.1.0 - 2020-07-06
### Build
- **clog:** updated change logs
- **clog:** updating changelog
- **clog:** updating changelog

### Feat
- **toolkit:** first Commit

### Refactor
- **lint:** lint go.mod
- **translog:** remove translog middleware


[Unreleased]: https://github.com/xmlking/toolkit/compare/v0.2.4...HEAD
[v0.2.4]: https://github.com/xmlking/toolkit/compare/logger/v0.2.3...v0.2.4
[logger/v0.2.3]: https://github.com/xmlking/toolkit/compare/broker/cloudevents/v0.2.3...logger/v0.2.3
[broker/cloudevents/v0.2.3]: https://github.com/xmlking/toolkit/compare/cmd/publish/v0.2.3...broker/cloudevents/v0.2.3
[cmd/publish/v0.2.3]: https://github.com/xmlking/toolkit/compare/broker/pubsub/v0.2.3...cmd/publish/v0.2.3
[broker/pubsub/v0.2.3]: https://github.com/xmlking/toolkit/compare/telemetry/v0.2.3...broker/pubsub/v0.2.3
[telemetry/v0.2.3]: https://github.com/xmlking/toolkit/compare/cmd/subscribe/v0.2.3...telemetry/v0.2.3
[cmd/subscribe/v0.2.3]: https://github.com/xmlking/toolkit/compare/cmd/emulator/v0.2.3...cmd/subscribe/v0.2.3
[cmd/emulator/v0.2.3]: https://github.com/xmlking/toolkit/compare/confy/v0.2.3...cmd/emulator/v0.2.3
[confy/v0.2.3]: https://github.com/xmlking/toolkit/compare/v0.2.3...confy/v0.2.3
[v0.2.3]: https://github.com/xmlking/toolkit/compare/broker/cloudevents/v0.2.2...v0.2.3
[broker/cloudevents/v0.2.2]: https://github.com/xmlking/toolkit/compare/cmd/subscribe/v0.2.2...broker/cloudevents/v0.2.2
[cmd/subscribe/v0.2.2]: https://github.com/xmlking/toolkit/compare/logger/v0.2.2...cmd/subscribe/v0.2.2
[logger/v0.2.2]: https://github.com/xmlking/toolkit/compare/telemetry/v0.2.2...logger/v0.2.2
[telemetry/v0.2.2]: https://github.com/xmlking/toolkit/compare/cmd/publish/v0.2.2...telemetry/v0.2.2
[cmd/publish/v0.2.2]: https://github.com/xmlking/toolkit/compare/broker/pubsub/v0.2.2...cmd/publish/v0.2.2
[broker/pubsub/v0.2.2]: https://github.com/xmlking/toolkit/compare/confy/v0.2.2...broker/pubsub/v0.2.2
[confy/v0.2.2]: https://github.com/xmlking/toolkit/compare/cmd/emulator/v0.2.2...confy/v0.2.2
[cmd/emulator/v0.2.2]: https://github.com/xmlking/toolkit/compare/v0.2.2...cmd/emulator/v0.2.2
[v0.2.2]: https://github.com/xmlking/toolkit/compare/cmd/subscribe/v0.2.1...v0.2.2
[cmd/subscribe/v0.2.1]: https://github.com/xmlking/toolkit/compare/cmd/publish/v0.2.1...cmd/subscribe/v0.2.1
[cmd/publish/v0.2.1]: https://github.com/xmlking/toolkit/compare/logger/v0.2.1...cmd/publish/v0.2.1
[logger/v0.2.1]: https://github.com/xmlking/toolkit/compare/broker/cloudevents/v0.2.1...logger/v0.2.1
[broker/cloudevents/v0.2.1]: https://github.com/xmlking/toolkit/compare/broker/pubsub/v0.2.1...broker/cloudevents/v0.2.1
[broker/pubsub/v0.2.1]: https://github.com/xmlking/toolkit/compare/telemetry/v0.2.1...broker/pubsub/v0.2.1
[telemetry/v0.2.1]: https://github.com/xmlking/toolkit/compare/cmd/emulator/v0.2.1...telemetry/v0.2.1
[cmd/emulator/v0.2.1]: https://github.com/xmlking/toolkit/compare/confy/v0.2.1...cmd/emulator/v0.2.1
[confy/v0.2.1]: https://github.com/xmlking/toolkit/compare/v0.2.1...confy/v0.2.1
[v0.2.1]: https://github.com/xmlking/toolkit/compare/logger/v0.2.0...v0.2.1
[logger/v0.2.0]: https://github.com/xmlking/toolkit/compare/confy/v0.2.0...logger/v0.2.0
[confy/v0.2.0]: https://github.com/xmlking/toolkit/compare/v0.2.0...confy/v0.2.0
[v0.2.0]: https://github.com/xmlking/toolkit/compare/logger/v0.1.6...v0.2.0
[logger/v0.1.6]: https://github.com/xmlking/toolkit/compare/confy/v0.1.6...logger/v0.1.6
[confy/v0.1.6]: https://github.com/xmlking/toolkit/compare/v0.1.6...confy/v0.1.6
[v0.1.6]: https://github.com/xmlking/toolkit/compare/confy/v0.1.5...v0.1.6
[confy/v0.1.5]: https://github.com/xmlking/toolkit/compare/logger/v0.1.5...confy/v0.1.5
[logger/v0.1.5]: https://github.com/xmlking/toolkit/compare/v0.1.5...logger/v0.1.5
[v0.1.5]: https://github.com/xmlking/toolkit/compare/confy/v0.1.4...v0.1.5
[confy/v0.1.4]: https://github.com/xmlking/toolkit/compare/logger/v0.1.4...confy/v0.1.4
[logger/v0.1.4]: https://github.com/xmlking/toolkit/compare/v0.1.4...logger/v0.1.4
[v0.1.4]: https://github.com/xmlking/toolkit/compare/logger/v0.1.3...v0.1.4
[logger/v0.1.3]: https://github.com/xmlking/toolkit/compare/confy/v0.1.3...logger/v0.1.3
[confy/v0.1.3]: https://github.com/xmlking/toolkit/compare/v0.1.3...confy/v0.1.3
[v0.1.3]: https://github.com/xmlking/toolkit/compare/v0.1.2...v0.1.3
[v0.1.2]: https://github.com/xmlking/toolkit/compare/confy/v0.1.2...v0.1.2
[confy/v0.1.2]: https://github.com/xmlking/toolkit/compare/v0.1.1...confy/v0.1.2
[v0.1.1]: https://github.com/xmlking/toolkit/compare/v0.1.0...v0.1.1

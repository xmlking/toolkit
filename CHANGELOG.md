# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

<a name="unreleased"></a>
## [Unreleased]


<a name="v0.1.4"></a>
## [v0.1.4] - 2021-03-26
### Build
- **go:** fix go mod

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


<a name="confy/v0.1.2"></a>
## [confy/v0.1.2] - 2021-02-19

<a name="v0.1.2"></a>
## [v0.1.2] - 2021-02-19
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


[Unreleased]: https://github.com/xmlking/toolkit/compare/v0.1.4...HEAD
[v0.1.4]: https://github.com/xmlking/toolkit/compare/logger/v0.1.3...v0.1.4
[logger/v0.1.3]: https://github.com/xmlking/toolkit/compare/confy/v0.1.3...logger/v0.1.3
[confy/v0.1.3]: https://github.com/xmlking/toolkit/compare/v0.1.3...confy/v0.1.3
[v0.1.3]: https://github.com/xmlking/toolkit/compare/confy/v0.1.2...v0.1.3
[confy/v0.1.2]: https://github.com/xmlking/toolkit/compare/v0.1.2...confy/v0.1.2
[v0.1.2]: https://github.com/xmlking/toolkit/compare/v0.1.1...v0.1.2
[v0.1.1]: https://github.com/xmlking/toolkit/compare/v0.1.0...v0.1.1

package xds

import (
	xlog "github.com/envoyproxy/go-control-plane/pkg/log"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type logger struct {
	log zerolog.Logger
}

var _ xlog.Logger = (*logger)(nil)

func newXdsLogger() xlog.Logger {
	return &logger{log: log.With().Str("component", "xds").Logger()}
}

func (x logger) Debugf(format string, args ...interface{}) {
	x.log.Debug().Msgf(format, args...)
}

func (x logger) Infof(format string, args ...interface{}) {
	x.log.Info().Msgf(format, args...)
}

func (x logger) Warnf(format string, args ...interface{}) {
	x.log.Warn().Msgf(format, args...)
}

func (x logger) Errorf(format string, args ...interface{}) {
	x.log.Error().Msgf(format, args...)
}

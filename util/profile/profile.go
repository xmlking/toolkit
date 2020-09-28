package profile

import (
	"time"

	"github.com/rs/zerolog/log"
)

func Duration(invocation time.Time, name string) {
	elapsed := time.Since(invocation)

	log.Debug().Msgf("%s lasted %s", name, elapsed)
}

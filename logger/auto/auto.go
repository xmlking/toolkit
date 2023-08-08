package auto

import (
	"os"
	"strconv"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/logger"
	"github.com/xmlking/toolkit/util/environs"
)

func init() {
	var opts []logger.Option

	if lvlStr := os.Getenv("CONFY_LOG_LEVEL"); len(lvlStr) > 0 {
		if lvl, err := zerolog.ParseLevel(lvlStr); err != nil {
			log.Fatal().Err(err).Send()
		} else {
			opts = append(opts, logger.WithLevel(lvl))
		}
	}

	if fmtStr := os.Getenv("CONFY_LOG_FORMAT"); len(fmtStr) > 0 {
		if logFmt, err := logger.ParseFormat(fmtStr); err != nil {
			log.Fatal().Err(err).Send()
		} else {
			opts = append(opts, logger.WithFormat(logFmt))
		}
	}

	if enableGrpcLog, _ := strconv.ParseBool(os.Getenv("CONFY_LOG_GRPC")); enableGrpcLog {
		opts = append(opts, logger.EnableGrpcLog(enableGrpcLog))
	}

	if logFileName := os.Getenv("CONFY_LOG_FILE"); len(logFileName) > 0 {
		lw := logger.FileWriter(logFileName, logger.FileConfig{
			MaxSize:    environs.GetEnvAsInt("CONFY_LOG_FILE_MAX_SIZE", 5),
			MaxBackups: environs.GetEnvAsInt("CONFY_LOG_FILE_MAX_BACKUPS", 10),
			MaxAge:     environs.GetEnvAsInt("CONFY_LOG_FILE_MAX_AGE", 14),
		})
		opts = append(opts, logger.WithOutput(lw))

		/**
		  // TODO defer file.Close()
		  if file, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, os.FileMode(0o666)); err != nil {
		    log.Fatal().Err(err).Send()
		  } else {
		    // Merging log writers: Stderr output and file output
		    multi := zerolog.MultiLevelWriter(os.Stderr, file)
		    opts = append(opts, logger.WithOutput(multi))
		  **/
	}

	logger.DefaultLogger = logger.NewLogger(opts...)
}

package config

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"path"
	"time"
)

// logFile is the file to write logs to
var logFile *os.File

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// Set the default logger level to Debug if the verbose flag is set
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if GlobalConfig.Verbose {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if GlobalConfig.EnableLogFile {
		var err error
		logFile, err = os.OpenFile(
			path.Join("logs", GlobalConfig.LogFileName),
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)
		if err != nil {
			log.Panic().Err(err).Msg("Failed to open log file")
		}

		consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.TimeOnly}
		multi := zerolog.MultiLevelWriter(consoleWriter, logFile)

		// Set the default logger to write to both console and file
		log.Logger = zerolog.New(multi).With().Timestamp().Logger()
	} else {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.TimeOnly,
		})
	}

	log.Debug().Msg("Logger initialized")
}

func CloseLogger() {
	if logFile != nil {
		if err := logFile.Close(); err != nil {
			log.Error().Err(err).Msg("Failed to close log file")
		}
		log.Debug().Msg("Logger stopped")
	}
}

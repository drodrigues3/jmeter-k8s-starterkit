package log

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/drodrigues3/jmeter-k8s-starterkit/config"
	"github.com/rs/zerolog"
)

func New() zerolog.Logger {

	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}

	// Output Log formatter
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf("%s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf("%s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s", i))
	}

	log := zerolog.New(output).With().Timestamp().Logger()

	cfg, err := config.LoadConfiguration()
	if err != nil {
		log.Panic().Msg("Error when try load configuration")
	}

	// Default level default is info
	zerolog.SetGlobalLevel(zerolog.Level(cfg.LogLevel.Level))

	log.Debug().Msg("Log level defined to " + zerolog.GlobalLevel().String())

	return log

}

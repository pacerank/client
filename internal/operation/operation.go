package operation

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"os"
	"runtime"
	"time"
)

// Nil placeholder for logging
var file *os.File

func formatMessage(i interface{}) string {
	if i == nil {
		return "|"
	}
	return fmt.Sprintf("| %-6s", i)
}

func Setup() error {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	var err error

	// Try to open file after setup to get correct logging structure
	file, err = os.OpenFile(fmt.Sprintf("%s_client.log", time.Now().Format("2006-01-02")), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	fileWriter := zerolog.ConsoleWriter{
		Out:           file,
		NoColor:       true,
		TimeFormat:    "2006-01-02T15:04:05Z",
		FormatMessage: formatMessage,
	}

	stdoutWriter := zerolog.ConsoleWriter{
		Out:           os.Stdout,
		NoColor:       runtime.GOOS != "linux",
		TimeFormat:    "2006-01-02T15:04:05Z",
		FormatMessage: formatMessage,
	}

	// Make a writer to both log types
	mw := io.MultiWriter(fileWriter, stdoutWriter)

	// Pretty log
	log.Logger = log.Output(mw)

	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	return err
}

func EnableDebug() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func Close() error {
	err := file.Close()
	if err != nil {
		return err
	}

	return nil
}

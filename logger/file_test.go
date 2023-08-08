package logger_test

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/xmlking/toolkit/logger"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"
	"time"
)

func TestFileWriter(t *testing.T) {
	filePath := filepath.Join(t.TempDir(), "testFile")
	lw := logger.FileWriter(filePath, logger.FileConfig{})

	logger.DefaultLogger = logger.NewLogger(logger.WithOutput(lw), logger.WithFormat(logger.JSON), logger.WithLevel(zerolog.WarnLevel))
	zerolog.TimeFieldFormat = "2006"

	log.Warn().Msg("msg")

	if c, ok := lw.(io.Closer); ok {
		if err := c.Close(); err != nil {
			t.Error(err)
		}
	}

	f, err := ioutil.ReadFile(filePath)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := string(f), fmt.Sprintf("{\"level\":\"warn\",\"time\":\"%s\",\"message\":\"msg\"}\n", time.Now().Format("2006")); got != want {
		t.Errorf("\ngot:\n%s\nwant:\n%s", got, want)
	}
}

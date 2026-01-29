package log

import (
	"os"
	"path/filepath"

	"simple_im/pkg/common/config"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

func InitLog(conf config.LoggerConfig) {
	// Ensure log directory exists
	logDir := filepath.Dir(conf.Filename)
	if logDir != "" && logDir != "." {
		os.MkdirAll(logDir, 0755)
	}

	rotatingLogger := &lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: 3,
		MaxAge:     28,
		Compress:   true,
	}

	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "2006-01-02 15:04:05"}
	multi := zerolog.MultiLevelWriter(consoleWriter, rotatingLogger)

	log.Logger = zerolog.New(multi).With().Caller().Timestamp().Logger()
}

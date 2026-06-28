package logger

import (
	"fmt"
	"io"
	"net"
	"os"

	"backend/internal/config"

	"github.com/sirupsen/logrus"
)

func ObfuscateIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	ip := net.ParseIP(host)
	if ip == nil {
		return host
	}
	if ip4 := ip.To4(); ip4 != nil {
		return fmt.Sprintf("%d.%d.%d.xxx", ip4[0], ip4[1], ip4[2])
	}
	return fmt.Sprintf("%x:%x:%x:%x:xxxx:xxxx:xxxx:xxxx", ip[0], ip[1], ip[2], ip[3])
}

func ObfuscateUA(ua string) string {
	if len(ua) > 60 {
		return ua[:60] + "..."
	}
	return ua
}

// Logger представляет логгер приложения
type Logger struct {
	*logrus.Logger
}

// New создает новый экземпляр логгера
func New(cfg *config.LoggerConfig) *Logger {
	log := logrus.New()

	// Установка уровня логирования
	level, err := logrus.ParseLevel(cfg.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	// Установка формата логов
	if cfg.Format == "json" {
		log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
		})
	} else {
		log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
		})
	}

	// Настройка вывода в файл
	if cfg.File != "" {
		file, err := os.OpenFile(cfg.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err == nil {
			log.SetOutput(io.MultiWriter(os.Stdout, file))
		} else {
			log.WithError(err).Error("Failed to open log file, using stdout only")
		}
	}

	return &Logger{Logger: log}
}

// WithField добавляет поле к логгеру
func (l *Logger) WithField(key string, value interface{}) *logrus.Entry {
	return l.Logger.WithField(key, value)
}

// WithFields добавляет несколько полей к логгеру
func (l *Logger) WithFields(fields logrus.Fields) *logrus.Entry {
	return l.Logger.WithFields(fields)
}

// WithError добавляет ошибку к логгеру
func (l *Logger) WithError(err error) *logrus.Entry {
	return l.Logger.WithError(err)
}
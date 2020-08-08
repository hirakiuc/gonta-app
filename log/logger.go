package log

import (
	"fmt"
	"os"

	"github.com/blendle/zapdriver"
	"go.uber.org/zap"
)

// GetLogger return a zap.Logger instance.
func GetLogger() *zap.Logger {
	logger, err := zapdriver.NewDevelopmentWithCore(
		zapdriver.WrapCore(
			zapdriver.ReportAllErrors(true),
			zapdriver.ServiceName("gonta"),
		),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger.")

		return nil
	}

	return logger
}

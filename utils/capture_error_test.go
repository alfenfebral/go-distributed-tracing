package utils_test

import (
	"errors"
	"go-clean-architecture/utils"
	"os"
	"testing"
)

func TestCaptureError(t *testing.T) {
	t.Run("when enabled", func(t *testing.T) {
		os.Setenv("ENABLE_SENTRY_LOG", "true")
		utils.CaptureError(errors.New("error"))
	})

	t.Run("when disabled", func(t *testing.T) {
		os.Setenv("ENABLE_SENTRY_LOG", "false")
		utils.CaptureError(errors.New("error"))
	})
}

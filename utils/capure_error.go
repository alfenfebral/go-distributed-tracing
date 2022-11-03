package utils

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func CaptureError(err error) {
	if os.Getenv("ENABLE_SENTRY_LOG") == "true" {
		sentry.CaptureException(err)
	}
	logrus.Error(err)
}

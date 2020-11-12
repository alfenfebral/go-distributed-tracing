package utils

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

func CaptureError(err error) {
	if os.Getenv("GO_ENV") != "test" {
		if os.Getenv("DEBUG") == "staging" || os.Getenv("DEBUG") == "dev" {
			logrus.Error(err)
			sentry.CaptureException(err)
		} else if os.Getenv("DEBUG") == "prod" {
			sentry.CaptureException(err)
		} else {
			logrus.Error(err)
		}
	}
}

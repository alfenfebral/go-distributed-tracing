package config

import (
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"

	mgo "gopkg.in/mgo.v2"
)

var (
	sessionMgo *mgo.Session
)

// InitMgo - initialize mgo database
func InitMgo() (*mgo.Session, error) {
	session, err := mgo.Dial(os.Getenv("DB_URL"))
	if err != nil {
		logrus.Error(err)
		sentry.CaptureException(err)

		return session, err
	}

	session.SetMode(mgo.Monotonic, true)
	logrus.Info("Connected to mongodb using mgo")

	sessionMgo = session

	return session, nil
}

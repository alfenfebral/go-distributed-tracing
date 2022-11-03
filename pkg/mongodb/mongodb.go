package mongodb

import (
	"context"
	"os"
	"time"

	"go-distributed-tracing/utils"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.opentelemetry.io/contrib/instrumentation/go.mongodb.org/mongo-driver/mongo/otelmongo"
)

// InitMongoDB - initialize mongo
func InitMongoDB() (context.Context, func(), *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := options.Client()

	// Mongo OpenTelemetry instrumentation
	opts.Monitor = otelmongo.NewMonitor()
	opts.ApplyURI(os.Getenv("DB_URL"))

	client, err := mongo.NewClient(opts)
	if err != nil {
		utils.CaptureError(err)
	}

	err = client.Connect(ctx)
	if err != nil {
		utils.CaptureError(err)
	}

	// Checking the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		utils.CaptureError(err)
	}
	logrus.Println("Database connected")

	return ctx, cancel, client
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/render"
	"github.com/joho/godotenv"
	"github.com/juju/mgosession"
	"github.com/sirupsen/logrus"

	apis "./apis"
	config "./config"
	repository "./repository"
	services "./services"
	"./utils/response"
)

func Routes() *chi.Mux {
	// Sentry
	InitializeSentry()

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	router := chi.NewRouter()
	router.Use(
		sentryHandler.Handle,
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger,                             // Log API request calls
		// middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	return router
}

func InitializeSentry() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: "https://7d1c4f676f3841b2ada881428f7014e3@o469167.ingest.sentry.io/5498121",
	})
	if err != nil {
		logrus.Fatalf("sentry.Init: %s", err)
	}
	// Flush buffered events before the program terminates.
	defer sentry.Flush(2 * time.Second)
}

// PrintAllRoutes - printing all routes
func PrintAllRoutes(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logrus.Panicf("Logging err: %s\n", err.Error()) // panic if there is an error
		sentry.CaptureException(err)
	}
}

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		logrus.Error("Error loading .env file")
		sentry.CaptureException(errors.New("Error loading .env file"))
	}

	// Connect Mgo Database
	session, err := config.InitMgo()
	if err != nil {
		logrus.Error(err)
		sentry.CaptureException(err)

		return
	}
	defer session.Close()

	// Mgo pooling
	configMCP, _ := strconv.Atoi(os.Getenv("MONGODB_CONNECTION_POOL"))
	mPool := mgosession.NewPool(nil, session, configMCP)
	defer mPool.Close()

	// Validator
	// binding.Validator = new(utils.DefaultValidator)

	router := Routes()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, response.H{
			"success": "true",
			"code":    200,
			"message": "Services run properly",
		})
	})

	// Repository
	todoRepo := repository.NewMongoTodoRepository(mPool)

	// Service
	todoService := services.NewTodoService(todoRepo)

	// Handler
	apis.NewTodoHTTPHandler(router, todoService)

	PrintAllRoutes(router)
	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s%s", ":", os.Getenv("PORT")), router)) // Note, the port is usually gotten from the environment.
}

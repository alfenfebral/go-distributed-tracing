package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/riandyrn/otelchi"
	"github.com/sirupsen/logrus"

	"go-distributed-tracing/pkg/config"
	pkg_mongodb "go-distributed-tracing/pkg/mongodb"
	pkg_tracing "go-distributed-tracing/pkg/tracing"
	handlers "go-distributed-tracing/todo/delivery/http"
	repository "go-distributed-tracing/todo/repository"
	services "go-distributed-tracing/todo/services"
	"go-distributed-tracing/utils"
	response "go-distributed-tracing/utils/response"
)

func Routes() *chi.Mux {
	// Sentry
	InitializeSentry()

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	router := chi.NewRouter()
	router.Use(otelchi.Middleware(os.Getenv("APP_NAME"), otelchi.WithChiRoutes(router)))
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
		Dsn: os.Getenv("SENTRY_URL"),
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
		utils.CaptureError(err)
	}
}

func main() {
	utils.InitializeValidator()

	// Load environment variables
	err := config.LoadConfig()
	if err != nil {
		utils.CaptureError(errors.New("error loading .env file"))
	}

	tp := pkg_tracing.InitializeTracing()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Init MongoDB
	_, cancel, client := pkg_mongodb.InitMongoDB()
	defer cancel()

	router := Routes()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, response.H{
			"success": "true",
			"code":    200,
			"message": "Services run properly",
		})
	})

	// Repository
	todoRepo := repository.NewMongoTodoRepository(client)

	// Service
	todoService := services.NewTodoService(todoRepo)

	// Handler
	todoHandler := handlers.NewTodoHTTPHandler(router, tp, todoService)
	todoHandler.RegisterRoutes()

	// Print
	PrintAllRoutes(router)

	logrus.Fatal(http.ListenAndServe(fmt.Sprintf("%s%s", ":", os.Getenv("PORT")), router)) // Note, the port is usually gotten from the environment.
}

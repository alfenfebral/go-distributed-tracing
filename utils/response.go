package utils

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/go-chi/render"

	response "./response"
	"github.com/sirupsen/logrus"
)

func ResponseErrorValidation(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, response.H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": "Validation errors in your request",
		"errors":  ValidatonError(err).Errors,
	})
}

func ResponseBodyError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusBadRequest)
	render.JSON(w, r, response.H{
		"success": false,
		"code":    http.StatusBadRequest,
		"message": "Validation errors in your request",
		"error":   "Check your body request",
	})
}

// ResponseError - send response error (500)
func ResponseError(w http.ResponseWriter, r *http.Request, err error) {
	logrus.Error(err)
	sentry.CaptureException(err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, response.H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": "There is something error",
	})
}

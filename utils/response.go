package utils

import (
	"net/http"

	"github.com/go-chi/render"

	response "go-clean-architecture/utils/response"
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
	CaptureError(err)

	render.Status(r, http.StatusInternalServerError)
	render.JSON(w, r, response.H{
		"success": false,
		"code":    http.StatusInternalServerError,
		"message": "There is something error",
	})
}

// ResponseNotFound - send response not found (404)
func ResponseNotFound(w http.ResponseWriter, r *http.Request, message string) {
	render.Status(r, http.StatusNotFound)
	render.JSON(w, r, response.H{
		"success": false,
		"code":    http.StatusNotFound,
		"message": message,
	})
}

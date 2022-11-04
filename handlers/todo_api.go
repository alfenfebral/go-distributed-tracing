package handlers

import (
	"net/http"

	"go-distributed-tracing/models"
	"go-distributed-tracing/services"
	"go-distributed-tracing/utils"
	response "go-distributed-tracing/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/sdk/trace"
)

// TodoHandler represent the httphandler for file
type TodoHandler struct {
	TodoService services.TodoService
	Tp          *trace.TracerProvider
}

// NewTodoHTTPHandler - make http handler
func NewTodoHTTPHandler(router *chi.Mux, service services.TodoService, tp *trace.TracerProvider) {
	handler := &TodoHandler{
		TodoService: service,
		Tp:          tp,
	}

	router.Get("/todo", handler.GetAll)
	router.Get("/todo/{id}", handler.GetByID)
	router.Post("/todo", handler.Create)
	router.Put("/todo/{id}", handler.Update)
	router.Delete("/todo/{id}", handler.Delete)
}

// GetAll - get all todo http handler
func (handler *TodoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.Tp.Tracer("TodoHandler").Start(r.Context(), "TodoHandler.GetAll")
	defer span.End()

	qQuery := r.URL.Query().Get("q")
	pageQuery := r.URL.Query().Get("page")
	perPageQuery := r.URL.Query().Get("per_page")

	err := utils.ValidateStruct(&models.TodoListRequest{
		Keywords: &models.SearchForm{
			Keywords: qQuery,
		},
		Page:    pageQuery,
		PerPage: perPageQuery,
	})
	if err != nil {
		response.ResponseErrorValidation(w, r, err)
		return
	}

	currentPage := utils.CurrentPage(pageQuery)
	perPage := utils.PerPage(perPageQuery)
	offset := utils.Offset(currentPage, perPage)

	results, totalData, err := handler.TodoService.GetAll(ctx, qQuery, perPage, offset)
	if err != nil {
		response.ResponseError(w, r, err)
		return
	}
	totalPages := utils.TotalPage(totalData, perPage)

	response.ResponseOKList(w, r, &response.ResponseSuccessList{
		Data: results,
		Meta: &response.Meta{
			PerPage:     perPage,
			CurrentPage: currentPage,
			TotalPage:   totalPages,
			TotalData:   totalData,
		},
	})
}

// GetByID - get todo by id http handler
func (handler *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.Tp.Tracer("TodoHandler").Start(r.Context(), "TodoHandler.GetByID")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Get detail
	result, err := handler.TodoService.GetByID(ctx, id)
	if err != nil {
		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: result,
	})

}

// Create - create todo http handler
func (handler *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.Tp.Tracer("TodoHandler").Start(r.Context(), "TodoHandler.Create")
	defer span.End()

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	result, err := handler.TodoService.Create(ctx, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})
	if err != nil {
		response.ResponseError(w, r, err)
		return
	}

	response.ResponseCreated(w, r, &response.ResponseSuccess{
		Data: result,
	})
}

// Update - update instance by id http handler
func (handler *TodoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.Tp.Tracer("TodoHandler").Start(r.Context(), "TodoHandler.Update")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	// Edit data
	_, err := handler.TodoService.Update(ctx, id, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})

	if err != nil {
		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: response.H{
			"id": id,
		},
	})
}

// Delete - delete instance by id http handler
func (handler *TodoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.Tp.Tracer("TodoHandler").Start(r.Context(), "TodoHandler.Delete")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Delete record
	err := handler.TodoService.Delete(ctx, id)
	if err != nil {
		if err.Error() == "not found" {
			response.ResponseNotFound(w, r, "Item not found")
			return
		}

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseOK(w, r, &response.ResponseSuccess{
		Data: response.H{
			"id": id,
		},
	})
}

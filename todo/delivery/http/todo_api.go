package handlers

import (
	"io"
	"net/http"

	"go-distributed-tracing/todo/models"
	"go-distributed-tracing/todo/services"
	"go-distributed-tracing/utils"
	response "go-distributed-tracing/utils/response"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/trace"
)

// todoHandler represent the http handler
type todoHandler struct {
	router      *chi.Mux
	tp          *trace.TracerProvider
	todoService services.TodoService
}

// NewTodoHTTPHandler - make http handler
func NewTodoHTTPHandler(router *chi.Mux, tp *trace.TracerProvider, service services.TodoService) *todoHandler {
	return &todoHandler{
		router:      router,
		tp:          tp,
		todoService: service,
	}
}

func (handler *todoHandler) RegisterRoutes() {
	handler.router.Get("/todo", handler.GetAll)
	handler.router.Get("/todo/{id}", handler.GetByID)
	handler.router.Post("/todo", handler.Create)
	handler.router.Put("/todo/{id}", handler.Update)
	handler.router.Delete("/todo/{id}", handler.Delete)
}

// GetAll - get all todo http handler
func (handler *todoHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tp.Tracer("todoHandler").Start(r.Context(), "todoHandler.GetAll")
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
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

		response.ResponseErrorValidation(w, r, err)
		return
	}

	currentPage := utils.CurrentPage(pageQuery)
	perPage := utils.PerPage(perPageQuery)
	offset := utils.Offset(currentPage, perPage)

	results, totalData, err := handler.todoService.GetAll(ctx, qQuery, perPage, offset)
	if err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

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
func (handler *todoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tp.Tracer("todoHandler").Start(r.Context(), "todoHandler.GetByID")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Get detail
	result, err := handler.todoService.GetByID(ctx, id)
	if err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

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
func (handler *todoHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tp.Tracer("todoHandler").Start(r.Context(), "todoHandler.Create")
	defer span.End()

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

		if err.Error() == io.EOF.Error() {
			response.ResponseBodyError(w, r, err)
			return
		}

		span.SetAttributes(attribute.Bool("validation.error", true))

		response.ResponseErrorValidation(w, r, err)
		return
	}

	result, err := handler.todoService.Create(ctx, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})
	if err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

		response.ResponseError(w, r, err)
		return
	}

	response.ResponseCreated(w, r, &response.ResponseSuccess{
		Data: result,
	})
}

// Update - update instance by id http handler
func (handler *todoHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tp.Tracer("todoHandler").Start(r.Context(), "todoHandler.Update")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	data := &models.TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

		if err.Error() == io.EOF.Error() {
			response.ResponseBodyError(w, r, err)
			return
		}

		response.ResponseErrorValidation(w, r, err)
		return
	}

	// Edit data
	_, err := handler.todoService.Update(ctx, id, &models.Todo{
		Title:       data.Title,
		Description: data.Description,
	})

	if err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

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
func (handler *todoHandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx, span := handler.tp.Tracer("todoHandler").Start(r.Context(), "todoHandler.Delete")
	defer span.End()

	// Get and filter id param
	id := chi.URLParam(r, "id")

	// Delete record
	err := handler.todoService.Delete(ctx, id)
	if err != nil {
		span.SetAttributes(
			attribute.Key("error").Bool(true),
		)
		span.RecordError(err)

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

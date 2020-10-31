package apis

import (
	"net/http"

	"../services"
	"../utils"
	response "../utils/response"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"github.com/go-playground/validator"
	"gopkg.in/mgo.v2/bson"
)

var validate *validator.Validate

// TodoRequest - todo request
type TodoRequest struct {
	Title       string `form:"title" json:"title" validate:"required"`
	Description string `form:"description" json:"description" validate:"required"`
}

func (a *TodoRequest) Bind(r *http.Request) error {
	validate = validator.New()
	err := validate.Struct(a)
	if err != nil {

		// this check is only needed when your code could produce
		// an invalid value for validation such as interface with nil
		// value most including myself do not usually have code like this.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		return err
	}

	return err
}

// Todohandler represent the httphandler for file
type Todohandler struct {
	TodoService services.TodoService
}

// NewTodoHTTPHandler - make http handler
func NewTodoHTTPHandler(router *chi.Mux, service services.TodoService) {
	handler := &Todohandler{
		TodoService: service,
	}

	router.Post("/todo", handler.Create)
	// r.GET("/file", handler.GetAll)
	// r.GET("/file/:id", handler.GetByID)
	// r.DELETE("/file/:id", handler.Delete)
}

// Create - create todo http handler
func (handler *Todohandler) Create(w http.ResponseWriter, r *http.Request) {
	data := &TodoRequest{}
	if err := render.Bind(r, data); err != nil {
		if err.Error() == "EOF" {
			utils.ResponseBodyError(w, r, err)
			return
		}

		utils.ResponseErrorValidation(w, r, err)
		return
	}
	timeNow := utils.GetTimeNow()

	result, err := handler.TodoService.Create(bson.M{
		"title":       data.Title,
		"description": data.Description,
		"createdAt":   timeNow,
		"updatedAt":   timeNow,
		"deletedAt":   timeNow,
	})
	if err != nil {
		utils.ResponseError(w, r, err)
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, response.H{
		"success": true,
		"code":    http.StatusCreated,
		"message": "Create Todo",
		"data":    result,
	})
}

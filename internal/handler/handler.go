package handler

import (
	"Gonoty/internal/handler/dto"
	"Gonoty/internal/models"
	"Gonoty/internal/storage"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TaskHandler struct {
	storage storage.Storage
}

func NewTaskHandler(storage storage.Storage) *TaskHandler {
	return &TaskHandler{
		storage: storage,
	}
}

func (h *TaskHandler) SendEmail(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	data := &dto.SendEmailRequest{}
	if err := render.Bind(r, data); err != nil {
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	task, err := CreateSendTask(r.Context(), data)
	if err != nil {
		//render.Render(w, r, ErrInternalServerError())
		return
	}

	err = h.storage.Add(ctx, task)
	if err != nil {
		log.Println(err)
		return
	}

	render.Render(w, r, dto.NewSendEmailResponse(
		task.ID,
		"pending",
		"Task accepted for processing",
	))
}

type ErrResponse struct {
	Err            error `json:"-"`
	HTTPStatusCode int   `json:"-"`

	StatusText string `json:"status"`
	AppCode    int64  `json:"code,omitempty"`
	ErrorText  string `json:"error,omitempty"`
}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 400,
		StatusText:     "Invalid request",
		ErrorText:      err.Error(),
	}
}

func CreateSendTask(ctx context.Context, req *dto.SendEmailRequest) (models.Task, error) {
	fromEmail := req.FromEmail
	if fromEmail == "" {
		fromEmail = "default@myapp.com" // get from config
	}

	task := models.Task{
		ID:         uuid.New().String(),
		Recipients: req.Recipients,
		Subject:    req.Subject,
		Body:       req.Body,
		FromEmail:  fromEmail,
		Status:     models.StatusPending,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return task, nil
}

package handler

import (
	"Gonoty/internal/handler/dto"
	"Gonoty/internal/models"
	"Gonoty/internal/queue"
	"context"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
)

type TaskHandler struct {
	ctx context.Context
	q   queue.Queue
}

func NewTaskHandler(context context.Context, queue queue.Queue) *TaskHandler {
	return &TaskHandler{
		ctx: context,
		q:   queue,
	}
}

func (h *TaskHandler) AddToQueue(w http.ResponseWriter, r *http.Request) {
	data := &dto.SendEmailRequest{}
	if err := render.Bind(r, data); err != nil {
		log.Printf("Invalid request: %v", err)
		render.Render(w, r, ErrInvalidRequest(err))
		return
	}

	task, err := CreateTask(r.Context(), data)
	if err != nil {
		log.Printf("Failed to create task: %v", err)
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	err = h.q.Enqueue(r.Context(), task)
	if err != nil {
		log.Printf("Failed to enqueue task %s: %v", task.ID, err)
		render.Render(w, r, ErrInternalServerError(err))
		return
	}

	log.Printf("Task %s enqueued", task.ID)

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
	if e.Err != nil {
		requestID := middleware.GetReqID(r.Context())
		log.Printf("[%s] HTTP %d: %v", requestID, e.HTTPStatusCode, e.Err)
	}

	render.Status(r, e.HTTPStatusCode)
	return nil
}

func ErrInvalidRequest(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Invalid request",
		ErrorText:      "Invalid request data",
	}
}

func ErrInternalServerError(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: http.StatusBadRequest,
		StatusText:     "Internal server error",
		ErrorText:      "Something went wrong",
	}
}

func CreateTask(ctx context.Context, req *dto.SendEmailRequest) (models.Task, error) {
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

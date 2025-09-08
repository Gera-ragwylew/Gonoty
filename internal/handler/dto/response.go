package dto

import (
	"net/http"

	"github.com/go-chi/render"
)

// SendEmailResponse DTO для ответа
type SendEmailResponse struct {
	TaskID  string `json:"task_id"`
	Status  string `json:"status"`
	Message string `json:"message,omitempty"`
}

func NewSendEmailResponse(taskID, status, message string) *SendEmailResponse {
	return &SendEmailResponse{
		TaskID:  taskID,
		Status:  status,
		Message: message,
	}
}

func (s *SendEmailResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, http.StatusAccepted)
	return nil
}

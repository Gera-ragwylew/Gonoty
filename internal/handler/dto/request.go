package dto

import (
	"Gonoty/internal/models"
	"errors"
	"net/http"
)

type SendEmailRequest struct {
	Recipients []models.Recipient `json:"recipients"`
	Subject    string             `json:"subject"`
	Body       models.EmailBody   `json:"body"`
	FromEmail  string             `json:"from_email,omitempty"`
}

func (r *SendEmailRequest) Bind(*http.Request) error {
	if len(r.Recipients) == 0 {
		// return render.ValidationError{Field: "recipients", Reason: "at least one recipient is required"}
		return errors.New("at least one recipient is required")
	}
	for _, recipient := range r.Recipients {
		if recipient.Email == "" {
			// return render.ValidationError{Field: "recipients", Reason: "email is required"}
			return errors.New("recipients email is required")
		}
	}
	if r.Subject == "" {
		// return render.ValidationError{Field: "subject", Reason: "is required"}
		return errors.New("subject is required")
	}
	if r.Body.Text == "" && r.Body.HTML == "" {
		// return render.ValidationError{Field: "body", Reason: "text or html content is required"}
		return errors.New("body: text or html content is required")
	}
	return nil
}

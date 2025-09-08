package models

import "time"

type TaskStatus string

const (
	StatusPending   TaskStatus = "pending"
	StatusProcessed TaskStatus = "processed"
	StatusFailed    TaskStatus = "failed"
	StatusCompleted TaskStatus = "completed"
)

type Recipient struct {
	Email string `json:"email" db:"email"`
}

type EmailBody struct {
	Text string `json:"text,omitempty" db:"text"`
	HTML string `json:"html,omitempty" db:"html"`
}

type Task struct {
	ID         string      `json:"id" db:"id"`
	Recipients []Recipient `json:"recipients" db:"recipients"`
	Subject    string      `json:"subject" db:"subject"`
	Body       EmailBody   `json:"body" db:"body"`
	FromEmail  string      `json:"from_email" db:"from_email"` // Важно для DMARC/SPF
	Status     TaskStatus  `json:"status" db:"status"`
	CreatedAt  time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" db:"updated_at"`
}

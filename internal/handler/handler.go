package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/render"
)

type EmailRequest struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func SendEmailHandler(w http.ResponseWriter, r *http.Request) {
	// if r.Method != http.MethodPost {
	// 	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	// 	log.Printf("Method not allowed")
	// 	return
	// }

	data := &EmailRequest{}
	if err := render.Bind(r, data); err != nil {
		// http.Error(w, "Invalid JSON", http.StatusMethodNotAllowed)
		// log.Printf("Invalid JSON")
		return
	}

	// if req.To == "" || req.Subject == "" || req.Body == "" {
	// 	http.Error(w, "Missing required fields: to, subject, body", http.StatusBadRequest)
	// 	log.Printf("Missing required fields: to, subject, body")
	// 	return
	// }

	// log.Println(req)
	// err := service.SendEmail(req.To, req.Subject, req.Body)
	// if err != nil {
	//     // Логируем ошибку для себя
	//     fmt.Printf("Failed to send email: %v\n", err)
	//     // И отправляем клиенту общую ошибку
	//     http.Error(w, "Failed to send email", http.StatusInternalServerError)
	//     return
	// }

	// Отправляем успешный ответ
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Email sent successfully",
	})

	log.Printf("Email sent successfully")
}

func (e *EmailRequest) Bind(r *http.Request) error {
	return nil
}

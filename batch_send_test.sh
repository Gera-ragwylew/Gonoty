#!/bin/bash

# Количество email для генерации
COUNT=1000

# Генерируем случайные email
recipients_json=""
for ((i=1; i<=COUNT; i++)); do
    # Генерируем случайное имя пользователя
    username="user$(shuf -i 1000-9999 -n 1)"
    
    if [ -n "$recipients_json" ]; then
        recipients_json="$recipients_json,"
    fi
    recipients_json="$recipients_json{\"email\": \"$username@example.com\"}"
done

echo "Отправка $COUNT писем..."

curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": ['"$recipients_json"'],
    "subject": "Test Subject - Batch Send",
    "body": {
      "text": "This is a test email content in plain text.",
      "html": "<p>This is a test email content in <b>HTML</b>.</p>"
    },
    "from_email": "noreply@myapp.com"
  }'

echo "Запрос отправлен!"
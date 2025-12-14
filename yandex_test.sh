#!/bin/bash

echo "Отправка yandex..."

curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "test@yandex.ru"}
    ],
    "subject": "Test Subject - Yandex Test",
    "body": {
      "text": "This is a test email content in plain text.",
      "html": "<p>This is a test email content in <b>HTML</b>.</p>"
    },
    "from_email": "test@yandex.ru"
  }'

echo "Запрос отправлен!"
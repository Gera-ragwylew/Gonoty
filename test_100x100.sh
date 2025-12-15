#!/bin/bash

# Отправляем 100 параллельных запросов по 100 email
for i in {1..100}; do
    (
        recipients=""
        for j in {1..100}; do
            num=$(( (i-1)*100 + j ))
            if [ $j -gt 1 ]; then
                recipients="$recipients,"
            fi
            recipients="$recipients{\"email\": \"test${num}@example.com\"}"
        done
        
        curl -X POST http://localhost:8080/send \
          -H "Content-Type: application/json" \
          -d '{
            "recipients": ['"$recipients"'],
            "subject": "Batch '"$i"'",
            "body": {
              "text": "Test email.",
              "html": "<p>Test email.</p>"
            },
            "from_email": "noreply@myapp.com"
          }' &
    ) &
done

wait
echo "10,000 emails отправлены в 100 параллельных запросах!"
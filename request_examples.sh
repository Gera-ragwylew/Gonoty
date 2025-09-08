# succsess 1 recipient
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "test1@example.com"}
    ],
    "subject": "Test Subject",
    "body": {
      "text": "This is a test email content in plain text.",
      "html": "<p>This is a test email content in <b>HTML</b>.</p>"
    },
    "from_email": "noreply@myapp.com"
  }' \

# succsess many recipient
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "user1@example.com"},
      {"email": "user2@example.com"},
      {"email": "user3@example.com"}
    ],
    "subject": "Weekly Newsletter",
    "body": {
      "text": "Hello! Here is your weekly update...",
      "html": "<h1>Hello!</h1><p>Here is your weekly update...</p>"
    }
  }' \

# failed no recipient
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [],
    "subject": "Test Subject",
    "body": {
      "text": "Content"
    }
  }' \

# failed empty email
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": ""}
    ],
    "subject": "Test Subject",
    "body": {
      "text": "Content"
    }
  }' \

# failed no subject
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "test@example.com"}
    ],
    "subject": "",
    "body": {
      "text": "Content"
    }
  }' \

# failed no body
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "test@example.com"}
    ],
    "subject": "Test Subject",
    "body": {}
  }' \

#failed error email format
curl -X POST http://localhost:8080/send \
  -H "Content-Type: application/json" \
  -d '{
    "recipients": [
      {"email": "invalid-email"}
    ],
    "subject": "Test Subject",
    "body": {
      "text": "Content"
    }
  }' \

#failed error content type
curl -X POST http://localhost:8080/send \
  -H "Content-Type: text/plain" \
  -d "plain text data" \

#failed error content type
curl -X GET http://localhost:8080/send \
  -H "Content-Type: application/json" \

#failed error endpoint
curl -X POST http://localhost:8080/nonexistent \
  -H "Content-Type: application/json" \
  -d '{"test": "data"}' \

curl -X POST http://localhost:8080/internal/create-license \
     -H "Content-Type: application/json" \
     -d '{"email": "user@example.com", "key": "ABCD-1234-EFGH-5678"}'

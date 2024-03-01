<!-- Auth Signup -->
curl -X POST http://localhost:8080/signup -H "Content-Type: application/json" -d '{"username": "testuser", "email": "your-email.com", "password": "admin1"}'

<!-- Auth Signin -->
curl -X POST http://localhost:8080/signin      -H "Content-Type: application/json"      -d '{"email":"zinko.adrian00@gmail.com", "password":"admintest"}'

<!-- Resources -->
curl -X GET http://localhost:8080/resources

curl -X GET http://localhost:8080/resources/1

curl -X POST http://localhost:8080/resources \
-H "Content-Type: application/json" \
-d '{"title":"New Resource", "category":"Programming", "description":"A new resource description", "url":"http://example.com", "resource_type":"Article", "completion_time":"1 hour"}'

curl -X PUT http://localhost:8080/resources/1 \
-H "Content-Type: application/json" \
-d '{"title":"Updated Resource", "category":"Math", "description":"Updated description", "url":"http://example.com", "resource_type":"Video", "completion_time":"2 hours"}'

curl -X DELETE http://localhost:8080/resources/1
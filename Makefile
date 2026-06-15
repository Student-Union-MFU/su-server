.PHONY: dev build up down logs test-health test-events

# Docker
up:
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f backend

# Dev
dev:
	go run cmd/main.go

build:
	go build -o bin/main cmd/main.go

# Routes
test-health:
	curl -s http://localhost:8080/ | jq

test-get-events:
	curl -s http://localhost:8080/su-server/events | jq

test-create-event:
	curl -s -X POST http://localhost:8080/su-server/events \
			-H "Content-Type: application/json" \
			-d "{\"title\":\"Student Night Market\",\"content\":\"A gathering place\",\"location\":\"M-Square Rooftop\",\"date\":\"27 March 2026\",\"time\":\"4:00 PM - 9:00 PM\",\"link\":\"#\"}" | jq

test-login:
	curl -s -X POST http://localhost:8080/auth/login \
		-H "Content-Type: application/json" \
		-d '{"username": "admin", "password": "yion"}' | jq

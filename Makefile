include .env
export

.PHONY: dev build up down logs test-health test-events

print-env:
	@echo "USER: $(DB_USER)"
	@echo "PASS: $(DB_PASS)"
	@echo "DB: $(DB_NAME)"

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

watch:
	air

build:
	go build -o bin/main cmd/main.go


# Routes

get-events:
	http GET $(BASE_URL)/events
 
get-event:
	http GET $(BASE_URL)/events/$(id)
 
create-event:
	http POST $(BASE_URL)/events \
		title="$(title)" \
		content="$(content)" \
		location="$(location)" \
		date="$(date)" \
		time="$(time)" \
		link="$(link)"

update-event:
	http PUT $(BASE_URL)/events/$(id) \
		title=$(title) \
		content=$(content)
 
delete-event:
	http DELETE $(BASE_URL)/events/$(id)

create-user:
	http POST $(BASE_URL)/users \
		name="$(name)" \
		email="$(email)" \
		usertype="$(usertype)" \
		student_id="$(student_id)" \
		major="$(major)" \
		school="$(school)" \
		avatar_url="$(avatar_url)" \
		oauth_subject="$(oauth_subject)"

get-user:
	http GET $(BASE_URL)/users/$(id)

get-user-email:
	http GET $(BASE_URL)/users/email/$(email)

update-user:
	http PATCH $(BASE_URL)/users/$(id) \
		major="$(major)" \
		school="$(school)" \
		student_id="$(student_id)


# migrations

migrate-up:
	migrate -path $(DB_MIGRATION_PATH) -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path $(DB_MIGRATION_PATH)/$(FILE) -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" down

migrate-drop:
	migrate $(DB_MIGRATION_PATH)/$(FILE) -database "postgres://$(DB_USER):$(DB_PASS)@localhost:5432/$(DB_NAME)?sslmode=disable" drop

tables:
	docker exec -it postgres-db psql -U $(DB_USER) -d $(DB_NAME) -c "\dt"

check-schema-migrations:
	docker exec -it postgres-db psql -U $(DB_USER) -d $(DB_NAME) -c "SELECT * FROM schema_migrations;" \dt

build:
	docker-compose build note-app

run:
	docker-compose up note-app

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@0.0.0.0:5432/postgres?sslmode=disable' up
	
swag:
	swag init -g cmd/main.go

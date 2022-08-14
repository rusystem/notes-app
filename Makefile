build:
	docker-compose build note-app

run:
	docker-compose up note-app

migrate:
	migrate -path ./schema -database 'postgres://postgres:qwerty@127.0.0.1:5432/postgres?sslmode=disable' up
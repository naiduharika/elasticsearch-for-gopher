composeup:
	docker compose up -d

composedown:
	docker compose down -v

run:
	go run main.go

.PHONY: composeup composedown
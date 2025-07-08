start:
	docker compose down --volumes && docker compose up -d --build

stop:
	docker compose down --volumes
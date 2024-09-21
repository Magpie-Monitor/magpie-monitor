watch:
	docker compose  -f docker-compose.yml -f docker-compose.dev.yml watch

logs:
	docker compose -f docker-compose.yml logs -f

restart:
	docker compose down && docker compose  -f docker-compose.yml -f docker-compose.dev.yml watch

down:
	docker compose down
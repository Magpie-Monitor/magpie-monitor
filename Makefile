watch:
	docker compose down && docker compose  -f docker-compose.yml -f docker-compose.dev.yml watch

elastic:
	docker compose -f docker-compose.es.yml up -d

logs:
	docker compose -f docker-compose.yml logs -f

restart:
	docker compose down && docker compose  -f docker-compose.yml -f docker-compose.dev.yml watch

down:
	docker compose down


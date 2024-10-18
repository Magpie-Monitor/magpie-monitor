all:
	docker compose down && docker compose  -f docker-compose.yml -f docker-compose.dev.yml -f docker-compose.es.yml up --watch  --build

watch:
	docker compose down && docker compose  -f docker-compose.yml -f docker-compose.dev.yml up --watch  --build

elastic:
	docker compose -f docker-compose.es.yml up -d

logs:
	docker compose -f docker-compose.yml logs -f

restart:
	docker compose down && docker compose  -f docker-compose.yml -f docker-compose.dev.yml watch

down:
	docker compose down


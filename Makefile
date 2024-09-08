watch:
	docker compose  -f docker-compose.yml -f docker-compose.dev.yml watch

eks:
	docker compose -f docker-compose.elk.yml up

logs:
	docker compose logs -f


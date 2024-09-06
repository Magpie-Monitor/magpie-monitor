watch:
	docker compose -f docker-compose.yml -f docker-compose.dev.yml watch & docker compose logs -f

logs:
	docker compose logs -f


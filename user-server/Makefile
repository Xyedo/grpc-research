MG_PATH=./migrations

docker-up:
	docker compose --env-file .env up -d
docker-down:
	docker compose down -t 1
migrate-up:
	docker exec db-migrator-user alembic upgrade head
migrate-revision:
	docker exec db-migrator-user alembic revision --autogenerate -m "$(message)"
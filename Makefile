.PHONY: build start logs

build:
	docker compose build app
start:
	docker compose up -d
restart:
	docker restart my-app
stop: #down
	@echo "=============Cleaning up============="
	docker compose down
	docker system prune -f
	docker volume prune -f
logs:
	docker logs -f my-app

ssh:
	docker exec -it my-app bash
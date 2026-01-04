.PHONY: build start logs swagger

build:
	docker compose build app
start:
	docker compose up -d
restart:
	docker restart nail-project
restart-logs:
	$(MAKE) start \
	${MAKE} logs
stop: #down
	@echo "=============Cleaning up============="
	docker compose down
	docker system prune -f
	docker volume prune -f
logs:
	docker logs -f nail-project
ssh:
	docker exec -it nail-project bash
swagger:
	@echo "=============Generating Swagger documentation============="
	swag init -g cmd/api/main.go --output docs
	@echo "=============Swagger documentation generated============="
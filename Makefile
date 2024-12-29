run: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
run:
	@docker compose up

build: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
build:
	@docker compose up --build

deploy: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
deploy:
	docker compose -f compose.prod.yaml up --build -d

status:
	@docker exec -it alcs-server goose postgres "$(shell docker exec -it alcs-server /bin/manage -c print-db-string)" -dir=internal/db/migrations status

up:
	@docker exec -it alcs-server goose postgres "$(shell docker exec -it alcs-server /bin/manage -c print-db-string)" -dir=internal/db/migrations up

reset:
	@docker exec -it alcs-server goose postgres "$(shell docker exec -it alcs-server /bin/manage -c print-db-string)" -dir=internal/db/migrations reset

templ:
	@templ generate

pkl_gen:
	@pkl-gen-go pkl/config.pkl --base-path github.com/tikhonp/alcs

create-super-admin:
	@docker exec -it alcs-server /bin/manage -c create-super-admin


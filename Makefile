run: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
run:
	@docker compose up --build

deploy: export SOURCE_COMMIT=$(shell git rev-parse HEAD)
deploy:
	docker compose -f compose.prod.yaml up --build -d

status:
	@goose postgres "$(shell /bin/manage -c print-db-string)" -dir=internal/db/migrations status

up:
	@goose postgres "$(shell /bin/manage -c print-db-string)" -dir=internal/db/migrations up

reset:
	@goose postgres "$(shell /bin/manage -c print-db-string)" -dir=internal/db/migrations reset

templ:
	@templ generate

pkl_gen:
	@pkl-gen-go pkl/config.pkl --base-path github.com/tikhonp/alcs


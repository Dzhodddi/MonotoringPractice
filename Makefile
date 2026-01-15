export COMPOSE_BAKE=true
.PHONY: up down
#tidy gen_docs lint

up:
	docker compose up --build

down:
	docker compose down

tidy:
	cd account && go mod tidy
	cd graphql && go mod tidy
	cd order && go mod tidy
	cd payment && go mod tidy
	cd pkg && go mod tidy
	cd product && go mod tidy

#
#gen_docs:
#	cd gateway && make docs
#	cd ingester && make docs
#	cd processor && make docs
#	cd reporter && make docs
#
#lint:
#	cd gateway && make lint
#	cd ingester && make lint
#	cd processor && make lint
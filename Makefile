.PHONY: keypair migrate-create migrate-up migrate-down migrate-force

PWD = $(shell pwd)
ACCOUNTPATH = $(PWD)/account
MPATH = $(ACCOUNTPATH)/migrations
PORT = 5432

# Default number of migrations to execute up or down.
N = 1

create-keypair:
	@echo "Creating an rsa 256 key pair"
	openssl genpkey -algorithm RSA -out $(ACCOUNTPATH)/rsa_private_$(ENV).pem -pkeyopt rsa_keygen_bits:2048
	openssl rsa -in $(ACCOUNTPATH)/rsa_private_$(ENV).pem -pubout -out $(ACCOUNTPATH)/rsa_public_$(ENV).pem


migrate-create:
	@echo "---Creating migration files---"
	migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)

migrate-up:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable up $(N)

migrate-down:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable down $(N)

migrate-force:
	migrate -source file://$(MPATH) -database postgres://postgres:password@localhost:$(PORT)/postgres?sslmode=disable force $(VERSION)

# create dev and test keys.
# run postgres containers in docker-compose.
# migrate down.
# migrate up.
# docker-compose down.
init:
	docker-compose up -d postgres-account && \
	$(MAKE) create-keypair ENV=dev && \
	$(MAKE) create-keypair ENV=test && \
	$(MAKE) migrate-down ACCOUNTPATH=account N= && \
	$(MAKE) migrate-up ACCOUNTPATH=account N= && \
	docker-compose down
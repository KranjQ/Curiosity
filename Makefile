SERVER_BINARY=graphql
DOCKER_DIR=docker

STORAGE_TYPE?=PSQL

build-graphql:
	go build -o ${SERVER_BINARY} ./cmd/main


.PHONY: builder-image
builder-image:
	docker build -t builder -f ${DOCKER_DIR}/builder.Dockerfile .

.PHONY: graphql-service-image
graphql-service-image:
	docker build -t graphql-service -f ${DOCKER_DIR}/graphql.Dockerfile .

.PHONY: curiosity-run
curiosity-run:
	make builder-image
	make graphql-service-image
	STORAGE_TYPE=${STORAGE_TYPE} docker-compose -f ${DOCKER_DIR}/docker-compose.yml up -d

.PHONY: curiosity-down
curiosity-down:
	docker-compose -f ${DOCKER_DIR}/docker-compose.yml down


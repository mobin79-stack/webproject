SERVICE_NAME=template
DOCKER_IMAGE_NAME=registry.ecobin.ir/ecomicro/${SERVICE_NAME}
BINARY=${SERVICE_NAME}.o

clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

build: clean
	go build -o ${BINARY} main.go

compose:
	docker compose up -d

compose-dev:
	docker compose -f docker-compose-dev.yml up -d

compose-test:
	docker compose -f docker-compose-test.yml up -d

run-foo:
	go run ./extra/foo-main.go

run: build
	$(build)
	./${BINARY}

watch: compose-dev 
	sleep 1
	air

docker:
	docker build --build-arg USERNAME=$(USERNAME) --build-arg APIKEY=$(APIKEY) -t ${DOCKER_IMAGE_NAME}:$(VERSION) .

docker-push:
	docker push ${DOCKER_IMAGE_NAME}:$(VERSION)

test: compose-test 
	sleep 1
	go test ./...


.PHONY: run
run:
	go run cmd/facade/main.go


.PHONY: build
build: .build

.build:
	go build -o ./bin/facade$(shell go env GOEXE) ./cmd/facade/main.go


.PHONY: docker
docker:
	docker build -t hablof/logistic-package-api-facade .

.PHONY: docker-run
docker-run:
	docker run -d  \
		-v $(PWD)/config.yml:/root/config.yml \
		hablof/logistic-package-api-facade
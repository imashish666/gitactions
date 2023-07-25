.PHONY: build, test

docker-push:
	sudo docker build -f prod.dockerfile --build-arg env=rtqa7 -t www-api:latest . && sudo docker tag www-api:latest 280563394466.dkr.ecr.us-west-1.amazonaws.com/rtqa7-www-api:latest && sudo docker push 280563394466.dkr.ecr.us-west-1.amazonaws.com/rtqa7-www-api:latest
	
docker-build:
	docker build -t www-api:latest .

docker-run:
	docker run --rm -it -p 8080:8080/tcp www-api:latest

test-report:
	go test ./... -coverprofile=cover.out && go tool cover -html=cover.out

test:
	go test ./... --cover

go-build:
	CGO_ENABLED=0 GOOS=linux go build main.go && rm main

vet:
	go fmt ./... && go vet ./...
	
doc:
	swag init

local-build:
	docker compose build

local-run:
	docker compose up -d
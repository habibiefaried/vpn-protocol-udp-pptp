integration-test:
	go fmt ./...
	go build .
	docker-compose down --rmi all
	docker-compose build
	docker-compose up -d
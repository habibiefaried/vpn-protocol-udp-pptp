integration-test:
	go fmt ./...
	go build .
	go test -parallel 4 -v integration_test.go
	go test -parallel 4 -v encrypt.go encrypt_test.go
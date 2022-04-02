all:
	go build server.go
	docker rm -f vpn
	docker build . -f Dockerfile.test -t vpn-protocol
	docker run --privileged --network=host --name vpn -dit vpn-protocol

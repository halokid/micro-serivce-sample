build-docker:
	GOOS=linux GOARCH=amd64 go build
	docker build -t consignment-cli .

run-docker:
	docker run -e MICRO_REGISTRY=mdns consignment-cli



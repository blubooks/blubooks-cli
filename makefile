help:  ## show help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[$$()% 0-9a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)


run-server:
	gin --appPort 4080 --port 3000 --all run main.go 

run-client:
	yarn --cwd webapp/default build --watch

protoc:
	yarn --cwd webapp/default run protoc
	protoc --go_out=internal/app   --go-grpc_out=internal/app  proto/person.proto


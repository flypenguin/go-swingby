# https://medium.com/@chaewonkong/simplifying-your-build-process-with-makefiles-in-golang-projects-b125af7a10c4
# https://earthly.dev/blog/golang-makefile/

.DEFAULT_GOAL := run
TARGET := swingby
GOFLAGS := 


$(TARGET): main.go
	go build $(GOFLAGS) -o $(TARGET)


$(TARGET)-linux-amd64: main.go
	export GOOS=linux ; \
	export GOARCH=amd64 ; \
	go build $(GOFLAGS) -o $(TARGET)-$$GOOS-$$GOARCH


$(TARGET)-darwin-arm64: main.go
	export GOOS=darwin ; \
	export GOARCH=arm64 ; \
	go build $(GOFLAGS) -o $(TARGET)-$$GOOS-$$GOARCH


run:
	go run main.go
.PHONY: run


build: $(TARGET)
build: $(TARGET)-linux-amd64
build: $(TARGET)-darwin-arm64
.PHONY: build


container-build: clean
container-build:
	podman build -t swingby .
	podman tag swingby docker.io/flypenguin/swingby
.PHONY: container-build


container-push:
	podman push docker.io/flypenguin/swingby
.PHONY: container-push


container: container-build container-push
.PHONY: container


clean:
	rm -f $(TARGET)
	rm -f $(TARGET)-darwin-arm64
	rm -f $(TARGET)-linux-amd64
.PHONY: clean

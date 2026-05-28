default: fmt lint test build

.PHONY: default fmt lint test build clean

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

build:
	CGO_ENABLED=0 go build -trimpath -o nullcloud-backend .

clean:
	rm -f nullcloud-backend

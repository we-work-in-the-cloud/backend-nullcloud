default: fmt lint test build

.PHONY: default fmt lint test build ui-build clean

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

ui-build:
	mkdir -p internal/api/ui-build
	cp internal/ui/index.html internal/api/ui-build/index.html
	cp internal/ui/style.css internal/api/ui-build/style.css
	cp internal/ui/app.js internal/api/ui-build/app.js

build: ui-build
	CGO_ENABLED=0 go build -trimpath -o nullcloud-backend .

clean:
	rm -f nullcloud-backend
	rm -rf internal/api/ui-build/

default: fmt lint test build

.PHONY: default fmt lint test build ui-build clean

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

test: ui-build
	go test -v -cover -coverprofile=coverage.out -timeout=120s -parallel=10 ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

ui-build:
	cd internal/ui && npm ci && npm run build
	mkdir -p internal/api/ui-build
	cp internal/ui/build/index.html internal/api/ui-build/index.html
	cp internal/ui/build/app.js internal/api/ui-build/app.js
	touch internal/api/ui-build/style.css

build: ui-build
	CGO_ENABLED=0 go build -trimpath -o nullcloud-backend .

clean:
	rm -f nullcloud-backend
	rm -rf internal/api/ui-build/

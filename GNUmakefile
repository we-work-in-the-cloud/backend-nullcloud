default: fmt lint test build

.PHONY: default fmt lint test build ui-build clean

fmt:
	gofmt -s -w -e .

lint:
	golangci-lint run

test:
	go test -v -cover -timeout=120s -parallel=10 ./...

ui-build:
	cd internal/ui && npm run build
	mkdir -p internal/api/ui-build
	cp internal/ui/build/index.html internal/api/ui-build/index.html
	cp internal/ui/build/app.js internal/api/ui-build/app.js
	touch internal/api/ui-build/style.css

build: ui-build
	CGO_ENABLED=0 go build -trimpath -o nullcloud-backend .

clean:
	rm -f nullcloud-backend
	rm -rf internal/api/ui-build/

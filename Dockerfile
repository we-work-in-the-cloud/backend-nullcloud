# Build the web UI assets embedded in the backend binary.
FROM node:22-alpine AS ui
WORKDIR /src/internal/ui
COPY internal/ui/package*.json ./
RUN npm ci
COPY internal/ui/ ./
RUN npm run build

# Compile a static Linux binary for the minimal runtime image.
FROM golang:1.27.1-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=ui /src/internal/ui/build/ ./internal/api/ui-build/
RUN touch internal/api/ui-build/style.css
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/nullcloud-backend .

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /out/nullcloud-backend /nullcloud-backend
EXPOSE 8080
VOLUME ["/data"]
USER nonroot:nonroot
ENTRYPOINT ["/nullcloud-backend"]
CMD ["--host", "0.0.0.0", "--port", "8080", "--store-file", "/data/store.json"]

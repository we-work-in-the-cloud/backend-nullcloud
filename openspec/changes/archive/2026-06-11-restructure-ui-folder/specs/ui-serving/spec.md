## ADDED Requirements

### Requirement: UI is served at /ui endpoint

The system SHALL serve the UI HTML at the `/ui/` endpoint. The endpoint behavior remains unchanged; only the source location of the embedded assets is updated.

#### Scenario: GET /ui/ returns HTML

- **WHEN** client sends `GET /ui/`
- **THEN** the server responds with HTTP 200
- **AND** the response body contains the UI HTML
- **AND** the Content-Type header is `text/html; charset=utf-8`

#### Scenario: GET /ui/style.css returns stylesheet

- **WHEN** client sends `GET /ui/style.css`
- **THEN** the server responds with HTTP 200
- **AND** the response body contains the UI CSS
- **AND** the Content-Type header is `text/css; charset=utf-8`

#### Scenario: GET /ui/app.js returns JavaScript

- **WHEN** client sends `GET /ui/app.js`
- **THEN** the server responds with HTTP 200
- **AND** the response body contains the UI JavaScript
- **AND** the Content-Type header is `application/javascript; charset=utf-8`

### Requirement: UI assets are embedded from new location

The system SHALL embed UI assets from `internal/ui/build/` instead of `internal/api/`.

#### Scenario: Embed directive points to new path

- **WHEN** the Go source code is compiled
- **THEN** the `//go:embed` directives in `internal/api/ui.go` reference `internal/ui/build/` files
- **AND** the compiled binary contains the assets from the new location

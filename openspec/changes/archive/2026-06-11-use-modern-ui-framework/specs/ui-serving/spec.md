## MODIFIED Requirements

### Requirement: UI is served at /ui endpoint

The system SHALL serve the UI HTML at the `/ui/` endpoint. The endpoint behavior remains unchanged; the UI is now built with React instead of vanilla JavaScript, but the HTTP interface is identical.

#### Scenario: GET /ui/ returns HTML

- **WHEN** client sends `GET /ui/`
- **THEN** the server responds with HTTP 200
- **AND** the response body contains the React-built UI HTML
- **AND** the Content-Type header is `text/html; charset=utf-8`

#### Scenario: GET /ui/style.css returns stylesheet

- **WHEN** client sends `GET /ui/style.css`
- **THEN** the server responds with HTTP 200
- **AND** the response body contains the compiled CSS
- **AND** the Content-Type header is `text/css; charset=utf-8`

#### Scenario: GET /ui/app.js returns JavaScript

- **WHEN** client sends `GET /ui/app.js`
- **THEN** the server responds with HTTP 200
- **AND** the response body contains the compiled React application
- **AND** the Content-Type header is `application/javascript; charset=utf-8`

### Requirement: UI assets are embedded from Vite build output

The system SHALL embed UI assets produced by Vite (npm build) from `internal/ui/build/`. The Vite output is optimized and minified before embedding.

#### Scenario: Embed directive points to Vite output

- **WHEN** the Go source code is compiled
- **THEN** the `//go:embed` directives in `internal/api/ui.go` reference files in `internal/api/ui-build/`
- **AND** those files are copied from `internal/ui/build/` by the Makefile before Go compilation
- **AND** the compiled binary contains the assets from the new location

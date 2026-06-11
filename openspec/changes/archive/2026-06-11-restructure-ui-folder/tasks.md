## 1. Folder Structure & Git Setup

- [x] 1.1 Create `internal/ui/` directory
- [x] 1.2 Create `internal/ui/build/` directory
- [x] 1.3 Add `internal/ui/build/` to `.gitignore` (build artifacts should not be committed)

## 2. Move UI Assets to Source Location

- [x] 2.1 Copy `internal/api/ui.html` to `internal/ui/index.html`
- [x] 2.2 Copy `internal/api/ui.css` to `internal/ui/style.css`
- [x] 2.3 Copy `internal/api/ui.js` to `internal/ui/app.js`
- [x] 2.4 Delete `internal/api/ui.html`
- [x] 2.5 Delete `internal/api/ui.css`
- [x] 2.6 Delete `internal/api/ui.js`

## 3. Update Makefile

- [x] 3.1 Add `ui-build` target that copies files from `internal/ui/` source to `internal/api/ui-build/` for embedding
- [x] 3.2 Make `build` target depend on `ui-build`
- [x] 3.3 Update `clean` target to remove `internal/api/ui-build/`

## 4. Update Embed Paths

- [x] 4.1 Update `//go:embed` directives in `internal/api/ui.go` to point to `ui-build/index.html`
- [x] 4.2 Update `//go:embed` directives in `internal/api/ui.go` to point to `ui-build/style.css`
- [x] 4.3 Update `//go:embed` directives in `internal/api/ui.go` to point to `ui-build/app.js`

## 5. Verify

- [x] 5.1 Run `make clean` to verify clean removes build artifacts
- [x] 5.2 Run `make build` to verify successful compilation
- [x] 5.3 Test `/ui/` endpoint returns HTML
- [x] 5.4 Test `/ui/style.css` endpoint returns CSS
- [x] 5.5 Test `/ui/app.js` endpoint returns JavaScript

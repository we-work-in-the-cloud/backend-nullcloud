## Why

The UI assets currently live in `internal/api/` alongside API handlers, mixing concerns. Moving UI to its own `internal/ui/` folder clarifies the separation between API and UI code, and establishes a dedicated build process that can evolve to support modern frameworks later.

## What Changes

- Create `internal/ui/` folder structure to house all UI code and assets
- Add a `build/` subdirectory to hold compiled/prepared assets
- Move ui.html, ui.css, ui.js from `internal/api/` to `internal/ui/build/`
- Add a UI build step to the Makefile (runs before Go compilation)
- Update embed directives in `internal/api/ui.go` to point to new paths
- Remove ui files from `internal/api/` (only keep the handler code)

## Capabilities

### New Capabilities
- `ui-build-process`: Establish a dedicated build step for UI assets that runs before Go compilation, enabling future framework integration

### Modified Capabilities
- `ui-serving`: The `/ui` endpoint behavior remains unchanged, but assets are now sourced from the restructured `internal/ui/build/` directory

## Impact

- `internal/ui/` — new folder with `build/` subdirectory
- `internal/api/ui.go` — updated embed paths only (logic unchanged)
- `GNUmakefile` — new UI build target
- Removed files: `internal/api/ui.html`, `internal/api/ui.css`, `internal/api/ui.js`

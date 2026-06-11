## MODIFIED Requirements

### Requirement: UI assets are built before Go compilation

The system SHALL have a dedicated build step that runs npm build to compile React code into `internal/ui/build/` before the Go binary is compiled. This enables framework migration while maintaining a clean separation between UI and API code.

#### Scenario: Make target builds UI assets via npm

- **WHEN** running `make build` or `make` (default target)
- **THEN** the UI build step executes first
- **AND** `npm run build` is executed in `internal/ui/`
- **AND** compiled React app outputs to `internal/ui/build/`
- **AND** the Go build step then embeds assets from the new location

#### Scenario: Clean removes UI build artifacts

- **WHEN** running `make clean`
- **THEN** the `internal/ui/build/` directory is removed
- **AND** the nullcloud-backend binary is removed

### Requirement: UI files are organized in internal/ui folder structure

The system SHALL organize UI source code in `internal/ui/src/` and build outputs in `internal/ui/build/`. The `internal/ui/` directory contains source code, package.json, Vite config, and build outputs.

#### Scenario: UI source files exist in src/

- **WHEN** the change is implemented
- **THEN** `internal/ui/src/` directory exists
- **AND** React component files are in `internal/ui/src/`
- **AND** `internal/ui/package.json` defines dependencies
- **AND** `internal/ui/vite.config.ts` configures the build

#### Scenario: Build outputs exist in build/

- **WHEN** the build completes successfully
- **THEN** `internal/ui/build/index.html` exists
- **AND** `internal/ui/build/app.js` exists
- **AND** `internal/ui/build/style.css` exists
- **AND** no other files are generated in build/

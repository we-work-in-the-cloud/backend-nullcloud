## ADDED Requirements

### Requirement: UI assets are built before Go compilation

The system SHALL have a dedicated build step that prepares UI assets in `internal/ui/build/` before the Go binary is compiled. This enables future expansion to support modern frameworks while maintaining a clean separation between UI and API code.

#### Scenario: Make target builds UI assets

- **WHEN** running `make build` or `make` (default target)
- **THEN** the UI build step executes first
- **AND** ui.html, ui.css, ui.js are copied/compiled into `internal/ui/build/`
- **AND** the Go build step then embeds assets from the new location

#### Scenario: Clean removes UI build artifacts

- **WHEN** running `make clean`
- **THEN** the `internal/ui/build/` directory is removed
- **AND** the nullcloud-backend binary is removed

### Requirement: UI files are organized in internal/ui folder structure

The system SHALL organize all UI assets under `internal/ui/build/` to clearly separate UI concerns from API logic.

#### Scenario: UI files exist in new location

- **WHEN** the build completes successfully
- **THEN** `internal/ui/build/ui.html` exists
- **AND** `internal/ui/build/ui.css` exists
- **AND** `internal/ui/build/ui.js` exists

#### Scenario: Old API folder no longer contains UI files

- **WHEN** the build completes successfully
- **THEN** `internal/api/ui.html` does not exist
- **AND** `internal/api/ui.css` does not exist
- **AND** `internal/api/ui.js` does not exist

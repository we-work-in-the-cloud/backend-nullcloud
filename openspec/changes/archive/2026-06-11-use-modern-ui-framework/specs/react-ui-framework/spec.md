## ADDED Requirements

### Requirement: UI is implemented with React and TypeScript

The system SHALL be built with React 18+ and TypeScript for component-based architecture and type safety. Components SHALL use hooks for state management and side effects. The build target MUST be a single application.html entry point with separate app.js and style.css outputs.

#### Scenario: React components load and render

- **WHEN** the application loads in a browser
- **THEN** React mounts successfully
- **AND** the root App component renders without errors
- **AND** the DOM contains the NullCloud UI elements

#### Scenario: Existing API endpoints are called

- **WHEN** the UI initializes or user interacts with it
- **THEN** all API calls target the same endpoints as before (`/v1/vpcs`, `/v1/subnets`, etc.)
- **AND** authorization headers are sent correctly
- **AND** responses are handled as before

### Requirement: UI is organized into reusable components

The system SHALL break the UI into composable React components. Major UI sections (ResourceTable, CreateModal, TokenInput, etc.) SHALL be separate components. Components SHALL be placed in `internal/ui/src/components/` with clear naming.

#### Scenario: Components are isolated and reusable

- **WHEN** a component is imported in multiple places
- **THEN** it renders consistently across all usages
- **AND** component state is isolated per instance
- **AND** props flow data correctly

#### Scenario: Modal and form components exist

- **WHEN** user clicks "Create VPC" or similar buttons
- **THEN** a modal component appears with the appropriate form
- **AND** the form component handles input validation
- **AND** form submission triggers the correct API call

### Requirement: Styling is maintainable and modern

The system SHALL use CSS (in `internal/ui/style.css`) organized by component or feature. Dark mode support SHALL be maintained. The visual design SHALL match or exceed the current UI in polish and usability.

#### Scenario: Light and dark modes work

- **WHEN** user toggles the theme
- **THEN** all components respond to the theme change
- **AND** colors update correctly for readability
- **AND** the preference is persisted

#### Scenario: Styling scales with responsive design

- **WHEN** the window is resized
- **THEN** the UI adapts to mobile, tablet, and desktop sizes
- **AND** layout remains functional on small screens

### Requirement: Build output is minimal (three files)

The system SHALL output exactly three files: `index.html`, `app.js`, and `style.css`. No separate chunk files, no source maps in production. The output MUST fit the existing embed pattern.

#### Scenario: Vite builds to three files

- **WHEN** running `make ui-build`
- **THEN** `internal/ui/build/index.html` exists
- **AND** `internal/ui/build/app.js` exists (minified, single bundle)
- **AND** `internal/ui/build/style.css` exists (all styles combined)
- **AND** no other files are generated in the build directory

#### Scenario: Binary embeds built files

- **WHEN** the Go binary is compiled after npm build
- **THEN** the `//go:embed` directive finds and embeds the three files
- **AND** the binary is generated successfully
- **AND** no build errors occur

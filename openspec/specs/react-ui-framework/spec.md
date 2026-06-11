# react-ui-framework Specification

## Purpose

Define requirements for the NullCloud UI modernization from vanilla JavaScript to React 18 with TypeScript, providing a modern, maintainable, and feature-rich management console.

## Requirements

### Requirement: Application shell with authentication

The UI SHALL provide a token-based authentication flow that persists across sessions.

#### Scenario: User enters API token

- **WHEN** user loads the UI for the first time
- **THEN** they see a welcome screen with token input
- **AND** token is persisted in localStorage
- **AND** subsequent visits auto-load the token

### Requirement: Resource management for all 7 resource types

The UI SHALL support full CRUD operations for VPCs, Subnets, Instances, Load Balancers, Buckets, Databases, and Kubernetes Clusters.

#### Scenario: User creates a resource

- **WHEN** user clicks Create button for a resource type
- **THEN** a modal opens with a form containing all required fields
- **AND** foreign key references show as dropdown selectors
- **AND** multi-select fields show as checkboxes
- **AND** form submission sends properly formatted JSON to the API
- **AND** success refreshes the resource list

#### Scenario: User edits a resource

- **WHEN** user clicks edit (✎) on a resource row
- **THEN** a modal opens with the resource name editable
- **AND** submission sends a PATCH request to update the resource

#### Scenario: User deletes a resource

- **WHEN** user clicks delete (✕) on a resource row
- **THEN** a confirmation dialog appears
- **AND** confirming sends a DELETE request
- **AND** successful deletion removes the resource from the list

### Requirement: Infrastructure hierarchy visualization

The UI SHALL display an interactive hierarchy view showing how infrastructure components relate to each other.

#### Scenario: User views the hierarchy

- **WHEN** user clicks the Hierarchy tab
- **THEN** they see regions as top-level containers
- **AND** VPCs and Buckets are shown per region
- **AND** Load Balancers are shown with their targets
- **AND** Zones contain subnets
- **AND** Subnets contain instances, databases, and clusters
- **AND** All components are displayed in an intuitive card-based layout

### Requirement: Theme support with persistence

The UI SHALL support both light and dark modes with persistence across sessions.

#### Scenario: User toggles dark mode

- **WHEN** user clicks the theme toggle in the header
- **THEN** the entire UI switches to dark mode
- **AND** the preference is saved to localStorage
- **AND** subsequent visits use the saved preference

### Requirement: Navigation via sidebar

The UI SHALL provide a sidebar navigation showing all resource types with current counts.

#### Scenario: User navigates between resource types

- **WHEN** user clicks a resource type in the sidebar
- **THEN** the main view switches to show that resource type
- **AND** the sidebar highlights the active type
- **AND** resource counts are displayed for each type

### Requirement: Proper API payload formatting

The UI SHALL send correctly formatted JSON to backend endpoints, including nested objects for foreign key references.

#### Scenario: Subnet creation with nested VPC reference

- **WHEN** user creates a subnet
- **THEN** the form sends `vpc: { id: "..." }` not `vpc_id: "..."`
- **AND** the API accepts the request and creates the resource

#### Scenario: Instance creation with nested references

- **WHEN** user creates an instance
- **THEN** the form sends `subnet: { id: "..." }`, `profile: { name: "..." }`, `image: { id: "..." }`
- **AND** the API accepts and processes the request correctly

### Requirement: Error handling and user feedback

The UI SHALL handle errors gracefully and provide clear feedback to users.

#### Scenario: Invalid API token

- **WHEN** user enters an invalid token
- **THEN** the API returns an error
- **AND** the UI displays an error message to the user
- **AND** user can try a different token

#### Scenario: Form validation

- **WHEN** user tries to submit a form with required fields missing
- **THEN** the form shows validation errors
- **AND** the submission is prevented until all fields are valid

### Requirement: 204 No Content response handling

The UI SHALL properly handle HTTP 204 responses from DELETE operations.

#### Scenario: Resource deletion completes

- **WHEN** a DELETE request succeeds with 204 No Content
- **THEN** the API response is properly handled
- **AND** no JSON parsing errors occur
- **AND** the resource is removed from the UI


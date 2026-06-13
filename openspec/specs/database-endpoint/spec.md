# Capability: Database Endpoint

## Purpose

Databases expose a deterministic, read-only endpoint attribute that clients can use to establish connections. The endpoint is computed from the database ID and engine type at creation time and remains constant throughout the database lifecycle.

## Requirements

### Requirement: Database has a computed endpoint
The Database object SHALL have a read-only `endpoint` attribute that provides a deterministic connection string in the format `db-{id}.db.nullcloud.internal:{port}`.

#### Scenario: Endpoint is computed on database creation
- **WHEN** a Database is created via the API
- **THEN** the endpoint is computed and stored with the Database record

#### Scenario: Endpoint format for PostgreSQL
- **WHEN** a Database with engine="postgres" is created
- **THEN** the endpoint uses port 5432 (e.g., `db-abc123.db.nullcloud.internal:5432`)

#### Scenario: Endpoint format for MySQL
- **WHEN** a Database with engine="mysql" is created
- **THEN** the endpoint uses port 3306 (e.g., `db-xyz789.db.nullcloud.internal:3306`)

#### Scenario: Endpoint format for MariaDB
- **WHEN** a Database with engine="mariadb" is created
- **THEN** the endpoint uses port 3306 (e.g., `db-xyz789.db.nullcloud.internal:3306`)

### Requirement: Endpoint is returned in all Database API responses
The endpoint attribute SHALL be included in all Database API responses (list, get, create).

#### Scenario: Endpoint in create response
- **WHEN** a new Database is created
- **THEN** the creation response includes the computed endpoint

#### Scenario: Endpoint in get response
- **WHEN** a Database is retrieved by ID
- **THEN** the response includes the endpoint attribute

#### Scenario: Endpoint in list response
- **WHEN** Databases are listed
- **THEN** each Database object includes the endpoint attribute

### Requirement: Endpoint is exposed in Terraform provider
The Terraform provider SHALL expose the endpoint as a computed (read-only) output attribute on the Database resource.

#### Scenario: Endpoint accessible in Terraform
- **WHEN** a Database resource is created in Terraform
- **THEN** users can access the endpoint via `nullcloud_database.*.endpoint` in their configuration

#### Scenario: Endpoint cannot be set by user
- **WHEN** a user attempts to set the endpoint field in Terraform
- **THEN** Terraform provider returns an error (endpoint is computed/read-only)

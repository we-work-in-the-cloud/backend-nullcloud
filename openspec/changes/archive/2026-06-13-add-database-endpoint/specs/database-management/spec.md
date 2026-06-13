## MODIFIED Requirements

### Requirement: Database resource provides connection information
The Database resource now provides a read-only endpoint attribute that clients can use to establish connections. The endpoint is deterministically computed from the database ID and engine type.

#### Scenario: Endpoint is available after database creation
- **WHEN** a Database has been created and reaches "available" status
- **THEN** the endpoint attribute is populated and accessible to clients

#### Scenario: Endpoint remains constant across queries
- **WHEN** a Database is retrieved multiple times
- **THEN** the endpoint value remains the same (deterministic)

#### Scenario: Different engines produce different ports
- **WHEN** two Databases are created with different engines (e.g., postgres and mysql)
- **THEN** each endpoint reflects the correct default port for its engine type

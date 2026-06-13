## Context

Currently, the Database model in `internal/model/model.go` has no connection endpoint. Users must manually construct connection strings or infer the infrastructure domain. By adding a computed `endpoint` field, we provide transparent, deterministic connectivity information across all surfaces: API, UI, Terraform, and documentation.

The endpoint will be computed in the API layer when a Database is created, then stored and returned with all subsequent queries. The format is deterministic: `db-{id}.db.nullcloud.internal:{port}`, where port is derived from the Engine type.

## Goals / Non-Goals

**Goals:**
- Expose a read-only, computed endpoint on all Database responses
- Make endpoint deterministic and based on database ID and engine type
- Support all existing engine types (postgres, mysql, mariadb) with appropriate default ports
- Surface endpoint in Terraform provider, UI, and documentation

**Non-Goals:**
- Allow users to customize the endpoint format or domain
- Support custom ports (use engine defaults only)
- Create a connection pooling or connection management system
- Change how databases are provisioned or created

## Decisions

### Decision 1: Compute endpoint in API layer (not on-read)
**Chosen Approach:** Compute the endpoint in the `createDatabase()` handler and store it with the Database record.

**Rationale:** 
- Consistency: The endpoint is always available and never changes
- Simplicity: One computation point, serialized directly in JSON responses
- Performance: No runtime computation for reads
- Terraform compatibility: Computed attributes are easier to expose as direct JSON fields

**Alternatives Considered:**
- Compute on-read in getDatabase/listDatabases: Would require a utility function called on every response, slightly slower. Benefit: source of truth is ID + Engine only. Drawback: string computation on every request.
- Computed property method on the struct (e.g., db.GetEndpoint()): Would require all consumers to call a method, not ideal for API serialization.

### Decision 2: Endpoint format: db-{id}.db.nullcloud.internal:{port}
**Chosen Approach:** Use the database ID (generated as `db-xxxxx`), the namespaced domain `db.nullcloud.internal`, and engine-specific port.

**Rationale:**
- ID-based is deterministic and unique per database
- Domain namespace `db.nullcloud.internal` clearly indicates the infrastructure scope
- Engine-specific ports follow database conventions (postgres=5432, mysql/mariadb=3306)

**Alternatives Considered:**
- Name-based (e.g., `{name}.db.nullcloud.internal`): Less reliable because names can be duplicated or changed. Would require uniqueness constraint.
- Random/assigned endpoint: Non-deterministic, harder to reason about.

### Decision 3: Port derivation from Engine field
**Chosen Approach:** Hardcode port mapping: postgres→5432, mysql→3306, mariadb→3306.

**Rationale:**
- Follows database industry standards
- Engine field already exists and is validated in createDatabase
- No additional user configuration needed

**Alternatives Considered:**
- Configurable per-database: Adds complexity, users rarely need to override default ports.
- Fixed port for all engines: Would break protocol assumptions (postgres and mysql use different wire protocols).

### Decision 4: Store endpoint in the model, not computed on read
**Chosen Approach:** Add `Endpoint string` field to the Database struct and persist it.

**Rationale:**
- Terraform provider serialization is simpler (direct JSON field)
- UI rendering gets the value directly from the API response
- No runtime computation or utility functions needed
- Straightforward API schema

**Alternatives Considered:**
- Don't store, compute on every response: Requires a utility function in the API layer; slightly higher CPU cost on list operations; Terraform provider would need custom logic.

## Risks / Trade-offs

**[Risk] Endpoint becomes stale if ID changes** → Mitigation: Database IDs are immutable; endpoint is computed once at creation and never changes.

**[Risk] Different engines might need non-standard ports in the future** → Mitigation: Currently hardcoded; if needed, we can refactor to a lookup table (e.g., config or database). For now, we accept the constraint.

**[Risk] Domain `nullcloud.internal` might change** → Mitigation: Currently hardcoded in endpoint computation. If infrastructure domain changes, we'd need to update endpoint computation logic. This is acceptable for now; consider making it configurable if multi-cloud support is planned.

**[Trade-off] Stored redundancy vs. computed on-read** → We chose to store for simplicity. Trade-off: The endpoint is technically derivable from ID and Engine, so we have data redundancy. Benefit outweighs cost in this case (simpler serialization, no runtime computation).

## Implementation Outline

1. **Model layer**: Add `Endpoint string` field to `Database` struct in `internal/model/model.go`
2. **API layer**: In `createDatabase()` handler, compute endpoint before storing the Database
3. **Terraform provider**: Expose endpoint as a computed attribute in the Database resource schema
4. **UI**: Display endpoint in database resource details view
5. **Examples**: Update database component examples to show endpoint usage
6. **Documentation**: Document the endpoint in API and provider reference docs

## Open Questions

- Should the `nullcloud.internal` domain be configurable via environment or config? (Deferring for now; assume hardcoded)
- Are there other resource types (LoadBalancer, Bucket, etc.) that should have similar endpoints? (Out of scope; Database is the pilot)

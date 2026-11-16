# Spec Delta

## Purpose

Provide a reproducible NullCloud Terraform Policy example that demonstrates resource rules and custom policy functions without provisioning infrastructure.

## ADDED Requirements

### Requirement: Local resource policy example
The example SHALL provide a bucket policy enforcing region `us-east` and a separate bucket policy checking CRN validity through a custom function. Policies SHALL reference the NullCloud provider schema and SHALL be locally validatable and testable without a running backend or credentials.

#### Scenario: Allowed region
- **WHEN** a mock bucket has region `us-east`
- **THEN** the region policy passes

#### Scenario: Disallowed region
- **WHEN** a mock bucket has region `us-west`
- **THEN** the region policy reports a violation identifying the region requirement

#### Scenario: CRN policy invokes the custom function
- **WHEN** a mock bucket has CRN `crn:nullcloud:bucket:bucket-1`
- **THEN** the CRN policy passes using `plugin::nullcloud::valid_crn` with expected type `bucket`

#### Scenario: CRN resource type mismatch
- **WHEN** a mock bucket has CRN `crn:nullcloud:vpc:vpc-1`
- **THEN** the CRN policy reports a violation identifying the CRN requirement

### Requirement: NullCloud CRN validation function
The policy plugin SHALL expose `valid_crn(crn, expected_type)` taking two strings and returning a boolean. It SHALL return true only for exactly four colon-separated components with prefix `crn:nullcloud`, a supported type matching `expected_type`, and a nonempty ID with no whitespace. Supported types SHALL be `vpc`, `subnet`, `instance`, `loadbalancer`, `bucket`, `database`, and `cluster`.

#### Scenario: Valid supported CRNs
- **WHEN** the function receives a correctly formatted CRN and matching expected type for any supported type
- **THEN** it returns true without network access

#### Scenario: Invalid identifiers fail closed
- **WHEN** a string argument has a wrong prefix, unsupported or mismatched type, empty ID, whitespace, or extra or missing colon-separated components
- **THEN** the function returns false without repairing or normalizing the input

#### Scenario: Invalid argument types
- **WHEN** the caller supplies a non-string or null argument
- **THEN** evaluation reports an argument diagnostic rather than silently accepting the value or crashing the plugin

### Requirement: Reproducible tests
The example SHALL include mock policy fixtures asserting both passing and violating outcomes, and automated function tests through the plugin registration/call interface. The documented local validation and test commands SHALL succeed on the recorded supported tool versions without creating infrastructure.

#### Scenario: Run the documented local workflow
- **WHEN** a developer builds the plugin and executes the documented policy validation and test commands on the supported tool versions
- **THEN** validation succeeds and all positive and expected-violation test cases pass

#### Scenario: Function interface coverage
- **WHEN** automated tests invoke the registered function
- **THEN** tests verify all supported types, malformed input, type mismatches, and argument diagnostics

### Requirement: Compatibility and evaluation limits
The example SHALL declare Terraform >= 1.16.0 for local testing, document Terraform >= 1.17 for HCP enforcement and tested tool versions, the distinction between policy-plugin and Terraform provider functions, and the experimental framework status. It SHALL explain that computed CRNs can defer enforcement until after apply, and that a post-apply violation does not undo created resources.

#### Scenario: Read prerequisites and safety limits
- **WHEN** a developer follows the example README
- **THEN** the minimum Terraform version, tested tool versions, computed-value timing, and experimental/non-production warning are explicit

### Requirement: Precompiled release distribution
The policy plugin SHALL be published as separately named platform-specific ZIP assets alongside the provider in the same GitHub release and at the same version. Each supported provider platform tuple SHALL have a matching plugin archive. Plugin archive digests SHALL be included in the release's signed SHA256SUMS manifest. Existing provider archive names, contents, and registry-compatible checksum/signature naming SHALL remain unchanged.

#### Scenario: Download a precompiled plugin
- **WHEN** a user selects `nullcloud-policy_<version>_<os>_<arch>.zip` from a provider release
- **THEN** the archive contains the matching standalone executable usable without installing Go or compiling source

#### Scenario: Verify a downloaded archive
- **WHEN** a user verifies the plugin archive against the signed release checksum manifest
- **THEN** the manifest contains its SHA256 digest and retains the existing release signature verification mechanism

#### Scenario: Provider installation stays separate
- **WHEN** the release packages the provider and policy plugin
- **THEN** provider ZIPs keep their existing names and executable contents and do not contain the policy-plugin executable

### Requirement: Isolated example and deployment guidance
The example SHALL leave existing provider schemas, backend contracts, and the provider root dependency graph unchanged. Its README SHALL explain release downloads, checksum verification, optional local plugin builds, and HCP target-platform binary requirements without requiring hosted policy-set configuration for local tests.

#### Scenario: Use the existing provider without policy examples
- **WHEN** a user builds or uses the provider without entering the policy example
- **THEN** no new policy runtime, backend configuration, or minimum Terraform requirement is imposed on that usage

#### Scenario: Prepare optional HCP deployment
- **WHEN** a developer reads the hosted deployment notes
- **THEN** they are told that deployment needs a matching platform binary in the policy VCS repository and that hosted enforcement was not verified by local mock tests

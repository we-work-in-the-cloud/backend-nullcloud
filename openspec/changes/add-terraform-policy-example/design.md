# Design

## Context

See `proposal.md` for motivation and `specs/terraform-policy/spec.md` for the behavior contract. The provider already exposes bucket `region` and computed `crn`. Both repositories define `crn:nullcloud:<resource-type>:<id>`; no model or schema changes are necessary. Earlier provider list-resource work was planned in this backend OpenSpec home.

Terraform provider functions and Terraform Policy plugins have different runtimes. `tfpolicy` loads an executable declared in `policy.plugins` and invokes functions as `plugin::<label>::<function>`. The current HashiCorp install documentation requires Terraform >= 1.16.0 for local usage; the HCP policy VCS documentation requires Terraform >= 1.17.

## Goals / Non-Goals

**Goals:** Demonstrate the full local policy/plugin loop using existing bucket attributes, mock resources, one custom function, and explicit failure cases.

**Non-Goals:** Provider function registration, API requests, credentials, provisioning, backend changes, general-purpose CRN parsing libraries, a new release service/installer, or a hosted policy-set deployment. No requirement bump for existing provider consumers.

## Decisions

### Approaches considered

1. Plain-HCL policies only: least code, but fails the requested custom-function demonstration.
2. A separate policy-plugin executable within the provider example: selected with the maintainer; matches the documented runtime and leaves the provider binary alone.
3. A Terraform provider function, optionally alongside the policy plugin: not selected; a second interface adds code and does not replace the policy-plugin integration.

### Keep source, policies, and tests together

Planned provider paths:

```
examples/policy/
  README.md
  policies/region/bucket-region.policy.hcl
  policies/crn/bucket-crn.policy.hcl
  tests/region/bucket-region.policytest.hcl
  tests/crn/bucket-crn.policytest.hcl
  plugins/src/nullcloud/
    go.mod
    go.sum
    main.go
    main_test.go
  plugins/bin/                 # locally built output, ignored
```

Use a separate example Go module with a pinned policy framework version. Do not add that dependency to the root provider module. Build commands live in the example README; add no root Make target or wrapper script unless actual execution proves it necessary. Update provider README discovery links and AGENTS.md so readers can locate the example and understand its separate runtime. Do not hand-edit generated provider docs: no provider schema or registered function is changing.

### Distribute through the existing provider release

Extend the existing `.goreleaser.yml` and `.github/workflows/release.yml` path. The policy plugin gets a distinct build/archive identity, built from its nested module, and shares the provider release tag/version. Produce `nullcloud-policy_<version>_<os>_<arch>.zip` with a `nullcloud-policy` executable (`.exe` on Windows), one per supported provider platform tuple. Reuse the current provider matrix and exclusions; if a target cannot compile, stop and document the supported subset for approval rather than silently dropping it. No separate release repository or version lifecycle.

Explicitly associate each archive with its build ID so adding the policy build never bundles it into a provider ZIP. Preserve existing provider archive and executable names and registry-compatible checksum/signature filenames. Include plugin ZIP digests in the existing signed SHA256SUMS manifest; verify compatibility with Terraform provider distribution and inspect the signature workflow. Prefer this existing trust path over a second checksum/signing system.

Add nested-module tests and a non-publishing release snapshot/package check to the existing CI path. Validate the installed GoReleaser configuration and exercise multi-module build/archive selection before choosing exact configuration syntax. Inspect generated ZIP members, full platform coverage, and SHA256SUMS entries; execute the host-compatible released plugin against the policy tests. Do not create a public tag or publish a release merely to test packaging.

The README leads with downloading a versioned plugin ZIP, verifying its digest against the release manifest (and signature using the existing release trust mechanism), extracting it to the policy plugin path, and preserving executable permissions. Building from source remains an alternative, not a prerequisite for binary consumers. HCP consumers must still commit the matching executable into their own policy VCS repository. One plugin path selects one binary: mixed-architecture agent pools must use the matching architecture. No automatic Terraform Registry installation or GitHub URL loading is implied.

### One pure custom function

Register `valid_crn(crn, expected_type)` in the executable and call `plugins.Serve()`. Use the standard library to split and inspect exactly four components, check exact case-sensitive prefix/type, reject unsupported types, and require a nonempty whitespace-free ID. Invalid string inputs return false, not an error; there is no normalization and no lookup for existence or ownership. This is validation of the documented CRN shape, not an invented UUID or numeric-ID constraint.

Test through the registered function call interface, not just the Go helper. Prefer the framework's simple typed registration; verify wrong-type and null diagnostics. If automatic conversion accepts values that violate the contract, use its direct typed registration interface rather than building a conversion wrapper. Pin the tested framework version after checking its API and dependency compatibility; do not assume a published tag contains current main's fixes.

### Two independent policies

Use a configured bucket `region == "us-east"` rule for the plan-time illustration and a separate `valid_crn(attrs.crn, "bucket")` rule for computed-value illustration. Separate files avoid making the configured-region rule wait for the unrelated computed CRN. Each file declares NullCloud `required_providers`; the CRN file also declares its plugin binary path. Both declare Terraform >= 1.16.0.

Use `.policytest.hcl` mocks to exercise both outcomes of each policy without backend access. Specify a tested available provider version rather than an unbounded requirement that may select an old schema. First confirm whether `tfpolicy` can load the existing published provider schema directly; if local provider installation is necessary, document and exercise its supported mechanism without adding a schema fork or fabricated provider responses.

### Keep planning in the backend, implementation in the sibling

Only proposal artifacts are created in this step. The installed OpenSpec 2.0.0 status reports `repo-local` scope and an `allowedEditRoots` entry for the backend only. Before an apply edits the provider sibling, use the installed CLI's supported scope/affected-area workflow and verify its reported permission to edit that sibling. If unavailable, stop and request an explicit scoped implementation handoff. Do not move the shared planning home, initialize duplicate specs, or silently treat a sibling as an allowed edit root.

## Risks / Trade-offs

- Experimental framework API/support status → Label the example non-production and record actual tested versions and dependency limitations.
- Computed CRNs are unavailable during initial planning → Explain deferred evaluation and that a post-apply policy failure does not roll back infrastructure; local mocks verify rule logic, not hosted lifecycle timing.
- Plugin binary architecture mismatch → Build locally for the host; document HCP linux/amd64 and agent-platform requirements. Do not commit a local macOS binary. Release archives provide matching binaries; hosted rollout still needs that binary in the policy VCS repository and remains unverified here.
- Separate module excluded from root `go test ./...` → Explicitly run its tests in addition to both repositories' existing checks.
- Provider schema download/version availability → Verify against a real published or supported locally installed provider during apply and report a blocker rather than replacing it with a mock schema.

## Migration Plan

No state migration or backend deployment. During apply, build and test the example in the provider sibling, document the verified commands, and run existing checks in both repositories. Rollback reverts the example, its README/agent-guideline references, and the added plugin release configuration; existing infrastructure and provider behavior remain unchanged.

## References

- Issue and maintainer clarification: https://github.com/we-work-in-the-cloud/terraform-provider-nullcloud/issues/5
- CLI minimum version: https://developer.hashicorp.com/terraform/policy/install
- Plugin runtime, builds, and deployment: https://developer.hashicorp.com/terraform/policy/write/plugins
- Policy plugin configuration: https://developer.hashicorp.com/terraform/policy/reference/policy/policy
- Local policy tests: https://developer.hashicorp.com/terraform/policy/write/tests
- HCP requirements and binary delivery: https://developer.hashicorp.com/terraform/cloud-docs/policy-enforcement/manage-policy-sets/vcs/terraform-policy-vcs
- Framework status and test entry point: https://github.com/hashicorp/terraform-policy-plugin-framework

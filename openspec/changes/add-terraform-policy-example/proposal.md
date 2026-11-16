# Proposal

## Why

Provider issue [#5](https://github.com/we-work-in-the-cloud/terraform-provider-nullcloud/issues/5) requests a Terraform Policy example and investigation of custom functions. The maintainer accepts Terraform Policy's minimum Terraform version and a NullCloud-specific function whose purpose is to demonstrate extensibility.

## What Changes

- Add a self-contained `examples/policy/` in the sibling Terraform provider repository, with policies and passing/failing mock tests runnable locally through `tfpolicy validate` and `tfpolicy test`.
- Demonstrate a rule on configured bucket region and a custom `valid_crn(crn, expected_type)` policy-plugin function checking the existing NullCloud CRN convention.
- Build the custom function as a separate policy-plugin executable in an isolated Go module; do not register a Terraform provider function or modify provider startup.
- Require Terraform >= 1.16.0 for local testing and document the current HCP Terraform >= 1.17 requirement, exact tested CLI/framework versions, computed-value evaluation, and experimental status.
- Publish precompiled policy-plugin ZIP archives and checksums alongside provider assets in the same GitHub release, using the same release version and supported provider platform matrix. Keep plugin assets distinctly named and separate from provider installation archives.
- Document download/checksum verification, local source builds, tests, and optional HCP deployment requirements, and link the example from the provider README.

## Capabilities

### New Capabilities

- `terraform-policy`: Locally testable NullCloud resource policies and a domain-specific policy-plugin function, with reproducible setup, precompiled release distribution, and explicit version/evaluation limitations.

### Modified Capabilities

None. Existing resources, list resources, backend APIs, and model contracts remain unchanged.

## Impact

- Planning artifacts live in this backend repository, following the existing provider-only OpenSpec precedent. Implementation targets `../terraform-provider-nullcloud/examples/policy/`, its README, its agent guidelines, `.goreleaser.yml`, and `.github/workflows/release.yml`.
- The example has its own `go.mod`/`go.sum` with a pinned `terraform-policy-plugin-framework` dependency; the provider root module and executable retain their current dependency and behavior contracts.
- Backend production code, Terraform state, and provider schemas are unchanged. Local policy tests use mocks and do not provision infrastructure or require backend credentials.
- Extend the existing GoReleaser/GitHub Actions release path rather than adding a repository, installer, or independent version lifecycle. Preserve provider ZIP names/content and registry-compatible signed checksum metadata.
- HCP enforcement setup and committing binaries into consumer policy repositories remain outside this change. Consumers download the appropriate release asset, verify its checksum, and include the executable in their policy VCS repository; Terraform does not fetch it automatically.
- OpenSpec currently reports repo-local edit scope limited to the backend. Before apply touches the provider sibling, resolve an explicit supported cross-repository implementation scope; do not silently bypass that boundary.

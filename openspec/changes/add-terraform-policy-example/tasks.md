# Tasks

## 1. Implementation scope and prerequisites

- [x] 1.1 Resolve the provider-sibling implementation scope using the installed OpenSpec supported workflow; verify allowed edit roots or obtain an explicit scoped implementation handoff before editing outside the backend. Stop if the boundary cannot be resolved.
- [x] 1.2 Verify Terraform >= 1.16.0, a working tfpolicy CLI, a real NullCloud provider schema source/version, and a compatible pinned policy-framework version; document the HCP Terraform >= 1.17 requirement and record actual local version output and schema-loading results in the example README, reporting blockers rather than fabricating schema responses.

## 2. Standalone NullCloud policy plugin

- [x] 2.1 Create the isolated module under `../terraform-provider-nullcloud/examples/policy/plugins/src/nullcloud/`, pin the framework dependency, and ignore local `plugins/bin/` output; verify module resolution and no changes to the provider root go.mod/go.sum.
- [x] 2.2 Add failing registered-function tests for every supported CRN type, mismatched/unsupported types, empty IDs, whitespace, wrong prefixes, extra/missing components, and non-string/null argument diagnostics; run module `go test ./...` and confirm failures reflect the missing contract.
- [x] 2.3 Implement and register `valid_crn(crn, expected_type)` and serve the plugin, using the minimum typed registration satisfying the spec; verify all module tests pass, `go vet ./...` succeeds, and the plugin builds for the local host.
- [x] 2.4 Document exact plugin build/test commands and the function's shape-only validation contract in the example README; execute those commands from the documented working directories and verify generated binaries remain ignored.

## 3. Policies and local mock tests

- [x] 3.1 Add separate bucket region and CRN `.policy.hcl` files with required provider/version declarations, Terraform >= 1.16.0, and the CRN plugin source; run real `tfpolicy validate` against the available provider schema and verify success.
- [x] 3.2 Add `.policytest.hcl` fixtures covering allowed/denied regions and valid/invalid/mismatched CRNs; run real `tfpolicy test` and verify expected-violation assertions as well as passing cases all succeed without a running backend or credentials.
- [x] 3.3 Finish the example README with exercised validation/test commands, recorded versions, independent plan-time versus computed-value policy timing, no automatic post-apply rollback, experimental/non-production warning, and optional HCP binary-platform/VCS requirements; verify commands run as written and distinguish local mock coverage from unverified hosted enforcement.
- [x] 3.4 Link the example from provider README and describe the standalone policy-plugin runtime in provider AGENTS.md; verify all new links/paths resolve and generated provider docs and backend API/model files are unchanged.

## 4. Release distribution

- [x] 4.1 Extend existing GoReleaser configuration for a distinct nested-module policy build and `nullcloud-policy_<version>_<os>_<arch>.zip` archives sharing the provider release version/platform matrix; validate the configuration and run non-publishing snapshot builds, verifying complete target coverage and unchanged provider ZIP names/content.
- [x] 4.2 Include policy archive digests in the existing signed SHA256SUMS release path and add nested-module tests plus non-publishing packaging verification to the existing workflow; verify every archive digest, archive/build separation, existing manifest/signature naming, and provider distribution compatibility without publishing a release.
- [x] 4.3 Document versioned binary download, checksum/signature verification, extraction and executable placement as the primary install path, source build as an alternative, and HCP consumer VCS delivery; test the host-compatible packaged executable with the policy fixtures and verify all documented paths and asset names match snapshot output.

## 5. Cross-repository integration checks

- [x] 5.1 Run `make test` in both repositories, provider `go build ./...`, and the separate plugin module's `go test ./...` and `go build`; verify actual exits and record any environmental blockers without claiming full success.
- [x] 5.2 Re-run the documented local plugin/policy workflow from a clean generated-output state, inspect both repository diffs, and run `openspec validate add-terraform-policy-example --strict`; verify there are no tracked binaries, provider schema/dependency changes, or out-of-scope backend changes.

## Workflow follow-up

- Review the finished implementation and local verification evidence before archiving the OpenSpec change.
- HCP policy-set deployment and hosted enforcement verification are not part of this example change.

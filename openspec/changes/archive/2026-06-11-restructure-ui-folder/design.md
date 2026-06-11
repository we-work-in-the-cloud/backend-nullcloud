## Context

Currently, UI assets (ui.html, ui.css, ui.js) are co-located with API handlers in `internal/api/`. They are embedded directly into the Go binary at build time using Go's `//go:embed` directive. The Makefile has a single `build` target that compiles the Go binary.

The project is prepared for modernization: a dedicated UI folder structure and build process will enable future migration to frameworks like React or Vue without architectural disruption.

## Goals / Non-Goals

**Goals:**
- Separate UI concerns from API logic through folder organization
- Establish a dedicated UI build step that runs before Go compilation
- Maintain identical endpoint behavior (`/ui` serves the same assets)
- Create a foundation for adding modern UI frameworks in Phase 2

**Non-Goals:**
- Migrate to a modern UI framework (deferred to Phase 2)
- Change API behavior or endpoints
- Modify authentication or routing logic
- Update the UI visually or functionally

## Decisions

### Decision 1: File Organization

**Chosen:** Store UI source files at `internal/ui/` root (index.html, style.css, app.js), and build outputs to `internal/ui/build/`. Add `internal/ui/build/` to `.gitignore`.

**Rationale:** The `build/` subdirectory signals "generated artifacts" and is excluded from git, while `internal/ui/` holds source files. This mirrors common project structures (e.g., `node_modules/`, `dist/`, `build/` are typically gitignored). When Phase 2 adds a framework, the pattern is already established—source files are committed, generated artifacts are not.

**Alternative:** Commit `internal/ui/build/` directly. Simpler for Phase 1, but requires gitignore adjustments in Phase 2 when builds become non-trivial.

### Decision 2: Build Process

**Chosen:** Add a `.PHONY` `ui-build` target in the Makefile that copies assets to `internal/ui/build/`, and make `build` depend on it.

**Rationale:** Explicit, reproducible, and easy to extend. When Phase 2 adds a Node build step, we'll replace the `cp` commands with framework-specific builds (e.g., `npm run build`).

**Alternative:** Use a shell script. More flexible but less integrated with Make workflow.

### Decision 3: Embed Path Updates

**Chosen:** Update `//go:embed` directives in `internal/api/ui.go` to reference relative paths in `internal/ui/build/`.

**Rationale:** Go's embed uses file paths relative to the Go source package. Updating paths is mechanical and requires only string changes, no logic rewrites.

**Alternative:** Move handlers to a new package in `internal/ui/`. More separation but introduces coupling and a new package to maintain.

### Decision 4: Backward Compatibility

**Chosen:** No changes to routing, request handling, or response behavior.

**Rationale:** The `/ui` endpoint must continue to work identically. This is a structural reorganization, not a functional change.

## Risks / Trade-offs

[Build dependency ordering] → Ensure `ui-build` runs before Go compilation. Mitigate by making `build` explicitly depend on `ui-build` in Makefile.

[Embed path errors] → If paths in `//go:embed` don't match files in `internal/ui/build/`, the build will fail with a cryptic error. Mitigate by verifying build succeeds and testing `/ui` endpoint after changes.

[Phase 2 friction] → When adding Node/framework builds, we'll need to integrate npm/framework tooling into the Makefile. Mitigate by keeping the pattern simple now (just `cp` commands).

## Migration Plan

1. Create `internal/ui/` and `internal/ui/build/` directories
2. Add `internal/ui/build/` to `.gitignore`
3. Copy ui.html → `internal/ui/index.html`, ui.css → `internal/ui/style.css`, ui.js → `internal/ui/app.js`
4. Remove old files from `internal/api/`
5. Add `ui-build` target to Makefile that copies from `internal/ui/` to `internal/ui/build/`
6. Make `build` depend on `ui-build`; update `clean` to remove `internal/ui/build/`
7. Update `//go:embed` paths in `internal/api/ui.go`
8. Test: run `make clean && make`, then verify `GET /ui/` works

No deployment complexity—this is a single binary rebuild. Rollback is a revert of git changes. When Phase 2 adds a framework, the Makefile's `ui-build` target simply swaps `cp` commands for `npm run build`.

## Open Questions

- Should we add a `.gitkeep` or similar to `internal/ui/` before Phase 2? (Nice-to-have, not blocking.)

## Why

The current vanilla JavaScript UI works but is difficult to maintain as complexity grows. A modern framework like React with component-based architecture will improve code organization, make the UI easier to extend with new features, and deliver a more polished user experience. Phase 1 established the build infrastructure; Phase 2 implements the framework migration.

## What Changes

- Replace vanilla JavaScript with React + TypeScript
- Organize UI into reusable components (ResourceTable, CreateModal, TokenInput, etc.)
- Add Vite as the build tool (replaces simple `cp` commands)
- Update Makefile to run `npm run build` instead of copying files
- Maintain the same output structure: index.html, app.js, style.css (three files as built by Vite)
- Enhance styling with modern CSS (keep existing design system, improve polish)
- Maintain 100% API compatibility with existing backend endpoints

## Capabilities

### New Capabilities
- `react-ui-framework`: Establish a modern, component-based React UI with TypeScript, improving maintainability and extensibility

### Modified Capabilities
- `ui-build-process`: Replace file copying with npm build pipeline (Vite)
- `ui-serving`: The `/ui` endpoint now serves a React-built application instead of vanilla JS, but behavior and endpoints remain identical

## Impact

- `internal/ui/src/` — new source directory with React components and TypeScript
- `internal/ui/package.json` — npm dependencies (React, Vite, etc.)
- `internal/ui/vite.config.ts` — Vite build configuration
- `internal/ui/index.html` — kept but now as Vite entry point
- `internal/ui/style.css` — shared stylesheet (enhanced)
- `GNUmakefile` — ui-build target updated to run npm build
- `.gitignore` — no changes (internal/ui/build/ already excluded)
- Removed: `internal/ui/app.js` (replaced by built app.js)
- No API changes (100% backward compatible)
- No database/backend changes

## Context

Phase 1 established a UI build infrastructure with a dedicated `internal/ui/` folder and build step. The current implementation copies static HTML/CSS/JS files. Phase 2 replaces this with a modern React application while maintaining the same three-file output constraint (index.html, app.js, style.css).

The UI currently communicates with the backend via RESTful API calls to `/v1/` endpoints. All backend logic remains unchanged; only the frontend is modernized. The goal is to improve maintainability, code organization, and visual polish without disrupting the backend.

## Goals / Non-Goals

**Goals:**
- Migrate UI to React 18+ with TypeScript for type safety and component reusability
- Organize UI into composable components (ResourceTable, CreateModal, TokenInput, etc.)
- Use Vite as the modern build tool (fast builds, minimal output)
- Maintain 100% API compatibility (no backend changes)
- Keep build output to exactly 3 files (index.html, app.js, style.css)
- Improve code maintainability and ease of adding new features
- Enhance UI/UX with modern styling while preserving the design system

**Non-Goals:**
- Change API endpoints or behavior (pure frontend modernization)
- Migrate backend logic or add new features to the API
- Introduce state management libraries (Redux, Zustand, etc.) — use React context and hooks
- Build a mobile app or SSR version
- Upgrade to a different backend language/framework
- Change authentication model or token handling

## Decisions

### Decision 1: Framework Choice – React 18

**Chosen:** React 18 with TypeScript.

**Rationale:** React has the largest ecosystem, most community resources, and best integration with Vite. TypeScript provides type safety, improving code quality and developer experience. React hooks are sufficient for state management in this UI complexity.

**Alternatives considered:**
- Vue 3: Equally capable, lighter weight, but smaller ecosystem
- Svelte: Compiles away framework, smallest bundle, but less community reach

### Decision 2: Build Tool – Vite

**Chosen:** Vite (v5+) with React preset.

**Rationale:** Vite produces clean, optimized output suitable for embedding. It can be configured to output exactly 3 files (no chunking, no source maps in production). Fast rebuild times improve DX. Well-integrated with React and TypeScript.

**Alternatives considered:**
- Webpack: Powerful but complex, harder to constrain output format
- Create React App: Simpler but less control over output, generates multiple chunks

### Decision 3: State Management – React Context + Hooks

**Chosen:** React Context API with useReducer for complex state, local component state for UI state.

**Rationale:** Sufficient for this application's complexity. Avoids external dependencies, simpler to understand, reduces bundle size. Global state (API tokens, resource lists) uses Context; local UI state (modal open/close, form inputs) uses useState.

**Alternatives considered:**
- Redux: Overkill for this use case, adds complexity and bundle size
- Zustand: Lightweight but unnecessary at current scope

### Decision 4: Component Structure

**Chosen:** Feature-based folder structure in `internal/ui/src/`:
```
src/
├── components/        (reusable UI components)
│   ├── ResourceTable.tsx
│   ├── CreateModal.tsx
│   ├── TokenInput.tsx
│   └── ...
├── features/          (feature-level components)
│   ├── VPCs/
│   ├── Subnets/
│   └── ...
├── hooks/             (custom React hooks)
│   ├── useAPI.ts
│   └── useTheme.ts
├── types/             (TypeScript types)
│   └── api.ts
├── App.tsx
└── main.tsx           (entry point)
```

**Rationale:** Clear separation of concerns. Reusable components in `components/`, feature-specific logic in `features/`. Easy to locate and modify code. Scales well as UI grows.

### Decision 5: Styling Approach

**Chosen:** Single CSS file (style.css) with BEM naming convention. Maintain existing color system and dark mode support. No CSS-in-JS library.

**Rationale:** Simpler, smaller output. Existing CSS system is clean. Single stylesheet embeds cleanly. Dark mode support via `:root` CSS variables and `[data-theme]` attributes (already established).

**Alternatives considered:**
- Tailwind CSS: Would require PostCSS, adds build complexity
- Styled-components: Adds runtime overhead, increases bundle size

### Decision 6: Bundling and Output Constraints

**Chosen:** Vite configured to output a single app.js bundle (no code splitting), single style.css, and index.html. No source maps in production.

**Rationale:** Meets the three-file constraint. Embedding into Go binary is simpler. Reduces complexity. App size is reasonable for the current feature set.

**Vite config:**
- `build.rollupOptions.output.inlineFormat = "umd"` to avoid dynamic imports
- `build.sourcemap = false` for production
- Single entry point (main.tsx)

### Decision 7: API Integration – No Changes to Backend

**Chosen:** Fetch API calls to existing `/v1/` endpoints. Same request/response format.

**Rationale:** Zero backend impact. Proven to work. Simple to implement. No new dependencies.

**No changes:**
- Authentication (token in Authorization header)
- Endpoint paths
- Request/response body formats
- Error handling patterns

## Risks / Trade-offs

[Bundle size increase] → React + Vite builds to ~100KB minified (vs. ~20KB vanilla JS). Mitigation: Acceptable tradeoff for maintainability; still small enough to embed and serve quickly.

[Build time increase] → npm build takes longer than `cp` commands (~10 seconds). Mitigation: Acceptable for development; npm install cache speeds up rebuild. Vite's dev server (optional for local development) is fast.

[Learning curve] → Team must learn React if unfamiliar. Mitigation: React is widely used and well-documented. Code review and documentation help.

[Single bundle approach] → Vite will warn about bundle size if exceeded. At current complexity, bundle fits easily. Mitigation: If future features significantly expand scope, consider code splitting or module federation.

[Node.js dependency] → npm is required to build. Mitigation: CI/CD pipelines will handle this; local development requires Node (standard practice).

## Migration Plan

1. Initialize npm project in `internal/ui/` (package.json, node_modules)
2. Create React app structure (src/, vite.config.ts, tsconfig.json)
3. Migrate UI logic from vanilla JS to React components
4. Update Makefile to run `npm run build` instead of `cp` commands
5. Test endpoints (`/ui/`, `/ui/style.css`, `/ui/app.js`)
6. Verify binary builds and embeds correctly
7. Full integration test (use UI in browser, test all features)

**Rollback strategy:** Revert git commits to restore Phase 1 state. The phase-based approach means we can always go back to the previous working state if needed.

## CI/CD Integration

Goreleaser's `before` hooks section must include the UI build step:

```yaml
before:
  hooks:
    - make ui-build    # Must run before go build
    - go mod tidy
```

This ensures the npm build happens before goreleaser invokes `go build`, so the embedded assets are available. The GitHub Actions workflow (`release.yml`) needs to:
1. Set up Node.js (npm)
2. Run the existing release workflow (it will call goreleaser, which calls the before hook)

The Makefile's `build` target already depends on `ui-build`, so local `make build` works. Goreleaser must explicitly call `make ui-build` in the before hook to ensure the same behavior in CI.

## Open Questions

- Should we set up a local Vite dev server (`npm run dev`) for faster development iteration? (Nice-to-have, not blocking.)
- Should we add ESLint/Prettier for code consistency? (Recommended but can be added later.)
- Should we add Jest/Vitest for component unit tests? (Good practice but not required for Phase 2 initial launch.)

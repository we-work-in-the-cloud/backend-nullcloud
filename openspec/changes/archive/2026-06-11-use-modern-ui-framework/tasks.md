## 1. Project Setup

- [x] 1.1 Initialize npm project in `internal/ui/` (npm init -y)
- [x] 1.2 Create package.json with React, Vite, TypeScript dependencies
- [x] 1.3 Create vite.config.ts with React preset
- [x] 1.4 Create tsconfig.json with strict mode enabled
- [x] 1.5 Run npm install to download dependencies
- [x] 1.6 Create .gitignore in internal/ui/ for node_modules, dist, .env

## 2. Project Structure

- [x] 2.1 Create `internal/ui/src/` directory
- [x] 2.2 Create `internal/ui/src/components/` for reusable components
- [x] 2.3 Create `internal/ui/src/features/` for feature-specific code
- [x] 2.4 Create `internal/ui/src/hooks/` for custom hooks
- [x] 2.5 Create `internal/ui/src/types/` for TypeScript types
- [x] 2.6 Create index.html template in `internal/ui/` (Vite entry point)
- [x] 2.7 Create src/main.tsx (React entry point)
- [x] 2.8 Create src/App.tsx (root component)

## 3. TypeScript Types & Utilities

- [x] 3.1 Create `src/types/api.ts` (API response types: VPC, Subnet, etc.)
- [x] 3.2 Create `src/types/index.ts` (export all types)
- [x] 3.3 Create `src/hooks/useAPI.ts` (fetch wrapper with auth)
- [x] 3.4 Create `src/hooks/useTheme.ts` (dark mode toggle logic)

## 4. Core Components

- [x] 4.1 Create `src/components/Header.tsx` (logo, theme toggle, connection pill)
- [x] 4.2 Create `src/components/TokenInput.tsx` (API token input, connect button)
- [x] 4.3 Create `src/components/ResourceTable.tsx` (table for resources in MainView)
- [x] 4.4 Create `src/components/Modal.tsx` (create/edit modal container)
- [x] 4.5 Create `src/components/FormField.tsx` (input field wrapper with TextInput, SelectField)
- [x] 4.6 Create `src/components/EditResourceModal.tsx` (edit resource modal)
- [x] 4.7 Create sidebar navigation (replaced StatsRow, integrated into MainView)
- [x] 4.8 Create `src/components/HierarchyView.tsx` (infrastructure hierarchy visualization)

## 5. Feature Components (VPCs, Subnets, etc.)

- [x] 5.1 Create `src/features/Forms.tsx` (create forms for all 7 resource types)
- [x] 5.2 Create create/edit/delete modals for all resource types
- [x] 5.3 Create resource table views with pagination and status display
- [x] 5.4 Implement form validation and error handling
- [x] 5.5 Add proper field types (dropdowns for references, checkboxes for multi-select)
- [x] 5.6 Create `src/components/HierarchyView.tsx` (infrastructure topology visualization)

## 6. Application State & Context

- [x] 6.1 Create `src/context/AuthContext.tsx` (manage API token, connection state)
- [x] 6.2 Create `src/context/ResourceContext.tsx` (manage resources: VPCs, Subnets, etc.)
- [x] 6.3 Create `src/App.tsx` main logic (token input, tab switching, resource loading)

## 7. Styling

- [x] 7.1 Create `internal/ui/src/style.css` with complete design system
- [x] 7.2 Implement CSS variables for light/dark mode (--brand, --bg, --surface, --border, etc.)
- [x] 7.3 Test light mode styling (theme works, all components styled)
- [x] 7.4 Test dark mode styling (toggle working, colors applied correctly)
- [x] 7.5 Verify responsive design (sidebar + main content, cards, tables)

## 8. Build & Integration

- [x] 8.1 Add npm build script to package.json (`vite build`)
- [x] 8.2 Configure Vite to output exactly 3 files (index.html, app.js, style.css)
- [x] 8.3 Update `GNUmakefile` ui-build target to run `npm run build` in `internal/ui/`
- [x] 8.4 Verify Makefile copies built files to `internal/api/ui-build/`
- [x] 8.5 Update Go embed paths in `internal/api/ui.go` (if needed)
- [x] 8.6 Run `make clean && make build` to verify full build succeeds

## 9. Testing & Verification

- [x] 9.1 Test `/ui/` endpoint returns React HTML
- [x] 9.2 Test `/ui/app.js` endpoint returns bundled React code
- [x] 9.3 Test `/ui/style.css` endpoint returns bundled CSS
- [x] 9.4 Test token input and API connection
- [x] 9.5 Test loading resources (VPCs, Subnets, Instances, etc.)
- [x] 9.6 Test creating all resource types via modals (forms with proper field types)
- [x] 9.7 Test editing and deleting resources (with proper API payloads)
- [x] 9.8 Test sidebar navigation between all resource types
- [x] 9.9 Test theme toggle (dark/light mode persistence via localStorage)
- [x] 9.10 Test responsive layout (sidebar + main content area)
- [x] 9.11 Test all API endpoints with correct payload formats (nested objects for references)
- [x] 9.12 Test error handling (form validation, API error messages)

## 10. CI/CD Integration

- [x] 10.1 Update `.goreleaser.yml` to add `make ui-build` in the `before.hooks` section
- [x] 10.2 Verify `.goreleaser.yml` runs `make ui-build` before Go build
- [x] 10.3 Update `.github/workflows/release.yml` to set up Node.js (actions/setup-node v4)
- [x] 10.4 Configure npm cache in the workflow for faster CI builds
- [x] 10.5 Test local `make build` works (builds UI then Go binary)
- [x] 10.6 Verify UI artifacts are embedded in binary and served at `/ui/`

## 11. Documentation & Cleanup

- [ ] 11.1 Add README to `internal/ui/` explaining build process and npm commands
- [ ] 11.2 Document component structure for future maintainers (src/components, src/features, src/hooks)
- [ ] 11.3 Document how to run the Vite dev server locally (`npm run dev`)
- [x] 11.4 Clean up old UI artifacts (removed old app.js, committed node_modules to .gitignore)
- [x] 11.5 Verify .gitignore excludes build artifacts (internal/ui/build/, internal/api/ui-build/)
- [ ] 11.6 Create CONTRIBUTING guide for adding new UI features

**Optional enhancements for future:**
- [ ] 11.7 Add unit tests for components (Jest + React Testing Library)
- [ ] 11.8 Add E2E tests (Cypress or Playwright)
- [ ] 11.9 Add Storybook for component documentation
- [ ] 11.10 Performance optimization (code splitting, lazy loading)

---

## Summary

**Phase 2 Complete:** React 18 + TypeScript UI fully implemented with:
- All 7 resource types (VPCs, Subnets, Instances, Load Balancers, Buckets, Databases, Clusters)
- Full CRUD operations (Create, Read, Update, Delete)
- Infrastructure hierarchy visualization
- Dark/light mode toggle
- Sidebar navigation
- Form validation with proper field types
- Token persistence
- Proper error handling

The UI maintains visual parity with the original vanilla JS design while adding new capabilities like the hierarchy view and improved navigation. All API integration is working correctly with proper payload formatting for each resource type.

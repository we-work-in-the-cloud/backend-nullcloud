# NullCloud Console UI

A modern React + TypeScript UI for the NullCloud backend, built with Vite.

## Project Structure

```
internal/ui/
├── src/
│   ├── components/      # Reusable React components
│   ├── features/        # Feature-specific components
│   ├── hooks/           # Custom React hooks (useAPI, useTheme)
│   ├── context/         # React Context (AuthContext, ResourceContext)
│   ├── types/           # TypeScript type definitions
│   ├── App.tsx          # Root component
│   ├── main.tsx         # React entry point
│   ├── MainView.tsx     # Main application view
│   └── style.css        # Global stylesheet
├── index.html           # Vite entry point
├── package.json         # npm dependencies and scripts
├── vite.config.ts       # Vite configuration
├── tsconfig.json        # TypeScript configuration
└── build/               # Build output (generated, not committed)
```

## Getting Started

### Install Dependencies

```bash
npm install
```

### Development

Run the Vite dev server:

```bash
npm run dev
```

The dev server runs at `http://localhost:5173` and hot-reloads on file changes.

### Build

```bash
npm run build
```

Produces optimized build in `build/` directory with:
- `index.html` - HTML entry point
- `app.js` - Bundled React application (minified)
- `style.css` - Combined styles (note: CSS is inlined in app.js)

### From Backend Root

Build the entire app with the Go backend:

```bash
make build
```

This runs:
1. `npm run build` in `internal/ui/`
2. Copies output to `internal/api/ui-build/`
3. Builds the Go binary with embedded assets

## Architecture

### Authentication

API token is stored in `AuthContext` and passed to all API calls via the `Authorization` header.

### Resource Management

Resources (VPCs, Subnets, etc.) are fetched via `ResourceContext` using the `useAPI` hook. Automatic fetching on connection.

### Styling

CSS uses CSS variables for theming:
- Light mode (default): `--bg`, `--surface`, `--text`, etc.
- Dark mode: Toggled via `[data-theme="dark"]` attribute
- Theme preference persisted to `localStorage`

### Components

- **Header.tsx** - Logo, theme toggle, connection status
- **TokenInput.tsx** - API token input, connect button
- **MainView.tsx** - Tab navigation and resource display
- Feature components can be added incrementally

## Adding New Components

1. Create component in `src/components/` or `src/features/`
2. Use `useAuth()` to access token
3. Use `useResources()` to access resource data
4. Use `useTheme()` for theme-aware styles
5. Export from component file

Example:

```typescript
import { useResources } from '../context/ResourceContext'

export default function MyComponent() {
  const { vpcs, loading } = useResources()
  
  return (
    <div>
      {loading ? <p>Loading...</p> : <p>{vpcs.length} VPCs</p>}
    </div>
  )
}
```

## Building Feature Components

To enhance the UI with detailed resource management:

1. **Create feature folders** under `src/features/`
2. **Add list components** for browsing resources
3. **Add modal components** for CRUD operations
4. **Use ResourceContext** for data access
5. **Test with API tokens**

## Troubleshooting

**Build fails with "terser not found"**
```bash
npm install --save-dev terser
```

**Styles not applying**
- Check that `style.css` is imported in `src/main.tsx`
- CSS is bundled into `app.js` by default

**API calls failing**
- Ensure token is provided to the app
- Check browser console for CORS/auth errors
- Verify backend is running on correct port

## Development Tips

- Use React DevTools browser extension for debugging
- `useTheme()` provides `theme` and `toggleTheme()`
- `useResources()` provides `refresh()` to manually fetch data
- TypeScript types defined in `src/types/api.ts`

## Build Output

The production build is a single `app.js` (163KB gzipped) bundled with:
- React 18
- All CSS and assets
- Minified and optimized

This is embedded in the Go binary at build time for zero external dependencies.

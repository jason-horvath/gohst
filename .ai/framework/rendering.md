# Rendering And Views

Related: `overview.md`, `application.md`, `controllers.md`, `sessions.md`

## Purpose

The current view layer is templ-based and is mediated by `internal/render/`. Controllers do not render raw template names. They render typed `render.Page` values or templ components.

Use this file for render flow, layouts, and templ context. Use `application.md` for application-layer UI composition rules such as component extraction, reusable field wrappers, and page-versus-component responsibilities.

## Current Rendering Flow

1. A controller calls `Render(w, r, page)` or `RenderPartial(w, r, component)`.
2. `BaseController` delegates to `render.View`.
3. `render.View` gathers session-backed request context such as CSRF, auth payload, flash messages, field errors, and request metadata.
4. The view looks up the active layout from the layout registry.
5. The layout wraps the page content with `templ.WithChildren(...)`.
6. The final HTML is written with the correct content type.

## Core Types

### `render.Page`

Use `render.Page` for full-page responses. It carries:

- `Title`
- `Content templ.Component`
- `Meta *render.PageMeta`

Page functions in `views/` should return `render.Page` values that controllers can pass directly to `Render`.

### `render.View`

`View` is created by `BaseController` and manages per-controller layout and metadata state.

Available operations include:

- `SetLayout(name)`
- `Render(w, r, page)`
- `RenderPartial(w, r, component)`
- `SetTitle(title)`
- `SetMeta(meta)`

Set layout defaults in controller constructors, not inside handlers.

## Layout Registration

Layouts are registered at boot in `cmd/web/main.go`:

- `layouts/default`
- `layouts/auth`

Those names are consumed by `View.SetLayout(...)`. If a layout is added, register it during application startup before serving requests.

## Request Context Available To Templ

`render.SetPageContext` injects shared page context into the templ render tree. Components should read from these helpers instead of threading common values through every function:

- `render.GetCSRFFromCtx(ctx)`
- `render.GetAuthFromCtx(ctx)`
- `render.GetFlashFromCtx(ctx)`
- `render.GetFieldErrorsFromCtx(ctx)`
- `render.GetRequestFromCtx(ctx)`

This is the mechanism that powers auth-aware nav, CSRF hidden inputs, flash rendering, and request metadata inside templ components.

## CSRF In Views

Forms should render the CSRF hidden input from render context:

```templ
@templ.Raw(string(render.GetCSRFFromCtx(ctx).Input))
```

This is already the current pattern in the auth forms and nav logout form.

## Asset Loading

Use `render.AssetsHead()` in layouts to inject Vite assets.

- Development mode uses the Vite dev server and entrypoints.
- Production mode reads the built manifest and emits script and stylesheet tags.

Layouts, not controllers, own asset injection.

## Request Metadata

`render.RequestProps` exposes the current path, method, and URL to the templ tree. Use that for active nav state or request-aware view behavior instead of reparsing request values in many components.

## Partial Rendering Rule

Use `RenderPartial` only for responses that should not be wrapped in a page layout, such as HTML fragments for progressive enhancement. Standard page requests should go through `Render` so shared layout, meta, assets, and request context remain consistent.

## Error Handling

`render.RenderError` is the framework fallback for render failures and recovered panics. Development mode can display detailed output; production mode should stay generic. Do not replace that with ad hoc HTML from controllers.

## Rendering Checklist

Use this when building or reviewing rendered responses:

- Return full-page HTML through `Render` with a `render.Page`.
- Return fragments through `RenderPartial` only when layout wrapping is intentionally omitted.
- Set layout defaults in controller constructors.
- Register any new layouts during web boot before requests are served.
- Keep pages composed from reusable components according to `application.md`.
- Use render context helpers for CSRF, auth, flash, field errors, and request metadata inside templ components.
- Keep asset injection inside layouts through `render.AssetsHead()`.
- Keep controller handlers focused on selecting pages and data, not hand-building HTML.
- Use framework error rendering for template failures instead of ad hoc controller fallbacks.

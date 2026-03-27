# Framework Overview

This folder documents the reusable Gohst framework conventions that should travel with every cloned project built on this repository.

## Layer Boundaries

- `internal/` is the framework core. Keep reusable HTTP, session, render, middleware, config, storage, and utility behavior here.
- `app/` is the application layer. Controllers, services, routes, config, and models here are free to adapt per cloned project.
- `views/` contains templ pages, layouts, partials, and UI components rendered by the application through the framework render package.
- `assets/` and `static/` contain front-end entrypoints and served static files.

Hard boundary rule:

- if a new feature, helper, service, interface, or utility is application-specific, keep it in `app/`
- if it is clearly app-agnostic and reusable across downstream projects, it belongs in `internal/`
- do not place code in `internal/` merely because it is convenient; it must earn its place by being framework-level reuse

## What These Docs Are For

Use `.ai/framework/` to capture contracts that should remain true across downstream projects:

- controller inheritance and handler style
- route composition and controller registration conventions
- middleware chaining behavior
- session lifecycle and request-context access
- templ rendering and layout registration

Do not use this folder to encode a specific project's business routes, domain modules, or app naming unless those choices are part of the reusable framework contract.

These files are not just descriptive reference. They are intended to act as build guidance and behavioral constraints for agents and contributors so feature work stays predictable across projects built on this framework.

## Document Map

- `application.md` covers application-layer boundaries, reusable UI composition, and view component rules.
- `controllers.md` covers controller inheritance and handler responsibilities.
- `routes.md` covers controller-owned route registration and router composition.
- `middleware.md` covers middleware contracts and ordering.
- `sessions.md` covers session lifecycle, flash, CSRF, and old-input rules.
- `rendering.md` covers templ pages, layouts, request context, and asset injection.

## Build Checklist

Use this checklist before and during feature work so implementation stays consistent:

- Confirm whether the work belongs in the framework layer or the application layer.
- If a requested feature is app-agnostic and reusable across projects, promote it into `internal/`; otherwise keep it in `app/`.
- Start from the framework docs before inventing a new pattern.
- Follow the application-layer and reusable component rules in `application.md`.
- Follow the controller inheritance rules in `controllers.md`.
- Follow the routing convention in `routes.md`, with controllers owning `RegisterRoutes()` and the app router mounting them.
- Apply middleware intentionally using `middleware.md` rather than ad hoc wrappers.
- Use session, flash, old-input, field-error, and regeneration rules from `sessions.md`.
- Use templ rendering, layout registration, and render-context rules from `rendering.md`.
- If abuse protection is relevant, attach the correct rate-limiter policy rather than skipping it.
- When a new rule is reusable across downstream projects, update `.ai/framework/`; when it is app-specific, place it in `.ai/project/`.

## Enforcement Intent

Agents should treat these docs as implementation constraints, not optional suggestions. The goal is that if the framework solves the same kind of problem twice, it should solve it in the same way unless there is a documented reason to diverge.

## Extension Rules

- Prefer extending `app/` before changing `internal/`.
- Only move behavior into `internal/` when it is reusable across many projects built on this framework and not tied to a specific application domain.
- Keep reusable application UI in `views/components/`; only move view helpers into `internal/` when they are framework-level contracts rather than app-level composition.
- If a cloned project needs business-specific conventions, document them in `.ai/project/` rather than editing framework docs unless the framework itself changed.

## Architectural Direction

- Feature controllers should own their route definitions through `RegisterRoutes() http.Handler`.
- `app/routes/routes.go` should compose and mount controllers rather than define feature handlers inline.
- `internal/routes.RegisterRouter()` wraps the application router with `middleware.Recover`.
- Controllers embed `AppController`, which embeds `BaseController`.
- Full-page rendering is templ-based and goes through `internal/render.View`.
- Layout functions are registered during application startup in `cmd/web/main.go`.
- Session state is injected by middleware and read from request context.

## Downstream Project Rule

Keep `AGENTS.md` as the stable top-level entry point. Future cloned projects should usually keep `.ai/framework/` intact, then add or adjust `.ai/project/` for app-specific rules.

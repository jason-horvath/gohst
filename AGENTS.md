# Gohst Agent Context

This repository contains two layers that must stay distinct:

- `internal/` is the reusable framework layer.
- `app/`, `views/`, `assets/`, `database/`, and project config are the application layer built on top of the framework.

Use this file as the stable entry point for all future cloned projects. Keep framework guidance in `.ai/framework/` and place project-specific guidance in `.ai/project/`. When a cloned project needs different business rules or route groupings, update the project docs first instead of rewriting the framework docs.

## Load Order

Start with `.ai/framework/overview.md` for repository boundaries and extension rules.

Load the matching framework doc before changing these areas:

- `.ai/framework/application.md` for `app/`, `views/`, and application-layer UI conventions
- `.ai/framework/controllers.md` for `app/controllers/` and `internal/controllers/`
- `.ai/framework/routes.md` for `app/routes/` and `internal/routes/`
- `.ai/framework/middleware.md` for `internal/middleware/`
- `.ai/framework/sessions.md` for `internal/session/` and auth/session flows
- `.ai/framework/rendering.md` for `internal/render/`, `views/`, layouts, and templ page composition

Then load any project-specific overrides from `.ai/project/` when they exist.

## Non-Negotiable Rules

- Do not move application-specific logic into `internal/`. Move code into `internal/` only when it is clearly app-agnostic and reusable across downstream projects built on this framework.
- Controller-owned `RegisterRoutes()` is the routing convention for this framework. Each feature controller should register its own routes and route-group middleware, and `app/routes/routes.go` should compose and mount controllers.
- All application controllers embed `AppController`, which embeds the framework `BaseController`.
- Reusable application UI belongs in `views/components/`; pages should compose reusable components rather than repeat markup.
- Full-page HTML responses render through templ components wrapped by registered layouts. Fragment responses use `RenderPartial`.
- Session access must come from request context via `session.FromContext(r.Context())`.
- Framework docs describe reusable contracts. Project docs may add stricter app conventions, but they should not contradict framework behavior without an intentional framework change.

## Routing Rule

Treat controller-owned route registration as the definitive framework pattern:

- each feature controller exposes a self-contained `RegisterRoutes() http.Handler`
- handler method definitions and route-group middleware stay with the controller
- `app/routes/routes.go` composes mount points and cross-controller structure
- outer framework concerns such as recovery stay at the router boundary

## Maintenance Rule

When the framework changes, update `.ai/framework/` and keep `.ai/project/` focused on the downstream application's business-specific architecture. That keeps cloned projects adaptable by changing one entry point and a small set of project docs instead of forking the whole framework guidance set.

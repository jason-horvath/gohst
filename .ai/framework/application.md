# Application Layer

Related: `overview.md`, `controllers.md`, `rendering.md`, `sessions.md`

## Purpose

This file defines the general rules for the application layer built on top of the framework. It covers how `app/` and `views/` should be organized so downstream projects build reusable, composable application code instead of repeating the same structures in multiple places.

These rules are framework guidance for application-layer work. They are not tied to a specific business domain.

## Application Boundary

- `internal/` is framework infrastructure and reusable framework contracts.
- `app/` is application behavior, orchestration, domain services, models, and configuration.
- `views/` is application UI expressed through templ pages, layouts, partials, and components.

The application layer is expected to vary per project. Its structure should still follow the same compositional rules across projects built on this framework.

Use `app/` as the default destination for new work. Promote code into `internal/` only when it is clearly app-agnostic and should become part of the reusable framework.

## Promotion Rule For `internal/`

This is a hard rule:

- when a new request produces code that is specific to the current application, keep it in `app/`
- when a new request produces code that is app-agnostic and can be applied across downstream projects, it should go into `internal/`
- if something belongs in `internal/`, treat that as a framework addition, not just a local implementation detail

Examples that may justify `internal/` placement:

- a reusable service abstraction that is not tied to one business domain
- a generic interface or helper that many applications built on this framework could use
- reusable middleware, render support, storage support, validation support, or HTTP infrastructure
- framework-safe UI render helpers that are truly cross-project rather than app-specific components

Examples that should stay in `app/`:

- business workflows
- project-specific controllers, services, and models
- app-specific page composition and UI components
- rules, policies, and behavior tied to one application domain

Do not move code into `internal/` just because it looks neat there. It belongs there only when it is genuinely framework-level reuse.

## General Application Rules

- Keep business-specific behavior in `app/`, not `internal/`.
- Keep HTTP orchestration in controllers and business workflows in services or other application packages.
- Organize `app/` in the closest practical shape to `internal/` so the application layer feels familiar and easy to navigate, while still allowing application-specific deviations where needed.
- Prefer small, focused packages and files over large mixed-responsibility modules.
- Reuse application patterns before introducing new ones.
- If the same application concern appears more than once, stop and decide whether it should become a reusable application primitive.

## Structure Mirror Rule

`app/` should generally mirror the framework shape where that improves clarity.

Examples:

- controller concerns in `app/controllers/`
- route composition in `app/routes/`
- application services in `app/services/`
- application models in `app/models/`
- application config in `app/config/`

This is not a requirement to force fake symmetry. It is a rule to keep the application layer organized in a way that feels structurally consistent with the framework. Mirror the framework where it helps; diverge only when the application actually needs a different shape.

## Views Structure

Use `views/` with clear responsibilities:

- `views/pages/` for full-page page builders that return `render.Page` values.
- `views/layouts/` for full-page wrappers and shared shell structure.
- `views/partials/` for small layout-adjacent fragments that are reused at layout or page level.
- `views/components/` for reusable UI primitives and composed UI building blocks.

Pages should assemble components. Components should not turn into page-sized templates.

## Reusable Component Rule

Any UI that is meaningfully reusable should become a component instead of being rewritten inline.

Common examples:

- buttons
- labels
- inputs
- selects
- checkboxes and radios
- file inputs
- error messages
- cards
- form rows
- field groups made of label, control, help text, and error output

Do not keep copying the same markup, classes, aria attributes, or error rendering from page to page. Extract it into a component and reuse it.

## Component Granularity

Prefer components that are small enough to be reused but large enough to remove real duplication.

Good examples:

- a button primitive
- an input primitive
- a label primitive
- an error-message primitive
- a field wrapper that composes label, input, help text, and field errors
- a card shell used in multiple screens

Avoid two extremes:

- page templates filled with repeated raw markup
- overly abstract components that hide simple markup without providing reuse value

Extract components where the composition becomes stable and obviously reusable.

## Form Composition Rule

Form UI should be built from reusable primitives and wrappers, not handcrafted independently in every page.

Preferred direction:

1. Use primitive components for controls such as input, select, checkbox, radio, file, label, and button.
2. Use composed field components for repeated control patterns such as label plus input plus error state.
3. Use a dedicated error component for field-level and form-level error rendering.
4. Feed old values, field errors, disabled state, and accessibility attributes into those components rather than rebuilding the same behavior inline.

If a page needs a labeled input with error output, the default assumption should be that it is composed from reusable field parts rather than written from scratch.

## Error Rendering Rule

Validation and feedback UI must be reusable.

- Keep error rendering in dedicated components when the same pattern appears across forms.
- Field-level components should be able to render their own errors consistently.
- Reuse the same visual and semantic structure for error output across pages.
- Do not invent a new error markup pattern for each form.

This applies to inline field errors, grouped form errors, and other repeated feedback blocks.

## Composition Over Copying

When building UI:

- compose pages from components
- compose complex components from simpler components
- extract repeated markup before it spreads across multiple pages
- standardize class combinations and structure in one place when they repeat

The goal is not abstraction for its own sake. The goal is that the same UI problem is solved once and reused consistently.

## Reuse Threshold

Use judgment, but default toward extraction when any of these are true:

- the same markup pattern appears in more than one page
- the same classes and structure are being copied with only text changes
- the same field shape is repeated across multiple forms
- the same accessibility or error-handling behavior would otherwise be duplicated
- the same card, panel, action row, or status block appears in more than one feature area

If reuse is plausible and the shape is stable, create the component early rather than waiting for copy-paste to spread.

## Application UI Ownership

- Keep app-specific reusable UI in `views/components/`.
- Keep page assembly in `views/pages/`.
- Keep layout shell concerns in `views/layouts/` and `views/partials/`.
- Keep framework render mechanics in `internal/render/`.

Only move something from `views/components/` into `internal/` when it has become a true framework-level contract that should exist across many downstream projects.

## Practical Rules For Agents

- Default new non-framework work into `app/`.
- When a requested abstraction is clearly reusable across many applications, raise it to `internal/` as a framework addition.
- Treat promotion into `internal/` as a deliberate design decision, not a convenience move.
- Before adding new page markup, check whether an existing component already solves the problem.
- If repeated markup is being introduced, stop and extract a reusable component.
- Prefer extending an existing component family before creating a parallel variant with overlapping responsibilities.
- Keep component APIs clear and intentional rather than passing arbitrary page state everywhere.
- Keep styling and semantic structure consistent across reused UI primitives.
- Treat form controls, field wrappers, and validation feedback as reusable building blocks by default.

## Application Checklist

Use this when building or reviewing application-layer work:

- Keep framework concerns in `internal/` and app concerns in `app/` and `views/`.
- Decide explicitly whether the new work is application-specific or framework-level reusable.
- Put app-agnostic, cross-project framework additions in `internal/`.
- Keep application-specific features in `app/` even when they are cleanly abstracted.
- Keep controllers thin and move reusable business logic into application services or packages.
- Keep `app/` organized in the closest practical structure to `internal/`.
- Build pages from reusable components instead of repeating markup.
- Put reusable UI primitives and composed UI blocks in `views/components/`.
- Extract repeated form controls and error rendering into reusable field components.
- Keep layout concerns in layouts and partials, not inside every page.
- Reuse existing components before introducing new markup patterns.
- Move something into `internal/` only when it is truly framework-wide rather than application-layer reuse.

## Enforcement Note

Agents should treat these as build rules. Repeated UI should be refactored into components rather than copied forward, and new application work should extend the existing component system whenever possible.

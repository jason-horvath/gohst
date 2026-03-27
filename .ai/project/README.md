# Project Overrides

This folder is the project-specific companion to `.ai/framework/`.

Use it in cloned applications to document rules that are specific to that application but still built on the same Gohst framework, for example:

- domain modules and route groups
- business-specific controller/service boundaries
- authorization rules tied to the project
- naming conventions for feature areas
- project-specific UI, content, or workflow constraints

## Precedence

- `.ai/framework/` defines reusable framework contracts.
- `.ai/project/` adds or narrows conventions for the current application.
- If a project rule appears to conflict with a framework rule, either the project rule is wrong or the framework has changed and `.ai/framework/` needs to be updated.

## Recommended Files

Add focused project docs here rather than one large file. Typical examples:

- `routing.md`
- `features.md`
- `auth-rules.md`
- `ui-patterns.md`

Keep `AGENTS.md` as the single top-level file that points agents to both framework and project context.

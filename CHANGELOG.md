# March 24 2022 - Initial release

# Feb 22 2026 
Adds parser instrumentation framework:
- Listener and Adjuster options.
- Policy/registry/reporter/action model for parser-time transformations and diagnostics.
- TransformTemplate support for source-to-source rewrite via patches.

Adds context-aware execution:

- est.State carries context.Context.
- Execution.ExecWithContext(ctx, state) introduced.
- Functions/methods with context.Context as first param are auto-injected from state context.

- Improves parser/planner capabilities:

- Span capture support for AST nodes.
- Expanded selector/call/index matching logic.
- Planner hook callbacks for variable/function/foreach resolution.

Adds foreach runtime context:

- $foreach data (index, count, hasNext, first, last) during loop execution. 
- 
Updates truthiness handling in if evaluation:

- More explicit coercion behavior for slices/strings/numbers/pointers.
- Special handling for pointer-to-struct with IsZero/Zero methods.

- Adds substantial test and benchmark coverage:

- New tests for context passing, foreach context, unary/if behavior, parser spans, pointer/slice/map selectors.
- Added benchmark artifacts/scripts.

Concurrency-safety hardening for parser spans (follow-up in this branch):

- Removed global parse span state.
- Switched to per-parse span ownership and explicit span propagation through hook paths.

Pool/session lifecycle fixes (follow-up in this branch):

- State.Reset() now restores context.Background().
- State.Reset() now marks pooled state reusable (isTaken = false).
- Added pool reuse regression test.

Performance outcome on your production VM corpus (*.vm):

- Compile behavior: all templates compile on both branches.
- Final compile benchmark is effectively at parity/slightly better than main (no material regression).
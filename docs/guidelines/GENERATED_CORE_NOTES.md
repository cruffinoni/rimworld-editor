# Generated Core Notes

## Context
The project core flow is:

1) Generate Go files from XML (intended as a one-time generation).
2) Use the generated structs to manipulate game data.
3) Reconstruct the XML file/savegame from those structs.

This is the current foundation of the project, with some concern that it may not be compatible with all savegames (especially modded ones).

## Risks Observed
- A "generate once" model is brittle when savegames evolve or include mods.
- Unmodeled fields can cause unmarshal failures or, worse, data loss on round-trip.

## Recommended Direction
Keep the generated structs as the core but make the system extensible and version-aware.

Suggested adjustments:
- Unknown/extra data preservation: capture unknown nodes/attributes and re-emit them on save to avoid data loss.
- Versioned generated schemas: `generated/v1_4`, `generated/v1_5`, etc., with a selector that chooses based on savegame header.
- Hybrid editing flow:
  - Typed editor for known fields.
  - XML tree fallback for unrecognized/modded data.

## Outcome
The generated types remain the central model for editing known data, while compatibility and data integrity improve for modded or future savegames.

## Open Decisions
- Whether to keep `generated` at repo root as public API or move to `internal/generated` and expose it via a versioned package.
- Whether to implement compatibility detection at runtime (savegame version, modlist signature) or via user selection.

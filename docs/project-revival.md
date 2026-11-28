# RimWorld Editor — Project Revival & Architectural Direction

*November 2026 · Project status and next steps*

## 1. Project Overview

**Repository:** [cruffinoni/rimworld-editor](https://github.com/cruffinoni/rimworld-editor)  
**Language:** Go 1.27  
**Status:** Existing project modernized; architectural exploration underway.

RimWorld Editor is a Go-based project designed to read, manipulate, and write RimWorld save games (`.rws` files).

The project originated from a technical challenge: **automatically generating a complete, strongly typed Go representation of a RimWorld save game from its XML structure.**

This generated type system became the foundation of the editor, allowing game entities such as pawns, skills, injuries, factions, and world objects to be manipulated through Go structures.

The original implementation demonstrated that complex save games could be successfully read, edited, and written back. However, compatibility was never universal.

## 2. Current State

The project saw most of its development during 2022–2023, followed by maintenance and refactoring in 2024 and December 2025.

The latest commit previously inspected, `e5a5881` (December 24, 2025), migrated the CLI to Cobra.

**The project has now been upgraded to Go 1.27.** The initial modernization phase is considered complete.

The repository includes:

- An XML parsing and serialization system.
- Go structure generators derived from XML.
- A terminal-based save editor.
- Editing capabilities for pawns, skills, injuries, factions, and other game data.
- Build tooling, tests, and CI configuration.
- Architectural notes in `docs/guidelines/`, particularly `GENERATED_CORE_NOTES.md`.

The existing architecture has a fundamental limitation: **generated Go structures can only represent the XML structures discovered during generation.**

RimWorld's save format is highly dynamic, particularly when mods are involved. Unknown structures, unexpected types, and variations between saves can cause compatibility problems and potentially compromise data integrity during serialization.

This architectural limitation was a major reason for losing motivation to continue development.

## 3. Architectural Discussions

### Separating responsibilities

The initial idea was to decouple save-file reading and Go structure generation from save-game editing logic.

This would allow the project to evolve into two independent components:

1. **Schema extractor / Go generator:** Responsible for discovering RimWorld's serialized structures and generating corresponding Go types.
2. **Save editor:** Responsible for reading, modifying, and writing save games through a reliable editing interface.

Each project should have a clear technical objective rather than maintaining two tightly coupled systems.

### Versioned schemas and mod support

One possibility was maintaining known vanilla structures for each RimWorld version, starting with 1.7, then progressively adding mod-specific schema patches.

However, this introduces substantial maintenance overhead and does not necessarily resolve the underlying compatibility problem.

This approach remains a possibility, but is not currently the preferred starting point.

### Generating types from C# instead of XML

A more promising direction emerged: **extracting serialization information directly from RimWorld's decompiled C# code instead of inferring structures from save-game XML.**

RimWorld uses the `Scribe` serialization system, with classes defining their serialization behavior primarily through `ExposeData()` methods.

Analyzing this code could reveal:

- Serialized field names and their actual C# types.
- Nested object relationships.
- Collections, references, and default values.
- Conditional serialization rules.
- Inheritance and polymorphic structures.

This would eliminate much of the trial-and-error process involved in generating types from representative save files.

A preliminary investigation of RimWorld 1.7's decompiled code suggests this approach is technically feasible, although custom serialization, conditional behavior, and runtime polymorphism remain significant challenges.

The objective is not to translate arbitrary C# into Go, but to extract its serialization contract as accurately as possible.

## 4. Proposed Direction

The immediate objective is **to validate the feasibility of generating Go types from RimWorld's C# serialization implementation** before committing to a larger redesign.

### Phase 1 — Restore the existing project

**Status: Completed**

The project has been upgraded to Go 1.27, completing the initial modernization phase.

The existing implementation provides the baseline for subsequent architectural experiments.

### Phase 2 — Prototype C# serialization analysis

**Status: Next priority**

Develop a small, independent proof of concept targeting RimWorld 1.7.

The prototype should:

1. Read locally decompiled RimWorld C# sources.
2. Parse the source using Roslyn or an equivalent analyzer.
3. Discover classes implementing `ExposeData()`.
4. Extract field declarations and supported `Scribe` calls.
5. Resolve basic type relationships, inheritance, and collections.
6. Generate corresponding Go structures.
7. Report unsupported or ambiguous serialization constructs explicitly.

Begin with a limited set of representative classes:

| Class | Purpose |
|---|---|
| `Pawn_HealthTracker` | Basic values and nested objects |
| `Faction` | References, collections, and defaults |
| `Pawn` | Inheritance and complex structures |
| `Map` | Custom serialization and analyzer limitations |

The goal is to evaluate the reliability and completeness of static extraction, not to generate an entire game model immediately.

The prototype should operate on locally available game assemblies or decompiled sources, without redistributing RimWorld's proprietary code.

### Phase 3 — Evaluate the results

**Status: Planned**

Compare generated Go structures against actual vanilla save data.

Determine:

- How much of the serialization schema can be discovered automatically?
- Which constructs require special handling?
- Can incomplete or unsupported definitions be identified reliably?
- How much additional work would mod support require?
- Can the generator become a genuinely independent and useful tool?

Only after these questions are answered should the existing save editor's architecture be reconsidered.

## 5. Long-Term Vision

The longer-term ambition remains to provide a powerful, convenient Go library for manipulating RimWorld saves while making the type-generation system useful independently.

Two distinct goals emerge:

**RimWorld Structure Generator**

A tool that analyzes RimWorld's serialization implementation and generates usable Go type definitions, potentially supporting multiple game versions and modded assemblies.

**RimWorld Save Editor**

A reliable library and CLI for reading, editing, and writing saves, with particular attention to preserving data that the editor does not explicitly understand.

The two projects may eventually integrate, but neither should be forced to depend on the other's completeness.

## 6. Immediate Milestone

**Prove that RimWorld 1.7 serialization declarations can be analyzed automatically and converted into correct, useful Go structures.**

Success means producing compilable types for a representative subset of classes, identifying unsupported serialization behavior, and validating the generated representations against actual save data.

The existing editor should remain functional during this exploration.

The intention is to make a technical decision based on a working experiment, rather than embark on another large architectural rewrite without knowing whether the underlying approach is viable.

# AGENTS.md

Instructions for AI coding agents working in this repository. Agents
load this file automatically; read it before doing anything else.

## Working agreement

- This project's code is hand-coded only.
- Any coding agent is used for **discussion only**:
  talking through approaches, explaining concepts, sketching
  pseudocode or type shapes on request, reviewing code the user
  wrote and pastes in. It is not for writing, generating, editing, or
  refactoring code automatically, and not for creating or editing
  project files directly.
- If asked something like "implement X" or "add Y", the right
  response is to talk through how to approach it (and sketch
  pseudocode inline in the chat if that helps), not to produce a
  diff, create a file, or edit one.
- This holds for the whole project.
- Exception: process/planning documents may be created and edited
  directly on request:
  - `DEVPLAN.md` - checkpoint status, known issues
  - `CONTEXT.md` - design decisions, architecture
  - `AGENTS.md` - this file
  - `README.md` - structure/skeleton only; the prose is hand-written
    by the user, so don't write or rewrite it unless explicitly asked
  This does not extend to anything under `internal/`, `cmd/`, or any
  `.go` file - those remain hand-coded only.
- The user is new to Go but not to programming. Be ready to explain
  syntax or Go-specific details, but assume there is no need to
  explain general programming concepts.

## Where things are

- `CONTEXT.md` - design: concept, biology, mechanics, aesthetic,
  current file architecture. Read for any design or code discussion.
- `DEVPLAN.md` - source of truth for build order, checkpoint status
  and known issues. Update it when status changes (rules for updating
  are at the top of that file). Note: gitignored.
- `README.md` - public-facing, hand-written by the user.

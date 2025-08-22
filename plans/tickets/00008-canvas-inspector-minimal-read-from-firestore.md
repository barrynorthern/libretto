# 00008 – Canvas/Inspector minimal read from Firestore (vertical slice UI)

Status: Won't Do
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
With GraphWrite persisting scenes to Firestore emulator, we can complete the vertical slice by adding a minimal UI that reads and displays scenes.

## Goal
Provide a basic UI page that lists scenes and shows a detail view sourced from Firestore emulator.

## Scope
- UI
  - Minimal page using existing stack (confirm tech; prefer Go or server-rendered where feasible; keep JS thin)
  - Endpoint(s) to read from Firestore emulator and render list + detail
  - Display title and summary for scenes
- Tests
  - Basic integration test: request page, contains scene title from emulator fixture
- Tooling
  - Make target to seed emulator with a scene document for local dev when needed

## Acceptance criteria
- Running the dev stack allows visiting a page that lists scenes stored via GraphWrite
- Selecting a scene shows title and summary
- Tests pass locally and in CI (using emulator)

## Non-functional requirements
- Keep load time snappy; minimal JS

## Out of scope
- Styling beyond bare minimum
- Complex navigation or editing


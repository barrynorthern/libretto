# 00007w – Firestore emulator store and UI skeleton

Status: Proposed
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
After collapsing to a monolith, we need persistence and a basic UI surface.

## Goal
Add a Firestore emulator-backed store implementation behind internal/graphwrite Store interface, and create a separate Next.js CSR app skeleton with a basic scenes list page.

## Scope
- Backend
  - Implement internal/graphwrite/store/firestore.go (emulator only)
  - Add Make targets to start/stop emulator
- Frontend (separate repo or subdirectory; prefer separate repo)
  - Scaffold Next.js CSR app with a design system (choose one: Mantine, Chakra, MUI)
  - Add a minimal /scenes page fetching from backend JSON API
- API
  - Expose a GET /api/scenes endpoint returning scenes (title, summary, id) from store

## Acceptance criteria
- Emulator starts via Make; backend writes/reads scenes
- UI loads /scenes and renders a list from the backend

## Out of scope
- Styling beyond basic components
- Auth and production Firestore


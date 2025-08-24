# 00009 – Wails desktop scaffold and shadcn/Tailwind setup

Status: Proposed
Owner: barrynorthern

## Context
We will ship a local-first desktop experience using Wails with a React frontend and shadcn/ui + Tailwind for fast, consistent UI.

## Goal
Scaffold the Wails app and a minimal UI page that invokes a Go binding.

## Scope
- Initialize Wails project under apps/desktop
- Set up React + Tailwind + shadcn/ui
- Create a "Scenes" page that calls a Go binding ListScenes() and renders results
- Add Make targets for wails dev/build

## Acceptance criteria
- `make desktop-dev` launches Wails dev window
- Scenes page renders with shadcn components and calls a Go binding

## Out of scope
- Advanced routing, state management, theming beyond defaults


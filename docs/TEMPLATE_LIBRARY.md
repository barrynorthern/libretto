# Template Library (v1)

Status: Draft (MVP)

## Plot Blueprints

### Three-Act Structure (v1)
- Parameters: genre, tone, theme
- Acts & beats:
  - Act I: Setup, Inciting Incident
  - Act II: Rising Action, Midpoint, Crisis
  - Act III: Climax, Resolution
- Output: arcs/main, beats per act (placeholders), initial scenes (empty content)

### Hero's Journey (v1)
- Parameters: genre, tone, theme
- Stages:
  - Ordinary World, Call to Adventure, Refusal, Meeting the Mentor, Crossing the Threshold
  - Tests, Allies, Enemies, Approach to the Inmost Cave, Ordeal
  - Reward, The Road Back, Resurrection, Return with the Elixir
- Output: arcs/main, beats per stage, initial scenes (empty content)

## Character Archetypes (v1)
- Protagonist: goals, flaw, transformation vector
- Antagonist: opposition vector, resources
- Mentor: guidance vector
- Ally, Trickster, Herald, Guardian, Shadow (lightweight fields)

## Data contract (skeleton)
- Template ID: e.g., three-act-v1, heros-journey-v1
- Parameters: { genre, tone, theme }
- Output: list of deltas for GraphWrite (arcs, beats, scene placeholders, character slots)

## Notes
- Templates produce a bootstrap branch (GraphVersion). Empty slots are intentional; Nudges suggest how to fill them.


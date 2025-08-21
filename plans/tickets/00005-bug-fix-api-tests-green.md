# 00005-bug – Fix API tests to green (compilation errors)

Status: Proposed
Owner: barrynorthern
Start: TBC
Date completed: pending

## Context
Running `go test ./...` shows API test failures unrelated to 00004 plotweaver work. These prevent a clean CI signal and should be addressed in their own branch/PR before proceeding with new feature work.

Observed errors:
- services/api/publisher/selector_test.go: unused import "os"
- services/api/server/baton_server_negative_test.go: handler return signature mismatch; connect-go now returns 2 values
- services/api/server/baton_server_negative_test.go: Attempt to call `MarshalJSON` on connect.Request which doesn't exist

## Goal
Restore a green test suite by fixing compilation errors in the API tests.

## Scope
- services/api/publisher/selector_test.go
  - Remove unused import(s) and ensure test compiles and asserts selector behavior.
- services/api/server/baton_server_negative_test.go
  - Update to current connect-go handler construction signature (capture both handler and path when needed)
  - Replace invalid `MarshalJSON` usage with proper request building using httptest and JSON body, or connect-go client request types
  - Keep the negative behavior under test intact, adapting to current server surface
- No production code changes unless strictly required to make tests compile; prefer amending tests.

## Acceptance criteria
- `go test ./...` passes locally
- CI shows green for API packages

## Non-functional notes
- Keep diffs small and focused; no dependency updates

## Out of scope
- Feature changes to API handlers beyond what's necessary for tests
- Broader refactors of publisher selection

## References
- Connect-Go docs: https://connect.build/docs/go/
- Current failing files:
  - services/api/publisher/selector_test.go
  - services/api/server/baton_server_negative_test.go


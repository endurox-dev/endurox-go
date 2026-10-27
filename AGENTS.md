# Working in endurox-go

Read [core's AGENTS.md](../endurox/AGENTS.md) for C, logging, error handling
and output compatibility rules. This file records Go-specific conventions and
maintainer preferences. Session instructions take precedence.

## Repository and compatibility

- This module wraps sibling `../endurox`; DLM services come from
  `../endurox-ee`, and the Python scripting provider from `../endurox-python`.
  `../endurox-rs` is a useful reference for native ownership and scripting;
  preserve the Go project's API and test conventions when consulting it.
- The v9 module path is `github.com/endurox-dev/endurox-go/v9`. Build with
  Go 1.26 or newer, cgo, and matching Enduro/X headers and runtime libraries.
  Native struct sizes and suspend tokens are ABI-dependent.
- Preserve released v8 method signatures and public struct contracts unless a
  breaking change is explicitly authorized. Unreleased v9 control types use the
  native-record/accessor design below. Do not keep redundant compatibility
  aliases for abandoned names in unreleased code.
- Inspect `git status` before editing. Preserve unrelated staged, unstaged and
  untracked work. Never stage binaries, generated field headers, logs, CMake
  output or local credentials. Do not push repositories; the operator pushes.

## Native records: one representation

- New public control and scripting types use mixed-case names such as
  `TpSagaCtl`, `TpTxStatus`, `TpDlmCtl`, `TpScrCfg` and `TpScrVm`. Initialisms
  stay upper case in types and methods alike (`XID`, `GetXID`, `GetCommitURL`).
  C ABI symbols keep their native spelling; pre-existing public types keep theirs.

- A new wrapper for a C control struct must keep that C struct as its
  authoritative storage. Getters and setters access its members directly;
  XATMI methods pass that same record. Do not create parallel exported Go
  fields and reconstruct a temporary C control for each operation.
- `TpSagaCtl`, `TpTxStatus`, `TpDlmCtl` and `TpScrCfg` follow this rule.
  Initialized wrapper copies share their backing record. A new zero value
  creates independent state. Access mutable controls serially.
- Validate fixed strings at the setter boundary: count UTF-8 bytes, reject
  embedded NUL, and leave room for the terminator. Invalid setters must preserve
  the previous value. Do not silently truncate. Validate native integer widths.
- Payload storage belongs to its control. `SetPayload` copies bytes and updates
  the native pointer and length together; `GetPayload` returns an independent
  copy. Clearing it releases native memory. Do not store a Go pointer in C-owned
  memory or let native code retain an unpinned Go allocation.
- Use `runtime.KeepAlive` through the last native use of contexts, controls,
  buffers and owners. Copies must not introduce double frees or allow an owner
  to be finalized while another wrapper still uses its memory.
- The `XID` wrapper contains `C.XID`; it is not a branch identifier or a
  suspend token. Use the native XID codec for global XIDs; Saga
  registration returns its branch identifier.
  `TPTRANID` already wraps the native suspend record. String conveniences
  encode/decode that record at the accessor boundary, not inside every API call.
- Opaque VMs retain their native handles.
  Callback and lifetime bookkeeping is separate from application control data;
  it is not a reason to duplicate a C struct's fields.

## Typed buffers and scripting

- `ATMIBuf` already owns `C_ptr` and `C_len`. Follow `TpCall` and `TpGetRply`.
  Do not add public pointer-to-pointer buffers, separate length pointers, or
  new slot objects to scripting APIs.
- Scripting execution and callbacks accept ordinary `TypedBuffer` values and
  return typed output. Native pointer/length temporaries are internal to cgo;
  write them back to the existing `ATMIBuf`, including on error.
- Aliased output reuses the existing buffer owner. Keep one owner/finalizer per
  allocation. Reallocation may move the pointer; failed reallocation preserves
  the old allocation. Returned bytes and snapshots must respect `C_len`.
- Callback buffer frames belong to the pinned native callback thread, not one
  Go context wrapper. Publish changes before provider reentry through any wrapper;
  native execution uses the existing buffer pointer/length fields so nested
  relocations remain visible to outer views.
- Callback buffers are borrowed for the invocation. Keep registered
  callbacks in stable native userdata via `runtime/cgo.Handle`, not Go pointers
  retained by C. Synchronize buffer edits before nested provider execution.
- A `TpScrVm` belongs to the `ATMICtx` that created it; its operations use
  that context, never a context argument. `FreeATMICtx` releases open VMs
  and does nothing while one is active. During cleanup, new VMs and
  compilation are refused; `TpScrUninit` invalidates the VM before shutdown.
  Callbacks receive the owning context object itself, so resources created
  there join the owner's registry; an invocation from another native context
  is rejected.
- Do not require or use `runtime.LockOSThread` for scripting. A single cgo call
  and a C-to-Go callback already stay on one OS thread. Read thread-local
  provider state (e.g. errors without a VM) in the same native call, in C.
- Destroy a VM explicitly through `vm.TpScrUninit`. Reject destruction during an
  active callback. A provider shutdown failure may already have consumed the
  handle: invalidate it and do not query or retry it. `TpScrComp` returns owned
  `[]byte`, valid after VM/context cleanup. Internally copy and immediately
  release native compiler output with `tpscrfree`, never `free` or `TpFree`.
  Preserve the compilation/copy error if native release also fails.
- Preserve native diagnostics before cleanup or error-query calls overwrite
  them. `TpBufCmp` returning -1 is successful ordering; only -2 is an error.
  `TpGetRply` must initialize the native descriptor from the requested value
  when `TPGETANY` is absent.

## Code and documentation

- Name source files and internal helpers by responsibility, not by the release
  that introduced them. Release numbers belong in module versions and changelogs.

- Follow nearby Go code and run `gofmt`. Use the standard file/license header.
  Keep tests and examples explicit: do not invent one-use helper wrappers just
  to shorten repeated API calls.
- Exported APIs have the project's `// @param` and `// @return` comments, as in
  `atmi.go`. These feed `doc/docupdate.pl`; godoc prose alone is insufficient.
  Document units, ownership, mutation, flags and failure behavior.
- Update the narrative examples and regenerate `doc/endurox-go-book.adoc` with
  `./docupdate.pl ../*.go` from the `doc` directory. Check that new receiver types
  have their own sections and that signatures and accessor examples render.
- Map API constants and flags from the matching C headers, including types.
  Generate UBF field IDs from field tables with `mkfldhdr` in application/test
  `ubftab` packages, including product fields from `Exfields`; do not maintain
  a binding-level `exfields.go` copy. Keep generated Go files out of Git.
  DLM operation/command names are strings; flags are numeric. Do not guess
  values or copy obsolete comments without checking the implementation.
- Use shared boundary helpers for repeated validation/ownership responsibilities.
  Do not add per-operation control mappers or generic frameworks for single uses.
  Keep new C glue in the core style: Allman braces, constant-first comparisons,
  documented arguments, explicitly initialized statics, and no embedded increments.

## Build and tests

- Functional tests belong in `tests/NN_functional_area/`, with their own
  `src/client`, `src/server`, configuration, Makefiles and `run.sh` where needed.
  Use neighboring cases as templates. Do not combine unrelated features in a
  catchall `v9` directory and do not put these integration tests in the package root.
- Register cases in `tests/Makefile` and `tests/run.sh`. Exercise real XATMI
  servers and native services, not only client-side argument validation.
  Fast scenarios use NDRX_TEST_DURATATION seconds (default 10); waiting scenarios may run once.
- Relevant areas: `13_buffer_utils`, `14_saga`, `15_tracing`, `16_scripting`,
  `17_dlm`, `18_multi_rm`. Saga tests use `tmsrv` and Go participant callbacks;
  DLM tests use `tpdlmsv` and a Go proxy; multi-RM tests use native XA fixtures.
- Scripting is separately licensed: default test runs must skip it. Explicit
  `-scripting` requires a configured provider and must fail if it is unavailable.
  DLM is independently enabled with `-dlm` and requires EE. Scope provider
  environment settings to the tests needing them.
- Never run shared-domain suites concurrently. Launchers source `~/ndrx_home`
  through the existing environment setup; inspect scripts before sourcing them.
  Other modules' environment scripts may stop all domains. Honor the session's
  test availability; do not request permission again after it was granted.
- Inspect effective binary/library paths. Legacy Makefiles can reuse stale
  executables because their dependencies omit the binding/native headers.
  After an ABI change, force rebuild the relevant cases, e.g. `make -B -C tests`.
  A stale v8 queue client with a v9 runtime can overwrite its stack.
- `go test ./...` and `go vet ./...` check packages; they do not execute the
  numbered integration launchers. Run relevant cases, inspect fresh diagnostics
  and preserve their exit status. Use `GOEXPERIMENT=cgocheck2` when checking new
  pointer ownership paths. Run `bash -n` for changed launchers and both ordinary
  and cached `git diff --check`.
- Record which binaries and libraries were actually used. A build alone is not
  runtime verification. Report blocked prerequisites and skipped tests accurately.
- Never print `~/ndrx_home`, credential files, full environment dumps or remote
  URLs containing tokens. Read only the specific nonsecret settings required.

Stress scenarios use NDRX_TEST_DURATATION in whole seconds, defaulting to 10
when unset or empty. Invalid, nonpositive or out-of-range values fail the test.
Each scenario completes its current operation/cycle before stopping and reports
the completed count and elapsed time. Waiting/recovery cases keep their bounded
single-pass sequence. Use NDRX_TEST_DURATATION=1 for a short check or 300 for load.
Native test clients suppress expected-error chatter; servers retain error logs.
Use the existing Enduro/X logging configuration for a short verbose diagnosis.

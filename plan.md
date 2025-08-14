# Table Library — Design & Plan

## Project overview

A tiny, dependency-free, high-performance table-rendering **library** for CLI apps. Not a standalone binary — intended to be embedded as a small module in larger codebases. Primary language target: Go (pure standard library), with an option to replace hotspots using a native backend (C/Rust) behind build tags if necessary.

## Goals & constraints (finalized)

* **No third-party dependencies.** Everything must compile using only the standard library.
* **Max speed / minimal allocations.** Prefer `[]byte`-first internals, reuse buffers, and minimize per-row allocations. Aim for zero-cost where possible.
* **Small size.** Keep the package surface minimal so downstream binaries don't bloat. Optional/native features behind build tags.
* **Flexible styling.** Support multiple border styles: single (rounded corners), double, rounded/curved, ASCII, and borderless output for clean columned text.
* **Correct-enough Unicode handling.** Handle ASCII well; use a compact, embedded heuristic/range table for CJK/emoji widths. No external width libraries.
* **ANSI-aware.** Preserve ANSI sequences in output but ignore them when measuring display width.
* **Library (not binary).** Small, single package that downstream code imports.

## Why this project matters

* Provides a tiny, consistent, and reusable table rendering utility across CLI apps.
* Avoids external dependency drag and reduces cross-project duplication.
* Lets downstream apps format output cleanly and efficiently without runtime penalties.

## Major problems & tradeoffs to be aware of

1. **Unicode display width complexity.** Perfect width calculation (emoji sequences, combining marks, and complex East-Asian rules) is hard. Without deps, we accept a pragmatic embedded-range heuristic that handles the common cases (CJK + many emoji) and document limitations.

2. **ANSI escape handling.** Must ignore SGR/CSI sequences while measuring widths but keep them in the rendered output.

3. **Cross-platform terminal size.** The provided ioctl snippet works only on Unix; Windows requires a different API. Terminal detection should be optional and behind a build-tag or a small platform switch.

4. **Allocation hotspots.** Converting strings repeatedly and using `fmt` in inner loops must be avoided; use byte buffers and pooling instead.

5. **Binary size (library footprint).** Go binaries are larger than C, but by keeping imports minimal and providing optional native backends behind build tags we can limit impact.

## Language choices & recommendations

* **Go (recommended first).** Fast to develop, good standard library, easy to test and publish. Use pure-Go by default with careful micro-optimizations for speed and low allocations.
* **Rust (optional backend).** Safe, zero-cost abstractions and faster performance in some string-handling hotspots; moderate binary size and cross-compilation complexity. Use behind a build tag / optional FFI.
* **C (optional backend).** Smallest binaries possible but manual memory/UTF-8 handling and more maintenance burden.

  Plan: implement a byte-centered pure-Go library first. Only replace critical hotspots after profiling.

## High-level architecture (library-focused)

* **Single public package** (example: `table`). Keep public API surface tiny.
* **Core types:** `Table`, `Column` (header, align, maxWidth), small `Style` type for border styles.
* **Rows representation:** accept rows as `[]interface{}` or `[]string` but convert once to `[][]byte` internally to operate on bytes.
* **Rendering modes:** two-pass rendering by default (measure widths then render). Also provide a streaming render mode for very large tables.
* **Width policy:** pluggable `WidthFunc` so callers can choose the heuristic or a more precise implementation.
* **Border styles:** encapsulate sets of characters for single/double/rounded/ascii/none and allow fallback to ASCII when Unicode unsupported.
* **Build tags:** optional native backend and platform-specific terminal detection behind build tags to keep default builds tiny.

## Width & Unicode strategy (no external deps)

* Use `utf8` decoding on `[]byte` to iterate runes.
* Implement a **compact embedded range table** for common wide rune ranges (CJK, Hangul, Hiragana, Katakana, and common emoji blocks). This keeps runtime zero-dep and small.
* Provide a fast heuristic that treats ASCII as width 1 and known ranges as width 2. Document edge cases (some emoji sequences and combining marks may be inaccurate).
* Make width calculation a pluggable function so a downstream consumer can swap in a more precise implementation if they want.

## ANSI handling

* During width measurement, detect and ignore ANSI CSI/SGR sequences (escape bracket sequences) so color codes do not affect column sizing.
* Preserve ANSI byte sequences in the final rendered output.
* Implement a small stateful scanner (no regex libs) to strip or skip ANSI sequences for measurement, avoiding allocations.

## Border styles & visual modes

Offer these styles (each as a small set of characters to draw table edges):

* **Single (Unicode single lines)** — compact and modern.
* **Double (Unicode double lines)** — bold boxed look.
* **Rounded / Curved** — soft corners.
* **ASCII** — pure ASCII fallback for terminals with poor Unicode support.
* **Borderless** — spacing-based alignment only (good for clean, compact text output).

Each style should be easily selectable via a `Style` option and default to a sensible modern style.

## API surface (small & focused)

Export a minimal API that’s easy to use and hard to misuse. Suggestion:

* `New(headers ...string) *Table` — create a new table instance.
* `(*Table) SetAlign(col int, align Align)` — optional per-column alignment.
* `(*Table) AddRow(vals ...interface{})` — add a row (convert once to bytes internally).
* `(*Table) Render(w io.Writer) error` — render the table to the provided writer.

Keep helper functions and heavy logic inside unexported/internal files.

## File layout (recommended)

* `table/` (package root)

  * `go.mod`
  * `README.md`
  * `LICENSE`
  * `table.go` — public API surface and orchestration
  * `style.go` — border style definitions and selection
  * `width.go` — width heuristics and embedded range table
  * `ansi.go` — ANSI scanner (no allocations)
  * `render.go` — rendering internals and streaming
  * `internal/` — small helpers and buffer pools
  * `*_test.go` — unit tests and example tests
  * `*_bench_test.go` — benchmarks

Keep the public package code minimal and put helpers under `internal/` or unexported files.

## Publishing to pkg.go.dev

To get published and visible:

1. Create a **public VCS repo** (GitHub/GitLab/Bitbucket) and push the module.
2. Ensure `go.mod`'s module path matches the repo URL (for example `module github.com/you/prettytable`).
3. Write godoc-style comments for package and exported symbols. pkg.go.dev renders these.
4. Add `Example` functions in test files so examples appear on pkg.go.dev.
5. Tag a semver release (recommended). For major version v2+, include `/v2` in the module path.
6. Optionally add CI (GitHub Actions) to run `go test` and `go vet`.

pkg.go.dev indexes public repositories automatically; tagging helps downstream `go get` pick stable versions.

## Performance techniques & micro-optimizations

* **Byte-first processing.** Convert incoming strings to `[]byte` once. Operate on bytes to avoid repeated allocations.
* **Buffer pools.** Use `sync.Pool` for temporary buffers and `bytes.Buffer` reuse.
* **Avoid fmt in hot loops.** Build byte slices and write them with `io.Writer` once per chunk.
* **Large ********`bufio.Writer`********.** Use a large buffered writer for final output to minimize syscalls.
* **Inlining & loop optimization.** Keep inner loops tight and avoid function call overhead in hot paths.
* **Fast width lookup.** Use a compact embedded range table and a binary search or well-structured switch for width determination.
* **Streaming render mode.** For huge tables, provide a streaming renderer to avoid holding all rows in memory.

## Milestones & step-by-step plan

1. **API & minimal prototype (library):** Single-file, no deps, byte-centric internals, ANSI-aware width heuristic, support for 4 border styles, two-pass render. Include `Example` tests and README.
2. **Micro-bench & profile:** Bench with 1k/10k/100k rows; identify hotspots (width calc, allocations, writes).
3. **Optimize hot paths:** Add `sync.Pool`, reuse buffers, inline critical loops, tune `bufio.Writer` size.
4. **Accuracy upgrades (if needed):** Replace heuristic with a more complete embedded range table for width; improve ANSI handling.
5. **Optional native backend:** If benchmarks show a measurable gap, implement a native width/formatting backend in Rust/C and expose via build tags. Default remains pure-Go.
6. **Polish & docs:** Add examples, CI, README, package-level godoc; tag release and publish.

## Acceptance criteria & benchmarks

* **Correctness:** Proper alignment and column widths for ASCII; good coverage for CJK/emoji using the heuristic with documented limitations.
* **Performance:** Rendering 100k small rows should be memory efficient and reasonably fast. Per-row extra allocations should be minimized and measured with `b.ReportAllocs()`.
* **Size:** Keep package imports minimal; no heavy standard-library imports that bloat downstream builds.

## Deliverables discussed

* **Primary deliverable (first):** Minimal, optimized single-package Go implementation (no external deps). Byte-oriented internals, ANSI-aware width heuristic, 4+ border styles, two-pass render, example tests, and benchmarks.
* **Optional deliverables:** Small single-file C implementation for extreme size/speed; Rust backend with FFI shim; generator for an embedded precise width table.

## Cross-platform terminal handling

* Unix: use ioctl-style approach (your provided snippet) for terminal width detection.
* Windows: implement `GetConsoleScreenBufferInfo` fallback under a Windows-specific build file.
* Make terminal-size detection optional (do not force it into every build). Provide a basic `GetTerminalSize()` helper that returns `(cols, rows, ok)`.

## Next steps (practical immediate actions)

1. Scaffold the repository with the suggested layout and minimal prototype files (public API, width heuristic, ANSI scanner, style definitions, example tests).
2. Add unit tests and two or three `Example` functions for pkg.go.dev.
3. Run micro-benchmarks and iterate on hotspots.
4. Publish to a public repo and tag

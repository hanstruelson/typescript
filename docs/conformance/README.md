# JavaScript compatibility audit

The full Test262 run was started on 2026-10-09. It targets this fork's native Go
emitter, with Node 24.19.0 checking that the harness can run each case.

- [Current full-suite checklist](prototype-full/CHECKLIST.md)
- [Current checkpoint totals](prototype-full/summary.json)
- [Individual executed results](prototype-full/results.jsonl)
- [Runner and reproduction instructions](../../tools/test262/README.md)

The checkout is `/workspace/test262`, revision
`2e0a56762801e275a9fdf96dc49d90ba0cddcf63`. The inventory contains 53,616 standalone
tests. After the initial textual dynamic-code exclusions, the full run schedules
100,058 strict/non-strict/raw/module variants. Module/host adapter limitations
are reported separately, not hidden or counted as passes. The complete inventory
is `full/inventory.json.gz`; it includes excluded cases.

**The report is a live checkpoint until `summary.json` says `complete: true`.**
This is not a claimed full-suite pass rate. Runtime compilation is expensive;
the run writes each result immediately and can resume after interruption.
The current process log is `/tmp/test262-prototype-full.log`. The original baseline in `full/` was stopped before the prototype implementation and is retained separately.

See the [prototype implementation notes](PROTOTYPES.md) and the [completed Date.now audit](prototype-date-now/summary.json): 8/12 variants pass after this change, compared with 0/12 before it.

## Confirmed initial blockers

- [x] Ordinary function objects now expose a working shared `.prototype` object. The original blocker was:
  The unmodified `harness/sta.js` alone fails before a test body can run:
  `Test262Error.prototype.toString = function () { ... }` results in
  `Async runtime errors: Expected an array`. This obstructs assertions across
  many unrelated categories; those failures do not establish the behavior of
  the individual test body.
- [ ] Standard property-verification helpers encounter unsupported emitter
  identifiers including `Function`, `Math`, `arguments`, and `TypeError`, and
  unsupported `delete` expressions in the initial Date.now audit.
- [ ] The runner still needs module-graph and `$262` realm/agent/GC/detachment
  adapters. These are test infrastructure gaps, distinct from proven language
  incompatibilities.
- [ ] Review conservative dynamic-code exclusions, including false positives
  from comments and indirect Function/eval aliases that textual scanning misses.

The first targeted audit executed all 12 variants in `built-ins/Date/now`:
4 emission errors and 8 runtime failures; none passed. That result exposed the
shared harness problems above, rather than providing an independent Date.now
compatibility score. The full-suite checklist records the broader evidence.

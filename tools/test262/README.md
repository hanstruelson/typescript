# Test262 conformance runner

This runs the upstream ECMAScript tests against JavaScript compiled to native Go.
The assertions and include files are taken from Test262 without substitutions.
Node is a harness/reference sanity check, not the implementation being tested.

```sh
git clone --depth 1 https://github.com/tc39/test262.git ../test262
python3 -m pip install -r tools/test262/requirements.txt
go build -o /tmp/test262emit ./tsc/cmd/test262emit
python3 tools/test262/run.py --all --jobs 2
```

Go 1.27+, Node, Python 3 and PyYAML must be installed. The Test262 checkout is a
sibling directory; its exact revision is recorded in each report. Test262 is
BSD licensed; keep its upstream LICENSE with the downloaded checkout.

`--all` selects every standalone test with Test262 metadata, rather than a
sample. Both strict and non-strict variants run unless metadata specifies one.
Without `--all`, two deterministic tests per category provide an initial audit.
Use `--filter built-ins/Date/` for a targeted run in a separate `--output`
directory. Each fixture emits, builds, and executes in its own temporary directory;
results are saved immediately and category checklists refresh every 100 variants.
Use `--resume` with the same revisions and output directory after an interruption.

Reports include a compressed inventory of **all** source tests, individual result
rows, totals and a category/failure checklist. The `complete` field stays false
until the selected run finishes. Exclusions and adapter limitations are never
counted as passes. Imported fixture files without metadata are not standalone tests.

Dynamic compilation via direct `eval` and Function-family constructor calls is
excluded. The initial scanner is conservative and textual: comments can match,
and indirect aliases can escape detection. Every exclusion is listed in the
inventory for review. Ordinary Function prototype tests remain in scope.

Module graphs and Test262 `$262` host operations (realms, agents, buffer detachment,
GC) need dedicated runner adapters and are reported as runner limitations. Node
may lack newer proposal features; its reference status remains separate and does
not prevent the native test from running. Parsing and
early-error tests pass only if the emitter reports a syntactic rejection matching
an expected SyntaxError. The runner does not treat unsupported Go emission as an
expected JavaScript exception. Runtime-negative tests must throw the expected
error name, and async tests must reach `$DONE` successfully. Error identity checks
inside upstream `assert.throws` are preserved.

The current compiler uses its default coercion policy. This measures that actual
mode, including any deviations from JavaScript, rather than switching semantics
to inflate compatibility. Full native compilation of tens of thousands of tests
can take hours. Repeated harness failures are real blockers and can obscure the
feature under test; investigate them before interpreting category totals.

For long runs on a workspace with little disk space, use a dedicated cache on the
same filesystem. Hard links reuse existing dependencies without duplicating them;
the runner prunes only new, older objects in this private cache periodically and every 25 results:

```sh
cp -al /workspace/go-build-cache /workspace/test262-go-cache
GOCACHE=/workspace/test262-go-cache \
TEST262_PRIVATE_GOCACHE=/workspace/test262-go-cache \
python3 tools/test262/run.py --all --jobs 2 --resume
```

Go 1.27 cache entries can be directories; the janitor handles those and preserves
every entry present in the original seed cache.

The full selection traverses categories in round-robin order so early checkpoints
cover more language areas. Resume retains completed results and runs the rest.

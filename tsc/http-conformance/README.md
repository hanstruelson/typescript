# HTTP conformance probes

The repository includes TypeScript HTTP/1 and HTTP/2 fixtures, scripts to emit
and build them, and repeatable protocol probes. The probes exercise emitted
server binaries and remain outside the shipped runtime.

The current recommendation is a layered check:

- **HTTP/1.1:** [Http11Probe](https://github.com/MDA2AV/Http11Probe), pinned
  here to commit `47883da3f4018ed7900d908e9133d249495c0104`. It exercises raw
  TCP framing, malformed input, ambiguity, normalization, and request-smuggling
  cases against RFC 9110 and RFC 9112. It requires the .NET 10 SDK. Its Pass,
  Warn, and Fail results are useful evidence, not an official conformance
  certificate.
- **HTTP/2:** [h2spec](https://github.com/summerwind/h2spec) v2.6.0 runs the
  HTTP/2, HPACK, and generic suites. Run it with strict mode enabled. This
  release is from 2020 and tests RFC 7540/7541; RFC 9113 now supersedes RFC 7540,
  so treat h2spec as an older broad probe, not complete current-spec coverage.
  The runner also uses Go's `golang.org/x/net/http2` client for ordinary GET and
  POST interoperability. That checks practical client/server interoperation,
  not frame-level conformance.
- **HTTP/0.9:** a small raw-socket probe checks the headerless `GET /\r\n`
  request shape, close behavior, and Node-compatible HTTP/1.1 response framing.

There is no single official HTTP/1.1 or HTTP/2 conformance suite that proves a
server correct. The RFCs remain authoritative. Keep failures tied to a named
test and RFC section, then add focused regression cases to the implementation's
own test harness when one is available. HTTP/2 work should add RFC 9113-focused
cases alongside h2spec rather than treating its older RFC 7540 results as
complete coverage.

## Prepare the tools

Build the server fixtures with the repository's supported Go toolchain:

```sh
GO_BIN=/path/to/go ./http-conformance/build-fixtures.sh
```

The default output directory is `/tmp/tsc-http-conformance`. Start each server
in its own terminal:

```sh
HTTP1_PORT=8080 /tmp/tsc-http-conformance/http1_server
HTTP2_PORT=8081 /tmp/tsc-http-conformance/http2_server
```

Run the curl matrix for HTTP/1.0, HTTP/1.1, cleartext HTTP/2, GET, POST, and
the raw HTTP/0.9 request:

```sh
HTTP1_PORT=8080 HTTP2_PORT=8081 ./http-conformance/run.sh --curl-only
```

This needs a curl build with HTTP/2 support and Python 3. The probe checks the
selected protocol, status, response body, and request-body length.

The repository's HTTP parser regressions are standard Go tests and run before
the external probes. They require the Go version in the repository's `go.mod`
(Go 1.27). Set `GO_BIN` if it is not named `go` in your `PATH`. To run only
those in-repository tests:

```sh
GO_BIN=/path/to/go ./http-conformance/run.sh --unit-only
```

Install h2spec v2.6.0 from its tagged release and make the binary available as
`h2spec` (or set `H2SPEC_BIN`). Checkout the pinned Http11Probe revision and
install the .NET 10 SDK:

```sh
git clone https://github.com/MDA2AV/Http11Probe.git /tmp/Http11Probe
git -C /tmp/Http11Probe checkout 47883da3f4018ed7900d908e9133d249495c0104
```

The target must provide `GET /` and `POST /` with status 200 and a non-empty
body for h2spec and the interoperability probe. Point the probes at the actual
HTTP/1 and HTTP/2 listeners. For a local self-signed TLS certificate, explicitly
set `HTTP2_INSECURE=1`.

```sh
HTTP11PROBE_DIR=/tmp/Http11Probe \
HTTP1_PORT=8080 \
HTTP2_PORT=8443 \
HTTP2_TLS=1 \
HTTP2_INSECURE=1 \
./http-conformance/run.sh
```

Use `--http1-only` or `--http2-only` to run one protocol. The HTTP/2 listener
can be cleartext prior-knowledge (h2c) by leaving `HTTP2_TLS` unset.

## Coverage boundaries

These tools are black-box probes. Their passing results do not establish full
Node compatibility, complete RFC 9113 compliance, or application-level Express
correctness. In particular, h2spec predates RFC 9113, and a passing GET/POST
interoperability check says little about malformed frames, flow-control edge
cases, stream cancellation, or concurrency. Preserve probe output and add a
first-party HTTP/2 regression suite with RFC 9113 cases as the transport gains a
test fixture.

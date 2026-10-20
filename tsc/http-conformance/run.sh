#!/usr/bin/env bash
set -euo pipefail

host=${HTTP_TEST_HOST:-127.0.0.1}
path=${HTTP_TEST_PATH:-/}
script_dir=$(cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(dirname "$script_dir")
go_bin=${GO_BIN:-go}
run_http1=1
run_http2=1
run_units=1
run_curl=0

case "${1:-}" in
  "") ;;
  --http1-only) run_http2=0 ;;
  --http2-only) run_http1=0 ;;
  --unit-only) run_http1=0; run_http2=0 ;;
  --curl-only) run_units=0; run_http1=0; run_http2=0; run_curl=1 ;;
  *) printf 'usage: %s [--http1-only|--http2-only|--unit-only|--curl-only]\n' "$0" >&2; exit 2 ;;
esac

if (( run_curl )); then
  HTTP_TEST_HOST="$host" HTTP_TEST_PATH="$path" "$script_dir/curl-smoke.sh"
fi

if (( run_units )); then
  (cd "$repo_root" && "$go_bin" test ./runtime -run '^TestHTTP' -count=1)
fi

if (( run_http1 )); then
  : "${HTTP1_PORT:?Set HTTP1_PORT to the target HTTP/1 listener}"
  : "${HTTP11PROBE_DIR:?Set HTTP11PROBE_DIR to the pinned Http11Probe checkout}"
  command -v dotnet >/dev/null || { printf 'dotnet is required for Http11Probe (.NET 10 SDK)\n' >&2; exit 2; }
  expected_probe_commit=47883da3f4018ed7900d908e9133d249495c0104
  actual_probe_commit=$(git -C "$HTTP11PROBE_DIR" rev-parse HEAD 2>/dev/null || true)
  [[ "$actual_probe_commit" == "$expected_probe_commit" ]] || {
    printf 'Http11Probe must be checked out at %s (found %s)\n' "$expected_probe_commit" "${actual_probe_commit:-unknown}" >&2
    exit 2
  }
  dotnet --list-sdks | awk '$1 ~ /^10\./ { found = 1 } END { exit !found }' || {
    printf 'Http11Probe requires the .NET 10 SDK\n' >&2
    exit 2
  }
  project="$HTTP11PROBE_DIR/src/Http11Probe.Cli/Http11Probe.Cli.csproj"
  [[ -f "$project" ]] || { printf 'Http11Probe project not found: %s\n' "$project" >&2; exit 2; }

  probe_args=(--host "$host" --port "$HTTP1_PORT")
  if [[ -n "${HTTP11PROBE_OUTPUT:-}" ]]; then
    probe_args+=(--output "$HTTP11PROBE_OUTPUT")
  fi
  dotnet run --project "$project" --configuration Release -- "${probe_args[@]}"
  python3 "$script_dir/http09_probe.py" --host "$host" --port "$HTTP1_PORT" --path "$path"
fi

if (( run_http2 )); then
  : "${HTTP2_PORT:?Set HTTP2_PORT to the target HTTP/2 listener}"
  h2spec_bin=${H2SPEC_BIN:-h2spec}
  command -v "$h2spec_bin" >/dev/null || { printf 'h2spec v2.6.0 is required (set H2SPEC_BIN if needed)\n' >&2; exit 2; }

  interop_args=(--host "$host" --port "$HTTP2_PORT" --path "$path")
  h2spec_args=(--strict --host "$host" --port "$HTTP2_PORT")
  if [[ "${HTTP2_TLS:-0}" == 1 ]]; then
    interop_args+=(--tls)
    h2spec_args+=(--tls)
    if [[ "${HTTP2_INSECURE:-0}" == 1 ]]; then
      interop_args+=(--insecure)
      h2spec_args+=(--insecure)
    fi
  fi
  if [[ -n "${H2SPEC_JUNIT:-}" ]]; then
    h2spec_args+=(--junit-report "$H2SPEC_JUNIT")
  fi
  (cd "$repo_root" && "$go_bin" run ./http-conformance/http2_interop.go "${interop_args[@]}")
  "$h2spec_bin" "${h2spec_args[@]}" http2 hpack generic
fi

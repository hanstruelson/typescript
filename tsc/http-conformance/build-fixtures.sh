#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "$0")" && pwd)
repo_root=$(dirname "$script_dir")
go_bin=${GO_BIN:-go}
build_dir=${HTTP_CONFORMANCE_BUILD_DIR:-${TMPDIR:-/tmp}/tsc-http-conformance}

mkdir -p -- "$build_dir"
build_dir=$(cd -- "$build_dir" && pwd)

for protocol in http1 http2; do
  source_file="$script_dir/fixtures/${protocol}_server.ts"
  output_dir="$build_dir/$protocol"
  mkdir -p -- "$output_dir"
  (cd "$repo_root" && "$go_bin" run ./cmd/tsc --target go --outDir "$output_dir" "$source_file")
  (cd "$repo_root" && "$go_bin" build -o "$build_dir/${protocol}_server" "$output_dir/${protocol}_server.go")
done

printf 'HTTP/1 server: %s\n' "$build_dir/http1_server"
printf 'HTTP/2 server: %s\n' "$build_dir/http2_server"

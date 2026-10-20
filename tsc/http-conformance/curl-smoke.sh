#!/usr/bin/env bash
set -euo pipefail

script_dir=$(cd -- "$(dirname -- "$0")" && pwd)
host=${HTTP_TEST_HOST:-127.0.0.1}
path=${HTTP_TEST_PATH:-/}
: "${HTTP1_PORT:?Set HTTP1_PORT to the HTTP/1 fixture listener}"
: "${HTTP2_PORT:?Set HTTP2_PORT to the HTTP/2 fixture listener}"

if [[ "$path" != /* || "$path" == *$'\r'* || "$path" == *$'\n'* ]]; then
  printf 'HTTP_TEST_PATH must be an origin-form path without line breaks\n' >&2
  exit 2
fi
command -v curl >/dev/null || { printf 'curl is required\n' >&2; exit 2; }
command -v python3 >/dev/null || { printf 'python3 is required for the HTTP/0.9 probe\n' >&2; exit 2; }
curl --version | grep -q 'Features:.*HTTP2' || {
  printf 'curl must be built with HTTP/2 support\n' >&2
  exit 2
}

url_host=$host
if [[ "$host" == *:* && "$host" != \[*\] ]]; then
  url_host="[$host]"
fi
body='curl-protocol-smoke-body'
body_bytes=$(LC_ALL=C printf '%s' "$body" | wc -c | tr -d '[:space:]')

assert_response() {
  local label=$1 expected=$2
  shift 2
  local actual
  actual=$(curl --silent --show-error --fail --max-time 5 --write-out $'\n%{http_version} %{http_code}' "$@") || {
    printf 'FAIL %s: curl request failed\n' "$label" >&2
    return 1
  }
  if [[ "$actual" != "$expected" ]]; then
    printf 'FAIL %s\nexpected: %q\nactual:   %q\n' "$label" "$expected" "$actual" >&2
    return 1
  fi
  printf 'PASS %s\n' "$label"
}

http1_url="http://${url_host}:${HTTP1_PORT}${path}"
http2_url="http://${url_host}:${HTTP2_PORT}${path}"

assert_response 'HTTP/1.0 GET' "GET $path version=1.0 bytes=0"$'\n1.1 200' --http1.0 "$http1_url"
assert_response 'HTTP/1.1 GET' "GET $path version=1.1 bytes=0"$'\n1.1 200' --http1.1 "$http1_url"
assert_response 'HTTP/1.1 POST' "POST $path version=1.1 bytes=$body_bytes"$'\n1.1 200' --http1.1 -X POST --data "$body" "$http1_url"
assert_response 'HTTP/2 GET (h2c prior knowledge)' "bytes=0 path=$path"$'\n2 200' --http2-prior-knowledge "$http2_url"
assert_response 'HTTP/2 POST (h2c prior knowledge)' "bytes=$body_bytes path=$path"$'\n2 200' --http2-prior-knowledge -X POST --data "$body" "$http2_url"
python3 "$script_dir/http09_probe.py" --host "$host" --port "$HTTP1_PORT" --path "$path"

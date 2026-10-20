#!/usr/bin/env python3
"""Probe the legacy HTTP/0.9 request form against an HTTP/1 server."""

import argparse
import socket
import sys


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--host", default="127.0.0.1")
    parser.add_argument("--port", type=int, required=True)
    parser.add_argument("--path", default="/")
    parser.add_argument("--timeout", type=float, default=5.0)
    args = parser.parse_args()

    if not args.path.startswith("/") or any(char in args.path for char in "\r\n\x00"):
        parser.error("--path must be an origin-form path without control characters")

    request = f"GET {args.path}\r\n".encode("ascii")
    response = bytearray()
    try:
        with socket.create_connection((args.host, args.port), timeout=args.timeout) as connection:
            connection.settimeout(args.timeout)
            connection.sendall(request)
            while True:
                chunk = connection.recv(64 * 1024)
                if not chunk:
                    break
                response.extend(chunk)
    except (OSError, UnicodeEncodeError) as error:
        print(f"HTTP/0.9 probe failed: {error}", file=sys.stderr)
        return 1

    header_end = response.find(b"\r\n\r\n")
    if header_end < 0:
        print("HTTP/0.9 probe failed: response did not contain an HTTP status/header block", file=sys.stderr)
        return 1
    header_lines = bytes(response[:header_end]).split(b"\r\n")
    status = header_lines[0]
    headers = {}
    for line in header_lines[1:]:
        name, separator, value = line.partition(b":")
        if separator:
            headers[name.strip().lower()] = value.strip().lower()

    if not status.startswith(b"HTTP/1.1 200 "):
        print(f"HTTP/0.9 probe failed: unexpected response status {status!r}", file=sys.stderr)
        return 1
    if headers.get(b"connection") != b"close":
        print("HTTP/0.9 probe failed: expected Connection: close", file=sys.stderr)
        return 1
    if b"chunked" in headers.get(b"transfer-encoding", b""):
        print("HTTP/0.9 probe failed: response must not use chunked framing", file=sys.stderr)
        return 1
    if len(response) == header_end + 4:
        print("HTTP/0.9 probe failed: target returned an empty response body", file=sys.stderr)
        return 1

    print(f"HTTP/0.9 probe passed: {status.decode('ascii', errors='replace')}")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())

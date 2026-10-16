#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
generated_dir="${script_dir}/generated"
mkdir -p "${generated_dir}"
download_dir="$(mktemp -d "${generated_dir}/.sqlite-amalgamation.XXXXXX")"
trap 'rm -rf "${download_dir}"' EXIT

base_url="https://raw.githubusercontent.com/mattn/go-sqlite3/v1.14.24"

fetch_and_verify() {
	local filename="$1"
	local expected_sha256="$2"
	local target="${download_dir}/${filename}"

	curl --fail --location --silent --show-error "${base_url}/${filename}" --output "${target}"
	if command -v sha256sum >/dev/null 2>&1; then
		printf '%s  %s\n' "${expected_sha256}" "${target}" | sha256sum --check --status
	elif command -v shasum >/dev/null 2>&1; then
		printf '%s  %s\n' "${expected_sha256}" "${target}" | shasum --algorithm 256 --check --status
	else
		echo "sha256sum or shasum is required to verify downloaded SQLite sources" >&2
		return 1
	fi
}

fetch_and_verify sqlite3-binding.c 955d1499fc1ca2b5419065a7be5a7f0f8e05291e9f04296796e9ab788129ef09
fetch_and_verify sqlite3-binding.h ae1e3c0b7832b7ce4fdf5b688c9cb94b719156b807cfca3a83db93330af8bd88
fetch_and_verify sqlite3ext.h 3dd4cdc4d5353b73d8c2db19d2dcc5b825ac2fb374fd0ab9a81720d82ec46428

mv "${download_dir}/sqlite3-binding.c" "${generated_dir}/sqlite3-binding.c"
mv "${download_dir}/sqlite3-binding.h" "${generated_dir}/sqlite3-binding.h"
mv "${download_dir}/sqlite3ext.h" "${generated_dir}/sqlite3ext.h"
echo "Installed the go-sqlite3 v1.14.24 SQLite amalgamation in ${generated_dir}"

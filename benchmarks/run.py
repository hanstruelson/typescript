#!/usr/bin/env python3
"""Build once, check results, and time complete processes in rotating order."""
import argparse
import datetime
import json
import math
import os
from pathlib import Path
import platform
import shutil
import statistics
import subprocess
import time

ROOT = Path(__file__).resolve().parent.parent
HERE = ROOT / "benchmarks"
BUILD = HERE / ".build"


def run(command, **kwargs):
    return subprocess.run(command, cwd=ROOT, check=True, text=True,
                          capture_output=True, **kwargs)


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("--samples", type=int, default=5)
    parser.add_argument("--cpus", type=int, default=8)
    parser.add_argument("--output", type=Path, default=HERE / "results.json")
    args = parser.parse_args()
    if args.samples < 1 or args.cpus < 1:
        parser.error("samples and cpus must be positive")
    BUILD.mkdir(exist_ok=True)
    go = shutil.which("go")
    node = shutil.which("node")
    bun = shutil.which("bun")
    if not bun:
        local_bun = BUILD / "tools" / "node_modules" / ".bin" / "bun"
        if local_bun.exists():
            bun = str(local_bun)
    if not go or not node or not bun:
        parser.error("Go 1.27+, Node 24+, and Bun must be installed (see benchmarks/README.md)")
    compiler = BUILD / "tsc"
    print("Building the compiler...", flush=True)
    run([go, "build", "-o", str(compiler), "./tsc/cmd/tsc"])
    arrays = (HERE / "arrays.ts").read_text()
    workers = (HERE / "workers.ts").read_text()
    compact = (HERE / "compact-array.ts").read_text()
    sources = {
        "array-native": arrays.replace("ELEMENT_TYPE", "number"),
        "array-dynamic": arrays.replace("ELEMENT_TYPE", "any"),
        "compact-array": compact,
        "workers-sequential": workers.replace("go compute(seed)", "Promise.resolve(compute(seed))"),
        "workers-parallel": workers,
    }
    workloads = []
    for name, source in sources.items():
        directory = BUILD / name
        directory.mkdir(exist_ok=True)
        (directory / "main.ts").write_text(source)
        (directory / "tsconfig.json").write_text(json.dumps({
            "compilerOptions": {"target": "go", "outDir": "out", "strict": True,
                                "noEmitOnError": True}, "files": ["main.ts"]}))
        print(f"Building {name}...", flush=True)
        try:
            run([str(compiler), "--project", str(directory)])
            executable = directory / "main"
            run([go, "build", "-o", str(executable), str(directory / "out" / "main.go")])
        except subprocess.CalledProcessError as error:
            raise SystemExit(error.stdout + error.stderr) from error
        group = "compact" if name == "compact-array" else "array" if name.startswith("array") else "workers"
        workloads.append((name, [str(executable)], group))
    # Node strips standard TypeScript annotations without changing the loop bodies.
    # Both array variants erase to the same JS, so only one JS baseline is needed.
    node_arrays = BUILD / "arrays-node.ts"
    node_arrays.write_text(sources["array-native"])
    node_serial = BUILD / "workers-node.ts"
    node_serial.write_text(sources["workers-sequential"])
    workloads.append(("array-node", [node, str(node_arrays)], "array"))
    workloads.append(("array-bun", [bun, str(node_arrays)], "array"))
    compact_source = BUILD / "compact-node.ts"
    compact_source.write_text(compact)
    workloads.append(("compact-array-node", [node, str(compact_source)], "compact"))
    workloads.append(("compact-array-bun", [bun, str(compact_source)], "compact"))
    workloads.append(("workers-node-sequential", [node, str(node_serial)], "workers"))
    workloads.append(("workers-bun-sequential", [bun, str(node_serial)], "workers"))
    # Generate the worker body from the same source rather than maintain a second algorithm.
    worker_body = workers[:workers.index("async function main")]
    erase = "const fs=require('node:fs');const {stripTypeScriptTypes}=require('node:module');process.stdout.write(stripTypeScriptTypes(fs.readFileSync(0,'utf8')));"
    worker_body = run([node, "-e", erase], input=worker_body).stdout
    (BUILD / "worker-body.cjs").write_text(worker_body + "\nmodule.exports = { compute };\n")
    shutil.copyfile(HERE / "node-workers.cjs", BUILD / "node-workers.cjs")
    workloads.append(("workers-node-parallel", [node, str(BUILD / "node-workers.cjs")], "workers"))
    workloads.append(("workers-bun-parallel", [bun, str(BUILD / "node-workers.cjs")], "workers"))
    env = dict(os.environ, GOMAXPROCS=str(args.cpus))
    checksums = {}
    rows = {name: {"group": group, "samples_ms": [], "stdout": ""}
            for name, _, group in workloads}
    # One complete, untimed warm-up process per workload populates filesystem caches.
    # Each timed Node process still pays its own startup/type stripping/JIT costs.
    for sample in range(-1, args.samples):
        order = workloads if sample < 0 else workloads[sample % len(workloads):] + workloads[:sample % len(workloads)]
        for name, command, group in order:
            start = time.perf_counter()
            try:
                result = run(command, env=env, timeout=180)
            except subprocess.CalledProcessError as error:
                raise SystemExit(f"{name} failed:\n" + error.stdout + error.stderr) from error
            elapsed = (time.perf_counter() - start) * 1000
            value = float(result.stdout.strip())
            if group in checksums and not math.isclose(value, checksums[group], rel_tol=1e-12):
                raise RuntimeError(f"Checksum mismatch for {name}: {value} versus {checksums[group]}")
            checksums[group] = value
            rows[name]["stdout"] = result.stdout.strip()
            if sample >= 0:
                rows[name]["samples_ms"].append(elapsed)
            print(f"{name}: {elapsed:.2f} ms ({'warm-up' if sample < 0 else sample + 1})", flush=True)
    cpu = "unknown"
    if Path("/proc/cpuinfo").exists():
        cpu = next((line.split(":", 1)[1].strip() for line in Path("/proc/cpuinfo").read_text().splitlines()
                    if line.startswith("model name")), cpu)
    for row in rows.values():
        row["median_ms"] = statistics.median(row["samples_ms"])
    report = {
        "measured_at_utc": datetime.datetime.now(datetime.timezone.utc).isoformat(),
        "platform": platform.platform(), "cpu": cpu, "logical_cpus": os.cpu_count(),
        "go": run([go, "version"]).stdout.strip(), "node": run([node, "--version"]).stdout.strip(),
        "bun": run([bun, "--version"]).stdout.strip(),
        "gomaxprocs": args.cpus, "samples": args.samples,
        "cpu_quota": Path("/sys/fs/cgroup/cpu.max").read_text().strip() if Path("/sys/fs/cgroup/cpu.max").exists() else "unavailable",
        "method": "Complete process wall time; compilation excluded; startup and setup included; rotating order; checksum verified.",
        "workloads": rows,
    }
    args.output.write_text(json.dumps(report, indent=2) + "\n")
    print(f"Results: {args.output}")


if __name__ == "__main__":
    main()

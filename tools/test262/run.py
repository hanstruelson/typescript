#!/usr/bin/env python3
"""Run upstream Test262 against the Go emitter; never count compilation errors as passes."""
import argparse, collections, concurrent.futures, datetime, gzip, hashlib, itertools, json, os
from pathlib import Path
import re, shutil, subprocess, tempfile, threading, time
import yaml
ROOT = Path(__file__).resolve().parents[2]
MARK = '__TEST262_SUCCESS_68cb62__'
DYNAMIC = re.compile(r'\b(?:eval\s*\(|(?:new\s+)?(?:Function|AsyncFunction|GeneratorFunction|AsyncGeneratorFunction)\s*\()')

def command(args, timeout, cwd=None):
    start = time.monotonic()
    try:
        r = subprocess.run(args, capture_output=True, text=True, encoding='utf-8', errors='replace', cwd=cwd, timeout=timeout,
                           env=dict(os.environ, TZ='UTC', GOMAXPROCS='2'))
        return r.returncode, (r.stdout + r.stderr)[-12000:], round(time.monotonic()-start, 3)
    except subprocess.TimeoutExpired as e:
        def text(v): return v.decode(errors='replace') if isinstance(v, bytes) else (v or '')
        return 124, (text(e.stdout)+text(e.stderr))[-12000:], round(time.monotonic()-start, 3)

def category(path):
    parts = path.split('/')
    return '/'.join(parts[:3] if parts[0]=='language' else parts[:2])

def native_source_identity():
    digest=hashlib.sha256()
    for directory in ('tsc/internal/goemit','tsc/internal/compiler','tsc/cmd/test262emit'):
        for path in sorted((ROOT/directory).rglob('*.go')):
            digest.update(str(path.relative_to(ROOT)).encode());digest.update(path.read_bytes())
    return digest.hexdigest()

def inventory(suite):
    cases=[]
    for p in sorted((suite/'test').rglob('*.js')):
        rel=p.relative_to(suite/'test').as_posix()
        source=p.read_text(encoding='utf-8')
        match=re.search(r'/\*---(.*?)---\*/', source, re.S)
        if not match: continue  # fixtures are imported inputs, not standalone tests
        meta=yaml.load(match.group(1), Loader=getattr(yaml, 'CSafeLoader', yaml.SafeLoader)) or {}
        meta={key:meta[key] for key in ('flags','includes','negative','features') if key in meta}
        flags=meta.get('flags', [])
        modes=['module'] if 'module' in flags else ['raw'] if 'raw' in flags else ['strict'] if 'onlyStrict' in flags else ['sloppy'] if 'noStrict' in flags else ['sloppy','strict']
        excluded = ('explicit-dynamic-code' if DYNAMIC.search(source) or rel.startswith('built-ins/eval/') or rel.startswith('built-ins/Function/constructor/') else None)
        cases.append(dict(path=rel,category=category(rel),metadata=meta,modes=modes,excluded=excluded))
    return cases

def execute(case, mode, args):
    result=dict(path=case['path'],category=case['category'],mode=mode,features=case['metadata'].get('features',[]))
    meta=case['metadata']; flags=meta.get('flags',[]); negative=meta.get('negative',{})
    if mode=='module': return dict(result,status='runner-module-unimplemented',detail='Module graph/realm harness needs a dedicated adapter; not a pass.')
    includes=meta.get('includes',[])
    if any(x in includes for x in ['agent.js','detachArrayBuffer.js']) or '$262' in (args.suite/'test'/case['path']).read_text():
        return dict(result,status='runner-host-unimplemented',detail='Requires Test262 $262 realm/agent/GC/detachment host API; not a pass.')
    source=(args.suite/'test'/case['path']).read_text()
    prefix='"use strict";\n' if mode=='strict' else ''
    harness=[] if mode=='raw' or negative.get('phase') in ('parse','early') else ['sta.js','assert.js']+includes
    try:
        source=prefix+'\n'.join((args.suite/'harness'/x).read_text() for x in dict.fromkeys(harness))+'\n'+source
    except FileNotFoundError as e: return dict(result,status='runner-missing-include',detail=str(e))
    parse_negative=negative.get('phase') in ('parse','early')
    if 'async' in flags:
        source+='\nfunction $DONE(error){if(error){throw error;}console.log("'+MARK+'");}\n'
    elif not parse_negative:
        if negative.get('phase')=='runtime':
            # Harness precedes the try, avoiding block-scope changes to helper declarations.
            split=source.rfind('/*---')
            source=source[:split]+'\ntry {\n'+source[split:]+'\nthrow new Test262Error("Expected runtime '+negative['type']+'");\n} catch (__caught) { assert.sameValue(__caught.name,"'+negative['type']+'"); }\n'
        source+='\nconsole.log("'+MARK+'");\n'
    with tempfile.TemporaryDirectory(prefix='t262-') as folder:
        folder=Path(folder); js=folder/'input.js'; go=folder/'main.go'; exe=folder/'run'
        js.write_text(source)
        # Run the same unmodified upstream helpers in Node, not rewritten assertion substitutes.
        oracle_cmd=[args.node,'--check' if parse_negative else str(js)]
        if parse_negative: oracle_cmd.append(str(js))
        rc,out,duration=command(oracle_cmd,args.timeout)
        result['reference_seconds']=duration
        oracle_ok=(rc!=0 and 'SyntaxError' in out) if parse_negative else rc==0 and MARK in out
        result['reference_status']='pass' if oracle_ok else 'fail'
        if not oracle_ok: result['reference_detail']=out or 'Reference did not reach completion marker'
        rc,out,duration=command([str(args.emitter),str(js),str(go)],args.timeout)
        result['emit_seconds']=duration
        if parse_negative:
            if rc==3 and negative.get('type')=='SyntaxError': return dict(result,status='pass',detail='Expected syntactic rejection')
            return dict(result,status='negative-not-rejected' if rc==0 else 'emit-error',detail=out or 'Expected early SyntaxError was not reported')
        if rc: return dict(result,status='emit-timeout' if rc==124 else 'emit-error',detail=out)
        rc,out,duration=command([args.go,'build','-o',str(exe),str(go)],args.build_timeout,cwd=ROOT/'tsc')
        result['build_seconds']=duration
        if rc: return dict(result,status='build-timeout' if rc==124 else 'build-error',detail=out)
        rc,out,duration=command([str(exe)],args.timeout)
        result['run_seconds']=duration
        if rc: return dict(result,status='runtime-timeout' if rc==124 else 'runtime-fail',detail=out)
        if out.splitlines().count(MARK)!=1: return dict(result,status='missing-completion',detail=out)
        return dict(result,status='pass',detail='')

def write_report(args,cases,results,selected):
    counts=collections.Counter(r['status'] for r in results)
    groups=collections.defaultdict(list)
    for r in results: groups[r['category']].append(r)
    inventory_counts=collections.Counter(c['category'] for c in cases)
    excluded=collections.Counter(c['category'] for c in cases if c['excluded'])
    revision=subprocess.check_output(['git','-C',str(args.suite),'rev-parse','HEAD'],text=True).strip()
    summary=dict(suite='Test262',revision=revision,generated=datetime.datetime.now(datetime.timezone.utc).isoformat(),
                 native_source_sha256=args.native_source_sha256,source_tests=len(cases),excluded_dynamic_tests=sum(excluded.values()),selected_tests=selected,
                 completed_variants=len(results),planned_variants=sum(len(c['modes']) for c in args.selected_cases),complete=len(results)==sum(len(c['modes']) for c in args.selected_cases),statuses=dict(counts),reference_failures=sum(r.get('reference_status')=='fail' for r in results),filters=args.filter,selection='all' if args.all else f'{args.per_category} stable hash-selected tests per category',
                 native_commit=subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip(),
                 node_version=subprocess.check_output([args.node,'--version'],text=True).strip(),
                 go_version=subprocess.check_output([args.go,'version'],text=True).strip())
    (args.output/'summary.json').write_text(json.dumps(summary,indent=2)+'\n')
    lines=['# JavaScript conformance checklist','',f'Test262 revision: `{revision}`.',
           f"Inventory: {len(cases):,} tests; {sum(excluded.values()):,} explicitly excluded for dynamic code.",
           f"Execution: {len(results):,} variants from {selected:,} selected tests; selection: {summary['selection']}.",
           '', ('**Full-suite run: results below are checkpoints until every selected variant completes.**' if args.all else '**Initial stratified audit, not a full-suite compatibility score.**') + ' Untested cases are not passes. The default coercion policy is unchanged; failures may expose differences from ECMAScript.',
           '', 'Module and host-adapter limitations are tracked separately from compiler/runtime failures. Node reference failures are recorded separately and do not prevent native execution. Negative tests pass only on syntactic rejection, never on an arbitrary build failure.',
           '', '## Results','']
    lines += [f'- {key}: {value:,}' for key,value in sorted(counts.items())]
    lines += ['', '## Category checklist','', '| Category | Suite tests | Dynamic exclusions | Executed variants | Pass | Native failures | Adapter/reference limits |', '|---|---:|---:|---:|---:|---:|---:|']
    native=lambda r:r['status'] not in ('pass','reference-limitation') and not r['status'].startswith('runner-')
    for cat,total in sorted(inventory_counts.items()):
        rows=groups[cat]
        lines.append(f"| {cat} | {total} | {excluded[cat]} | {len(rows)} | {sum(r['status']=='pass' for r in rows)} | {sum(native(r) for r in rows)} | {sum(r['status']=='reference-limitation' or r['status'].startswith('runner-') for r in rows)} |")
    lines += ['', '## Failure checklist','']
    for cat,rows in sorted(groups.items()):
        failures=[r for r in rows if r['status']!='pass']
        if not failures: continue
        lines += [f'### {cat}','']
        for r in failures:
            detail=r['detail'].replace(str(args.suite),'<test262>')
            detail=re.sub(r'/tmp/t262-[^/\s]+','<case>',detail)
            lines += [f"- [ ] `{r['path']}` ({r['mode']}): **{r['status']}**",'','```text',detail[:1600] or '(no output)','```','']
    (args.output/'CHECKLIST.md').write_text('\n'.join(lines)+'\n')
    return summary

def trim_private_cache():
    # Go 1.27 uses directories for some cache entries. Preserve every seeded
    # entry; only remove additional entries from the dedicated task cache.
    cache=os.environ.get('TEST262_PRIVATE_GOCACHE')
    if not cache: return
    root=Path(cache)
    if os.environ.get('GOCACHE')!=str(root) or root.name!='test262-go-cache':
        raise RuntimeError('Private cache configuration does not match task cache')
    seed=Path(os.environ.get('TEST262_CACHE_SEED',str(root.parent/'go-build-cache')))
    if seed==root or not seed.is_dir(): raise RuntimeError('Missing independent cache seed')
    cutoff=time.time()-20
    for directory in root.iterdir():
        if not directory.is_dir(): continue
        for path in directory.iterdir():
            try:
                if (seed/path.relative_to(root)).exists(): continue
                if path.stat().st_mtime<cutoff:
                    if path.is_dir(): shutil.rmtree(path)
                    else: path.unlink()
            except FileNotFoundError: pass


def main():
    p=argparse.ArgumentParser(description=__doc__)
    p.add_argument('--suite',type=Path,default=ROOT.parent/'test262')
    p.add_argument('--output',type=Path,default=ROOT/'docs/conformance/latest')
    p.add_argument('--per-category',type=int,default=2)
    p.add_argument('--all',action='store_true')
    p.add_argument('--resume',action='store_true',help='Keep existing result rows and run unfinished variants')
    p.add_argument('--filter',action='append',default=[],help='Path substring; repeat to include multiple feature groups')
    p.add_argument('--jobs',type=int,default=4)
    p.add_argument('--timeout',type=float,default=15)
    p.add_argument('--build-timeout',type=float,default=90)
    p.add_argument('--emitter',type=Path,default=Path('/tmp/test262emit'))
    args=p.parse_args()
    if args.jobs<1 or args.per_category<1: p.error('jobs and per-category must be positive')
    args.suite=args.suite.resolve();args.output=args.output.resolve();args.emitter=args.emitter.resolve()
    args.go=shutil.which('go');args.node=shutil.which('node')
    if not args.go or not args.node: p.error('Go 1.27 and Node must be on PATH')
    args.output.mkdir(parents=True,exist_ok=True)
    if os.environ.get('TEST262_PRIVATE_GOCACHE'):
        def cache_janitor():
            while True:
                trim_private_cache()
                time.sleep(5)
        threading.Thread(target=cache_janitor,daemon=True).start()
    args.native_source_sha256=native_source_identity()
    cases=inventory(args.suite)
    with gzip.GzipFile(filename=str(args.output/'inventory.json.gz'),mode='wb',mtime=0) as output:
        output.write(json.dumps(cases,separators=(',',':')).encode())
    groups=collections.defaultdict(list)
    for c in cases:
        if not c['excluded'] and (not args.filter or any(fragment in c['path'] for fragment in args.filter)): groups[c['category']].append(c)
    selected=[]
    buckets=[]
    for cat,rows in sorted(groups.items()):
        rows.sort(key=lambda c:hashlib.sha256(c['path'].encode()).hexdigest())
        buckets.append(rows if args.all else rows[:args.per_category])
    selected=[c for row in itertools.zip_longest(*buckets) for c in row if c is not None]
    args.selected_cases=selected
    tasks=[(c,m) for c in selected for m in c['modes']]
    results=[]
    if args.resume and (args.output/'results.jsonl').exists():
        previous=json.loads((args.output/'summary.json').read_text())
        revision=subprocess.check_output(['git','-C',str(args.suite),'rev-parse','HEAD'],text=True).strip()
        if previous.get('native_source_sha256')!=args.native_source_sha256 or previous['revision']!=revision or previous['native_commit']!=subprocess.check_output(['git','rev-parse','HEAD'],cwd=ROOT,text=True).strip():
            p.error('Cannot resume across different suite or native revisions; use a new output directory')
        for line in (args.output/'results.jsonl').read_text().splitlines():
            try:
                row=json.loads(line)
                if row['status']!='reference-limitation': results.append(row)
            except json.JSONDecodeError: pass
        # Discard an interrupted trailing write before appending new complete rows.
        (args.output/'results.jsonl').write_text(''.join(json.dumps(r)+'\n' for r in results))
        done={(r['path'],r['mode']) for r in results}
        tasks=[(c,m) for c,m in tasks if (c['path'],m) not in done]
    write_report(args,cases,results,len(selected))
    total_variants=sum(len(c['modes']) for c in selected)
    print(f'Inventory {len(cases)} tests; {total_variants} planned variants, {len(tasks)} remaining, across {len(groups)} categories',flush=True)
    with (args.output/'results.jsonl').open('a' if args.resume else 'w') as log, concurrent.futures.ThreadPoolExecutor(max_workers=args.jobs) as pool:
        futures=[pool.submit(execute,c,m,args) for c,m in tasks]
        for future in concurrent.futures.as_completed(futures):
            r=future.result();results.append(r);log.write(json.dumps(r)+'\n');log.flush()
            if len(results)%25==0:
                trim_private_cache()
            if len(results)%100==0:
                print(f'{len(results)}/{total_variants}: {dict(collections.Counter(x["status"] for x in results))}',flush=True)
                write_report(args,cases,results,len(selected))
    summary=write_report(args,cases,results,len(selected))
    print(json.dumps(summary,indent=2),flush=True)
if __name__=='__main__': main()

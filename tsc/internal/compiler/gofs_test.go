package compiler_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"testing"
)

// Exercise the entire I/O boundary through emitted code, including callback
// identity, descriptor offsets, worker loops, and iterator cleanup.
func TestGoNodeFS(t *testing.T) {
	root := strconv.Quote(t.TempDir())
	source := `import fs from "node:fs";import * as fsp from "node:fs/promises";import { Buffer } from "node:buffer";
 const root=` + root + `;
 async function main(){
 const dir=fs.mkdtempSync(root+"/fs-");const file=dir+"/data";
 fs.writeFileSync(file,"hello");fs.appendFileSync(file,"!");
 console.log(fs.readFileSync(file,"utf8"),Buffer.isBuffer(fs.readFileSync(file)));
 const fd=fs.openSync(file,"r+");const bytes=Buffer.alloc(3);
 console.log(fs.readSync(fd,bytes,0,3,1),bytes.toString());
 console.log(fs.writeSync(fd,"OK",0,"utf8"));fs.closeSync(fd);
 const handle=await fsp.open(file,"r+");
 console.log((await handle.stat()).isFile(),(await handle.readFile("utf8")));
 await handle.close();await handle.close();
 await new Promise((resolve,reject)=>fs.readFile(file,(err,data)=>{if(err){reject(err);}else{console.log(data.toString());resolve(0);}}));
 fs.mkdirSync(dir+"/sub");fs.writeFileSync(dir+"/sub/a","a");
 await fsp.cp(dir+"/sub",dir+"/copy",{recursive:true,filter:(src,dest)=>true});
 console.log(fs.readdirSync(dir+"/copy")[0],fs.statSync(file).size);
 const entries=await fsp.opendir(dir+"/copy");for await(const entry of entries){console.log(entry.name,entry.isFile());break;}
 try{entries.readSync();}catch(e){console.log(e.code);}
 for await(const name of fsp.glob("**/a",{cwd:dir})){console.log(name);}
 fs.linkSync(file,dir+"/hard");fs.symlinkSync("data",dir+"/link");
 console.log(fs.lstatSync(dir+"/link").isSymbolicLink(),fs.readlinkSync(dir+"/link"));
 try{await fsp.readFile(dir+"/missing");}catch(e){console.log(e.code,e.path===dir+"/missing");}
 async function worker(path){const mod=await import("node:fs/promises");return await mod.readFile(path,"utf8");}
 console.log(await go worker(file));
 await fsp.rm(dir,{recursive:true});console.log(fs.existsSync(dir));
 }main();`
	text := emitGoProgram(t, source)

	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "hello! true\n3 ell\n2\ntrue OKllo!\nOKllo!\na 6\na true\nERR_DIR_CLOSED\ncopy/a\nsub/a\ntrue data\nENOENT true\nOKllo!\nfalse\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSWatch(t *testing.T) {
	source := `const fs=require("node:fs");const root=` + strconv.Quote(t.TempDir()) + `;
 async function main(){
 const path=root+"/watch";fs.writeFileSync(path,"before");
 await new Promise((resolve,reject)=>{const watcher=fs.watch(path,(event,filename)=>{console.log(event,filename);watcher.close();resolve(0);});watcher.on("error",reject);fs.writeFileSync(path,"after");});
 await new Promise((resolve,reject)=>{const watcher=fs.watchFile(path,{interval:5},(curr,prev)=>{console.log(curr.size,prev.size);fs.unwatchFile(path);resolve(0);});setTimeout(()=>fs.writeFileSync(path,"longer-value"),30);});
 const idle=fs.watch(path,{persistent:false});idle.unref();console.log("finished");
 }main();`
	text := emitGoProgram(t, source)
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != "change watch\n12 5\nfinished\n" {
		t.Fatalf("unexpected watcher output %q", got)
	}
}

func TestGoNodeFSStreams(t *testing.T) {
	source := `const fs=require("node:fs");const root=` + strconv.Quote(t.TempDir()) + `;
 async function main(){
 const path=root+"/stream";const writer=fs.createWriteStream(path,{highWaterMark:2});
 await new Promise((resolve,reject)=>{writer.on("error",reject);writer.on("close",()=>resolve(0));console.log(writer.write("ab"));writer.end("cdef");});
 console.log(writer.bytesWritten,fs.readFileSync(path,"utf8"));
 let total="";for await(const chunk of fs.createReadStream(path,{highWaterMark:2,start:1,end:4})){total+=chunk.toString();}console.log(total);
 await new Promise((resolve,reject)=>{const input=fs.createReadStream(path,{highWaterMark:2});const output=fs.createWriteStream(root+"/copy",{highWaterMark:2});input.on("error",reject);output.on("error",reject);output.on("close",()=>resolve(0));input.pipe(output);});
 console.log(fs.readFileSync(root+"/copy","utf8"));
 fs.writeFileSync(root+"/unicode","é🙂終");let decoded="";for await(const piece of fs.createReadStream(root+"/unicode",{highWaterMark:1,encoding:"utf8"})){decoded+=piece;}console.log(decoded);
 }main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != "false\n6 abcdef\nbcde\nabcdef\né🙂終\n" {
		t.Fatalf("unexpected stream output %q", got)
	}
}

func TestGoNodeFSDescriptorLifetime(t *testing.T) {
	source := `const fs=require("fs");const fsp=require("fs/promises");const Buffer=require("buffer").Buffer;const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 async function main(){
 fs.writeFileSync(file,Buffer.from("00ff807f","hex"));const fd=fs.openSync(file,"r+");const a=Buffer.alloc(2),b=Buffer.alloc(2);console.log(fs.readvSync(fd,[a,b],0),a.toString("hex"),b.toString("hex"));
 await new Promise((resolve,reject)=>fs.read(fd,a,0,2,0,(error,n,same)=>{if(error){reject(error);}else{console.log(n,same===a);resolve(0);}}));fs.closeSync(fd);
 const handle=await fsp.open(file,"r");const reading=handle.readFile();const closing=handle.close();console.log((await reading).toString("hex"));await closing;console.log(handle.fd);
 await fsp.cp(file,file+".copy",{filter:async(src,dest)=>{await 0;return true;}});console.log(fs.readFileSync(file+".copy").equals(fs.readFileSync(file)));
 try{fs.copyFileSync(file,file);}catch(e){console.log(e.code);}
 try{fs.readSync(fd,a,0,1,0);}catch(e){console.log(e.code);}
 console.log(fs.realpathSync.native(file)===fs.realpathSync(file),require("fs")===fs,require("buffer")===require("node:buffer"));
 }main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "4 00ff 807f\n2 true\n00ff807f\n-1\ntrue\nEINVAL\nEBADF\ntrue true true\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

// A shared JavaScript fixture checks behavior against Node, not just against
// assumptions encoded in the Go implementation. Node is optional for contributors.
func TestGoNodeFSDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node reference executable unavailable")
	}
	source := `const fs=require("fs");const fsp=require("fs/promises");const root=` + strconv.Quote(t.TempDir()) + `;
 async function main(){const file=root+"/file";fs.writeFileSync(file,"hello");console.log(fs.readFileSync(file).toString(),fs.statSync(file).isFile());const handle=await fsp.open(file,"r+");const buffer=Buffer.alloc(3);const result=await handle.read(buffer,0,3,1);console.log(result.bytesRead,result.buffer===buffer,buffer.toString());await handle.close();console.log(handle.fd);const out=fs.createWriteStream(root+"/stream",{highWaterMark:1});await new Promise((resolve,reject)=>{out.on("error",reject);out.on("close",resolve);out.end("é🙂終");});let text="";for await(const chunk of fs.createReadStream(root+"/stream",{encoding:"utf8",highWaterMark:1})){text+=chunk;}console.log(text);try{await fsp.readFile(root+"/absent");}catch(e){console.log(e.code,e.path===root+"/absent");}}main();`
	path := filepath.Join(t.TempDir(), "reference.cjs")
	if err := os.WriteFile(path, []byte(source), 0600); err != nil {
		t.Fatal(err)
	}
	want, err := exec.Command(node, path).CombinedOutput()
	if err != nil {
		t.Fatalf("Node reference failed: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("Go %q differs from Node %q", got, want)
	}
}

func TestGoNodeFSAbort(t *testing.T) {
	source := `const fs=require("fs");const fsp=require("fs/promises");const root=` + strconv.Quote(t.TempDir()) + `;
 async function main(){
 const file=root+"/data";fs.writeFileSync(file,"original");const controller=new AbortController();let events=0;controller.signal.addEventListener("abort",()=>events++,{once:true});controller.abort("stop");controller.abort("again");console.log(controller.signal.aborted,controller.signal.reason,events);
 try{controller.signal.throwIfAborted();}catch(e){console.log(e);}
 try{await fsp.readFile(file,{signal:controller.signal});}catch(e){console.log(e.name,e.code,e.cause);}
 try{await fsp.writeFile(file,"changed",{signal:controller.signal});}catch(e){console.log(e.code,fs.readFileSync(file,"utf8"));}
 await new Promise((resolve,reject)=>fs.readFile(file,{signal:controller.signal},err=>{console.log(err.code);resolve(0);}));
 const pre=fsp.watch(root+"/missing",{signal:controller.signal});try{await pre.next();}catch(e){console.log(e.code,e.cause);}console.log((await pre.next()).done);
 const watchController=new AbortController();const watching=fsp.watch(file,{signal:watchController.signal});setTimeout(()=>watchController.abort("watch-stop"),10);try{await watching.next();}catch(e){console.log(e.code,e.cause);}
 const timed=fsp.watch(file,{signal:AbortSignal.timeout(10)});try{await timed.next();}catch(e){console.log(e.name,e.code,e.cause.name);}
 const streamController=new AbortController();await new Promise((resolve,reject)=>{const reader=fs.createReadStream(file,{signal:streamController.signal});reader.on("error",error=>console.log(error.name,error.code,error.cause));reader.on("close",()=>resolve(0));streamController.abort("stream-stop");});
 }main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true stop 1\nstop\nAbortError ABORT_ERR stop\nABORT_ERR original\nABORT_ERR\nABORT_ERR stop\ntrue\nABORT_ERR watch-stop\nAbortError ABORT_ERR TimeoutError\nAbortError ABORT_ERR stream-stop\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSBigInt(t *testing.T) {
	source := `const fs=require("fs");const fsp=require("fs/promises");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 async function main(){fs.writeFileSync(file,"four");const stat=fs.statSync(file,{bigint:true});console.log(typeof stat.size,stat.size===4n,typeof stat.ino,typeof stat.mtimeNs,stat.mtimeNs/1000000n===stat.mtimeMs,stat.birthtimeNs>0n);const handle=await fsp.open(file,"r");console.log((await handle.stat({bigint:true})).size===4n);await handle.close();const disk=fs.statfsSync(file,{bigint:true});console.log(typeof disk.blocks,disk.blocks>0n);let big=9007199254740993n;console.log(big>9007199254740992,big==9007199254740992,big===9007199254740992,big.toString(16));console.log(0n==null,0n==false,BigInt("010"),BigInt.asIntN(8,255n),BigInt.asUintN(8,-1n));console.log(7n/2n,-7n%2n,1n<<65n,~0n);big++;console.log(big);console.log(big.int64());try{console.log(1n+1);}catch(e){console.log(e.name);}try{console.log(+1n);}catch(e){console.log(e.name);}console.log(String(123n),Number(12n));}main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "bigint true bigint bigint true true\ntrue\nbigint true\ntrue false false 20000000000001\nfalse true 10 -1 255\n3 -1 36893488147419103232 -1\n9007199254740994\n9007199254740994\nTypeError\nTypeError\n123 12\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSHandleStreamLifetime(t *testing.T) {
	source := `import fs = require("node:fs");import fsp = require("node:fs/promises");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 async function main(){fs.writeFileSync(file,"abcdef");const handle=await fsp.open(file,"r");const reader=handle.createReadStream({highWaterMark:2});const closing=handle.close();let total="";for await(const chunk of reader){total+=chunk.toString();}await closing;await handle.close();console.log(total,handle.fd);const second=await fsp.open(file,"r");const retained=second.createReadStream({autoClose:false});for await(const chunk of retained){console.log(chunk.toString());}console.log(second.fd>=0);await second.close();}main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "abcdef -1\nabcdef\ntrue\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSBufferMethods(t *testing.T) {
	source := `const fs=require("fs");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 const value=Buffer.alloc(40);value.writeUIntLE(0x123456,0,3);value.writeInt16BE(-1234,3);value.writeFloatLE(1.5,5);value.writeDoubleBE(-0.125,9);value.writeBigUInt64LE(18446744073709551615n,17);value.writeBigInt64BE(-123n,25);
 fs.writeFileSync(file,value);const read=fs.readFileSync(file);console.log(read.readUIntLE(0,3),read.readInt16BE(3),read.readFloatLE(5),read.readDoubleBE(9),read.readBigUInt64LE(17).toString(),read.readBigInt64BE(25).toString());
 console.log(Buffer.from([1,2,3,4]).swap16().toString("hex"));console.log(Buffer.from("abcdef").toString("utf8",-1,3));console.log(Buffer.from("abcdef").indexOf("cd"),Buffer.from("ababa").lastIndexOf("ba"),Buffer.from("abc").includes(98));
 const json=Buffer.from("bytes").toJSON();console.log(json.type,Buffer.from(json).toString(),Buffer.isEncoding("UTF-8"),Buffer.isEncoding("buffer"));console.log(Buffer.copyBytesFrom(new Uint16Array([0x1234])).toString("hex"));
 try{read.readInt32LE(38);}catch(e){console.log(e.code);}try{read.writeUInt8(256);}catch(e){console.log(e.code);}console.log(Buffer.alloc(2,"").toString("hex"));`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "1193046 -1234 1.5 -0.125 18446744073709551615 -123\n02010403\nabc\n2 3 true\nBuffer bytes true false\n3412\nERR_OUT_OF_RANGE\nERR_OUT_OF_RANGE\n0000\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSNativeClasses(t *testing.T) {
	source := `const fs=require("fs");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 fs.writeFileSync(file,"abc");const stat=fs.statSync(file);console.log(stat instanceof fs.Stats,stat.constructor.name,Object.getPrototypeOf(stat)===fs.Stats.prototype,fs.Stats.prototype.isFile.call(stat));const detached=stat.isFile;try{detached();}catch(e){console.log(e.name);}console.log(detached.bind(stat)());stat.mode=0;console.log(stat.isFile());const entry=new fs.Dirent("a",1,"root");console.log(entry instanceof fs.Dirent,entry.name,entry.parentPath,entry.isFile());const reader=new fs.ReadStream(file);console.log(reader instanceof fs.ReadStream,fs.FileReadStream===fs.ReadStream);reader.destroy();const buffer=new Buffer("test");console.log(buffer instanceof Buffer,buffer instanceof Uint8Array,buffer.constructor===Buffer,typeof Buffer,Buffer.prototype.toString.call(buffer),Buffer.prototype.readUInt8.call(buffer));console.log(fs.statSync(file,{bigint:true}) instanceof fs.Stats);`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true Stats true true\nTypeError\ntrue\nfalse\ntrue a root true\ntrue true\ntrue true true function test 116\nfalse\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSDateMetadata(t *testing.T) {
	source := `const fs=require("fs");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 fs.writeFileSync(file,"x");const date=new Date("2020-01-02T03:04:05.678Z");fs.utimesSync(file,date,date);const stat=fs.statSync(file);console.log(stat.mtime instanceof Date,stat.mtime.getTime(),stat.mtime.toISOString(),stat.mtime.getUTCFullYear(),stat.mtime===stat.mtime,Number(stat.mtime));console.log(Object.hasOwn(stat,"mtime"),Object.hasOwn(stat,"isFile"));const empty=new fs.Stats();console.log(empty.isFile(),empty.atime.getTime(),empty.atime.toString(),empty.atime.toJSON());try{empty.atime.toISOString();}catch(e){console.log(e.name);}console.log(new Date(0).getTime(),Date.UTC(2020,0,2,3,4,5,678));`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true 1577934245678 2020-01-02T03:04:05.678Z 2020 true 1577934245678\nfalse false\nfalse NaN Invalid Date null\nRangeError\n0 1577934245678\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSBlob(t *testing.T) {
	source := `const fs=require("fs");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 async function main(){fs.writeFileSync(file,"hello é🙂");const blob=await fs.openAsBlob(file,{type:"TEXT/PLAIN"});console.log(blob instanceof Blob,blob.size,blob.type,await blob.text());console.log(await blob.slice(0,5).text());const bytes=await blob.bytes();console.log(bytes instanceof Uint8Array,Buffer.isBuffer(bytes),bytes.length);const stream=blob.stream();console.log(stream instanceof ReadableStream,stream.locked);const reader=stream.getReader();console.log(stream.locked,(await reader.read()).value.length,(await reader.read()).done);await reader.closed;reader.releaseLock();console.log(stream.locked);let count=0;for await(const chunk of blob.stream()){count+=chunk.length;}console.log(count);const immutable=new Blob([Buffer.from("abc"),"def"],{type:"TeXT/X"});console.log(immutable.size,immutable.type,await immutable.text());const changed=await fs.openAsBlob(file);fs.writeFileSync(file,"modified");try{await changed.text();}catch(e){console.log(e.name,e.code);}}main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "true 12 text/plain hello é🙂\nhello\ntrue false 12\ntrue false\ntrue 12 true\nfalse\n12\n6 text/x abcdef\nNotReadableError 0\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNodeFSReadLines(t *testing.T) {
	source := `const fs=require("fs");const fsp=require("fs/promises");const file=` + strconv.Quote(filepath.Join(t.TempDir(), "data")) + `;
 async function main(){fs.writeFileSync(file,"one\r\né🙂\rtwo\n\nlast");const handle=await fsp.open(file,"r");let lines=[];for await(const line of handle.readLines({highWaterMark:1})){lines.push(line);}console.log(lines.length,lines.join("|"));await handle.close();console.log(handle.fd);const second=await fsp.open(file,"r");for await(const line of second.readLines({highWaterMark:1})){console.log(line);break;}await second.close();console.log(second.fd);}main();`
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "5 one|é🙂|two||last\n-1\none\n-1\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

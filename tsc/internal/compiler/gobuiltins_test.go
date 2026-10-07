package compiler_test

import (
	"os/exec"
	"testing"
)

func TestGoNativeObjectJSON(t *testing.T) {
	source := `const object:any={a:1};let stored=2;
 const receiver:any={value:5,get doubled(){return this.value*2;},set doubled(value){this.value=value/2;},read(){return this.value;}};receiver.doubled=12;console.log(receiver.doubled,receiver.read(),receiver.read.call({value:9}));
 Object.defineProperty(object,"hidden",{value:7});
 Object.defineProperty(object,"access",{get:()=>stored,set:(value)=>{stored=value;},enumerable:true,configurable:true});
 object.access=4;console.log(Object.keys(object).join(","),object.access,Object.hasOwn(object,"access"));
 console.log(JSON.stringify(object));console.log(JSON.stringify([undefined,NaN,Infinity,"\ud800"]));
 const parsed=JSON.parse('{"__proto__":3,"2":2,"a":1,"a":9}');console.log(JSON.stringify(parsed));
 console.log(JSON.stringify({b:2,a:1},["a","b","a"],2));
 console.log(JSON.stringify(JSON.parse('{"a":1,"b":2}',(key,value)=>key==="a"?undefined:value)));
 console.log(JSON.stringify({date:new Date("2024-02-29T00:00:00Z")}));
 console.log(JSON.stringify({raw:JSON.rawJSON("12345678901234567890")}),JSON.isRawJSON(JSON.rawJSON("null")));
 const ancestor:any={x:8};const child:any=Object.create(ancestor);console.log(child.x,"x" in child,Object.hasOwn(child,"x"));child.x=9;console.log(child.x,ancestor.x);
 console.log(JSON.stringify({a:1,b:2},function(key,value){if(key==="a"){return this.b;}return value;}));
 const atomic:any={};try{Object.defineProperties(atomic,{first:{value:1},second:{get:2}});}catch(e){console.log(e.name,Object.hasOwn(atomic,"first"));}
 const inherited:any=Object.create(receiver);console.log(inherited.doubled);inherited.doubled=20;console.log(inherited.value,receiver.value);
 Object.defineProperty(Date.prototype,"custom",{get:function(){return this.getTime();},configurable:true});console.log(new Date(123).custom);
Object.freeze(object);console.log(Object.isFrozen(object),Object.isExtensible(object),Object.getOwnPropertyDescriptor(object,"a").writable);
 try{object.a=2;}catch(e){console.log(e.name);}
 try{JSON.parse("01");}catch(e){console.log(e.name);}
 const circular:any={};circular.self=circular;try{JSON.stringify(circular);}catch(e){console.log(e.name);}`
	text := emitGoProgram(t, source)
	got, err := runGoProgram(t, text, true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	want := "12 6 9\na,access 4 true\n{\"a\":1,\"access\":4}\n[null,null,null,\"\\ud800\"]\n{\"2\":2,\"__proto__\":3,\"a\":9}\n{\n  \"a\": 1,\n  \"b\": 2\n}\n{\"b\":2}\n{\"date\":\"2024-02-29T00:00:00.000Z\"}\n{\"raw\":12345678901234567890} true\n8 true false\n9 8\n{\"a\":2,\"b\":2}\nTypeError false\n12\n10 6\n123\ntrue false false\nTypeError\nSyntaxError\nTypeError\n"
	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestGoNativeDateDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node not installed")
	}
	source := `for(const text of ["2024","2024-02","2024-02-29","2024-02-30T24:00:00Z","2024-03-10T02:30:00","2024-11-03T01:30:00","+010000-01-01T00:00:00Z","-000001-01-01T00:00:00Z","2024-01-01T00:00:00.1234567890123Z","2024T00:00","2024-01T00:00","2024-01-01T00:00:00,1Z","bad"]){const date=new Date(text);console.log(date.getTime(),date.toJSON());}
 console.log(Date.UTC(2024),Date.UTC(2024,2,10,2,30),Date.UTC(99,0,1));
 const date=new Date("2024-01-31T12:34:56.789Z");console.log(date.setUTCMonth(1),date.toISOString());console.log(date.setUTCHours(25,2),date.toISOString());console.log(date.toUTCString(),date.toGMTString());
 console.log(Date.UTC(2024,0,1,0,0,0,1000000000000));const wide=new Date(0);console.log(wide.setUTCSeconds(100000000000),wide.toISOString());
 const gap=new Date(2024,2,10,2,30);console.log(gap.getTime(),gap.getHours());const fold=new Date(2024,10,3,1,30);console.log(fold.getTime());
 const invalid=new Date(NaN);invalid.setFullYear(2000);console.log(invalid.getTime());
 for(const method of ["getFullYear","getUTCFullYear","getMonth","getUTCMonth","getDate","getUTCDate","getDay","getUTCDay","getHours","getUTCHours","getMinutes","getUTCMinutes","getSeconds","getUTCSeconds","getMilliseconds","getUTCMilliseconds","getTimezoneOffset","getYear","valueOf","getTime"]){const item:any=new Date("2024-02-29T12:34:56.789Z");console.log(method,item[method]());}
 for(const method of ["setDate","setUTCDate","setMonth","setUTCMonth","setFullYear","setUTCFullYear","setHours","setUTCHours","setMinutes","setUTCMinutes","setSeconds","setUTCSeconds","setMilliseconds","setUTCMilliseconds","setYear","setTime"]){const item:any=new Date("2024-02-29T12:34:56.789Z");console.log(method,item[method](14),item.toJSON());}
 for(const value of [8640000000000000,-8640000000000000,8640000000000001,-8640000000000001,-0,0.9,-0.9]){const clipped=new Date(value);console.log(clipped.getTime(),Object.is(clipped.getTime(),-0),clipped.toJSON());}
 const equivalent=new Date(0);console.log(equivalent==equivalent.toString(),equivalent==0,equivalent+"",equivalent-0);
 const epoch=new Date(0);console.log(epoch.toString(),epoch.toDateString(),epoch.toTimeString());`
	for _, zone := range []string{"UTC", "America/New_York", "Australia/Lord_Howe"} {
		t.Run(zone, func(t *testing.T) {
			t.Setenv("TZ", zone)
			cmd := exec.Command(node, "-e", source)
			want, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("Node: %v\n%s", err, want)
			}
			got, err := runGoProgram(t, emitGoProgram(t, source), true)
			if err != nil {
				t.Fatalf("%v\n%s", err, got)
			}
			if got != string(want) {
				t.Fatalf("got:\n%s\nwant:\n%s", got, want)
			}
		})
	}
}

func TestGoNativeDateLocaleDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node not installed")
	}
	t.Setenv("TZ", "UTC")
	source := `const date=new Date("2024-02-29T03:04:05.123Z");for(const locale of ["en-US","en-GB"]){console.log(date.toLocaleString(locale),date.toLocaleDateString(locale),date.toLocaleTimeString(locale));console.log(date.toLocaleString(locale,{year:"numeric",month:"long",day:"numeric",hour:"numeric",minute:"numeric",timeZone:"UTC"}));console.log(date.toLocaleString(locale,{weekday:"long",hour:"numeric"}));console.log(date.toLocaleString(locale,{year:"2-digit",month:"2-digit",day:"2-digit",hour:"2-digit",minute:"2-digit",second:"2-digit",fractionalSecondDigits:3,hour12:false}));console.log(date.toLocaleDateString(locale,{era:"short"}),date.toLocaleTimeString(locale,{timeZoneName:"short"}));console.log(date.toLocaleString(locale,{dateStyle:"full",timeStyle:"full",timeZone:"UTC"}));console.log(date.toLocaleTimeString(locale,{hour:"numeric",minute:"numeric",timeZone:"America/New_York",timeZoneName:"long"}));}
 console.log(Date.prototype.toGMTString===Date.prototype.toUTCString,Object.keys(Date.prototype).length,Date.prototype.getTime.length,Date.prototype.setHours.length,Date.prototype.toJSON.name,Date.now.length,Date.parse.length,Date.UTC.length,Date.UTC.name);
 try{date.toLocaleDateString("en-US",{timeStyle:"short"});}catch(e){console.log(e.name);}
 try{date.toLocaleString("en-US",{month:"banana"});}catch(e){console.log(e.name);}`
	cmd := exec.Command(node, "-e", source)
	want, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestGoNativeJSONDifferential(t *testing.T) {
	node, err := exec.LookPath("node")
	if err != nil {
		t.Skip("Node not installed")
	}
	source := `for(const text of ["null","true","false","-0","1e400","01","[1,]","{\"x\":1,}","\"\\ud800\"","\"\\u0000\"","{\"b\":1,\"1\":2,\"b\":3}"]){try{const value=JSON.parse(text);console.log(JSON.stringify(value),Object.is(value,-0));}catch(e){console.log(e.name);}}
 console.log(JSON.parse("12345678901234567890",(key,value,context)=>context.source));
 console.log(JSON.stringify(JSON.stringify({a:1,b:[2,3]},undefined,"\ud800")));
 console.log(JSON.parse===JSON.parse,JSON.parse.length,JSON.stringify.length,JSON.stringify.name,JSON.stringify(undefined)===undefined);
 const namespace=JSON;const parse=namespace.parse;console.log(parse("5"));namespace.parse=(value)=>7;console.log(JSON.parse("5"));`
	cmd := exec.Command(node, "-e", source)
	want, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Node: %v\n%s", err, want)
	}
	got, err := runGoProgram(t, emitGoProgram(t, source), true)
	if err != nil {
		t.Fatalf("%v\n%s", err, got)
	}
	if got != string(want) {
		t.Fatalf("got:\n%s\nwant:\n%s", got, want)
	}
}

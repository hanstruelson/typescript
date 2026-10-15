package goemit

// PathRuntime implements Node's path.posix and path.win32 APIs for the
// generated runtime. The string boundary preserves lone UTF-16 surrogates,
// which Go strings cannot represent as ordinary Unicode.
//
// Copyright Joyent, Inc. and other Node contributors.
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
const PathRuntime = `
var tsNodePathModules struct { sync.Once; posix, win32 *tsObject }

func tsPathError(message string) { panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:message})}) }
func tsPathText(value *tsString) string {
	var out strings.Builder
	for i:=0;i<len(value.units);i++ {
		unit:=value.units[i]
		if unit>=0xd800&&unit<=0xdbff&&i+1<len(value.units)&&value.units[i+1]>=0xdc00&&value.units[i+1]<=0xdfff {
			out.WriteRune(utf16.DecodeRune(rune(unit),rune(value.units[i+1])));i++;continue
		}
		if unit>=0xd800&&unit<=0xdfff {out.Write([]byte{0xed,0x80|byte(unit>>6&0x3f),0x80|byte(unit&0x3f)})} else {out.WriteRune(rune(unit))}
	}
	return out.String()
}
func tsPathInputNamed(value tsValue, argument string) string {
	if value.kind!=tsStringKind { tsPathError("The \""+argument+"\" argument must be of type string.") }
	return tsPathText((*tsString)(value.ref))
}
func tsPathInput(value tsValue) string { return tsPathInputNamed(value,"path") }
func tsPathCoerceStringObject(loop *tsLoop,value tsValue) tsValue {
	if value.kind==tsObjectKind&&(*tsObject)(value.ref).boxedPrimitive.kind==tsStringKind {return tsStringValueWithLoop(loop,value)}
	return value
}
func tsPathIsStringObject(value tsValue) bool {return value.kind==tsObjectKind&&(*tsObject)(value.ref).boxedPrimitive.kind==tsStringKind}
func tsPathCoerceArgs(loop *tsLoop,args []tsValue) []tsValue {
	for _,arg:=range args {if tsPathIsStringObject(arg) {converted:=make([]tsValue,0,len(args));for _,item:=range args{converted=append(converted,tsPathCoerceStringObject(loop,item))};return converted}}
	return args
}
func tsPathOut(value string) tsValue {
	units:=[]uint16{}
	for i:=0;i<len(value); {
		if i+2<len(value)&&value[i]==0xed&&value[i+1]>=0xa0&&value[i+1]<=0xbf&&value[i+2]>=0x80&&value[i+2]<=0xbf {
			unit:=uint16(value[i]&0x0f)<<12|uint16(value[i+1]&0x3f)<<6|uint16(value[i+2]&0x3f)
			if unit>=0xd800&&unit<=0xdfff {units=append(units,unit);i+=3;continue}
		}
		r,size:=utf8.DecodeRuneInString(value[i:])
		if size==0 {break}
		if r==utf8.RuneError&&size==1 {units=append(units,0xfffd);i++;continue}
		if r<=0xffff {units=append(units,uint16(r))} else {high,low:=utf16.EncodeRune(r);units=append(units,uint16(high),uint16(low))}
		i+=size
	}
	return tsStringReference(tsStringUnits(units))
}

func tsPathPosixResolveTexts(paths []string) string {
	resolved:="";absolute:=false
	for i:=len(paths)-1;i>=-1&&!absolute;i-- {
		value:=""
		if i>=0 {value=paths[i]} else {var err error;value,err=os.Getwd();if err!=nil{tsPathError("Unable to read current working directory")};if runtime.GOOS=="windows"{value=strings.ReplaceAll(value,"\\","/");if slash:=strings.IndexByte(value,'/');slash>=0{value=value[slash:]}else if len(value)>0{value=value[len(value)-1:]}}}
		if value=="" {continue}
		resolved=value+"/"+resolved
		absolute=value[0]=='/'
	}
	resolved=path.Clean(resolved)
	if absolute {if !strings.HasPrefix(resolved,"/"){resolved="/"+resolved};return resolved}
	if resolved=="" {return "."};return resolved
}
func tsPathPosixResolve(args []tsValue) string {
	paths:=make([]string,len(args))
	for i:=len(args)-1;i>=0;i-- {
		value:=tsPathInputNamed(args[i],"paths["+strconv.Itoa(i)+"]")
		paths[i]=value
		if strings.HasPrefix(value,"/") { break }
	}
	return tsPathPosixResolveTexts(paths)
}
func tsPathPosixNormalize(value string) string {
	if value=="" {return "."};absolute:=strings.HasPrefix(value,"/");trailing:=strings.HasSuffix(value,"/")
	value=path.Clean(value)
	if trailing&&value!="/" {value+="/"};if absolute&&!strings.HasPrefix(value,"/"){value="/"+value};return value
}
func tsPathPosixJoin(args []tsValue) string {
	parts:=make([]string,0,len(args));for i,arg:=range args{value:=tsPathInputNamed(arg,"paths["+strconv.Itoa(i)+"]");if value!=""{parts=append(parts,value)}}
	if len(parts)==0{return "."};return tsPathPosixNormalize(strings.Join(parts,"/"))
}
func tsPathPosixRelative(from,to string) string {
	from=tsPathPosixResolveTexts([]string{from});to=tsPathPosixResolveTexts([]string{to});if from==to{return ""}
	a:=strings.Split(strings.Trim(from,"/"),"/");b:=strings.Split(strings.Trim(to,"/"),"/")
	if len(a)==1&&a[0]==""{a=nil};if len(b)==1&&b[0]==""{b=nil}
	common:=0;for common<len(a)&&common<len(b)&&a[common]==b[common]{common++}
	parts:=make([]string,0,len(a)-common+len(b)-common);for i:=common;i<len(a);i++{parts=append(parts,"..")} ;parts=append(parts,b[common:]...);return strings.Join(parts,"/")
}
func tsPathPosixDirname(value string) string {
	if value=="" {return "."};rooted:=value[0]=='/';matched:=true
	for i:=len(value)-1;i>=1;i-- {if value[i]=='/' {if !matched {if rooted&&i==1{return "//"};return value[:i]}} else {matched=false}}
	if rooted{return "/"};return "."
}
func tsPathPosixBasename(value,ext string,hasExt bool) string {
	if value=="" {return ""};end:=len(value);for end>0&&value[end-1]=='/'{end--};if end==0{return ""}
	start:=strings.LastIndexByte(value[:end],'/')+1;base:=value[start:end]
	if hasExt&&ext!=""&&len(ext)<=len(base)&&strings.HasSuffix(base,ext){base=strings.TrimSuffix(base,ext)};return base
}
func tsPathExtnameBase(base string) string {
	if base==""||base=="."||base==".." {return ""};dot:=strings.LastIndexByte(base,'.');if dot<0{return ""};if dot==0&&strings.IndexByte(base[1:],'.')<0{return ""};return base[dot:]
}
func tsPathPosixParse(value string) (root,dir,base,ext,name string) {
	if value=="" {return};rootEnd:=0;if value[0]=='/'{rootEnd=1;root="/"}
	end:=len(value);for end>rootEnd&&value[end-1]=='/'{end--};if end<=rootEnd {dir=root;return}
	start:=rootEnd;separator:=strings.LastIndexByte(value[rootEnd:end],'/');if separator>=0 {separator+=rootEnd;start=separator+1;if separator>=rootEnd {if separator==0 {dir=root}else{dir=value[:separator]}}} else {dir=root}
	base=value[start:end];ext=tsPathExtnameBase(base);name=base;if ext!=""{name=strings.TrimSuffix(base,ext)};return
}

func tsPathWinSep(c byte) bool {return c=='/'||c=='\\'}
func tsPathWinDeviceRoot(c byte) bool {return c>='A'&&c<='Z'||c>='a'&&c<='z'}
// rootEnd includes the separator after a UNC share when one is present.
func tsPathWinRoot(value string) (rootEnd int,device string,absolute,unc bool) {
	n:=len(value);if n==0{return}
	if n>=2&&tsPathWinSep(value[0])&&tsPathWinSep(value[1]) {
		serverStart:=2;serverEnd:=serverStart;for serverEnd<n&&!tsPathWinSep(value[serverEnd]){serverEnd++}
		if serverEnd>serverStart&&serverEnd<n {
			shareStart:=serverEnd;for shareStart<n&&tsPathWinSep(value[shareStart]){shareStart++};shareEnd:=shareStart
			for shareEnd<n&&!tsPathWinSep(value[shareEnd]){shareEnd++}
			if shareEnd>shareStart {device="\\\\"+value[serverStart:serverEnd]+"\\"+value[shareStart:shareEnd];rootEnd=shareEnd;if rootEnd<n {rootEnd++};return rootEnd,device,true,true}
		}
	}
	if n>=2&&tsPathWinDeviceRoot(value[0])&&value[1]==':' {device=value[:2];rootEnd=2;if n>2&&tsPathWinSep(value[2]){rootEnd=3;absolute=true};return}
	if tsPathWinSep(value[0]) {rootEnd=1;absolute=true}
	return
}
func tsPathWinSegments(value string,allowAbove bool) []string {
	parts:=strings.FieldsFunc(value,func(c rune)bool{return c=='/'||c=='\\'});out:=make([]string,0,len(parts))
	for _,part:=range parts {switch part{case "",".":continue;case "..":if len(out)>0&&out[len(out)-1]!=".."{out=out[:len(out)-1]}else if allowAbove{out=append(out,"..")};default:out=append(out,part)}};return out
}
func tsPathWinNormalize(value string) string {
	if value=="" {return "."};value=strings.ReplaceAll(value,"/","\\");rootEnd,device,absolute,unc:=tsPathWinRoot(value);trailing:=tsPathWinSep(value[len(value)-1])
	var root string;bodyStart:=rootEnd
	if unc {root=device+"\\";bodyStart=rootEnd;for bodyStart<len(value)&&tsPathWinSep(value[bodyStart]){bodyStart++}} else if device!="" {if absolute{root=device+"\\"}else{root=device}} else if absolute {root="\\"}
	segments:=tsPathWinSegments(value[bodyStart:],!absolute);tail:=strings.Join(segments,"\\")
	if tail=="" {if device!=""&&!absolute&&!unc{return device+"."};if root!="" {return root};if device!="" {return device+"."};return "."}
	out:=root+tail;if !absolute&&device!="" {out=device+tail};if trailing&&!strings.HasSuffix(out,"\\"){out+="\\"};return out
}
func tsPathWinJoin(args []tsValue) string {
	parts:=make([]string,0,len(args));for i,arg:=range args{value:=tsPathInputNamed(arg,"paths["+strconv.Itoa(i)+"]");if value!=""{parts=append(parts,strings.ReplaceAll(value,"/","\\"))}}
	if len(parts)==0{return "."};joined:=strings.Join(parts,"\\");firstPartLen:=len(parts[0])
	// Match Node's win32.join UNC guard: preserve an intentional leading pair
	// of separators in the first nonempty argument, collapse accidental runs.
	needsReplace:=true;slashCount:=0
	if len(joined)>0&&tsPathWinSep(joined[0]) {slashCount=1;if firstPartLen>1&&len(joined)>1&&tsPathWinSep(joined[1]) {slashCount=2;if firstPartLen>2&&len(joined)>2 {if tsPathWinSep(joined[2]) {slashCount++}else{needsReplace=false}}}}
	if needsReplace {for slashCount<len(joined)&&tsPathWinSep(joined[slashCount]){slashCount++};if slashCount>=2{joined="\\"+joined[slashCount:]}}
	return tsPathWinNormalize(joined)
}
func tsPathWinResolveTexts(paths []string) string {
	device:="";absolute:=false;segments:=[]string{}
	for i:=len(paths)-1;i>=-1;i-- {
		value:=""
		if i>=0 {value=strings.ReplaceAll(paths[i],"/","\\")} else {cwd,err:=os.Getwd();if err!=nil{tsPathError("Unable to read current working directory")};value=strings.ReplaceAll(cwd,"/","\\")}
		if value=="" {continue}
		rootEnd,currentDevice,isAbsolute,unc:=tsPathWinRoot(value)
		if currentDevice!="" {if device!=""&&!strings.EqualFold(device,currentDevice){continue};if device==""{device=currentDevice}}
		if absolute {if device!="" {break};continue}
		body:=value[rootEnd:];newParts:=tsPathWinSegments(body,true);segments=append(newParts,segments...)
		if isAbsolute {absolute=true;if device!=""||unc{break}}
	}
	if !absolute { // A drive-relative path uses that drive's cwd when available.
		absolute=true
	}
	clean:=tsPathWinSegments(strings.Join(segments,"\\"),false);tail:=strings.Join(clean,"\\")
	if absolute {if device!="" {if tail=="" {return device+"\\"};return device+"\\"+tail};if tail==""{return "\\"};return "\\"+tail}
	if device!="" {if tail==""{return device+"."};return device+tail};if tail==""{return "."};return tail
}
func tsPathWinResolve(args []tsValue) string {paths:=make([]string,len(args));for i,arg:=range args{paths[i]=tsPathInputNamed(arg,"paths["+strconv.Itoa(i)+"]")};return tsPathWinResolveTexts(paths)}
func tsPathWinRelative(from,to string) string {
	from=tsPathWinResolveTexts([]string{from});to=tsPathWinResolveTexts([]string{to});if strings.EqualFold(from,to){return ""}
	_,fromDevice,_,_:=tsPathWinRoot(from);_,toDevice,_,_:=tsPathWinRoot(to);if !strings.EqualFold(fromDevice,toDevice){return to}
	strip:=func(value,device string)string{if device!=""&&len(value)>=len(device){value=value[len(device):]};return strings.Trim(value,"\\")}
	a:=strings.Split(strip(from,fromDevice),"\\");b:=strings.Split(strip(to,toDevice),"\\");if len(a)==1&&a[0]==""{a=nil};if len(b)==1&&b[0]==""{b=nil}
	common:=0;for common<len(a)&&common<len(b)&&strings.EqualFold(a[common],b[common]){common++};parts:=make([]string,0,len(a)-common+len(b)-common)
	for i:=common;i<len(a);i++{parts=append(parts,"..")};parts=append(parts,b[common:]...);return strings.Join(parts,"\\")
}
func tsPathWinBasename(value,ext string,hasExt bool) string {
	value=strings.ReplaceAll(value,"/","\\");start:=0;if len(value)>=2&&tsPathWinDeviceRoot(value[0])&&value[1]==':'{start=2};end:=len(value);for end>start&&tsPathWinSep(value[end-1]){end--};if end<=start{return ""}
	for i:=end-1;i>=start;i--{if tsPathWinSep(value[i]){start=i+1;break}};base:=value[start:end]
	if hasExt&&ext!=""&&len(ext)<=len(base)&&strings.HasSuffix(base,ext){base=strings.TrimSuffix(base,ext)};return base
}
func tsPathWinParse(value string) (root,dir,base,ext,name string) {
	if value=="" {return};rootEnd,_,_,_:=tsPathWinRoot(value);if rootEnd>0{root=value[:rootEnd]}
	end:=len(value);for end>rootEnd&&tsPathWinSep(value[end-1]){end--};if end<=rootEnd {dir=root;return}
	start:=rootEnd;separator:=-1;for i:=rootEnd;i<end;i++{if tsPathWinSep(value[i]){separator=i}}
	if separator>=0 {start=separator+1;if separator>=rootEnd {dir=value[:separator]}else{dir=root}} else {dir=root}
	base=value[start:end];ext=tsPathExtnameBase(base);name=base;if ext!=""{name=strings.TrimSuffix(base,ext)};return
}
func tsPathWinDirname(value string) string {_,dir,_,_,_:=tsPathWinParse(value);if dir==""{return "."};return dir}
func tsPathWinIsAbsolute(value string) bool {_,_,absolute,_:=tsPathWinRoot(value);return absolute}
func tsPathWinToNamespaced(value string) string {
	resolved:=tsPathWinResolveTexts([]string{value});if len(resolved)<3{return resolved}
	if strings.HasPrefix(resolved,"\\\\?\\")||strings.HasPrefix(resolved,"\\\\.\\"){return resolved}
	if strings.HasPrefix(resolved,"\\\\") {return "\\\\?\\UNC\\"+resolved[2:]}
	if tsPathWinDeviceRoot(resolved[0])&&resolved[1]==':'&&resolved[2]=='\\' {return "\\\\?\\"+resolved};return resolved
}

func tsPathFormat(loop *tsLoop,value tsValue,windows bool) tsValue {
	if tsIsPrimitiveValue(value)||value.kind==tsArrayKind||value.kind==tsFunctionKind {tsPathError("The \"pathObject\" argument must be of type object.")}
	sep:="/";if windows{sep="\\"}
	field:=func(name string)string{v:=tsGet(value,name);if !tsTruthy(v){return ""};return tsPathText(tsStringValueWithLoop(loop,v))}
	rootText:=field("root");dirText:=field("dir");baseText:=field("base");nameText:=field("name");extText:=field("ext")
	if baseText=="" {if extText!=""&&extText[0]!='.'{extText="."+extText};baseText=nameText+extText}
	if dirText==""{dirText=rootText};if dirText==""{return tsPathOut(baseText)};if rootText!=""&&dirText==rootText{return tsPathOut(dirText+baseText)};return tsPathOut(dirText+sep+baseText)
}
func tsPathInvoke(loop *tsLoop,windows bool,name string,args []tsValue) tsValue {
	// Path.cpp stringifies boxed String arguments before its Rust validators,
	// except for format() and toNamespacedPath().
	if name!="format"&&name!="toNamespacedPath"&&name!="_makeLong" {args=tsPathCoerceArgs(loop,args)}
	input:=func(i int)string{return tsPathInput(tsArg(args,i))}
	switch name {
	case "resolve":if windows{return tsPathOut(tsPathWinResolve(args))};return tsPathOut(tsPathPosixResolve(args))
	case "normalize":if windows{return tsPathOut(tsPathWinNormalize(input(0)))};return tsPathOut(tsPathPosixNormalize(input(0)))
	case "isAbsolute":if windows{return tsBooleanValue(tsPathWinIsAbsolute(input(0)))};return tsBooleanValue(strings.HasPrefix(input(0),"/"))
	case "join":if windows{return tsPathOut(tsPathWinJoin(args))};return tsPathOut(tsPathPosixJoin(args))
	case "relative":from:=tsPathInputNamed(tsArg(args,0),"from");to:=tsPathInputNamed(tsArg(args,1),"to");if windows{return tsPathOut(tsPathWinRelative(from,to))};return tsPathOut(tsPathPosixRelative(from,to))
	case "toNamespacedPath","_makeLong":
		value:=tsArg(args,0);if value.kind!=tsStringKind{return value};if windows{return tsPathOut(tsPathWinToNamespaced(tsPathText((*tsString)(value.ref))))};return value
	case "dirname":if windows{return tsPathOut(tsPathWinDirname(input(0)))};return tsPathOut(tsPathPosixDirname(input(0)))
	case "basename":
		ext:=tsArg(args,1);hasExt:=len(args)>1&&!tsIsUndefined(ext);if hasExt&&ext.kind!=tsStringKind{tsPathError("The \"ext\" argument must be of type string.")};extension:="";if hasExt{extension=tsPathText((*tsString)(ext.ref))}
		value:=input(0);if hasExt&&value==extension{return tsPathOut("")};if windows{return tsPathOut(tsPathWinBasename(value,extension,hasExt))};return tsPathOut(tsPathPosixBasename(value,extension,hasExt))
	case "extname":value:=input(0);base:=tsPathPosixBasename(value,"",false);if windows{base=tsPathWinBasename(value,"",false)};return tsPathOut(tsPathExtnameBase(base))
	case "format":return tsPathFormat(loop,tsArg(args,0),windows)
	case "parse":
		value:=input(0);var root,dir,base,ext,name string;if windows{root,dir,base,ext,name=tsPathWinParse(value)}else{root,dir,base,ext,name=tsPathPosixParse(value)}
		out:=tsNewObject();out.set("root",tsPathOut(root));out.set("dir",tsPathOut(dir));out.set("base",tsPathOut(base));out.set("ext",tsPathOut(ext));out.set("name",tsPathOut(name));return tsObjectValue(out)
	case "matchesGlob":
		value:=tsPathInputNamed(tsArg(args,0),"path");pattern:=tsPathInputNamed(tsArg(args,1),"pattern");if windows{value=strings.ReplaceAll(value,"\\","/");pattern=strings.ReplaceAll(pattern,"\\","/")};return tsBooleanValue(tsFSGlobMatch(pattern,value))
	}
	panic("Unsupported path operation")
}
func tsPathModuleFor(windows bool) tsValue {
	tsNodePathModules.Do(func(){
		makeModule:=func(isWindows bool)*tsObject {
			out:=tsNewObject();for _,entry:=range []struct{name string;length int}{{"resolve",0},{"normalize",1},{"isAbsolute",1},{"join",0},{"relative",2},{"toNamespacedPath",1},{"dirname",1},{"basename",1},{"extname",1},{"format",1},{"parse",1},{"matchesGlob",2}} {
				operation:=entry.name;out.set(operation,tsNamedNativeMethod(operation,entry.length,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsPathInvoke(loop,isWindows,operation,args)}))
			}
			sep,delimiter:="/",":";if isWindows{sep,delimiter="\\",";"};out.set("sep",tsStringReference(tsStringUTF8(sep)));out.set("delimiter",tsStringReference(tsStringUTF8(delimiter)));return out
		}
		posix:=makeModule(false);win32:=makeModule(true);posix.set("win32",tsObjectValue(win32));posix.set("posix",tsObjectValue(posix));win32.set("win32",tsObjectValue(win32));win32.set("posix",tsObjectValue(posix));posix.set("_makeLong",posix.values["toNamespacedPath"]);win32.set("_makeLong",win32.values["toNamespacedPath"]);tsNodePathModules.posix=posix;tsNodePathModules.win32=win32
	})
	if windows {return tsObjectValue(tsNodePathModules.win32)};return tsObjectValue(tsNodePathModules.posix)
}
func tsPathModule() tsValue {return tsPathModuleFor(runtime.GOOS=="windows")}
`

package goemit

// internalAssertMyersRuntime ports Bun's shared Myers diff algorithm and the
// JS-facing list/render helpers used by Node's AssertionError implementation.
const internalAssertMyersRuntime = `
type tsAssertDiffKind uint8
const (
	tsAssertDiffInsert tsAssertDiffKind = iota
	tsAssertDiffDelete
	tsAssertDiffEqual
)
type tsAssertDiffEntry struct { kind tsAssertDiffKind; value []uint16 }

func tsAssertDiffLines(value []uint16) [][]uint16 {
	lines:=make([][]uint16,0,len(value)>>4)
	start:=0
	for index,unit:=range value { if unit=='\n' { lines=append(lines,value[start:index]);start=index+1 } }
	return append(lines,value[start:])
}
func tsAssertDiffLineEqual(a,b []uint16,checkComma bool)bool {
	if len(a)==len(b){for index:=range a{if a[index]!=b[index]{return false}};return true}
	if !checkComma { return false }
	larger,smaller:=a,b;if len(larger)<len(smaller){larger,smaller=smaller,larger}
	if len(larger)-len(smaller)!=1||larger[len(larger)-1]!=',' { return false }
	for index:=range smaller{if larger[index]!=smaller[index]{return false}}
	return true
}
func tsAssertDiffInputsTooLarge()tsValue{return tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Input strings are too large to diff. Please open a bug report on GitHub.",fields:map[string]tsValue{"code":tsStringReference(tsStringUTF8("ERR_INVALID_ARG_VALUE"))}})}
func tsAssertDiffTooLarge()tsValue{return tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Diffing these two values would create a string that is too large. If this was intentional, please open a bug report on GitHub.",fields:map[string]tsValue{"code":tsStringReference(tsStringUTF8("ERR_INVALID_ARG_VALUE"))}})}
func tsAssertMyers(actual,expected []uint16,lines,checkComma bool)([]tsAssertDiffEntry,tsValue){
	if len(actual)==0&&len(expected)==0{return []tsAssertDiffEntry{},tsU}
	var a,b [][]uint16
	if lines { a,b=tsAssertDiffLines(actual),tsAssertDiffLines(expected) } else { a=make([][]uint16,len(actual));b=make([][]uint16,len(expected));for index,value:=range actual{a[index]=[]uint16{value}};for index,value:=range expected{b[index]=[]uint16{value}};checkComma=false }
	maxInt:=int(^uint(0)>>1)
	if len(a)>maxInt-len(b){return nil,tsAssertDiffInputsTooLarge()}
	max:=len(a)+len(b);if uint64(max)>uint64(^uint32(0)){return nil,tsAssertDiffInputsTooLarge()}
	if max>(maxInt-1)/2{return nil,tsAssertDiffTooLarge()}
	graphSize:=2*max+1;if uint64(graphSize)>uint64(^uint32(0)){return nil,tsAssertDiffTooLarge()}
	const maxTraceBytes=64*1024*1024
	graph:=make([]uint32,graphSize);trace:=make([][]uint32,0,max+1)
	for distance:=0;distance<=max;distance++{
		if graphSize>maxTraceBytes/4/(distance+1){return nil,tsAssertDiffTooLarge()}
		trace=append(trace,append([]uint32(nil),graph...))
		diagonalStart,diagonalEnd:=-distance,distance
		for diagonal:=diagonalStart;diagonal<=diagonalEnd;diagonal+=2{
			index:=diagonal+max
			var x uint32
			if diagonal==diagonalStart||(diagonal!=diagonalEnd&&graph[index-1]<graph[index+1]){x=graph[index+1]}else{x=graph[index-1]+1}
			y:=int(x)-diagonal
			for int(x)<len(a)&&y<len(b)&&tsAssertDiffLineEqual(a[int(x)],b[y],checkComma){x++;y++}
			graph[index]=x
			if int(x)>=len(a)&&y>=len(b){return tsAssertMyersBacktrack(trace,a,b,checkComma),tsU}
		}
	}
	return nil,tsAssertDiffTooLarge()
}
func tsAssertMyersBacktrack(trace [][]uint32,actual,expected [][]uint16,checkComma bool)[]tsAssertDiffEntry{
	max:=len(actual)+len(expected);x,y:=len(actual),len(expected);result:=make([]tsAssertDiffEntry,0)
	for distance:=len(trace);distance>0;{
		distance--;graph:=trace[distance];diagonal:=x-y;offset:=diagonal+max
		previousDiagonal:=diagonal-1
		if diagonal==-distance||(diagonal!=distance&&graph[offset-1]<graph[offset+1]){previousDiagonal=diagonal+1}
		previousX:=int(graph[previousDiagonal+max]);previousY:=previousX-previousDiagonal
		for x>previousX&&y>previousY{
			value:=actual[x-1]
			if checkComma&&len(value)>0&&value[len(value)-1]!=','{value=expected[y-1]}
			result=append(result,tsAssertDiffEntry{kind:tsAssertDiffEqual,value:value});x--;y--
		}
		if distance>0{if x>previousX{result=append(result,tsAssertDiffEntry{kind:tsAssertDiffInsert,value:actual[x-1]});x--}else{result=append(result,tsAssertDiffEntry{kind:tsAssertDiffDelete,value:expected[y-1]});y--}}
	}
	return result
}
func tsAssertDiffString(value []uint16)tsValue{return tsStringReference(tsStringUnits(append([]uint16(nil),value...)))}
func tsAssertDiffText(value tsValue,name string)[]uint16{if value.kind!=tsStringKind{tsFSInvalid(name,"a string")};return append([]uint16(nil),(*tsString)(value.ref).units...)}
func tsAssertDiffColors(loop *tsLoop,value tsValue)map[string][]uint16{
	colors:=map[string][]uint16{"green":nil,"red":nil,"white":nil,"blue":nil}
	if value.kind!=tsObjectKind&&value.kind!=tsInstanceKind&&value.kind!=tsArrayKind{return colors}
	for _,name:=range []string{"green","red","white","blue"}{item:=tsGet(loop,value,tsStringReference(tsStringUTF8(name)));if item.kind==tsStringKind{colors[name]=append([]uint16(nil),(*tsString)(item.ref).units...)}}
	return colors
}
func tsAssertDiffAppend(out *[]uint16,value []uint16){*out=append(*out,value...)}
func tsAssertDiffAppendASCII(out *[]uint16,value string){for index:=0;index<len(value);index++{*out=append(*out,uint16(value[index]))}}
func tsAssertDiffAppendColor(out *[]uint16,colors map[string][]uint16,name string){tsAssertDiffAppend(out,colors[name])}
func tsAssertDiffValue(entry tsAssertDiffEntry,lines bool)tsValue{if lines{return tsAssertDiffString(entry.value)};if len(entry.value)==0{return tsNumberValue(0)};return tsNumberValue(float64(entry.value[0]))}
func tsAssertDiffListValue(entries []tsAssertDiffEntry,lines bool)tsValue{
	values:=make([]tsValue,0,len(entries));for _,entry:=range entries{object:=tsNewObject();object.prototype=tsNull;object.set("kind",tsNumberValue(float64(entry.kind)));object.set("value",tsAssertDiffValue(entry,lines));values=append(values,tsObjectValue(object))};return tsArrayValue(&tsArray{values:values})
}
func tsAssertDiffSimple(entries []tsAssertDiffEntry,colors map[string][]uint16)tsValue{
	out:=[]uint16{'\n'};for index:=len(entries)-1;index>=0;index--{entry:=entries[index];color:="white";if entry.kind==tsAssertDiffInsert{color="green"}else if entry.kind==tsAssertDiffDelete{color="red"};tsAssertDiffAppendColor(&out,colors,color);tsAssertDiffAppend(&out,entry.value);tsAssertDiffAppendColor(&out,colors,"white")};return tsAssertDiffString(out)
}
func tsAssertDiffEqualLine(out *[]uint16,value []uint16,colors map[string][]uint16){tsAssertDiffAppendColor(out,colors,"white");tsAssertDiffAppendASCII(out,"  ");tsAssertDiffAppend(out,value);tsAssertDiffAppendASCII(out,"\n")}
func tsAssertDiffJSWhitespace(value uint16)bool{return value>=0x09&&value<=0x0d||value==0x20||value==0x00a0||value==0x1680||value>=0x2000&&value<=0x200a||value==0x2028||value==0x2029||value==0x202f||value==0x205f||value==0x3000||value==0xfeff}
func tsAssertDiffLinesText(entries []tsAssertDiffEntry,colors map[string][]uint16)(tsValue,bool){
	const collapseAfter=5
	out:=[]uint16{'\n'};skipped:=false;unchanged:=0
	for index:=len(entries)-1;index>=0;index--{
		entry:=entries[index];previousEqual:=index+1<len(entries)&&entries[index+1].kind==tsAssertDiffEqual
		if previousEqual&&entry.kind!=tsAssertDiffEqual{
			if unchanged==collapseAfter+1{tsAssertDiffEqualLine(&out,entries[index+1].value,colors)}else if unchanged==collapseAfter+2{tsAssertDiffEqualLine(&out,entries[index+2].value,colors);tsAssertDiffEqualLine(&out,entries[index+1].value,colors)}
			if unchanged>=collapseAfter+3{tsAssertDiffAppendColor(&out,colors,"blue");tsAssertDiffAppendASCII(&out,"...");tsAssertDiffAppendColor(&out,colors,"white");tsAssertDiffAppendASCII(&out,"\n");tsAssertDiffEqualLine(&out,entries[index+1].value,colors);skipped=true};unchanged=0
		}
		switch entry.kind{
		case tsAssertDiffInsert:tsAssertDiffAppendColor(&out,colors,"green");tsAssertDiffAppendASCII(&out,"+");tsAssertDiffAppendColor(&out,colors,"white");tsAssertDiffAppendASCII(&out," ");tsAssertDiffAppend(&out,entry.value);tsAssertDiffAppendASCII(&out,"\n")
		case tsAssertDiffDelete:tsAssertDiffAppendColor(&out,colors,"red");tsAssertDiffAppendASCII(&out,"-");tsAssertDiffAppendColor(&out,colors,"white");tsAssertDiffAppendASCII(&out," ");tsAssertDiffAppend(&out,entry.value);tsAssertDiffAppendASCII(&out,"\n")
		case tsAssertDiffEqual:if unchanged<collapseAfter{tsAssertDiffEqualLine(&out,entry.value,colors)};unchanged++
		}
	}
	for len(out)>1&&tsAssertDiffJSWhitespace(out[len(out)-1]){out=out[:len(out)-1]};return tsAssertDiffString(out),skipped
}
var tsInternalAssertMyersState struct{sync.Once;module tsValue}
func tsInternalAssertMyersModule(loop *tsLoop)tsValue{
	tsInternalAssertMyersState.Do(func(){module:=tsNewObject();myers:=tsNamedNativeMethod("myersDiff",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if len(args)<2{tsPropertyFailure("myersDiff requires actual and expected strings")};actual,expected:=tsAssertDiffText(tsArg(args,0),"actual"),tsAssertDiffText(tsArg(args,1),"expected");entries,failure:=tsAssertMyers(actual,expected,tsTruthy(tsArg(args,3)),tsTruthy(tsArg(args,2)));if !tsIsUndefined(failure){panic(tsThrown{failure})};return tsAssertDiffListValue(entries,tsTruthy(tsArg(args,3)))})
		simple:=tsNamedNativeMethod("printSimpleMyersDiff",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if len(args)<2{tsPropertyFailure("printSimpleMyersDiff requires actual and expected strings")};entries,failure:=tsAssertMyers(tsAssertDiffText(tsArg(args,0),"actual"),tsAssertDiffText(tsArg(args,1),"expected"),false,false);if !tsIsUndefined(failure){panic(tsThrown{failure})};return tsAssertDiffSimple(entries,tsAssertDiffColors(loop,tsArg(args,2)))})
		lines:=tsNamedNativeMethod("printMyersDiff",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{if len(args)<2{tsPropertyFailure("printMyersDiff requires actual and expected strings")};entries,failure:=tsAssertMyers(tsAssertDiffText(tsArg(args,0),"actual"),tsAssertDiffText(tsArg(args,1),"expected"),true,tsTruthy(tsArg(args,2)));if !tsIsUndefined(failure){panic(tsThrown{failure})};message,skipped:=tsAssertDiffLinesText(entries,tsAssertDiffColors(loop,tsArg(args,3)));result:=tsNewObject();result.set("message",message);result.set("skipped",tsBooleanValue(skipped));return tsObjectValue(result)})
		module.set("myersDiff",myers);module.set("printSimpleMyersDiff",simple);module.set("printMyersDiff",lines);module.set("default",tsObjectValue(module));tsInternalAssertMyersState.module=tsObjectValue(module)
	});return tsInternalAssertMyersState.module
}
`

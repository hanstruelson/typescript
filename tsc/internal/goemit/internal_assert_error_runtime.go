package goemit

// internalAssertErrorRuntime ports Node's AssertionError message, diff, and
// metadata construction over the local inspect and Myers helpers.
const internalAssertErrorRuntime = `
var tsInternalAssertionErrorState struct{sync.Once;constructor tsValue}

var tsAssertReadableOperators=map[string]string{
	"deepStrictEqual":"Expected values to be strictly deep-equal:",
	"strictEqual":"Expected values to be strictly equal:",
	"strictEqualObject":"Expected \"actual\" to be reference-equal to \"expected\":",
	"deepEqual":"Expected values to be loosely deep-equal:",
	"notDeepStrictEqual":"Expected \"actual\" not to be strictly deep-equal to:",
	"notStrictEqual":"Expected \"actual\" to be strictly unequal to:",
	"notStrictEqualObject":"Expected \"actual\" not to be reference-equal to \"expected\":",
	"notDeepEqual":"Expected \"actual\" not to be loosely deep-equal to:",
	"notIdentical":"Values have same structure but are not reference-equal:",
	"notDeepEqualUnequal":"Expected values not to be loosely deep-equal:",
	"partialDeepStrictEqual":"Expected values to be partially and strictly deep-equal:",
}
func tsAssertErrorString(loop *tsLoop,value tsValue)string{return tsStringValueWithLoop(loop,value).String()}
func tsAssertErrorObjectLike(value tsValue)bool{return tsNodeUtilIsObject(value)}
func tsAssertErrorCopy(loop *tsLoop,value tsValue)tsValue{
	if value.kind!=tsErrorKind{return value};source:=(*tsRuntimeError)(value.ref);fields:=map[string]tsValue{};for name,item:=range source.fields{fields[name]=item};cause:=fields["cause"];if cause.kind==tsErrorKind{fields["cause"]=tsAssertErrorCopy(loop,cause)};return tsErrorValue(&tsRuntimeError{name:source.name,message:source.message,fields:fields})
}
func tsAssertErrorInspect(loop *tsLoop,value tsValue)string{
	module:=tsInternalInspectModule(loop);inspect:=tsGet(loop,module,tsStringReference(tsStringUTF8("inspect")));options:=tsNewObject()
	options.set("compact",tsBooleanValue(false));options.set("customInspect",tsBooleanValue(false));options.set("depth",tsNumberValue(1000));options.set("maxArrayLength",tsNull);options.set("showHidden",tsBooleanValue(false));options.set("showProxy",tsBooleanValue(false));options.set("sorted",tsBooleanValue(true));options.set("getters",tsBooleanValue(true))
	return tsText(tsCallReceiver(loop,inspect,module,value,tsObjectValue(options)))
}
func tsAssertErrorColors(loop *tsLoop)map[string][]uint16{
	module:=tsInternalUtilColorsModule(loop);refresh:=tsGet(loop,module,tsStringReference(tsStringUTF8("refresh")));if refresh.kind==tsFunctionKind{tsCallReceiver(loop,refresh,module)}
	return tsAssertDiffColors(loop,module)
}
func tsAssertErrorCurrentColors(loop *tsLoop)map[string][]uint16{return tsAssertDiffColors(loop,tsInternalUtilColorsModule(loop))}
func tsAssertErrorInspectLines(value string)[]string{return strings.Split(value,"\n")}
func tsAssertErrorGetMessage(operator string,message tsValue)string{if tsTruthy(message){return tsText(message)};if value,ok:=tsAssertReadableOperators[operator];ok{return value};return "undefined"}
func tsAssertErrorCheckOperator(actual,expected tsValue,operator string)string{
	if operator=="strictEqual"&&((tsAssertErrorObjectLike(actual)&&!tsIsNull(actual)&&tsAssertErrorObjectLike(expected)&&!tsIsNull(expected))||actual.kind==tsFunctionKind&&expected.kind==tsFunctionKind){return "strictEqualObject"};return operator
}
func tsAssertErrorStackedDiff(actual,expected string,isStringComparison bool,loop *tsLoop,colors map[string][]uint16)string{
	out:="\n"+tsText(tsStringReference(tsStringUnits(colors["green"])))+"+"+tsText(tsStringReference(tsStringUnits(colors["white"])))+" "+actual+"\n"+tsText(tsStringReference(tsStringUnits(colors["red"])))+"- "+tsText(tsStringReference(tsStringUnits(colors["white"])))+expected
	if !isStringComparison{return out}
	process:=tsNodeProcessModule(loop);stderr:=tsGet(loop,process,tsStringReference(tsStringUTF8("stderr")));terminalLength:=80;if tsTruthy(tsGet(loop,stderr,tsStringReference(tsStringUTF8("isTTY")))){if columns:=tsGet(loop,stderr,tsStringReference(tsStringUTF8("columns")));tsIsNumeric(columns){terminalLength=int(tsNumber(columns))}}
	a,b:=tsStringUTF8(actual).units,tsStringUTF8(expected).units;if len(a)+len(b)>terminalLength{return out};indicator:=-1;for index:=range a{if index>=len(b)||a[index]!=b[index]{if index>=3{indicator=index};break}};if indicator>=0{out+= "\n"+strings.Repeat(" ",indicator+2)+"^"};return out
}
func tsAssertErrorSimpleDiff(loop *tsLoop,actual,expected tsValue,inspectedActual,inspectedExpected string,colors map[string][]uint16)(string,string,bool,bool){
	stringsLength:=len([]uint16(tsStringUTF8(inspectedActual).units))+len([]uint16(tsStringUTF8(inspectedExpected).units));if actual.kind==tsStringKind{stringsLength-=2};if expected.kind==tsStringKind{stringsLength-=2}
	if stringsLength<=12&&!(tsIsNumeric(actual)&&tsNumber(actual)==0&&tsIsNumeric(expected)&&tsNumber(expected)==0){return inspectedActual+" !== "+inspectedExpected,"",true,false}
	if actual.kind==tsStringKind&&expected.kind==tsStringKind {
		colorsModule:=tsInternalUtilColorsModule(loop);if tsTruthy(tsGet(loop,colorsModule,tsStringReference(tsStringUTF8("hasColors")))){
			myers:=tsInternalAssertMyersModule(loop);renderer:=tsGet(loop,myers,tsStringReference(tsStringUTF8("printSimpleMyersDiff")))
			var rendered tsValue;failed:=false;func(){defer func(){if recover()!=nil{failed=true}}();rendered=tsCall(loop,renderer,tsStringReference(tsStringUTF8(inspectedActual)),tsStringReference(tsStringUTF8(inspectedExpected)),tsObjectValue(func()*tsObject{o:=tsNewObject();for _,name:=range []string{"green","red","white","blue"}{v:=tsGet(loop,colorsModule,tsStringReference(tsStringUTF8(name)));o.set(name,v)};return o}()))}()
			if !failed {header:=tsStringValueWithLoop(loop,tsGet(loop,colorsModule,tsStringReference(tsStringUTF8("green")))).String()+"actual"+tsStringValueWithLoop(loop,tsGet(loop,colorsModule,tsStringReference(tsStringUTF8("white")))).String()+" "+tsStringValueWithLoop(loop,tsGet(loop,colorsModule,tsStringReference(tsStringUTF8("red")))).String()+"expected"+tsStringValueWithLoop(loop,tsGet(loop,colorsModule,tsStringReference(tsStringUTF8("white")))).String();return tsText(rendered),header,true,false}
		}
	}
	return tsAssertErrorStackedDiff(inspectedActual,inspectedExpected,true,loop,colors),"",false,false
}
func tsAssertErrorIsSimple(actual,expected tsValue,actualLines,expectedLines []string)bool{if len(actualLines)>1||len(expectedLines)>1{return false};return !tsAssertErrorObjectLike(actual)||tsIsNull(actual)||!tsAssertErrorObjectLike(expected)||tsIsNull(expected)}
func tsAssertErrorCreateDiff(loop *tsLoop,actual,expected,message tsValue,operator,diffType string)string{
	operator=tsAssertErrorCheckOperator(actual,expected,operator);colors:=tsAssertErrorCurrentColors(loop);inspectedActual,inspectedExpected:=tsAssertErrorInspect(loop,actual),tsAssertErrorInspect(loop,expected);actualLines,expectedLines:=tsAssertErrorInspectLines(inspectedActual),tsAssertErrorInspectLines(inspectedExpected);showSimple:=tsAssertErrorIsSimple(actual,expected,actualLines,expectedLines)
	header:="+ actual - expected";if len(colors["green"])>0||len(colors["red"])>0{header=tsText(tsStringReference(tsStringUnits(colors["green"])))+"+ actual"+tsText(tsStringReference(tsStringUnits(colors["white"])))+" "+tsText(tsStringReference(tsStringUnits(colors["red"])))+"- expected"+tsText(tsStringReference(tsStringUnits(colors["white"]))) }
	messageText:="";skipped:=false
	if showSimple {var setHeader bool;messageText,header,setHeader,_=tsAssertErrorSimpleDiff(loop,actual,expected,actualLines[0],expectedLines[0],colors);if !setHeader{header="+ actual - expected";if len(colors["green"])>0||len(colors["red"])>0{header=tsText(tsStringReference(tsStringUnits(colors["green"])))+"+ actual"+tsText(tsStringReference(tsStringUnits(colors["white"])))+" "+tsText(tsStringReference(tsStringUnits(colors["red"])))+"- expected"+tsText(tsStringReference(tsStringUnits(colors["white"])))}}
	}else if inspectedActual==inspectedExpected {
		operator="notIdentical";if len(actualLines)>50&&diffType!="full"{messageText=strings.Join(actualLines[:50],"\n")+"\n...}";skipped=true}else{messageText=strings.Join(actualLines,"\n")};header=""
	}else{
		comma:=!tsNullish(actual)&&tsAssertErrorObjectLike(actual);myers:=tsInternalAssertMyersModule(loop);renderer:=tsGet(loop,myers,tsStringReference(tsStringUTF8("printMyersDiff")));var result tsValue;failed:=false
		func(){defer func(){if recover()!=nil{failed=true}}();colorsObject:=tsNewObject();for _,name:=range []string{"green","red","white","blue"}{colorsObject.set(name,tsGet(loop,tsInternalUtilColorsModule(loop),tsStringReference(tsStringUTF8(name))))};result=tsCall(loop,renderer,tsStringReference(tsStringUTF8(inspectedActual)),tsStringReference(tsStringUTF8(inspectedExpected)),tsBooleanValue(comma),tsObjectValue(colorsObject))}()
		if failed||result.kind!=tsObjectKind{if len(actualLines)>50{actualLines=actualLines[:50]};messageText=strings.Join(actualLines,"\n")+"\n...";header="";skipped=true}else{messageText=tsText(tsGet(loop,result,tsStringReference(tsStringUTF8("message"))));skipped=tsTruthy(tsGet(loop,result,tsStringReference(tsStringUTF8("skipped"))))}
	}
	headerMessage:=tsAssertErrorGetMessage(operator,message)+"\n"+header;skippedMessage:="";if skipped{skippedMessage="\n... Skipped lines"};return headerMessage+skippedMessage+"\n"+messageText+"\n"
}
func tsAssertErrorAddEllipsis(value string)string{lines:=strings.SplitN(value,"\n",11);if len(lines)>10{return strings.Join(lines[:10],"\n")+"\n..."};units:=tsStringUTF8(value).units;if len(units)>512{return tsText(tsStringReference(tsStringUnits(units[512:])))+"..."};return value}
func tsAssertErrorNew(loop *tsLoop,options tsValue)tsValue{
	validate:=tsGet(loop,tsInternalValidatorsModule(loop),tsStringReference(tsStringUTF8("validateObject")));if validate.kind==tsFunctionKind{tsCall(loop,validate,options,tsStringReference(tsStringUTF8("options")))}
	message,operator:=tsGet(loop,options,tsStringReference(tsStringUTF8("message"))),tsText(tsGet(loop,options,tsStringReference(tsStringUTF8("operator"))));actual,expected:=tsGet(loop,options,tsStringReference(tsStringUTF8("actual"))),tsGet(loop,options,tsStringReference(tsStringUTF8("expected")));diff:=tsGet(loop,options,tsStringReference(tsStringUTF8("diff")));if tsIsUndefined(diff){diff=tsStringReference(tsStringUTF8("simple"))};diffType:=tsText(diff);messageText:=""
	custom:=operator=="deepStrictEqual"||operator=="strictEqual"||operator=="partialDeepStrictEqual"
	if !tsNullish(message){if custom{messageText=tsAssertErrorCreateDiff(loop,actual,expected,message,operator,diffType)}else{messageText=tsAssertErrorString(loop,message)}}else{
		colors:=tsAssertErrorColors(loop);_ = colors
		if actual.kind==tsErrorKind&&expected.kind==tsErrorKind{actual=tsAssertErrorCopy(loop,actual);expected=tsAssertErrorCopy(loop,expected)}
		if custom{messageText=tsAssertErrorCreateDiff(loop,actual,expected,tsU,operator,diffType)}else if operator=="notDeepStrictEqual"||operator=="notStrictEqual"{
			base,ok:=tsAssertReadableOperators[operator];if !ok{base="undefined"};actualText:=tsAssertErrorInspect(loop,actual);actualLines:=strings.Split(actualText,"\n");if operator=="notStrictEqual"&&(tsAssertErrorObjectLike(actual)&&!tsIsNull(actual)||actual.kind==tsFunctionKind){base=tsAssertReadableOperators["notStrictEqualObject"]};if len(actualLines)>50&&diffType!="full"{actualLines=append(actualLines[:46],"..." );actualLines=actualLines[:47]};if len(actualLines)==1{sep:=" ";if len(actualLines[0])>5{sep="\n\n"};messageText=base+sep+actualLines[0]}else{messageText=base+"\n\n"+strings.Join(actualLines,"\n")+"\n"}
		}else{
			actualText,expectedText:=tsAssertErrorInspect(loop,actual),tsAssertErrorInspect(loop,expected);base,known:=tsAssertReadableOperators[operator];if !known{base="undefined"};if operator=="notDeepEqual"&&actualText==expectedText{messageText=base+"\n\n"+actualText;if len(messageText)>1024&&diffType!="full"{messageText=tsText(tsStringReference(tsStringUnits(tsStringUTF8(messageText).units[:1021])))+"..."}}else{if len(tsStringUTF8(actualText).units)>512&&diffType!="full"{actualText=tsText(tsStringReference(tsStringUnits(tsStringUTF8(actualText).units[:509])))+"..."};if len(tsStringUTF8(expectedText).units)>512&&diffType!="full"{expectedText=tsText(tsStringReference(tsStringUnits(tsStringUTF8(expectedText).units[:509])))+"..."};if operator=="deepEqual"{actualText=base+"\n\n"+actualText+"\n\nshould loosely deep-equal\n\n"}else if _,ok:=tsAssertReadableOperators[operator+"Unequal"];ok{actualText=tsAssertReadableOperators[operator+"Unequal"]+"\n\n"+actualText+"\n\nshould not loosely deep-equal\n\n"}else{expectedText=" "+operator+" "+expectedText};messageText=actualText+expectedText}
		}
	}
	generated:=!tsTruthy(message);fields:=map[string]tsValue{"code":tsStringReference(tsStringUTF8("ERR_ASSERTION")),"generatedMessage":tsBooleanValue(generated),"diff":diff}
	details:=tsGet(loop,options,tsStringReference(tsStringUTF8("details")));if details.kind==tsArrayKind{for index,item:=range (*tsArray)(details.ref).values{for _,name:=range []string{"message","actual","expected","operator","stack"}{key:=name+" "+strconv.Itoa(index);fields[key]=tsGet(loop,item,tsStringReference(tsStringUTF8(name)))}}}else{fields["actual"],fields["expected"],fields["operator"]=actual,expected,tsStringReference(tsStringUTF8(operator))}
	runtimeError:=&tsRuntimeError{name:"AssertionError",message:messageText,fields:fields};value:=tsErrorValue(runtimeError);object:=tsErrorObject(loop,runtimeError)
	object.set("toString",tsNamedNativeMethod("toString",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{name:=tsGet(loop,receiver,tsStringReference(tsStringUTF8("name")));code:=tsGet(loop,receiver,tsStringReference(tsStringUTF8("code")));message:=tsGet(loop,receiver,tsStringReference(tsStringUTF8("message")));return tsStringReference(tsStringUTF8(tsText(name)+" ["+tsText(code)+"]: "+tsText(message)))}))
	inspectModule:=tsInternalInspectModule(loop);inspect:=tsGet(loop,inspectModule,tsStringReference(tsStringUTF8("inspect")));custom:=tsGet(loop,inspect,tsStringReference(tsStringUTF8("custom")));if custom.kind==tsSymbolKind{symbolKey:=(*tsSymbol)(custom.ref).key;object.set(symbolKey,tsNamedNativeMethod("[inspect.custom]",2,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if receiver.kind!=tsErrorKind{return tsCall(loop,inspect,receiver,tsObjectValue(func()*tsObject{o:=tsNewObject();o.set("customInspect",tsBooleanValue(false));o.set("depth",tsNumberValue(0));return o}()))};err:=(*tsRuntimeError)(receiver.ref);originalActual,originalExpected:=tsGet(loop,receiver,tsStringReference(tsStringUTF8("actual"))),tsGet(loop,receiver,tsStringReference(tsStringUTF8("expected")));if originalActual.kind==tsStringKind{tsSet(loop,receiver,tsStringReference(tsStringUTF8("actual")),tsStringReference(tsStringUTF8(tsAssertErrorAddEllipsis(tsText(originalActual)))))};if originalExpected.kind==tsStringKind{tsSet(loop,receiver,tsStringReference(tsStringUTF8("expected")),tsStringReference(tsStringUTF8(tsAssertErrorAddEllipsis(tsText(originalExpected)))))};options:=tsNewObject();options.set("customInspect",tsBooleanValue(false));options.set("depth",tsNumberValue(0));result:=tsCallReceiver(loop,inspect,inspectModule,receiver,tsObjectValue(options));if originalActual.kind!=tsUndefinedKind{tsSet(loop,receiver,tsStringReference(tsStringUTF8("actual")),originalActual))};if originalExpected.kind!=tsUndefinedKind{tsSet(loop,receiver,tsStringReference(tsStringUTF8("expected")),originalExpected)};_ = err;return result}))}
	return value
}
func tsInternalAssertionErrorConstructor(loop *tsLoop)tsValue{
	tsInternalAssertionErrorState.Do(func(){fn:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsAssertErrorNew(loop,tsArg(args,0))});fn.name="AssertionError";fn.length=1;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsAssertErrorNew(loop,tsArg(args,0))};fn.properties=tsNewObject();prototype:=tsNewObject();fn.properties.set("prototype",tsObjectValue(prototype));fn.properties.set("name",tsStringReference(tsStringUTF8("AssertionError")));tsInternalAssertionErrorState.constructor=tsFunctionValue(fn)});return tsInternalAssertionErrorState.constructor
}
`

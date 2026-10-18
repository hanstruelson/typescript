package goemit

// nodeParseArgsImplementation translates parse_args.rs's validation,
// tokenization, value accumulation, defaults, and optional token log.
const nodeParseArgsImplementation = `
func tsNodeParseSlice(value string,start,end int)string{runes:=[]rune(value);if start<0{start+=len(runes)};if start<0{start=0};if start>len(runes){start=len(runes)};if end<0{end+=len(runes)};if end<start{end=start};if end>len(runes){end=len(runes)};return string(runes[start:end])}

func tsNodeParseAdd(loop *tsLoop, token tsNodeParseToken, options []tsNodeParseOption, strict, allowPositionals bool, values *tsObject, positionals, tokens *tsArray) {
	if token.kind == tsNodeParseOptionKind {
		option := token.option
		if strict { tsNodeParseValidateOption(option, options, allowPositionals) }
		if option.name != "__proto__" {
			value := option.value
			if tsIsUndefined(value) { value = tsBooleanValue(!option.negative) }
			if option.optionIndex >= 0 && options[option.optionIndex].multiple {
				current := tsGet(loop, tsObjectValue(values), tsStringReference(tsStringUTF8(option.name)))
				if tsIsUndefined(current) { values.set(option.name, tsArrayValue(&tsArray{values:[]tsValue{value}})) } else { (*tsArray)(current.ref).appendItems(loop, value) }
			} else { values.set(option.name, value) }
		}
	} else if token.kind == tsNodeParsePositional {
		if !allowPositionals { tsNodeParseError("ERR_PARSE_ARGS_UNEXPECTED_POSITIONAL", "Unexpected argument '"+tsText(token.value)+"'. This command does not take positional arguments") }
		positionals.appendItems(loop, token.value)
	}
	if tokens == nil { return }
	entry := tsNewObject()
	switch token.kind {
	case tsNodeParseOptionKind:
		entry.set("kind", tsStringReference(tsStringUTF8("option")))
		entry.set("name", tsStringReference(tsStringUTF8(token.option.name)))
		entry.set("rawName", tsStringReference(tsStringUTF8(tsNodeParseRawName(token.option))))
		entry.set("index", tsNumberValue(float64(token.option.index)))
		entry.set("value", token.option.value)
		inline := tsU
		if !tsIsUndefined(token.option.value) { inline = tsBooleanValue(token.option.inlineValue) }
		entry.set("inlineValue", inline)
	case tsNodeParsePositional:
		entry.set("kind", tsStringReference(tsStringUTF8("positional")))
		entry.set("index", tsNumberValue(float64(token.index)))
		entry.set("value", token.value)
	case tsNodeParseTerminator:
		entry.set("kind", tsStringReference(tsStringUTF8("option-terminator")))
		entry.set("index", tsNumberValue(float64(token.index)))
	}
	tokens.appendItems(loop, tsObjectValue(entry))
}
func tsNodeParseEmit(loop *tsLoop, index int, raw, name string, value tsValue, inline bool, optionIndex int, negative bool, groupIndex int, subtype tsNodeParseTokenSubtype, options []tsNodeParseOption, strict, allowPositionals bool, values *tsObject, positionals, tokens *tsArray) {
	option := tsNodeParseOptionToken{index:index,name:name,raw:raw,value:value,inlineValue:inline,optionIndex:optionIndex,negative:negative,groupIndex:groupIndex,subtype:subtype}
	tsNodeParseAdd(loop,tsNodeParseToken{kind:tsNodeParseOptionKind,index:index,option:option},options,strict,allowPositionals,values,positionals,tokens)
}
func tsNodeParseTokenize(loop *tsLoop, args []tsValue, options []tsNodeParseOption, allowNegative, strict, allowPositionals bool, values *tsObject, positionals, tokens *tsArray) {
	for index := 0; index < len(args); index++ {
		raw := tsStringValueWithLoop(loop,args[index]).String()
		switch tsNodeParseClassifyToken(raw,options) {
		case tsNodeParseSubtypeTerminator:
			tsNodeParseAdd(loop,tsNodeParseToken{kind:tsNodeParseTerminator,index:index},options,strict,allowPositionals,values,positionals,tokens)
			for index++; index<len(args); index++ { tsNodeParseAdd(loop,tsNodeParseToken{kind:tsNodeParsePositional,index:index,value:args[index]},options,strict,allowPositionals,values,positionals,tokens) }
			return
		case tsNodeParseSubtypeLoneShort:
			tokenIndex:=index; short:=tsNodeParseSlice(raw,1,2); optionIndex:=tsNodeParseFindOptionByShortName(short,options); value,inline:=tsU,true
			if tsNodeParseOptionType(optionIndex,options)=="string"&&index+1<len(args){index++;value,inline=args[index],false}
			tsNodeParseEmit(loop,tokenIndex,raw,tsNodeParseShortName(short,options),value,inline,optionIndex,false,-1,tsNodeParseSubtypeLoneShort,options,strict,allowPositionals,values,positionals,tokens)
		case tsNodeParseSubtypeShortGroup:
			original:=index; runes:=[]rune(raw)
			for group:=1;group<len(runes);group++ { short:=string(runes[group]); optionIndex:=tsNodeParseFindOptionByShortName(short,options); valueType:=tsNodeParseOptionType(optionIndex,options)
				if valueType!="string"||group==len(runes)-1 { value,inline:=tsU,true;if valueType=="string"&&index+1<len(args){index++;value,inline=args[index],false};tsNodeParseEmit(loop,original,raw,tsNodeParseShortName(short,options),value,inline,optionIndex,false,group,tsNodeParseSubtypeLoneShort,options,strict,allowPositionals,values,positionals,tokens)
				} else { value:=tsStringReference(tsStringUTF8(tsNodeParseSlice(raw,group+1,len(runes))));tsNodeParseEmit(loop,original,raw,tsNodeParseShortName(short,options),value,true,optionIndex,false,group,tsNodeParseSubtypeShortValue,options,strict,allowPositionals,values,positionals,tokens);break }
			}
		case tsNodeParseSubtypeShortValue:
			short:=tsNodeParseSlice(raw,1,2);optionIndex:=tsNodeParseFindOptionByShortName(short,options);value:=tsStringReference(tsStringUTF8(tsNodeParseSlice(raw,2,len([]rune(raw)))))
			tsNodeParseEmit(loop,index,tsNodeParseSlice(raw,0,2),tsNodeParseShortName(short,options),value,true,optionIndex,false,-1,tsNodeParseSubtypeShortValue,options,strict,allowPositionals,values,positionals,tokens)
		case tsNodeParseSubtypeLoneLong:
			tokenIndex:=index;name:=tsNodeParseSlice(raw,2,len([]rune(raw)));negative:=false
			if allowNegative&&strings.HasPrefix(name,"no-"){name,negative=strings.TrimPrefix(name,"no-"),true}
			optionIndex:=tsNodeParseLongOption(name,options);value,inline:=tsU,true
			if tsNodeParseOptionType(optionIndex,options)=="string"&&index+1<len(args)&&!negative{index++;value,inline=args[index],false}
			tsNodeParseEmit(loop,tokenIndex,raw,name,value,inline,optionIndex,negative,-1,tsNodeParseSubtypeLoneLong,options,strict,allowPositionals,values,positionals,tokens)
		case tsNodeParseSubtypeLongValue:
			equal:=strings.IndexByte(raw,'=');name,text:=raw[2:equal],raw[equal+1:];value:=tsStringReference(tsStringUTF8(text))
			tsNodeParseEmit(loop,index,raw,name,value,true,tsNodeParseLongOption(name,options),false,-1,tsNodeParseSubtypeLongValue,options,strict,allowPositionals,values,positionals,tokens)
		case tsNodeParseSubtypePositional:
			tsNodeParseAdd(loop,tsNodeParseToken{kind:tsNodeParsePositional,index:index,value:args[index]},options,strict,allowPositionals,values,positionals,tokens)
		}
	}
}
func tsNodeParseDefaultArgs(loop *tsLoop) []tsValue {
	process:=tsNodeProcessModule(loop);argv:=tsGet(loop,process,tsStringReference(tsStringUTF8("argv")));execArgv:=tsGet(loop,process,tsStringReference(tsStringUTF8("execArgv")))
	if argv.kind!=tsArrayKind||execArgv.kind!=tsArrayKind{return nil}
	eval:=false;for _,value:=range (*tsArray)(execArgv.ref).values{if value.kind==tsStringKind{switch tsText(value){case "-e","--eval","-p","--print":eval=true}}}
	input:=(*tsArray)(argv.ref);start:=2;if eval{start=1};if start>input.length(){start=input.length()}
	return append([]tsValue(nil),input.values[start:]...)
}
func tsNodeParseArgs(loop *tsLoop,config tsValue) tsValue {
	if !tsIsUndefined(config){tsInternalValidatorObject(config,tsStringReference(tsStringUTF8("config")),tsU)}
	argsValue:=tsNodeParseOwn(loop,config,"args");var args []tsValue
	if !tsNullish(argsValue){tsInternalValidatorArray(argsValue,tsStringReference(tsStringUTF8("args")),tsU);args=append([]tsValue(nil),(*tsArray)(argsValue.ref).values...)}else{args=tsNodeParseDefaultArgs(loop)}
	strictValue:=tsNodeParseOwn(loop,config,"strict");if tsNullish(strictValue){strictValue=tsBooleanValue(true)};tsNodeParseCheckBoolean("strict",strictValue);strict:=tsTruthy(strictValue)
	allowPositionalsValue:=tsNodeParseOwn(loop,config,"allowPositionals");if tsNullish(allowPositionalsValue){allowPositionalsValue=tsBooleanValue(!strict)};tsNodeParseCheckBoolean("allowPositionals",allowPositionalsValue)
	tokensValue:=tsNodeParseOwn(loop,config,"tokens");if tsNullish(tokensValue){tokensValue=tsBooleanValue(false)};tsNodeParseCheckBoolean("tokens",tokensValue)
	negativeValue:=tsNodeParseOwn(loop,config,"allowNegative");if tsNullish(negativeValue){negativeValue=tsBooleanValue(false)};tsNodeParseCheckBoolean("allowNegative",negativeValue)
	options:=tsNodeParseReadOptions(loop,tsNodeParseOwn(loop,config,"options"))
	values:=tsNewObject();values.prototype=tsNull;positionals:=&tsArray{};var tokens *tsArray;if tsTruthy(tokensValue){tokens=&tsArray{}}
	tsNodeParseTokenize(loop,args,options,tsTruthy(negativeValue),strict,tsTruthy(allowPositionalsValue),values,positionals,tokens)
	for _,option:=range options{if option.hasDefault&&option.longName!="__proto__"&&!tsObjectHasOwn(tsObjectValue(values),option.longName){values.set(option.longName,option.defaultValue)}}
	result:=tsNewObject();result.set("values",tsObjectValue(values));result.set("positionals",tsArrayValue(positionals));if tokens!=nil{result.set("tokens",tsArrayValue(tokens))}
	return tsObjectValue(result)
}
func tsNodeParseArgsFunction() tsValue {
	return tsNamedNativeMethod("parseArgs",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeParseArgs(loop,tsArg(args,0))})
}
`

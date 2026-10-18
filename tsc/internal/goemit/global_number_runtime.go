package goemit

const globalNumberRuntime = `
func tsGlobalIsFiniteFunction()tsValue{return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
	value:=tsNumber(tsPrimitiveHint(loop,tsArg(args,0),false))
	return tsBooleanValue(!math.IsNaN(value)&&!math.IsInf(value,0))
}))}

func tsGlobalParseIntFunction()tsValue{return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
	text:=strings.TrimSpace(tsStringValueWithLoop(loop,tsArg(args,0)).String())
	text=strings.TrimLeft(text,"\uFEFF")
	if text==""{return tsNumberValue(math.NaN())}
	negative:=false
	if text[0]=='+'||text[0]=='-'{negative=text[0]=='-';text=text[1:]}
	radixValue:=tsArg(args,1);radix:=int32(0)
	if !tsIsUndefined(radixValue){radix=int32(tsToUint32(tsPrimitiveHint(loop,radixValue,false)))}
	if radix==0{if strings.HasPrefix(text,"0x")||strings.HasPrefix(text,"0X"){radix=16;text=text[2:]}else{radix=10}}else if radix==16&&(strings.HasPrefix(text,"0x")||strings.HasPrefix(text,"0X")){text=text[2:]}
	if radix<2||radix>36{return tsNumberValue(math.NaN())}
	digitCount:=0
	for digitCount<len(text){character:=text[digitCount];valid:=character>='0'&&character<='9'&&int32(character-'0')<radix||character>='a'&&character<='z'&&int32(character-'a'+10)<radix||character>='A'&&character<='Z'&&int32(character-'A'+10)<radix;if !valid{break};digitCount++}
	if digitCount==0{return tsNumberValue(math.NaN())}
	integer,ok:=new(big.Int).SetString(text[:digitCount],int(radix));if !ok{return tsNumberValue(math.NaN())};if negative{integer.Neg(integer)}
	value,_:=new(big.Float).SetInt(integer).Float64()
	return tsNumberValue(value)
}))}

func tsGlobalParseFloatFunction()tsValue{return tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
	text:=strings.TrimSpace(tsStringValueWithLoop(loop,tsArg(args,0)).String())
	text=strings.TrimLeft(text,"\uFEFF")
	if strings.HasPrefix(text,"Infinity")||strings.HasPrefix(text,"+Infinity"){return tsNumberValue(math.Inf(1))}
	if strings.HasPrefix(text,"-Infinity"){return tsNumberValue(math.Inf(-1))}
	index:=0;if index<len(text)&&(text[index]=='+'||text[index]=='-'){index++}
	integerDigits:=0;for index<len(text)&&text[index]>='0'&&text[index]<='9'{index++;integerDigits++}
	if index<len(text)&&text[index]=='.'{index++;for index<len(text)&&text[index]>='0'&&text[index]<='9'{index++;integerDigits++}}
	if integerDigits==0{return tsNumberValue(math.NaN())}
	end:=index
	if index<len(text)&&(text[index]=='e'||text[index]=='E'){
		exponent:=index+1;if exponent<len(text)&&(text[exponent]=='+'||text[exponent]=='-'){exponent++};digits:=exponent;for exponent<len(text)&&text[exponent]>='0'&&text[exponent]<='9'{exponent++};if exponent>digits{end=exponent}
	}
	value,err:=strconv.ParseFloat(text[:end],64);if err!=nil{if numberError,ok:=err.(*strconv.NumError);ok&&numberError.Err==strconv.ErrRange{return tsNumberValue(value)};return tsNumberValue(math.NaN())};return tsNumberValue(value)
}))}
`

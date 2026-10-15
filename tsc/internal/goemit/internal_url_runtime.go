package goemit

// internalURLRuntime translates internal/url.ts and adapts NodeURL.cpp's IDNA
// boundary with the local RFC 3492 implementation and Go's URL parser.
const internalURLRuntime = `
var tsInternalURLSingleton struct { sync.Once; value *tsObject }
func tsInternalURLMapDomain(text string,decode bool)(result string,ok bool){
	defer func(){if recover()!=nil{result="";ok=false}}()
	mapped:=tsNodePunycodeMapDomain(tsStringUTF8(text),decode).String()
	if !decode{mapped=strings.ToLower(mapped)}
	return mapped,true
}
func tsInternalURLHost(text string)(string,bool){
	end:=len(text);if at:=strings.IndexAny(text,"/?#\\");at>=0{end=at};host:=text[:end];if host==""{return "",false}
	if strings.HasPrefix(host,"["){if !strings.HasSuffix(host,"]"){return "",false}}else if strings.ContainsAny(host,":@%"){return "",false}
	parsed,err:=url.Parse("ws://"+host+"/");if err!=nil||parsed.Host==""{return "",false};ascii,ok:=tsInternalURLMapDomain(parsed.Host,true);_ = ascii
	if !ok{return "",false};ascii,ok=tsInternalURLMapDomain(parsed.Host,false);if !ok{return "",false}
	if strings.Contains(ascii,"xn--"){decoded,ok:=tsInternalURLMapDomain(ascii,true);if !ok{return "",false};roundTrip,ok:=tsInternalURLMapDomain(decoded,false);if !ok||!strings.EqualFold(roundTrip,ascii){return "",false}}
	return ascii,true
}
func tsInternalURLDomainASCII(loop *tsLoop,args []tsValue)tsValue{if len(args)==0{panic(tsThrown{tsFSException("TypeError","ERR_MISSING_ARGS","The \"domain\" argument must be specified")})};host,ok:=tsInternalURLHost(tsText(args[0]));if !ok{return tsStringReference(tsStringUTF8(""))};return tsStringReference(tsStringUTF8(host))}
func tsInternalURLDomainUnicode(loop *tsLoop,args []tsValue)tsValue{if len(args)==0{panic(tsThrown{tsFSException("TypeError","ERR_MISSING_ARGS","The \"domain\" argument must be specified")})};host,ok:=tsInternalURLHost(tsText(args[0]));if !ok{return tsStringReference(tsStringUTF8(""))};unicode,ok:=tsInternalURLMapDomain(host,true);if !ok{return tsStringReference(tsStringUTF8(""))};return tsStringReference(tsStringUTF8(unicode))}
func tsInternalURLIDNAToASCII(loop *tsLoop,args []tsValue)tsValue{value:=tsText(tsArg(args,0));result,ok:=tsInternalURLMapDomain(value,false);if !ok{return tsStringReference(tsStringUTF8(""))};return tsStringReference(tsStringUTF8(result))}
func tsInternalURLToHTTPOptions(loop *tsLoop,args []tsValue)tsValue{
	value:=tsArg(args,0);out:=tsNewObject();for _,key:=range tsObjectKeys(value).values{name:=tsPropertyKey(key);out.set(name,tsGet(loop,value,tsStringReference(tsStringUTF8(name))))}
	for _,name:=range []string{"protocol","hash","search","pathname","href"}{out.set(name,tsGet(loop,value,tsStringReference(tsStringUTF8(name))))}
	hostname:=tsGet(loop,value,tsStringReference(tsStringUTF8("hostname")));if hostname.kind==tsStringKind{text:=tsText(hostname);if strings.HasPrefix(text,"[")&&strings.HasSuffix(text,"]"){hostname=tsStringReference(tsStringUTF8(text[1:len(text)-1]))}};out.set("hostname",hostname)
	search:=tsGet(loop,value,tsStringReference(tsStringUTF8("search")));pathname:=tsGet(loop,value,tsStringReference(tsStringUTF8("pathname")));pathText:="";if !tsNullish(pathname){pathText=tsText(pathname)};if !tsNullish(search){pathText+=tsText(search)};out.set("path",tsStringReference(tsStringUTF8(pathText)))
	port:=tsGet(loop,value,tsStringReference(tsStringUTF8("port")));if port.kind==tsStringKind&&tsText(port)!=""{n,err:=strconv.Atoi(tsText(port));if err==nil{out.set("port",tsNumberValue(float64(n)))}}
	username:=tsGet(loop,value,tsStringReference(tsStringUTF8("username")));password:=tsGet(loop,value,tsStringReference(tsStringUTF8("password")));if tsTruthy(username)||tsTruthy(password){user,_:=url.PathUnescape(tsText(username));pass,_:=url.PathUnescape(tsText(password));out.set("auth",tsStringReference(tsStringUTF8(user+":"+pass)))}
	return tsObjectValue(out)
}
func tsInternalURLModule(loop *tsLoop)tsValue{
	tsInternalURLSingleton.Do(func(){module:=tsNewObject();module.set("domainToASCII",tsNamedNativeMethod("domainToASCII",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalURLDomainASCII(loop,args)}));module.set("domainToUnicode",tsNamedNativeMethod("domainToUnicode",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalURLDomainUnicode(loop,args)}));module.set("idnaToASCII",tsNamedNativeMethod("idnaToASCII",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalURLIDNAToASCII(loop,args)}));module.set("urlToHttpOptions",tsNamedNativeMethod("urlToHttpOptions",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalURLToHTTPOptions(loop,args)}));tsInternalURLSingleton.value=module})
	return tsObjectValue(tsInternalURLSingleton.value)
}
`

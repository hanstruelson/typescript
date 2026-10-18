package goemit

// nodeURLRuntime maps node:url's legacy and WHATWG APIs onto net/url,
// filepath, and the internal IDNA/Buffer adapters.
const nodeURLRuntime = `
var tsNodeURLSingleton struct { sync.Once; module *tsObject; urlClass,paramsClass,patternClass,legacyClass *tsClass }

func tsNodeURLFail(code,message string){panic(tsThrown{tsFSException("TypeError",code,message)})}
func tsNodeURLString(value tsValue)string{return tsText(value)}
func tsNodeURLParse(text string,base *url.URL)(*url.URL,error){
	parsed,err:=url.Parse(text);if err!=nil{return nil,err};if base!=nil{parsed=base.ResolveReference(parsed)}else if !parsed.IsAbs(){return nil,fmt.Errorf("Invalid URL")}
	parsed.Scheme=strings.ToLower(parsed.Scheme);if parsed.Host!=""{hostname:=parsed.Hostname();ascii,ok:=tsInternalURLMapDomain(hostname,false);if ok&&ascii!=""{port:=parsed.Port();if port!=""{parsed.Host=net.JoinHostPort(ascii,port);if strings.Contains(ascii,":"){parsed.Host="["+ascii+"]:"+port}}else if strings.Contains(ascii,":"){parsed.Host="["+ascii+"]"}else{parsed.Host=ascii}}}
	if (parsed.Scheme=="http"&&parsed.Port()=="80")||(parsed.Scheme=="https"&&parsed.Port()=="443")||(parsed.Scheme=="ws"&&parsed.Port()=="80")||(parsed.Scheme=="wss"&&parsed.Port()=="443"){host:=parsed.Hostname();if strings.Contains(host,":"){parsed.Host="["+host+"]"}else{parsed.Host=host}}
	return parsed,nil
}
func tsNodeURLField(record *tsNodeURLRecord,name string)tsValue{
	u:=record.parsed;switch name{
	case "href":return tsStringReference(tsStringUTF8(u.String()))
	case "protocol":if u.Scheme==""{return tsStringReference(tsStringUTF8(":"))};return tsStringReference(tsStringUTF8(u.Scheme+":"))
	case "host":return tsStringReference(tsStringUTF8(u.Host))
	case "hostname":host:=u.Hostname();if strings.Contains(host,":"){host="["+host+"]"};return tsStringReference(tsStringUTF8(host))
	case "port":return tsStringReference(tsStringUTF8(u.Port()))
	case "pathname":return tsStringReference(tsStringUTF8(u.EscapedPath()))
	case "search":if u.RawQuery==""&&!u.ForceQuery{return tsStringReference(tsStringUTF8(""))};return tsStringReference(tsStringUTF8("?"+u.RawQuery))
	case "hash":if u.Fragment==""&&u.RawFragment==""{return tsStringReference(tsStringUTF8(""))};return tsStringReference(tsStringUTF8("#"+u.EscapedFragment()))
	case "username":if u.User==nil{return tsStringReference(tsStringUTF8(""))};return tsStringReference(tsStringUTF8(u.User.Username()))
	case "password":if u.User==nil{return tsStringReference(tsStringUTF8(""))};password,_:=u.User.Password();return tsStringReference(tsStringUTF8(password))
	case "origin":if u.Scheme=="http"||u.Scheme=="https"||u.Scheme=="ws"||u.Scheme=="wss"||u.Scheme=="ftp"{return tsStringReference(tsStringUTF8(u.Scheme+"://"+u.Host))};return tsStringReference(tsStringUTF8("null"))
	case "searchParams":return record.params
	};return tsU
}
func tsNodeURLSync(object *tsObject){if object==nil||object.nodeURL==nil{return};for _,name:=range []string{"href","protocol","host","hostname","port","pathname","search","hash","username","password","origin"}{object.values[name]=tsNodeURLField(object.nodeURL,name)}}
func tsNodeURLSetField(loop *tsLoop,object *tsObject,name string,value tsValue)tsValue{
	if object.nodeURL==nil{return value};record:=object.nodeURL;u:=record.parsed;text:=tsStringValueWithLoop(loop,value).String()
	switch name{
	case "href":next,err:=tsNodeURLParse(text,nil);if err!=nil{tsNodeURLFail("ERR_INVALID_URL",text)};record.parsed=next;if object.nodeURLParams!=nil{object.nodeURLParams.pairs=tsNodeURLParsePairs(next.RawQuery);tsNodeURLParamsChanged(object.nodeURLParams)}
	case "protocol":u.Scheme=strings.TrimSuffix(strings.ToLower(text),":")
	case "hostname":port:=u.Port();host:=strings.Trim(text,"[]");if strings.Contains(host,":"){host="["+host+"]"};if port!=""{host=net.JoinHostPort(strings.Trim(host,"[]"),port);if strings.Contains(strings.Trim(host,"[]"),":"){host="["+strings.Trim(host,"[]")+"]:"+port}};u.Host=host
	case "host":parsed,err:=url.Parse("//"+text);if err!=nil{tsNodeURLFail("ERR_INVALID_URL",text)};u.Host=parsed.Host
	case "port":hostname:=u.Hostname();if text!=""{if _,err:=strconv.Atoi(text);err!=nil{tsNodeURLFail("ERR_INVALID_URL",text)};if strings.Contains(hostname,":"){u.Host="["+hostname+"]:"+text}else{u.Host=hostname+":"+text}}else if strings.Contains(hostname,":"){u.Host="["+hostname+"]"}else{u.Host=hostname}
	case "pathname":parsed,err:=url.Parse(text);if err!=nil{tsNodeURLFail("ERR_INVALID_URL",text)};u.Path=parsed.Path;u.RawPath=parsed.RawPath
	case "search":hasQuery:=strings.HasPrefix(text,"?");text=strings.TrimPrefix(text,"?");u.RawQuery=text;u.ForceQuery=hasQuery&&text==""
	case "hash":text=strings.TrimPrefix(text,"#");u.Fragment=text;u.RawFragment=""
	case "username":password:="";if u.User!=nil{password,_=u.User.Password()};if password==""{u.User=url.User(text)}else{u.User=url.UserPassword(text,password)}
	case "password":username:="";if u.User!=nil{username=u.User.Username()};u.User=url.UserPassword(username,text)
	}
	tsNodeURLSync(object);if name=="search"&&object.nodeURLParams!=nil{object.nodeURLParams.pairs=tsNodeURLParsePairs(u.RawQuery)};return value
}
func tsNodeURLAccessor(object *tsObject,name string){if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};getter:=tsNamedNativeMethod("get "+name,0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if receiver.kind!=tsObjectKind{return tsU};return (*tsObject)(receiver.ref).values[name]});setter:=tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsObjectKind{tsNodeURLFail("ERR_INVALID_THIS","URL accessor called on incompatible receiver")};return tsNodeURLSetField(loop,(*tsObject)(receiver.ref),name,tsArg(args,0))});object.descriptors[name]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:getter,setter:setter}}
func tsNodeURLObject(parsed *url.URL)tsValue{
	class:=tsNodeURLSingleton.urlClass;object:=tsNewObject();object.nativeClass=class;object.prototype=tsObjectValue(class.nativePrototype);object.nodeURL=&tsNodeURLRecord{parsed:parsed}
	for _,name:=range []string{"href","protocol","host","hostname","port","pathname","search","hash","username","password","origin"}{object.set(name,tsU)};tsNodeURLSync(object);for _,name:=range []string{"href","protocol","host","hostname","port","pathname","search","hash","username","password"}{tsNodeURLAccessor(object,name)}
	params:=tsNodeURLParamsObject(tsNodeURLParsePairs(parsed.RawQuery),object);object.nodeURL.params=params;object.nodeURLParams=(*tsObject)(params.ref).nodeURLParams
	object.descriptors["searchParams"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get searchParams",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return object.nodeURL.params})}
	return tsObjectValue(object)
}
func tsNodeURLParsePairs(query string)[][2]string{
	pairs:=[][2]string{};query=strings.TrimPrefix(query,"?");if query==""{return pairs};for _,part:=range strings.Split(query,"&"){if part==""{continue};key,value,has:=strings.Cut(part,"=");if !has{value=""};k,err:=url.QueryUnescape(key);if err!=nil{k=key};v,err:=url.QueryUnescape(value);if err!=nil{v=value};pairs=append(pairs,[2]string{k,v})};return pairs
}
func tsNodeURLFormEscape(value string)string{escaped:=url.QueryEscape(value);return strings.ReplaceAll(escaped,"%2A","*")}
func tsNodeURLParamsString(state *tsNodeURLSearchParamsState)string{parts:=make([]string,0,len(state.pairs));for _,pair:=range state.pairs{parts=append(parts,tsNodeURLFormEscape(pair[0])+"="+tsNodeURLFormEscape(pair[1]))};return strings.Join(parts,"&")}
func tsNodeURLParamsObject(pairs [][2]string,owner *tsObject)tsValue{
	object:=tsNewObject();object.nativeClass=tsNodeURLSingleton.paramsClass;object.prototype=tsObjectValue(tsNodeURLSingleton.paramsClass.nativePrototype);state:=&tsNodeURLSearchParamsState{pairs:pairs,owner:owner,object:object};object.nodeURLParams=state;object.set("size",tsNumberValue(float64(len(pairs))))
	method:=func(name string,length int,call func(*tsLoop,*tsNodeURLSearchParamsState,[]tsValue)tsValue){object.set(name,tsNamedNativeMethod(name,length,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if receiver.kind!=tsObjectKind||(*tsObject)(receiver.ref).nodeURLParams==nil{tsNodeURLFail("ERR_INVALID_THIS","URLSearchParams called on incompatible receiver")};return call(loop,(*tsObject)(receiver.ref).nodeURLParams,args)}))}
	method("append",2,func(_ *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{s.pairs=append(s.pairs,[2]string{tsText(tsArg(args,0)),tsText(tsArg(args,1))});tsNodeURLParamsChanged(s);return tsU})
	method("delete",1,func(_ *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{name:=tsText(tsArg(args,0));value:=tsArg(args,1);kept:=s.pairs[:0];for _,pair:=range s.pairs{if pair[0]!=name||len(args)>1&&pair[1]!=tsText(value){kept=append(kept,pair)}};s.pairs=kept;tsNodeURLParamsChanged(s);return tsU})
	method("get",1,func(_ *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{name:=tsText(tsArg(args,0));for _,pair:=range s.pairs{if pair[0]==name{return tsStringReference(tsStringUTF8(pair[1]))}};return tsNull})
	method("getAll",1,func(_ *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{name:=tsText(tsArg(args,0));out:=&tsArray{};for _,pair:=range s.pairs{if pair[0]==name{out.values=append(out.values,tsStringReference(tsStringUTF8(pair[1])))}};return tsArrayValue(out)})
	method("has",1,func(_ *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{name:=tsText(tsArg(args,0));for _,pair:=range s.pairs{if pair[0]==name&&(len(args)<2||pair[1]==tsText(args[1])){return tsBooleanValue(true)}};return tsBooleanValue(false)})
	method("set",2,func(_ *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{name,value:=tsText(tsArg(args,0)),tsText(tsArg(args,1));next:=[][2]string{};inserted:=false;for _,pair:=range s.pairs{if pair[0]==name{if !inserted{next=append(next,[2]string{name,value});inserted=true};continue};next=append(next,pair)};if !inserted{next=append(next,[2]string{name,value})};s.pairs=next;tsNodeURLParamsChanged(s);return tsU})
	method("sort",0,func(_ *tsLoop,s *tsNodeURLSearchParamsState,_ []tsValue)tsValue{sort.SliceStable(s.pairs,func(i,j int)bool{return s.pairs[i][0]<s.pairs[j][0]});tsNodeURLParamsChanged(s);return tsU})
	method("toString",0,func(_ *tsLoop,s *tsNodeURLSearchParamsState,_ []tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeURLParamsString(s)))})
	method("toJSON",0,func(_ *tsLoop,s *tsNodeURLSearchParamsState,_ []tsValue)tsValue{out:=tsNewObject();for _,pair:=range s.pairs{previous,ok:=out.values[pair[0]];if !ok{out.set(pair[0],tsStringReference(tsStringUTF8(pair[1])));continue};if previous.kind==tsArrayKind{(*tsArray)(previous.ref).values=append((*tsArray)(previous.ref).values,tsStringReference(tsStringUTF8(pair[1])))}else{out.set(pair[0],tsArrayValue(&tsArray{values:[]tsValue{previous,tsStringReference(tsStringUTF8(pair[1]))}}))}};return tsObjectValue(out)})
	method("entries",0,func(_ *tsLoop,s *tsNodeURLSearchParamsState,_ []tsValue)tsValue{return tsIteratorValueReference(tsIterate(tsNodeURLParamArray(s,2)))})
	method("keys",0,func(_ *tsLoop,s *tsNodeURLSearchParamsState,_ []tsValue)tsValue{return tsIteratorValueReference(tsIterate(tsNodeURLParamArray(s,0)))})
	method("values",0,func(_ *tsLoop,s *tsNodeURLSearchParamsState,_ []tsValue)tsValue{return tsIteratorValueReference(tsIterate(tsNodeURLParamArray(s,1)))})
	method("forEach",1,func(loop *tsLoop,s *tsNodeURLSearchParamsState,args []tsValue)tsValue{callback:=tsArg(args,0);if callback.kind!=tsFunctionKind{tsNodeURLFail("ERR_INVALID_ARG_TYPE","The callback argument must be a function")};thisArg:=tsArg(args,1);for _,pair:=range append([][2]string{},s.pairs...){tsCallReceiver(loop,callback,thisArg,tsStringReference(tsStringUTF8(pair[1])),tsStringReference(tsStringUTF8(pair[0])),tsObjectValue(object))};return tsU})
	iterator:=tsPropertyKey(tsWellKnownSymbol("iterator"));object.set(iterator,object.values["entries"]);return tsObjectValue(object)
}
func tsNodeURLParamsChanged(state *tsNodeURLSearchParamsState){if state==nil{return};if state.object!=nil{state.object.values["size"]=tsNumberValue(float64(len(state.pairs)))};if state.owner!=nil&&state.owner.nodeURL!=nil{state.owner.nodeURL.parsed.RawQuery=tsNodeURLParamsString(state);state.owner.nodeURL.parsed.ForceQuery=false;tsNodeURLSync(state.owner)}}
func tsNodeURLParamArray(state *tsNodeURLSearchParamsState,mode int)tsValue{out:=&tsArray{};for _,pair:=range state.pairs{if mode==2{out.values=append(out.values,tsArrayValue(&tsArray{values:[]tsValue{tsStringReference(tsStringUTF8(pair[0])),tsStringReference(tsStringUTF8(pair[1]))}}))}else if mode==0{out.values=append(out.values,tsStringReference(tsStringUTF8(pair[0])))}else{out.values=append(out.values,tsStringReference(tsStringUTF8(pair[1])))}};return tsArrayValue(out)}
func tsNodeURLParamsFrom(value tsValue)[][2]string{
	if tsIsUndefined(value)||value.kind==tsNullKind{return nil};if value.kind==tsStringKind{return tsNodeURLParsePairs(tsText(value))};pairs:=[][2]string{}
	if value.kind==tsArrayKind{array:=(*tsArray)(value.ref);for i:=0;i<array.length();i++{entry:=array.at(i);if entry.kind!=tsArrayKind||(*tsArray)(entry.ref).length()!=2{tsNodeURLFail("ERR_INVALID_TUPLE","Each query pair must be a two-item sequence")};pairs=append(pairs,[2]string{tsText((*tsArray)(entry.ref).at(0)),tsText((*tsArray)(entry.ref).at(1))})};return pairs}
	if value.kind==tsObjectKind{for _,key:=range tsObjectKeys(value).values{name:=tsPropertyKey(key);item:=tsGet(value,tsStringReference(tsStringUTF8(name)));if item.kind==tsArrayKind{for _,part:=range (*tsArray)(item.ref).values{pairs=append(pairs,[2]string{name,tsText(part)})}}else{pairs=append(pairs,[2]string{name,tsText(item)})}};return pairs}
	it:=tsIterate(value);for it.next(){entry:=it.value;if entry.kind!=tsArrayKind||(*tsArray)(entry.ref).length()!=2{tsNodeURLFail("ERR_INVALID_TUPLE","Each query pair must be a two-item sequence")};pairs=append(pairs,[2]string{tsText((*tsArray)(entry.ref).at(0)),tsText((*tsArray)(entry.ref).at(1))})};return pairs
}
func tsNodeURLClassStatic(class *tsClass)*tsProperties{return tsInstanceProperties(class.static)}
func tsNodeURLBuildClasses(loop *tsLoop){
	_ = loop
	tsNodeURLSingleton.urlClass=&tsClass{builtin:"URL",static:tsInstanceValue(tsNewProperties()),nativePrototype:tsNewObject()}
	tsNodeURLSingleton.paramsClass=&tsClass{builtin:"URLSearchParams",static:tsInstanceValue(tsNewProperties()),nativePrototype:tsNewObject()}
	tsNodeURLSingleton.patternClass=&tsClass{builtin:"URLPattern",static:tsInstanceValue(tsNewProperties()),nativePrototype:tsNewObject()}
	tsNodeURLSingleton.legacyClass=&tsClass{builtin:"Url",static:tsInstanceValue(tsNewProperties()),nativePrototype:tsNewObject()}
	urlClass:=tsNodeURLSingleton.urlClass;urlClass.nativePrototype.set("constructor",tsClassValue(urlClass));urlClass.construct=func(_ *tsLoop, args ...tsValue)tsValue{input:=tsText(tsArg(args,0));var base *url.URL;if len(args)>1&&!tsIsUndefined(args[1]){b:=tsArg(args,1);if b.kind==tsObjectKind&&(*tsObject)(b.ref).nodeURL!=nil{base=(*tsObject)(b.ref).nodeURL.parsed}else{var err error;base,err=url.Parse(tsText(b));if err!=nil{tsNodeURLFail("ERR_INVALID_URL",tsText(b))}}};parsed,err:=tsNodeURLParse(input,base);if err!=nil{tsNodeURLFail("ERR_INVALID_URL",input)};return tsNodeURLObject(parsed)}
	for _,entry:=range []struct{name string;length int}{{"toString",0},{"toJSON",0},{"toLocaleString",0}}{name:=entry.name;urlClass.nativePrototype.set(name,tsNamedNativeMethod(name,entry.length,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if receiver.kind!=tsObjectKind||(*tsObject)(receiver.ref).nodeURL==nil{tsNodeURLFail("ERR_INVALID_THIS","URL method called on incompatible receiver")};return tsNodeURLField((*tsObject)(receiver.ref).nodeURL,"href")}))}
	classStatic:=tsNodeURLClassStatic(urlClass);classStatic.set("canParse",tsNamedNativeMethod("canParse",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{input:=tsText(tsArg(args,0));var base *url.URL;if len(args)>1&&!tsIsUndefined(args[1]){parsed,err:=url.Parse(tsText(args[1]));if err!=nil{return tsBooleanValue(false)};base=parsed};_,err:=tsNodeURLParse(input,base);return tsBooleanValue(err==nil)}));classStatic.set("parse",tsNamedNativeMethod("parse",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{input:=tsText(tsArg(args,0));var base *url.URL;if len(args)>1&&!tsIsUndefined(args[1]){parsed,err:=url.Parse(tsText(args[1]));if err!=nil{return tsNull};base=parsed};parsed,err:=tsNodeURLParse(input,base);if err!=nil{return tsNull};return tsNodeURLObject(parsed)}))
	paramsClass:=tsNodeURLSingleton.paramsClass;paramsClass.nativePrototype.set("constructor",tsClassValue(paramsClass));paramsClass.construct=func(_ *tsLoop, args ...tsValue)tsValue{return tsNodeURLParamsObject(tsNodeURLParamsFrom(tsArg(args,0)),nil)}
	patternClass:=tsNodeURLSingleton.patternClass;patternClass.nativePrototype.set("constructor",tsClassValue(patternClass));patternClass.construct=func(_ *tsLoop, args ...tsValue)tsValue{object:=tsNewObject();object.nativeClass=patternClass;object.prototype=tsObjectValue(patternClass.nativePrototype);object.set("pathname",tsStringReference(tsStringUTF8(tsText(tsArg(args,0)))));return tsObjectValue(object)}
	patternClass.nativePrototype.set("test",tsNamedNativeMethod("test",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{pattern:=tsText(tsGet(receiver,tsStringReference(tsStringUTF8("pathname"))));inputValue:=tsArg(args,0);input:=tsText(inputValue);if inputValue.kind==tsObjectKind&&(*tsObject)(inputValue.ref).nodeURL!=nil{input=tsText(tsNodeURLField((*tsObject)(inputValue.ref).nodeURL,"pathname"))};return tsBooleanValue(strings.Contains(input,pattern))}))
	patternClass.nativePrototype.set("exec",tsNamedNativeMethod("exec",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{test:=tsGet(tsClassValue(patternClass),tsStringReference(tsStringUTF8("test")));if !tsTruthy(tsCallReceiver(loop,test,receiver,tsArg(args,0))){return tsNull};out:=tsNewObject();out.set("pathname",tsNewObject());return tsObjectValue(out)}))
	legacyClass:=tsNodeURLSingleton.legacyClass;legacyClass.nativePrototype.set("constructor",tsClassValue(legacyClass));legacyClass.construct=func(_ *tsLoop, args ...tsValue)tsValue{object:=tsNewObject();object.nativeClass=legacyClass;object.prototype=tsObjectValue(legacyClass.nativePrototype);for _,name:=range []string{"protocol","slashes","auth","host","port","hostname","hash","search","query","pathname","path","href"}{object.set(name,tsNull)};return tsObjectValue(object)}
	legacyClass.nativePrototype.set("parse",tsNamedNativeMethod("parse",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeURLLegacyParseInto(loop,receiver,args)}));legacyClass.nativePrototype.set("format",tsNamedNativeMethod("format",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeURLLegacyFormat(loop,receiver)))}));legacyClass.nativePrototype.set("resolve",tsNamedNativeMethod("resolve",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeURLResolve(tsText(tsGet(receiver,tsStringReference(tsStringUTF8("href")))),tsText(tsArg(args,0)))))}));legacyClass.nativePrototype.set("resolveObject",tsNamedNativeMethod("resolveObject",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeURLLegacyParseValue(loop,tsNodeURLResolve(tsText(tsGet(receiver,tsStringReference(tsStringUTF8("href")))),tsText(tsArg(args,0))),false,false)}));legacyClass.nativePrototype.set("parseHost",tsNamedNativeMethod("parseHost",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsNodeURLLegacyParseHost(receiver);return tsU}))
}
func tsNodeURLLegacySet(object tsValue,name string,value tsValue){(*tsObject)(object.ref).values[name]=value}
func tsNodeURLLegacyParseValue(loop *tsLoop,text string,parseQuery,slashes bool)tsValue{value:=tsConstruct(tsClassValue(tsNodeURLSingleton.legacyClass));tsNodeURLLegacyParseInto(loop,value,[]tsValue{tsStringReference(tsStringUTF8(text)),tsBooleanValue(parseQuery),tsBooleanValue(slashes)});return value}
func tsNodeURLLegacyParseInto(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{
	if receiver.kind!=tsObjectKind{tsNodeURLFail("ERR_INVALID_THIS","Url.parse called on incompatible receiver")};text:=tsText(tsArg(args,0));if strings.TrimSpace(text)==""{tsNodeURLFail("ERR_INVALID_URL","Invalid URL")};parsed,err:=url.Parse(strings.TrimSpace(strings.ReplaceAll(text,"\\","/")));if err!=nil{tsNodeURLFail("ERR_INVALID_URL",text)};parseQuery:=len(args)>1&&tsTruthy(args[1]);slashesArg:=len(args)>2&&tsTruthy(args[2]);hasSlashes:=strings.HasPrefix(strings.TrimPrefix(text,parsed.Scheme+":"),"//")
	protocol:=tsNull;if parsed.Scheme!=""{protocol=tsStringReference(tsStringUTF8(strings.ToLower(parsed.Scheme)+":"))};slashes:=tsNull;if hasSlashes{slashes=tsBooleanValue(true)};auth:=tsNull;if parsed.User!=nil{user:=parsed.User.Username();password,ok:=parsed.User.Password();if ok{user+=":"+password};auth=tsStringReference(tsStringUTF8(user))}
	host:=tsNull;hostname:=tsNull;port:=tsNull;if parsed.Host!=""||slashesArg{host=tsStringReference(tsStringUTF8(parsed.Host));hostName:=parsed.Hostname();if strings.Contains(hostName,":"){hostName="["+hostName+"]"};hostname=tsStringReference(tsStringUTF8(strings.ToLower(hostName)));if parsed.Port()!=""{port=tsStringReference(tsStringUTF8(parsed.Port()))}}
	hash:=tsNull;if parsed.Fragment!=""||parsed.RawFragment!=""{hash=tsStringReference(tsStringUTF8("#"+parsed.EscapedFragment()))};search:=tsNull;query:=tsNull;if parsed.RawQuery!=""||parsed.ForceQuery{search=tsStringReference(tsStringUTF8("?"+parsed.RawQuery));if parseQuery{query=tsObjectValue(tsNodeURLQueryObject(parsed.RawQuery))}else{query=tsStringReference(tsStringUTF8(parsed.RawQuery))}}else if parseQuery{query=tsObjectValue(tsNewObject())}
	pathname:=tsNull;if parsed.Path!=""{pathname=tsStringReference(tsStringUTF8(parsed.EscapedPath()))};pathValue:=tsNull;if pathname.kind==tsStringKind||search.kind==tsStringKind{p:="";if pathname.kind==tsStringKind{p=tsText(pathname)};if search.kind==tsStringKind{p+=tsText(search)};pathValue=tsStringReference(tsStringUTF8(p))}
	for name,value:=range map[string]tsValue{"protocol":protocol,"slashes":slashes,"auth":auth,"host":host,"port":port,"hostname":hostname,"hash":hash,"search":search,"query":query,"pathname":pathname,"path":pathValue,"href":tsStringReference(tsStringUTF8(text))}{tsNodeURLLegacySet(receiver,name,value)}
	return receiver
}
func tsNodeURLQueryObject(query string)*tsObject{out:=tsNewObject();for _,pair:=range tsNodeURLParsePairs(query){value,exists:=out.values[pair[0]];if !exists{out.set(pair[0],tsStringReference(tsStringUTF8(pair[1])))}else if value.kind==tsArrayKind{(*tsArray)(value.ref).values=append((*tsArray)(value.ref).values,tsStringReference(tsStringUTF8(pair[1])))}else{out.set(pair[0],tsArrayValue(&tsArray{values:[]tsValue{value,tsStringReference(tsStringUTF8(pair[1]))}}))}};return out}
func tsNodeURLLegacyParseHost(receiver tsValue){object:=(*tsObject)(receiver.ref);host:=tsText(object.values["host"]);if host==""||host=="null"{return};if strings.HasPrefix(host,"["){if close:=strings.Index(host,"]");close>=0{object.values["hostname"]=tsStringReference(tsStringUTF8(host[:close+1]));if len(host)>close+1&&host[close+1]==':'{object.values["port"]=tsStringReference(tsStringUTF8(host[close+2:]))};return}};if index:=strings.LastIndex(host,":");index>=0&&strings.Index(host,":")==index{object.values["hostname"]=tsStringReference(tsStringUTF8(host[:index]));if index+1<len(host){object.values["port"]=tsStringReference(tsStringUTF8(host[index+1:]))}}else{object.values["hostname"]=tsStringReference(tsStringUTF8(host))}}
func tsNodeURLLegacyFormat(loop *tsLoop,value tsValue)string{
	get:=func(name string)tsValue{return tsGet(value,tsStringReference(tsStringUTF8(name)))};auth:=tsText(get("auth"));if auth!="null"&&auth!=""{auth=url.QueryEscape(auth);auth=strings.ReplaceAll(auth,"+","%20");auth+="@"}else{auth=""}
	protocol:=tsText(get("protocol"));if protocol=="null"{protocol=""};if protocol!=""&&!strings.HasSuffix(protocol,":"){protocol+=":"};host:=tsText(get("host"));if host=="null"||host==""{host=tsText(get("hostname"));if host=="null"{host=""};if strings.Contains(host,":")&&!strings.HasPrefix(host,"["){host="["+host+"]"};port:=tsText(get("port"));if port!="null"&&port!=""{host+=":"+port}}
	pathname:=tsText(get("pathname"));if pathname=="null"{pathname=""};pathname=strings.ReplaceAll(strings.ReplaceAll(pathname,"?","%3F"),"#","%23");search:=tsText(get("search"));query:=get("query");if (search=="null"||search=="")&&query.kind==tsObjectKind{params:=tsNodeURLParamsObject(tsNodeURLParamsFrom(query),nil);search=tsNodeURLParamsString((*tsObject)(params.ref).nodeURLParams);if search!=""{search="?"+search}};if search=="null"{search=""};if search!=""&&!strings.HasPrefix(search,"?"){search="?"+search};search=strings.ReplaceAll(search,"#","%23")
	hash:=tsText(get("hash"));if hash=="null"{hash=""};if hash!=""&&!strings.HasPrefix(hash,"#"){hash="#"+hash};slashes:=tsTruthy(get("slashes"));slashProto:=map[string]bool{"http:":true,"https:":true,"ftp:":true,"file:":true,"ws:":true,"wss:":true,"gopher:":true};if slashes||slashProto[protocol]{if slashes||host!=""{if pathname!=""&&!strings.HasPrefix(pathname,"/"){pathname="/"+pathname};host="//"+host}else if protocol=="file:"{host="//"}}
	return protocol+host+auth+pathname+search+hash
}
func tsNodeURLResolve(source,relative string)string{base,err:=url.Parse(source);if err!=nil{tsNodeURLFail("ERR_INVALID_URL",source)};ref,err:=url.Parse(relative);if err!=nil{tsNodeURLFail("ERR_INVALID_URL",relative)};return base.ResolveReference(ref).String()}

func tsNodeURLLegacyResolveObject(loop *tsLoop,source,relative string)tsValue{return tsNodeURLLegacyParseValue(loop,tsNodeURLResolve(source,relative),false,false)}
func tsNodeURLAsParsed(value tsValue)(*url.URL,bool){
	if value.kind==tsObjectKind&&(*tsObject)(value.ref).nodeURL!=nil{return (*tsObject)(value.ref).nodeURL.parsed,true}
	parsed,err:=url.Parse(tsText(value));if err!=nil{return nil,false};return parsed,true
}
func tsNodeURLAsFilePath(loop *tsLoop,value tsValue,windows bool)string{
	var path string
	if value.kind==tsStringKind{
		if !strings.HasPrefix(strings.ToLower(tsText(value)),"file:"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_URL_SCHEME","The URL must be of scheme file")})}
		parsed,err:=url.Parse(tsText(value));if err!=nil||parsed.Scheme!="file"{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_URL","Invalid file URL")})}
		if parsed.Host!=""&&parsed.Host!="localhost"&&!windows{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_HOST","File URL host must be empty or localhost")})}
		encoded:=parsed.EscapedPath();if strings.Contains(strings.ToLower(encoded),"%2f")||windows&&strings.Contains(strings.ToLower(encoded),"%5c"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_PATH","File URL path contains an encoded path separator")})}
		decoded,err:=url.PathUnescape(encoded);if err!=nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_PATH",err.Error())})};path=decoded
		if windows{path=strings.ReplaceAll(path,"/","\\");if strings.HasPrefix(path,"\\")&&len(path)>2&&path[2]==':'{path=path[1:]};if parsed.Host!=""&&parsed.Host!="localhost"{path="\\\\"+parsed.Host+path}}
	}else if value.kind==tsObjectKind&&(*tsObject)(value.ref).nodeURL!=nil{
		parsed:=(*tsObject)(value.ref).nodeURL.parsed;if parsed.Scheme!="file"{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_URL_SCHEME","The URL must be of scheme file")})};if parsed.Host!=""&&parsed.Host!="localhost"&&!windows{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_HOST","File URL host must be empty or localhost")})}
		encoded:=parsed.EscapedPath();if strings.Contains(strings.ToLower(encoded),"%2f")||windows&&strings.Contains(strings.ToLower(encoded),"%5c"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_PATH","File URL path contains an encoded path separator")})};decoded,err:=url.PathUnescape(encoded);if err!=nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_FILE_URL_PATH",err.Error())})};path=decoded
		if windows{path=strings.ReplaceAll(path,"/","\\");if strings.HasPrefix(path,"\\")&&len(path)>2&&path[2]==':'{path=path[1:]};if parsed.Host!=""&&parsed.Host!="localhost"{path="\\\\"+parsed.Host+path}}
	}else{path=tsFSPath(value);if windows{path=strings.ReplaceAll(path,"/","\\")}}
	if strings.ContainsRune(path,0){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The file URL path must not contain encoded null bytes")})};_ = loop;return path
}
func tsNodeURLPathToFileURL(path string,windows bool)tsValue{
	if windows{path=strings.ReplaceAll(path,"/","\\");path=tsPathWinResolveTexts([]string{path});if strings.HasPrefix(path,"\\\\"){trimmed:=strings.TrimLeft(path,"\\");parts:=strings.SplitN(trimmed,"\\",2);host:=parts[0];rest:="/";if len(parts)>1{rest+="/"+strings.ReplaceAll(parts[1],"\\","/")};return tsNodeURLObject(&url.URL{Scheme:"file",Host:host,Path:rest})};path=strings.ReplaceAll(path,"\\","/");if len(path)>=2&&path[1]==':'{path="/"+path};return tsNodeURLObject(&url.URL{Scheme:"file",Path:path})}
	path=tsPathPosixResolveTexts([]string{path});return tsNodeURLObject(&url.URL{Scheme:"file",Path:path})
}
func tsNodeURLFilePathBuffer(loop *tsLoop,args []tsValue)tsValue{windows:=false;if len(args)>1{windows=tsTruthy(tsGet(loop,tsArg(args,1),tsStringReference(tsStringUTF8("windows"))))};return tsNodeBuffer([]byte(tsNodeURLAsFilePath(loop,tsArg(args,0),windows)))}
func tsNodeURLFormat(loop *tsLoop,args []tsValue)tsValue{
	value:=tsArg(args,0);if value.kind==tsObjectKind&&(*tsObject)(value.ref).nodeURL!=nil{return tsNodeURLField((*tsObject)(value.ref).nodeURL,"href")}
	if value.kind==tsStringKind{parsed,err:=url.Parse(tsText(value));if err==nil{return tsStringReference(tsStringUTF8(parsed.String()))}}
	text:=tsNodeURLLegacyFormat(loop,value);options:=tsArg(args,1);if options.kind==tsObjectKind{
		if !tsTruthy(tsGet(loop,options,tsStringReference(tsStringUTF8("auth")))){if at:=strings.LastIndex(text,"@");at>=0{schemeEnd:=strings.Index(text,"://");if schemeEnd>=0&&at>schemeEnd+2{text=text[:schemeEnd+3]+text[at+1:]}}}
		if !tsTruthy(tsGet(loop,options,tsStringReference(tsStringUTF8("search")))){if at:=strings.IndexByte(text,'?');at>=0{hash:=strings.IndexByte(text[at:],'#');if hash>=0{text=text[:at]+text[at+hash:]}else{text=text[:at]}}}
		if !tsTruthy(tsGet(loop,options,tsStringReference(tsStringUTF8("fragment")))){if at:=strings.IndexByte(text,'#');at>=0{text=text[:at]}}
	};return tsStringReference(tsStringUTF8(text))
}
func tsNodeURLPatternMatch(pattern,input string)bool{
	if pattern==""||pattern=="*"||pattern=="/*"{return true};if strings.Contains(pattern,"*"){parts:=strings.Split(pattern,"*");at:=0;for _,part:=range parts{if part==""{continue};index:=strings.Index(input[at:],part);if index<0{return false};at+=index+len(part)};return true};return pattern==input
}
func tsNodeURLModule(loop *tsLoop)tsValue{
	tsNodeURLSingleton.Do(func(){tsNodeURLBuildClasses(loop);module:=tsNewObject();tsNodeURLSingleton.module=module
		urlClass:=tsClassValue(tsNodeURLSingleton.urlClass);paramsClass:=tsClassValue(tsNodeURLSingleton.paramsClass);patternClass:=tsClassValue(tsNodeURLSingleton.patternClass);legacyClass:=tsClassValue(tsNodeURLSingleton.legacyClass)
		module.set("URL",urlClass);module.set("URLSearchParams",paramsClass);module.set("URLPattern",patternClass);module.set("Url",legacyClass)
		module.set("parse",tsNamedNativeMethod("parse",3,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeURLLegacyParseValue(loop,tsText(tsArg(args,0)),len(args)>1&&tsTruthy(args[1]),len(args)>2&&tsTruthy(args[2]))}))
		module.set("resolve",tsNamedNativeMethod("resolve",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsStringReference(tsStringUTF8(tsNodeURLResolve(tsText(tsArg(args,0)),tsText(tsArg(args,1)))))}))
		module.set("resolveObject",tsNamedNativeMethod("resolveObject",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeURLLegacyResolveObject(loop,tsText(tsArg(args,0)),tsText(tsArg(args,1)))}))
		module.set("format",tsNamedNativeMethod("format",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeURLFormat(loop,args)}))
		module.set("pathToFileURL",tsNamedNativeMethod("pathToFileURL",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{windows:=false;if len(args)>1{windows=tsTruthy(tsGet(loop,tsArg(args,1),tsStringReference(tsStringUTF8("windows"))))};return tsNodeURLPathToFileURL(tsFSPath(tsArg(args,0)),windows)}))
		module.set("fileURLToPath",tsNamedNativeMethod("fileURLToPath",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{windows:=false;if len(args)>1{windows=tsTruthy(tsGet(loop,tsArg(args,1),tsStringReference(tsStringUTF8("windows"))))};return tsStringReference(tsStringUTF8(tsNodeURLAsFilePath(loop,tsArg(args,0),windows)))}))
		module.set("fileURLToPathBuffer",tsNamedNativeMethod("fileURLToPathBuffer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeURLFilePathBuffer(loop,args)}))
		internal:=tsInternalURLModule(loop);for _,name:=range []string{"domainToASCII","domainToUnicode","urlToHttpOptions"}{module.set(name,tsGet(loop,internal,tsStringReference(tsStringUTF8(name))))}
		module.set("default",tsObjectValue(module))
	});return tsObjectValue(tsNodeURLSingleton.module)
}
`

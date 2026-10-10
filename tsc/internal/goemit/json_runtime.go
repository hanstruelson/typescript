package goemit

// Native UTF-16 JSON operations. No JavaScript engine or Go interface boxing.
const jsonRuntime = `
var tsJSONNamespace struct{sync.Once;object *tsObject}
func tsJSONModule()tsValue{tsJSONNamespace.Do(func(){object:=tsNewObject();object.descriptors=map[string]*tsDescriptor{};for _,name:=range []string{"parse","stringify","rawJSON","isRawJSON"}{operation:=name;arity:=1;if name=="parse"{arity=2};if name=="stringify"{arity=3};object.set(name,tsNamedNativeMethod(name,arity,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsJSONBuiltin(loop,operation,args)}));object.descriptors[name]=&tsDescriptor{writable:true,configurable:true}};tsJSONNamespace.object=object});return tsObjectValue(tsJSONNamespace.object)}
type tsJSONNode struct {value tsValue;start,end int;children map[string]*tsJSONNode}
type tsJSONParser struct {units []uint16;index,depth int}
func tsJSONSyntax(){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"SyntaxError",message:"Invalid JSON text"})})}
func(p *tsJSONParser)white(){for p.index<len(p.units){c:=p.units[p.index];if c!=32&&c!=9&&c!=10&&c!=13{break};p.index++}}
func(p *tsJSONParser)consume(c uint16)bool{if p.index<len(p.units)&&p.units[p.index]==c{p.index++;return true};return false}
func(p *tsJSONParser)quoted()*tsString{if !p.consume('"'){tsJSONSyntax()};out:=[]uint16{};for p.index<len(p.units){c:=p.units[p.index];p.index++;if c=='"'{return tsStringUnits(out)};if c<32{tsJSONSyntax()};if c=='\\'{if p.index>=len(p.units){tsJSONSyntax()};c=p.units[p.index];p.index++;switch c{case '"','\\','/':case 'b':c=8;case 'f':c=12;case 'n':c=10;case 'r':c=13;case 't':c=9;case 'u':c=0;for i:=0;i<4;i++{if p.index>=len(p.units){tsJSONSyntax()};digit:=p.units[p.index];p.index++;var n uint16;switch{case digit>='0'&&digit<='9':n=digit-'0';case digit>='a'&&digit<='f':n=digit-'a'+10;case digit>='A'&&digit<='F':n=digit-'A'+10;default:tsJSONSyntax()};c=c*16+n};default:tsJSONSyntax()}};out=append(out,c)};tsJSONSyntax();return nil}
func(p *tsJSONParser)node()*tsJSONNode{
 p.white();p.depth++;if p.depth>10000{tsJSONSyntax()};defer func(){p.depth--}();start:=p.index;if start>=len(p.units){tsJSONSyntax()};node:=&tsJSONNode{start:start}
 switch p.units[p.index]{
 case '"':node.value=tsStringReference(p.quoted())
 case '{':p.index++;object:=tsNewObject();node.children=map[string]*tsJSONNode{};p.white();if !p.consume('}'){for{p.white();key:=tsPropertyKey(tsStringReference(p.quoted()));p.white();if !p.consume(':'){tsJSONSyntax()};child:=p.node();object.set(key,child.value);node.children[key]=child;p.white();if p.consume('}'){break};if !p.consume(','){tsJSONSyntax()}}};node.value=tsObjectValue(object)
 case '[':p.index++;array:=&tsArray{};node.children=map[string]*tsJSONNode{};p.white();if !p.consume(']'){for{child:=p.node();node.children[strconv.Itoa(len(array.values))]=child;array.values=append(array.values,child.value);p.white();if p.consume(']'){break};if !p.consume(','){tsJSONSyntax()}}};node.value=tsArrayValue(array)
 case 't','f','n':literal:="null";node.value=tsNull;if p.units[p.index]=='t'{literal="true";node.value=tsBooleanValue(true)}else if p.units[p.index]=='f'{literal="false";node.value=tsBooleanValue(false)};for _,c:=range literal{if !p.consume(uint16(c)){tsJSONSyntax()}}
 default:
  p.consume('-');if p.consume('0'){}else{if p.index>=len(p.units)||p.units[p.index]<'1'||p.units[p.index]>'9'{tsJSONSyntax()};for p.index<len(p.units)&&p.units[p.index]>='0'&&p.units[p.index]<='9'{p.index++}}
  if p.consume('.') {first:=p.index;for p.index<len(p.units)&&p.units[p.index]>='0'&&p.units[p.index]<='9'{p.index++};if first==p.index{tsJSONSyntax()}}
  if p.consume('e')||p.consume('E'){if !p.consume('+'){p.consume('-')};first:=p.index;for p.index<len(p.units)&&p.units[p.index]>='0'&&p.units[p.index]<='9'{p.index++};if first==p.index{tsJSONSyntax()}}
  node.value=tsNumberValue(tsStringNumber(tsStringUnits(p.units[start:p.index])))
 };node.end=p.index;return node
}
func tsJSONDelete(holder tsValue,key string){if holder.kind==tsObjectKind{object:=(*tsObject)(holder.ref);delete(object.values,key);delete(object.descriptors,key);for i,name:=range object.order{if name==key{object.order=append(object.order[:i],object.order[i+1:]...);break}}}else if holder.kind==tsArrayKind{array:=(*tsArray)(holder.ref);index,_:=strconv.Atoi(key);if array.holes==nil{array.holes=map[int]bool{}};array.holes[index]=true;array.values[index]=tsU}}
func tsJSONRevive(loop *tsLoop,holder tsValue,key string,node *tsJSONNode,source []uint16,reviver tsValue)tsValue{
 value:=tsGet(loop,holder,tsStringReference(tsStringKey(key)));if value.kind==tsArrayKind||value.kind==tsObjectKind{keys:=tsObjectKeys(value);if value.kind==tsArrayKind{keys=&tsArray{};length:=(*tsArray)(value.ref).length();for i:=0;i<length;i++{keys.values=append(keys.values,tsStringReference(tsStringUTF8(strconv.Itoa(i))))}};for _,name:=range keys.values{childKey:=tsPropertyKey(name);var child *tsJSONNode;if node!=nil{child=node.children[childKey]};replacement:=tsJSONRevive(loop,value,childKey,child,source,reviver);if tsIsUndefined(replacement){tsJSONDelete(value,childKey)}else{tsSet(loop,value,name,replacement)}}}
 context:=tsNewObject();if node!=nil&&tsIsPrimitiveValue(value)&&tsSameValue(value,node.value){context.set("source",tsStringReference(tsStringUnits(source[node.start:node.end])))};return tsCallReceiver(loop,reviver,holder,tsStringReference(tsStringKey(key)),value,tsObjectValue(context))
}
func tsJSONQuote(value *tsString)string{
 var out strings.Builder;out.WriteByte('"');for i:=0;i<len(value.units);i++{c:=value.units[i];switch c{case '"':out.WriteString("\\\"");case '\\':out.WriteString("\\\\");case 8:out.WriteString("\\b");case 12:out.WriteString("\\f");case 10:out.WriteString("\\n");case 13:out.WriteString("\\r");case 9:out.WriteString("\\t");default:if c<32||c>=0xd800&&c<=0xdfff{if c>=0xd800&&c<=0xdbff&&i+1<len(value.units)&&value.units[i+1]>=0xdc00&&value.units[i+1]<=0xdfff{out.WriteRune(utf16.DecodeRune(rune(c),rune(value.units[i+1])));i++}else{fmt.Fprintf(&out,"\\u%04x",c)}}else{out.WriteRune(rune(c))}}};out.WriteByte('"');return out.String()
}
type tsJSONWriter struct {replacer tsValue;keys []string;propertyList bool;gap string;stack map[unsafe.Pointer]bool;depth int}
func(w *tsJSONWriter)property(loop *tsLoop,holder tsValue,key string)(string,bool){
 value:=tsGet(loop,holder,tsStringReference(tsStringKey(key)))
 if !tsIsPrimitiveValue(value)||value.kind==tsBigIntKind{method:=tsGet(loop,value,tsStringReference(tsStringUTF8("toJSON")));if method.kind==tsFunctionKind{value=tsCallReceiver(loop,method,value,tsStringReference(tsStringKey(key)))}}
 if w.replacer.kind==tsFunctionKind{value=tsCallReceiver(loop,w.replacer,holder,tsStringReference(tsStringKey(key)),value)}
 if tsIsNumeric(value){if math.IsNaN(tsNumber(value))||math.IsInf(tsNumber(value),0){return "null",true};return tsText(value),true}
 switch value.kind{case tsNullKind:return "null",true;case tsBooleanKind:return tsText(value),true;case tsStringKind:return tsJSONQuote((*tsString)(value.ref)),true;case tsSymbolKind,tsUndefinedKind,tsFunctionKind,tsClassKind:return "",false;case tsBigIntKind:tsPropertyFailure("Do not know how to serialize a BigInt")}
 if value.kind==tsObjectKind&&(*tsObject)(value.ref).rawJSON!=nil{return tsPropertyKey(tsStringReference((*tsObject)(value.ref).rawJSON)),true}
 if w.stack[value.ref]{tsPropertyFailure("Converting circular structure to JSON")};w.stack[value.ref]=true;defer delete(w.stack,value.ref);w.depth++;defer func(){w.depth--}();if w.depth>10000{tsPropertyFailure("JSON nesting is too deep")}
 parts:=[]string{};open,close:="{","}";if value.kind==tsArrayKind{open,close="[","]";array:=(*tsArray)(value.ref);for i:=0;i<array.length();i++{part,ok:=w.property(loop,value,strconv.Itoa(i));if !ok{part="null"};parts=append(parts,part)}}else{keys:=w.keys;if !w.propertyList{keys=[]string{};for _,name:=range tsObjectKeys(value).values{keys=append(keys,tsPropertyKey(name))}};for _,name:=range keys{part,ok:=w.property(loop,value,name);if ok{separator:=":";if w.gap!=""{separator=": "};parts=append(parts,tsJSONQuote(tsStringKey(name))+separator+part)}}}
 if len(parts)==0{return open+close,true};if w.gap==""{return open+strings.Join(parts,",")+close,true};indent:=strings.Repeat(w.gap,w.depth);return open+"\n"+indent+strings.Join(parts,",\n"+indent)+"\n"+strings.Repeat(w.gap,w.depth-1)+close,true
}
func tsJSONBuiltin(loop *tsLoop,name string,args []tsValue)tsValue{
 value:=tsArg(args,0);switch name{
 case "parse","rawJSON":text:=tsStringValue(loop,value);parser:=&tsJSONParser{units:text.units};node:=parser.node();parser.white();if parser.index!=len(parser.units){tsJSONSyntax()};if name=="rawJSON"{if !tsIsPrimitiveValue(node.value)||node.value.kind==tsUndefinedKind||node.start!=0||node.end!=len(text.units){tsJSONSyntax()};out:=tsNewObject();out.rawJSON=text;out.set("rawJSON",tsStringReference(text));result:=tsObjectValue(out);tsDescriptorBuiltin(loop,"freeze",[]tsValue{result});return result};reviver:=tsArg(args,1);if reviver.kind==tsFunctionKind{holder:=tsNewObject();holder.set("",node.value);return tsJSONRevive(loop,tsObjectValue(holder),"",node,text.units,reviver)};return node.value
 case "isRawJSON":return tsBooleanValue(value.kind==tsObjectKind&&(*tsObject)(value.ref).rawJSON!=nil)
 case "stringify":writer:=&tsJSONWriter{replacer:tsArg(args,1),stack:map[unsafe.Pointer]bool{}};if writer.replacer.kind==tsArrayKind{writer.propertyList=true;seen:=map[string]bool{};array:=(*tsArray)(writer.replacer.ref);for i:=0;i<array.length();i++{item:=array.at(i);if item.kind==tsStringKind||tsIsNumeric(item){key:=tsPropertyKey(item);if !seen[key]{seen[key]=true;writer.keys=append(writer.keys,key)}}}}
  space:=tsArg(args,2);if tsIsNumeric(space){count:=tsNumber(space);if count>10{count=10};if count>0{writer.gap=strings.Repeat(" ",int(count))}}else if space.kind==tsStringKind{units:=(*tsString)(space.ref).units;if len(units)>10{units=units[:10]};writer.gap=tsPropertyKey(tsStringReference(tsStringUnits(units)))};holder:=tsNewObject();holder.set("",value);text,ok:=writer.property(loop,tsObjectValue(holder),"");if !ok{return tsU};return tsStringReference(tsStringKey(text))
 };panic("Unknown JSON operation")
}
`

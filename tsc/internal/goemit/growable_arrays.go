package goemit

import (
	"fmt"
	"github.com/microsoft/TypeScript/tsc/internal/ast"
	"strings"
)

func (e *emitter) arrayElementPrimitive(node *ast.Node) primitive {
	if resolver, ok := e.resolver.(interface {
		GetEmitArrayElementType(*ast.Node, []*ast.Node) (string, uint8)
	}); ok {
		kind, nulls := resolver.GetEmitArrayElementType(node, e.specializationCalls)
		return primitive{kind, nulls}
	}
	return primitive{}
}
func (b *machineBuilder) arrayBoundary(value string, node *ast.Node) string {
	p := b.e.arrayElementPrimitive(node)
	if p.kind == "" {
		return value
	}
	return fmt.Sprintf("tsArrayBoundary(%s,%q,%d,%t)", value, p.kind, p.nulls, b.e.coerce)
}

func (b *machineBuilder) denseStorage(node *ast.Node, p primitive) string {
	if node.Kind == ast.KindIdentifier {
		if cell := b.e.binding(node); cell != nil && cell.arrayElement.kind == p.kind && cell.arrayElement.nulls == p.nulls {
			return b.typedTemp(cell.name+".readStorage()", "*tsGrowableStorage["+p.goType()+"]")
		}
	}
	receiver := b.expression(node)
	receiver = b.temp(fmt.Sprintf("tsArrayBoundary(%s,%q,%d,%t)", receiver, p.kind, p.nulls, b.e.coerce))
	return b.typedTemp("tsDenseStorage["+p.goType()+"]("+receiver+","+fmt.Sprintf("%q,%d", p.kind, p.nulls)+")", "*tsGrowableStorage["+p.goType()+"]")
}

func (b *machineBuilder) denseRead(storage, key string, p primitive) string {
	if p.nulls != 0 {
		return "tsNullableRead[" + p.goType() + "](" + storage + ",float64(" + key + "))"
	}
	kind := b.tempType(key)
	if strings.HasPrefix(kind, "int") || strings.HasPrefix(kind, "uint") {
		return "tsDenseReadInteger[" + p.goType() + "," + kind + "](" + storage + "," + key + ")"
	}
	return "tsDenseRead[" + p.goType() + "](" + storage + ",float64(" + key + "))"
}
func (b *machineBuilder) denseWrite(storage, key, value string, p primitive) string {
	if p.nulls != 0 {
		return "tsNullableWrite[" + p.goType() + "](" + storage + ",float64(" + key + ")," + value + ")"
	}
	kind := b.tempType(key)
	if strings.HasPrefix(kind, "int") || strings.HasPrefix(kind, "uint") {
		return "tsDenseWriteInteger[" + p.goType() + "," + kind + "](" + storage + "," + key + "," + value + ")"
	}
	return "tsDenseWrite[" + p.goType() + "](" + storage + ",float64(" + key + ")," + value + ")"
}

const GrowableArrayRuntime = growableArrayBaseRuntime + NullableArrayRuntime + NativeCallbackRuntime

const growableArrayBaseRuntime = `
type tsGrowableStorage[T tsPrimitive]struct {values []T;tags []uint8;owner *tsArray}
func tsArrayBinding(initialized,constant bool,kind string,nulls uint8,coerce bool)*tsCell{return &tsCell{value:tsU,initialized:initialized,constant:constant,arrayKind:kind,arrayNulls:nulls,arrayCoerce:coerce}}
func(array *tsArray)length()int{if array.growLength!=nil{return array.growLength()};return len(array.values)}
func(array *tsArray)at(index int)tsValue{if index<0||index>=array.length()||array.holes[index]{return tsU};if array.growGet!=nil{return array.growGet(index)};return array.values[index]}
func(array *tsArray)clearSlot(index int){if array.growClear!=nil{array.growClear(index)}else{array.values[index]=tsU}}
func(array *tsArray)resize(length int){if array.growResize!=nil{array.growResize(length)}else{if length<len(array.values){clear(array.values[length:]);array.values=array.values[:length]}else{array.values=append(array.values,make([]tsValue,length-len(array.values))...)}}}
func(array *tsArray)put(loop *tsLoop,index int,item tsValue){if array.growPut!=nil{array.growPut(loop,index,item)}else{if array.elementKind!=""{item=tsBoundary(item,array.elementKind,array.elementNulls,array.elementCoerce)};array.values[index]=item};delete(array.holes,index)}
func(array *tsArray)appendItems(loop *tsLoop,items ...tsValue){for _,item:=range items{length:=array.length();if array.growSplice!=nil{array.growSplice(loop,length,0,[]tsValue{item})}else{if array.elementKind!=""{item=tsBoundary(item,array.elementKind,array.elementNulls,array.elementCoerce)};array.values=append(array.values,item)}}}
func tsInstallGrowable[T tsPrimitive](array *tsArray,kind string,nulls uint8,coerce bool){storage:=&tsGrowableStorage[T]{values:make([]T,len(array.values)),owner:array};if nulls!=0{storage.tags=make([]uint8,len(array.values))};for i,value:=range array.values{if array.holes[i]{tsArrayRangeFailure("Primitive arrays cannot contain holes")};checked:=tsOptionalFrom[T](tsBoundary(value,kind,nulls,coerce));storage.values[i]=checked.value;if nulls!=0{storage.tags[i]=checked.tag}};array.native=unsafe.Pointer(storage);array.elementKind=kind;array.elementNulls=nulls;array.elementCoerce=coerce;array.growLength=func()int{return len(storage.values)};array.growClear=func(i int){var zero T;storage.values[i]=zero;if nulls!=0{storage.tags[i]=0}};array.growGet=func(i int)tsValue{if nulls!=0{return tsOptional[T]{value:storage.values[i],tag:storage.tags[i]}.raw()};return tsPrimitiveValue(storage.values[i])};array.growPut=func(loop *tsLoop,i int,value tsValue){checked:=tsOptionalFrom[T](tsBoundary(value,kind,nulls,coerce));storage.values[i]=checked.value;if nulls!=0{storage.tags[i]=checked.tag}};array.growResize=func(length int){if length<len(storage.values){clear(storage.values[length:]);storage.values=storage.values[:length];if nulls!=0{clear(storage.tags[length:]);storage.tags=storage.tags[:length]}}else{storage.values=append(storage.values,make([]T,length-len(storage.values))...);if nulls!=0{storage.tags=append(storage.tags,make([]uint8,length-len(storage.tags))...)}}};array.growSplice=func(loop *tsLoop,start,remove int,items []tsValue){if nulls!=0{converted:=make([]tsOptional[T],len(items));for i,value:=range items{converted[i]=tsOptionalFrom[T](tsBoundary(value,kind,nulls,coerce))};tsNullableSplice(storage,kind,coerce,float64(start),float64(remove),converted...);return};converted:=make([]T,len(items));for i,value:=range items{converted[i]=tsNative[T](tsBoundary(value,kind,0,coerce))};tsDenseSplice(storage,kind,coerce,float64(start),float64(remove),converted...)};array.values=nil}

func tsArrayBoundary(value tsValue,kind string,nulls uint8,coerce bool)tsValue{if value.kind!=tsArrayKind{tsConversionTypeFailure(kind+"[]",tsValueTypeName(value))};array:=(*tsArray)(value.ref);if array.elementKind!=""{if array.elementKind!=kind||array.elementNulls!=nulls{tsConversionTypeFailure(kind+"[]",array.elementKind+"[]")};return value};if len(array.holes)!=0{tsArrayRangeFailure("Primitive arrays cannot contain holes")};switch kind{case "number","float64":tsInstallGrowable[float64](array,kind,nulls,coerce);case "float32":tsInstallGrowable[float32](array,kind,nulls,coerce);case "int":tsInstallGrowable[int](array,kind,nulls,coerce);case "int8":tsInstallGrowable[int8](array,kind,nulls,coerce);case "int16":tsInstallGrowable[int16](array,kind,nulls,coerce);case "int32":tsInstallGrowable[int32](array,kind,nulls,coerce);case "int64":tsInstallGrowable[int64](array,kind,nulls,coerce);case "uint":tsInstallGrowable[uint](array,kind,nulls,coerce);case "uint8":tsInstallGrowable[uint8](array,kind,nulls,coerce);case "uint16":tsInstallGrowable[uint16](array,kind,nulls,coerce);case "uint32":tsInstallGrowable[uint32](array,kind,nulls,coerce);case "uint64":tsInstallGrowable[uint64](array,kind,nulls,coerce);case "string":tsInstallGrowable[*tsString](array,kind,nulls,coerce);case "boolean":tsInstallGrowable[bool](array,kind,nulls,coerce)};return value}
func tsNewGrowableArray(kind string,nulls uint8,coerce bool)*tsArray{value:=tsArrayBoundary(tsArrayValue(&tsArray{}),kind,nulls,coerce);return (*tsArray)(value.ref)}

func tsGrowableIndex(key tsValue)(int,bool){if !tsIsNumeric(key){return 0,false};n:=tsNumber(key);if n<0||n>=4294967295||n>float64(^uint(0)>>1)||math.Trunc(n)!=n||math.IsNaN(n){return 0,false};return int(n),true}
func tsGrowableArrayRead[T tsPrimitive](value,key tsValue,kind string)tsValue{if value.kind==tsArrayKind{array:=(*tsArray)(value.ref);if array.native!=nil&&array.elementKind==kind&&array.elementNulls==0{if index,ok:=tsGrowableIndex(key);ok{storage:=(*tsGrowableStorage[T])(array.native);if index>=len(storage.values)||array.holes[index]{return tsU};return tsPrimitiveValue(storage.values[index])}}};return tsGet(value,key)}
func tsGrowableArrayWrite[T tsPrimitive](value,key,item tsValue,kind string)tsValue{if value.kind==tsArrayKind{array:=(*tsArray)(value.ref);if array.native!=nil&&array.elementKind==kind&&array.elementNulls==0{if index,ok:=tsGrowableIndex(key);ok{storage:=(*tsGrowableStorage[T])(array.native);if index<=len(storage.values){checked:=tsBoundary(item,kind,0,array.elementCoerce);native:=tsNative[T](checked);if index==len(storage.values){storage.values=append(storage.values,native)}else{storage.values[index]=native};delete(array.holes,index);return item}}}};return tsSet(value,key,item)}
func tsGrowableArrayPush[T tsPrimitive](value tsValue,kind string,args ...tsValue)float64{if value.kind==tsArrayKind{array:=(*tsArray)(value.ref);_,overridden:=array.properties["push"];if !overridden&&array.native!=nil&&array.elementKind==kind&&array.elementNulls==0{storage:=(*tsGrowableStorage[T])(array.native);for _,item:=range args{checked:=tsBoundary(item,kind,0,array.elementCoerce);storage.values=append(storage.values,tsNative[T](checked))};return float64(len(storage.values))}};return tsNumber(tsCall(tsGet(value,"push"),args...))}
func tsGrowableArrayLength[T tsPrimitive](value tsValue,kind string)float64{if value.kind==tsArrayKind{array:=(*tsArray)(value.ref);if array.native!=nil&&array.elementKind==kind&&array.elementNulls==0{return float64(len((*tsGrowableStorage[T])(array.native).values))}};return tsNumber(tsGet(value,"length"))}

func tsDenseStorage[T tsPrimitive](value tsValue,kind string,nulls ...uint8)*tsGrowableStorage[T]{mask:=uint8(0);if len(nulls)>0{mask=nulls[0]};if value.kind!=tsArrayKind{tsArrayTypeFailure("Expected a dense "+kind+" array")};array:=(*tsArray)(value.ref);if array.native==nil||array.elementKind!=kind||array.elementNulls!=mask{tsArrayTypeFailure("Unexpected primitive array storage")};return (*tsGrowableStorage[T])(array.native)}
func tsDensePush[T tsPrimitive](storage *tsGrowableStorage[T],items ...T)float64{storage.values=append(storage.values,items...);return float64(len(storage.values))}
func tsDenseRead[T tsPrimitive](storage *tsGrowableStorage[T],index float64)T{if math.IsNaN(index)||index<0||index>=float64(len(storage.values))||math.Trunc(index)!=index{tsArrayRangeFailure("Dense array index out of bounds")};return storage.values[int(index)]}
func tsDenseWrite[T tsPrimitive](storage *tsGrowableStorage[T],index float64,value T)T{if math.IsNaN(index)||index<0||index>float64(len(storage.values))||math.Trunc(index)!=index{tsArrayRangeFailure("Dense array writes cannot create holes")};if int(index)==len(storage.values){storage.values=append(storage.values,value)}else{storage.values[int(index)]=value};return value}

func tsDensePresent[T tsPrimitive](value tsOptional[T])T{if value.tag!=0{tsArrayTypeFailure("Expected an array element")};return value.value}
func tsDensePop[T tsPrimitive](storage *tsGrowableStorage[T],front bool)tsOptional[T]{length:=len(storage.values);if length==0{return tsOptional[T]{tag:2}};index:=length-1;if front{index=0};value:=storage.values[index];if front{copy(storage.values,storage.values[1:])};var zero T;storage.values[length-1]=zero;storage.values=storage.values[:length-1];return tsOptional[T]{value:value}}
func tsDenseUnshift[T tsPrimitive](storage *tsGrowableStorage[T],items ...T)float64{old:=len(storage.values);storage.values=append(storage.values,items...);copy(storage.values[len(items):],storage.values[:old]);copy(storage.values,items);return float64(len(storage.values))}

type tsDenseArrayCell[T tsPrimitive]struct{array *tsArray;storage *tsGrowableStorage[T];initialized,constant bool;kind string;nulls uint8;coerce bool}
func tsDenseArrayBinding[T tsPrimitive](initialized,constant bool,kind string,coerce bool,nulls ...uint8)*tsDenseArrayCell[T]{mask:=uint8(0);if len(nulls)>0{mask=nulls[0]};return &tsDenseArrayCell[T]{initialized:initialized,constant:constant,kind:kind,nulls:mask,coerce:coerce}}
func(cell *tsDenseArrayCell[T])raw()tsValue{if cell.array==nil{return tsU};return tsArrayValue(cell.array)}
func(cell *tsDenseArrayCell[T])initializedState()bool{return cell.initialized}
func(cell *tsDenseArrayCell[T])bindingRef()tsBindingCell{return tsBindingCell{raw:cell.raw,initializedState:cell.initializedState}}
func(cell *tsDenseArrayCell[T])get()tsValue{if !cell.initialized{panic("Cannot access binding before initialization")};return cell.raw()}
func(cell *tsDenseArrayCell[T])init(value tsValue)tsValue{value=tsArrayBoundary(value,cell.kind,cell.nulls,cell.coerce);cell.array=(*tsArray)(value.ref);cell.storage=(*tsGrowableStorage[T])(cell.array.native);cell.initialized=true;return value}
func(cell *tsDenseArrayCell[T])set(value tsValue)tsValue{if !cell.initialized{panic("Cannot access binding before initialization")};if cell.constant{panic("Assignment to constant variable")};return cell.init(value)}
func(cell *tsDenseArrayCell[T])readStorage()*tsGrowableStorage[T]{if !cell.initialized{panic("Cannot access binding before initialization")};if cell.storage==nil{tsArrayTypeFailure("Expected an initialized primitive array")};return cell.storage}

func tsDenseArrayWrapper[T tsPrimitive](storage *tsGrowableStorage[T])*tsArray{return storage.owner}
func tsDenseBound(index float64,length int)int{index=math.Trunc(index);if math.IsNaN(index){index=0};if index<0{index+=float64(length)};if index<0{return 0};if index>float64(length){return length};return int(index)}
func tsDenseSlice[T tsPrimitive](storage *tsGrowableStorage[T],kind string,coerce bool,start,end float64)*tsArray{first,last:=tsDenseBound(start,len(storage.values)),tsDenseBound(end,len(storage.values));if last<first{last=first};result:=tsNewGrowableArray(kind,storage.owner.elementNulls,coerce);out:=(*tsGrowableStorage[T])(result.native);out.values=append(out.values,storage.values[first:last]...);if storage.owner.elementNulls!=0{out.tags=append(out.tags,storage.tags[first:last]...)};return result}
func tsDenseSplice[T tsPrimitive](storage *tsGrowableStorage[T],kind string,coerce bool,start,remove float64,items ...T)*tsArray{old:=len(storage.values);first:=tsDenseBound(start,old);remove=math.Trunc(remove);if math.IsNaN(remove)||remove<0{remove=0};count:=old-first;if remove<float64(count){count=int(remove)};result:=tsDenseSlice(storage,kind,coerce,float64(first),float64(first+count));length:=old-count+len(items);if length>old{storage.values=append(storage.values,make([]T,length-old)...)};copy(storage.values[first+len(items):],storage.values[first+count:old]);if length<old{clear(storage.values[length:]);storage.values=storage.values[:length]};copy(storage.values[first:first+len(items)],items);return result}
func tsDenseFill[T tsPrimitive](storage *tsGrowableStorage[T],value T,start,end float64){first,last:=tsDenseBound(start,len(storage.values)),tsDenseBound(end,len(storage.values));for i:=first;i<last;i++{storage.values[i]=value}}
func tsDenseReverse[T tsPrimitive](storage *tsGrowableStorage[T]){for i,j:=0,len(storage.values)-1;i<j;i,j=i+1,j-1{storage.values[i],storage.values[j]=storage.values[j],storage.values[i];if storage.owner.elementNulls!=0{storage.tags[i],storage.tags[j]=storage.tags[j],storage.tags[i]}}}
func tsDenseCopyWithin[T tsPrimitive](storage *tsGrowableStorage[T],target,start,end float64){length:=len(storage.values);dest,first,last:=tsDenseBound(target,length),tsDenseBound(start,length),tsDenseBound(end,length);count:=min(last-first,length-dest);if count>0{copy(storage.values[dest:dest+count],storage.values[first:first+count]);if storage.owner.elementNulls!=0{copy(storage.tags[dest:dest+count],storage.tags[first:first+count])}}}

func tsDenseFromValues[T tsPrimitive](kind string,coerce bool,items ...T)*tsArray{result:=tsNewGrowableArray(kind,0,coerce);storage:=(*tsGrowableStorage[T])(result.native);storage.values=append(storage.values,items...);return result}
func tsDenseNewDynamic[T tsPrimitive](kind string,coerce bool,item tsValue,zero T)*tsArray{if tsIsNumeric(item){return tsDenseNewSized(kind,coerce,tsNumber(item),zero)};return tsDenseFromValues(kind,coerce,tsNative[T](tsBoundary(item,kind,0,coerce)))}
func tsDenseNewSized[T tsPrimitive](kind string,coerce bool,length float64,zero T)*tsArray{if math.IsNaN(length)||length<0||length>4294967295||length>float64(^uint(0)>>1)||math.Trunc(length)!=length{tsArrayRangeFailure("Invalid dense array length")};result:=tsNewGrowableArray(kind,0,coerce);storage:=(*tsGrowableStorage[T])(result.native);storage.values=make([]T,int(length));var naturalZero T;if zero!=naturalZero{for i:=range storage.values{storage.values[i]=zero}};return result}
func(cell *tsDenseArrayCell[T])initStorage(storage *tsGrowableStorage[T])*tsGrowableStorage[T]{cell.storage=storage;cell.array=storage.owner;cell.initialized=true;return storage}

func tsDenseReadInteger[T tsPrimitive,I tsNativeInteger](storage *tsGrowableStorage[T],index I)T{if index<0||uint64(index)>=uint64(len(storage.values)){tsArrayRangeFailure("Dense array index out of bounds")};return storage.values[index]}
func tsDenseWriteInteger[T tsPrimitive,I tsNativeInteger](storage *tsGrowableStorage[T],index I,value T)T{if index<0||uint64(index)>uint64(len(storage.values)){tsArrayRangeFailure("Dense array writes cannot create holes")};if uint64(index)==uint64(len(storage.values)){storage.values=append(storage.values,value)}else{storage.values[index]=value};return value}
func tsCloneDense[T tsPrimitive](cell *tsDenseArrayCell[T])*tsDenseArrayCell[T]{copy:=*cell;return &copy}
`

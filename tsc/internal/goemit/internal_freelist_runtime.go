package goemit

// internalFreeListRuntime adapts Bun's tiny reusable-object free-list class.
const internalFreeListRuntime = `
func tsInternalFreeListConstructor(loop *tsLoop)tsValue{
	fn:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor FreeList cannot be invoked without new");return tsU})
	fn.name="FreeList";fn.length=3;fn.constructible=true
	fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{
		name,max,ctor:=tsArg(args,0),tsArg(args,1),tsArg(args,2);object:=tsNewObject()
		object.set("name",name);object.set("max",max);object.set("ctor",ctor);object.set("list",tsArrayValue(&tsArray{}))
		object.set("alloc",tsNamedNativeMethod("alloc",0,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{
			list:=tsInternalHTTPGet(loop,receiver,"list");if list.kind==tsArrayKind&&len((*tsArray)(list.ref).values)>0{array:=(*tsArray)(list.ref);last:=len(array.values)-1;value:=array.values[last];array.values=array.values[:last];return value}
			constructor:=tsInternalHTTPGet(loop,receiver,"ctor");return tsConstruct(loop,constructor,values...)
		}))
		object.set("free",tsNamedNativeMethod("free",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{
			list:=tsInternalHTTPGet(loop,receiver,"list");max:=tsNumber(tsInternalHTTPGet(loop,receiver,"max"));if list.kind==tsArrayKind&&float64(len((*tsArray)(list.ref).values))<max{array:=(*tsArray)(list.ref);array.values=append(array.values,tsArg(values,0));return tsBooleanValue(true)};return tsBooleanValue(false)
		}))
		return tsObjectValue(object)
	}
	fn.properties=tsNewObject();fn.properties.set("prototype",tsObjectValue(tsNewObject()))
	return tsFunctionValue(fn)
}
`

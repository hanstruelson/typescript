package goemit

const deleteRuntime = `
func tsDelete(loop *tsLoop,receiver,key tsValue,strict bool)bool{
 if tsNullish(receiver){tsPropertyFailure("Cannot delete property of null or undefined")};name:=tsPropertyKey(key);ok:=true
 switch receiver.kind{
 case tsErrorKind:return tsDelete(loop,tsObjectValue(tsErrorObject(loop,(*tsRuntimeError)(receiver.ref))),key,strict)
 case tsObjectKind:object:=(*tsObject)(receiver.ref);if descriptor:=object.descriptors[name];descriptor!=nil&&!descriptor.configurable{ok=false;break};delete(object.argumentMap,name);delete(object.values,name);delete(object.descriptors,name);for i,k:=range object.order{if k==name{object.order=append(object.order[:i],object.order[i+1:]...);break}}
 case tsFunctionKind:object:=(*tsFunction)(receiver.ref).metadata();return tsDelete(loop,tsObjectValue(object),key,strict)
 case tsInstanceKind:properties:=(*tsProperties)(receiver.ref);if descriptor:=properties.accessors[name];descriptor!=nil&&!descriptor.configurable{ok=false;break};delete(properties.accessors,name);if property,exists:=properties.declared[name];exists&&properties.own[name]{if property.clearField!=nil{property.clearField()}};delete(properties.own,name);delete(properties.extra,name);delete(properties.methods,name);for i,k:=range properties.order{if k==name{properties.order=append(properties.order[:i],properties.order[i+1:]...);break}}
 case tsClassKind:return tsDelete(loop,(*tsClass)(receiver.ref).static,key,strict)
 case tsArrayKind:array:=(*tsArray)(receiver.ref);if descriptor:=array.descriptors[name];descriptor!=nil&&!descriptor.configurable{ok=false;break};if name=="length"{ok=false;break};if index,valid:=tsGrowableIndex(key);valid&&index<array.length(){if array.native!=nil{tsPropertyFailure("Sparse deletion of native arrays is not supported")};if array.holes==nil{array.holes=map[int]bool{}};array.holes[index]=true;array.values[index]=tsU}else{delete(array.properties,name);delete(array.descriptors,name);for i,key:=range array.order{if key==name{array.order=append(array.order[:i],array.order[i+1:]...);break}}}
 case tsStringKind:if name=="length"{ok=false}else if index,valid:=tsGrowableIndex(key);valid&&index<len((*tsString)(receiver.ref).units){ok=false}
 case tsTypedArrayKind:array:=(*tsTypedArray)(receiver.ref);if index,valid:=tsGrowableIndex(key);valid&&index<array.length{ok=false}else{delete(array.properties,name)}
 }
 if !ok&&strict{tsPropertyFailure("Cannot delete non-configurable property")};return ok
}
`

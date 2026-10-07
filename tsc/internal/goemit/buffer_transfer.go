package goemit

// Weak view registration allows transfer to invalidate existing views without
// adding a detached-buffer branch to every native element load and store.
const bufferTransferRuntime = `
func tsArrayBufferRegister(buffer *tsArrayBuffer,view *tsTypedArray){
    if buffer.detached{tsArrayTypeFailure("Cannot create a view of a detached ArrayBuffer")}
    if len(buffer.views)>=64&&len(buffer.views)%64==0{live:=buffer.views[:0];for _,reference:=range buffer.views{if reference.Value()!=nil{live=append(live,reference)}};buffer.views=live}
    buffer.views=append(buffer.views,weak.Make(view))
}
func tsArrayBufferTransfer(buffer *tsArrayBuffer,length int)*tsArrayBuffer{
    if buffer.detached{tsArrayTypeFailure("Cannot transfer a detached ArrayBuffer")};data:=buffer.data
    if length!=len(data){resized:=make([]byte,length);copy(resized,data);data=resized}
    output:=&tsArrayBuffer{data:data}
    for _,reference:=range buffer.views{var view *tsTypedArray;view=reference.Value();if view==nil{continue};view.length=0;view.offset=0;view.storage=nil;view.read=func(i int)float64{tsArrayTypeFailure("Cannot read a detached TypedArray");return 0};view.write=func(i int,n float64){tsArrayTypeFailure("Cannot write a detached TypedArray")};view.view=func(start,end int)*tsTypedArray{tsArrayTypeFailure("Cannot slice a detached TypedArray");return nil}}
    buffer.views=nil;buffer.data=nil;buffer.detached=true;return output
}
func tsBYOBTransfer(value tsValue)tsValue{
    if value.kind!=tsTypedArrayKind{tsFSInvalid("view","a TypedArray")};array:=(*tsTypedArray)(value.ref)
    if array.buffer.detached||array.length==0{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_STATE","BYOB view is empty or detached")})}
    offset,length,name,bufferMarker:=array.offset,array.length,array.name,array.nodeBuffer
    buffer:=tsArrayBufferTransfer(array.buffer,len(array.buffer.data));output:=tsTypedArrayBuffer(name,[]tsValue{tsArrayBufferValue(buffer),tsNumberValue(float64(offset)),tsNumberValue(float64(length))});(*tsTypedArray)(output.ref).nodeBuffer=bufferMarker;return output
}
func tsByteView(data []byte)tsValue{value:=tsNodeBuffer(data);(*tsTypedArray)(value.ref).nodeBuffer=false;return value}
`

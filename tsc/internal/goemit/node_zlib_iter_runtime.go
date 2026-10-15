package goemit

// nodeZlibIterRuntime exports Bun's experimental compression iterable API
// backed by the local standard-library stream/iter transform adapters.
const nodeZlibIterRuntime = `
var tsNodeZlibIter struct{sync.Once;module tsValue}
func tsNodeZlibIterKey(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeZlibIterModule(loop *tsLoop)tsValue{tsNodeZlibIter.Do(func(){process:=tsNodeProcessModule(loop);if emit:=tsGet(process,tsNodeZlibIterKey("emitWarning"));emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,process,tsNodeZlibIterKey("zlib/iter is an experimental feature and might change at any time"),tsNodeZlibIterKey("ExperimentalWarning"))};source:=tsInternalStreamIterTransformModule(loop);module:=tsNewObject();names:=[]string{"compressGzip","compressGzipSync","compressDeflate","compressDeflateSync","compressBrotli","compressBrotliSync","compressZstd","compressZstdSync","decompressGzip","decompressGzipSync","decompressDeflate","decompressDeflateSync","decompressBrotli","decompressBrotliSync","decompressZstd","decompressZstdSync"};for _,name:=range names{module.set(name,tsGet(source,tsNodeZlibIterKey(name)))};module.set("default",tsObjectValue(module));tsNodeZlibIter.module=tsObjectValue(module)});return tsNodeZlibIter.module}
func tsNodeZlibIterAddon(loop *tsLoop,name string)(tsValue,bool){if name=="zlib/iter"{return tsNodeZlibIterModule(loop),true};return tsU,false}
`

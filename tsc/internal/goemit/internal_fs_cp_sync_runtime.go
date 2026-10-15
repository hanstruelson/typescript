package goemit

// internalFSCPCopySyncRuntime exposes Node's shared fs.cp validation and error
// helpers over the runtime's recursive Go filesystem copy implementation.
const internalFSCPCopySyncRuntime = `
var tsInternalFSCPCopySyncSingleton struct { sync.Once; value tsValue }

func tsInternalFSCPCopySyncError(loop *tsLoop, context tsValue, code, prefix string) tsValue {
	syscallName := tsText(tsFSProperty(context, "syscall"))
	ctxMessage := tsText(tsFSProperty(context, "message"))
	message := prefix + ": " + syscallName + " returned " + code + " (" + ctxMessage + ")"
	path, dest := tsFSProperty(context, "path"), tsFSProperty(context, "dest")
	if !tsIsUndefined(path) { message += " " + tsText(path) }
	if !tsIsUndefined(dest) { message += " => " + tsText(dest) }
	out := tsFSException("SystemError", code, message)
	if out.kind == tsErrorKind {
		err := (*tsRuntimeError)(out.ref)
		for name, value := range map[string]tsValue{
			"path": path, "dest": dest, "syscall": tsStringReference(tsStringUTF8(syscallName)),
			"errno": tsFSProperty(context, "errno"), "info": context,
		} {
			if !tsIsUndefined(value) { err.fields[name] = value }
		}
	}
	return out
}

func tsInternalFSCPCopyError(loop *tsLoop, args []tsValue, code, prefix string) tsValue {
	context := tsArg(args, 0)
	return tsInternalFSCPCopySyncError(loop, context, code, prefix)
}

func tsInternalFSCPCopySyncModule(loop *tsLoop) tsValue {
	tsInternalFSCPCopySyncSingleton.Do(func() {
		module := tsNewObject()
		module.set("validateCpOptions", tsNamedNativeMethod("validateCpOptions", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue {
			options := tsArg(args, 0)
			if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
			if options.kind != tsObjectKind { tsFSInvalid("options", "an object") }
			input := (*tsObject)(options.ref)
			out := tsNewObject()
			for _, item := range []struct{name string; value tsValue}{
				{"dereference",tsBooleanValue(false)}, {"errorOnExist",tsBooleanValue(false)}, {"filter",tsU},
				{"force",tsBooleanValue(true)}, {"preserveTimestamps",tsBooleanValue(false)},
				{"recursive",tsBooleanValue(false)}, {"verbatimSymlinks",tsBooleanValue(false)}, {"mode",tsNumberValue(0)},
			} { out.set(item.name,item.value) }
			for _, key := range input.order { out.set(key,input.get(key)) }
			for _, key := range []string{"dereference","errorOnExist","force","preserveTimestamps","recursive","verbatimSymlinks"} {
				value := out.get(key); if value.kind != tsBooleanKind { tsFSInvalid("options."+key,"a boolean") }
			}
			mode := out.get("mode")
			if !tsIsNull(mode) && !tsIsUndefined(mode) {
				value := tsFSInteger(mode,"mode",0,0); if value > 7 { panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE","mode is out of range" )}) }; out.set("mode",tsNumberValue(float64(value)))
			}
			filter := out.get("filter"); if !tsIsUndefined(filter) && filter.kind != tsFunctionKind { tsFSInvalid("options.filter","a function") }
			if tsTruthy(out.get("dereference")) && tsTruthy(out.get("verbatimSymlinks")) { panic(tsThrown{tsFSException("TypeError","ERR_INCOMPATIBLE_OPTION_PAIR","Option \"dereference\" cannot be used in combination with option \"verbatimSymlinks\"")}) }
			return tsObjectValue(out)
		}))
		module.set("tryNativeFastPathSync", tsNamedNativeMethod("tryNativeFastPathSync", 3, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue {
			out := tsNewObject(); out.set("ok",tsBooleanValue(false)); out.set("checked",tsU); return tsObjectValue(out)
		}))
		module.set("cpSyncFn", tsNamedNativeMethod("cpSyncFn", 3, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue {
			source, destination := tsFSPath(tsArg(args,0)), tsFSPath(tsArg(args,1))
			if err := tsFSCopy(loop,source,destination,tsArg(args,2),false); err != nil { panic(tsThrown{tsFSError(err,"cp",source,destination)}) }
			return tsU
		}))
		for _, item := range []struct{name,code,prefix string}{
			{"fsCpDirToNonDirError","ERR_FS_CP_DIR_TO_NON_DIR","Cannot overwrite non-directory with directory"},
			{"fsCpEExistError","ERR_FS_CP_EEXIST","Target already exists"},
			{"fsCpEinvalError","ERR_FS_CP_EINVAL","Invalid src or dest"},
			{"fsCpFifoPipeError","ERR_FS_CP_FIFO_PIPE","Cannot copy a FIFO pipe"},
			{"fsCpNonDirToDirError","ERR_FS_CP_NON_DIR_TO_DIR","Cannot overwrite directory with non-directory"},
			{"fsCpSocketError","ERR_FS_CP_SOCKET","Cannot copy a socket file"},
			{"fsCpSymlinkToSubdirectoryError","ERR_FS_CP_SYMLINK_TO_SUBDIRECTORY","Cannot overwrite symlink in subdirectory of self"},
			{"fsCpUnknownError","ERR_FS_CP_UNKNOWN","Cannot copy an unknown file type"},
			{"fsEisdirError","ERR_FS_EISDIR","Path is a directory"},
		} {
			entry := item
			module.set(entry.name,tsNamedNativeMethod(entry.name,1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsInternalFSCPCopyError(loop,args,entry.code,entry.prefix)}))
		}
		errnos := tsNewObject()
		for _, item := range []struct{name string; value syscall.Errno}{{"EEXIST",syscall.EEXIST},{"EISDIR",syscall.EISDIR},{"EINVAL",syscall.EINVAL},{"ENOTDIR",syscall.ENOTDIR}} {errnos.set(item.name,tsNumberValue(float64(item.value)))}
		module.set("errno",tsObjectValue(errnos))
		module.set("areIdentical",tsNamedNativeMethod("areIdentical",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{
			a,b:=tsArg(args,0),tsArg(args,1);return tsBooleanValue(tsStrictEqual(tsFSProperty(a,"ino"),tsFSProperty(b,"ino"))&&tsStrictEqual(tsFSProperty(a,"dev"),tsFSProperty(b,"dev")))
		}))
		module.set("isSrcSubdir",tsNamedNativeMethod("isSrcSubdir",2,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{
		src,_:=filepath.Abs(tsFSPath(tsArg(args,0)));dest,_:=filepath.Abs(tsFSPath(tsArg(args,1)));rel,err:=filepath.Rel(src,dest);return tsBooleanValue(err==nil&&rel!=".."&&!strings.HasPrefix(rel,".."+string(filepath.Separator))&&!filepath.IsAbs(rel))
		}))
		module.set("default",tsObjectValue(module))
		tsInternalFSCPCopySyncSingleton.value=tsObjectValue(module)
	})
	return tsInternalFSCPCopySyncSingleton.value
}
`

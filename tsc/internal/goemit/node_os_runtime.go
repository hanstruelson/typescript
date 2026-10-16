//go:build !no_node_modules

package goemit

// nodeOSRuntime maps Bun's Node OS wrapper to Go's standard library. Linux
// procfs supplies the host statistics for which Go has no portable API.
const nodeOSRuntime = `
var tsNodeOSSingleton struct { sync.Once; value *tsObject }

func tsNodeOSFailure(name string, err error) {
	panic(tsThrown{tsErrorValue(&tsRuntimeError{name:name, message:err.Error()})})
}

func tsNodeOSString(name, value string) tsValue {
	return tsNamedNativeMethod(name, 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsStringReference(tsStringUTF8(value)) })
}

func tsNodeOSNumber(name string, value float64) tsValue {
	return tsNamedNativeMethod(name, 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNumberValue(value) })
}

func tsNodeOSInt32(value tsValue, name string, minimum, maximum tsValue) int {
	tsInternalValidatorInt32(value, tsStringReference(tsStringUTF8(name)), minimum, maximum)
	return int(tsNumber(value))
}

func tsNodeOSArch() string {
	switch runtime.GOARCH {
	case "amd64": return "x64"
	case "386": return "ia32"
	case "arm": return "arm"
	case "arm64": return "arm64"
	case "ppc64", "ppc64le": return "ppc64"
	case "mips", "mipsle": return "mips"
	case "mips64", "mips64le": return "mips64"
	case "s390x": return "s390x"
	case "riscv64": return "riscv64"
	default: return runtime.GOARCH
	}
}

func tsNodeOSPlatform() string {
	switch runtime.GOOS {
	case "darwin": return "darwin"
	case "windows": return "win32"
	case "freebsd": return "freebsd"
	case "openbsd": return "openbsd"
	case "netbsd": return "netbsd"
	case "dragonfly": return "dragonfly"
	case "solaris": return "sunos"
	case "illumos": return "sunos"
	default: return runtime.GOOS
	}
}

func tsNodeOSType() string {
	switch runtime.GOOS {
	case "windows": return "Windows_NT"
	case "darwin": return "Darwin"
	case "linux", "android": return "Linux"
	case "freebsd": return "FreeBSD"
	case "openbsd": return "OpenBSD"
	case "netbsd": return "NetBSD"
	case "dragonfly": return "DragonFly"
	case "solaris", "illumos": return "SunOS"
	case "aix": return "AIX"
	default: return runtime.GOOS
	}
}

func tsNodeOSMachine() string {
	arch := tsNodeOSArch()
	switch arch {
	case "x64": if runtime.GOOS == "freebsd" { return "amd64" }; return "x86_64"
	case "ia32": return "i686"
	case "arm64": if runtime.GOOS == "android" { return "aarch64" }; return "arm64"
	case "arm": return "arm"
	default: return arch
	}
}

func tsNodeOSTmpdir() string {
	if runtime.GOOS == "windows" {
		for _, key := range []string{"TEMP", "TMP"} {
			if value := os.Getenv(key); value != "" && len(value) > 1 && !strings.HasSuffix(value, ":\\") { return strings.TrimRight(value, "\\") }
		}
		root := os.Getenv("SystemRoot"); if root == "" { root = os.Getenv("windir") }
		if root != "" { return strings.TrimRight(root, "\\") + "\\temp" }
		return "C:\\temp"
	}
	for _, key := range []string{"TMPDIR", "TMP", "TEMP"} {
		if value := os.Getenv(key); value != "" { if len(value) > 1 { return strings.TrimRight(value, "/") }; return value }
	}
	if runtime.GOOS == "android" { return "/data/local/tmp" }
	return "/tmp"
}

func tsNodeOSRelease() string {
	if runtime.GOOS == "linux" || runtime.GOOS == "android" {
		if data, err := os.ReadFile("/proc/sys/kernel/osrelease"); err == nil { return strings.TrimSpace(string(data)) }
	}
	return ""
}

func tsNodeOSVersion() string {
	if data, err := os.ReadFile("/proc/version"); err == nil {
		text := strings.TrimSpace(string(data))
		release := tsNodeOSRelease()
		if strings.HasPrefix(text, "Linux version "+release+" ") { return strings.TrimPrefix(text, "Linux version "+release+" ") }
		return text
	}
	return ""
}

func tsNodeOSUptime() float64 {
	if data, err := os.ReadFile("/proc/uptime"); err == nil {
		fields := strings.Fields(string(data)); if len(fields) > 0 { if value, err := strconv.ParseFloat(fields[0], 64); err == nil { return float64(uint64(value)) } }
	}
	return 0
}

func tsNodeOSMemory() (uint64, uint64) {
	if data, err := os.ReadFile("/proc/meminfo"); err == nil {
		values := map[string]uint64{}
		for _, line := range strings.Split(string(data), "\n") {
			fields := strings.Fields(line); if len(fields) < 2 { continue }
			value, err := strconv.ParseUint(fields[1], 10, 64); if err == nil { values[strings.TrimSuffix(fields[0], ":")] = value * 1024 }
		}
		return values["MemTotal"], values["MemFree"]
	}
	return 0, 0
}

func tsNodeOSLoadavg() []tsValue {
	values := []tsValue{tsNumberValue(0), tsNumberValue(0), tsNumberValue(0)}
	if data, err := os.ReadFile("/proc/loadavg"); err == nil {
		fields := strings.Fields(string(data))
		for i := 0; i < len(values) && i < len(fields); i++ { if value, err := strconv.ParseFloat(fields[i], 64); err == nil { values[i] = tsNumberValue(value) } }
	}
	return values
}

func tsNodeOSCPUData() (string, float64, map[int][]float64) {
	model, speed := runtime.GOARCH, float64(0)
	if data, err := os.ReadFile("/proc/cpuinfo"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			parts := strings.SplitN(line, ":", 2); if len(parts) != 2 { continue }
			switch strings.TrimSpace(parts[0]) {
			case "model name", "Hardware": if model == runtime.GOARCH { model = strings.TrimSpace(parts[1]) }
			case "cpu MHz": if speed == 0 { speed, _ = strconv.ParseFloat(strings.TrimSpace(parts[1]), 64) }
			}
		}
	}
	times := map[int][]float64{}
	if data, err := os.ReadFile("/proc/stat"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			fields := strings.Fields(line); if len(fields) < 5 || !strings.HasPrefix(fields[0], "cpu") || fields[0] == "cpu" { continue }
			index, err := strconv.Atoi(strings.TrimPrefix(fields[0], "cpu")); if err != nil { continue }
			values := make([]float64, 5)
			for i, fieldIndex := range []int{1, 2, 3, 4, 6} { if fieldIndex < len(fields) { ticks, _ := strconv.ParseFloat(fields[fieldIndex], 64); values[i] = ticks * 10 } }
			times[index] = values
		}
	}
	return model, speed, times
}

func tsNodeOSCPUs() tsValue {
	model, speed, times := tsNodeOSCPUData()
	cpuIndexes:=make([]int,0,len(times));for index:=range times{cpuIndexes=append(cpuIndexes,index)};sort.Ints(cpuIndexes)
	count := runtime.NumCPU(); if len(cpuIndexes)>0 { count=len(cpuIndexes) } else { for index:=0;index<count;index++{cpuIndexes=append(cpuIndexes,index)} }; values := make([]tsValue, count)
	for index, cpuIndex := range cpuIndexes {
		cpuModel:=model;if cpuModel==runtime.GOARCH{cpuModel="unknown"}
		entry := tsNewObject(); entry.set("model", tsStringReference(tsStringUTF8(cpuModel))); entry.set("speed", tsNumberValue(speed))
		parts := times[cpuIndex]; if parts == nil { parts = make([]float64, 5) }
		cpuTimes := tsNewObject()
		for i, key := range []string{"user", "nice", "sys", "idle", "irq"} { cpuTimes.set(key, tsNumberValue(parts[i])) }
		entry.set("times", tsObjectValue(cpuTimes)); values[index] = tsObjectValue(entry)
	}
	return tsArrayValue(&tsArray{values:values})
}

func tsNodeOSNetworkInterfaces() tsValue {
	interfaces, err := net.Interfaces(); if err != nil { tsNodeOSFailure("Error", err) }
	out := tsNewObject()
	for _, iface := range interfaces {
		if runtime.GOOS!="windows"&&(iface.Flags&net.FlagUp==0||iface.Flags&net.FlagRunning==0){continue}
		addresses, err := iface.Addrs(); if err != nil { continue }
		entries := make([]tsValue, 0, len(addresses))
		for _, address := range addresses {
			ip, network, err := net.ParseCIDR(address.String()); if err != nil { continue }
			family, mask, scope := "IPv6", net.IP(network.Mask).String(), float64(0)
			if ip.IsLinkLocalUnicast() { scope = float64(iface.Index) }
			if v4 := ip.To4(); v4 != nil { family, mask, scope = "IPv4", net.IP(network.Mask).String(), 0; ip = v4 }
			item := tsNewObject(); item.set("address", tsStringReference(tsStringUTF8(ip.String()))); item.set("netmask", tsStringReference(tsStringUTF8(mask))); item.set("family", tsStringReference(tsStringUTF8(family)))
			mac:="00:00:00:00:00:00";if len(iface.HardwareAddr)>=6{mac=iface.HardwareAddr[:6].String()}
			item.set("mac", tsStringReference(tsStringUTF8(mac))); item.set("internal", tsBooleanValue(iface.Flags&net.FlagLoopback != 0)); item.set("cidr", tsStringReference(tsStringUTF8(network.String()))); if family=="IPv6"{item.set("scopeid", tsNumberValue(scope))}
			entries = append(entries, tsObjectValue(item))
		}
		if len(entries) > 0 { out.set(iface.Name, tsArrayValue(&tsArray{values:entries})) }
	}
	return tsObjectValue(out)
}

func tsNodeOSUserInfo(loop *tsLoop, args []tsValue) tsValue {
	current, err := osuser.Current(); if err != nil { tsNodeOSFailure("Error", err) }
	username := os.Getenv("USER"); if username == "" { username = current.Username }; if username == "" { username = "unknown" }
	shell := os.Getenv("SHELL"); if shell == "" { shell = os.Getenv("COMSPEC") }; if shell == "" { shell = "unknown" }
	encoding := tsGet(loop, tsArg(args, 0), "encoding")
	bufferEncoding := encoding.kind == tsStringKind && (*tsString)(encoding.ref).String() == "buffer"
	text := func(value string) tsValue { if bufferEncoding { return tsNodeBuffer([]byte(value)) }; return tsStringReference(tsStringUTF8(value)) }
	uid, gid := int64(-1), int64(-1)
	if parsed, parseErr := strconv.ParseInt(current.Uid, 10, 64); parseErr == nil { uid = parsed }
	if parsed, parseErr := strconv.ParseInt(current.Gid, 10, 64); parseErr == nil { gid = parsed }
	out := tsNewObject(); out.set("username", text(username)); out.set("uid", tsNumberValue(float64(uid))); out.set("gid", tsNumberValue(float64(gid))); out.set("homedir", text(current.HomeDir)); out.set("shell", text(shell))
	return tsObjectValue(out)
}

func tsNodeOSConstants() tsValue {
	out := tsNewObject(); errno := tsNewObject(); signals := tsNewObject(); priority := tsNewObject()
	setValues := func(target *tsObject, values map[string]float64) {
		names := make([]string, 0, len(values)); for name := range values { names = append(names, name) }; sort.Strings(names)
		for _, name := range names { target.set(name, tsNumberValue(values[name])) }
	}
	setValues(errno, map[string]float64{"EPERM":1,"ENOENT":2,"EINTR":4,"EIO":5,"ENXIO":6,"E2BIG":7,"ENOEXEC":8,"EBADF":9,"ECHILD":10,"EAGAIN":11,"ENOMEM":12,"EACCES":13,"EFAULT":14,"EBUSY":16,"EEXIST":17,"EXDEV":18,"ENODEV":19,"ENOTDIR":20,"EISDIR":21,"EINVAL":22,"ENFILE":23,"EMFILE":24,"ENOTTY":25,"EFBIG":27,"ENOSPC":28,"ESPIPE":29,"EROFS":30,"EMLINK":31,"EPIPE":32,"EDOM":33,"ERANGE":34,"EDEADLK":35,"ENAMETOOLONG":36,"ENOLCK":37,"ENOSYS":38,"ENOTEMPTY":39,"ELOOP":40,"ENOMSG":42,"EIDRM":43,"ECONNREFUSED":111,"ETIMEDOUT":110})
	setValues(signals, map[string]float64{"SIGHUP":1,"SIGINT":2,"SIGQUIT":3,"SIGILL":4,"SIGTRAP":5,"SIGABRT":6,"SIGBUS":7,"SIGFPE":8,"SIGKILL":9,"SIGUSR1":10,"SIGSEGV":11,"SIGUSR2":12,"SIGPIPE":13,"SIGALRM":14,"SIGTERM":15,"SIGCHLD":17,"SIGCONT":18,"SIGSTOP":19,"SIGTSTP":20,"SIGTTIN":21,"SIGTTOU":22,"SIGURG":23,"SIGXCPU":24,"SIGXFSZ":25,"SIGVTALRM":26,"SIGPROF":27,"SIGWINCH":28,"SIGIO":29,"SIGSYS":31})
	setValues(priority, map[string]float64{"PRIORITY_LOW":19,"PRIORITY_BELOW_NORMAL":10,"PRIORITY_NORMAL":0,"PRIORITY_ABOVE_NORMAL":-7,"PRIORITY_HIGH":-14,"PRIORITY_HIGHEST":-20})
	out.set("errno", tsObjectValue(errno)); out.set("signals", tsObjectValue(signals)); out.set("priority", tsObjectValue(priority))
	return tsObjectValue(out)
}

func tsNodeOSModule() tsValue {
	tsNodeOSSingleton.Do(func() {
		out := tsNewObject(); total, free := tsNodeOSMemory()
		out.set("EOL", tsStringReference(tsStringUTF8("\n"))); if runtime.GOOS == "windows" { out.set("EOL", tsStringReference(tsStringUTF8("\r\n"))) }
		out.set("devNull", tsStringReference(tsStringUTF8("/dev/null"))); if runtime.GOOS == "windows" { out.set("devNull", tsStringReference(tsStringUTF8("\\\\.\\nul"))) }
		out.set("constants", tsNodeOSConstants())
		out.set("arch", tsNodeOSString("arch", tsNodeOSArch())); out.set("platform", tsNodeOSString("platform", tsNodeOSPlatform())); out.set("type", tsNodeOSString("type", tsNodeOSType())); out.set("machine", tsNodeOSString("machine", tsNodeOSMachine())); out.set("release", tsNodeOSString("release", tsNodeOSRelease())); out.set("version", tsNodeOSString("version", tsNodeOSVersion()))
		out.set("tmpdir", tsNamedNativeMethod("tmpdir", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsStringReference(tsStringUTF8(tsNodeOSTmpdir())) }))
		out.set("homedir", tsNamedNativeMethod("homedir", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { value, err := os.UserHomeDir(); if err != nil { tsNodeOSFailure("Error", err) }; return tsStringReference(tsStringUTF8(value)) }))
		out.set("hostname", tsNamedNativeMethod("hostname", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { value, err := os.Hostname(); if err != nil { tsNodeOSFailure("Error", err) }; return tsStringReference(tsStringUTF8(value)) }))
		out.set("availableParallelism", tsNodeOSNumber("availableParallelism", float64(runtime.NumCPU())))
		out.set("totalmem", tsNodeOSNumber("totalmem", float64(total))); out.set("freemem", tsNodeOSNumber("freemem", float64(free)))
		out.set("uptime", tsNamedNativeMethod("uptime", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNumberValue(tsNodeOSUptime()) }))
		out.set("loadavg", tsNamedNativeMethod("loadavg", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsArrayValue(&tsArray{values:tsNodeOSLoadavg()}) }))
		out.set("cpus", tsNamedNativeMethod("cpus", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNodeOSCPUs() }))
		out.set("networkInterfaces", tsNamedNativeMethod("networkInterfaces", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { return tsNodeOSNetworkInterfaces() }))
		out.set("userInfo", tsNamedNativeMethod("userInfo", 0, func(loop *tsLoop, _ tsValue, args ...tsValue) tsValue { return tsNodeOSUserInfo(loop, args) }))
		out.set("endianness", tsNamedNativeMethod("endianness", 0, func(_ *tsLoop, _ tsValue, _ ...tsValue) tsValue { endian := "LE"; switch runtime.GOARCH { case "ppc64", "s390x", "mips", "mips64": endian = "BE" }; return tsStringReference(tsStringUTF8(endian)) }))
		out.set("getPriority", tsNamedNativeMethod("getPriority", 1, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { pid := 0; if len(args)>0 && !tsIsUndefined(args[0]) { pid = tsNodeOSInt32(args[0],"pid",tsU,tsU) }; value, err := syscall.Getpriority(syscall.PRIO_PROCESS, pid); if err != nil { tsNodeOSFailure("Error", err) }; return tsNumberValue(float64(value)) }))
		out.set("setPriority", tsNamedNativeMethod("setPriority", 2, func(_ *tsLoop, _ tsValue, args ...tsValue) tsValue { pid, priority := 0, tsU; if len(args)==1 { priority = tsNodeOSInt32(args[0],"priority",tsNumberValue(-20),tsNumberValue(19)) } else { pid=tsNodeOSInt32(tsArg(args,0),"pid",tsU,tsU); priority=tsNodeOSInt32(tsArg(args,1),"priority",tsNumberValue(-20),tsNumberValue(19)) }; if err := syscall.Setpriority(syscall.PRIO_PROCESS, pid, priority); err != nil { tsNodeOSFailure("Error", err) }; return tsU }))
		tsNodeOSSingleton.value = out
	})
	return tsObjectValue(tsNodeOSSingleton.value)
}
`

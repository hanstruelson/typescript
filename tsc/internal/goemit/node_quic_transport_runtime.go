package goemit

// nodeQUICTransportRuntime supplies the Go transport behind the Bun-compatible
// QUIC objects. The QUIC wire implementation is quic-go; runtime values stay
// on their owning event loop and socket work runs in its task pool.
const nodeQUICTransportRuntime = `
type tsNodeQUICTransportEndpoint struct {
	listener *quicgo.Listener
	cancel context.CancelFunc
}

type tsNodeQUICTransportSession struct {
	conn *quicgo.Conn
	control *quicgo.SendStream
	controlMu sync.Mutex
	http3 bool
	ready chan struct{}
	readyOnce sync.Once
	peerMaxFieldSectionSize uint64
	peerH3Datagrams *bool
}

type tsNodeQUICTransportStream struct {
	reader io.Reader
	writer io.Writer
	tryWrite func([]byte) error
	closeWrite func() error
	cancelRead func(quicgo.StreamErrorCode)
	cancelWrite func(quicgo.StreamErrorCode)
	canRead, canWrite bool
	http3 bool
	headersSent bool
	http3Buffer []byte
	http3FrameType, http3FrameRemaining, http3MaxHeaderLength uint64
	http3MaxHeaderPairs, http3MaxFieldSectionSize uint64
	http3FrameActive bool
	http3HeaderBlock []byte
	writeQueue []tsNodeQUICWriteRequest
	writeBusy bool
	queuedBytes int
	totalBytes int64
}

type tsNodeQUICWriteRequest struct {
	data []byte
	callback tsValue
	count int
	end bool
}

var tsNodeQUICTransportRegistry struct {
	sync.Once
	endpoints map[*tsObject]*tsNodeQUICTransportEndpoint
	sessions map[*tsObject]*tsNodeQUICTransportSession
	streams map[*tsObject]*tsNodeQUICTransportStream
}

func tsNodeQUICTransportInit() {
	tsNodeQUICTransportRegistry.Do(func() {
		tsNodeQUICTransportRegistry.endpoints = map[*tsObject]*tsNodeQUICTransportEndpoint{}
		tsNodeQUICTransportRegistry.sessions = map[*tsObject]*tsNodeQUICTransportSession{}
		tsNodeQUICTransportRegistry.streams = map[*tsObject]*tsNodeQUICTransportStream{}
	})
}

func tsNodeQUICOption(loop *tsLoop, value tsValue, name string) tsValue {
	return tsNodeQUICGet(loop, value, name)
}

func tsNodeQUICOptionNumber(loop *tsLoop, value tsValue, name string) (float64, bool) {
	option := tsNodeQUICOption(loop, value, name)
	if option.kind == tsBigIntKind { number, _ := new(big.Float).SetInt((*big.Int)(option.ref)).Float64(); return number, true }
	if tsIsNumeric(option) { return tsNumber(option), true }
	return 0, false
}

func tsNodeQUICTLSOptions(loop *tsLoop, options tsValue) tsValue {
	tlsOptions := tsNodeQUICOption(loop, options, "tls")
	if tsNodeQUICObject(tlsOptions) {
		return tlsOptions
	}
	return options
}

func tsNodeQUICInputBytes(loop *tsLoop, value tsValue) []byte {
	if value.kind == tsStringKind {
		return []byte(tsStringValueWithLoop(loop, value).String())
	}
	return append([]byte(nil), tsFSBytes(value)...)
}

func tsNodeQUICInputList(value tsValue) []tsValue {
	if value.kind == tsArrayKind {
		array := (*tsArray)(value.ref)
		values := make([]tsValue, 0, array.length())
		for index := 0; index < array.length(); index++ {
			values = append(values, array.at(index))
		}
		return values
	}
	if tsIsUndefined(value) || tsIsNull(value) {
		return nil
	}
	return []tsValue{value}
}

func tsNodeQUICCertificatePEM(loop *tsLoop, values []tsValue) []byte {
	var output []byte
	for _, value := range values {
		data := tsNodeQUICInputBytes(loop, value)
		if len(data) == 0 {
			continue
		}
		if bytes.Contains(data, []byte("-----BEGIN")) {
			output = append(output, data...)
			if data[len(data)-1] != '\n' {
				output = append(output, '\n')
			}
			continue
		}
		output = append(output, pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: data})...)
	}
	return output
}

func tsNodeQUICFirstInput(loop *tsLoop, value tsValue) []byte {
	values := tsNodeQUICInputList(value)
	if len(values) == 0 {
		return nil
	}
	return tsNodeQUICInputBytes(loop, values[0])
}

func tsNodeQUICTLSIdentity(loop *tsLoop, identity tsValue) (cryptotls.Certificate, error) {
	certs := tsNodeQUICInputList(tsNodeQUICOption(loop, identity, "certs"))
	keys := tsNodeQUICInputList(tsNodeQUICOption(loop, identity, "keys"))
	if len(certs) == 0 || len(keys) == 0 { return cryptotls.Certificate{}, fmt.Errorf("SNI identity requires certs and keys") }
	keyPEM, err := tsNodeQUICPrivateKeyPEM(loop, keys[0])
	if err != nil { return cryptotls.Certificate{}, err }
	return cryptotls.X509KeyPair(tsNodeQUICCertificatePEM(loop, certs), keyPEM)
}

func tsNodeQUICCertPool(loop *tsLoop, value tsValue) (*x509.CertPool, error) {
	values := tsNodeQUICInputList(value)
	if len(values) == 0 { return nil, nil }
	pool := cryptotls.NewCertPool()
	for _, item := range values {
		data := tsNodeQUICInputBytes(loop, item)
		if !bytes.Contains(data, []byte("-----BEGIN")) {
			data = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: data})
		}
		if !pool.AppendCertsFromPEM(data) { return nil, fmt.Errorf("invalid QUIC CA certificate") }
	}
	return pool, nil
}

func tsNodeQUICPrivateKeyPEM(loop *tsLoop, value tsValue) ([]byte, error) {
	if key := tsNodeCryptoKeyFor(value); key != nil && key.kind == "private" && key.private != nil {
		der, blockType, err := tsNodeCryptoMarshalPrivate(key.private, "pkcs8")
		if err != nil { return nil, err }
		return pem.EncodeToMemory(&pem.Block{Type: blockType, Bytes: der}), nil
	}
	data := tsNodeQUICInputBytes(loop, value)
	if block, _ := pem.Decode(data); block != nil {
		return data, nil
	}
	if _, err := x509.ParsePKCS8PrivateKey(data); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: data}), nil
	}
	if _, err := x509.ParsePKCS1PrivateKey(data); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: data}), nil
	}
	if _, err := x509.ParseECPrivateKey(data); err == nil {
		return pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: data}), nil
	}
	return nil, fmt.Errorf("invalid QUIC private key")
}

func tsNodeQUICSetSNIContexts(loop *tsLoop, endpoint tsValue, entries, options tsValue) tsValue {
	state := tsNodeQUICEndpointFor(endpoint)
	if state == nil || state.closed { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "Endpoint is destroyed")}) }
	if entries.kind != tsObjectKind && entries.kind != tsInstanceKind { tsFSInvalid("entries", "an object") }
	replace := tsNodeQUICOption(loop, options, "replace")
	if !tsIsUndefined(replace) && replace.kind != tsBooleanKind { tsFSInvalid("options.replace", "a boolean") }
	identityMap := (*tsObject)(entries.ref)
	contexts := make(map[string]cryptotls.Certificate, len(identityMap.order))
	for _, hostname := range identityMap.order {
		identityOptions := identityMap.values[hostname]
		if !tsNodeQUICObject(identityOptions) { tsFSInvalid("entries["+strconv.Quote(hostname)+"]", "an object") }
		if verify := tsNodeQUICOption(loop, identityOptions, "verifyPrivateKey"); !tsIsUndefined(verify) && verify.kind != tsBooleanKind { tsFSInvalid("entries["+strconv.Quote(hostname)+"].verifyPrivateKey", "a boolean") }
		keys := tsNodeQUICInputList(tsNodeQUICOption(loop, identityOptions, "keys"))
		if len(keys) == 0 { tsFSInvalid("entries["+strconv.Quote(hostname)+"].keys", "a private KeyObject") }
		for _, keyValue := range keys { key := tsNodeCryptoKeyFor(keyValue); if key == nil || key.kind != "private" { tsFSInvalid("entries["+strconv.Quote(hostname)+"].keys", "a private KeyObject") } }
		certsValue := tsNodeQUICOption(loop, identityOptions, "certs")
		if tsIsUndefined(certsValue) { tsFSInvalid("entries["+strconv.Quote(hostname)+"].certs", "an ArrayBuffer or ArrayBufferView") }
		for _, cert := range tsNodeQUICInputList(certsValue) { if cert.kind != tsArrayBufferKind && cert.kind != tsTypedArrayKind && cert.kind != tsDataViewKind { tsFSInvalid("entries["+strconv.Quote(hostname)+"].certs", "an ArrayBuffer or ArrayBufferView") } }
		identity, err := tsNodeQUICTLSIdentity(loop, identityOptions)
		if err != nil { panic(tsThrown{tsNodeQUICError("ERR_TLS_CERT_ALTNAME_INVALID", fmt.Sprintf("entries[%q]: %v", hostname, err))}) }
		contexts[strings.ToLower(strings.TrimSuffix(hostname, "."))] = identity
	}
	state.sniMu.Lock()
	if state.sniCertificates == nil { state.sniCertificates = map[string]cryptotls.Certificate{} }
	if tsTruthy(replace) { state.sniCertificates = contexts } else { for hostname, identity := range contexts { state.sniCertificates[hostname] = identity } }
	state.sniMu.Unlock()
	return tsU
}

func tsNodeQUICInstallSNICallback(loop *tsLoop, endpoint, options tsValue, config *cryptotls.Config) {
	state := tsNodeQUICEndpointFor(endpoint)
	if state == nil { return }
	if identities := tsNodeQUICOption(loop, tsNodeQUICTLSOptions(loop, options), "sni"); tsNodeQUICObject(identities) {
		_ = tsNodeQUICSetSNIContexts(loop, endpoint, identities, tsU)
	}
	config.GetCertificate = func(hello *cryptotls.ClientHelloInfo) (*cryptotls.Certificate, error) {
		hostname := strings.ToLower(strings.TrimSuffix(hello.ServerName, "."))
		state.sniMu.RLock()
		identity, found := state.sniCertificates[hostname]
		if !found {
			patterns := make([]string, 0, len(state.sniCertificates))
			for pattern := range state.sniCertificates { if strings.HasPrefix(pattern, "*.") { patterns = append(patterns, pattern) } }
			sort.Slice(patterns, func(i, j int) bool { return len(patterns[i]) > len(patterns[j]) })
			for _, pattern := range patterns {
				prefix := strings.TrimSuffix(hostname, pattern[1:])
				if prefix != hostname && prefix != "" && !strings.Contains(prefix, ".") { identity, found = state.sniCertificates[pattern], true; break }
			}
		}
		if !found { identity, found = state.sniCertificates["*"] }
		state.sniMu.RUnlock()
		if found { return &identity, nil }
		if len(config.Certificates) != 0 { return &config.Certificates[0], nil }
		return nil, fmt.Errorf("no QUIC certificate for server name %q", hello.ServerName)
	}
}

func tsNodeQUICMakeTLSConfig(loop *tsLoop, options tsValue, server bool, remoteHost string) (*cryptotls.Config, error) {
	tlsOptions := tsNodeQUICTLSOptions(loop, options)
	identityOptions := tlsOptions
	config := &cryptotls.Config{MinVersion: cryptotls.VersionTLS13}
	groupList := tsNodeQUICOption(loop, tlsOptions, "groups")
	if tsIsUndefined(groupList) { groupList = tsNodeQUICKey("X25519:P-256:P-384:P-521") }
	if groupList.kind == tsStringKind {
		for _, group := range strings.Split(tsText(groupList), ":") {
			curve := map[string]cryptotls.CurveID{"X25519": cryptotls.X25519, "P-256": cryptotls.CurveP256, "P-384": cryptotls.CurveP384, "P-521": cryptotls.CurveP521}[group]
			if curve == 0 { return nil, fmt.Errorf("unsupported QUIC TLS group %q", group) }
			config.CurvePreferences = append(config.CurvePreferences, curve)
		}
	}
	if server {
		identity := tsNodeQUICOption(loop, tlsOptions, "sni")
		if tsNodeQUICObject(identity) {
			identityObject := (*tsObject)(identity.ref)
			selected := identityObject.values["*"]
			if !tsNodeQUICObject(selected) && len(identityObject.order) != 0 {
				selected = identityObject.values[identityObject.order[0]]
			}
			if tsNodeQUICObject(selected) {
				identityOptions = selected
			}
		}
	}
	certValues := tsNodeQUICInputList(tsNodeQUICOption(loop, identityOptions, "certs"))
	keyValues := tsNodeQUICInputList(tsNodeQUICOption(loop, identityOptions, "keys"))
	if len(certValues) != 0 && len(keyValues) != 0 {
		certificatePEM := tsNodeQUICCertificatePEM(loop, certValues)
		keyPEM, err := tsNodeQUICPrivateKeyPEM(loop, keyValues[0])
		if err != nil { return nil, err }
		certificate, err := cryptotls.X509KeyPair(certificatePEM, keyPEM)
		if err != nil {
			return nil, err
		}
		config.Certificates = []cryptotls.Certificate{certificate}
	}
	if server && len(config.Certificates) == 0 {
		return nil, fmt.Errorf("server QUIC sessions require a certificate and private key")
	}
	if server {
		if caPool, err := tsNodeQUICCertPool(loop, tsNodeQUICOption(loop, tlsOptions, "ca")); err != nil {
			return nil, err
		} else if caPool != nil {
			config.ClientCAs = caPool
		}
		if tsTruthy(tsNodeQUICOption(loop, options, "verifyClient")) || tsTruthy(tsNodeQUICOption(loop, tlsOptions, "verifyClient")) {
			config.ClientAuth = cryptotls.RequireAndVerifyClientCert
		}
	}
	if !server {
		serverName := tsNodeQUICOption(loop, tlsOptions, "servername")
		if tsIsUndefined(serverName) || tsIsNull(serverName) || tsText(serverName) == "" {
			config.ServerName = remoteHost
		} else {
			config.ServerName = tsText(serverName)
		}
		config.InsecureSkipVerify = tsStrictEqual(tsNodeQUICOption(loop, tlsOptions, "rejectUnauthorized"), tsBooleanValue(false)) || tsText(tsNodeQUICOption(loop, options, "verifyPeer")) == "manual"
		caPool, err := tsNodeQUICCertPool(loop, tsNodeQUICOption(loop, tlsOptions, "ca"))
		if err != nil { return nil, err }
		config.RootCAs = caPool
	}
	config.NextProtos = tsNodeQUICALPN(loop, tlsOptions, server)
	if len(config.NextProtos) == 0 {
		config.NextProtos = []string{"h3"}
	}
	return config, nil
}

func tsNodeQUICLocalCertificate(loop *tsLoop, options tsValue, serverNames ...string) *x509.Certificate {
	tlsOptions := tsNodeQUICTLSOptions(loop, options)
	identityOptions := tlsOptions
	identities := tsNodeQUICOption(loop, tlsOptions, "sni")
	if tsNodeQUICObject(identities) {
		name := ""
		if len(serverNames) != 0 { name = strings.ToLower(strings.TrimSuffix(serverNames[0], ".")) }
		identityMap := (*tsObject)(identities.ref)
		selected := identityMap.values[name]
		if !tsNodeQUICObject(selected) {
			for pattern, candidate := range identityMap.values {
				wildcard := strings.ToLower(pattern)
				if !strings.HasPrefix(wildcard, "*.") || !strings.HasSuffix(name, wildcard[1:]) { continue }
				label := strings.TrimSuffix(name, wildcard[1:])
				if label != "" && !strings.Contains(label, ".") { selected = candidate; break }
			}
		}
		if !tsNodeQUICObject(selected) { selected = identityMap.values["*"] }
		if !tsNodeQUICObject(selected) && len(identityMap.order) != 0 { selected = identityMap.values[identityMap.order[0]] }
		if tsNodeQUICObject(selected) { identityOptions = selected }
	}
	certs := tsNodeQUICInputList(tsNodeQUICOption(loop, identityOptions, "certs"))
	if len(certs) == 0 { return nil }
	data := tsNodeQUICCertificatePEM(loop, certs)
	block, _ := pem.Decode(data)
	if block != nil && block.Type == "CERTIFICATE" { certificate, err := x509.ParseCertificate(block.Bytes); if err == nil { return certificate }; return nil }
	if parsed, err := x509.ParseCertificate(tsNodeQUICFirstInput(loop, tsArrayValue(&tsArray{values: certs}))); err == nil { return parsed }
	return nil
}

func tsNodeQUICALPN(loop *tsLoop, options tsValue, server bool) []string {
	value := tsNodeQUICOption(loop, options, "alpn")
	if tsIsUndefined(value) || tsIsNull(value) {
		return nil
	}
	if value.kind == tsTypedArrayKind || value.kind == tsDataViewKind || value.kind == tsArrayBufferKind {
		data := tsFSBytes(value)
		var protocols []string
		for offset := 0; offset < len(data); {
			length := int(data[offset])
			offset++
			if length == 0 || offset+length > len(data) {
				return nil
			}
			protocols = append(protocols, string(data[offset:offset+length]))
			offset += length
		}
		return protocols
	}
	if value.kind == tsArrayKind && server {
		array := (*tsArray)(value.ref)
		protocols := make([]string, 0, array.length())
		for index := 0; index < array.length(); index++ {
			protocols = append(protocols, tsText(array.at(index)))
		}
		return protocols
	}
	return []string{tsText(value)}
}

func tsNodeQUICValidateIntegerOption(loop *tsLoop, options tsValue, name string, allowBigInt bool, maximum uint64) {
	if tsNullish(options) { return }
	value := tsNodeQUICOption(loop, options, name)
	if tsIsUndefined(value) {
		return
	}
	if value.kind == tsBigIntKind && allowBigInt {
		integer := (*big.Int)(value.ref)
		if integer.Sign() >= 0 && (maximum == 0 || integer.IsUint64() && integer.Uint64() <= maximum) {
			return
		}
	}
	if tsIsNumeric(value) {
		number := tsNumber(value)
		if number >= 0 && math.Trunc(number) == number && !math.IsInf(number, 0) && !math.IsNaN(number) && (maximum == 0 || number <= float64(maximum)) {
			return
		}
	}
	tsFSInvalid("options."+name, "a non-negative integer")
}

func tsNodeQUICValidateTLSBytes(loop *tsLoop, options tsValue, name string) {
	value := tsNodeQUICOption(loop, options, name)
	if tsIsUndefined(value) { return }
	values := tsNodeQUICInputList(value)
	for _, item := range values {
		if item.kind != tsArrayBufferKind && item.kind != tsTypedArrayKind && item.kind != tsDataViewKind {
			tsFSInvalid("options."+name, "an ArrayBuffer or ArrayBufferView")
		}
	}
}

func tsNodeQUICValidateTLSKeys(loop *tsLoop, options tsValue, name string) {
	value := tsNodeQUICOption(loop, options, name)
	if tsIsUndefined(value) { return }
	for _, item := range tsNodeQUICInputList(value) {
		key := tsNodeCryptoKeyFor(item)
		if key == nil || key.kind != "private" { tsFSInvalid("options."+name, "a private KeyObject") }
	}
}

func tsNodeQUICReadUint32BE(data []byte, offset int) uint32 {
	return uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
}

func tsNodeQUICSessionTicketValid(value tsValue) bool {
	data := tsFSBytes(value)
	if len(data) < 16 || tsNodeQUICReadUint32BE(data, 4) != 1 {
		return false
	}
	ticketSize := uint64(tsNodeQUICReadUint32BE(data, 8))
	transportOffset := uint64(12) + ticketSize
	if transportOffset+4 > uint64(len(data)) {
		return false
	}
	transportSize := uint64(tsNodeQUICReadUint32BE(data, int(transportOffset)))
	return transportOffset+4+transportSize == uint64(len(data))
}

func tsNodeQUICValidateSessionOptions(loop *tsLoop, options tsValue, server bool) {
	if tsIsUndefined(options) {
		return
	}
	if !tsNodeQUICObject(options) {
		tsFSInvalid("options", "an object")
	}
	if policy := tsNodeQUICOption(loop, options, "preferredAddressPolicy"); !tsIsUndefined(policy) && (policy.kind != tsStringKind || (tsText(policy) != "use" && tsText(policy) != "ignore" && tsText(policy) != "default")) {
		tsFSInvalid("options.preferredAddressPolicy", "'use', 'ignore', or 'default'")
	}
	if value := tsNodeQUICOption(loop, options, "cc"); !tsIsUndefined(value) && (value.kind != tsStringKind || (tsText(value) != "reno" && tsText(value) != "cubic" && tsText(value) != "bbr")) {
		tsFSInvalid("options.cc", "'reno', 'cubic', or 'bbr'")
	}
	if value := tsNodeQUICOption(loop, options, "datagramDropPolicy"); !tsIsUndefined(value) && (value.kind != tsStringKind || (tsText(value) != "drop-oldest" && tsText(value) != "drop-newest")) {
		tsFSInvalid("options.datagramDropPolicy", "'drop-oldest' or 'drop-newest'")
	}
	for _, option := range []struct{name string; minimum, maximum uint64}{
		{"drainingPeriodMultiplier", 3, 255}, {"maxDatagramSendAttempts", 1, 255},
	} {
		value := tsNodeQUICOption(loop, options, option.name)
		if !tsIsUndefined(value) {
			if !tsIsNumeric(value) || math.Trunc(tsNumber(value)) != tsNumber(value) || tsNumber(value) < float64(option.minimum) || tsNumber(value) > float64(option.maximum) {
				tsFSInvalid("options."+option.name, "an integer in the supported range")
			}
		}
	}
	if value := tsNodeQUICOption(loop, options, "verifyPeer"); !tsIsUndefined(value) && (value.kind != tsStringKind || (tsText(value) != "strict" && tsText(value) != "auto" && tsText(value) != "manual")) {
		tsFSInvalid("options.verifyPeer", "'strict', 'auto', or 'manual'")
	}
	if value := tsNodeQUICOption(loop, options, "transportParams"); !tsIsUndefined(value) && !tsNodeQUICObject(value) {
		tsFSInvalid("options.transportParams", "an object")
	}
	transportParams := tsNodeQUICOption(loop, options, "transportParams")
	for _, name := range []string{"initialMaxStreamDataBidiLocal", "initialMaxStreamDataBidiRemote", "initialMaxStreamDataUni", "initialMaxData", "initialMaxStreamsBidi", "initialMaxStreamsUni", "maxIdleTimeout", "maxUdpPayloadSize", "maxAckDelay", "ackDelayExponent", "activeConnectionIDLimit", "maxDatagramFrameSize"} {
		maximum := uint64(0)
		if name == "maxDatagramFrameSize" { maximum = 65535 }
		tsNodeQUICValidateIntegerOption(loop, transportParams, name, true, maximum)
	}
	for _, preferred := range []struct{name, family string}{
		{"preferredAddressIpv4", "ipv4"}, {"preferredAddressIpv6", "ipv6"},
	} {
		address := tsNodeQUICOption(loop, transportParams, preferred.name)
		if tsIsUndefined(address) {
			continue
		}
		if !tsNodeNetIsSocketAddress(address) {
			tsFSInvalid("options.transportParams."+preferred.name, "a SocketAddress")
		}
		if family := tsText(tsNodeQUICGet(loop, address, "family")); family != preferred.family {
			tsFSInvalid("options.transportParams."+preferred.name, "a "+preferred.family+" address")
		}
	}
	application := tsNodeQUICOption(loop, options, "application")
	if !tsIsUndefined(application) && !tsNodeQUICObject(application) {
		tsFSInvalid("options.application", "an object")
	}
	for _, name := range []string{"maxHeaderPairs", "maxHeaderLength", "maxFieldSectionSize", "qpackMaxDTableCapacity", "qpackEncoderMaxDTableCapacity", "qpackBlockedStreams"} {
		tsNodeQUICValidateIntegerOption(loop, application, name, true, 0)
	}
	for _, name := range []string{"sessionTicket", "token"} {
		value := tsNodeQUICOption(loop, options, name)
		if !tsIsUndefined(value) && value.kind != tsTypedArrayKind && value.kind != tsDataViewKind {
			tsFSInvalid("options."+name, "an ArrayBufferView")
		}
	}
	if ticket := tsNodeQUICOption(loop, options, "sessionTicket"); !tsIsUndefined(ticket) && !tsNodeQUICSessionTicketValid(ticket) {
		panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_VALUE", "The \"options.sessionTicket\" argument is not a valid session ticket")})
	}
	for _, name := range []string{"onerror", "onstream", "ondatagram", "ondatagramstatus", "onpathvalidation", "onsessionticket", "onversionnegotiation", "onhandshake", "onnewtoken", "onearlyrejected", "onorigin", "ongoaway", "onkeylog", "onqlog", "onheaders", "ontrailers", "oninfo", "onwanttrailers"} {
		callback := tsNodeQUICOption(loop, options, name)
		if tsTruthy(callback) && callback.kind != tsFunctionKind {
			tsFSInvalid("options."+name, "a function")
		}
	}
	tlsOptions := tsNodeQUICTLSOptions(loop, options)
	if !tsNodeQUICObject(tlsOptions) {
		tsFSInvalid("options", "an object")
	}
	for _, name := range []string{"servername", "ciphers", "groups"} {
		value := tsNodeQUICOption(loop, tlsOptions, name)
		if !tsIsUndefined(value) && value.kind != tsStringKind {
			tsFSInvalid("options."+name, "a string")
		}
	}
	for _, name := range []string{"keylog", "verifyClient", "rejectUnauthorized", "enableEarlyData", "tlsTrace", "verifyPrivateKey"} {
		value := tsNodeQUICOption(loop, tlsOptions, name)
		if !tsIsUndefined(value) && value.kind != tsBooleanKind {
			tsFSInvalid("options."+name, "a boolean")
		}
	}
	if value := tsNodeQUICOption(loop, tlsOptions, "alpn"); !tsIsUndefined(value) && !tsIsNull(value) {
		protocols := []tsValue{value}
		if value.kind == tsArrayKind && server {
			protocols = (*tsArray)(value.ref).values
		}
		var serverMode *bool
		for index, protocol := range protocols {
			if protocol.kind != tsStringKind || tsText(protocol) == "" || len([]rune(tsText(protocol))) > 255 {
				tsFSInvalid(fmt.Sprintf("options.alpn[%d]", index), "a non-empty string of at most 255 characters")
			}
			if server {
				name := tsText(protocol)
				isHTTP := name == "h3" || strings.HasPrefix(name, "h3-")
				if serverMode != nil && *serverMode != isHTTP {
					panic(tsThrown{tsFSException("TypeError", "ERR_INVALID_ARG_VALUE", "options.alpn cannot mix HTTP/3 and non-HTTP/3 protocols on one endpoint; use a separate QuicEndpoint for each")})
				}
				serverMode = &isHTTP
			}
		}
		if value.kind == tsArrayKind && !server { tsFSInvalid("options.alpn", "a string") }
	}
	if server {
		sni := tsNodeQUICOption(loop, tlsOptions, "sni")
		if !tsNodeQUICObject(sni) || len((*tsObject)(sni.ref).order) == 0 {
			tsFSInvalid("options.sni", "a non-empty object")
		}
		for _, hostname := range (*tsObject)(sni.ref).order {
			identity := (*tsObject)(sni.ref).values[hostname]
			if !tsNodeQUICObject(identity) { tsFSInvalid("options.sni["+strconv.Quote(hostname)+"]", "an object") }
			if tsIsUndefined(tsNodeQUICOption(loop, identity, "keys")) { tsFSInvalid("options.sni["+strconv.Quote(hostname)+"].keys", "a private KeyObject") }
			if tsIsUndefined(tsNodeQUICOption(loop, identity, "certs")) { tsFSInvalid("options.sni["+strconv.Quote(hostname)+"].certs", "an ArrayBuffer or ArrayBufferView") }
			tsNodeQUICValidateTLSKeys(loop, identity, "keys")
			tsNodeQUICValidateTLSBytes(loop, identity, "certs")
		}
	} else {
		tsNodeQUICValidateTLSKeys(loop, tlsOptions, "keys")
		tsNodeQUICValidateTLSBytes(loop, tlsOptions, "certs")
	}
	tsNodeQUICValidateTLSBytes(loop, tlsOptions, "ca")
	tsNodeQUICValidateTLSBytes(loop, tlsOptions, "crl")
}

func tsNodeQUICTransportConfig(loop *tsLoop, options tsValue) *quicgo.Config {
	// Bun's Node-compatible default is 10 seconds; quic-go's library default
	// is longer, so pin it even when no transport parameter was supplied.
	config := &quicgo.Config{MaxIdleTimeout: 10 * time.Second, EnableDatagrams: true}
	if timeout, _ := tsNodeQUICOptionNumber(loop, options, "handshakeTimeout"); timeout > 0 {
		config.HandshakeIdleTimeout = time.Duration(timeout) * time.Millisecond
	}
	transportParams := tsNodeQUICOption(loop, options, "transportParams")
	if value := tsNodeQUICOption(loop, transportParams, "maxDatagramFrameSize"); !tsIsUndefined(value) {
		frameSize, _ := tsNodeQUICOptionNumber(loop, transportParams, "maxDatagramFrameSize")
		config.EnableDatagrams = frameSize != 0
	}
	if timeout, ok := tsNodeQUICOptionNumber(loop, transportParams, "maxIdleTimeout"); ok {
		// Node exposes maxIdleTimeout in seconds; handshakeTimeout and
		// keepAlive above are milliseconds.
		if timeout == 0 {
			// quic-go interprets zero as its own default, so the wire-level
			// disabled timeout cannot be requested exactly through this field.
			config.MaxIdleTimeout = 0
		} else if timeout > 0 && timeout < float64(math.MaxInt64)/float64(time.Second) {
			config.MaxIdleTimeout = time.Duration(timeout * float64(time.Second))
		}
	}
	if timeout, _ := tsNodeQUICOptionNumber(loop, options, "keepAlive"); timeout > 0 {
		config.KeepAlivePeriod = time.Duration(timeout) * time.Millisecond
	}
	if value := tsNodeQUICOption(loop, transportParams, "initialMaxStreamsBidi"); !tsIsUndefined(value) {
		count, _ := tsNodeQUICOptionNumber(loop, transportParams, "initialMaxStreamsBidi")
		config.MaxIncomingStreams = int64(count)
	}
	if value := tsNodeQUICOption(loop, transportParams, "initialMaxStreamsUni"); !tsIsUndefined(value) {
		count, _ := tsNodeQUICOptionNumber(loop, transportParams, "initialMaxStreamsUni")
		config.MaxIncomingUniStreams = int64(count)
	}
	if window, _ := tsNodeQUICOptionNumber(loop, options, "maxStreamWindow"); window > 0 {
		config.MaxStreamReceiveWindow = uint64(window)
		config.InitialStreamReceiveWindow = uint64(window)
	}
	if window, _ := tsNodeQUICOptionNumber(loop, options, "maxWindow"); window > 0 {
		config.MaxConnectionReceiveWindow = uint64(window)
		config.InitialConnectionReceiveWindow = uint64(window)
	}
	version, _ := tsNodeQUICOptionNumber(loop, options, "version")
	minVersion, _ := tsNodeQUICOptionNumber(loop, options, "minVersion")
	if version != 0 || minVersion != 0 {
		for _, supported := range quicgo.SupportedVersions() {
			value := float64(supported)
			if version != 0 && value != version { continue }
			if minVersion != 0 && value < minVersion { continue }
			config.Versions = append(config.Versions, supported)
		}
	}
	return config
}

func tsNodeQUICSplitAddress(loop *tsLoop, address tsValue, defaultPort string) (string, string, error) {
	host, port := "", defaultPort
	if address.kind == tsStringKind {
		text := tsText(address)
		if parsed, err := url.Parse(text); err == nil && parsed.Host != "" {
			text = parsed.Host
		}
		if parsedHost, parsedPort, err := net.SplitHostPort(text); err == nil {
			host, port = parsedHost, parsedPort
		} else if strings.HasPrefix(text, "[") && strings.HasSuffix(text, "]") {
			host = strings.Trim(text, "[]")
		} else if net.ParseIP(text) != nil {
			host = text
		} else if strings.Count(text, ":") == 1 {
			parts := strings.SplitN(text, ":", 2)
			host, port = parts[0], parts[1]
		} else {
			host = text
		}
	} else if tsNodeQUICObject(address) {
		hostValue := tsNodeQUICOption(loop, address, "address")
		if tsIsUndefined(hostValue) {
			hostValue = tsNodeQUICOption(loop, address, "host")
		}
		host = tsText(hostValue)
		portValue := tsNodeQUICOption(loop, address, "port")
		if !tsIsUndefined(portValue) && !tsIsNull(portValue) {
			port = strconv.Itoa(int(tsNumber(portValue)))
		}
	} else {
		return "", "", fmt.Errorf("address must be a string or SocketAddress")
	}
	if host == "" {
		host = "localhost"
	}
	if port == "" {
		return "", "", fmt.Errorf("address must include a port")
	}
	if _, err := strconv.Atoi(port); err != nil {
		return "", "", fmt.Errorf("invalid QUIC port %q", port)
	}
	return host, port, nil
}

func tsNodeQUICAddressValue(address net.Addr) tsValue {
	result := tsNewObject()
	if udp, ok := address.(*net.UDPAddr); ok {
		family := "IPv6"
		if udp.IP.To4() != nil {
			family = "IPv4"
		}
		result.set("address", tsNodeQUICKey(udp.IP.String()))
		result.set("family", tsNodeQUICKey(family))
		result.set("port", tsNumberValue(float64(udp.Port)))
		return tsObjectValue(result)
	}
	host, port, err := net.SplitHostPort(address.String())
	if err == nil {
		portNumber, _ := strconv.Atoi(port)
		result.set("address", tsNodeQUICKey(host))
		result.set("family", tsNodeQUICKey("IPv4"))
		result.set("port", tsNumberValue(float64(portNumber)))
	}
	return tsObjectValue(result)
}

func tsNodeQUICSetEndpointAddress(loop *tsLoop, endpoint tsValue, address net.Addr) {
	state := tsNodeQUICEndpointFor(endpoint)
	if state == nil {
		return
	}
	state.address = tsNodeQUICAddressValue(address)
	tsNodeQUICSet(loop, endpoint, "address", state.address)
}

func tsNodeQUICConnect(loop *tsLoop, address, options tsValue) tsValue {
	if address.kind != tsStringKind && !tsNodeQUICObject(address) {
		return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The \"address\" argument must be a string or SocketAddress"), true))
	}
	host, port, err := tsNodeQUICSplitAddress(loop, address, "443")
	if err != nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", err.Error()), true))
	}
	if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
	tsNodeQUICValidateSessionOptions(loop, options, false)
	remote := net.JoinHostPort(host, port)
	tlsConfig, err := tsNodeQUICMakeTLSConfig(loop, options, false, host)
	if err != nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", err.Error()), true))
	}
	transportConfig := tsNodeQUICTransportConfig(loop, options)
	endpointValue := tsNodeQUICOption(loop, options, "endpoint")
	if tsIsUndefined(endpointValue) {
		wantHTTP := true
		if alpn := tsNodeQUICOption(loop, options, "alpn"); alpn.kind == tsStringKind {
			wantHTTP = strings.HasPrefix(tsText(alpn), "h3")
		}
		reuse := tsNodeQUICOption(loop, options, "reuseEndpoint")
		if tsIsUndefined(reuse) || tsTruthy(reuse) {
			var selected *tsNodeQUICEndpointState
			for _, candidate := range tsNodeQUICSingleton.endpoints {
				if candidate == nil || candidate.closed || candidate.closing || candidate.busy || candidate.listening { continue }
				if candidate.clientHTTPKnown && candidate.clientHTTP != wantHTTP { continue }
				if selected == nil || candidate.creationOrder < selected.creationOrder { selected = candidate }
			}
			if selected != nil { endpointValue = tsObjectValue(selected.object) }
		}
		if tsIsUndefined(endpointValue) { endpointValue = tsNodeQUICCreateEndpoint(loop, tsObjectValue(tsNewObject())) }
	} else if tsNodeQUICEndpointFor(endpointValue) == nil {
		if !tsNodeQUICObject(endpointValue) { tsFSInvalid("options.endpoint", "a QuicEndpoint or an object") }
		endpointValue = tsNodeQUICCreateEndpoint(loop, endpointValue)
	}
	if endpointState := tsNodeQUICEndpointFor(endpointValue); endpointState == nil || endpointState.closed || endpointState.closing || endpointState.busy {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", "QUIC endpoint is unavailable"), true))
	} else if !endpointState.clientHTTPKnown {
		alpn := tsNodeQUICOption(loop, options, "alpn")
		endpointState.clientHTTP = tsIsUndefined(alpn) || alpn.kind != tsStringKind || strings.HasPrefix(tsText(alpn), "h3")
		endpointState.clientHTTPKnown = true
	}
	var remoteIP net.Addr
	if ip := net.ParseIP(host); ip != nil {
		remoteIP = &net.UDPAddr{IP: ip}
	}
	if tsNodeQUICConnectionLimitReached(endpointValue, remoteIP) {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", "QUIC endpoint connection limit reached"), true))
	}
	tsNodeQUICCancelEndpointIdle(tsNodeQUICEndpointFor(endpointValue))
	tsNodeQUICPublish(loop, "onEndpointConnectChannel", tsNodeQUICData(map[string]tsValue{"endpoint":endpointValue,"address":address,"options":options}))
	promise := loop.promise()
	var connection *quicgo.Conn
	var dialError error
	loop.submit(func() tsResult {
		connection, dialError = quicgo.DialAddrEarly(context.Background(), remote, tlsConfig, transportConfig)
		return tsResult{tsU, false}
	}, func(_ tsResult) {
		if dialError != nil {
			promise.settle(tsResult{tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", dialError.Error()), true})
			return
		}
		session := tsNodeQUICCreateSession(loop, endpointValue, options)
		if sessionState := tsNodeQUICSessionFor(session); sessionState != nil {
			sessionState.client = true
			if endpointState := tsNodeQUICEndpointFor(endpointValue); endpointState != nil { endpointState.activeSessions++; tsNodeQUICStatsAdd(endpointState.stats, "clientSessions", 1) }
		}
		tsNodeQUICPublish(loop, "onEndpointClientSessionChannel", tsNodeQUICData(map[string]tsValue{"endpoint":endpointValue,"session":session,"address":tsNodeQUICAddressValue(connection.RemoteAddr()),"options":options}))
		tsNodeQUICStartTransportSession(loop, session, connection)
		promise.settle(tsResult{session, false})
	})
	return tsPromiseValue(promise)
}

func tsNodeQUICListen(loop *tsLoop, callback, options tsValue) tsValue {
	if callback.kind != tsFunctionKind {
		return tsPromiseValue(loop.resolved(tsFSException("TypeError", "ERR_INVALID_ARG_TYPE", "The \"callback\" argument must be of type function"), true))
	}
	if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
	tsNodeQUICValidateSessionOptions(loop, options, true)
	endpoint := tsNodeQUICOption(loop, options, "endpoint")
	if tsIsUndefined(endpoint) {
		endpointOptions := endpoint
		if tsNodeQUICObject(endpointOptions) { endpointOptions = tsNodeQUICCloneProperties(loop, endpointOptions) } else { endpointOptions = tsObjectValue(tsNewObject()) }
		if address := tsNodeQUICOption(loop, options, "address"); !tsIsUndefined(address) {
			tsNodeQUICSet(loop, endpointOptions, "address", address)
		}
		if port := tsNodeQUICOption(loop, options, "port"); !tsIsUndefined(port) {
			tsNodeQUICSet(loop, endpointOptions, "port", port)
		}
		endpoint = tsNodeQUICCreateEndpoint(loop, endpointOptions)
	} else if tsNodeQUICEndpointFor(endpoint) == nil {
		if !tsNodeQUICObject(endpoint) { tsFSInvalid("options.endpoint", "a QuicEndpoint or an object") }
		endpointOptions := tsNodeQUICCloneProperties(loop, endpoint)
		if address := tsNodeQUICOption(loop, options, "address"); !tsIsUndefined(address) { tsNodeQUICSet(loop, endpointOptions, "address", address) }
		if port := tsNodeQUICOption(loop, options, "port"); !tsIsUndefined(port) { tsNodeQUICSet(loop, endpointOptions, "port", port) }
		endpoint = tsNodeQUICCreateEndpoint(loop, endpointOptions)
	}
	endpointState := tsNodeQUICEndpointFor(endpoint)
	if endpointState == nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Invalid QUIC endpoint"), true))
	}
	if endpointState.closed || endpointState.closing || endpointState.busy || endpointState.listening {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Endpoint is closed or already listening"), true))
	}
	tsNodeQUICCancelEndpointIdle(endpointState)
	tlsConfig, err := tsNodeQUICMakeTLSConfig(loop, options, true, "")
	if err != nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", err.Error()), true))
	}
	tsNodeQUICInstallSNICallback(loop, endpoint, options, tlsConfig)
	bindHost, bindPort, err := tsNodeQUICSplitAddress(loop, tsNodeQUICEndpointBindAddress(loop, endpointState), "0")
	if err != nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_ENDPOINT_CLOSED", err.Error()), true))
	}
	listener, err := quicgo.ListenAddr(net.JoinHostPort(bindHost, bindPort), tlsConfig, tsNodeQUICTransportConfig(loop, options))
	if err != nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_ENDPOINT_CLOSED", err.Error()), true))
	}
	state := tsNodeQUICEndpointFor(endpoint)
	state.listening = true
	state.onSession = callback
	tsNodeQUICSetEndpointAddress(loop, endpoint, listener.Addr())
	ctx, cancel := context.WithCancel(context.Background())
	tsNodeQUICTransportInit()
	tsNodeQUICTransportRegistry.endpoints[(*tsObject)(endpoint.ref)] = &tsNodeQUICTransportEndpoint{listener: listener, cancel: cancel}
	loop.resourceClosers = append(loop.resourceClosers, func() { cancel(); _ = listener.Close() })
	tsNodeQUICSet(loop, endpoint, "listening", tsBooleanValue(true))
	tsNodeNetEmitSafe(loop, endpoint, "listening")
	tsNodeQUICPublish(loop, "onEndpointListeningChannel", tsNodeQUICData(map[string]tsValue{"endpoint":endpoint,"options":options}))
	go tsNodeQUICAcceptConnections(loop, endpoint, options, listener, ctx)
	endpointTask := loop.submit(func() tsResult {
		<-ctx.Done()
		return tsResult{tsU, false}
	}, func(_ tsResult) {
		if current := tsNodeQUICEndpointFor(endpoint); current != nil && !current.closed {
			current.listening = false
			if !current.closing { current.closing = true }
			tsNodeQUICFinishEndpoint(loop, endpoint)
		}
	})
	endpointTask.unref = !state.refed
	state.task = endpointTask
	return tsPromiseValue(loop.resolved(endpoint, false))
}

func tsNodeQUICEndpointBindAddress(loop *tsLoop, state *tsNodeQUICEndpointState) tsValue {
	if state == nil {
		return tsU
	}
	address := tsNodeQUICGet(loop, state.options, "address")
	if !tsIsUndefined(address) && !tsIsNull(address) {
		if tsNodeQUICObject(address) && !tsIsUndefined(tsNodeQUICGet(loop, address, "port")) {
			return address
		}
		port := tsNodeQUICGet(loop, state.options, "port")
		if tsIsUndefined(port) || tsIsNull(port) {
			return address
		}
		result := tsNewObject()
		result.set("address", address)
		result.set("port", port)
		return tsObjectValue(result)
	}
	port := tsNodeQUICGet(loop, state.options, "port")
	if tsIsUndefined(port) || tsIsNull(port) {
		port = tsNumberValue(0)
	}
	result := tsNewObject()
	// Bun's BindConfig defaults to loopback IPv4 when no address is supplied.
	result.set("address", tsNodeQUICKey("127.0.0.1"))
	result.set("port", port)
	return tsObjectValue(result)
}

func tsNodeQUICAcceptConnections(loop *tsLoop, endpoint, options tsValue, listener *quicgo.Listener, ctx context.Context) {
	for {
		connection, err := listener.Accept(ctx)
		if err != nil {
			return
		}
		accepted := connection
		handled := make(chan struct{})
		loop.post(func() {
			defer close(handled)
			state := tsNodeQUICEndpointFor(endpoint)
			if state == nil || state.closed || state.closing || state.busy {
				_ = accepted.CloseWithError(0, "endpoint closed")
				return
			}
			if blockList := tsNodeQUICOption(loop, state.options, "blockList"); !tsIsUndefined(blockList) {
				remoteHost, _, _ := net.SplitHostPort(accepted.RemoteAddr().String())
				family := "ipv6"
				if parsed := net.ParseIP(remoteHost); parsed != nil && parsed.To4() != nil { family = "ipv4" }
				check := tsNodeQUICGet(loop, blockList, "check")
				listed := check.kind == tsFunctionKind && tsTruthy(tsCallReceiver(loop, check, blockList, tsNodeQUICKey(remoteHost), tsNodeQUICKey(family)))
				policy := tsNodeQUICGet(loop, state.options, "blockListPolicy")
				blocked := listed
				if tsStrictEqual(policy, tsNodeQUICKey("allow")) { blocked = !listed }
				if blocked {
					tsNodeQUICStatsAdd(state.stats, "packetsBlocked", 1)
					_ = accepted.CloseWithError(0, "address blocked by endpoint policy")
					return
				}
			}
			if tsNodeQUICConnectionLimitReached(endpoint, accepted.RemoteAddr()) {
				tsNodeQUICStatsAdd(state.stats, "serverBusyCount", 1)
				_ = accepted.CloseWithError(0, "endpoint connection limit reached")
				return
			}
			session := tsNodeQUICCreateSession(loop, endpoint, options)
			if endpointState := tsNodeQUICEndpointFor(endpoint); endpointState != nil {
				endpointState.activeSessions++
				tsNodeQUICStatsAdd(endpointState.stats, "serverSessions", 1)
			}
			tsNodeQUICPublish(loop, "onEndpointServerSessionChannel", tsNodeQUICData(map[string]tsValue{"endpoint":endpoint,"session":session,"address":tsNodeQUICAddressValue(accepted.RemoteAddr())}))
			tsNodeQUICStartTransportSession(loop, session, accepted)
			callback := state.onSession
			if callback.kind == tsFunctionKind {
				tsCallReceiver(loop, callback, endpoint, session)
			}
		})
		select {
		case <-handled:
		case <-ctx.Done():
			_ = accepted.CloseWithError(0, "endpoint closed")
			return
		}
	}
}

func tsNodeQUICConnectionLimitReached(endpoint tsValue, remote net.Addr) bool {
	endpointState := tsNodeQUICEndpointFor(endpoint)
	if endpointState == nil {
		return true
	}
	total, host := 0, 0
	remoteHost := ""
	if remote != nil {
		remoteHost, _, _ = net.SplitHostPort(remote.String())
	}
	for _, session := range tsNodeQUICSingleton.sessions {
		if session == nil || session.closed || !tsStrictEqual(session.endpoint, endpoint) {
			continue
		}
		total++
		if transport := tsNodeQUICTransportRegistry.sessions[session.object]; transport != nil && transport.conn != nil {
			peer, _, _ := net.SplitHostPort(transport.conn.RemoteAddr().String())
			if remoteHost != "" && peer == remoteHost {
				host++
			}
		}
	}
	return endpointState.maxTotal > 0 && total >= endpointState.maxTotal || endpointState.maxPerHost > 0 && host >= endpointState.maxPerHost
}

func tsNodeQUICStartTransportSession(loop *tsLoop, session tsValue, connection *quicgo.Conn) {
	tsNodeQUICTransportInit()
	state := tsNodeQUICSessionFor(session)
	if state == nil {
		_ = connection.CloseWithError(0, "invalid session")
		return
	}
	transport := &tsNodeQUICTransportSession{conn: connection, ready: make(chan struct{}), peerMaxFieldSectionSize: ^uint64(0)}
	tsNodeQUICTransportRegistry.sessions[state.object] = transport
	loop.resourceClosers = append(loop.resourceClosers, func() { _ = connection.CloseWithError(0, "event loop closed") })
	if connection.ConnectionState().TLS.HandshakeComplete {
		tsNodeQUICHandshakeReady(loop, session, connection, transport)
	} else {
		go func() {
			<-connection.HandshakeComplete()
			loop.post(func() { tsNodeQUICHandshakeReady(loop, session, connection, transport) })
		}()
	}
	loop.submit(func() tsResult {
		<-connection.Context().Done()
		return tsResult{tsU, false}
	}, func(_ tsResult) {
		tsNodeQUICTransportSessionClosed(loop, session, connection.Context().Err())
	})
}

func tsNodeQUICHandshakeReady(loop *tsLoop, session tsValue, connection *quicgo.Conn, transport *tsNodeQUICTransportSession) {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.closed || connection == nil || transport == nil {
		if transport != nil { transport.readyOnce.Do(func() { close(transport.ready) }) }
		return
	}
	tlsState := connection.ConnectionState().TLS
	if !tlsState.HandshakeComplete {
		transport.readyOnce.Do(func() { close(transport.ready) })
		message := "QUIC TLS handshake failed"
		if err := connection.Context().Err(); err != nil { message = err.Error() }
		tsNodeQUICCloseSession(loop, session, tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", message))
		return
	}
	transport.http3 = strings.HasPrefix(tlsState.NegotiatedProtocol, "h3")
	state.localCertificate = tsNodeQUICLocalCertificate(loop, state.options, tlsState.ServerName)
	tsNodeQUICHandshakeInfo(loop, session, connection)
	go tsNodeQUICAcceptStreams(loop, session, connection)
	go tsNodeQUICAcceptUniStreams(loop, session, connection, transport)
	if transport.http3 {
		tsNodeQUICStartHTTP3(loop, session, transport)
	} else {
		transport.readyOnce.Do(func() { close(transport.ready) })
	}
	connectionState := connection.ConnectionState()
	if connectionState.SupportsDatagrams.Local && connectionState.SupportsDatagrams.Remote {
		go tsNodeQUICReceiveDatagrams(loop, session, connection)
	}
}

func tsNodeQUICAcceptStreams(loop *tsLoop, session tsValue, connection *quicgo.Conn) {
	for {
		wire, err := connection.AcceptStream(connection.Context())
		if err != nil {
			return
		}
		accepted := wire
		loop.post(func() { tsNodeQUICReceiveStream(loop, session, accepted) })
	}
}

func tsNodeQUICAcceptUniStreams(loop *tsLoop, session tsValue, connection *quicgo.Conn, sessionTransport *tsNodeQUICTransportSession) {
	controlSeen, encoderSeen, decoderSeen := false, false, false
	for {
		wire, err := connection.AcceptUniStream(connection.Context())
		if err != nil {
			return
		}
		accepted := wire
		var reader io.Reader = accepted
		if sessionTransport != nil && sessionTransport.http3 {
			streamType, _, readErr := tsNodeQUICReadVarIntFrom(accepted)
			if readErr != nil {
				accepted.CancelRead(0)
				continue
			}
			switch streamType {
			case 0:
				if controlSeen {
					tsNodeQUICHTTP3Close(loop, session, sessionTransport, 0x103, "peer opened a second HTTP/3 control stream")
					accepted.CancelRead(0)
					continue
				}
				controlSeen = true
				go tsNodeQUICReadHTTP3ControlStream(loop, session, sessionTransport, accepted)
				continue
			case 2:
				if encoderSeen {
					tsNodeQUICHTTP3Close(loop, session, sessionTransport, 0x103, "peer opened a second QPACK encoder stream")
					accepted.CancelRead(0)
					continue
				}
				encoderSeen = true
				go tsNodeQUICMonitorQPACKStream(loop, session, sessionTransport, accepted, 0x201)
				continue
			case 3:
				if decoderSeen {
					tsNodeQUICHTTP3Close(loop, session, sessionTransport, 0x103, "peer opened a second QPACK decoder stream")
					accepted.CancelRead(0)
					continue
				}
				decoderSeen = true
				go tsNodeQUICMonitorQPACKStream(loop, session, sessionTransport, accepted, 0x202)
				continue
			case 0x21:
				reader = accepted
			default:
				go tsNodeQUICDiscardUniStream(accepted)
				continue
			}
		}
		streamReader := reader
		loop.post(func() {
			stream := &tsNodeQUICTransportStream{
				reader: streamReader,
				cancelRead: accepted.CancelRead,
				canRead: true,
			}
			tsNodeQUICAddTransportStream(loop, session, int(accepted.StreamID()), "unidirectional", stream)
		})
	}
}

func tsNodeQUICHTTP3Close(loop *tsLoop, session tsValue, transport *tsNodeQUICTransportSession, code uint64, reason string) {
	loop.post(func() {
		if state := tsNodeQUICSessionFor(session); state == nil || state.closed {
			return
		}
		errorValue := tsNodeQUICError("ERR_HTTP3_PROTOCOL_ERROR", reason)
		if callback := tsNodeQUICSessionFor(session).callbacks["onerror"]; callback.kind == tsFunctionKind {
			tsCallReceiver(loop, callback, session, errorValue)
		}
		tsNodeNetEmitSafe(loop, session, "error", errorValue)
		if transport != nil && transport.conn != nil {
			_ = transport.conn.CloseWithError(quicgo.ApplicationErrorCode(code), reason)
		}
	})
}

func tsNodeQUICReadHTTP3ControlStream(loop *tsLoop, session tsValue, transport *tsNodeQUICTransportSession, stream *quicgo.ReceiveStream) {
	settingsSeen := false
	for {
		frameType, _, err := tsNodeQUICReadVarIntFrom(stream)
		if err != nil {
			if sessionTransport := tsNodeQUICSessionTransportState(session); sessionTransport != nil && sessionTransport.conn != nil && sessionTransport.conn.Context().Err() == nil {
				tsNodeQUICHTTP3Close(loop, session, transport, 0x104, "peer closed its HTTP/3 control stream")
			}
			return
		}
		frameLength, _, err := tsNodeQUICReadVarIntFrom(stream)
		if err != nil || frameLength > 1<<20 {
			tsNodeQUICHTTP3Close(loop, session, transport, 0x106, "invalid HTTP/3 control frame length")
			return
		}
		payload := make([]byte, int(frameLength))
		if _, err = io.ReadFull(stream, payload); err != nil {
			tsNodeQUICHTTP3Close(loop, session, transport, 0x106, "truncated HTTP/3 control frame")
			return
		}
		if !settingsSeen && frameType != 0x4 {
			tsNodeQUICHTTP3Close(loop, session, transport, 0x10a, "the first HTTP/3 control frame must be SETTINGS")
			return
		}
		switch frameType {
		case 0x4:
			if settingsSeen {
				tsNodeQUICHTTP3Close(loop, session, transport, 0x109, "peer sent more than one HTTP/3 SETTINGS frame")
				return
			}
			settingsSeen = true
			settings := make(map[uint64]uint64)
			for len(payload) != 0 {
				id, used, complete := tsNodeQUICReadVarInt(payload)
				if !complete {
					tsNodeQUICHTTP3Close(loop, session, transport, 0x109, "malformed HTTP/3 SETTINGS frame")
					return
				}
				payload = payload[used:]
				value, used, complete := tsNodeQUICReadVarInt(payload)
				if !complete {
					tsNodeQUICHTTP3Close(loop, session, transport, 0x109, "malformed HTTP/3 SETTINGS frame")
					return
				}
				payload = payload[used:]
				if _, duplicate := settings[id]; duplicate {
					tsNodeQUICHTTP3Close(loop, session, transport, 0x109, "peer sent a duplicate HTTP/3 setting")
					return
				}
				if (id == 0x8 || id == 0x33) && value > 1 {
					tsNodeQUICHTTP3Close(loop, session, transport, 0x109, "peer sent an invalid boolean HTTP/3 setting")
					return
				}
				settings[id] = value
			}
			loop.post(func() {
				if value, ok := settings[0x6]; ok {
					transport.peerMaxFieldSectionSize = value
				}
				if value, ok := settings[0x33]; ok {
					enabled := value != 0
					transport.peerH3Datagrams = &enabled
				}
			})
		case 0x7:
			lastStreamID, used, complete := tsNodeQUICReadVarInt(payload)
			if !complete || used != len(payload) {
				tsNodeQUICHTTP3Close(loop, session, transport, 0x106, "invalid HTTP/3 GOAWAY frame")
				return
			}
			goawayID := lastStreamID
			loop.post(func() {
				if state := tsNodeQUICSessionFor(session); state != nil && !state.closed {
					state.closing = true
					lastStreamIDValue := tsBigIntUnsigned(goawayID)
					tsNodeQUICPublish(loop, "onSessionClosingChannel", tsNodeQUICData(map[string]tsValue{"session":session}))
					tsNodeQUICPublish(loop, "onSessionGoawayChannel", tsNodeQUICData(map[string]tsValue{"session":session,"lastStreamId":lastStreamIDValue}))
					if callback := state.callbacks["ongoaway"]; callback.kind == tsFunctionKind {
						tsCallReceiver(loop, callback, session, lastStreamIDValue)
					}
				}
			})
		case 0xF0700:
			streamID, used, complete := tsNodeQUICReadVarInt(payload)
			if !complete || streamID&0x3 != 0 || used == len(payload) {
				tsNodeQUICHTTP3Close(loop, session, transport, 0x106, "invalid HTTP/3 PRIORITY_UPDATE frame")
				return
			}
			fieldValue := string(payload[used:])
			urgency := -1
			for _, part := range strings.Split(fieldValue, ",") {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "u=") {
					parsed, parseErr := strconv.Atoi(strings.TrimPrefix(part, "u="))
					if parseErr != nil || parsed < 0 || parsed > 7 { urgency = -1; break }
					urgency = parsed
				}
			}
			if urgency < 0 {
				tsNodeQUICHTTP3Close(loop, session, transport, 0x109, "invalid HTTP/3 priority field value")
				return
			}
		case 0xC:
			origins := make([]tsValue, 0)
			for len(payload) != 0 {
				if len(payload) < 2 {
					tsNodeQUICHTTP3Close(loop, session, transport, 0x106, "invalid HTTP/3 ORIGIN frame")
					return
				}
				length := int(payload[0])<<8 | int(payload[1])
				payload = payload[2:]
				if length == 0 || length > len(payload) {
					tsNodeQUICHTTP3Close(loop, session, transport, 0x106, "invalid HTTP/3 ORIGIN frame")
					return
				}
				origins = append(origins, tsNodeQUICKey(string(payload[:length])))
				payload = payload[length:]
			}
			loop.post(func() {
				if state := tsNodeQUICSessionFor(session); state != nil && !state.closed {
					originValues := tsArrayValue(&tsArray{values: origins})
					tsNodeQUICPublish(loop, "onSessionOriginChannel", tsNodeQUICData(map[string]tsValue{"session":session,"origins":originValues}))
					if callback := state.callbacks["onorigin"]; callback.kind == tsFunctionKind {
						tsCallReceiver(loop, callback, session, originValues)
					}
				}
			})
		}
	}
}

func tsNodeQUICMonitorQPACKStream(loop *tsLoop, session tsValue, transport *tsNodeQUICTransportSession, stream *quicgo.ReceiveStream, errorCode uint64) {
	buffer := make([]byte, 1024)
	for {
		count, err := stream.Read(buffer)
		if count != 0 {
			tsNodeQUICHTTP3Close(loop, session, transport, errorCode, "dynamic QPACK instructions are unsupported")
			return
		}
		if err != nil {
			if transport != nil && transport.conn != nil && transport.conn.Context().Err() == nil {
				tsNodeQUICHTTP3Close(loop, session, transport, 0x104, "peer closed a critical QPACK stream")
			}
			return
		}
	}
}

func tsNodeQUICReceiveStream(loop *tsLoop, session tsValue, wire *quicgo.Stream) {
	sessionTransport := tsNodeQUICSessionTransportState(session)
	transport := &tsNodeQUICTransportStream{
		reader: wire,
		writer: wire,
		tryWrite: wire.TryWriteAll,
		closeWrite: wire.Close,
		cancelRead: wire.CancelRead,
		cancelWrite: wire.CancelWrite,
		canRead: true,
		canWrite: true,
		http3: sessionTransport != nil && sessionTransport.http3,
	}
	tsNodeQUICAddTransportStream(loop, session, int(wire.StreamID()), "bidirectional", transport)
}

func tsNodeQUICReceiveDatagrams(loop *tsLoop, session tsValue, connection *quicgo.Conn) {
	for {
		data, err := connection.ReceiveDatagram(connection.Context())
		if err != nil {
			return
		}
		copyOfData := append([]byte(nil), data...)
		loop.post(func() {
			state := tsNodeQUICSessionFor(session)
			if state == nil || state.closed {
				return
			}
			tsNodeQUICStatsAdd(state.stats, "datagramsReceived", 1)
			if endpointState := tsNodeQUICEndpointFor(state.endpoint); endpointState != nil { tsNodeQUICStatsAdd(endpointState.stats, "bytesReceived", uint64(len(copyOfData))) }
			tsNodeQUICPublish(loop, "onSessionReceiveDatagramChannel", tsNodeQUICData(map[string]tsValue{"length":tsNumberValue(float64(len(copyOfData))),"early":tsBooleanValue(false),"session":session}))
			callback := state.callbacks["ondatagram"]
			if callback.kind == tsFunctionKind {
				tsCallReceiver(loop, callback, session, tsNodeBuffer(copyOfData), tsBooleanValue(false))
			} else {
				tsNodeNetEmitSafe(loop, session, "datagram", tsNodeBuffer(copyOfData), tsBooleanValue(false))
			}
		})
	}
}

func tsNodeQUICAddTransportStream(loop *tsLoop, session tsValue, id int, direction string, transport *tsNodeQUICTransportStream) {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.closed || state.closing {
		if transport.cancelRead != nil {
			transport.cancelRead(0)
		}
		if transport.cancelWrite != nil {
			transport.cancelWrite(0)
		}
		return
	}
	if transport.http3 {
		if sessionState := tsNodeQUICSessionFor(session); sessionState != nil {
			application := tsNodeQUICGet(loop, tsObjectValue(sessionState.object), "applicationOptions")
			transport.http3MaxHeaderLength = tsNodeQUICValueCode(loop, tsNodeQUICGet(loop, application, "maxHeaderLength"), 16384)
			transport.http3MaxHeaderPairs = tsNodeQUICValueCode(loop, tsNodeQUICGet(loop, application, "maxHeaderPairs"), 128)
			transport.http3MaxFieldSectionSize = tsNodeQUICValueCode(loop, tsNodeQUICGet(loop, application, "maxFieldSectionSize"), 0)
		}
	}
	stream := tsNodeQUICCreateStream(loop, session, id, direction)
	streamState := tsNodeQUICStreamFor(stream)
	streamTransportState := tsNodeQUICWireStateFor(stream)
	if streamState == nil || streamTransportState == nil {
		return
	}
	streamTransportState.transport = transport
	if direction == "bidirectional" { tsNodeQUICStatsAdd(state.stats, "bidiInStreamCount", 1) } else { tsNodeQUICStatsAdd(state.stats, "uniInStreamCount", 1) }
	if !transport.canRead {
		streamTransportState.readEnded = true
		if readable := tsNodeReadableStateFor(stream); readable != nil {
			readable.readable = false
			readable.ended = true
			readable.endEmitted = true
			tsNodeReadableSyncState(readable)
		}
	}
	if !transport.canWrite {
		streamTransportState.writeEnded = true
		if writable := tsNodeWritableStateFor(stream); writable != nil {
			writable.writableEnabled = false
			writable.ending = true
			writable.ended = true
			writable.finished = true
			tsNodeWritableSync(writable)
		}
	}
	if sessionState := tsNodeQUICSessionFor(session); sessionState != nil {
		name := "bidiInStreamCount"
		if direction == "unidirectional" { name = "uniInStreamCount" }
		tsNodeQUICStatsAdd(sessionState.stats, name, 1)
	}
	tsNodeQUICStartTransportRead(loop, stream, transport)
	tsNodeQUICScheduleStreamIdleTimeout(loop, stream)
	callback := state.callbacks["onstream"]
	if callback.kind == tsFunctionKind {
		directionName := "bidi"
		if direction == "unidirectional" { directionName = "uni" }
		tsNodeQUICPublish(loop, "onSessionReceivedStreamChannel", tsNodeQUICData(map[string]tsValue{"stream":stream,"session":session,"direction":tsNodeQUICKey(directionName)}))
		tsCallReceiver(loop, callback, session, stream)
	} else {
		warning := tsNodeQUICGet(loop, tsNodeProcessModule(loop), "emitWarning")
		if warning.kind == tsFunctionKind { tsCall(loop, warning, tsNodeQUICKey("A new stream was received but no onstream callback was provided"), tsNodeQUICKey("ExperimentalWarning")) }
		tsNodeQUICEndStream(loop, stream, tsU)
	}
}

func tsNodeQUICStartTransportRead(loop *tsLoop, stream tsValue, transport *tsNodeQUICTransportStream) {
	if !transport.canRead || transport.reader == nil {
		return
	}
	go func() {
		buffer := make([]byte, 32*1024)
		for {
			n, err := transport.reader.Read(buffer)
			if n != 0 {
				data := append([]byte(nil), buffer[:n]...)
				loop.post(func() {
					state := tsNodeQUICStreamFor(stream)
					if state == nil || state.closed {
						return
					}
					tsNodeQUICStatsAdd(state.stats, "bytesReceived", uint64(len(data)))
					tsNodeQUICStatsAdd(state.stats, "bytesAccumulated", uint64(len(data)))
					if sessionState := tsNodeQUICSessionFor(state.session); sessionState != nil {
						tsNodeQUICStatsAdd(sessionState.stats, "bytesReceived", uint64(len(data)))
						if endpointState := tsNodeQUICEndpointFor(sessionState.endpoint); endpointState != nil { tsNodeQUICStatsAdd(endpointState.stats, "bytesReceived", uint64(len(data))) }
					}
					if transport.http3 {
						tsNodeQUICConsumeHTTP3(loop, stream, transport, data)
					} else if readable := tsNodeReadableStateFor(stream); readable != nil {
						tsNodeReadablePush(loop, readable, tsNodeBuffer(data), tsU, false)
					}
				})
			}
			if err != nil {
				loop.post(func() { tsNodeQUICTransportReadEnded(loop, stream, err) })
				return
			}
		}
	}()
}

func tsNodeQUICConsumeHTTP3(loop *tsLoop, stream tsValue, transport *tsNodeQUICTransportStream, data []byte) {
	transport.http3Buffer = append(transport.http3Buffer, data...)
	for {
		if !transport.http3FrameActive {
			frameType, typeBytes, complete := tsNodeQUICReadVarInt(transport.http3Buffer)
			if !complete { return }
			frameLength, lengthBytes, complete := tsNodeQUICReadVarInt(transport.http3Buffer[typeBytes:])
			if !complete { return }
			transport.http3Buffer = transport.http3Buffer[typeBytes+lengthBytes:]
			transport.http3FrameType, transport.http3FrameRemaining, transport.http3FrameActive = frameType, frameLength, true
			transport.http3HeaderBlock = transport.http3HeaderBlock[:0]
			if frameType == 0x1 && frameLength > transport.http3MaxHeaderLength {
				tsNodeQUICEndStream(loop, stream, tsNodeQUICError("ERR_HTTP3_EXCESSIVE_LOAD", "HTTP/3 header section exceeds the configured limit"))
				return
			}
		}
		if transport.http3FrameRemaining == 0 {
			if transport.http3FrameType == 0x1 { tsNodeQUICReceiveHTTP3Headers(loop, stream, transport, transport.http3HeaderBlock) }
			transport.http3FrameActive = false
			transport.http3HeaderBlock = nil
			continue
		}
		if len(transport.http3Buffer) == 0 { return }
		take := uint64(len(transport.http3Buffer))
		if take > transport.http3FrameRemaining { take = transport.http3FrameRemaining }
		count := int(take)
		payload := transport.http3Buffer[:count]
		switch transport.http3FrameType {
		case 0x0:
			if readable := tsNodeReadableStateFor(stream); readable != nil && count != 0 { tsNodeReadablePush(loop, readable, tsNodeBuffer(append([]byte(nil), payload...)), tsU, false) }
		case 0x1:
			transport.http3HeaderBlock = append(transport.http3HeaderBlock, payload...)
		}
		transport.http3Buffer = transport.http3Buffer[count:]
		transport.http3FrameRemaining -= take
	}
}

func tsNodeQUICReceiveHTTP3Headers(loop *tsLoop, stream tsValue, transport *tsNodeQUICTransportStream, block []byte) {
	decoder := quicqpack.NewDecoder().Decode(block)
	headers := tsNewObject()
	fieldCount, fieldSectionSize := uint64(0), uint64(0)
	for {
		field, err := decoder()
		if err == io.EOF {
			break
		}
		if err != nil {
		tsNodeQUICHTTP3Close(loop, stateSession(stream), tsNodeQUICSessionTransportState(stateSession(stream)), 0x200, err.Error())
			return
		}
		fieldCount++
		fieldSectionSize += uint64(len(field.Name) + len(field.Value) + 32)
		if fieldCount > transport.http3MaxHeaderPairs || (transport.http3MaxFieldSectionSize != 0 && fieldSectionSize > transport.http3MaxFieldSectionSize) {
			tsNodeQUICHTTP3Close(loop, stateSession(stream), tsNodeQUICSessionTransportState(stateSession(stream)), 0x107, "HTTP/3 header section exceeds the configured limit")
			return
		}
		if previous, exists := headers.values[field.Name]; exists {
			if previous.kind == tsArrayKind {
				array := (*tsArray)(previous.ref)
				array.values = append(array.values, tsNodeQUICKey(field.Value))
			} else {
				headers.set(field.Name, tsArrayValue(&tsArray{values: []tsValue{previous, tsNodeQUICKey(field.Value)}}))
			}
		} else {
			headers.set(field.Name, tsNodeQUICKey(field.Value))
		}
	}
	streamState := tsNodeQUICStreamFor(stream)
	if streamState == nil {
		return
	}
	callbackName := "onheaders"
	if !tsIsUndefined(streamState.headers) && !tsIsNull(streamState.headers) {
		callbackName = "ontrailers"
		streamState.pendingTrailers = tsObjectValue(headers)
	} else {
		status := tsText(tsNodeQUICGet(loop, tsObjectValue(headers), ":status"))
		if len(status) == 3 && status[0] == '1' {
			callbackName = "oninfo"
		} else {
			streamState.headers = tsObjectValue(headers)
		}
	}
	callback := streamState.callbacks[callbackName]
	if callback.kind != tsFunctionKind {
		if session := tsNodeQUICSessionFor(streamState.session); session != nil {
			callback = session.callbacks[callbackName]
		}
	}
	channel, field := "onStreamHeadersChannel", "headers"
	if callbackName == "ontrailers" { channel, field = "onStreamTrailersChannel", "trailers" }
	if callbackName == "oninfo" { channel, field = "onStreamInfoChannel", "headers" }
	if callbackName != "onheaders" || callback.kind == tsFunctionKind {
		tsNodeQUICPublish(loop, channel, tsNodeQUICData(map[string]tsValue{"stream":stream,"session":streamState.session,field:tsObjectValue(headers)}))
	}
	if callback.kind == tsFunctionKind { tsCallReceiver(loop, callback, stream, tsObjectValue(headers)) }
}

func tsNodeQUICTransportReadEnded(loop *tsLoop, stream tsValue, err error) {
	state := tsNodeQUICStreamFor(stream)
	if state == nil || state.closed {
		return
	}
	if readable := tsNodeReadableStateFor(stream); readable != nil {
		tsNodeReadablePush(loop, readable, tsNull, tsU, false)
	}
	transportState := tsNodeQUICWireStateFor(stream)
	if transportState != nil {
		transportState.readEnded = true
	}
	if err != nil && err != io.EOF {
		errorValue := tsNodeQUICError("ERR_QUIC_STREAM_RESET", err.Error())
		tsNodeQUICPublish(loop, "onStreamResetChannel", tsNodeQUICData(map[string]tsValue{"stream":stream,"session":state.session,"error":errorValue}))
		if streamState := tsNodeQUICStreamFor(stream); streamState != nil {
			if callback := streamState.callbacks["onreset"]; callback.kind == tsFunctionKind { tsCallReceiver(loop, callback, stream, errorValue) }
			if callback := streamState.callbacks["onerror"]; callback.kind == tsFunctionKind { tsCallReceiver(loop, callback, stream, errorValue) }
		}
		tsNodeNetEmitSafe(loop, stream, "error", errorValue)
	}
	tsNodeQUICMaybeCloseTransportStream(loop, stream)
}

func tsNodeQUICMaybeCloseTransportStream(loop *tsLoop, stream tsValue) {
	state := tsNodeQUICStreamFor(stream)
	transportState := tsNodeQUICWireStateFor(stream)
	if state == nil || transportState == nil || state.closed {
		return
	}
	if transportState.readEnded && transportState.writeEnded {
		session := state.session
		state.closed = true
		state.destroyed = true
		if state.idleTimer != nil { state.idleTimer.Stop(); state.idleTimer = nil }
		tsNodeQUICCloseFileHandle(loop, state)
		tsNodeQUICStatsSet(state.stats, "destroyedAt", time.Now().UnixNano())
		if state.closedPromise != nil { state.closedPromise.settle(tsResult{tsU, false}) }
		tsNodeNetEmitSafe(loop, stream, "close")
		tsNodeQUICPublish(loop, "onStreamClosedChannel", tsNodeQUICData(map[string]tsValue{"stream":stream,"session":session,"error":tsU,"stats":state.stats}))
		state.session = tsNull
		if sessionState := tsNodeQUICSessionFor(session); sessionState != nil && sessionState.graceful { tsNodeQUICMaybeFinishSession(loop, session) }
	}
}

func tsNodeQUICTransportSessionClosed(loop *tsLoop, session tsValue, err error) {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.closed {
		return
	}
	var errorValue tsValue = tsU
	if err != nil && err != context.Canceled {
		errorValue = tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", err.Error())
	}
	tsNodeQUICCloseSession(loop, session, errorValue)
}

func tsNodeQUICOpenStream(loop *tsLoop, session tsValue, direction string, options tsValue) tsValue {
	state := tsNodeQUICSessionFor(session)
	transport := tsNodeQUICSessionTransportState(session)
	if state == nil || state.closed || state.closing || transport == nil || transport.conn == nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_OPEN_STREAM_FAILED", "Failed to open QUIC stream"), true))
	}
	if tsIsUndefined(options) { options = tsObjectValue(tsNewObject()) }
	if !tsNodeQUICObject(options) { tsFSInvalid("options", "an object") }
	tsNodeQUICValidateOpenStreamOptions(loop, session, options)
	promise := loop.promise()
	var bidirectional *quicgo.Stream
	var unidirectional *quicgo.SendStream
	var openError error
	connection := transport.conn
	http3 := transport.http3
	ready := transport.ready
	loop.submit(func() tsResult {
		if ready != nil {
			<-ready
		}
		if direction == "bidirectional" {
			bidirectional, openError = connection.OpenStreamSync(connection.Context())
		} else {
			unidirectional, openError = connection.OpenUniStreamSync(connection.Context())
			if openError == nil && http3 {
				_, openError = unidirectional.Write(tsNodeQUICAppendVarInt(nil, 0x21))
			}
		}
		return tsResult{tsU, false}
	}, func(_ tsResult) {
		if openError != nil {
			promise.settle(tsResult{tsNodeQUICError("ERR_QUIC_OPEN_STREAM_FAILED", openError.Error()), true})
			return
		}
		var id int
		wire := &tsNodeQUICTransportStream{canWrite: true}
		if bidirectional != nil {
			id = int(bidirectional.StreamID())
			wire.reader = bidirectional
			wire.writer = bidirectional
			wire.tryWrite = bidirectional.TryWriteAll
			wire.closeWrite = bidirectional.Close
			wire.cancelRead = bidirectional.CancelRead
			wire.cancelWrite = bidirectional.CancelWrite
			wire.canRead = true
		} else {
			id = int(unidirectional.StreamID())
			wire.writer = unidirectional
			wire.tryWrite = unidirectional.TryWriteAll
			wire.closeWrite = unidirectional.Close
			wire.cancelWrite = unidirectional.CancelWrite
		}
		wire.http3 = http3 && direction == "bidirectional"
		stream := tsNodeQUICAddOpenedTransportStream(loop, session, id, direction, wire)
		tsNodeQUICApplyOpenStreamOptions(loop, session, stream, options)
		directionName := "bidi"
		if direction == "unidirectional" { directionName = "uni" }
		tsNodeQUICPublish(loop, "onSessionOpenStreamChannel", tsNodeQUICData(map[string]tsValue{"stream":stream,"session":session,"direction":tsNodeQUICKey(directionName)}))
		promise.settle(tsResult{stream, false})
	})
	return tsPromiseValue(promise)
}

func tsNodeQUICValidateOpenStreamOptions(loop *tsLoop, session, options tsValue) {
	if value := tsNodeQUICOption(loop, options, "highWaterMark"); !tsIsUndefined(value) {
		if !tsIsNumeric(value) || tsNumber(value) < 0 || tsNumber(value) > 4294967295 || math.Trunc(tsNumber(value)) != tsNumber(value) { tsFSInvalid("options.highWaterMark", "an integer from 0 to 4294967295") }
	}
	priority := tsNodeQUICOption(loop, options, "priority")
	if !tsIsUndefined(priority) {
		name := tsText(priority)
		if name != "default" && name != "low" && name != "high" { tsFSInvalid("options.priority", "'default', 'low', or 'high'") }
	}
	if value := tsNodeQUICOption(loop, options, "incremental"); !tsIsUndefined(value) && value.kind != tsBooleanKind { tsFSInvalid("options.incremental", "a boolean") }
	for _, name := range []string{"onheaders", "ontrailers", "oninfo", "onwanttrailers"} {
		if callback := tsNodeQUICOption(loop, options, name); tsTruthy(callback) {
			if callback.kind != tsFunctionKind { tsFSInvalid("options."+name, "a function") }
			if !tsNodeQUICHeadersSupported(session) { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "The negotiated QUIC application protocol does not support headers")}) }
		}
	}
	if headers := tsNodeQUICOption(loop, options, "headers"); !tsIsUndefined(headers) {
		if !tsNodeQUICObject(headers) { tsFSInvalid("headers", "an object") }
		tsNodeQUICValidateHeaders(headers)
	}
	if body := tsNodeQUICOption(loop, options, "body"); !tsIsUndefined(body) {
		valid := body.kind == tsStringKind || body.kind == tsArrayBufferKind || body.kind == tsTypedArrayKind || body.kind == tsDataViewKind || body.kind == tsPromiseKind || tsIsNull(body) || tsNodeQUICIsFileHandle(body) || (body.kind == tsObjectKind && (*tsObject)(body.ref).blob != nil) || tsNodeIterFromIsSyncIterable(loop, body) || tsNodeIterFromIsAsyncIterable(loop, body)
		if !valid { tsNodeIterFromInvalid("body", "a string, binary view, Blob, FileHandle, Promise, or iterable", body) }
	}
}

func tsNodeQUICApplyOpenStreamOptions(loop *tsLoop, session, stream, options tsValue) {
	streamState := tsNodeQUICStreamFor(stream)
	if streamState == nil { return }
	if value := tsNodeQUICOption(loop, options, "highWaterMark"); !tsIsUndefined(value) {
		if !tsIsNumeric(value) || tsNumber(value) < 0 || tsNumber(value) > 4294967295 || math.Trunc(tsNumber(value)) != tsNumber(value) { tsFSInvalid("options.highWaterMark", "an integer from 0 to 4294967295") }
		streamState.highWaterMark = int(tsNumber(value))
	}
	priority := tsNodeQUICOption(loop, options, "priority")
	if tsIsUndefined(priority) { priority = tsNodeQUICKey("default") }
	if !tsIsUndefined(priority) {
		name := tsText(priority)
		if name != "default" && name != "low" && name != "high" { tsFSInvalid("options.priority", "'default', 'low', or 'high'") }
	}
	incremental := tsNodeQUICOption(loop, options, "incremental")
	if tsIsUndefined(incremental) { incremental = tsBooleanValue(false) } else if incremental.kind != tsBooleanKind { tsFSInvalid("options.incremental", "a boolean") }
	priorityValue := tsNewObject()
	priorityValue.set("level", priority)
	priorityValue.set("incremental", incremental)
	streamState.priority = tsObjectValue(priorityValue)
	urgency := uint64(3)
	switch tsText(priority) { case "high": urgency = 0; case "low": urgency = 7 }
	if urgency != 3 || tsTruthy(incremental) { tsNodeQUICSendPriorityUpdate(stream, urgency, tsTruthy(incremental)) }
	for _, name := range []string{"onheaders", "ontrailers", "oninfo", "onwanttrailers"} {
		if callback := tsNodeQUICOption(loop, options, name); tsTruthy(callback) {
			if callback.kind != tsFunctionKind { tsFSInvalid(name, "a function") }
			streamState.callbacks[name] = callback
		}
	}
	body := tsNodeQUICOption(loop, options, "body")
	if headers := tsNodeQUICOption(loop, options, "headers"); !tsIsUndefined(headers) {
		tsNodeQUICValidateHeaders(headers)
		tsNodeQUICSendHeaders(loop, stream, headers, tsIsUndefined(body), "initial")
	}
	if !tsIsUndefined(body) { tsNodeQUICSetOutbound(loop, stream, body) }
	_ = session
}

func tsNodeQUICAddOpenedTransportStream(loop *tsLoop, session tsValue, id int, direction string, transport *tsNodeQUICTransportStream) tsValue {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.closed {
		return tsU
	}
	stream := tsNodeQUICCreateStream(loop, session, id, direction)
	if direction == "bidirectional" { tsNodeQUICStatsAdd(state.stats, "bidiOutStreamCount", 1) } else { tsNodeQUICStatsAdd(state.stats, "uniOutStreamCount", 1) }
	if streamState := tsNodeQUICWireStateFor(stream); streamState != nil {
		streamState.transport = transport
	}
	if !transport.canRead {
		if state := tsNodeQUICWireStateFor(stream); state != nil {
			state.readEnded = true
		}
		if readable := tsNodeReadableStateFor(stream); readable != nil {
			readable.readable = false
			readable.ended = true
			readable.endEmitted = true
			tsNodeReadableSyncState(readable)
		}
	}
	if !transport.canWrite {
		if state := tsNodeQUICWireStateFor(stream); state != nil {
			state.writeEnded = true
		}
		if writable := tsNodeWritableStateFor(stream); writable != nil {
			writable.writableEnabled = false
			writable.ending = true
			writable.ended = true
			writable.finished = true
			tsNodeWritableSync(writable)
		}
	}
	tsNodeQUICStartTransportRead(loop, stream, transport)
	tsNodeQUICScheduleStreamIdleTimeout(loop, stream)
	return stream
}

func tsNodeQUICScheduleStreamIdleTimeout(loop *tsLoop, stream tsValue) {
	streamState := tsNodeQUICStreamFor(stream)
	if streamState == nil { return }
	sessionState := tsNodeQUICSessionFor(streamState.session)
	if sessionState == nil || sessionState.streamIdleTimeout <= 0 { return }
	streamState.idleTimer = time.AfterFunc(sessionState.streamIdleTimeout, func() {
		loop.post(func() {
			if current := tsNodeQUICStreamFor(stream); current != nil && !current.closed {
				tsNodeQUICEndStream(loop, stream, tsNodeQUICError("ERR_QUIC_TRANSPORT_ERROR", "QUIC transport error: stream idle timeout"))
			}
		})
	})
}

type tsNodeQUICStreamWireState struct {
	transport *tsNodeQUICTransportStream
	readEnded, writeEnded bool
}

var tsNodeQUICStreamTransportStates map[uintptr]*tsNodeQUICStreamWireState

func tsNodeQUICStreamTransportStateForInit() {
	if tsNodeQUICStreamTransportStates == nil {
		tsNodeQUICStreamTransportStates = map[uintptr]*tsNodeQUICStreamWireState{}
	}
}

func tsNodeQUICWireStateFor(stream tsValue) *tsNodeQUICStreamWireState {
	if stream.kind != tsObjectKind {
		return nil
	}
	tsNodeQUICStreamTransportStateForInit()
	return tsNodeQUICStreamTransportStates[uintptr(stream.ref)]
}

func tsNodeQUICSessionTransportState(session tsValue) *tsNodeQUICTransportSession {
	if session.kind != tsObjectKind {
		return nil
	}
	tsNodeQUICTransportInit()
	return tsNodeQUICTransportRegistry.sessions[(*tsObject)(session.ref)]
}

func tsNodeQUICHeadersSupported(session tsValue) bool {
	transport := tsNodeQUICSessionTransportState(session)
	return transport != nil && transport.http3
}

func tsNodeQUICWriteTransport(loop *tsLoop, stream tsValue, chunk, encoding, callback tsValue) {
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil || !state.transport.canWrite || state.transport.writer == nil {
		if callback.kind == tsFunctionKind {
			tsCall(loop, callback, tsNodeQUICError("ERR_QUIC_STREAM_ABORTED", "The QUIC stream is not writable"))
		}
		return
	}
	data := tsNodeQUICInputBytes(loop, chunk)
	if chunk.kind == tsStringKind {
		data = append([]byte(nil), tsNodeBufferVisibleBytes(tsNodeBufferFrom(loop, []tsValue{chunk, encoding}))...)
	}
	data = append([]byte(nil), data...)
	count := len(data)
	data = tsNodeQUICEncodeTransportPayload(stream, state.transport, data)
	tsNodeQUICEnqueueTransportBytes(loop, stream, data, count, callback, false)
}

func tsNodeQUICFinalizeTransportWrite(loop *tsLoop, stream tsValue, callback tsValue) {
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil || !state.transport.canWrite {
		if callback.kind == tsFunctionKind {
			tsCall(loop, callback, tsU)
		}
		return
	}
	if state.transport.http3 && state.transport.headersSent {
		if streamState := tsNodeQUICStreamFor(stream); streamState != nil {
			if !tsIsUndefined(streamState.pendingTrailers) && !tsIsNull(streamState.pendingTrailers) {
				trailers := streamState.pendingTrailers
				streamState.pendingTrailers = tsU
				tsNodeQUICSendHeaders(loop, stream, trailers, true, "trailing")
				if callback.kind == tsFunctionKind { tsCall(loop, callback, tsU) }
				return
			}
			if wantTrailers := streamState.callbacks["onwanttrailers"]; wantTrailers.kind == tsFunctionKind {
				tsCallReceiver(loop, wantTrailers, stream)
			}
			if !tsIsUndefined(streamState.pendingTrailers) && !tsIsNull(streamState.pendingTrailers) {
				trailers := streamState.pendingTrailers
				streamState.pendingTrailers = tsU
				tsNodeQUICSendHeaders(loop, stream, trailers, true, "trailing")
				if callback.kind == tsFunctionKind { tsCall(loop, callback, tsU) }
				return
			}
		}
	}
	tsNodeQUICEnqueueTransportBytes(loop, stream, nil, 0, callback, true)
}

func tsNodeQUICEncodeTransportPayload(stream tsValue, transport *tsNodeQUICTransportStream, data []byte) []byte {
	if transport == nil || !transport.http3 || !transport.headersSent {
		return data
	}
	frame := tsNodeQUICAppendVarInt(nil, 0x0)
	frame = tsNodeQUICAppendVarInt(frame, uint64(len(data)))
	return append(frame, data...)
}

func tsNodeQUICSyncWrite(loop *tsLoop, stream tsValue, data []byte) tsValue {
	state := tsNodeQUICWireStateFor(stream)
	streamState := tsNodeQUICStreamFor(stream)
	if state == nil || state.transport == nil || state.transport.tryWrite == nil || !state.transport.canWrite || state.writeEnded || state.transport.writeBusy || len(state.transport.writeQueue) != 0 || streamState == nil || (len(data) != 0 && tsNodeQUICWriterAvailable(streamState, state.transport) == 0) {
		return tsBooleanValue(false)
	}
	encoded := tsNodeQUICEncodeTransportPayload(stream, state.transport, data)
	if err := state.transport.tryWrite(encoded); err != nil {
		return tsBooleanValue(false)
	}
	state.transport.totalBytes += int64(len(data))
	tsNodeQUICRecordSent(stream, uint64(len(data)))
	return tsBooleanValue(true)
}

func tsNodeQUICWriterAvailable(state *tsNodeQUICStreamState, transport *tsNodeQUICTransportStream) int {
	if state == nil { return 0 }
	queued := 0
	if transport != nil { queued = transport.queuedBytes }
	available := state.highWaterMark - queued
	if available < 0 { return 0 }
	return available
}

func tsNodeQUICRecordSent(stream tsValue, count uint64) {
	streamState := tsNodeQUICStreamFor(stream)
	if streamState == nil { return }
	tsNodeQUICStatsAdd(streamState.stats, "bytesSent", count)
	if sessionState := tsNodeQUICSessionFor(streamState.session); sessionState != nil {
		tsNodeQUICStatsAdd(sessionState.stats, "bytesSent", count)
		if endpointState := tsNodeQUICEndpointFor(sessionState.endpoint); endpointState != nil { tsNodeQUICStatsAdd(endpointState.stats, "bytesSent", count) }
	}
}

func tsNodeQUICWakeWriter(loop *tsLoop, state *tsNodeQUICStreamState) {
	if state == nil || len(state.drains) == 0 { return }
	wire := tsNodeQUICWireStateFor(tsObjectValue(state.object))
	if wire == nil || wire.transport == nil || tsNodeQUICWriterAvailable(state, wire.transport) == 0 { return }
	for _, pending := range state.drains { pending.settle(tsResult{tsBooleanValue(true), false}) }
	state.drains = nil
	_ = loop
}

func tsNodeQUICSyncEnd(loop *tsLoop, stream tsValue) tsValue {
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil || !state.transport.canWrite {
		return tsNumberValue(-1)
	}
	transport := state.transport
	if transport.writeBusy || len(transport.writeQueue) != 0 {
		return tsNumberValue(-1)
	}
	if !state.writeEnded && transport.closeWrite != nil {
		if err := transport.closeWrite(); err != nil {
			return tsNumberValue(-1)
		}
	}
	state.writeEnded = true
	tsNodeQUICMaybeCloseTransportStream(loop, stream)
	return tsNumberValue(float64(transport.totalBytes))
}

func tsNodeQUICValueCode(loop *tsLoop, value tsValue, fallback uint64) uint64 {
	if value.kind == tsBigIntKind {
		integer := (*big.Int)(value.ref)
		if integer.Sign() >= 0 && integer.IsUint64() { return integer.Uint64() }
	}
	if tsIsNumeric(value) && tsNumber(value) >= 0 && math.Trunc(tsNumber(value)) == tsNumber(value) { return uint64(tsNumber(value)) }
	if !tsNodeQUICObject(value) { return fallback }
	code := tsNodeQUICGet(loop, value, "errorCode")
	if code.kind == tsBigIntKind {
		integer := (*big.Int)(code.ref)
		if integer.Sign() >= 0 && integer.IsUint64() { return integer.Uint64() }
	}
	if tsIsNumeric(code) && tsNumber(code) >= 0 { return uint64(tsNumber(code)) }
	return fallback
}

func tsNodeQUICValidateWireCode(value tsValue, name string) uint64 {
	if value.kind == tsBigIntKind {
		integer := (*big.Int)(value.ref)
		if integer.Sign() >= 0 && integer.BitLen() <= 62 { return integer.Uint64() }
	} else if tsIsNumeric(value) {
		number := tsNumber(value)
		if !math.IsNaN(number) && number >= 0 && number < math.Ldexp(1, 62) && math.Trunc(number) == number { return uint64(number) }
	}
	tsFSInvalid(name, "a non-negative 62-bit integer")
	return 0
}

func tsNodeQUICStreamErrorCode(loop *tsLoop, value tsValue, fallback uint64) uint64 {
	if value.kind != tsObjectKind || !tsNodeQUICSingleton.quicErrors[(*tsObject)(value.ref)] { return fallback }
	if tsText(tsNodeQUICGet(loop, value, "type")) != "application" { return fallback }
	return tsNodeQUICValidateWireCode(tsNodeQUICGet(loop, value, "errorCode"), "errorCode")
}

func tsNodeQUICResetTransportStream(loop *tsLoop, stream tsValue, errorValue tsValue, override ...uint64) {
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil {
		return
	}
	code := tsNodeQUICStreamErrorCode(loop, errorValue, 0)
	if len(override) != 0 { code = override[0] }
	if state.transport.cancelRead != nil {
		state.transport.cancelRead(quicgo.StreamErrorCode(code))
	}
	if state.transport.cancelWrite != nil {
		state.transport.cancelWrite(quicgo.StreamErrorCode(code))
	}
	state.readEnded = true
	state.writeEnded = true
	_ = loop
}

func tsNodeQUICStopSending(loop *tsLoop, stream, codeValue tsValue) tsValue {
	if tsIsUndefined(codeValue) { codeValue = tsBigIntUnsigned(0) }
	code := tsNodeQUICValidateWireCode(codeValue, "code")
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil || state.readEnded { return tsU }
	if state.transport.cancelRead != nil { state.transport.cancelRead(quicgo.StreamErrorCode(code)) }
	state.readEnded = true
	if streamState := tsNodeQUICStreamFor(stream); streamState != nil { tsNodeQUICMaybeCloseTransportStream(loop, stream) }
	return tsU
}

func tsNodeQUICResetWriteSide(loop *tsLoop, stream, codeValue tsValue) tsValue {
	if tsIsUndefined(codeValue) { codeValue = tsBigIntUnsigned(0) }
	code := tsNodeQUICValidateWireCode(codeValue, "code")
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil || state.writeEnded { return tsU }
	if state.transport.cancelWrite != nil { state.transport.cancelWrite(quicgo.StreamErrorCode(code)) }
	state.writeEnded = true
	if streamState := tsNodeQUICStreamFor(stream); streamState != nil { tsNodeQUICMaybeCloseTransportStream(loop, stream) }
	return tsU
}

func tsNodeQUICCloseTransportSession(session tsValue, application bool, code uint64, reason string) {
	if state := tsNodeQUICSessionTransportState(session); state != nil && state.conn != nil {
		if application {
			_ = state.conn.CloseWithError(quicgo.ApplicationErrorCode(code), reason)
		} else {
			_ = state.conn.CloseWithTransportError(quicgo.TransportErrorCode(code), reason)
		}
	}
}

func tsNodeQUICErrorCloseOptions(loop *tsLoop, value tsValue) tsValue {
	reasonValue := value
	if tsNodeQUICObject(value) {
		if message := tsNodeQUICGet(loop, value, "message"); tsTruthy(message) {
			reasonValue = message
		}
	}
	reason := tsText(tsStringValueWithLoop(loop, reasonValue))
	options := tsNewObject()
	options.set("code", tsBigIntUnsigned(1))
	options.set("type", tsNodeQUICKey("transport"))
	options.set("reason", tsNodeQUICKey(reason))
	return tsObjectValue(options)
}

func tsNodeQUICCloseTransportEndpoint(endpoint tsValue) {
	if endpoint.kind != tsObjectKind {
		return
	}
	tsNodeQUICTransportInit()
	key := (*tsObject)(endpoint.ref)
	if transport := tsNodeQUICTransportRegistry.endpoints[key]; transport != nil {
		if transport.cancel != nil {
			transport.cancel()
		}
		if transport.listener != nil {
			_ = transport.listener.Close()
		}
		delete(tsNodeQUICTransportRegistry.endpoints, key)
	}
}

func tsNodeQUICStartHTTP3(loop *tsLoop, session tsValue, transport *tsNodeQUICTransportSession) {
	application := tsNodeQUICGet(loop, session, "applicationOptions")
	fieldSectionSize := tsNodeQUICValueCode(loop, tsNodeQUICGet(loop, application, "maxFieldSectionSize"), 0)
	connectProtocol := tsTruthy(tsNodeQUICGet(loop, application, "enableConnectProtocol"))
	datagramsEnabled := tsTruthy(tsNodeQUICGet(loop, application, "enableDatagrams"))
	origins := tsNodeQUICServerOrigins(loop, session)
	go func() {
		var openError error
		control, err := transport.conn.OpenUniStreamSync(transport.conn.Context())
		if err == nil {
			transport.controlMu.Lock()
			transport.control = control
			transport.controlMu.Unlock()
			settings := append([]byte{}, tsNodeQUICAppendVarInt(nil, 0x1)...)
			settings = append(settings, tsNodeQUICAppendVarInt(nil, 0)...)
			if fieldSectionSize > 0 { settings = append(settings, tsNodeQUICAppendVarInt(nil, 0x6)...); settings = append(settings, tsNodeQUICAppendVarInt(nil, fieldSectionSize)...) }
			settings = append(settings, tsNodeQUICAppendVarInt(nil, 0x7)...)
			settings = append(settings, tsNodeQUICAppendVarInt(nil, 0)...)
			if connectProtocol { settings = append(settings, tsNodeQUICAppendVarInt(nil, 0x8)...); settings = append(settings, tsNodeQUICAppendVarInt(nil, 1)...) }
			if datagramsEnabled && transport.conn.ConnectionState().SupportsDatagrams.Local {
				settings = append(settings, tsNodeQUICAppendVarInt(nil, 0x33)...)
				settings = append(settings, tsNodeQUICAppendVarInt(nil, 1)...)
			}
			frame := append(tsNodeQUICAppendVarInt(nil, 0x4), tsNodeQUICAppendVarInt(nil, uint64(len(settings)))...)
			frame = append(frame, settings...)
			_, openError = control.Write(append(tsNodeQUICAppendVarInt(nil, 0), frame...))
			if openError == nil && len(origins) != 0 {
				payload := make([]byte, 0)
				for _, origin := range origins {
					payload = append(payload, byte(len(origin)>>8), byte(len(origin)))
					payload = append(payload, origin...)
				}
				originFrame := tsNodeQUICAppendVarInt(nil, 0xC)
				originFrame = tsNodeQUICAppendVarInt(originFrame, uint64(len(payload)))
				originFrame = append(originFrame, payload...)
				_, openError = control.Write(originFrame)
			}
		}
		if err == nil && openError == nil {
			var encoder, decoder *quicgo.SendStream
			encoder, err = transport.conn.OpenUniStreamSync(transport.conn.Context())
			if err == nil {
				_, openError = encoder.Write(tsNodeQUICAppendVarInt(nil, 0x2))
			}
			if err == nil && openError == nil {
				decoder, err = transport.conn.OpenUniStreamSync(transport.conn.Context())
				if err == nil {
					_, openError = decoder.Write(tsNodeQUICAppendVarInt(nil, 0x3))
				}
			}
		}
		transport.readyOnce.Do(func() { close(transport.ready) })
		if err != nil || openError != nil {
			problem := err
			if problem == nil {
				problem = openError
			}
			loop.post(func() {
				if state := tsNodeQUICSessionFor(session); state != nil && !state.closed {
					tsNodeNetEmitSafe(loop, session, "error", tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", problem.Error()))
				}
			})
		}
	}()
}

func tsNodeQUICServerOrigins(loop *tsLoop, session tsValue) []string {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.client {
		return nil
	}
	tlsOptions := tsNodeQUICTLSOptions(loop, state.options)
	identities := tsNodeQUICOption(loop, tlsOptions, "sni")
	if !tsNodeQUICObject(identities) {
		return nil
	}
	result := make([]string, 0)
	for _, hostname := range (*tsObject)(identities.ref).order {
		if hostname == "*" || strings.Contains(hostname, "*") {
			continue
		}
		identity := (*tsObject)(identities.ref).values[hostname]
		if tsStrictEqual(tsNodeQUICOption(loop, identity, "authoritative"), tsBooleanValue(false)) {
			continue
		}
		origin := "https://" + hostname
		port := tsNodeQUICOption(loop, identity, "port")
		if !tsIsUndefined(port) && !(tsIsNumeric(port) && tsNumber(port) == 443) {
			origin += ":" + tsText(port)
		}
		if len(origin) == 0 || len(origin) > 65535 {
			continue
		}
		result = append(result, origin)
	}
	return result
}

func tsNodeQUICBeginGracefulTransportSession(session tsValue) {
	state := tsNodeQUICSessionFor(session)
	transport := tsNodeQUICSessionTransportState(session)
	if state == nil || state.client || transport == nil || !transport.http3 { return }
	transport.controlMu.Lock()
	control := transport.control
	if control == nil { transport.controlMu.Unlock(); return }
	lastStreamID := uint64(0)
	for _, stream := range state.streams {
		streamState := tsNodeQUICStreamFor(stream)
		if streamState != nil && streamState.direction == "bidi" && streamState.id >= 0 && uint64(streamState.id)%4 == 0 && uint64(streamState.id) > lastStreamID {
			lastStreamID = uint64(streamState.id)
		}
	}
	payload := tsNodeQUICAppendVarInt(nil, lastStreamID)
	frame := tsNodeQUICAppendVarInt(nil, 0x7)
	frame = append(frame, tsNodeQUICAppendVarInt(nil, uint64(len(payload)))...)
	frame = append(frame, payload...)
	_, err := control.Write(frame)
	transport.controlMu.Unlock()
	if err != nil && transport.conn != nil { _ = transport.conn.CloseWithError(0, "failed to send HTTP/3 GOAWAY") }
}

func tsNodeQUICSendPriorityUpdate(stream tsValue, urgency uint64, incremental bool) {
	streamState := tsNodeQUICStreamFor(stream)
	if streamState == nil || streamState.id < 0 { return }
	transport := tsNodeQUICSessionTransportState(streamState.session)
	if transport == nil || !transport.http3 { return }
	transport.controlMu.Lock()
	control := transport.control
	if control == nil { transport.controlMu.Unlock(); return }
	priority := "u=" + strconv.FormatUint(urgency, 10)
	if incremental { priority += ", i" }
	payload := tsNodeQUICAppendVarInt(nil, uint64(streamState.id))
	payload = append(payload, priority...)
	frame := tsNodeQUICAppendVarInt(nil, 0xF0700)
	frame = append(frame, tsNodeQUICAppendVarInt(nil, uint64(len(payload)))...)
	frame = append(frame, payload...)
	_, err := control.Write(frame)
	transport.controlMu.Unlock()
	if err != nil && transport.conn != nil { _ = transport.conn.CloseWithError(0, "failed to send HTTP/3 PRIORITY_UPDATE") }
}

func tsNodeQUICDiscardUniStream(stream *quicgo.ReceiveStream) {
	buffer := make([]byte, 4096)
	for {
		if _, err := stream.Read(buffer); err != nil {
			return
		}
	}
}

func tsNodeQUICAppendVarInt(buffer []byte, value uint64) []byte {
	switch {
	case value < 1<<6:
		return append(buffer, byte(value))
	case value < 1<<14:
		return append(buffer, byte(value>>8)|0x40, byte(value))
	case value < 1<<30:
		return append(buffer, byte(value>>24)|0x80, byte(value>>16), byte(value>>8), byte(value))
	default:
		return append(buffer, byte(value>>56)|0xc0, byte(value>>48), byte(value>>40), byte(value>>32), byte(value>>24), byte(value>>16), byte(value>>8), byte(value))
	}
}

func tsNodeQUICReadVarInt(buffer []byte) (uint64, int, bool) {
	if len(buffer) == 0 {
		return 0, 0, false
	}
	length := 1 << (buffer[0] >> 6)
	if len(buffer) < length {
		return 0, 0, false
	}
	value := uint64(buffer[0] & 0x3f)
	for index := 1; index < length; index++ {
		value = (value << 8) | uint64(buffer[index])
	}
	return value, length, true
}

func tsNodeQUICReadVarIntFrom(reader io.Reader) (uint64, []byte, error) {
	first := []byte{0}
	if _, err := io.ReadFull(reader, first); err != nil {
		return 0, nil, err
	}
	length := 1 << (first[0] >> 6)
	data := make([]byte, length)
	data[0] = first[0]
	if length > 1 {
		if _, err := io.ReadFull(reader, data[1:]); err != nil {
			return 0, nil, err
		}
	}
	value, _, _ := tsNodeQUICReadVarInt(data)
	return value, data, nil
}

func tsNodeQUICEncodeHTTP3Headers(loop *tsLoop, headers tsValue) ([]byte, error) {
	var block bytes.Buffer
	encoder := quicqpack.NewEncoder(&block)
	writeField := func(name string, value tsValue) error {
		if value.kind == tsArrayKind {
			array := (*tsArray)(value.ref)
			for index := 0; index < array.length(); index++ {
				if err := encoder.WriteField(quicqpack.HeaderField{Name: strings.ToLower(name), Value: tsStringValueWithLoop(loop, array.at(index)).String()}); err != nil {
					return err
				}
			}
			return nil
		}
		return encoder.WriteField(quicqpack.HeaderField{Name: strings.ToLower(name), Value: tsStringValueWithLoop(loop, value).String()})
	}
	validPseudo := map[string]bool{":status": true, ":method": true, ":authority": true, ":scheme": true, ":path": true, ":protocol": true}
	_ = tsNodeQUICBuildHeaders(loop, headers, tsU, tsU, validPseudo)
	for _, pseudo := range []bool{true, false} {
		if headers.kind == tsArrayKind {
			array := (*tsArray)(headers.ref)
			for index := 0; index+1 < array.length(); index += 2 {
				name := tsText(array.at(index))
				if strings.HasPrefix(name, ":") != pseudo { continue }
				if err := writeField(name, array.at(index+1)); err != nil {
					return nil, err
				}
			}
		} else if tsNodeQUICObject(headers) {
			object := (*tsObject)(headers.ref)
			for _, name := range object.order {
				if strings.HasPrefix(name, ":") != pseudo { continue }
				if err := writeField(name, object.values[name]); err != nil {
					return nil, err
				}
			}
		}
	}
	if err := encoder.Close(); err != nil {
		return nil, err
	}
	frame := tsNodeQUICAppendVarInt(nil, 0x1)
	frame = tsNodeQUICAppendVarInt(frame, uint64(block.Len()))
	return append(frame, block.Bytes()...), nil
}

func tsNodeQUICHeaderSectionSize(loop *tsLoop, headers tsValue) uint64 {
	var total uint64
	add := func(name string, value tsValue) {
		if value.kind == tsArrayKind {
			array := (*tsArray)(value.ref)
			for index := 0; index < array.length(); index++ {
				total += uint64(len(name) + len(tsStringValueWithLoop(loop, array.at(index)).String()) + 32)
			}
			return
		}
		total += uint64(len(name) + len(tsStringValueWithLoop(loop, value).String()) + 32)
	}
	if headers.kind == tsArrayKind {
		array := (*tsArray)(headers.ref)
		for index := 0; index+1 < array.length(); index += 2 {
			add(strings.ToLower(tsText(array.at(index))), array.at(index+1))
		}
	} else if tsNodeQUICObject(headers) {
		object := (*tsObject)(headers.ref)
		for _, name := range object.order {
			add(strings.ToLower(name), object.values[name])
		}
	}
	return total
}

func tsNodeQUICSendHeaders(loop *tsLoop, stream tsValue, headers tsValue, terminal bool, kind string) tsValue {
	state := tsNodeQUICStreamFor(stream)
	wireState := tsNodeQUICWireStateFor(stream)
	if state == nil || wireState == nil || wireState.transport == nil || state.closed {
		return tsBooleanValue(false)
	}
	if !wireState.transport.http3 {
		panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "The negotiated QUIC application protocol does not support headers")})
	}
	if transportSession := tsNodeQUICSessionTransportState(state.session); transportSession != nil && tsNodeQUICHeaderSectionSize(loop, headers) > transportSession.peerMaxFieldSectionSize {
		panic(tsThrown{tsNodeQUICError("ERR_HTTP3_EXCESSIVE_LOAD", "Header section exceeds the peer's configured limit")})
	}
	if kind == "trailing" {
		if headers.kind == tsObjectKind {
			for _, name := range (*tsObject)(headers.ref).order {
				if strings.HasPrefix(name, ":") {
					panic(tsThrown{tsNodeQUICError("ERR_HTTP3_INVALID_HEADER", "Pseudo-headers are not permitted in trailers")})
				}
			}
		}
	}
	frame, err := tsNodeQUICEncodeHTTP3Headers(loop, headers)
	if err != nil {
		panic(tsThrown{tsNodeQUICError("ERR_HTTP3_QPACK_ERROR", err.Error())})
	}
	if kind == "initial" && (tsIsUndefined(state.headers) || tsIsNull(state.headers)) {
		state.headers = headers
	} else if kind == "trailing" {
		state.pendingTrailers = tsU
	}
	if kind == "initial" {
		wireState.headersSent = true
		wireState.transport.headersSent = true
	}
	tsNodeQUICEnqueueTransportBytes(loop, stream, frame, 0, tsU, terminal)
	return tsBooleanValue(true)
}

func tsNodeQUICSetOutbound(loop *tsLoop, stream, value tsValue) tsValue {
	state := tsNodeQUICStreamFor(stream)
	if state == nil || state.closed {
		panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "Stream is destroyed" )})
	}
	if state.outboundSet { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "Stream outbound already configured")}) }
	if !tsIsUndefined(state.writer) { panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "Stream writer already accessed")}) }
	state.outboundSet = true
	if value.kind == tsPromiseKind {
		(*tsPromise)(value.ref).then(func(result tsResult) {
			if result.rejected { tsNodeQUICEndStream(loop, stream, result.value); return }
			state.outboundSet = false
			tsNodeQUICSetOutbound(loop, stream, result.value)
		})
		return stream
	}
	if tsNodeQUICIsFileHandle(value) {
		tsNodeQUICStartFileHandleBody(loop, stream, value)
		return stream
	}
	if value.kind == tsObjectKind && (*tsObject)(value.ref).blob != nil {
		state.blob = (*tsObject)(value.ref).blob
		state.blobOffset = 0
		tsNodeQUICPumpBlobBody(loop, stream)
		return stream
	}
	if tsIsNull(value) || tsIsUndefined(value) {
		tsNodeQUICFinalizeTransportWrite(loop, stream, tsU)
		return stream
	}
	switch value.kind {
	case tsStringKind, tsTypedArrayKind, tsDataViewKind, tsArrayBufferKind:
		data := tsNodeQUICInputBytes(loop, value)
		completion := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			if problem := tsArg(args, 0); !tsNullish(problem) && !tsIsUndefined(problem) {
				tsNodeNetEmitSafe(loop, stream, "error", problem)
				return tsU
			}
			tsNodeQUICFinalizeTransportWrite(loop, stream, tsU)
			return tsU
		}))
		tsNodeQUICWriteTransport(loop, stream, tsNodeBuffer(data), tsU, completion)
		return stream
	default:
		tsNodeQUICStartIterableBody(loop, stream, value)
		return stream
	}
}

func tsNodeQUICStartIterableBody(loop *tsLoop, stream, source tsValue) {
	state := tsNodeQUICStreamFor(stream)
	if state == nil || state.closed {
		return
	}
	normalized := tsNodeIterFromFrom(loop, source)
	state.outboundIterator = tsAsyncIterate(normalized)
	tsNodeQUICPumpIterableBody(loop, stream, state.outboundIterator)
}

func tsNodeQUICPumpIterableBody(loop *tsLoop, stream tsValue, iterator *tsAsyncIterator) {
	state := tsNodeQUICStreamFor(stream)
	if state == nil || state.closed || state.outboundIterator != iterator {
		return
	}
	loop.await(tsPromiseValue(iterator.next(loop)), func(next tsResult) {
		current := tsNodeQUICStreamFor(stream)
		if current == nil || current.closed || current.outboundIterator != iterator {
			return
		}
		if next.rejected {
			current.outboundIterator = nil
			tsNodeQUICEndStream(loop, stream, next.value)
			return
		}
		if tsTruthy(tsNodeQUICGet(loop, next.value, "done")) {
			current.outboundIterator = nil
			tsNodeQUICFinalizeTransportWrite(loop, stream, tsU)
			return
		}
		batch := tsNodeQUICGet(loop, next.value, "value")
		var bytesValue tsValue
		if batch.kind == tsArrayKind {
			bytesValue = tsInternalStreamIterUtilsConcatBytes(batch)
		} else {
			bytesValue = tsInternalStreamIterUtilsToUint8Array(loop, batch)
		}
		data := tsFSBytes(bytesValue)
		if len(data) == 0 {
			tsNodeQUICPumpIterableBody(loop, stream, iterator)
			return
		}
		completion := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			if problem := tsArg(args, 0); !tsIsUndefined(problem) && !tsNullish(problem) {
				tsNodeQUICEndStream(loop, stream, problem)
			} else {
				tsNodeQUICPumpIterableBody(loop, stream, iterator)
			}
			return tsU
		}))
		tsNodeQUICWriteTransport(loop, stream, tsNodeBuffer(data), tsU, completion)
	})
}

func tsNodeQUICPumpBlobBody(loop *tsLoop, stream tsValue) {
	state := tsNodeQUICStreamFor(stream)
	if state == nil || state.closed || state.blob == nil {
		return
	}
	remaining := state.blob.length - state.blobOffset
	if remaining <= 0 {
		state.blob = nil
		tsNodeQUICFinalizeTransportWrite(loop, stream, tsU)
		return
	}
	size := remaining
	if size > 32*1024 {
		size = 32 * 1024
	}
	offset := state.blobOffset
	state.blobOffset += size
	blob := state.blob
	loop.await(tsPromiseValue(loop.start(func() tsResult { return blob.readRange(offset, size) })), func(readResult tsResult) {
		current := tsNodeQUICStreamFor(stream)
		if current == nil || current.closed || current.blob != blob {
			return
		}
		if readResult.rejected {
			current.blob = nil
			tsNodeQUICEndStream(loop, stream, readResult.value)
			return
		}
		data := tsNodeQUICInputBytes(loop, readResult.value)
		completion := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			if problem := tsArg(args, 0); !tsIsUndefined(problem) && !tsNullish(problem) {
				tsNodeQUICEndStream(loop, stream, problem)
			} else {
				tsNodeQUICPumpBlobBody(loop, stream)
			}
			return tsU
		}))
		tsNodeQUICWriteTransport(loop, stream, tsNodeBuffer(data), tsU, completion)
	})
}

var tsNodeQUICFileHandleLocks sync.Map

func tsNodeQUICIsFileHandle(value tsValue) bool {
	return value.kind == tsObjectKind && (*tsObject)(value.ref).nativeClass != nil && (*tsObject)(value.ref).nativeClass.builtin == "FileHandle"
}

func tsNodeQUICStartFileHandleBody(loop *tsLoop, stream, fileHandle tsValue) {
	state := tsNodeQUICStreamFor(stream)
	object := (*tsObject)(fileHandle.ref)
	if state == nil || state.closed || object.nativeClass == nil || object.nativeClass.builtin != "FileHandle" {
		tsFSInvalid("options.body", "a FileHandle")
	}
	if _, loaded := tsNodeQUICFileHandleLocks.LoadOrStore(object, true); loaded {
		panic(tsThrown{tsNodeQUICError("ERR_INVALID_STATE", "FileHandle is locked")})
	}
	state.fileHandle = fileHandle
	tsNodeQUICPumpFileHandle(loop, stream, fileHandle)
}

func tsNodeQUICPumpFileHandle(loop *tsLoop, stream, fileHandle tsValue) {
	state := tsNodeQUICStreamFor(stream)
	if state == nil || state.closed || !tsStrictEqual(state.fileHandle, fileHandle) {
		return
	}
	buffer := tsNodeBuffer(make([]byte, 32*1024))
	read := tsNodeQUICGet(loop, fileHandle, "read")
	if read.kind != tsFunctionKind {
		tsNodeQUICEndStream(loop, stream, tsNodeQUICError("ERR_INVALID_STATE", "FileHandle cannot be read"))
		return
	}
	result := tsCallReceiver(loop, read, fileHandle, buffer, tsNumberValue(0), tsNumberValue(32768), tsNull)
	loop.await(result, func(readResult tsResult) {
		current := tsNodeQUICStreamFor(stream)
		if current == nil || current.closed || !tsStrictEqual(current.fileHandle, fileHandle) {
			return
		}
		if readResult.rejected {
			tsNodeQUICEndStream(loop, stream, readResult.value)
			return
		}
		count := int(tsNumber(tsNodeQUICGet(loop, readResult.value, "bytesRead")))
		if count <= 0 {
			tsNodeQUICFinalizeTransportWrite(loop, stream, tsU)
			return
		}
		data := append([]byte(nil), tsFSBytes(buffer)[:count]...)
		completion := tsFunctionValue(tsFunc(func(loop *tsLoop, args ...tsValue) tsValue {
			if problem := tsArg(args, 0); !tsIsUndefined(problem) && !tsNullish(problem) {
				tsNodeQUICEndStream(loop, stream, problem)
			} else {
				tsNodeQUICPumpFileHandle(loop, stream, fileHandle)
			}
			return tsU
		}))
		tsNodeQUICWriteTransport(loop, stream, tsNodeBuffer(data), tsU, completion)
	})
}

func tsNodeQUICCloseFileHandle(loop *tsLoop, state *tsNodeQUICStreamState) {
	if state == nil {
		return
	}
	if iterator := state.outboundIterator; iterator != nil {
		state.outboundIterator = nil
		loop.await(tsPromiseValue(iterator.close(loop)), func(tsResult) {})
	}
	state.blob = nil
	if tsIsUndefined(state.fileHandle) || tsIsNull(state.fileHandle) { return }
	fileHandle := state.fileHandle
	state.fileHandle = tsU
	if fileHandle.kind == tsObjectKind {
		tsNodeQUICFileHandleLocks.Delete((*tsObject)(fileHandle.ref))
	}
	close := tsNodeQUICGet(loop, fileHandle, "close")
	if close.kind == tsFunctionKind {
		result := tsCallReceiver(loop, close, fileHandle)
		if result.kind == tsPromiseKind {
			loop.await(result, func(tsResult) {})
		}
	}
}

func tsNodeQUICEnqueueTransportBytes(loop *tsLoop, stream tsValue, data []byte, count int, callback tsValue, end bool) {
	state := tsNodeQUICWireStateFor(stream)
	if state == nil || state.transport == nil || !state.transport.canWrite || state.transport.writer == nil {
		if callback.kind == tsFunctionKind {
			tsCall(loop, callback, tsNodeQUICError("ERR_QUIC_STREAM_ABORTED", "The QUIC stream is not writable"))
		}
		return
	}
	state.transport.writeQueue = append(state.transport.writeQueue, tsNodeQUICWriteRequest{data: append([]byte(nil), data...), callback: callback, count: count, end: end})
	state.transport.queuedBytes += len(data)
	tsNodeQUICWriteNext(loop, stream, state)
}

func tsNodeQUICWriteNext(loop *tsLoop, stream tsValue, state *tsNodeQUICStreamWireState) {
	if state == nil || state.transport == nil || state.transport.writeBusy || len(state.transport.writeQueue) == 0 {
		return
	}
	transport := state.transport
	request := transport.writeQueue[0]
	transport.writeQueue = transport.writeQueue[1:]
	transport.writeBusy = true
	var ready <-chan struct{}
	if transport.http3 { if sessionTransport := tsNodeQUICSessionTransportState(stateSession(stream)); sessionTransport != nil { ready = sessionTransport.ready } }
	var writeError error
	loop.submit(func() tsResult {
		if ready != nil { <-ready }
		if len(request.data) != 0 {
			_, writeError = transport.writer.Write(request.data)
		}
		if writeError == nil && request.end && transport.closeWrite != nil {
			writeError = transport.closeWrite()
		}
		return tsResult{tsU, false}
	}, func(_ tsResult) {
		transport.writeBusy = false
		transport.queuedBytes -= len(request.data)
		if transport.queuedBytes < 0 { transport.queuedBytes = 0 }
		if writeError == nil {
			transport.totalBytes += int64(request.count)
			tsNodeQUICRecordSent(stream, uint64(request.count))
		}
		if request.end {
			state.writeEnded = true
		}
		if request.callback.kind == tsFunctionKind {
			if writeError != nil {
				tsCall(loop, request.callback, tsNodeQUICError("ERR_QUIC_STREAM_ABORTED", writeError.Error()))
			} else {
				tsCall(loop, request.callback, tsU)
			}
		}
		if writeError != nil {
			tsNodeNetEmitSafe(loop, stream, "error", tsNodeQUICError("ERR_QUIC_STREAM_ABORTED", writeError.Error()))
			if streamState := tsNodeQUICStreamFor(stream); streamState != nil {
				streamState.writerErrored = true
				streamState.writerError = tsNodeQUICError("ERR_QUIC_STREAM_ABORTED", writeError.Error())
				for _, pending := range streamState.drains { pending.settle(tsResult{streamState.writerError, true}) }
				streamState.drains = nil
			}
		}
		if streamState := tsNodeQUICStreamFor(stream); streamState != nil { tsNodeQUICWakeWriter(loop, streamState) }
		tsNodeQUICMaybeCloseTransportStream(loop, stream)
		tsNodeQUICWriteNext(loop, stream, state)
	})
}

func stateSession(stream tsValue) tsValue {
	if state := tsNodeQUICStreamFor(stream); state != nil {
		return state.session
	}
	return tsU
}

func tsNodeQUICApplySessionCallbacks(loop *tsLoop, state *tsNodeQUICSessionState, options tsValue) {
	if state == nil {
		return
	}
	for _, name := range []string{"onerror", "onstream", "ondatagram", "ondatagramstatus", "onpathvalidation", "onkeylog", "onqlog", "onsessionticket", "onversionnegotiation", "onhandshake", "onnewtoken", "onearlyrejected", "onorigin", "ongoaway", "onheaders", "ontrailers", "oninfo", "onwanttrailers"} {
		callback := tsNodeQUICGet(loop, options, name)
		if callback.kind == tsFunctionKind {
			state.callbacks[name] = callback
		}
	}
}

func tsNodeQUICHandshakeInfo(loop *tsLoop, session tsValue, connection *quicgo.Conn) {
	state := tsNodeQUICSessionFor(session)
	if state == nil || connection == nil {
		return
	}
	tlsState := connection.ConnectionState().TLS
	tsNodeQUICStatsSet(state.stats, "handshakeCompletedAt", time.Now().UnixNano())
	tsNodeQUICStatsSet(state.stats, "handshakeConfirmedAt", time.Now().UnixNano())
	info := tsNewObject()
	info.set("servername", tsNodeQUICKey(tlsState.ServerName))
	info.set("protocol", tsNodeQUICKey(tlsState.NegotiatedProtocol))
	info.set("cipher", tsNodeQUICKey(cryptotls.CipherSuiteName(tlsState.CipherSuite)))
	info.set("cipherVersion", tsNodeQUICKey(cryptotls.VersionName(tlsState.Version)))
	info.set("validationErrorReason", tsU)
	info.set("validationErrorCode", tsU)
	info.set("earlyDataAttempted", tsBooleanValue(false))
	info.set("earlyDataAccepted", tsBooleanValue(false))
	info.set("local", tsNodeQUICAddressValue(connection.LocalAddr()))
	info.set("remote", tsNodeQUICAddressValue(connection.RemoteAddr()))
	state.opened = true
	if state.openedPromise != nil { state.openedPromise.settle(tsResult{tsObjectValue(info), false}) }
	callback := state.callbacks["onhandshake"]
	if callback.kind == tsFunctionKind {
		tsCallReceiver(loop, callback, session, tsObjectValue(info))
	}
	tsNodeNetEmitSafe(loop, session, "handshake", tsObjectValue(info))
	diagnosticInfo := tsNodeQUICCloneProperties(loop, tsObjectValue(info))
	(*tsObject)(diagnosticInfo.ref).set("session", session)
	tsNodeQUICPublish(loop, "onSessionHandshakeChannel", diagnosticInfo)
	path := tsNewObject()
	path.set("local", tsNodeQUICAddressValue(connection.LocalAddr()))
	path.set("remote", tsNodeQUICAddressValue(connection.RemoteAddr()))
	state.path = tsObjectValue(path)
}

func tsNodeQUICSessionCertificate(loop *tsLoop, session tsValue) tsValue {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.closed || state.localCertificate == nil { return tsU }
	return tsNodeTLSCertificate(loop, state.localCertificate)
}

func tsNodeQUICSessionPeerCertificate(loop *tsLoop, session tsValue) tsValue {
	state := tsNodeQUICSessionFor(session)
	if state == nil || state.closed { return tsU }
	transport := tsNodeQUICSessionTransportState(session)
	if transport == nil || transport.conn == nil {
		return tsU
	}
	certificates := transport.conn.ConnectionState().TLS.PeerCertificates
	if len(certificates) == 0 {
		return tsU
	}
	return tsNodeTLSCertificate(loop, certificates[0])
}

func tsNodeQUICEphemeralKeyInfo(loop *tsLoop, session tsValue) tsValue {
	if state := tsNodeQUICSessionFor(session); state == nil || state.closed {
		return tsU
	}
	transport := tsNodeQUICSessionTransportState(session)
	if transport == nil || transport.conn == nil {
		return tsU
	}
	curve := transport.conn.ConnectionState().TLS.CurveID
	if curve == 0 {
		return tsNull
	}
	name := map[cryptotls.CurveID]string{
		cryptotls.X25519: "X25519",
		cryptotls.CurveP256: "prime256v1",
		cryptotls.CurveP384: "secp384r1",
		cryptotls.CurveP521: "secp521r1",
	}[curve]
	if name == "" {
		name = fmt.Sprintf("%d", curve)
	}
	result := tsNewObject()
	result.set("type", tsNodeQUICKey("ECDH"))
	result.set("name", tsNodeQUICKey(name))
	return tsObjectValue(result)
}

func tsNodeQUICSessionMaxDatagramSize(session tsValue) tsValue {
	transport := tsNodeQUICSessionTransportState(session)
	if transport == nil || transport.conn == nil { return tsBigIntUnsigned(0) }
	connectionState := transport.conn.ConnectionState()
	state := connectionState.SupportsDatagrams
	if !connectionState.TLS.HandshakeComplete || !state.Local || !state.Remote { return tsBigIntUnsigned(0) }
	// Bun subtracts the DATAGRAM frame overhead and caps the payload budget
	// at 1150 bytes; quic-go exposes capability but not the peer's exact limit.
	return tsBigIntUnsigned(1150)
}

func tsNodeQUICSendDatagram(loop *tsLoop, session, value, encoding tsValue) tsValue {
	state := tsNodeQUICSessionFor(session)
	transport := tsNodeQUICSessionTransportState(session)
	if state == nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", "QUIC connection failed"), true))
	}
	if state.closed || state.closing { return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_INVALID_STATE", "Session is closed"), true)) }
	if value.kind == tsPromiseKind {
		result := loop.promise()
		(*tsPromise)(value.ref).then(func(resolved tsResult) {
			if resolved.rejected { result.settle(resolved); return }
			if current := tsNodeQUICSessionFor(session); current == nil || current.closed || current.closing { result.settle(tsResult{tsBigIntUnsigned(0), false}); return }
			sent := tsNodeQUICSendDatagram(loop, session, resolved.value, encoding)
			if sent.kind == tsPromiseKind { (*tsPromise)(sent.ref).then(result.settle) } else { result.settle(tsResult{sent, false}) }
		})
		return tsPromiseValue(result)
	}
	if transport == nil || transport.conn == nil {
		return tsPromiseValue(loop.resolved(tsNodeQUICError("ERR_QUIC_CONNECTION_FAILED", "QUIC connection failed"), true))
	}
	if transport.http3 && transport.peerH3Datagrams != nil && !*transport.peerH3Datagrams {
		return tsPromiseValue(loop.resolved(tsBigIntUnsigned(0), false))
	}
	if tsNodeQUICValueCode(loop, tsNodeQUICSessionMaxDatagramSize(session), 0) == 0 {
		return tsPromiseValue(loop.resolved(tsBigIntUnsigned(0), false))
	}
	if value.kind != tsStringKind && value.kind != tsTypedArrayKind && value.kind != tsDataViewKind && value.kind != tsArrayBufferKind {
		tsFSInvalid("datagram", "a string or ArrayBufferView")
	}
	var data []byte
	if value.kind == tsStringKind {
		if tsIsUndefined(encoding) { encoding = tsNodeQUICKey("utf8") }
		data = append([]byte(nil), tsNodeBufferVisibleBytes(tsNodeBufferFrom(loop, []tsValue{value, encoding}))...)
	} else {
		data = tsNodeQUICInputBytes(loop, value)
	}
	maxDatagramSize := tsNodeQUICValueCode(loop, tsNodeQUICSessionMaxDatagramSize(session), 0)
	if len(data) == 0 || uint64(len(data)) > maxDatagramSize {
		return tsPromiseValue(loop.resolved(tsBigIntUnsigned(0), false))
	}
	if err := transport.conn.SendDatagram(data); err != nil {
		return tsPromiseValue(loop.resolved(tsBigIntUnsigned(0), false))
	}
	id := state.nextDatagram
	state.nextDatagram++
	tsNodeQUICStatsAdd(state.stats, "datagramsSent", 1)
	tsNodeQUICStatsAdd(state.stats, "bytesSent", uint64(len(data)))
	if endpointState := tsNodeQUICEndpointFor(state.endpoint); endpointState != nil { tsNodeQUICStatsAdd(endpointState.stats, "bytesSent", uint64(len(data))) }
	tsNodeQUICPublish(loop, "onSessionSendDatagramChannel", tsNodeQUICData(map[string]tsValue{"id":tsBigIntUnsigned(id),"length":tsNumberValue(float64(len(data))),"session":session}))
	return tsPromiseValue(loop.resolved(tsBigIntUnsigned(id), false))
}

func tsNodeQUICStreamStateObject(stream tsValue) tsValue {
	state := tsNodeQUICStreamFor(stream)
	if state == nil {
		tsFSInvalid("stream", "a QuicStream")
	}
	wire := tsNodeQUICWireStateFor(stream)
	readEnded, writeEnded, reset := state.ended, state.ended, state.destroyed
	writeDesired := float64(state.highWaterMark)
	if wire != nil && wire.transport != nil {
		readEnded, writeEnded = wire.readEnded, wire.writeEnded
		reset = state.destroyed
		queued := 0
		for _, request := range wire.transport.writeQueue { queued += len(request.data) }
		if wire.transport.writeBusy { writeDesired = 0 } else if queued >= state.highWaterMark { writeDesired = 0 } else { writeDesired = float64(state.highWaterMark-queued) }
	}
	result := tsNewObject()
	for name, value := range map[string]tsValue{
		"id": tsBigIntUnsigned(uint64(state.id)),
		"pending": tsBooleanValue(!state.opened && !state.closed),
		"finSent": tsBooleanValue(writeEnded),
		"finReceived": tsBooleanValue(readEnded),
		"readEnded": tsBooleanValue(readEnded),
		"writeEnded": tsBooleanValue(writeEnded),
		"reset": tsBooleanValue(reset),
		"early": tsBooleanValue(state.early),
		"resetCode": tsNumberValue(0),
		"highWaterMark": tsNumberValue(float64(state.highWaterMark)),
		"writeDesiredSize": tsNumberValue(writeDesired),
		"hasOutbound": tsBooleanValue(state.outboundSet),
		"hasReader": tsBooleanValue(wire != nil && wire.transport != nil && wire.transport.canRead),
		"wantsBlock": tsBooleanValue(state.callbacks["onblocked"].kind == tsFunctionKind),
		"wantsHeaders": tsBooleanValue(state.callbacks["onheaders"].kind == tsFunctionKind),
		"wantsReset": tsBooleanValue(state.callbacks["onreset"].kind == tsFunctionKind),
		"wantsTrailers": tsBooleanValue(state.callbacks["onwanttrailers"].kind == tsFunctionKind),
	} { result.set(name, value) }
	return tsObjectValue(result)
}

func tsNodeQUICRefreshSessionStats(session tsValue) {
	state := tsNodeQUICSessionFor(session)
	transport := tsNodeQUICSessionTransportState(session)
	if state == nil || transport == nil || transport.conn == nil { return }
	stats := transport.conn.ConnectionStats()
	for name, value := range map[string]uint64{
		"bytesReceived": stats.BytesReceived,
		"bytesSent": stats.BytesSent,
		"bytesRecv": stats.BytesReceived,
		"pktSent": stats.PacketsSent,
		"pktRecv": stats.PacketsReceived,
		"pktLost": stats.PacketsLost,
		"bytesLost": stats.BytesLost,
		"cwnd": stats.CongestionWindow,
		"latestRtt": uint64(max(stats.LatestRTT, 0)),
		"minRtt": uint64(max(stats.MinRTT, 0)),
		"rttVar": uint64(max(stats.MeanDeviation, 0)),
		"smoothedRtt": uint64(max(stats.SmoothedRTT, 0)),
	} { tsNodeQUICStatsSetUnsigned(state.stats, name, value) }
}

func tsNodeQUICSessionStateObject(session tsValue) tsValue {
	state := tsNodeQUICSessionFor(session)
	if state == nil {
		tsFSInvalid("session", "a QuicSession")
	}
	tsNodeQUICRefreshSessionStats(session)
	transport := tsNodeQUICSessionTransportState(session)
	h3, handshake, datagrams := false, false, false
	if transport != nil && transport.conn != nil {
		connectionState := transport.conn.ConnectionState()
		h3 = transport.http3
		handshake = connectionState.TLS.HandshakeComplete
		datagrams = connectionState.SupportsDatagrams.Local && connectionState.SupportsDatagrams.Remote
	}
	result := tsNewObject()
	result.set("isClosing", tsBooleanValue(state.closing))
	result.set("isGracefulClose", tsBooleanValue(state.gracefulClose))
	for _, listener := range []struct{name string; callback string}{
		{"hasPathValidationListener", "onpathvalidation"},
		{"hasDatagramListener", "ondatagram"},
		{"hasDatagramStatusListener", "ondatagramstatus"},
		{"hasSessionTicketListener", "onsessionticket"},
		{"hasNewTokenListener", "onnewtoken"},
		{"hasOriginListener", "onorigin"},
	} {
		result.set(listener.name, tsBooleanValue(state.callbacks[listener.callback].kind == tsFunctionKind))
	}
	result.set("isSilentClose", tsBooleanValue(false))
	result.set("isStatelessReset", tsBooleanValue(false))
	result.set("isHandshakeCompleted", tsBooleanValue(handshake))
	result.set("isHandshakeConfirmed", tsBooleanValue(handshake))
	result.set("isStreamOpenAllowed", tsBooleanValue(!state.closed && !state.closing))
	result.set("isPrioritySupported", tsBooleanValue(h3))
	result.set("headersSupported", tsNumberValue(map[bool]float64{true: 1, false: 2}[h3]))
	result.set("isWrapped", tsBooleanValue(true))
	if h3 { result.set("applicationType", tsNumberValue(1)) } else { result.set("applicationType", tsNumberValue(0)) }
	result.set("noErrorCode", tsBigIntUnsigned(0))
	result.set("internalErrorCode", tsBigIntUnsigned(1))
	result.set("lastDatagramId", tsBigIntUnsigned(state.nextDatagram - 1))
	if datagrams { result.set("maxDatagramSize", tsNumberValue(1200)) } else { result.set("maxDatagramSize", tsNumberValue(0)) }
	result.set("maxPendingDatagrams", tsNumberValue(float64(state.maxPendingDatagrams)))
	return tsObjectValue(result)
}

func tsNodeQUICEndpointStateObject(endpoint tsValue) tsValue {
	state := tsNodeQUICEndpointFor(endpoint)
	if state == nil {
		tsFSInvalid("endpoint", "a QuicEndpoint")
	}
	result := tsNewObject()
	result.set("isBound", tsBooleanValue(state.listening))
	result.set("isReceiving", tsBooleanValue(state.listening))
	result.set("isListening", tsBooleanValue(state.listening))
	result.set("isClosing", tsBooleanValue(state.closing))
	result.set("isBusy", tsBooleanValue(state.busy))
	result.set("maxConnectionsPerHost", tsNumberValue(float64(state.maxPerHost)))
	result.set("maxConnectionsTotal", tsNumberValue(float64(state.maxTotal)))
	result.set("pendingCallbacks", tsNumberValue(0))
	return tsObjectValue(result)
}
`

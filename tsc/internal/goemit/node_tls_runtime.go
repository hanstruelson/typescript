package goemit

// nodeTLSRuntime adapts Node's TLS API to Go's standard crypto/tls stack.
// Certificate parsing, sockets, and listeners stay inside the standard library.
const nodeTLSRuntime = `
type tsNodeTLSSecureState struct { object *tsObject; options tsValue; config *cryptotls.Config; digest [32]byte; shared bool; extraMemory int }
type tsNodeTLSSocketState struct { object *tsObject; conn *cryptotls.Conn; server bool; options tsValue; session []byte; pendingSession []byte; sessionCache *tsNodeTLSClientSessionCache }
type tsNodeTLSServerState struct { object *tsObject; loop *tsLoop; server tsValue; config *cryptotls.Config; options tsValue; contexts map[string]*cryptotls.Config; mu sync.RWMutex }
type tsNodeTLSHandshakeListener struct { net.Listener; config *cryptotls.Config }
func (listener *tsNodeTLSHandshakeListener) Accept() (net.Conn,error) { conn,err:=listener.Listener.Accept();if err!=nil{return nil,err};return cryptotls.Server(conn,listener.config),nil }
var tsNodeTLS struct { sync.Once; module tsValue; socketPrototype,serverPrototype *tsObject; secureConstructor,socketConstructor,serverConstructor tsValue; serverProperties *tsObject; secure map[*tsObject]*tsNodeTLSSecureState; sockets map[*tsObject]*tsNodeTLSSocketState; servers map[*tsObject]*tsNodeTLSServerState; minVersion,maxVersion,defaultCiphers,defaultCurve tsValue; defaultCAs []tsValue; defaultCAsSet bool }
func tsNodeTLSKey(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeTLSGet(loop *tsLoop,value tsValue,name string)tsValue{return tsGet(loop,value,tsNodeTLSKey(name))}
func tsNodeTLSSet(loop *tsLoop,value tsValue,name string,item tsValue){tsSet(loop,value,tsNodeTLSKey(name),item)}
func tsNodeTLSObject(value tsValue)bool{return value.kind==tsObjectKind||value.kind==tsInstanceKind}
func tsNodeTLSObjectCopy(value tsValue)*tsObject{out:=tsNewObject();if tsNodeTLSObject(value){source:=(*tsObject)(value.ref);for _,key:=range source.order{out.set(key,source.values[key])}};return out}
func tsNodeTLSBytes(value tsValue)[]byte{if value.kind==tsStringKind{return []byte(tsText(value))};return append([]byte(nil),tsFSBytes(value)...)}
func tsNodeTLSList(value tsValue)[]tsValue{if value.kind==tsArrayKind{return append([]tsValue(nil),(*tsArray)(value.ref).values...)};if tsNullish(value)||tsIsUndefined(value){return nil};return []tsValue{value}}
func tsNodeTLSVersion(value tsValue, fallback uint16)uint16{switch tsText(value){case "TLSv1","TLSv1.0":return cryptotls.VersionTLS10;case "TLSv1.1":return cryptotls.VersionTLS11;case "TLSv1.2":return cryptotls.VersionTLS12;case "TLSv1.3":return cryptotls.VersionTLS13};if tsIsNumeric(value){version:=int(tsNumber(value));if version>=cryptotls.VersionTLS10&&version<=cryptotls.VersionTLS13{return uint16(version)}};return fallback}
func tsNodeTLSPEMValues(loop *tsLoop,value tsValue)[]byte{if tsNullish(value)||!tsTruthy(value){return nil};var out []byte;for _,item:=range tsNodeTLSList(value){if tsNodeTLSObject(item)&&!tsIsUndefined(tsNodeTLSGet(loop,item,"pem")){item=tsNodeTLSGet(loop,item,"pem")};out=append(out,tsNodeTLSBytes(item)...);if len(out)>0&&out[len(out)-1]!='\n'{out=append(out,'\n')}};return out}
func tsNodeTLSParseCACertificates(loop *tsLoop,certs tsValue)([]tsValue,error){
	if certs.kind!=tsArrayKind{return nil,fmt.Errorf("expected an array of certificates")}
	seen:=map[string]bool{};out:=[]tsValue{}
	for _,input:=range (*tsArray)(certs.ref).values{
		var bytes []byte
		switch input.kind{case tsStringKind:bytes=[]byte(tsText(input));case tsTypedArrayKind,tsDataViewKind:bytes=tsFSBytes(input);default:return nil,fmt.Errorf("expected a string or ArrayBufferView")}
		for len(bytes)>0{
			block,rest:=pem.Decode(bytes);if block==nil{if strings.Contains(string(bytes),"-----BEGIN CERTIFICATE-----"){return nil,fmt.Errorf("failed to parse PEM certificate")};break};bytes=rest
			if block.Type!="CERTIFICATE"{continue};certificate,err:=x509.ParseCertificate(block.Bytes);if err!=nil{return nil,err};normalized:=string(pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE",Bytes:certificate.Raw}));if !seen[normalized]{seen[normalized]=true;out=append(out,tsStringReference(tsStringUTF8(normalized)))}
		}
	}
	return out,nil
}
func tsNodeTLSParseCRLs(loop *tsLoop,value tsValue)([]*x509.RevocationList,error){
	crls:=[]*x509.RevocationList{}
	for _,item:=range tsNodeTLSList(value){if tsNodeTLSObject(item)&&!tsIsUndefined(tsNodeTLSGet(loop,item,"pem")){item=tsNodeTLSGet(loop,item,"pem")};data:=tsNodeTLSBytes(item);parsedPEM:=false
		for len(data)>0{block,rest:=pem.Decode(data);if block==nil{break};data=rest;parsedPEM=true;if block.Type!="X509 CRL"&&block.Type!="CRL"{continue};crl,err:=x509.ParseRevocationList(block.Bytes);if err!=nil{return nil,err};crls=append(crls,crl)}
		if !parsedPEM&&len(data)>0{crl,err:=x509.ParseRevocationList(data);if err!=nil{return nil,err};crls=append(crls,crl)}
	}
	return crls,nil
}
func tsNodeTLSCheckCRLs(state cryptotls.ConnectionState,crls []*x509.RevocationList)error{
	if len(crls)==0{return nil};chain:=state.PeerCertificates;if len(state.VerifiedChains)>0&&len(state.VerifiedChains[0])>0{chain=state.VerifiedChains[0]}
	for index,certificate:=range chain{var issuer *x509.Certificate;if index+1<len(chain){issuer=chain[index+1]}else if bytes.Equal(certificate.RawIssuer,certificate.RawSubject){issuer=certificate};if issuer==nil{continue}
		for _,crl:=range crls{if !bytes.Equal(crl.RawIssuer,issuer.RawSubject)||crl.CheckSignatureFrom(issuer)!=nil{continue};if !crl.NextUpdate.IsZero()&&time.Now().After(crl.NextUpdate){return fmt.Errorf("certificate revocation list has expired")};for _,entry:=range crl.RevokedCertificateEntries{if entry.SerialNumber.Cmp(certificate.SerialNumber)==0{return fmt.Errorf("certificate has been revoked")}}}
	}
	return nil
}
func tsNodeTLSOptionPEM(loop *tsLoop,options tsValue,field,fileField string)([]byte,error){
	data:=tsNodeTLSPEMValues(loop,tsNodeTLSGet(loop,options,field))
	path:=tsNodeTLSGet(loop,options,fileField)
	if !tsIsUndefined(path)&&!tsNullish(path){
		if path.kind!=tsStringKind{return nil,fmt.Errorf("TLSOptions.%s must be a path string",fileField)}
		name:=tsText(path);if strings.IndexByte(name,0)>=0{return nil,fmt.Errorf("TLSOptions.%s must be a path without null bytes",fileField)}
		contents,err:=os.ReadFile(name);if err!=nil{return nil,fmt.Errorf("unable to access %s path: %w",fileField,err)}
		data=append(data,contents...)
	}
	return data,nil
}
func tsNodeTLSDecryptKeyPEM(loop *tsLoop,options tsValue,keyPEM []byte)([]byte,error){
	passphrase:=tsNodeTLSGet(loop,options,"passphrase");var password []byte
	if passphrase.kind==tsStringKind{password=[]byte(tsText(passphrase))}else if passphrase.kind==tsTypedArrayKind||passphrase.kind==tsDataViewKind{password=tsFSBytes(passphrase)}
	var output []byte
	for len(keyPEM)>0{
		block,rest:=pem.Decode(keyPEM);if block==nil{output=append(output,keyPEM...);break};keyPEM=rest
		if x509.IsEncryptedPEMBlock(block){plain,err:=x509.DecryptPEMBlock(block,password);if err!=nil{return nil,fmt.Errorf("failed to decrypt private key: %w",err)};block.Bytes=plain;block.Headers=nil}else if block.Type=="ENCRYPTED PRIVATE KEY"{plain,err:=tsNodeCryptoDecryptPKCS8(block.Bytes,password);if err!=nil{return nil,fmt.Errorf("failed to decrypt private key: %w",err)};block.Type="PRIVATE KEY";block.Bytes=plain;block.Headers=nil}
		output=append(output,pem.EncodeToMemory(block)...)
	}
	return output,nil
}
func tsNodeTLSPKCS12Certificates(loop *tsLoop,options tsValue)([]cryptotls.Certificate,error){
	pfx:=tsNodeTLSGet(loop,options,"pfx");if tsNullish(pfx)||!tsTruthy(pfx){return nil,nil}
	defaultPassphrase:=tsNodeTLSGet(loop,options,"passphrase");certificates:=[]cryptotls.Certificate{}
	for _,entry:=range tsNodeTLSList(pfx){buffer:=entry;passphrase:=defaultPassphrase
		if tsNodeTLSObject(entry){buffer=tsNodeTLSGet(loop,entry,"buf");if tsIsUndefined(buffer){return nil,fmt.Errorf("PKCS#12 entry must include a buf value")};entryPassphrase:=tsNodeTLSGet(loop,entry,"passphrase");if !tsIsUndefined(entryPassphrase){passphrase=entryPassphrase}}
		key,leaf,chain,err:=tsNodeTLSDecodePKCS12(tsNodeTLSBytes(buffer),string(tsNodeTLSBytes(passphrase)));if err!=nil{return nil,fmt.Errorf("failed to parse PKCS#12 data: %w",err)}
		certificate:=cryptotls.Certificate{Certificate:[][]byte{append([]byte(nil),leaf.Raw...)},PrivateKey:key,Leaf:leaf};for _,ca:=range chain{certificate.Certificate=append(certificate.Certificate,append([]byte(nil),ca.Raw...))};certificates=append(certificates,certificate)
	}
	return certificates,nil
}
func tsNodeTLSDecodePKCS12(data []byte,password string)(any,*x509.Certificate,[]*x509.Certificate,error){
	key,leaf,err:=cryptopkcs12.Decode(data,password);if err!=nil{return nil,nil,nil,err};blocks,err:=cryptopkcs12.ToPEM(data,password);if err!=nil{return nil,nil,nil,err};chain:=[]*x509.Certificate{};seen:=map[string]bool{string(leaf.Raw):true}
	for _,block:=range blocks{if block.Type!="CERTIFICATE"{continue};cert,parseErr:=x509.ParseCertificate(block.Bytes);if parseErr!=nil{return nil,nil,nil,parseErr};fingerprint:=string(cert.Raw);if seen[fingerprint]{continue};seen[fingerprint]=true;chain=append(chain,cert)}
	return key,leaf,chain,nil
}
func tsNodeTLSPKCS12Parse(loop *tsLoop,args []tsValue)tsValue{
	input:=tsArg(args,0);if tsNullish(input){panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED","PFX certificate argument is mandatory")})};data:=tsNodeTLSBytes(input);if len(data)==0{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED","PFX certificate argument is mandatory")})}
	passphrase:=tsArg(args,1);if tsNullish(passphrase){passphrase=tsNodeTLSKey("")};key,leaf,chain,err:=tsNodeTLSDecodePKCS12(data,string(tsNodeTLSBytes(passphrase)));if err!=nil{message:="Unable to load PFX certificate";lower:=strings.ToLower(err.Error());if strings.Contains(lower,"mac"){message="PFX MAC verification failed - is the passphrase correct?"}else if strings.Contains(lower,"private key"){message="Unable to load private key from PFX data"}else if strings.Contains(lower,"certificate"){message="Unable to load certificate from PFX data"};panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED",message)})}
	keyDER,err:=x509.MarshalPKCS8PrivateKey(key);if err!=nil{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED",err.Error())})};keyPEM:=pem.EncodeToMemory(&pem.Block{Type:"PRIVATE KEY",Bytes:keyDER});certPEM:=pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE",Bytes:leaf.Raw});out:=tsNewObject();out.set("key",tsStringReference(tsStringUTF8(string(keyPEM))));out.set("cert",tsStringReference(tsStringUTF8(string(certPEM))));if len(chain)>0{var caPEM []byte;for _,ca:=range chain{caPEM=append(caPEM,pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE",Bytes:ca.Raw})...)};out.set("ca",tsStringReference(tsStringUTF8(string(caPEM))))};return tsObjectValue(out)
}
func tsNodeTLSCurvePreferences(value tsValue)([]cryptotls.CurveID,error){
	if value.kind!=tsStringKind||tsText(value)==""||tsText(value)=="auto"{return nil,nil}
	var curves []cryptotls.CurveID
	for _,name:=range strings.Split(tsText(value),":"){
		switch strings.ToUpper(strings.TrimSpace(name)){
		case "X25519":curves=append(curves,cryptotls.X25519)
		case "P-256","PRIME256V1","SECP256R1":curves=append(curves,cryptotls.CurveP256)
		case "P-384","SECP384R1":curves=append(curves,cryptotls.CurveP384)
		case "P-521","SECP521R1":curves=append(curves,cryptotls.CurveP521)
		default:return nil,fmt.Errorf("unknown ecdhCurve %q",name)
		}
	}
	return curves,nil
}
func tsNodeTLSConfig(loop *tsLoop,options tsValue,server bool)(*cryptotls.Config,error){
	if state:=tsNodeTLSSecureFor(options);state!=nil&&state.config!=nil{return state.config.Clone(),nil}
	if tsNullish(options)||tsIsUndefined(options){options=tsObjectValue(tsNewObject())}
	config:=&cryptotls.Config{MinVersion:tsNodeTLSVersion(tsNodeTLSGet(loop,options,"minVersion"),tsNodeTLSVersion(tsNodeTLS.minVersion,cryptotls.VersionTLS12)),MaxVersion:tsNodeTLSVersion(tsNodeTLSGet(loop,options,"maxVersion"),tsNodeTLSVersion(tsNodeTLS.maxVersion,cryptotls.VersionTLS13))}
	pfxCertificates,err:=tsNodeTLSPKCS12Certificates(loop,options);if err!=nil{return nil,err};config.Certificates=append(config.Certificates,pfxCertificates...)
	crls,err:=tsNodeTLSParseCRLs(loop,tsNodeTLSGet(loop,options,"crl"));if err!=nil{return nil,fmt.Errorf("failed to parse CRL: %w",err)}
	if keylog:=tsNodeNetKeyLogWriter(loop);keylog!=nil{config.KeyLogWriter=keylog}
	servername:=tsNodeTLSGet(loop,options,"servername");if tsIsUndefined(servername){servername=tsNodeTLSGet(loop,options,"serverName")};if servername.kind==tsStringKind{config.ServerName=tsText(servername);if strings.IndexByte(config.ServerName,0)>=0{return nil,fmt.Errorf("serverName must not contain null bytes")}}
	secureOptions:=tsNodeTLSGet(loop,options,"secureOptions");if tsIsNumeric(secureOptions){flags:=uint64(tsNumber(secureOptions));if flags&0x20000000!=0&&config.MaxVersion>cryptotls.VersionTLS12{config.MaxVersion=cryptotls.VersionTLS12};if flags&0x04000000!=0&&config.MinVersion<=cryptotls.VersionTLS10{config.MinVersion=cryptotls.VersionTLS11};if flags&0x10000000!=0&&config.MinVersion<=cryptotls.VersionTLS11{config.MinVersion=cryptotls.VersionTLS12};if flags&0x08000000!=0&&config.MinVersion<=cryptotls.VersionTLS12{config.MinVersion=cryptotls.VersionTLS13};if flags&0x00004000!=0{config.SessionTicketsDisabled=true}}
	curves,curveErr:=tsNodeTLSCurvePreferences(tsNodeTLSGet(loop,options,"ecdhCurve"));if curveErr!=nil{return nil,curveErr};config.CurvePreferences=curves
	ca,err:=tsNodeTLSOptionPEM(loop,options,"ca","caFile");if err!=nil{return nil,err}
	caOption:=tsNodeTLSGet(loop,options,"ca")
	if tsNodeTLS.defaultCAsSet&&(tsNullish(caOption)||!tsTruthy(caOption)){fileCA,fileErr:=tsNodeTLSOptionPEM(loop,options,"","caFile");if fileErr!=nil{return nil,fileErr};ca=append(tsNodeTLSPEMValues(loop,tsArrayValue(&tsArray{values:tsNodeTLS.defaultCAs})),fileCA...)}
	if !server{
		reject:=tsNodeTLSGet(loop,options,"rejectUnauthorized");rejectUnauthorized:=!tsStrictEqual(reject,tsBooleanValue(false));standalone:=tsTruthy(tsNodeTLSGet(loop,options,"_standaloneWrap"));config.InsecureSkipVerify=!rejectUnauthorized||standalone
		if tsIsUndefined(reject){if value,ok:=os.LookupEnv("NODE_TLS_REJECT_UNAUTHORIZED");ok&&value=="0"{config.InsecureSkipVerify=true;rejectUnauthorized=false}}
		config.RootCAs,_=cryptotls.SystemCertPool();if config.RootCAs==nil{config.RootCAs=x509.NewCertPool()}
		if tsNodeTLS.defaultCAsSet&&(tsNullish(caOption)||!tsTruthy(caOption)){config.RootCAs=x509.NewCertPool()}
		if len(ca)>0&&!config.RootCAs.AppendCertsFromPEM(ca){return nil,fmt.Errorf("failed to parse CA certificate")}
		if (tsNodeTLSGet(loop,options,"checkServerIdentity").kind==tsFunctionKind||standalone)&&rejectUnauthorized{config.InsecureSkipVerify=true;roots:=config.RootCAs;config.VerifyConnection=func(state cryptotls.ConnectionState)error{return tsNodeTLSVerifyChain(state,roots)}}
	}else{
		config.ClientAuth=cryptotls.NoClientCert;if tsTruthy(tsNodeTLSGet(loop,options,"requestCert")){if tsStrictEqual(tsNodeTLSGet(loop,options,"rejectUnauthorized"),tsBooleanValue(false)){config.ClientAuth=cryptotls.RequestClientCert}else{config.ClientAuth=cryptotls.RequireAndVerifyClientCert}}
		config.ClientCAs=x509.NewCertPool();if len(ca)>0&&!config.ClientCAs.AppendCertsFromPEM(ca){return nil,fmt.Errorf("failed to parse CA certificate")}
	}
	protos:=tsNodeTLSALPN(loop,tsNodeTLSGet(loop,options,"ALPNProtocols"));if len(protos)>0{config.NextProtos=protos}
	certPEM,err:=tsNodeTLSOptionPEM(loop,options,"cert","certFile");if err!=nil{return nil,err};keyPEM,err:=tsNodeTLSOptionPEM(loop,options,"key","keyFile");if err!=nil{return nil,err}
	if len(certPEM)>0&&len(keyPEM)>0{keyPEM,err=tsNodeTLSDecryptKeyPEM(loop,options,keyPEM);if err!=nil{return nil,err};certificate,err:=cryptotls.X509KeyPair(certPEM,keyPEM);if err!=nil{return nil,err};config.Certificates=append(config.Certificates,certificate)}
	if len(crls)>0&&(!config.InsecureSkipVerify||config.VerifyConnection!=nil){previousVerify:=config.VerifyConnection;config.VerifyConnection=func(state cryptotls.ConnectionState)error{if previousVerify!=nil{if verifyErr:=previousVerify(state);verifyErr!=nil{return verifyErr}};return tsNodeTLSCheckCRLs(state,crls)}}
	suites:=tsNodeTLSGet(loop,options,"ciphers");if tsIsUndefined(suites)&&tsNodeTLS.defaultCiphers.kind==tsStringKind{suites=tsNodeTLS.defaultCiphers};if suites.kind==tsStringKind{if tsText(suites)==""{config.MinVersion=max(config.MinVersion,cryptotls.VersionTLS13)}else{allowed:=map[string]uint16{};for _,suite:=range append(cryptotls.CipherSuites(),cryptotls.InsecureCipherSuites()...){allowed[suite.Name]=suite.ID};for _,name:=range strings.Split(tsText(suites),":"){if id,ok:=allowed[name];ok{config.CipherSuites=append(config.CipherSuites,id)}}}}
	return config,nil
}
func tsNodeTLSSecureFor(value tsValue)*tsNodeTLSSecureState{if value.kind!=tsObjectKind{return nil};return tsNodeTLS.secure[(*tsObject)(value.ref)]}
func tsNodeTLSSocketFor(value tsValue)*tsNodeTLSSocketState{if value.kind!=tsObjectKind{return nil};return tsNodeTLS.sockets[(*tsObject)(value.ref)]}
func tsNodeTLSServerFor(value tsValue)*tsNodeTLSServerState{if value.kind!=tsObjectKind{return nil};return tsNodeTLS.servers[(*tsObject)(value.ref)]}
func tsNodeTLSServerNameKey(name string)string{return strings.TrimSuffix(strings.ToLower(name),".")}
func tsNodeTLSServerContextForName(contexts map[string]*cryptotls.Config,servername string)*cryptotls.Config{
	host:=tsNodeTLSServerNameKey(servername);if selected:=contexts[host];selected!=nil{return selected}
	for pattern,selected:=range contexts{if !strings.HasPrefix(pattern,"*."){continue};suffix:=pattern[1:];if strings.HasSuffix(host,suffix){prefix:=strings.TrimSuffix(host,suffix);if prefix!=""&&!strings.Contains(prefix,"."){return selected}}}
	return nil
}
type tsNodeTLSSNIError struct { value tsValue; message string }
func (failure *tsNodeTLSSNIError) Error()string{return failure.message}
type tsNodeTLSSNIResult struct { called bool; config *cryptotls.Config; err error }
func tsNodeTLSErrorValue(value tsValue)bool{if value.kind==tsErrorKind{return true};if value.kind==tsObjectKind{return (*tsObject)(value.ref).nativeError};if value.kind==tsInstanceKind{return (*tsProperties)(value.ref).nativeError};return false}
func tsNodeTLSSNIErrorValue(loop *tsLoop,value tsValue)tsValue{if tsNodeTLSErrorValue(value){return value};err:=tsFSException("Error","ERR_TLS_HANDSHAKE_FAILURE","SNI callback error");tsNodeTLSSet(loop,err,"reason",value);return err}
func tsNodeTLSSNIContext(loop *tsLoop,value tsValue)(*cryptotls.Config,error){
	if tsNullish(value)||tsIsUndefined(value){return nil,nil}
	if inner:=tsNodeTLSGet(loop,value,"context");!tsNullish(inner)&&!tsIsUndefined(inner){value=inner}
	if secure:=tsNodeTLSSecureFor(value);secure!=nil&&secure.config!=nil{return secure.config.Clone(),nil}
	return nil,&tsNodeTLSSNIError{message:"Invalid SNI context",value:tsFSException("Error","ERR_TLS_HANDSHAKE_FAILURE","Invalid SNI context")}
}
func tsNodeTLSSNICallback(state *tsNodeTLSServerState,servername string,handshakeTimeout time.Duration)(*cryptotls.Config,bool,error){
	if servername==""{return nil,false,nil}
	result:=make(chan tsNodeTLSSNIResult,1)
	state.loop.post(func(){
		callback:=tsNodeTLSGet(state.loop,state.server,"_SNICallback");if callback.kind!=tsFunctionKind{callback=tsNodeTLSGet(state.loop,state.options,"SNICallback")};if callback.kind!=tsFunctionKind{result<-tsNodeTLSSNIResult{called:false};return}
		settled:=false
		complete:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{
			if settled{return tsU};settled=true;errorValue,contextValue:=tsArg(args,0),tsArg(args,1)
			if tsTruthy(errorValue){normalized:=tsNodeTLSSNIErrorValue(loop,errorValue);message:=tsText(tsNodeTLSGet(loop,normalized,"message"));result<-tsNodeTLSSNIResult{called:true,err:&tsNodeTLSSNIError{value:normalized,message:message}};return tsU}
			context,err:=tsNodeTLSSNIContext(loop,contextValue);result<-tsNodeTLSSNIResult{called:true,config:context,err:err};return tsU
		}))
		defer func(){if thrown:=recover();thrown!=nil&&!settled{settled=true;value:=tsUnwrap(thrown);normalized:=tsNodeTLSSNIErrorValue(state.loop,value);message:=tsText(tsNodeTLSGet(state.loop,normalized,"message"));result<-tsNodeTLSSNIResult{called:true,err:&tsNodeTLSSNIError{value:normalized,message:message}}}}()
		tsCallReceiver(state.loop,callback,state.server,tsNodeTLSKey(servername),complete)
	})
	var timeout *time.Timer;var timeoutChannel <-chan time.Time;if handshakeTimeout>0{timeout=time.NewTimer(handshakeTimeout);timeoutChannel=timeout.C;defer timeout.Stop()}
	select{case resolved:=<-result:return resolved.config,resolved.called,resolved.err
	case <-timeoutChannel:return nil,true,&tsNodeTLSSNIError{value:tsFSException("Error","ERR_TLS_HANDSHAKE_TIMEOUT","TLS handshake timeout"),message:"TLS handshake timeout"}
	case <-state.loop.stop:return nil,true,errors.New("TLS handshake stopped")}
}
type tsNodeTLSALPNResult struct { called bool; protocol string; err error }
func tsNodeTLSServerALPN(state *tsNodeTLSServerState,servername string,protocols []string,handshakeTimeout time.Duration)(string,bool,error){
	result:=make(chan tsNodeTLSALPNResult,1)
	state.loop.post(func(){
		callback:=tsNodeTLSGet(state.loop,state.server,"_ALPNCallback");if callback.kind!=tsFunctionKind{callback=tsNodeTLSGet(state.loop,state.options,"ALPNCallback")};if callback.kind!=tsFunctionKind{result<-tsNodeTLSALPNResult{called:false};return}
		protocolValues:=&tsArray{};for _,protocol:=range protocols{protocolValues.values=append(protocolValues.values,tsNodeTLSKey(protocol))}
		servernameValue:=tsValue(tsBooleanValue(false));if servername!=""{servernameValue=tsNodeTLSKey(servername)}
		info:=tsNewObject();info.set("servername",servernameValue);info.set("protocols",tsArrayValue(protocolValues))
		defer func(){if thrown:=recover();thrown!=nil{value:=tsUnwrap(thrown);result<-tsNodeTLSALPNResult{called:true,err:&tsNodeTLSSNIError{value:value,message:tsText(value)}}}}()
		selected:=tsCallReceiver(state.loop,callback,state.server,tsObjectValue(info))
		if tsIsUndefined(selected){result<-tsNodeTLSALPNResult{called:true,err:&tsNodeTLSSNIError{value:tsFSException("Error","EPROTO","no application protocol"),message:"no application protocol"}};return}
		for _,protocol:=range protocols{if tsStrictEqual(selected,tsNodeTLSKey(protocol)){result<-tsNodeTLSALPNResult{called:true,protocol:protocol};return}}
		joined:=strings.Join(protocols,", ");message:=fmt.Sprintf("ALPN callback returned a value (%s) that did not match any of the client's offered protocols (%s)",tsText(selected),joined)
		value:=tsFSException("Error","ERR_TLS_ALPN_CALLBACK_INVALID_RESULT",message);result<-tsNodeTLSALPNResult{called:true,err:&tsNodeTLSSNIError{value:value,message:message}}
	})
	var timeout *time.Timer;var timeoutChannel <-chan time.Time;if handshakeTimeout>0{timeout=time.NewTimer(handshakeTimeout);timeoutChannel=timeout.C;defer timeout.Stop()}
	select{case resolved:=<-result:return resolved.protocol,resolved.called,resolved.err
	case <-timeoutChannel:return "",true,&tsNodeTLSSNIError{value:tsFSException("Error","ERR_TLS_HANDSHAKE_TIMEOUT","TLS handshake timeout"),message:"TLS handshake timeout"}
	case <-state.loop.stop:return "",true,errors.New("TLS handshake stopped")}
}
func tsNodeTLSServerConfig(state *tsNodeTLSServerState,config *cryptotls.Config)*cryptotls.Config{handshakeTimeout:=tsNodeTLSHandshakeTimeout(state.loop,state.options);out:=config.Clone();out.GetConfigForClient=func(hello *cryptotls.ClientHelloInfo)(*cryptotls.Config,error){
		var base *cryptotls.Config;var selected *cryptotls.Config
		state.mu.RLock();selected=tsNodeTLSServerContextForName(state.contexts,hello.ServerName);base=state.config;state.mu.RUnlock()
		context,called,err:=tsNodeTLSSNICallback(state,hello.ServerName,handshakeTimeout);if err!=nil{return nil,err};if called{if context!=nil{return context,nil};selected=nil}
		var fallback *cryptotls.Config;if selected!=nil{fallback=selected.Clone()}else if base!=nil{fallback=base.Clone();fallback.GetConfigForClient=nil}else{return nil,nil}
		protocol,alpnCalled,err:=tsNodeTLSServerALPN(state,hello.ServerName,hello.SupportedProtos,handshakeTimeout);if err!=nil{return nil,err};if alpnCalled{fallback=fallback.Clone();fallback.NextProtos=[]string{protocol}}
		return fallback,nil
	};return out}
func tsNodeTLSFailure(err error)tsValue{if err==nil{return tsNull};message:=strings.ToLower(err.Error());code:="ERR_TLS_HANDSHAKE_FAILURE";if strings.Contains(message,"pkcs#12"){code="ERR_CRYPTO_OPERATION_FAILED"}else if strings.Contains(message,"certificate")||strings.Contains(message,"x509"){code="ERR_TLS_CERT_ALTNAME_INVALID"};return tsFSException("Error",code,err.Error())}
func tsNodeTLSHandshakeFailure(loop *tsLoop,err error)tsValue{
	if err==nil{return tsNull};var sniFailure *tsNodeTLSSNIError;if errors.As(err,&sniFailure)&&!tsIsUndefined(sniFailure.value){return sniFailure.value};code:="EPROTO";lower:=strings.ToLower(err.Error())
	if errors.Is(err,io.EOF)||strings.Contains(lower,"connection reset"){code="ECONNRESET"}else if strings.Contains(lower,"not valid for")||strings.Contains(lower,"hostname"){code="ERR_TLS_CERT_ALTNAME_INVALID"}else if strings.Contains(lower,"expired"){code="CERT_HAS_EXPIRED"}else if strings.Contains(lower,"self-signed"){code="SELF_SIGNED_CERT_IN_CHAIN"}else if strings.Contains(lower,"certificate")||strings.Contains(lower,"x509"){code="UNABLE_TO_VERIFY_LEAF_SIGNATURE"}
	verify:=tsNewObject();verify.set("code",tsNodeTLSKey(code));verify.set("message",tsNodeTLSKey(err.Error()));return tsInternalTLSHandshakeError(loop,tsObjectValue(verify))
}
type tsNodeTLSKeyLogSocketWriter struct { loop *tsLoop; socket tsValue; base io.Writer }
func (writer tsNodeTLSKeyLogSocketWriter) Write(line []byte)(int,error){if writer.base!=nil{if _,err:=writer.base.Write(line);err!=nil{return 0,err}};copyOfLine:=append([]byte(nil),line...);writer.loop.post(func(){tsNodeNetEmitSafe(writer.loop,writer.socket,"keylog",tsNodeBuffer(copyOfLine))});return len(line),nil}
type tsNodeTLSClientSessionCache struct { mu sync.Mutex; entries map[string]*cryptotls.ClientSessionState; loop *tsLoop; socket tsValue }
func (cache *tsNodeTLSClientSessionCache) Get(key string)(*cryptotls.ClientSessionState,bool){cache.mu.Lock();defer cache.mu.Unlock();session,ok:=cache.entries[key];return session,ok}
func (cache *tsNodeTLSClientSessionCache) Put(key string,session *cryptotls.ClientSessionState){if session==nil{return};cache.mu.Lock();if cache.entries==nil{cache.entries=map[string]*cryptotls.ClientSessionState{}};cache.entries[key]=session;cache.mu.Unlock();ticket,_,err:=session.ResumptionState();if err!=nil||len(ticket)==0{return};copyOfTicket:=append([]byte(nil),ticket...);cache.loop.post(func(){state:=tsNodeTLSSocketFor(cache.socket);if state==nil{return};state.session=copyOfTicket;if tsTruthy(tsNodeTLSGet(cache.loop,cache.socket,"_secureEstablished")){tsNodeNetEmitSafe(cache.loop,cache.socket,"session",tsNodeBuffer(copyOfTicket))}else{state.pendingSession=copyOfTicket}})}
func tsNodeTLSFlushPendingSession(loop *tsLoop,socket tsValue){state:=tsNodeTLSSocketFor(socket);if state==nil||len(state.pendingSession)==0{return};session:=append([]byte(nil),state.pendingSession...);state.pendingSession=nil;tsNodeNetEmitSafe(loop,socket,"session",tsNodeBuffer(session))}
func tsNodeTLSOnConnectEnd(loop *tsLoop,socket,options tsValue){
	if tsTruthy(tsNodeTLSGet(loop,socket,"_hadError"))||!tsTruthy(tsNodeTLSGet(loop,socket,"secureConnecting")){return}
	tsNodeTLSSet(loop,socket,"_hadError",tsBooleanValue(true));err:=tsFSException("Error","ECONNRESET","Client network socket disconnected before secure TLS connection was established")
	for _,name:=range []string{"path","host","port","localAddress"}{tsNodeTLSSet(loop,err,name,tsNodeTLSGet(loop,options,name))}
	destroy:=tsNodeTLSGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,err)}
}
func tsNodeTLSVerifyFailure(err error)tsValue{
	if err==nil{return tsNull};lower:=strings.ToLower(err.Error());code:="UNABLE_TO_VERIFY_LEAF_SIGNATURE";if strings.Contains(lower,"not valid for")||strings.Contains(lower,"hostname"){code="ERR_TLS_CERT_ALTNAME_INVALID"}else if strings.Contains(lower,"expired"){code="CERT_HAS_EXPIRED"}else if strings.Contains(lower,"self-signed"){code="SELF_SIGNED_CERT_IN_CHAIN"};return tsFSException("Error",code,err.Error())
}
func tsNodeTLSEstablishedError(loop *tsLoop,socket tsValue,err error,operation string)tsValue{
	if err==nil{return tsNull};if tsNodeTLSSocketFor(socket)==nil||!tsTruthy(tsNodeTLSGet(loop,socket,"_secureEstablished"))||errors.Is(err,io.EOF){return tsNodeNetError(err,operation)}
	var op *net.OpError;var errno syscall.Errno;if errors.As(err,&op)||errors.As(err,&errno){return tsNodeNetError(err,operation)}
	verify:=tsNewObject();verify.set("code",tsNodeTLSKey("EPROTO"));verify.set("message",tsNodeTLSKey(err.Error()));return tsInternalTLSHandshakeError(loop,tsObjectValue(verify))
}
func tsNodeTLSVerifyChain(state cryptotls.ConnectionState,roots *x509.CertPool)error{
	if len(state.PeerCertificates)==0{return errors.New("peer did not return a certificate")}
	intermediates:=x509.NewCertPool();for _,certificate:=range state.PeerCertificates[1:]{intermediates.AddCert(certificate)}
	_,err:=state.PeerCertificates[0].Verify(x509.VerifyOptions{Roots:roots,Intermediates:intermediates,KeyUsages:[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}});return err
}
func tsNodeTLSVerifyClientChain(state cryptotls.ConnectionState,roots *x509.CertPool)error{
	if len(state.PeerCertificates)==0{return errors.New("peer did not return a certificate")}
	intermediates:=x509.NewCertPool();for _,certificate:=range state.PeerCertificates[1:]{intermediates.AddCert(certificate)}
	_,err:=state.PeerCertificates[0].Verify(x509.VerifyOptions{Roots:roots,Intermediates:intermediates,KeyUsages:[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}});return err
}
func tsNodeTLSIdentityError(loop *tsLoop,socket tsValue,conn *cryptotls.Conn,options tsValue,config *cryptotls.Config,server,standalone bool)tsValue{
	if server||conn==nil{return tsNull};state:=conn.ConnectionState();if len(state.PeerCertificates)==0{return tsNull}
	if len(state.VerifiedChains)==0{if err:=tsNodeTLSVerifyChain(state,config.RootCAs);err!=nil{return tsNodeTLSVerifyFailure(err)}}
	if standalone{return tsNull}
	hostValue:=tsNodeTLSGet(loop,options,"servername");if tsNullish(hostValue){hostValue=tsNodeTLSGet(loop,options,"host")};if tsNullish(hostValue){if netState:=tsNodeNetSocketFor(socket);netState!=nil&&netState.remoteAddress!=""{hostValue=tsNodeTLSKey(netState.remoteAddress)}else{hostValue=tsNodeTLSKey("localhost")}}
	certificate:=tsNodeTLSCertificate(loop,state.PeerCertificates[0]);if check:=tsNodeTLSGet(loop,options,"checkServerIdentity");check.kind==tsFunctionKind{var result tsValue;func(){defer func(){if thrown:=recover();thrown!=nil{result=tsUnwrap(thrown)}}();result=tsCallReceiver(loop,check,options,hostValue,certificate)}();if tsTruthy(result){return result};return tsNull}
	return tsNodeTLSCheckServerIdentity(loop,hostValue,certificate)
}
func tsNodeTLSHandshakeComplete(loop *tsLoop,socket tsValue,conn *cryptotls.Conn,options tsValue,server bool,config *cryptotls.Config)bool{
	standalone:=tsTruthy(tsNodeTLSGet(loop,socket,"_standaloneWrap"));verifyError:=tsNodeTLSIdentityError(loop,socket,conn,options,config,server,standalone)
	tsNodeTLSAttachSocket(loop,socket,conn,server,options);rejectUnauthorized:=tsNodeTLSGet(loop,options,"rejectUnauthorized");reject:=tsStrictEqual(rejectUnauthorized,tsBooleanValue(true))||!server&&tsIsUndefined(rejectUnauthorized)
	if !server&&tsIsUndefined(rejectUnauthorized){if value,ok:=os.LookupEnv("NODE_TLS_REJECT_UNAUTHORIZED");ok&&value=="0"{reject=false}}
	if server{tsNodeNetEmit(loop,socket,"secure",socket);return true}
	if standalone{if tsTruthy(verifyError)&&reject{destroy:=tsNodeTLSGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,verifyError)};return false};tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeTLSFlushPendingSession(loop,socket);tsNodeNetEmit(loop,socket,"secure",socket);return true}
	if tsTruthy(verifyError){tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));code:=tsNodeTLSGet(loop,verifyError,"code");if tsNullish(code){code=tsNodeTLSGet(loop,verifyError,"message")};tsNodeTLSSet(loop,socket,"authorizationError",code);if reject{destroy:=tsNodeTLSGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,verifyError)};tsNodeNetEmit(loop,socket,"secure",socket);return false}}else{tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorizationError",tsNull)}
	tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeNetEmit(loop,socket,"secureConnect",verifyError);tsNodeTLSFlushPendingSession(loop,socket);tsNodeNetEmit(loop,socket,"secure",socket);return true
}
func tsNodeTLSSecureDigest(loop *tsLoop,options tsValue,config *cryptotls.Config)[32]byte{
	var value strings.Builder
	fmt.Fprintf(&value,"%d|%d|%t|%d|%s|",config.MinVersion,config.MaxVersion,config.InsecureSkipVerify,config.ClientAuth,config.ServerName)
	for _,protocol:=range config.NextProtos{fmt.Fprintf(&value,"p:%d:%s|",len(protocol),protocol)}
	for _,suite:=range config.CipherSuites{fmt.Fprintf(&value,"c:%d|",suite)}
	for _,certificate:=range config.Certificates{for _,der:=range certificate.Certificate{fmt.Fprintf(&value,"x:%d:",len(der));value.Write(der)}}
	for _,pool:=range []*x509.CertPool{config.RootCAs,config.ClientCAs}{if pool!=nil{for _,subject:=range pool.Subjects(){fmt.Fprintf(&value,"a:%d:",len(subject));value.Write(subject)}}}
	for _,name:=range []string{"key","cert","ca","crl"}{bytes:=tsNodeTLSPEMValues(loop,tsNodeTLSGet(loop,options,name));fmt.Fprintf(&value,"%s:%d:",name,len(bytes));value.Write(bytes)}
	for _,name:=range []string{"minVersion","maxVersion","secureOptions","sessionTimeout","allowPartialTrustChain","sigalgs","ecdhCurve","ciphers","passphrase"}{item:=tsNodeTLSGet(loop,options,name);fmt.Fprintf(&value,"%s:%d:",name,item.kind);if item.kind==tsStringKind||tsIsNumeric(item)||item.kind==tsBooleanKind{value.WriteString(tsText(item))}else if tsIsNull(item)||tsIsUndefined(item){value.WriteString("nullish")};value.WriteByte('|')}
	return sha256.Sum256([]byte(value.String()))
}
func tsNodeTLSSecureAddCACert(loop *tsLoop,state *tsNodeTLSSecureState,certificate tsValue)tsValue{
	if state==nil{return tsU}
	if state.shared{panic(tsThrown{tsFSException("Error","ERR_TLS_INVALID_CONTEXT","cannot mutate a shared SecureContext; use tls.createSecureContext()")})}
	if certificate.kind!=tsStringKind&&certificate.kind!=tsTypedArrayKind&&certificate.kind!=tsDataViewKind{tsFSInvalid("cert","string or an instance of ArrayBufferView")}
	bytes:=tsNodeTLSBytes(certificate)
	if len(bytes)==0{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","addCACert requires a certificate")})}
	parsed:=0
	remaining:=bytes
	for len(remaining)>0{block,rest:=pem.Decode(remaining);if block==nil{break};remaining=rest;if block.Type!="CERTIFICATE"{continue};cert,err:=x509.ParseCertificate(block.Bytes);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};if state.config.RootCAs==nil{state.config.RootCAs=x509.NewCertPool()};if state.config.ClientCAs==nil{state.config.ClientCAs=x509.NewCertPool()};state.config.RootCAs.AddCert(cert);state.config.ClientCAs.AddCert(cert);parsed++}
	if parsed==0{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED","Invalid CA certificate")})}
	folded:=make([]byte,0,len(state.digest)+len(bytes));folded=append(folded,state.digest[:]...);folded=append(folded,bytes...);state.digest=sha256.Sum256(folded);state.extraMemory+=len(bytes);return tsU
}
func tsNodeTLSSecureCreate(loop *tsLoop,options tsValue)tsValue{
	config,err:=tsNodeTLSConfig(loop,options,true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})}
	object:=tsNewObject();value:=tsObjectValue(object);state:=&tsNodeTLSSecureState{object:object,options:options,config:config,digest:tsNodeTLSSecureDigest(loop,options,config),extraMemory:50*1024}
	for _,name:=range []string{"key","cert","ca","crl"}{state.extraMemory+=len(tsNodeTLSPEMValues(loop,tsNodeTLSGet(loop,options,name)))}
	for _,certificate:=range config.Certificates{for _,der:=range certificate.Certificate{state.extraMemory+=len(der)}}
	tsNodeTLS.secure[object]=state;object.set("context",value);object.set("options",options);object.set("default",value)
	object.set("addCACert",tsNamedNativeMethod("addCACert",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSecureAddCACert(loop,state,tsArg(args,0))}))
	if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}}
	object.descriptors["_external"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get _external",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU})}
	return value
}
func tsNodeTLSALPN(loop *tsLoop,value tsValue)[]string{if value.kind==tsArrayKind{out:=[]string{};for _,item:=range (*tsArray)(value.ref).values{out=append(out,tsText(item))};return out};if value.kind==tsTypedArrayKind||value.kind==tsDataViewKind||value.kind==tsArrayBufferKind{bytes:=tsFSBytes(value);out:=[]string{};for i:=0;i<len(bytes);{size:=int(bytes[i]);i++;if size==0||i+size>len(bytes){break};out=append(out,string(bytes[i:i+size]));i+=size};return out};if value.kind==tsStringKind{return []string{tsText(value)}};return nil}
func tsNodeTLSALPNValue(loop *tsLoop,protocols []string)tsValue{out:=&tsArray{};for _,protocol:=range protocols{out.values=append(out.values,tsNodeTLSKey(protocol))};return tsArrayValue(out)}
func tsNodeTLSAttachSocket(loop *tsLoop,socket tsValue,conn *cryptotls.Conn,server bool,options tsValue)tsValue{if socket.kind!=tsObjectKind{return socket};object:=(*tsObject)(socket.ref);previous:=tsNodeTLSSocketFor(socket);state:=&tsNodeTLSSocketState{object:object,conn:conn,server:server,options:options};if previous!=nil{state.session=previous.session;state.pendingSession=previous.pendingSession;state.sessionCache=previous.sessionCache};tsNodeTLS.sockets[object]=state;tsNodeTLSSet(loop,socket,"encrypted",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorizationError",tsNull);tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_securePending",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_secureEstablished",tsBooleanValue(true));if conn!=nil{connection:=conn.ConnectionState();if connection.ServerName!=""{tsNodeTLSSet(loop,socket,"servername",tsNodeTLSKey(connection.ServerName))};if connection.NegotiatedProtocol!=""{tsNodeTLSSet(loop,socket,"alpnProtocol",tsNodeTLSKey(connection.NegotiatedProtocol))}else{tsNodeTLSSet(loop,socket,"alpnProtocol",tsNull)}};return socket}
// Bun's net.ts links a second TLS Socket to the raw Socket and then adopts a
// native TLS handle pair. Go's crypto/tls wraps the existing net.Conn directly;
// the read pump is briefly interrupted so bytes it already consumed can be
// replayed into the TLS handshake before the same stream resumes with plaintext.
func tsNodeTLSUpgradeSocket(loop *tsLoop,socket,options tsValue,server bool){
	raw:=tsNodeNetSocketFor(socket);if raw==nil{return}
	tsNodeTLSSet(loop,socket,"encrypted",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_securePending",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_secureEstablished",tsBooleanValue(false))
	config,err:=tsNodeTLSConfig(loop,options,server);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})}
	if !server{tlsState:=tsNodeTLSSocketFor(socket);if tlsState==nil&&socket.kind==tsObjectKind{object:=(*tsObject)(socket.ref);tlsState=&tsNodeTLSSocketState{object:object,options:options};tsNodeTLS.sockets[object]=tlsState};if tlsState!=nil{if tlsState.sessionCache==nil{tlsState.sessionCache=&tsNodeTLSClientSessionCache{loop:loop,socket:socket}};config.ClientSessionCache=tlsState.sessionCache};config.KeyLogWriter=tsNodeTLSKeyLogSocketWriter{loop:loop,socket:socket,base:config.KeyLogWriter}}
	start:=func(){
		if raw.closed||raw.conn==nil{tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));return}
		raw.upgrading=true
		raw.upgradeContinue=func(){
			if raw.closed||raw.conn==nil{return}
			transport:=&tsNodeNetPrefixConn{Conn:raw.conn,prefix:append([]byte(nil),raw.upgradeData...)};raw.upgradeData=nil
			secure:=cryptotls.Client(transport,config);if server{secure=cryptotls.Server(transport,config)};raw.conn=secure
			var handshakeErr error
			task:=loop.submit(func()tsResult{handshakeErr=secure.Handshake();if handshakeErr!=nil{return tsResult{tsNodeTLSHandshakeFailure(loop,handshakeErr),true}};return tsResult{tsU,false}},func(result tsResult){
				if raw.closed{return}
				if result.rejected{if tsText(tsNodeTLSGet(loop,result.value,"code"))=="ECONNRESET"{tsNodeTLSOnConnectEnd(loop,socket,options);return};tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_hadError",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_closeAfterHandlingError",tsBooleanValue(true));tsNodeNetEmit(loop,socket,"error",result.value);tsNodeNetSocketDestroy(loop,socket,result.value);return}
				if wasConnecting:=raw.connecting;wasConnecting{raw.connecting=false;tsNodeNetSyncSocket(raw);tsNodeNetEmit(loop,socket,"connect")};if !tsNodeTLSHandshakeComplete(loop,socket,secure,options,server,config){return};tsNodeNetSyncSocket(raw);tsNodeNetStartRead(raw);tsNodeNetResetTimeout(raw)
			});task.unref=raw.unref
		}
		if raw.conn!=nil{_ = raw.conn.SetReadDeadline(time.Now())}
		if !raw.reading{tsNodeNetStartRead(raw)}
	}
	if raw.conn==nil&&raw.connecting{wrapped:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{start();return tsU}));failed:=tsFunctionValue(tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{failure:=tsArg(args,0);tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));if !tsTruthy(tsNodeTLSGet(loop,socket,"_hadError")){tsNodeTLSSet(loop,socket,"_hadError",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_closeAfterHandlingError",tsBooleanValue(true));destroy:=tsNodeTLSGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,failure)}};return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,socket,"once"),socket,tsNodeTLSKey("connect"),wrapped);tsCallReceiver(loop,tsNodeNetGet(loop,socket,"once"),socket,tsNodeTLSKey("error"),failed);return}
	start()
}
func tsNodeTLSCertificate(loop *tsLoop,certificate *x509.Certificate)tsValue{if certificate==nil{return tsNull};var subject tsValue;subject=tsNodeTLSObjectValue(loop,map[string]tsValue{"CN":tsNodeTLSKey(certificate.Subject.CommonName),"O":tsNodeTLSStringArray(loop,certificate.Subject.Organization),"OU":tsNodeTLSStringArray(loop,certificate.Subject.OrganizationalUnit),"C":tsNodeTLSStringArray(loop,certificate.Subject.Country),"ST":tsNodeTLSStringArray(loop,certificate.Subject.Province),"L":tsNodeTLSStringArray(loop,certificate.Subject.Locality)});var issuer tsValue;issuer=tsNodeTLSObjectValue(loop,map[string]tsValue{"CN":tsNodeTLSKey(certificate.Issuer.CommonName),"O":tsNodeTLSStringArray(loop,certificate.Issuer.Organization)});sans:=append([]string(nil),certificate.DNSNames...);for _,ip:=range certificate.IPAddresses{sans=append(sans,ip.String())};fingerprint:=sha256.Sum256(certificate.Raw);out:=tsNodeTLSObjectValue(loop,map[string]tsValue{"subject":subject,"issuer":issuer,"subjectaltname":tsNodeTLSKey(strings.Join(sans,", ")),"valid_from":tsNodeTLSKey(certificate.NotBefore.UTC().Format(time.RFC3339)),"valid_to":tsNodeTLSKey(certificate.NotAfter.UTC().Format(time.RFC3339)),"serialNumber":tsNodeTLSKey(strings.ToUpper(certificate.SerialNumber.Text(16))),"fingerprint256":tsNodeTLSKey(strings.ToUpper(hex.EncodeToString(fingerprint[:]))),"raw":tsNodeBuffer(certificate.Raw),"bits":tsNumberValue(float64(certificate.PublicKeyAlgorithm))});return out}
func tsNodeTLSObjectValue(loop *tsLoop,values map[string]tsValue)tsValue{object:=tsNewObject();for name,value:=range values{object.set(name,value)};return tsObjectValue(object)}
func tsNodeTLSStringArray(loop *tsLoop,values []string)tsValue{array:=&tsArray{};for _,value:=range values{array.values=append(array.values,tsNodeTLSKey(value))};return tsArrayValue(array)}
func tsNodeTLSVersionName(version uint16)string{switch version{case cryptotls.VersionTLS10:return "TLSv1";case cryptotls.VersionTLS11:return "TLSv1.1";case cryptotls.VersionTLS12:return "TLSv1.2";case cryptotls.VersionTLS13:return "TLSv1.3"};return cryptotls.VersionName(version)}
func tsNodeTLSSocketCipher(loop *tsLoop,receiver tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNull};connection:=state.conn.ConnectionState();name:=cryptotls.CipherSuiteName(connection.CipherSuite);version:=tsNodeTLSVersionName(connection.Version);return tsNodeTLSObjectValue(loop,map[string]tsValue{"name":tsNodeTLSKey(name),"standardName":tsNodeTLSKey(name),"version":tsNodeTLSKey(version)})}
func tsNodeTLSSocketPeer(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNodeTLSObjectValue(loop,nil)};certs:=state.conn.ConnectionState().PeerCertificates;if len(certs)==0{return tsNodeTLSObjectValue(loop,nil)};value:=tsNodeTLSCertificate(loop,certs[0]);if len(args)>0&&tsTruthy(args[0]){out:=(*tsObject)(value.ref);chain:=&tsArray{};for _,cert:=range certs{chain.values=append(chain.values,tsNodeTLSCertificate(loop,cert))};out.set("issuerCertificate",tsArrayValue(chain))};return value}
func tsNodeTLSSetupSocketPrototype(loop *tsLoop){
	if tsNodeTLS.socketPrototype!=nil{return}
	netModule:=tsNodeNetModule(loop);socketCtor:=tsNodeTLSGet(loop,netModule,"Socket")
	prototype:=tsNewObject();if fn:=tsGet(loop,socketCtor,tsNodeTLSKey("prototype"));fn.kind==tsObjectKind{prototype.prototype=fn};tsNodeTLS.socketPrototype=prototype
	prototype.set("getCipher",tsNamedNativeMethod("getCipher",0,tsNodeTLSSocketCipher));prototype.set("getProtocol",tsNamedNativeMethod("getProtocol",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNull};return tsNodeTLSKey(tsNodeTLSVersionName(state.conn.ConnectionState().Version))}))
	prototype.set("getPeerCertificate",tsNamedNativeMethod("getPeerCertificate",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeTLSSocketPeer(loop,receiver,args)}));prototype.set("getCertificate",tsNamedNativeMethod("getCertificate",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNodeTLSObjectValue(loop,nil)};certs:=state.conn.ConnectionState().PeerCertificates;if len(certs)==0{return tsNodeTLSObjectValue(loop,nil)};return tsNodeTLSCertificate(loop,certs[0])}))
	prototype.set("getSession",tsNamedNativeMethod("getSession",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||len(state.session)==0{return tsNull};return tsNodeBuffer(state.session)}));prototype.set("setSession",tsNamedNativeMethod("setSession",1,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{if state:=tsNodeTLSSocketFor(receiver);state!=nil{state.session=tsNodeTLSBytes(tsArg(args,0))};return tsU}));prototype.set("isSessionReused",tsNamedNativeMethod("isSessionReused",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsBooleanValue(false)}))
	prototype.set("getFinished",tsNamedNativeMethod("getFinished",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNull}));prototype.set("getPeerFinished",tsNamedNativeMethod("getPeerFinished",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNull}));prototype.set("getEphemeralKeyInfo",tsNamedNativeMethod("getEphemeralKeyInfo",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNull}));prototype.set("getSharedSigalgs",tsNamedNativeMethod("getSharedSigalgs",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsArrayValue(&tsArray{})}))
	prototype.set("exportKeyingMaterial",tsNamedNativeMethod("exportKeyingMaterial",3,func(_ *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{panic(tsThrown{tsFSException("Error","ERR_TLS_INVALID_STATE","TLS socket is not established")})};length:=int(tsNumber(tsArg(args,0)));if length<0||length>16*1024*1024{tsFSInvalid("length","a non-negative safe integer")};label:=tsText(tsArg(args,1));contextValue:=tsArg(args,2);var context []byte;if !tsIsUndefined(contextValue)&&!tsIsNull(contextValue){context=tsNodeTLSBytes(contextValue)};out,err:=state.conn.ConnectionState().ExportKeyingMaterial(label,context,length);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};return tsNodeBuffer(out)}))
	prototype.set("renegotiate",tsNamedNativeMethod("renegotiate",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{callback:=tsArg(args,1);err:=tsFSException("Error","ERR_TLS_RENEGOTIATION_DISABLED","TLS renegotiation is unavailable in Go's crypto/tls");if callback.kind==tsFunctionKind{tsCall(loop,callback,err)};return tsBooleanValue(false)}));prototype.set("disableRenegotiation",tsNamedNativeMethod("disableRenegotiation",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));prototype.set("setMaxSendFragment",tsNamedNativeMethod("setMaxSendFragment",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{n:=tsNumber(tsArg(args,0));return tsBooleanValue(n>=512&&n<=16384)}))
}
func tsNodeTLSSocketCreate(loop *tsLoop,socket,options tsValue)tsValue{
	tsNodeTLSSetupSocketPrototype(loop);if tsIsUndefined(socket)||tsIsNull(socket){socket=tsNodeNetSocketCreate(loop,options)};if socket.kind==tsObjectKind{(*tsObject)(socket.ref).prototype=tsObjectValue(tsNodeTLS.socketPrototype)}
	if tsNodeTLSObject(options){options=tsObjectValue(tsNodeTLSObjectCopy(options))}else{options=tsObjectValue(tsNewObject())};server:=tsTruthy(tsNodeTLSGet(loop,options,"isServer"));if !server{tsNodeTLSSet(loop,options,"_standaloneWrap",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_standaloneWrap",tsBooleanValue(true))};tsNodeTLSUpgradeSocket(loop,socket,options,server);return socket
}
func tsNodeTLSConnectOptions(loop *tsLoop,args []tsValue)(tsValue,tsValue){options:=tsObjectValue(tsNewObject());callback:=tsU;if len(args)>0&&tsNodeTLSObject(args[0]){options=tsObjectValue(tsNodeTLSObjectCopy(args[0]));if len(args)>1&&args[1].kind==tsFunctionKind{callback=args[1]};return options,callback};if len(args)>0&&tsIsNumeric(args[0]){options=tsObjectValue(tsNewObject());tsNodeTLSSet(loop,options,"port",args[0]);if len(args)>1&&args[1].kind==tsStringKind{tsNodeTLSSet(loop,options,"host",args[1])};for i:=1;i<len(args);i++{if tsNodeTLSObject(args[i]){for _,name:=range (*tsObject)(args[i].ref).order{tsNodeTLSSet(loop,options,name,(*tsObject)(args[i].ref).values[name])}};if args[i].kind==tsFunctionKind{callback=args[i]}};return options,callback};if len(args)>0&&args[0].kind==tsStringKind{options=tsObjectValue(tsNewObject());tsNodeTLSSet(loop,options,"path",args[0]);for i:=1;i<len(args);i++{if tsNodeTLSObject(args[i]){for _,name:=range (*tsObject)(args[i].ref).order{tsNodeTLSSet(loop,options,name,(*tsObject)(args[i].ref).values[name])}};if args[i].kind==tsFunctionKind{callback=args[i]}}};return options,callback}
func tsNodeTLSConnect(loop *tsLoop,args []tsValue)tsValue{
	options,callback:=tsNodeTLSConnectOptions(loop,args);hostValue:=tsNodeTLSGet(loop,options,"host");host:="localhost";if hostValue.kind==tsStringKind{host=tsText(hostValue)}
	servername:=tsNodeTLSGet(loop,options,"servername");if tsIsUndefined(servername)&&net.ParseIP(host)==nil{servername=tsNodeTLSKey(host);tsNodeTLSSet(loop,options,"servername",servername)};if servername.kind==tsStringKind&&net.ParseIP(tsText(servername))!=nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","Setting the TLS ServerName to an IP address is not permitted.")})}
	if _,err:=tsNodeTLSConfig(loop,options,false);err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};providedSocket:=tsNodeTLSGet(loop,options,"socket");socket:=tsNodeTLSSetupClientSocket(loop,options,providedSocket);if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,socket,"once"),socket,tsNodeTLSKey("secureConnect"),callback)}
	if !tsNullish(providedSocket){if state:=tsNodeNetSocketFor(socket);state!=nil&&(state.conn!=nil||state.connecting){tsNodeTLSUpgradeSocket(loop,socket,options,false);return socket};panic(tsThrown{tsFSException("Error","ERR_SOCKET_CLOSED","The supplied socket is not connected")})}
	state:=tsNodeNetSocketFor(socket);if state==nil{return socket};tsNodeNetSocketConnect(loop,socket,[]tsValue{options}...);tsNodeTLSUpgradeSocket(loop,socket,options,false);return socket
}
func tsNodeTLSSetupClientSocket(loop *tsLoop,options,socket tsValue)tsValue{tsNodeNetModule(loop);if tsIsNull(socket)||tsIsUndefined(socket){socket=tsNodeNetSocketCreate(loop,options)};tsNodeTLSSetupSocketPrototype(loop);if socket.kind==tsObjectKind{(*tsObject)(socket.ref).prototype=tsObjectValue(tsNodeTLS.socketPrototype)};tsNodeTLSSet(loop,socket,"encrypted",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_securePending",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_secureEstablished",tsBooleanValue(false));return socket}
func tsNodeTLSHandshakeTimeout(loop *tsLoop,options tsValue)time.Duration{
	value:=tsNodeTLSGet(loop,options,"handshakeTimeout");if tsIsUndefined(value)||tsIsNull(value){return 120*time.Second};if !tsIsNumeric(value)||tsNumber(value)<=0{return 0};return time.Duration(tsNumber(value))*time.Millisecond
}
func tsNodeTLSServerHandshake(loop *tsLoop,server,socket tsValue,conn *cryptotls.Conn,options tsValue){
	raw:=tsNodeNetSocketFor(socket);if raw==nil||conn==nil{return};raw.tlsHandshaking=true
	if socket.kind==tsObjectKind{object:=(*tsObject)(socket.ref);if tsNodeTLS.socketPrototype!=nil{object.prototype=tsObjectValue(tsNodeTLS.socketPrototype)};tsNodeTLS.sockets[object]=&tsNodeTLSSocketState{object:object,conn:conn,server:true,options:options}}
	tsNodeTLSSet(loop,socket,"encrypted",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"isServer",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"server",server);tsNodeTLSSet(loop,socket,"_requestCert",tsBooleanValue(tsTruthy(tsNodeTLSGet(loop,options,"requestCert"))));tsNodeTLSSet(loop,socket,"_rejectUnauthorized",tsBooleanValue(!tsStrictEqual(tsNodeTLSGet(loop,options,"rejectUnauthorized"),tsBooleanValue(false))));tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_securePending",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_secureEstablished",tsBooleanValue(false))
	timeout:=tsNodeTLSHandshakeTimeout(loop,options);var handshakeErr error
	task:=loop.submit(func()tsResult{
		if timeout>0{_ = conn.SetDeadline(time.Now().Add(timeout))}
		handshakeErr=conn.Handshake();_ = conn.SetDeadline(time.Time{})
		return tsResult{tsU,handshakeErr!=nil}
	},func(result tsResult){
		raw.tlsHandshaking=false;if raw.closed{return}
		if result.rejected{
			failure:=tsNodeTLSHandshakeFailure(loop,handshakeErr);if netError,ok:=handshakeErr.(net.Error);ok&&netError.Timeout(){failure=tsFSException("Error","ERR_TLS_HANDSHAKE_TIMEOUT","TLS handshake timeout")}
			tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_securePending",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_secureEstablished",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_hadError",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_closeAfterHandlingError",tsBooleanValue(true))
			tsNodeNetEmitSafe(loop,server,"tlsClientError",failure,socket);destroy:=tsNodeNetGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket)};return
		}
		tsNodeTLSAttachSocket(loop,socket,conn,true,options);tsNodeTLSSet(loop,socket,"isServer",tsBooleanValue(true))
		if tsTruthy(tsNodeTLSGet(loop,options,"requestCert")){var verifyErr error;if secure:=tsNodeTLSServerFor(server);secure!=nil&&secure.config!=nil{verifyErr=tsNodeTLSVerifyClientChain(conn.ConnectionState(),secure.config.ClientCAs)}else{verifyErr=errors.New("peer did not return a certificate")};if verifyErr!=nil{authorization:=tsNodeTLSVerifyFailure(verifyErr);code:=tsNodeTLSGet(loop,authorization,"code");if tsNullish(code){code=tsNodeTLSGet(loop,authorization,"message")};tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"authorizationError",code)}else{tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorizationError",tsNull)}}
		tsNodeNetEmitSafe(loop,server,"secureConnection",socket)
		if !raw.closed&&!tsTruthy(tsNodeNetGet(loop,server,"pauseOnConnect")){tsNodeNetStartRead(raw)}
	})
	task.unref=raw.unref
}
func tsNodeTLSServerListen(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{return tsNodeNetListen(loop,receiver,args...)}
func tsNodeTLSCreateServer(loop *tsLoop,args []tsValue)tsValue{options,listener:=tsArg(args,0),tsArg(args,1);if options.kind==tsFunctionKind{listener=options;options=tsObjectValue(tsNewObject())};config,err:=tsNodeTLSConfig(loop,options,true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};tsNodeNetModule(loop);server:=tsNodeNetServerCreate(loop,options);object:=(*tsObject)(server.ref);if object.prototype!=nil{if tsNodeTLS.serverPrototype==nil{prototype:=tsNewObject();prototype.prototype=object.prototype;tsNodeTLS.serverPrototype=prototype};object.prototype=tsObjectValue(tsNodeTLS.serverPrototype)};secure:=&tsNodeTLSServerState{object:object,loop:loop,server:server,config:config,options:options,contexts:map[string]*cryptotls.Config{}};secure.config=tsNodeTLSServerConfig(secure,config);tsNodeTLS.servers[object]=secure;serverListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{socket:=tsArg(values,0);if raw:=tsNodeNetSocketFor(socket);raw!=nil{if conn,ok:=raw.conn.(*cryptotls.Conn);ok{tsNodeTLSServerHandshake(loop,server,socket,conn,options);return tsU}};tsNodeNetEmit(loop,server,"secureConnection",socket);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,server,"on"),server,tsNodeTLSKey("connection"),serverListener);if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,server,"on"),server,tsNodeTLSKey("secureConnection"),listener)}
	if tsNodeTLS.serverPrototype!=nil{tsNodeTLS.serverProperties.set("prototype",tsObjectValue(tsNodeTLS.serverPrototype));tsNodeTLS.serverPrototype.set("listen",tsNamedNativeMethod("listen",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{return tsNodeTLSServerListen(loop,receiver,values)}));tsNodeTLS.serverPrototype.set("setSecureContext",tsNamedNativeMethod("setSecureContext",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{current:=tsNodeTLSServerFor(receiver);if current!=nil{config,err:=tsNodeTLSConfig(loop,tsArg(values,0),true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};current.mu.Lock();current.options=tsArg(values,0);current.config=tsNodeTLSServerConfig(current,config);current.mu.Unlock()};return tsU}));tsNodeTLS.serverPrototype.set("addContext",tsNamedNativeMethod("addContext",2,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{current:=tsNodeTLSServerFor(receiver);if current!=nil{config,err:=tsNodeTLSConfig(loop,tsArg(values,1),true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};current.mu.Lock();current.contexts[tsNodeTLSServerNameKey(tsText(tsArg(values,0)))]=config;current.mu.Unlock()};return tsU}));tsNodeTLS.serverPrototype.set("getTicketKeys",tsNamedNativeMethod("getTicketKeys",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeBuffer(make([]byte,48))}));tsNodeTLS.serverPrototype.set("setTicketKeys",tsNamedNativeMethod("setTicketKeys",1,func(_ *tsLoop,_ tsValue,values ...tsValue)tsValue{data:=tsNodeTLSBytes(tsArg(values,0));if len(data)!=48{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","Session ticket keys must be a 48-byte buffer")})};return tsU}))}
	return server}
func tsNodeTLSParseCertString(loop *tsLoop,value tsValue)tsValue{out:=tsNewObject();for _,part:=range strings.Split(tsText(value),"\n"){part=strings.TrimSpace(part);if index:=strings.IndexByte(part,'=');index>0{out.set(strings.TrimSpace(part[:index]),tsNodeTLSKey(strings.TrimSpace(part[index+1:])))}};return tsObjectValue(out)}
func tsNodeTLSCheckServerIdentity(loop *tsLoop,hostValue,cert tsValue)tsValue{host:=tsText(hostValue);if tsNodeTLSObject(cert){sans:=tsText(tsNodeTLSGet(loop,cert,"subjectaltname"));for _,entry:=range strings.Split(sans,", "){entry=strings.TrimPrefix(entry,"DNS:");entry=strings.TrimPrefix(entry,"IP Address:");if entry==host{return tsU};if strings.HasPrefix(entry,"*."){suffix:=entry[1:];if strings.HasSuffix(host,suffix)&&strings.Count(host,".")==strings.Count(entry,"."){return tsU}}};subject:=tsNodeTLSGet(loop,cert,"subject");if tsText(tsNodeTLSGet(loop,subject,"CN"))==host{return tsU}};return tsFSException("Error","ERR_TLS_CERT_ALTNAME_INVALID","Hostname/IP does not match certificate's altnames: "+host)}
func tsNodeTLSRootCerts()[]string{for _,path:=range []string{"/etc/ssl/certs/ca-certificates.crt","/etc/pki/tls/certs/ca-bundle.crt","/etc/ssl/ca-bundle.pem"}{if data,err:=os.ReadFile(path);err==nil{out:=[]string{};remaining:=data;for len(remaining)>0{block,rest:=pem.Decode(remaining);if block==nil{break};if block.Type=="CERTIFICATE"{out=append(out,string(pem.EncodeToMemory(block)))};remaining=rest};if len(out)>0{return out}}};return nil}
func tsNodeTLSExtraRootCerts()[]string{path:=os.Getenv("NODE_EXTRA_CA_CERTS");if path==""{return nil};data,err:=os.ReadFile(path);if err!=nil{return nil};out:=[]string{};for len(data)>0{block,rest:=pem.Decode(data);if block==nil{break};data=rest;if block.Type!="CERTIFICATE"{continue};certificate,parseErr:=x509.ParseCertificate(block.Bytes);if parseErr==nil{out=append(out,string(pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE",Bytes:certificate.Raw})))}};return out}
func tsNodeTLSDefaultCiphers()string{names:=[]string{"TLS_AES_256_GCM_SHA384","TLS_CHACHA20_POLY1305_SHA256","TLS_AES_128_GCM_SHA256"};legacy:="";if tsNodeTLS.defaultCiphers.kind==tsStringKind{legacy=tsText(tsNodeTLS.defaultCiphers)}else{values:=[]string{};for _,suite:=range cryptotls.CipherSuites(){values=append(values,suite.Name)};legacy=strings.Join(values,":")};if legacy!=""{names=append(names,legacy)};return strings.Join(names,":")}
func tsNodeTLSModule(loop *tsLoop)tsValue{tsNodeTLS.Do(func(){tsNodeNetModule(loop);tsNodeTLS.secure=map[*tsObject]*tsNodeTLSSecureState{};tsNodeTLS.sockets=map[*tsObject]*tsNodeTLSSocketState{};tsNodeTLS.servers=map[*tsObject]*tsNodeTLSServerState{};tsNodeTLS.minVersion=tsNodeTLSKey("TLSv1.2");tsNodeTLS.maxVersion=tsNodeTLSKey("TLSv1.3");tsNodeTLS.defaultCurve=tsNodeTLSKey("auto");tsNodeTLS.defaultCAs=nil;tsNodeTLSSetupSocketPrototype(loop)
		secure:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))});secure.name="SecureContext";secure.constructible=true;secure.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))};secure.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))};secure.properties=tsNewObject();secure.properties.set("parsePkcs12",tsNamedNativeMethod("parsePkcs12",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSPKCS12Parse(loop,args)}));tsNodeTLS.secureConstructor=tsFunctionValue(secure)
		socket:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSocketCreate(loop,tsArg(args,0),tsArg(args,1))});socket.name="TLSSocket";socket.constructible=true;socket.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSocketCreate(loop,tsArg(args,0),tsArg(args,1))};socket.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSocketCreate(loop,tsArg(args,0),tsArg(args,1))};socket.properties=tsNewObject();socket.properties.set("prototype",tsObjectValue(tsNodeTLS.socketPrototype));tsNodeTLS.socketConstructor=tsFunctionValue(socket)
		server:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)});server.name="Server";server.constructible=true;server.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)};server.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)};server.properties=tsNewObject();tsNodeTLS.serverProperties=server.properties;tsNodeTLS.serverConstructor=tsFunctionValue(server)
		module:=tsNewObject();module.set("CLIENT_RENEG_LIMIT",tsNumberValue(3));module.set("CLIENT_RENEG_WINDOW",tsNumberValue(600));module.set("connect",tsNamedNativeMethod("connect",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSConnect(loop,args)}));module.set("createConnection",module.get("connect"));module.set("createSecureContext",tsNamedNativeMethod("createSecureContext",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))}));module.set("SecureContext",tsNodeTLS.secureConstructor);module.set("TLSSocket",tsNodeTLS.socketConstructor);module.set("Server",tsNodeTLS.serverConstructor);module.set("createServer",tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)}));module.set("DEFAULT_CIPHERS",tsNodeTLSKey(tsNodeTLSDefaultCiphers()));module.set("DEFAULT_ECDH_CURVE",tsNodeTLS.defaultCurve);module.set("DEFAULT_MIN_VERSION",tsNodeTLS.minVersion);module.set("DEFAULT_MAX_VERSION",tsNodeTLS.maxVersion);if module.descriptors==nil{module.descriptors=map[string]*tsDescriptor{}};module.descriptors["DEFAULT_CIPHERS"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_CIPHERS",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLSKey(tsNodeTLSDefaultCiphers())}),setter:tsNamedNativeMethod("set DEFAULT_CIPHERS",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if tsTruthy(value){if value.kind!=tsStringKind{tsFSInvalid("value","a string")};filtered:=tsInternalTLSStripCipherNames(value);value=tsInternalTLSGet(loop,filtered,"cipherList")}else{value=tsNodeTLSKey("")};tsNodeTLS.defaultCiphers=value;return tsU})};module.descriptors["DEFAULT_ECDH_CURVE"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_ECDH_CURVE",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLS.defaultCurve}),setter:tsNamedNativeMethod("set DEFAULT_ECDH_CURVE",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeTLS.defaultCurve=tsArg(args,0);return tsU})};module.descriptors["DEFAULT_MIN_VERSION"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_MIN_VERSION",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLS.minVersion}),setter:tsNamedNativeMethod("set DEFAULT_MIN_VERSION",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeTLS.minVersion=tsArg(args,0);return tsU})};module.descriptors["DEFAULT_MAX_VERSION"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_MAX_VERSION",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLS.maxVersion}),setter:tsNamedNativeMethod("set DEFAULT_MAX_VERSION",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeTLS.maxVersion=tsArg(args,0);return tsU})}
		module.set("convertALPNProtocols",tsNamedNativeMethod("convertALPNProtocols",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSALPNValue(loop,tsNodeTLSALPN(loop,tsArg(args,0)))}));module.set("getCiphers",tsNamedNativeMethod("getCiphers",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{names:=[]string{"TLS_AES_128_GCM_SHA256","TLS_AES_256_GCM_SHA384","TLS_CHACHA20_POLY1305_SHA256"};for _,suite:=range append(cryptotls.CipherSuites(),cryptotls.InsecureCipherSuites()...){names=append(names,strings.ToLower(suite.Name))};sort.Strings(names);return tsNodeTLSStringArray(loop,names)}));module.set("getDefaultCiphers",tsNamedNativeMethod("getDefaultCiphers",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLSKey(tsNodeTLSDefaultCiphers())}));module.set("parseCertString",tsNamedNativeMethod("parseCertString",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSParseCertString(loop,tsArg(args,0))}));module.set("checkServerIdentity",tsNamedNativeMethod("checkServerIdentity",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSCheckServerIdentity(loop,tsArg(args,0),tsArg(args,1))}))
		module.set("getCACertificates",tsNamedNativeMethod("getCACertificates",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{kind:=tsArg(args,0);if tsIsUndefined(kind){kind=tsNodeTLSKey("default")};if kind.kind!=tsStringKind{tsFSInvalid("type","a string")};var certs []string;switch tsText(kind){case "default":if tsNodeTLS.defaultCAsSet{out:=&tsArray{};for _,cert:=range tsNodeTLS.defaultCAs{out.values=append(out.values,cert)};return tsArrayValue(out)};certs=tsNodeTLSRootCerts();if os.Getenv("NODE_EXTRA_CA_CERTS")!=""{certs=append(certs,tsNodeTLSExtraRootCerts()...)};case "bundled","system":certs=tsNodeTLSRootCerts();case "extra":certs=tsNodeTLSExtraRootCerts();default:panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument \"type\" is invalid")})};out:=&tsArray{};for _,cert:=range certs{out.values=append(out.values,tsNodeTLSKey(cert))};return tsArrayValue(out)}));module.set("setDefaultCACertificates",tsNamedNativeMethod("setDefaultCACertificates",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsArrayKind{tsFSInvalid("certs","an instance of Array")};array:=(*tsArray)(value.ref);for index,item:=range array.values{if item.kind!=tsStringKind&&item.kind!=tsTypedArrayKind&&item.kind!=tsDataViewKind{tsFSInvalid("certs["+strconv.Itoa(index)+"]","string or an instance of ArrayBufferView")}};normalized,err:=tsNodeTLSParseCACertificates(loop,value);if err!=nil{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED",err.Error())})};if len(normalized)==0&&array.length()>0{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED","No valid certificates found in the provided array")})};tsNodeTLS.defaultCAs=normalized;tsNodeTLS.defaultCAsSet=true;return tsU}));roots:=&tsArray{};for _,cert:=range tsNodeTLSRootCerts(){roots.values=append(roots.values,tsNodeTLSKey(cert))};module.set("rootCertificates",tsArrayValue(roots));module.set("default",tsObjectValue(module));tsNodeTLS.module=tsObjectValue(module)
	});return tsNodeTLS.module}
`

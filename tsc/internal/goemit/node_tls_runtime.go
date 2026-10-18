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
	tsInternalTLSValidateContext(loop,options)
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
		config.RootCAs,_=x509.SystemCertPool();if config.RootCAs==nil{config.RootCAs=x509.NewCertPool()}
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
	suites:=tsNodeTLSGet(loop,options,"ciphers");if tsIsUndefined(suites)&&tsNodeTLS.defaultCiphers.kind==tsStringKind{suites=tsNodeTLS.defaultCiphers};if suites.kind==tsStringKind{if tsText(suites)==""{config.MinVersion=max(config.MinVersion,cryptotls.VersionTLS13)}else{for _,name:=range strings.Split(tsText(suites),":"){if id,ok:=tsNodeTLSCipherID(name);ok{config.CipherSuites=append(config.CipherSuites,id)}else if strings.Contains(name,"PSK-"){return nil,fmt.Errorf("cipher %s is not supported by Go crypto/tls",name)}}}}
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
func tsNodeTLSEmitSecureConnectDone(loop *tsLoop,socket tsValue){symbols:=tsInternalNetSymbolsModule();event:=tsInternalNetGet(loop,symbols,"kSecureConnectDone");emit:=tsNodeTLSGet(loop,socket,"emit");if event.kind==tsSymbolKind&&emit.kind==tsFunctionKind{tsCallReceiver(loop,emit,socket,event)}}
func tsNodeTLSVerifyClientChain(state cryptotls.ConnectionState,roots *x509.CertPool)error{
	if len(state.PeerCertificates)==0{return errors.New("peer did not return a certificate")}
	intermediates:=x509.NewCertPool();for _,certificate:=range state.PeerCertificates[1:]{intermediates.AddCert(certificate)}
	_,err:=state.PeerCertificates[0].Verify(x509.VerifyOptions{Roots:roots,Intermediates:intermediates,KeyUsages:[]x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}});return err
}
func tsNodeTLSConnectionIdentityError(loop *tsLoop,socket tsValue,conn *cryptotls.Conn,options tsValue,config *cryptotls.Config,server,standalone bool)tsValue{
	if server||conn==nil{return tsNull};state:=conn.ConnectionState();if len(state.PeerCertificates)==0{return tsNull}
	if len(state.VerifiedChains)==0{if err:=tsNodeTLSVerifyChain(state,config.RootCAs);err!=nil{return tsNodeTLSVerifyFailure(err)}}
	if standalone{return tsNull}
	hostValue:=tsNodeTLSGet(loop,options,"servername");if tsNullish(hostValue){hostValue=tsNodeTLSGet(loop,options,"host")};if tsNullish(hostValue){if netState:=tsNodeNetSocketFor(socket);netState!=nil&&netState.remoteAddress!=""{hostValue=tsNodeTLSKey(netState.remoteAddress)}else{hostValue=tsNodeTLSKey("localhost")}}
	certificate:=tsNodeTLSCertificate(loop,state.PeerCertificates[0]);if check:=tsNodeTLSGet(loop,options,"checkServerIdentity");check.kind==tsFunctionKind{var result tsValue;func(){defer func(){if thrown:=recover();thrown!=nil{result=tsUnwrap(thrown)}}();result=tsCallReceiver(loop,check,options,hostValue,certificate)}();if tsTruthy(result){return result};return tsNull}
	return tsNodeTLSCheckServerIdentity(loop,hostValue,certificate)
}
func tsNodeTLSHandshakeComplete(loop *tsLoop,socket tsValue,conn *cryptotls.Conn,options tsValue,server bool,config *cryptotls.Config)bool{
	standalone:=tsTruthy(tsNodeTLSGet(loop,socket,"_standaloneWrap"));verifyError:=tsNodeTLSConnectionIdentityError(loop,socket,conn,options,config,server,standalone)
	tsNodeTLSAttachSocket(loop,socket,conn,server,options);rejectUnauthorized:=tsNodeTLSGet(loop,options,"rejectUnauthorized");reject:=tsStrictEqual(rejectUnauthorized,tsBooleanValue(true))||!server&&tsIsUndefined(rejectUnauthorized)
	if !server&&tsIsUndefined(rejectUnauthorized){if value,ok:=os.LookupEnv("NODE_TLS_REJECT_UNAUTHORIZED");ok&&value=="0"{reject=false}}
	if server{tsNodeTLSEmitSecureConnectDone(loop,socket);tsNodeNetEmit(loop,socket,"secure",socket);return true}
	if standalone{if tsTruthy(verifyError)&&reject{destroy:=tsNodeTLSGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,verifyError)};return false};tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeTLSEmitSecureConnectDone(loop,socket);tsNodeTLSFlushPendingSession(loop,socket);tsNodeNetEmit(loop,socket,"secure",socket);return true}
	if tsTruthy(verifyError){tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));code:=tsNodeTLSGet(loop,verifyError,"code");if tsNullish(code){code=tsNodeTLSGet(loop,verifyError,"message")};tsNodeTLSSet(loop,socket,"authorizationError",code);if reject{destroy:=tsNodeTLSGet(loop,socket,"destroy");if destroy.kind==tsFunctionKind{tsCallReceiver(loop,destroy,socket,verifyError)};tsNodeNetEmit(loop,socket,"secure",socket);return false}}else{tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorizationError",tsNull)}
	tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(false));tsNodeTLSEmitSecureConnectDone(loop,socket);tsNodeNetEmit(loop,socket,"secureConnect",verifyError);tsNodeTLSFlushPendingSession(loop,socket);tsNodeNetEmit(loop,socket,"secure",socket);return true
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
func tsNodeTLSALPN(loop *tsLoop,value tsValue)[]string{
	if value.kind==tsArrayKind{out:=[]string{};for index,item:=range (*tsArray)(value.ref).values{protocol:=tsText(item);length:=len([]byte(protocol));if length>255{message:=fmt.Sprintf("The byte length of the protocol at index %d exceeds the maximum length. It must be <= 255. Received %d",index,length);panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",message)})};out=append(out,protocol)};return out}
	if value.kind==tsTypedArrayKind||value.kind==tsDataViewKind||value.kind==tsArrayBufferKind{bytes:=tsFSBytes(value);out:=[]string{};for i:=0;i<len(bytes);{size:=int(bytes[i]);i++;if size==0||i+size>len(bytes){break};out=append(out,string(bytes[i:i+size]));i+=size};return out}
	if value.kind==tsStringKind{return []string{tsText(value)}};return nil
}
func tsNodeTLSConvertALPNProtocols(loop *tsLoop,protocols,out tsValue)tsValue{
	if protocols.kind==tsArrayKind{var encoded []byte;for index,item:=range (*tsArray)(protocols.ref).values{if item.kind!=tsStringKind{tsInternalValidatorInvalidType("protocols["+strconv.Itoa(index)+"]","a string",item)};protocol:=[]byte(tsText(item));if len(protocol)>255{message:=fmt.Sprintf("The byte length of the protocol at index %d exceeds the maximum length. It must be <= 255. Received %d",index,len(protocol));panic(tsThrown{tsFSException("RangeError","ERR_OUT_OF_RANGE",message)})};encoded=append(encoded,byte(len(protocol)));encoded=append(encoded,protocol...)};tsNodeTLSSet(loop,out,"ALPNProtocols",tsNodeBuffer(encoded))
	}else if protocols.kind==tsTypedArrayKind||protocols.kind==tsDataViewKind{tsNodeTLSSet(loop,out,"ALPNProtocols",tsNodeBuffer(tsNodeTLSBytes(protocols)))}
	return tsU
}
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
func tsNodeTLSCertificate(loop *tsLoop,certificate *x509.Certificate)tsValue{if certificate==nil{return tsNull};var subject tsValue;subject=tsNodeTLSObjectValue(loop,map[string]tsValue{"CN":tsNodeTLSKey(certificate.Subject.CommonName),"O":tsNodeTLSStringArray(loop,certificate.Subject.Organization),"OU":tsNodeTLSStringArray(loop,certificate.Subject.OrganizationalUnit),"C":tsNodeTLSStringArray(loop,certificate.Subject.Country),"ST":tsNodeTLSStringArray(loop,certificate.Subject.Province),"L":tsNodeTLSStringArray(loop,certificate.Subject.Locality)});var issuer tsValue;issuer=tsNodeTLSObjectValue(loop,map[string]tsValue{"CN":tsNodeTLSKey(certificate.Issuer.CommonName),"O":tsNodeTLSStringArray(loop,certificate.Issuer.Organization)});sans:=append([]string(nil),certificate.DNSNames...);for _,ip:=range certificate.IPAddresses{sans=append(sans,ip.String())};fingerprint:=sha256.Sum256(certificate.Raw);bits:=0;switch public:=certificate.PublicKey.(type){case *rsa.PublicKey:bits=public.N.BitLen();case *ecdsa.PublicKey:bits=public.Params().BitSize;case ed25519.PublicKey:bits=256};out:=tsNodeTLSObjectValue(loop,map[string]tsValue{"subject":subject,"issuer":issuer,"subjectaltname":tsNodeTLSKey(strings.Join(sans,", ")),"valid_from":tsNodeTLSKey(certificate.NotBefore.UTC().Format(time.RFC3339)),"valid_to":tsNodeTLSKey(certificate.NotAfter.UTC().Format(time.RFC3339)),"serialNumber":tsNodeTLSKey(strings.ToUpper(certificate.SerialNumber.Text(16))),"fingerprint256":tsNodeTLSKey(strings.ToUpper(hex.EncodeToString(fingerprint[:]))),"raw":tsNodeBuffer(certificate.Raw),"bits":tsNumberValue(float64(bits))});return out}
func tsNodeTLSObjectValue(loop *tsLoop,values map[string]tsValue)tsValue{object:=tsNewObject();for name,value:=range values{object.set(name,value)};return tsObjectValue(object)}
func tsNodeTLSStringArray(loop *tsLoop,values []string)tsValue{array:=&tsArray{};for _,value:=range values{array.values=append(array.values,tsNodeTLSKey(value))};return tsArrayValue(array)}
func tsNodeTLSVersionName(version uint16)string{switch version{case cryptotls.VersionTLS10:return "TLSv1";case cryptotls.VersionTLS11:return "TLSv1.1";case cryptotls.VersionTLS12:return "TLSv1.2";case cryptotls.VersionTLS13:return "TLSv1.3"};return cryptotls.VersionName(version)}
func tsNodeTLSCipherID(name string)(uint16,bool){
	switch name{
	case "DES-CBC3-SHA":return cryptotls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,true
	case "AES128-SHA":return cryptotls.TLS_RSA_WITH_AES_128_CBC_SHA,true
	case "AES256-SHA":return cryptotls.TLS_RSA_WITH_AES_256_CBC_SHA,true
	case "AES128-GCM-SHA256":return cryptotls.TLS_RSA_WITH_AES_128_GCM_SHA256,true
	case "AES256-GCM-SHA384":return cryptotls.TLS_RSA_WITH_AES_256_GCM_SHA384,true
	case "ECDHE-ECDSA-AES128-SHA":return cryptotls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,true
	case "ECDHE-ECDSA-AES256-SHA":return cryptotls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,true
	case "ECDHE-RSA-AES128-SHA":return cryptotls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,true
	case "ECDHE-RSA-AES256-SHA":return cryptotls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,true
	case "ECDHE-ECDSA-AES128-SHA256":return cryptotls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,true
	case "ECDHE-RSA-AES128-SHA256":return cryptotls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256,true
	case "ECDHE-ECDSA-AES128-GCM-SHA256":return cryptotls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,true
	case "ECDHE-ECDSA-AES256-GCM-SHA384":return cryptotls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,true
	case "ECDHE-RSA-AES128-GCM-SHA256":return cryptotls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,true
	case "ECDHE-RSA-AES256-GCM-SHA384":return cryptotls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,true
	case "ECDHE-RSA-CHACHA20-POLY1305":return cryptotls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,true
	case "ECDHE-ECDSA-CHACHA20-POLY1305":return cryptotls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,true
	}
	for _,suite:=range append(cryptotls.CipherSuites(),cryptotls.InsecureCipherSuites()...){if suite.Name==name{return suite.ID,true}}
	return 0,false
}
func tsNodeTLSOpenSSLCipherName(id uint16)string{
	switch id{
	case cryptotls.TLS_RSA_WITH_3DES_EDE_CBC_SHA:return "DES-CBC3-SHA"
	case cryptotls.TLS_RSA_WITH_AES_128_CBC_SHA:return "AES128-SHA"
	case cryptotls.TLS_RSA_WITH_AES_256_CBC_SHA:return "AES256-SHA"
	case cryptotls.TLS_RSA_WITH_AES_128_GCM_SHA256:return "AES128-GCM-SHA256"
	case cryptotls.TLS_RSA_WITH_AES_256_GCM_SHA384:return "AES256-GCM-SHA384"
	case cryptotls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA:return "ECDHE-ECDSA-AES128-SHA"
	case cryptotls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA:return "ECDHE-ECDSA-AES256-SHA"
	case cryptotls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA:return "ECDHE-RSA-AES128-SHA"
	case cryptotls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA:return "ECDHE-RSA-AES256-SHA"
	case cryptotls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256:return "ECDHE-ECDSA-AES128-SHA256"
	case cryptotls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256:return "ECDHE-RSA-AES128-SHA256"
	case cryptotls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256:return "ECDHE-ECDSA-AES128-GCM-SHA256"
	case cryptotls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384:return "ECDHE-ECDSA-AES256-GCM-SHA384"
	case cryptotls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256:return "ECDHE-RSA-AES128-GCM-SHA256"
	case cryptotls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384:return "ECDHE-RSA-AES256-GCM-SHA384"
	case cryptotls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305:return "ECDHE-RSA-CHACHA20-POLY1305"
	case cryptotls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305:return "ECDHE-ECDSA-CHACHA20-POLY1305"
	}
	return ""
}
func tsNodeTLSLocalCertificate(loop *tsLoop,state *tsNodeTLSSocketState)*x509.Certificate{if state==nil{return nil};config,err:=tsNodeTLSConfig(loop,state.options,state.server);if err!=nil||len(config.Certificates)==0{return nil};certificate:=config.Certificates[0];if certificate.Leaf!=nil{return certificate.Leaf};if len(certificate.Certificate)>0{leaf,_:=x509.ParseCertificate(certificate.Certificate[0]);return leaf};return nil}
func tsNodeTLSX509Value(loop *tsLoop,certificates []*x509.Certificate,index int)tsValue{if index<0||index>=len(certificates)||certificates[index]==nil{return tsU};tsNodeCryptoModule(loop);value:=tsNodeCryptoX509Create(loop,tsNodeBuffer(certificates[index].Raw));if index+1<len(certificates){issuer:=tsNodeTLSX509Value(loop,certificates,index+1);if value.kind==tsObjectKind{object:=(*tsObject)(value.ref);object.set("issuerCertificate",issuer);if object.descriptors==nil{object.descriptors=map[string]*tsDescriptor{}};object.descriptors["issuerCertificate"]=&tsDescriptor{enumerable:false,configurable:true}}};return value}
func tsNodeTLSSocketCipher(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsU};connection:=state.conn.ConnectionState();standard:=cryptotls.CipherSuiteName(connection.CipherSuite);name:=tsNodeTLSOpenSSLCipherName(connection.CipherSuite);if name==""{name=standard};version:=tsNodeTLSVersionName(connection.Version);return tsNodeTLSObjectValue(loop,map[string]tsValue{"name":tsNodeTLSKey(name),"standardName":tsNodeTLSKey(standard),"version":tsNodeTLSKey(version)})}
func tsNodeTLSSocketPeer(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNull};certs:=state.conn.ConnectionState().PeerCertificates;if len(certs)==0{return tsNodeTLSObjectValue(loop,nil)};value:=tsNodeTLSCertificate(loop,certs[0]);if len(args)>0&&tsTruthy(args[0]){chain:=make([]tsValue,len(certs));for i,cert:=range certs{chain[i]=tsNodeTLSCertificate(loop,cert)};for index:=len(chain)-2;index>=0;index--{object:=(*tsObject)(chain[index].ref);object.set("issuerCertificate",chain[index+1])};return chain[0]};return value}
func tsNodeTLSSocketLocalLegacy(loop *tsLoop,receiver tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNull};certificate:=tsNodeTLSLocalCertificate(loop,state);if certificate==nil{return tsNodeTLSObjectValue(loop,nil)};return tsNodeTLSCertificate(loop,certificate)}
func tsNodeTLSSocketPeerX509(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsU};return tsNodeTLSX509Value(loop,state.conn.ConnectionState().PeerCertificates,0)}
func tsNodeTLSSocketX509(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsU};certificate:=tsNodeTLSLocalCertificate(loop,state);if certificate==nil{return tsU};return tsNodeTLSX509Value(loop,[]*x509.Certificate{certificate},0)}
func tsNodeTLSSocketEphemeral(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsU};var name string;size:=0;switch state.conn.ConnectionState().CurveID{case cryptotls.X25519:name="X25519";size=253;case cryptotls.CurveP256:name="prime256v1";size=256;case cryptotls.CurveP384:name="secp384r1";size=384;case cryptotls.CurveP521:name="secp521r1";size=521};if name==""{return tsU};return tsNodeTLSObjectValue(loop,map[string]tsValue{"type":tsNodeTLSKey("ECDH"),"name":tsNodeTLSKey(name),"size":tsNumberValue(float64(size))})}
func tsNodeTLSSetupSocketPrototype(loop *tsLoop){
	if tsNodeTLS.socketPrototype!=nil{return}
	netModule:=tsNodeNetModule(loop);socketCtor:=tsNodeTLSGet(loop,netModule,"Socket")
	prototype:=tsNewObject();if fn:=tsGet(loop,socketCtor,tsNodeTLSKey("prototype"));fn.kind==tsObjectKind{prototype.prototype=fn};tsNodeTLS.socketPrototype=prototype
	prototype.set("getCipher",tsNamedNativeMethod("getCipher",0,tsNodeTLSSocketCipher));prototype.set("getProtocol",tsNamedNativeMethod("getProtocol",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsNull};return tsNodeTLSKey(tsNodeTLSVersionName(state.conn.ConnectionState().Version))}))
	prototype.set("getPeerCertificate",tsNamedNativeMethod("getPeerCertificate",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{return tsNodeTLSSocketPeer(loop,receiver,args)}));prototype.set("getCertificate",tsNamedNativeMethod("getCertificate",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return tsNodeTLSSocketLocalLegacy(loop,receiver)}))
	prototype.set("getSession",tsNamedNativeMethod("getSession",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||len(state.session)==0{return tsU};return tsNodeBuffer(state.session)}));prototype.set("setSession",tsNamedNativeMethod("setSession",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{if state:=tsNodeTLSSocketFor(receiver);state!=nil{session:=tsArg(args,0);if session.kind==tsStringKind{session=tsNodeBufferFrom(loop,[]tsValue{session,tsNodeTLSKey("latin1")})};state.session=tsNodeTLSBytes(session)};return tsU}));prototype.set("isSessionReused",tsNamedNativeMethod("isSessionReused",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);return tsBooleanValue(state!=nil&&state.conn!=nil&&state.conn.ConnectionState().DidResume)}))
	prototype.set("getFinished",tsNamedNativeMethod("getFinished",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));prototype.set("getPeerFinished",tsNamedNativeMethod("getPeerFinished",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));prototype.set("getEphemeralKeyInfo",tsNamedNativeMethod("getEphemeralKeyInfo",0,tsNodeTLSSocketEphemeral));prototype.set("getSharedSigalgs",tsNamedNativeMethod("getSharedSigalgs",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));prototype.set("getTLSTicket",tsNamedNativeMethod("getTLSTicket",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}));prototype.set("getPeerX509Certificate",tsNamedNativeMethod("getPeerX509Certificate",0,tsNodeTLSSocketPeerX509));prototype.set("getX509Certificate",tsNamedNativeMethod("getX509Certificate",0,tsNodeTLSSocketX509))
	prototype.set("exportKeyingMaterial",tsNamedNativeMethod("exportKeyingMaterial",3,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{panic(tsThrown{tsFSException("Error","ERR_TLS_INVALID_STATE","TLS socket is not established")})};lengthValue,labelValue,contextValue:=tsArg(args,0),tsArg(args,1),tsArg(args,2);if !tsIsNumeric(lengthValue)||math.Trunc(tsNumber(lengthValue))!=tsNumber(lengthValue)||tsNumber(lengthValue)<0||tsNumber(lengthValue)>16*1024*1024{tsFSInvalid("length","a non-negative safe integer")};if labelValue.kind!=tsStringKind{tsInternalValidatorInvalidType("label","a string",labelValue)};if !tsIsUndefined(contextValue){tsInternalValidatorBuffer(contextValue,tsNodeTLSKey("context"))};length:=int(tsNumber(lengthValue));label:=tsText(labelValue);var context []byte;if !tsIsUndefined(contextValue){context=tsNodeTLSBytes(contextValue)};connectionState:=state.conn.ConnectionState();out,err:=connectionState.ExportKeyingMaterial(label,context,length);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};return tsNodeBuffer(out)}))
	prototype.set("renegotiate",tsNamedNativeMethod("renegotiate",2,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
		options,callback:=tsArg(args,0),tsArg(args,1)
		if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind&&options.kind!=tsArrayKind{tsInternalValidatorInvalidType("options","an object",options)}
		if !tsIsUndefined(callback)&&callback.kind!=tsFunctionKind{tsInternalValidatorInvalidType("callback","a function",callback)}
		if tsTruthy(tsNodeTLSGet(loop,receiver,"destroyed")){return tsU}
		if tsTruthy(tsNodeTLSGet(loop,receiver,"_renegotiationDisabled")){
			if callback.kind==tsFunctionKind{failure:=tsFSException("Error","ERR_TLS_RENEGOTIATION_DISABLED","TLS session renegotiation disabled for this socket");process:=tsNodeProcessModule(loop);tsNodeProcessNextTick(loop,(*tsObject)(process.ref),[]tsValue{callback,failure})}
			return tsBooleanValue(false)
		}
		state:=tsNodeTLSSocketFor(receiver);if state==nil||state.conn==nil{return tsU}
		if callback.kind==tsFunctionKind{failure:=tsFSException("Error","ERR_TLS_RENEGOTIATION_UNSUPPORTED","TLS session renegotiation is unsupported by this TLS implementation");process:=tsNodeProcessModule(loop);tsNodeProcessNextTick(loop,(*tsObject)(process.ref),[]tsValue{callback,failure})}
		return tsBooleanValue(false)
	}));prototype.set("disableRenegotiation",tsNamedNativeMethod("disableRenegotiation",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsNodeTLSSet(loop,receiver,"_renegotiationDisabled",tsBooleanValue(true));return tsU}));prototype.set("setMaxSendFragment",tsNamedNativeMethod("setMaxSendFragment",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{size:=tsArg(args,0);tsInternalValidatorInt32(size,tsNodeTLSKey("size"),tsU,tsU);return tsBooleanValue(false)}))
	prototype.set("enableTrace",tsNamedNativeMethod("enableTrace",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
	prototype.set("setServername",tsNamedNativeMethod("setServername",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{name:=tsArg(args,0);if name.kind!=tsStringKind{tsInternalValidatorInvalidType("name","a string",name)};if tsTruthy(tsNodeTLSGet(loop,receiver,"isServer")){panic(tsThrown{tsFSException("Error","ERR_TLS_SNI_FROM_SERVER","Cannot issue SNI from a TLS server-side socket")})};tsNodeTLSSet(loop,receiver,"servername",name);if state:=tsNodeTLSSocketFor(receiver);state!=nil&&tsNodeTLSObject(state.options){tsNodeTLSSet(loop,state.options,"servername",name)};return tsU}))
	prototype.set("setKeyCert",tsNamedNativeMethod("setKeyCert",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{context:=tsArg(args,0);if tsNodeTLSSecureFor(context)==nil{inner:=tsNodeTLSGet(loop,context,"context");if tsNodeTLSSecureFor(inner)==nil{tsNodeTLSSecureCreate(loop,context)}};tsRuntimeUnsupported("TLSSocket.setKeyCert","Go crypto/tls does not expose per-connection certificate replacement");return tsU}))
	prototype.set("_tlsError",tsNamedNativeMethod("_tlsError",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{failure:=tsArg(args,0);tsNodeNetEmit(loop,receiver,"_tlsError",failure);if tsTruthy(tsNodeTLSGet(loop,receiver,"_controlReleased")){return failure};return tsNull}))
	prototype.set("_emitTLSError",tsNamedNativeMethod("_emitTLSError",1,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{failure:=tsArg(args,0);handler:=tsNodeTLSGet(loop,receiver,"_tlsError");value:=tsCallReceiver(loop,handler,receiver,failure);if tsTruthy(value){tsNodeNetEmit(loop,receiver,"error",value)};return tsU}))
	prototype.set("_final",tsNamedNativeMethod("_final",1,tsNodeTLSSocketFinal))
	prototype.set("_releaseControl",tsNamedNativeMethod("_releaseControl",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{if tsTruthy(tsNodeTLSGet(loop,receiver,"_controlReleased")){return tsBooleanValue(false)};tsNodeTLSSet(loop,receiver,"_controlReleased",tsBooleanValue(true));return tsBooleanValue(true)}))
	prototype.set("_destroySSL",tsNamedNativeMethod("_destroySSL",0,func(loop *tsLoop,receiver tsValue,_ ...tsValue)tsValue{tsNodeTLSSet(loop,receiver,"secureConnecting",tsBooleanValue(false));tsNodeTLSSet(loop,receiver,"_secureEstablished",tsBooleanValue(false));return tsU}))
	prototype.set("_start",tsNamedNativeMethod("_start",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsU}))
}
func tsNodeTLSSocketFinal(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{callback:=tsArg(args,0);if tsTruthy(tsNodeTLSGet(loop,receiver,"secureConnecting")){symbols:=tsInternalNetSymbolsModule();preHandshakeWrite:=tsInternalNetGet(loop,symbols,"kPreHandshakeWrite");if preHandshakeWrite.kind==tsSymbolKind&&tsTruthy(tsGet(loop,receiver,tsPropertyKey(preHandshakeWrite))){secureConnectDone:=tsInternalNetGet(loop,symbols,"kSecureConnectDone");if secureConnectDone.kind==tsSymbolKind{resume:=tsFunctionValue(tsFunc(func(loop *tsLoop,_ ...tsValue)tsValue{return tsNodeNetSocketFinal(loop,receiver,[]tsValue{callback}...)}));once:=tsNodeTLSGet(loop,receiver,"once");if once.kind==tsFunctionKind{tsCallReceiver(loop,once,receiver,secureConnectDone,resume);return tsU}}}};return tsNodeNetSocketFinal(loop,receiver,[]tsValue{callback}...)}
func tsNodeTLSSocketCreate(loop *tsLoop,socket,options tsValue)tsValue{
	tsNodeTLSSetupSocketPrototype(loop);if tsIsUndefined(socket)||tsIsNull(socket){socket=tsNodeNetSocketCreate(loop,options)};if socket.kind==tsObjectKind{(*tsObject)(socket.ref).prototype=tsObjectValue(tsNodeTLS.socketPrototype)}
	if tsNodeTLSObject(options){options=tsObjectValue(tsNodeTLSObjectCopy(options))}else{options=tsObjectValue(tsNewObject())};server:=tsTruthy(tsNodeTLSGet(loop,options,"isServer"));tsNodeTLSSet(loop,socket,"isServer",tsBooleanValue(server));tsNodeTLSSet(loop,socket,"_controlReleased",tsBooleanValue(false));if !server{tsNodeTLSSet(loop,options,"_standaloneWrap",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_standaloneWrap",tsBooleanValue(true))};tsNodeTLSUpgradeSocket(loop,socket,options,server);return socket
}
func tsNodeTLSConnectOptions(loop *tsLoop,args []tsValue)(tsValue,tsValue){options:=tsObjectValue(tsNewObject());callback:=tsU;if len(args)>0&&tsNodeTLSObject(args[0]){options=tsObjectValue(tsNodeTLSObjectCopy(args[0]));if len(args)>1&&args[1].kind==tsFunctionKind{callback=args[1]};return options,callback};if len(args)>0&&tsIsNumeric(args[0]){options=tsObjectValue(tsNewObject());tsNodeTLSSet(loop,options,"port",args[0]);if len(args)>1&&args[1].kind==tsStringKind{tsNodeTLSSet(loop,options,"host",args[1])};for i:=1;i<len(args);i++{if tsNodeTLSObject(args[i]){for _,name:=range (*tsObject)(args[i].ref).order{tsNodeTLSSet(loop,options,name,(*tsObject)(args[i].ref).values[name])}};if args[i].kind==tsFunctionKind{callback=args[i]}};return options,callback};if len(args)>0&&args[0].kind==tsStringKind{options=tsObjectValue(tsNewObject());tsNodeTLSSet(loop,options,"path",args[0]);for i:=1;i<len(args);i++{if tsNodeTLSObject(args[i]){for _,name:=range (*tsObject)(args[i].ref).order{tsNodeTLSSet(loop,options,name,(*tsObject)(args[i].ref).values[name])}};if args[i].kind==tsFunctionKind{callback=args[i]}}};return options,callback}
func tsNodeTLSConnect(loop *tsLoop,args []tsValue)tsValue{
	options,callback:=tsNodeTLSConnectOptions(loop,args);hostValue:=tsNodeTLSGet(loop,options,"host");host:="localhost";if hostValue.kind==tsStringKind{host=tsText(hostValue)}
	servername:=tsNodeTLSGet(loop,options,"servername");if tsIsUndefined(servername)&&net.ParseIP(host)==nil{servername=tsNodeTLSKey(host);tsNodeTLSSet(loop,options,"servername",servername)};if servername.kind==tsStringKind&&net.ParseIP(tsText(servername))!=nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","Setting the TLS ServerName to an IP address is not permitted.")})}
	if _,err:=tsNodeTLSConfig(loop,options,false);err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};providedSocket:=tsNodeTLSGet(loop,options,"socket");socket:=tsNodeTLSSetupClientSocket(loop,options,providedSocket);if callback.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,socket,"once"),socket,tsNodeTLSKey("secureConnect"),callback)}
	if !tsNullish(providedSocket){if state:=tsNodeNetSocketFor(socket);state!=nil&&(state.conn!=nil||state.connecting){tsNodeTLSUpgradeSocket(loop,socket,options,false);return socket};panic(tsThrown{tsFSException("Error","ERR_SOCKET_CLOSED","The supplied socket is not connected")})}
	state:=tsNodeNetSocketFor(socket);if state==nil{return socket};tsNodeNetSocketConnect(loop,socket,[]tsValue{options}...);tsNodeTLSUpgradeSocket(loop,socket,options,false);return socket
}
func tsNodeTLSSetupClientSocket(loop *tsLoop,options,socket tsValue)tsValue{tsNodeNetModule(loop);if tsIsNull(socket)||tsIsUndefined(socket){socket=tsNodeNetSocketCreate(loop,options)};tsNodeTLSSetupSocketPrototype(loop);if socket.kind==tsObjectKind{(*tsObject)(socket.ref).prototype=tsObjectValue(tsNodeTLS.socketPrototype)};tsNodeTLSSet(loop,socket,"encrypted",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"authorized",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"secureConnecting",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_securePending",tsBooleanValue(true));tsNodeTLSSet(loop,socket,"_secureEstablished",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"_controlReleased",tsBooleanValue(false));tsNodeTLSSet(loop,socket,"isServer",tsBooleanValue(false));return socket}
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
		if !raw.closed{tsNodeTLSEmitSecureConnectDone(loop,socket)}
		if !raw.closed&&!tsTruthy(tsNodeNetGet(loop,server,"pauseOnConnect")){tsNodeNetStartRead(raw)}
	})
	task.unref=raw.unref
}
func tsNodeTLSServerListen(loop *tsLoop,receiver tsValue,args []tsValue)tsValue{return tsNodeNetListen(loop,receiver,args...)}
func tsNodeTLSCreateServer(loop *tsLoop,args []tsValue)tsValue{options,listener:=tsArg(args,0),tsArg(args,1);if options.kind==tsFunctionKind{listener=options;options=tsObjectValue(tsNewObject())};config,err:=tsNodeTLSConfig(loop,options,true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};tsNodeNetModule(loop);server:=tsNodeNetServerCreate(loop,options);object:=(*tsObject)(server.ref);if object.prototype!=nil{if tsNodeTLS.serverPrototype==nil{prototype:=tsNewObject();prototype.prototype=object.prototype;tsNodeTLS.serverPrototype=prototype};object.prototype=tsObjectValue(tsNodeTLS.serverPrototype)};secure:=&tsNodeTLSServerState{object:object,loop:loop,server:server,config:config,options:options,contexts:map[string]*cryptotls.Config{}};secure.config=tsNodeTLSServerConfig(secure,config);tsNodeTLS.servers[object]=secure;serverListener:=tsFunctionValue(tsFunc(func(loop *tsLoop,values ...tsValue)tsValue{socket:=tsArg(values,0);if raw:=tsNodeNetSocketFor(socket);raw!=nil{if conn,ok:=raw.conn.(*cryptotls.Conn);ok{tsNodeTLSServerHandshake(loop,server,socket,conn,options);return tsU}};tsNodeNetEmit(loop,server,"secureConnection",socket);return tsU}));tsCallReceiver(loop,tsNodeNetGet(loop,server,"on"),server,tsNodeTLSKey("connection"),serverListener);if listener.kind==tsFunctionKind{tsCallReceiver(loop,tsNodeNetGet(loop,server,"on"),server,tsNodeTLSKey("secureConnection"),listener)}
	if tsNodeTLS.serverPrototype!=nil{tsNodeTLS.serverProperties.set("prototype",tsObjectValue(tsNodeTLS.serverPrototype));tsNodeTLS.serverPrototype.set("listen",tsNamedNativeMethod("listen",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{return tsNodeTLSServerListen(loop,receiver,values)}));tsNodeTLS.serverPrototype.set("setSecureContext",tsNamedNativeMethod("setSecureContext",1,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{current:=tsNodeTLSServerFor(receiver);if current!=nil{config,err:=tsNodeTLSConfig(loop,tsArg(values,0),true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};current.mu.Lock();current.options=tsArg(values,0);current.config=tsNodeTLSServerConfig(current,config);current.mu.Unlock()};return tsU}));tsNodeTLS.serverPrototype.set("addContext",tsNamedNativeMethod("addContext",2,func(loop *tsLoop,receiver tsValue,values ...tsValue)tsValue{current:=tsNodeTLSServerFor(receiver);if current!=nil{config,err:=tsNodeTLSConfig(loop,tsArg(values,1),true);if err!=nil{panic(tsThrown{tsNodeTLSFailure(err)})};current.mu.Lock();current.contexts[tsNodeTLSServerNameKey(tsText(tsArg(values,0)))]=config;current.mu.Unlock()};return tsU}));tsNodeTLS.serverPrototype.set("getTicketKeys",tsNamedNativeMethod("getTicketKeys",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"Error",message:"Not implented in Bun yet"})})}));tsNodeTLS.serverPrototype.set("setTicketKeys",tsNamedNativeMethod("setTicketKeys",1,func(_ *tsLoop,_ tsValue,values ...tsValue)tsValue{keys:=tsArg(values,0);if keys.kind!=tsTypedArrayKind&&keys.kind!=tsDataViewKind{tsInternalValidatorInvalidType("buffer","an instance of Buffer, TypedArray, or DataView",keys)};if len(tsNodeTLSBytes(keys))!=48{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument 'buffer' is invalid. Session ticket keys must be a 48-byte buffer")})};panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"Error",message:"Not implented in Bun yet"})})}))}
	return server}
func tsNodeTLSIdentityError(loop *tsLoop,host,cert tsValue,reason string)tsValue{return tsErrorValue(&tsRuntimeError{name:"Error",message:"Hostname/IP does not match certificate's altnames: "+tsText(host),fields:map[string]tsValue{"code":tsNodeTLSKey("ERR_TLS_CERT_ALTNAME_INVALID"),"reason":tsNodeTLSKey(reason),"host":host,"cert":cert}})}
func tsNodeTLSCheckServerIdentity(loop *tsLoop,hostValue,cert tsValue)tsValue{
	host:=tsText(hostValue);normalized:=strings.TrimSuffix(host,".")
	if net.ParseIP(normalized)==nil&&strings.IndexFunc(normalized,func(r rune)bool{return r>127})>=0{ascii:=tsInternalURLIDNAToASCII(loop,[]tsValue{tsNodeTLSKey(normalized)});if tsText(ascii)!=""{normalized=tsText(ascii)}}
	if tsNodeTLSObject(cert){
		raw:=tsNodeTLSGet(loop,cert,"raw")
		if !tsNullish(raw){if parsed,err:=x509.ParseCertificate(tsNodeTLSBytes(raw));err==nil{
			verifyErr:=parsed.VerifyHostname(normalized);if verifyErr==nil{return tsU};if net.ParseIP(normalized)==nil&&len(parsed.DNSNames)==0&&parsed.Subject.CommonName!=""{commonName:=*parsed;commonName.DNSNames=[]string{parsed.Subject.CommonName};if commonName.VerifyHostname(normalized)==nil{return tsU}}
			return tsNodeTLSIdentityError(loop,hostValue,cert,verifyErr.Error())
		}}
		sans:=tsText(tsNodeTLSGet(loop,cert,"subjectaltname"));for _,entry:=range strings.Split(sans,", "){entry=strings.TrimPrefix(entry,"DNS:");entry=strings.TrimPrefix(entry,"IP Address:");if entry==normalized{return tsU};if strings.HasPrefix(entry,"*."){suffix:=entry[1:];if strings.HasSuffix(normalized,suffix)&&strings.Count(normalized,".")==strings.Count(entry,"."){return tsU}}};subject:=tsNodeTLSGet(loop,cert,"subject");if tsText(tsNodeTLSGet(loop,subject,"CN"))==normalized{return tsU}
	}
	return tsNodeTLSIdentityError(loop,hostValue,cert,"Hostname/IP does not match certificate's altnames")
}
func tsNodeTLSRootCerts()[]string{for _,path:=range []string{"/etc/ssl/certs/ca-certificates.crt","/etc/pki/tls/certs/ca-bundle.crt","/etc/ssl/ca-bundle.pem"}{if data,err:=os.ReadFile(path);err==nil{out:=[]string{};remaining:=data;for len(remaining)>0{block,rest:=pem.Decode(remaining);if block==nil{break};if block.Type=="CERTIFICATE"{out=append(out,string(pem.EncodeToMemory(block)))};remaining=rest};if len(out)>0{return out}}};return nil}
func tsNodeTLSExtraRootCerts()[]string{path:=os.Getenv("NODE_EXTRA_CA_CERTS");if path==""{return nil};data,err:=os.ReadFile(path);if err!=nil{return nil};out:=[]string{};for len(data)>0{block,rest:=pem.Decode(data);if block==nil{break};data=rest;if block.Type!="CERTIFICATE"{continue};certificate,parseErr:=x509.ParseCertificate(block.Bytes);if parseErr==nil{out=append(out,string(pem.EncodeToMemory(&pem.Block{Type:"CERTIFICATE",Bytes:certificate.Raw})))}};return out}
func tsNodeTLSDefaultCiphers()string{names:=[]string{"TLS_AES_256_GCM_SHA384","TLS_CHACHA20_POLY1305_SHA256","TLS_AES_128_GCM_SHA256"};legacy:="";if tsNodeTLS.defaultCiphers.kind==tsStringKind{legacy=tsText(tsNodeTLS.defaultCiphers)}else{values:=[]string{};for _,suite:=range cryptotls.CipherSuites(){if name:=tsNodeTLSOpenSSLCipherName(suite.ID);name!=""{values=append(values,name)}};legacy=strings.Join(values,":")};if legacy!=""{names=append(names,legacy)};return strings.Join(names,":")}
func tsNodeTLSModule(loop *tsLoop)tsValue{tsNodeTLS.Do(func(){tsNodeNetModule(loop);tsNodeTLS.secure=map[*tsObject]*tsNodeTLSSecureState{};tsNodeTLS.sockets=map[*tsObject]*tsNodeTLSSocketState{};tsNodeTLS.servers=map[*tsObject]*tsNodeTLSServerState{};tsNodeTLS.minVersion=tsNodeTLSKey("TLSv1.2");tsNodeTLS.maxVersion=tsNodeTLSKey("TLSv1.3");tsNodeTLS.defaultCurve=tsNodeTLSKey("auto");tsNodeTLS.defaultCAs=nil;tsNodeTLSSetupSocketPrototype(loop)
		secure:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))});secure.name="SecureContext";secure.constructible=true;secure.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))};secure.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))};secure.properties=tsNewObject();secure.properties.set("parsePkcs12",tsNamedNativeMethod("parsePkcs12",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSPKCS12Parse(loop,args)}));tsNodeTLS.secureConstructor=tsFunctionValue(secure)
		socket:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSocketCreate(loop,tsArg(args,0),tsArg(args,1))});socket.name="TLSSocket";socket.constructible=true;socket.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSSocketCreate(loop,tsArg(args,0),tsArg(args,1))};socket.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSocketCreate(loop,tsArg(args,0),tsArg(args,1))};socket.properties=tsNewObject();socket.properties.set("prototype",tsObjectValue(tsNodeTLS.socketPrototype));tsNodeTLS.socketConstructor=tsFunctionValue(socket)
		server:=tsFunc(func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)});server.name="Server";server.constructible=true;server.construct=func(loop *tsLoop,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)};server.receiverCall=func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)};server.properties=tsNewObject();tsNodeTLS.serverProperties=server.properties;tsNodeTLS.serverConstructor=tsFunctionValue(server)
		module:=tsNewObject();module.set("CLIENT_RENEG_LIMIT",tsNumberValue(3));module.set("CLIENT_RENEG_WINDOW",tsNumberValue(600));module.set("connect",tsNamedNativeMethod("connect",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSConnect(loop,args)}));module.set("createConnection",module.get("connect"));module.set("createSecureContext",tsNamedNativeMethod("createSecureContext",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSSecureCreate(loop,tsArg(args,0))}));module.set("SecureContext",tsNodeTLS.secureConstructor);module.set("TLSSocket",tsNodeTLS.socketConstructor);module.set("Server",tsNodeTLS.serverConstructor);module.set("createServer",tsNamedNativeMethod("createServer",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSCreateServer(loop,args)}));module.set("DEFAULT_CIPHERS",tsNodeTLSKey(tsNodeTLSDefaultCiphers()));module.set("DEFAULT_ECDH_CURVE",tsNodeTLS.defaultCurve);module.set("DEFAULT_MIN_VERSION",tsNodeTLS.minVersion);module.set("DEFAULT_MAX_VERSION",tsNodeTLS.maxVersion);if module.descriptors==nil{module.descriptors=map[string]*tsDescriptor{}};module.descriptors["DEFAULT_CIPHERS"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_CIPHERS",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLSKey(tsNodeTLSDefaultCiphers())}),setter:tsNamedNativeMethod("set DEFAULT_CIPHERS",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if tsTruthy(value){tsInternalTLSValidateCiphers(value,"value");filtered:=tsInternalTLSStripCipherNames(value);value=tsInternalTLSGet(loop,filtered,"cipherList")}else{value=tsNodeTLSKey("")};tsNodeTLS.defaultCiphers=value;return tsU})};module.descriptors["DEFAULT_ECDH_CURVE"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_ECDH_CURVE",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLS.defaultCurve}),setter:tsNamedNativeMethod("set DEFAULT_ECDH_CURVE",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeTLS.defaultCurve=tsArg(args,0);return tsU})};module.descriptors["DEFAULT_MIN_VERSION"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_MIN_VERSION",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLS.minVersion}),setter:tsNamedNativeMethod("set DEFAULT_MIN_VERSION",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeTLS.minVersion=tsArg(args,0);return tsU})};module.descriptors["DEFAULT_MAX_VERSION"]=&tsDescriptor{accessor:true,enumerable:true,configurable:true,getter:tsNamedNativeMethod("get DEFAULT_MAX_VERSION",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLS.maxVersion}),setter:tsNamedNativeMethod("set DEFAULT_MAX_VERSION",1,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{tsNodeTLS.maxVersion=tsArg(args,0);return tsU})}
		module.set("convertALPNProtocols",tsNamedNativeMethod("convertALPNProtocols",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSConvertALPNProtocols(loop,tsArg(args,0),tsArg(args,1))}));module.set("getCiphers",tsNamedNativeMethod("getCiphers",0,func(loop *tsLoop,_ tsValue,_ ...tsValue)tsValue{names:=[]string{"tls_aes_128_gcm_sha256","tls_aes_256_gcm_sha384","tls_chacha20_poly1305_sha256","des-cbc3-sha","aes128-sha","aes256-sha","psk-aes128-cbc-sha","psk-aes256-cbc-sha","aes128-gcm-sha256","aes256-gcm-sha384","ecdhe-ecdsa-aes128-sha","ecdhe-ecdsa-aes256-sha","ecdhe-rsa-aes128-sha","ecdhe-rsa-aes256-sha","ecdhe-ecdsa-aes128-sha256","ecdhe-rsa-aes128-sha256","ecdhe-ecdsa-aes128-gcm-sha256","ecdhe-ecdsa-aes256-gcm-sha384","ecdhe-rsa-aes128-gcm-sha256","ecdhe-rsa-aes256-gcm-sha384","ecdhe-psk-aes128-cbc-sha","ecdhe-psk-aes256-cbc-sha","ecdhe-rsa-chacha20-poly1305","ecdhe-ecdsa-chacha20-poly1305","ecdhe-psk-chacha20-poly1305"};sort.Strings(names);return tsNodeTLSStringArray(loop,names)}));module.set("getDefaultCiphers",tsNamedNativeMethod("getDefaultCiphers",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return tsNodeTLSKey(tsNodeTLSDefaultCiphers())}));module.set("parseCertString",tsNamedNativeMethod("parseCertString",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{panic(tsThrown{tsInternalSharedNotImplemented(tsNodeTLSKey("Not implemented"),tsU,tsU)})}));module.set("checkServerIdentity",tsNamedNativeMethod("checkServerIdentity",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeTLSCheckServerIdentity(loop,tsArg(args,0),tsArg(args,1))}))
		module.set("getCACertificates",tsNamedNativeMethod("getCACertificates",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{kind:=tsArg(args,0);if tsIsUndefined(kind){kind=tsNodeTLSKey("default")};if kind.kind!=tsStringKind{tsFSInvalid("type","a string")};var certs []string;switch tsText(kind){case "default":if tsNodeTLS.defaultCAsSet{out:=&tsArray{};for _,cert:=range tsNodeTLS.defaultCAs{out.values=append(out.values,cert)};return tsArrayValue(out)};certs=tsNodeTLSRootCerts();if os.Getenv("NODE_EXTRA_CA_CERTS")!=""{certs=append(certs,tsNodeTLSExtraRootCerts()...)};case "bundled","system":certs=tsNodeTLSRootCerts();case "extra":certs=tsNodeTLSExtraRootCerts();default:panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The argument \"type\" is invalid")})};out:=&tsArray{};for _,cert:=range certs{out.values=append(out.values,tsNodeTLSKey(cert))};return tsArrayValue(out)}));module.set("setDefaultCACertificates",tsNamedNativeMethod("setDefaultCACertificates",1,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{value:=tsArg(args,0);if value.kind!=tsArrayKind{tsFSInvalid("certs","an instance of Array")};array:=(*tsArray)(value.ref);for index,item:=range array.values{if item.kind!=tsStringKind&&item.kind!=tsTypedArrayKind&&item.kind!=tsDataViewKind{tsFSInvalid("certs["+strconv.Itoa(index)+"]","string or an instance of ArrayBufferView")}};normalized,err:=tsNodeTLSParseCACertificates(loop,value);if err!=nil{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED",err.Error())})};if len(normalized)==0&&array.length()>0{panic(tsThrown{tsFSException("Error","ERR_CRYPTO_OPERATION_FAILED","No valid certificates found in the provided array")})};tsNodeTLS.defaultCAs=normalized;tsNodeTLS.defaultCAsSet=true;return tsU}));roots:=&tsArray{};for _,cert:=range tsNodeTLSRootCerts(){roots.values=append(roots.values,tsNodeTLSKey(cert))};rootValues:=tsArrayValue(roots);module.set("rootCertificates",rootValues);if module.descriptors==nil{module.descriptors=map[string]*tsDescriptor{}};module.descriptors["rootCertificates"]=&tsDescriptor{accessor:true,enumerable:true,configurable:false,getter:tsNamedNativeMethod("get rootCertificates",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{return rootValues})};module.set("default",tsObjectValue(module));tsNodeTLS.module=tsObjectValue(module)
	});return tsNodeTLS.module}
`

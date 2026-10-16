package goemit

// nodeSQLiteRuntime maps Bun's native node:sqlite binding onto SQLite through
// the bundled Cgo driver. SQLite's engine internals are supplied by
// github.com/mattn/go-sqlite3; this fragment owns the Node API and value bridge.
const nodeSQLiteRuntime = `
type tsNodeSQLiteDB struct{object,limitsObject *tsObject;location string;opened,readonly,transaction,defensive,allowExtension,extensions,foreignKeys,doubleQuoted bool;timeout int;readBigInts,returnArrays,allowBare,allowUnknown bool;totalChanges,lastInsertRowID int64;callbackFailure any;functions,functionOptions,aggregates map[string]tsValue;authorizer tsValue;limits,limitMax map[string]int;statements map[*tsObject]bool;pool *sql.DB;conn *sql.Conn;native *sqlite3.SQLiteConn}
type tsNodeSQLiteStatement struct{object *tsObject;db tsValue;sql string;finalized,safeIntegers,returnArrays,allowBare,allowUnknown,stepping,readonly,executed,tagStore bool;columns,declared,nullable,databases,tables,origins []string;parameterNames []string;lastBindings []any;params int;stmt *sql.Stmt;iteratorClose func()}
type tsNodeSQLiteSession struct{object *tsObject;db tsValue;closed bool;native *sqlite3.SQLiteSession}
type tsNodeSQLiteTagStore struct{object *tsObject;db tsValue;capacity uint64;cache map[string]tsValue;order []string}
var tsNodeSQLiteSingleton struct{sync.Once;module tsValue;dbProto,stmtProto,sessionProto,tagProto *tsObject;dbClass,stmtClass,sessionClass tsValue;dbs map[*tsObject]*tsNodeSQLiteDB;statements map[*tsObject]*tsNodeSQLiteStatement;sessions map[*tsObject]*tsNodeSQLiteSession;tags map[*tsObject]*tsNodeSQLiteTagStore;constants tsValue}
func tsNodeSQLiteKey(name string)tsValue{return tsStringReference(tsStringUTF8(name))}
func tsNodeSQLiteGet(loop *tsLoop,value tsValue,name string)tsValue{return tsGet(loop,value,tsNodeSQLiteKey(name))}
func tsNodeSQLiteSet(loop *tsLoop,value tsValue,name string,item tsValue){tsSet(loop,value,tsNodeSQLiteKey(name),item)}
func tsNodeSQLiteError(code,message string)tsValue{name:="Error";switch code{case "ERR_INVALID_ARG_TYPE","ERR_INVALID_ARG_VALUE":name="TypeError";case "ERR_OUT_OF_RANGE":name="RangeError"};return tsFSException(name,code,message)}
func tsNodeSQLiteDBFor(value tsValue)*tsNodeSQLiteDB{if value.kind!=tsObjectKind{return nil};return tsNodeSQLiteSingleton.dbs[(*tsObject)(value.ref)]}
func tsNodeSQLiteStatementFor(value tsValue)*tsNodeSQLiteStatement{if value.kind!=tsObjectKind{return nil};return tsNodeSQLiteSingleton.statements[(*tsObject)(value.ref)]}
func tsNodeSQLiteSessionFor(value tsValue)*tsNodeSQLiteSession{if value.kind!=tsObjectKind{return nil};return tsNodeSQLiteSingleton.sessions[(*tsObject)(value.ref)]}
func tsNodeSQLiteTagFor(value tsValue)*tsNodeSQLiteTagStore{if value.kind!=tsObjectKind{return nil};return tsNodeSQLiteSingleton.tags[(*tsObject)(value.ref)]}
func tsNodeSQLiteRequireDB(value tsValue)*tsNodeSQLiteDB{db:=tsNodeSQLiteDBFor(value);if db==nil{tsPropertyFailure("DatabaseSync method called on an incompatible receiver")};if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};return db}
func tsNodeSQLiteRethrowCallbackFailure(db *tsNodeSQLiteDB){if db!=nil&&db.callbackFailure!=nil{failure:=db.callbackFailure;db.callbackFailure=nil;panic(failure)}}
func tsNodeSQLiteAccessor(proto *tsObject,name string,get func(*tsLoop,tsValue)tsValue,set func(*tsLoop,tsValue,tsValue)tsValue){if proto.descriptors==nil{proto.descriptors=map[string]*tsDescriptor{}};descriptor:=&tsDescriptor{accessor:true,configurable:true,getter:tsNamedNativeMethod("get "+name,0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return get(loop,r)})};if set!=nil{descriptor.setter=tsNamedNativeMethod("set "+name,1,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return set(loop,r,tsArg(args,0))})};proto.descriptors[name]=descriptor}
func tsNodeSQLiteDefaultLimits()map[string]int{return map[string]int{"length":1000000000,"sqlLength":1000000000,"column":2000,"exprDepth":1000,"compoundSelect":500,"vdbeOp":250000,"functionArg":127,"attach":10,"likePatternLength":50000,"variableNumber":32766,"triggerDepth":1000}}
func tsNodeSQLiteLimits(loop *tsLoop,receiver tsValue)tsValue{
	db:=tsNodeSQLiteDBFor(receiver);if db==nil{tsPropertyFailure("DatabaseSync method called on an incompatible receiver")}
	if db.limitsObject!=nil{return tsObjectValue(db.limitsObject)}
	object:=tsNewObject();object.prototype=tsObjectValue(tsObjectPrototype());names:=[]string{"length","sqlLength","column","exprDepth","compoundSelect","vdbeOp","functionArg","attach","likePatternLength","variableNumber","triggerDepth"}
	for _,name:=range names{
		property:=name;object.set(property,tsNumberValue(float64(db.limits[property])))
		tsNodeSQLiteAccessor(object,property,func(loop *tsLoop,_ tsValue)tsValue{current:=tsNodeSQLiteRequireDB(receiver);return tsNumberValue(float64(current.limits[property]))},func(_ *tsLoop,_ tsValue,value tsValue)tsValue{
			current:=tsNodeSQLiteRequireDB(receiver);if !tsIsNumeric(value){panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_TYPE","Limit value must be a non-negative integer or Infinity.")})}
			number:=tsNumber(value);limit:=0
			if math.IsInf(number,1){limit=current.limitMax[property]}else if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<math.MinInt32||number>math.MaxInt32{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_TYPE","Limit value must be a non-negative integer or Infinity.")})}else{limit=int(number);if limit<0{panic(tsThrown{tsNodeSQLiteError("ERR_OUT_OF_RANGE","Limit value must be non-negative.")})}}
			if limit>current.limitMax[property]{limit=current.limitMax[property]};if current.native!=nil{current.native.SetLimit(tsNodeSQLiteLimitID(property),limit)};current.limits[property]=limit;return value
		})
		object.descriptors[property].enumerable=true;object.descriptors[property].configurable=false
	}
	db.limitsObject=object;return tsObjectValue(object)
}
func tsNodeSQLiteConstructor(name string,proto *tsObject,length int,construct func(*tsLoop,[]tsValue)tsValue)tsValue{fn:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});fn.name=name;fn.length=length;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return construct(loop,args)};fn.properties=tsNewObject();fn.properties.set("prototype",tsObjectValue(proto));fn.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};proto.set("constructor",tsFunctionValue(fn));return tsFunctionValue(fn)}
func tsNodeSQLiteRowObject(loop *tsLoop,columns []string,row []tsValue,arrays bool)tsValue{if arrays{return tsArrayValue(&tsArray{values:append([]tsValue(nil),row...)})};o:=tsNewObject();o.prototype=tsNull;for i,name:=range columns{if i<len(row){o.set(name,row[i])}};return tsObjectValue(o)}
func tsNodeSQLiteExecute(loop *tsLoop,statement tsValue,args []tsValue,mode string)tsValue{
	state:=tsNodeSQLiteStatementFor(statement)
	if state==nil{tsPropertyFailure("Invalid StatementSync receiver")}
	db:=tsNodeSQLiteRequireDB(state.db)
	if state.finalized{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})}
	if state.stepping{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement is currently executing")})}
	if state.stmt==nil{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})}
	state.stepping=true
	releaseStepping:=true
	defer func(){if releaseStepping{state.stepping=false}}()
	bindings,err:=tsNodeSQLiteBindings(loop,args,state)
	if err!=nil{code:="ERR_INVALID_ARG_TYPE";if strings.HasPrefix(err.Error(),"Unknown named parameter")||strings.HasPrefix(err.Error(),"Cannot create bare named parameter"){code="ERR_INVALID_STATE"}else if strings.HasPrefix(err.Error(),"BigInt value is too large to bind"){code="ERR_INVALID_ARG_VALUE"}else if strings.HasPrefix(err.Error(),"Too many SQL parameters"){code="ERR_SQLITE_ERROR"};panic(tsThrown{tsNodeSQLiteError(code,err.Error())})}
	if !db.opened||state.finalized||state.stmt==nil{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})}
	state.lastBindings=append(state.lastBindings[:0],bindings...)
	if mode=="run"{
		result,err:=state.stmt.ExecContext(context.Background(),bindings...)
		if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		tsNodeSQLiteRethrowCallbackFailure(db)
		state.executed=true
		changes,_:=result.RowsAffected()
		lastInsertRowID,_:=result.LastInsertId()
		db.lastInsertRowID=lastInsertRowID
		if db.native!=nil{db.transaction=!db.native.AutoCommit()}
		out:=tsNewObject()
		changeValue:=tsNumberValue(float64(changes))
		lastValue:=tsNumberValue(float64(lastInsertRowID))
		if state.safeIntegers{changeValue=tsBigIntSigned(changes);lastValue=tsBigIntSigned(lastInsertRowID)}
		out.set("changes",changeValue);out.set("lastInsertRowid",lastValue)
		return tsObjectValue(out)
	}
	rows,columns,err:=tsNodeSQLiteQueryRows(state,bindings)
	if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
	tsNodeSQLiteRethrowCallbackFailure(db)
	state.executed=true
	if mode=="iterate"{releaseStepping=false;return tsNodeSQLiteRowsIterator(loop,rows,state,columns)}
	defer rows.Close()
	out:=&tsArray{}
	for rows.Next(){
		value,scanErr:=tsNodeSQLiteScanRow(rows,columns,state)
		if scanErr!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",scanErr.Error())})}
		item:=tsNodeSQLiteRowObject(loop,state.columns,value,state.returnArrays)
		if mode=="get"{return item}
		out.values=append(out.values,item)
	}
	if err:=rows.Err();err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
	if mode=="get"{return tsU}
	return tsArrayValue(out)
}
func tsNodeSQLiteStatement(loop *tsLoop,dbValue tsValue,sqlValue,options tsValue)tsValue{
	db:=tsNodeSQLiteRequireDB(dbValue)
	if sqlValue.kind!=tsStringKind{tsFSInvalid("sql","a string")}
	sqlText:=tsText(sqlValue)
	state:=&tsNodeSQLiteStatement{db:dbValue,sql:sqlText,allowBare:db.allowBare,allowUnknown:db.allowUnknown,returnArrays:db.returnArrays,safeIntegers:db.readBigInts}
	if !tsIsUndefined(options){if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};for _,entry:=range []struct{name string;target *bool}{{"readBigInts",&state.safeIntegers},{"returnArrays",&state.returnArrays},{"allowBareNamedParameters",&state.allowBare},{"allowUnknownNamedParameters",&state.allowUnknown}}{value:=tsNodeSQLiteGet(loop,options,entry.name);if !tsIsUndefined(value){if value.kind!=tsBooleanKind{tsFSInvalid("options."+entry.name,"a boolean")};*entry.target=tsTruthy(value)}};db=tsNodeSQLiteRequireDB(dbValue)}
	var parameterCount int
	var parameterNames []string
	var columnMetadata []sqlite3.SQLiteColumnMetadata
	var readonly bool
	if db.authorizer.kind==tsFunctionKind{db.native.RegisterAuthorizerWithSource(nil)}
	err:=db.conn.Raw(func(driverConn any)error{
		native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)}
		prepared,err:=native.Prepare(sqlText);if err!=nil{return err}
		compiled,ok:=prepared.(*sqlite3.SQLiteStmt);if !ok{_ = prepared.Close();return fmt.Errorf("unexpected SQLite statement type %T",prepared)}
		parameterCount=compiled.NumInput()
		parameterNames=compiled.BindParameterNames()
		columnMetadata=compiled.ColumnMetadata()
		readonly=compiled.Readonly()
		return compiled.Close()
	})
	if db.authorizer.kind==tsFunctionKind{tsNodeSQLiteInstallAuthorizer(loop,db)}
	if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
	state.params=parameterCount;state.parameterNames=parameterNames;state.readonly=readonly
	for _,column:=range columnMetadata{state.columns=append(state.columns,column.Name);state.declared=append(state.declared,column.DeclaredType);state.databases=append(state.databases,column.Database);state.tables=append(state.tables,column.Table);state.origins=append(state.origins,column.Origin);state.nullable=append(state.nullable,true)}
	state.stmt,err=db.conn.PrepareContext(context.Background(),sqlText)
	if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
	obj:=tsNewObject();obj.prototype=tsObjectValue(tsNodeSQLiteSingleton.stmtProto);state.object=obj
	tsNodeSQLiteSingleton.statements[obj]=state;db.statements[obj]=true
	_ = loop
	return tsObjectValue(obj)
}
func tsNodeSQLitePath(loop *tsLoop,pathValue tsValue)string{
	location:=""
	switch pathValue.kind {
	case tsStringKind:
		location=tsText(pathValue)
	case tsTypedArrayKind:
		view:=(*tsTypedArray)(pathValue.ref)
		if view.name!="Uint8Array"&&!view.nodeBuffer{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"path\" argument must be a string, Uint8Array, or URL without null bytes.")})}
		data:=tsNodeBufferVisibleBytes(pathValue)
		if !utf8.Valid(data){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE","The \"path\" argument must be a Uint8Array containing a valid UTF-8 byte sequence.")})}
		location=string(data)
	case tsObjectKind,tsInstanceKind:
		href,protocol:=tsNodeSQLiteGet(loop,pathValue,"href"),tsNodeSQLiteGet(loop,pathValue,"protocol")
		if href.kind!=tsStringKind||protocol.kind!=tsStringKind{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"path\" argument must be a string, Uint8Array, or URL without null bytes.")})}
		if tsText(protocol)!="file:"||!strings.HasPrefix(tsText(href),"file:"){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_URL_SCHEME","The URL must be of scheme file:")})}
		location=tsText(href)
	default:
		panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"path\" argument must be a string, Uint8Array, or URL without null bytes.")})
	}
	if strings.ContainsRune(location,0){panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_TYPE","The \"path\" argument must be a string, Uint8Array, or URL without null bytes.")})}
	return location
}
func tsNodeSQLiteLimitID(name string)int{switch name{case "length":return sqlite3.SQLITE_LIMIT_LENGTH;case "sqlLength":return sqlite3.SQLITE_LIMIT_SQL_LENGTH;case "column":return sqlite3.SQLITE_LIMIT_COLUMN;case "exprDepth":return sqlite3.SQLITE_LIMIT_EXPR_DEPTH;case "compoundSelect":return sqlite3.SQLITE_LIMIT_COMPOUND_SELECT;case "vdbeOp":return sqlite3.SQLITE_LIMIT_VDBE_OP;case "functionArg":return sqlite3.SQLITE_LIMIT_FUNCTION_ARG;case "attach":return sqlite3.SQLITE_LIMIT_ATTACHED;case "likePatternLength":return sqlite3.SQLITE_LIMIT_LIKE_PATTERN_LENGTH;case "variableNumber":return sqlite3.SQLITE_LIMIT_VARIABLE_NUMBER;case "triggerDepth":return sqlite3.SQLITE_LIMIT_TRIGGER_DEPTH};return -1}
func tsNodeSQLiteDSN(location string,readonly,foreignKeys bool,timeout int)string{
	base:=location;query:=url.Values{}
	if parsed,err:=url.Parse(location);err==nil&&parsed.Scheme=="file"{base=parsed.Scheme+":"+parsed.Opaque;if parsed.Opaque==""{base=parsed.Scheme+":"+parsed.EscapedPath()};query=parsed.Query()}else if location!=":memory:"{base=(&url.URL{Scheme:"file",Path:filepath.ToSlash(location)}).String()}
	query.Set("_busy_timeout",strconv.Itoa(timeout));query.Set("_foreign_keys",strconv.FormatBool(foreignKeys));if readonly{query.Set("mode","ro");query.Set("_query_only","on")};if strings.Contains(base,"?"){parts:=strings.SplitN(base,"?",2);base=parts[0];if parsed,err:=url.ParseQuery(parts[1]);err==nil{for key,values:=range parsed{query[key]=values}}};return base+"?"+query.Encode()
}
func tsNodeSQLiteOpen(loop *tsLoop,state *tsNodeSQLiteDB)error{
	if state==nil{return fmt.Errorf("invalid database state")};if state.conn!=nil{state.opened=true;return nil}
	pool,err:=sql.Open("sqlite3",tsNodeSQLiteDSN(state.location,state.readonly,state.foreignKeys,state.timeout));if err!=nil{return err};pool.SetMaxOpenConns(1);pool.SetMaxIdleConns(1)
	conn,err:=pool.Conn(context.Background());if err!=nil{pool.Close();return err};var native *sqlite3.SQLiteConn
	err=conn.Raw(func(driverConn any)error{var ok bool;native,ok=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};return nil});if err!=nil{conn.Close();pool.Close();return err}
	state.pool,state.conn,state.native=pool,conn,native
	native.RegisterUpdateHook(func(_ int,_ string,_ string,_ int64){state.totalChanges++})
	for _,name:=range []string{"length","sqlLength","column","exprDepth","compoundSelect","vdbeOp","functionArg","attach","likePatternLength","variableNumber","triggerDepth"}{id:=tsNodeSQLiteLimitID(name);maximum:=native.GetLimit(id);state.limitMax[name]=maximum;if state.limits[name]>maximum{state.limits[name]=maximum};native.SetLimit(id,state.limits[name])}
	if err:=native.SetDefensive(state.defensive);err!=nil{conn.Close();pool.Close();state.conn=nil;state.pool=nil;state.native=nil;return err}
	if err:=native.SetDoubleQuotedStringLiterals(state.doubleQuoted);err!=nil{conn.Close();pool.Close();state.conn=nil;state.pool=nil;state.native=nil;return err}
	if err:=native.EnableLoadExtension(state.extensions);err!=nil{conn.Close();pool.Close();state.conn=nil;state.pool=nil;state.native=nil;return err}
	tsNodeSQLiteInstallAuthorizer(loop,state)
	for name,fn:=range state.functions{if err:=tsNodeSQLiteRegisterFunction(loop,state,name,fn,state.functionOptions[name]);err!=nil{conn.Close();pool.Close();state.conn=nil;state.pool=nil;state.native=nil;return err}}
	for name,options:=range state.aggregates{if err:=tsNodeSQLiteRegisterAggregate(loop,state,name,options);err!=nil{conn.Close();pool.Close();state.conn=nil;state.pool=nil;state.native=nil;return err}}
	if err:=conn.PingContext(context.Background());err!=nil{conn.Close();pool.Close();state.conn=nil;state.pool=nil;state.native=nil;return err}
	state.opened=true;return nil
}
func tsNodeSQLiteClose(state *tsNodeSQLiteDB){if state==nil{return};for _,session:=range tsNodeSQLiteSingleton.sessions{if session!=nil&&tsNodeSQLiteDBFor(session.db)==state&&!session.closed{if session.native!=nil{_ = session.native.Close();session.native=nil};session.closed=true}};for object:=range state.statements{if statement:=tsNodeSQLiteSingleton.statements[object];statement!=nil{if statement.iteratorClose!=nil{statement.iteratorClose()};statement.finalized=true;if statement.stmt!=nil{statement.stmt.Close();statement.stmt=nil}}};if state.conn!=nil{state.conn.Close();state.conn=nil};if state.pool!=nil{state.pool.Close();state.pool=nil};state.native=nil;state.opened=false;state.transaction=false}
func tsNodeSQLiteInstallAuthorizer(loop *tsLoop,state *tsNodeSQLiteDB){if state==nil||state.native==nil{return};if state.authorizer.kind==tsNullKind||tsIsUndefined(state.authorizer){state.native.RegisterAuthorizerWithSource(nil);return};callback:=state.authorizer;state.native.RegisterAuthorizerWithSource(func(action int,first,second,database,source string) (result int){defer func(){if failure:=recover();failure!=nil{state.callbackFailure=failure;result=sqlite3.SQLITE_DENY}}();toValue:=func(text string)tsValue{if text==""{return tsNull};return tsNodeSQLiteKey(text)};value:=tsCall(loop,callback,tsNumberValue(float64(action)),toValue(first),toValue(second),toValue(database),toValue(source));if !tsIsNumeric(value){state.callbackFailure=tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_TYPE","Authorizer callback must return an integer authorization code")};return sqlite3.SQLITE_DENY};number:=tsNumber(value);if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<math.MinInt32||number>math.MaxInt32{state.callbackFailure=tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_TYPE","Authorizer callback must return an integer authorization code")};return sqlite3.SQLITE_DENY};code:=int(number);if code!=sqlite3.SQLITE_OK&&code!=sqlite3.SQLITE_DENY&&code!=sqlite3.SQLITE_IGNORE{state.callbackFailure=tsThrown{tsNodeSQLiteError("ERR_OUT_OF_RANGE","Authorizer callback returned an invalid authorization code")};return sqlite3.SQLITE_DENY};return code})}
func tsNodeSQLiteDriverValue(loop *tsLoop,value tsValue)(any,error){
	switch value.kind{
	case tsUndefinedKind,tsNullKind:return nil,nil
	case tsStringKind:return tsText(value),nil
case tsBigIntKind:number:=(*big.Int)(value.ref);if !number.IsInt64(){return nil,fmt.Errorf("BigInt value is too large to bind")};return number.Int64(),nil
case tsTypedArrayKind,tsDataViewKind:return append([]byte(nil),tsNodeBufferVisibleBytes(value)...),nil
	default:
		if tsIsNumeric(value){return tsNumber(value),nil};return nil,fmt.Errorf("Returned JavaScript value cannot be converted to a SQLite value")
	}
}
func tsNodeSQLiteBindingValue(loop *tsLoop,value tsValue)(any,error){if value.kind==tsUndefinedKind||value.kind==tsBooleanKind||value.kind==tsArrayBufferKind||value.kind==tsSymbolKind{return nil,fmt.Errorf("Provided value cannot be bound to SQLite parameter")};return tsNodeSQLiteDriverValue(loop,value)}
func tsNodeSQLiteBindings(loop *tsLoop,args []tsValue,state *tsNodeSQLiteStatement)([]any,error){
	bindings:=make([]any,state.params);start:=0
	if !state.tagStore&&len(args)>0&&(args[0].kind==tsObjectKind||args[0].kind==tsInstanceKind||args[0].kind==tsArrayKind||args[0].kind==tsArrayBufferKind){
		named:=args[0];start=1;bare:=map[string]string{}
		if state.allowBare{for _,full:=range state.parameterNames{if full==""||full[0]=='?'{continue};key:=full[1:];if existing,ok:=bare[key];ok&&existing!=full{return nil,fmt.Errorf("Cannot create bare named parameter '%s' because of conflicting names '%s' and '%s'.",key,existing,full)};bare[key]=full}}
		for _,keyValue:=range tsEnumerableOwnKeys(loop,named).values{
			key:=tsPropertyKey(keyValue);index:=0
			for i,full:=range state.parameterNames{if full==key{index=i+1;break}}
			if index==0&&state.allowBare{if full,ok:=bare[key];ok{for i,name:=range state.parameterNames{if name==full{index=i+1;break}}}}
			if index==0{if state.allowUnknown{continue};return nil,fmt.Errorf("Unknown named parameter '%s'",key)}
			native,err:=tsNodeSQLiteBindingValue(loop,tsGet(loop,named,keyValue));if err!=nil{return nil,err};bindings[index-1]=native
		}
	}
	index:=0
	for i:=start;i<len(args);i++{
		for index<len(state.parameterNames){name:=state.parameterNames[index];if name==""||name[0]=='?'{break};index++}
		if index>=state.params{return nil,fmt.Errorf("Too many SQL parameters were provided")}
		native,err:=tsNodeSQLiteBindingValue(loop,args[i]);if err!=nil{return nil,err};bindings[index]=native;index++
	}
	return bindings,nil
}
func tsNodeSQLiteJSValue(value any,safeIntegers bool)tsValue{switch item:=value.(type){case nil:return tsNull;case int64:if safeIntegers{return tsBigIntSigned(item)};return tsNumberValue(float64(item));case int:return tsNumberValue(float64(item));case float64:return tsNumberValue(item);case bool:return tsBooleanValue(item);case string:return tsNodeSQLiteKey(item);case []byte:return tsNodeBuffer(append([]byte(nil),item...));case time.Time:return tsNodeSQLiteKey(item.UTC().Format(time.RFC3339Nano));default:return tsU}}
func tsNodeSQLiteCollectMetadata(rows *sql.Rows,state *tsNodeSQLiteStatement)[]string{columns,err:=rows.Columns();if err!=nil{return nil};state.columns=append(state.columns[:0],columns...);state.declared=make([]string,len(columns));state.nullable=make([]bool,len(columns));types,typeErr:=rows.ColumnTypes();if typeErr==nil{for i,columnType:=range types{if i<len(state.declared){state.declared[i]=columnType.DatabaseTypeName();nullable,ok:=columnType.Nullable();state.nullable[i]=!ok||nullable}}};return columns}
func tsNodeSQLiteScanRow(rows *sql.Rows,columns []string,state *tsNodeSQLiteStatement)([]tsValue,error){raw:=make([]any,len(columns));targets:=make([]any,len(columns));for i:=range raw{targets[i]=&raw[i]};if err:=rows.Scan(targets...);err!=nil{return nil,err};out:=make([]tsValue,len(raw));for i,value:=range raw{out[i]=tsNodeSQLiteJSValue(value,state.safeIntegers)};return out,nil}
func tsNodeSQLiteRegisterFunction(loop *tsLoop,state *tsNodeSQLiteDB,name string,fn,options tsValue)error{
	if state==nil||state.native==nil{return nil}
	useBigInt:=tsTruthy(tsNodeSQLiteGet(loop,options,"useBigIntArguments"))
	varargs:=tsTruthy(tsNodeSQLiteGet(loop,options,"varargs"))
	deterministic:=tsTruthy(tsNodeSQLiteGet(loop,options,"deterministic"))
	directOnly:=tsTruthy(tsNodeSQLiteGet(loop,options,"directOnly"))
	length:=tsNodeSQLiteGet(loop,fn,"length");arity:=0;if tsIsNumeric(length){arity=int(tsNumber(length))}
	return state.native.RegisterFuncWithOptions(name,func(args ...any)(result any,err error){
		defer func(){if failure:=recover();failure!=nil{state.callbackFailure=failure;err=fmt.Errorf("SQLite function callback failed: %v",failure)}}()
		if !varargs&&len(args)!=arity{return nil,fmt.Errorf("SQLite function %s expected %d arguments, received %d",name,arity,len(args))}
		values:=make([]tsValue,len(args));for i,arg:=range args{values[i]=tsNodeSQLiteJSValue(arg,useBigInt)}
		value:=tsCall(loop,fn,values...)
		return tsNodeSQLiteDriverValue(loop,value)
	},deterministic,directOnly)
}
type tsNodeSQLiteAggregateRunner struct{loop *tsLoop;db *tsNodeSQLiteDB;start,step,inverse,result tsValue;useBigInt,varargs bool;arity int;initialized bool;value tsValue}
func(runner *tsNodeSQLiteAggregateRunner)initialize(){if runner.initialized{return};runner.initialized=true;runner.value=runner.start;if runner.start.kind==tsFunctionKind{runner.value=tsCall(runner.loop,runner.start)}}
func(runner *tsNodeSQLiteAggregateRunner)Step(args ...any)(err error){defer func(){if failure:=recover();failure!=nil{runner.db.callbackFailure=failure;err=fmt.Errorf("aggregate step callback failed: %v",failure)}}();runner.initialize();if !runner.varargs&&len(args)!=runner.arity{return fmt.Errorf("aggregate step expected %d SQL arguments, received %d",runner.arity,len(args))};values:=make([]tsValue,1,len(args)+1);values[0]=runner.value;for _,arg:=range args{values=append(values,tsNodeSQLiteJSValue(arg,runner.useBigInt))};runner.value=tsCall(runner.loop,runner.step,values...);return nil}
func(runner *tsNodeSQLiteAggregateRunner)Inverse(args ...any)(err error){defer func(){if failure:=recover();failure!=nil{runner.db.callbackFailure=failure;err=fmt.Errorf("aggregate inverse callback failed: %v",failure)}}();runner.initialize();if !runner.varargs&&len(args)!=runner.arity{return fmt.Errorf("aggregate inverse expected %d SQL arguments, received %d",runner.arity,len(args))};values:=make([]tsValue,1,len(args)+1);values[0]=runner.value;for _,arg:=range args{values=append(values,tsNodeSQLiteJSValue(arg,runner.useBigInt))};runner.value=tsCall(runner.loop,runner.inverse,values...);return nil}
func(runner *tsNodeSQLiteAggregateRunner)Value()(result any,err error){defer func(){if failure:=recover();failure!=nil{runner.db.callbackFailure=failure;err=fmt.Errorf("aggregate result callback failed: %v",failure)}}();runner.initialize();value:=runner.value;if runner.result.kind==tsFunctionKind{value=tsCall(runner.loop,runner.result,value)};return tsNodeSQLiteDriverValue(runner.loop,value)}
func(runner *tsNodeSQLiteAggregateRunner)Done()(any,error){return runner.Value()}
func tsNodeSQLiteRegisterAggregate(loop *tsLoop,state *tsNodeSQLiteDB,name string,options tsValue)error{
	start:=tsNodeSQLiteGet(loop,options,"start");if tsIsUndefined(start){return fmt.Errorf("The \"options.start\" argument must be a function or a primitive value.")};if start.kind!=tsFunctionKind&&!tsIsPrimitiveValue(start){return fmt.Errorf("The \"options.start\" argument must be a function or a primitive value.")}
	step:=tsNodeSQLiteGet(loop,options,"step");if step.kind!=tsFunctionKind{return fmt.Errorf("The \"options.step\" argument must be a function")}
	result:=tsNodeSQLiteGet(loop,options,"result");if !tsIsUndefined(result)&&result.kind!=tsFunctionKind{return fmt.Errorf("The \"options.result\" argument must be a function")}
	inverse:=tsNodeSQLiteGet(loop,options,"inverse");if !tsIsUndefined(inverse)&&inverse.kind!=tsFunctionKind{return fmt.Errorf("The \"options.inverse\" argument must be a function")}
	useBigInt:=tsTruthy(tsNodeSQLiteGet(loop,options,"useBigIntArguments"));varargs:=tsTruthy(tsNodeSQLiteGet(loop,options,"varargs"))
	length:=tsNodeSQLiteGet(loop,step,"length");arity:=0;if tsIsNumeric(length){arity=max(int(tsNumber(length))-1,0)};inverseLength:=tsNodeSQLiteGet(loop,inverse,"length");if inverse.kind==tsFunctionKind&&tsIsNumeric(inverseLength){arity=max(arity,max(int(tsNumber(inverseLength))-1,0))}
	state.aggregates[strings.ToLower(name)]=options
	directOnly:=tsTruthy(tsNodeSQLiteGet(loop,options,"directOnly"))
	factory:=func()*tsNodeSQLiteAggregateRunner{return &tsNodeSQLiteAggregateRunner{loop:loop,db:state,start:start,step:step,inverse:inverse,result:result,useBigInt:useBigInt,varargs:varargs,arity:arity}}
	if inverse.kind==tsFunctionKind{return state.native.RegisterWindowAggregator(name,factory,directOnly)}
	return state.native.RegisterAggregatorWithOptions(name,factory,false,directOnly)
}
func tsNodeSQLiteQueryRows(state *tsNodeSQLiteStatement,args []any)(*sql.Rows,[]string,error){if state==nil||state.stmt==nil{return nil,nil,fmt.Errorf("statement is not prepared")};rows,err:=state.stmt.QueryContext(context.Background(),args...);if err!=nil{return nil,nil,err};return rows,tsNodeSQLiteCollectMetadata(rows,state),nil}
func tsNodeSQLiteRowsIterator(loop *tsLoop,rows *sql.Rows,state *tsNodeSQLiteStatement,columns []string)tsValue{
	object:=tsNewObject();closed:=false;closeRows:=func(){if !closed{closed=true;rows.Close();state.iteratorClose=nil;state.stepping=false}};state.iteratorClose=closeRows
	object.set("next",tsNamedNativeMethod("next",0,func(_ *tsLoop,_ tsValue,_ ...tsValue)tsValue{if closed{return tsIteratorResult(tsU,true)};if !rows.Next(){if err:=rows.Err();err!=nil{closeRows();tsNodeSQLiteRethrowCallbackFailure(tsNodeSQLiteDBFor(state.db));panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};closeRows();return tsIteratorResult(tsU,true)};value,err:=tsNodeSQLiteScanRow(rows,columns,state);if err!=nil{closeRows();tsNodeSQLiteRethrowCallbackFailure(tsNodeSQLiteDBFor(state.db));panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};return tsIteratorResult(tsNodeSQLiteRowObject(loop,state.columns,value,state.returnArrays),false)}))
	object.set("return",tsNamedNativeMethod("return",0,func(_ *tsLoop,_ tsValue,args ...tsValue)tsValue{closeRows();return tsIteratorResult(tsArg(args,0),true)}));tsSet(loop,tsObjectValue(object),tsWellKnownSymbol("iterator"),tsNamedNativeMethod("[Symbol.iterator]",0,func(_ *tsLoop,receiver tsValue,_ ...tsValue)tsValue{return receiver}));return tsObjectValue(object)
}
func tsNodeSQLiteDatabase(loop *tsLoop,args []tsValue)tsValue{
	location:=tsNodeSQLitePath(loop,tsArg(args,0))
	options:=tsObjectValue(tsNewObject())
	if len(args)>1{options=tsArg(args,1);if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}}
	db:=&tsNodeSQLiteDB{location:location,opened:true,functions:map[string]tsValue{},functionOptions:map[string]tsValue{},aggregates:map[string]tsValue{},authorizer:tsNull,limits:tsNodeSQLiteDefaultLimits(),limitMax:tsNodeSQLiteDefaultLimits(),statements:map[*tsObject]bool{},defensive:true,allowBare:true}
	readBool:=func(name string,def bool)bool{v:=tsNodeSQLiteGet(loop,options,name);if tsIsUndefined(v){return def};if v.kind!=tsBooleanKind{tsFSInvalid("options."+name,"a boolean")};return tsTruthy(v)}
	db.opened=readBool("open",true)
	db.readonly=readBool("readOnly",false)
	db.foreignKeys=readBool("enableForeignKeyConstraints",true)
	db.doubleQuoted=readBool("enableDoubleQuotedStringLiterals",false)
	db.allowExtension=readBool("allowExtension",false)
	db.extensions=db.allowExtension
	db.readBigInts=readBool("readBigInts",false)
	db.returnArrays=readBool("returnArrays",false)
	db.allowBare=readBool("allowBareNamedParameters",true)
	db.allowUnknown=readBool("allowUnknownNamedParameters",false)
	db.defensive=readBool("defensive",true)
	limits:=tsNodeSQLiteGet(loop,options,"limits")
	if !tsIsUndefined(limits){
		if limits.kind!=tsObjectKind&&limits.kind!=tsInstanceKind{tsFSInvalid("options.limits","an object")}
		for _,name:=range []string{"length","sqlLength","column","exprDepth","compoundSelect","vdbeOp","functionArg","attach","likePatternLength","variableNumber","triggerDepth"}{
			value:=tsNodeSQLiteGet(loop,limits,name);if tsIsUndefined(value){continue}
			if !tsIsNumeric(value)||math.IsNaN(tsNumber(value))||math.IsInf(tsNumber(value),0)||math.Trunc(tsNumber(value))!=tsNumber(value)||tsNumber(value)<math.MinInt32||tsNumber(value)>math.MaxInt32{tsFSInvalid("options.limits."+name,"an integer")}
			limit:=int(tsNumber(value));if limit<0{panic(tsThrown{tsNodeSQLiteError("ERR_OUT_OF_RANGE","The \"options.limits."+name+"\" argument must be non-negative.")})};if limit>db.limitMax[name]{limit=db.limitMax[name]};db.limits[name]=limit
		}
	}
	timeout:=tsNodeSQLiteGet(loop,options,"timeout")
	if !tsIsUndefined(timeout){if !tsIsNumeric(timeout)||math.IsNaN(tsNumber(timeout))||math.IsInf(tsNumber(timeout),0)||math.Trunc(tsNumber(timeout))!=tsNumber(timeout)||tsNumber(timeout)<math.MinInt32||tsNumber(timeout)>math.MaxInt32{tsFSInvalid("options.timeout","an integer")};db.timeout=int(tsNumber(timeout));if db.timeout<0{panic(tsThrown{tsNodeSQLiteError("ERR_OUT_OF_RANGE","The \"options.timeout\" argument must be a non-negative integer.")})}}
	obj:=tsNewObject();obj.prototype=tsObjectValue(tsNodeSQLiteSingleton.dbProto);db.object=obj
	for _,property:=range []string{"isOpen","isTransaction","limits"}{name:=property;obj.order=append(obj.order,name);tsNodeSQLiteAccessor(obj,name,func(loop *tsLoop,r tsValue)tsValue{current:=tsNodeSQLiteDBFor(r);if current==nil{return tsU};switch name{case "isOpen":return tsBooleanValue(current.opened);case "isTransaction":if !current.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};return tsBooleanValue(current.transaction);case "limits":return tsNodeSQLiteLimits(loop,r)};return tsU},nil);obj.descriptors[name].enumerable=true;obj.descriptors[name].configurable=false}
	tsNodeSQLiteSingleton.dbs[obj]=db
	if db.opened{if err:=tsNodeSQLiteOpen(loop,db);err!=nil{delete(tsNodeSQLiteSingleton.dbs,obj);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}}
	return tsObjectValue(obj)
}
func tsNodeSQLiteDatabaseMethod(loop *tsLoop,receiver tsValue,name string,args []tsValue)tsValue{
	db:=tsNodeSQLiteDBFor(receiver)
	if db==nil{tsPropertyFailure("DatabaseSync method called on an incompatible receiver")}
	if name!="open"&&!db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})}
	switch name{
	case "open":
		if db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is already open")})}
		if err:=tsNodeSQLiteOpen(loop,db);err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		return tsU
	case "close":
		if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})}
		tsNodeSQLiteClose(db)
		for _,state:=range tsNodeSQLiteSingleton.sessions{if state!=nil&&tsNodeSQLiteDBFor(state.db)==db{state.closed=true}}
		return tsU
	case "exec":
		sqlValue:=tsArg(args,0);if sqlValue.kind!=tsStringKind{tsFSInvalid("sql","a string")}
		result,err:=db.conn.ExecContext(context.Background(),tsText(sqlValue))
		if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		tsNodeSQLiteRethrowCallbackFailure(db)
		if result!=nil{db.lastInsertRowID,_=result.LastInsertId()};if db.native!=nil{db.transaction=!db.native.AutoCommit()}
		return tsU
	case "prepare":
		return tsNodeSQLiteStatement(loop,receiver,tsArg(args,0),tsArg(args,1))
	case "location":
		return tsNodeSQLiteKey(db.location)
	case "enableDefensive":
		active:=tsArg(args,0);if active.kind!=tsBooleanKind{tsFSInvalid("active","a boolean")}
		enabled:=tsTruthy(active);err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};return native.SetDefensive(enabled)});if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};db.defensive=enabled
		return tsU
	case "enableLoadExtension":
		allow:=tsArg(args,0);if allow.kind!=tsBooleanKind{tsFSInvalid("allow","a boolean")}
		enabled:=tsTruthy(allow);if enabled&&!db.allowExtension{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","Cannot enable extension loading because it was disabled at database creation.")})}
		err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};return native.EnableLoadExtension(enabled)});if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};db.extensions=enabled
		return tsU
	case "loadExtension":
		if !db.extensions{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","extension loading is not allowed")})}
		path:=tsArg(args,0);if path.kind!=tsStringKind{tsFSInvalid("path","a string")}
		entry:=tsArg(args,1);if !tsIsUndefined(entry)&&entry.kind!=tsStringKind{tsFSInvalid("entryPoint","a string")}
		err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};if tsIsUndefined(entry){return native.LoadExtensionDefault(tsText(path))};return native.LoadExtensionWithEntry(tsText(path),tsText(entry))});if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_LOAD_SQLITE_EXTENSION",err.Error())})};return tsU
	case "function":
		if len(args)<2||tsArg(args,0).kind!=tsStringKind{tsFSInvalid("name","a string")}
		nameValue:=tsText(tsArg(args,0));fnIndex:=1;options:=tsObjectValue(tsNewObject())
		if len(args)>2{
			options=tsArg(args,1);if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}
			normalized:=tsNewObject();for _,key:=range []string{"useBigIntArguments","varargs","deterministic","directOnly"}{value:=tsNodeSQLiteGet(loop,options,key);if !tsIsUndefined(value){if value.kind!=tsBooleanKind{tsFSInvalid("options."+key,"a boolean")};normalized.set(key,value)}};options=tsObjectValue(normalized)
			fnIndex=2
		}
		fn:=tsArg(args,fnIndex);if fn.kind!=tsFunctionKind{tsFSInvalid("function","a function")}
		db=tsNodeSQLiteRequireDB(receiver)
		key:=strings.ToLower(nameValue);db.functions[key]=fn;db.functionOptions[key]=options
		if err:=tsNodeSQLiteRegisterFunction(loop,db,nameValue,fn,options);err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		return tsU
	case "aggregate":
		if len(args)<2||tsArg(args,0).kind!=tsStringKind{tsFSInvalid("name","a string")}
		options:=tsArg(args,1);if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}
		normalized:=tsNewObject()
		for _,key:=range []string{"start","step","result"}{value:=tsNodeSQLiteGet(loop,options,key);if key=="start"{if tsIsUndefined(value){panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_TYPE","The \"options.start\" argument must be a function or a primitive value.")})};if value.kind!=tsFunctionKind&&!tsIsPrimitiveValue(value){tsFSInvalid("options.start","a function or a primitive value")}}else if key=="step"{if value.kind!=tsFunctionKind{tsFSInvalid("options.step","a function")}}else if !tsIsUndefined(value)&&value.kind!=tsFunctionKind{tsFSInvalid("options.result","a function")};normalized.set(key,value)}
		for _,key:=range []string{"useBigIntArguments","varargs","directOnly"}{value:=tsNodeSQLiteGet(loop,options,key);if !tsIsUndefined(value){if value.kind!=tsBooleanKind{tsFSInvalid("options."+key,"a boolean")};normalized.set(key,value)}}
		inverse:=tsNodeSQLiteGet(loop,options,"inverse");if !tsIsUndefined(inverse)&&inverse.kind!=tsFunctionKind{tsFSInvalid("options.inverse","a function")};normalized.set("inverse",inverse)
		db=tsNodeSQLiteRequireDB(receiver)
		if _,err:=tsNodeSQLiteRegisterAggregate(loop,db,tsText(tsArg(args,0)),tsObjectValue(normalized));err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		return tsU
	case "createSession":
		options:=tsArg(args,0);table,database:="","main"
		if len(args)>0{if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};for _,key:=range []string{"table","db"}{value:=tsNodeSQLiteGet(loop,options,key);if !tsIsUndefined(value)&&value.kind!=tsStringKind{tsFSInvalid("options."+key,"a string")};if key=="table"&&!tsIsUndefined(value){table=tsText(value)};if key=="db"&&!tsIsUndefined(value){database=tsText(value)}};db=tsNodeSQLiteRequireDB(receiver)}
		var nativeSession *sqlite3.SQLiteSession
		err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};var err error;nativeSession,err=native.Session(database);if err==nil{err=nativeSession.Attach(table);if err!=nil{_ = nativeSession.Close()}};return err})
		if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		session:=tsNewObject();session.prototype=tsObjectValue(tsNodeSQLiteSingleton.sessionProto)
		state:=&tsNodeSQLiteSession{object:session,db:receiver,native:nativeSession}
		tsNodeSQLiteSingleton.sessions[session]=state
		return tsObjectValue(session)
	case "applyChangeset":
		changeset:=tsArg(args,0);if changeset.kind!=tsTypedArrayKind{tsFSInvalid("changeset","a Uint8Array")}
		view:=(*tsTypedArray)(changeset.ref);if view.name!="Uint8Array"&&!view.nodeBuffer{tsFSInvalid("changeset","a Uint8Array")}
		options:=tsArg(args,1);onConflict,filter:=tsU,tsU
		if !tsIsUndefined(options){if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};onConflict=tsNodeSQLiteGet(loop,options,"onConflict");if !tsIsUndefined(onConflict)&&onConflict.kind!=tsFunctionKind{tsFSInvalid("options.onConflict","a function")};filter=tsNodeSQLiteGet(loop,options,"filter");if !tsIsUndefined(filter)&&filter.kind!=tsFunctionKind{tsFSInvalid("options.filter","a function")}}
		if view.buffer.detached{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_VALUE","The \"changeset\" argument must not be detached.")})}
		db=tsNodeSQLiteRequireDB(receiver)
		callbacks:=&sqlite3.SQLiteChangesetCallbacks{}
		if filter.kind==tsFunctionKind{callbacks.Filter=func(table string)bool{return tsTruthy(tsCall(loop,filter,tsNodeSQLiteKey(table)))}}
		if onConflict.kind==tsFunctionKind{callbacks.Conflict=func(kind int)int{value:=tsCall(loop,onConflict,tsNumberValue(float64(kind)));if !tsIsNumeric(value){return -1};result:=tsNumber(value);if math.IsNaN(result)||math.IsInf(result,0)||math.Trunc(result)!=result||result<math.MinInt32||result>math.MaxInt32{return -1};decision:=int(result);if decision!=sqlite3.SQLITE_CHANGESET_OMIT&&decision!=sqlite3.SQLITE_CHANGESET_REPLACE&&decision!=sqlite3.SQLITE_CHANGESET_ABORT{return -1};return decision}}
		if len(tsNodeBufferVisibleBytes(changeset))>math.MaxInt32{panic(tsThrown{tsNodeSQLiteError("ERR_OUT_OF_RANGE","The \"changeset\" argument is too large for SQLite.")})}
		err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};return native.ApplyChangeset(tsNodeBufferVisibleBytes(changeset),callbacks)})
		if callbacks.Panic!=nil{panic(callbacks.Panic)}
		if err!=nil{var sqliteError sqlite3.Error;if errors.As(err,&sqliteError)&&sqliteError.Code==sqlite3.ErrAbort{return tsBooleanValue(false)};panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		return tsBooleanValue(true)
	case "setAuthorizer":
		callback:=tsArg(args,0);if callback.kind!=tsNullKind&&callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function or null")}
		db.authorizer=callback;tsNodeSQLiteInstallAuthorizer(loop,db)
		return tsU
	case "serialize":
		database:="main";if value:=tsArg(args,0);!tsIsUndefined(value){if value.kind!=tsStringKind{tsFSInvalid("dbName","a string")};database=tsText(value)}
		var data []byte
		err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};var err error;data,err=native.Serialize(database);return err})
		if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		tsNodeSQLiteRethrowCallbackFailure(db)
		return tsNodeBuffer(data)
	case "deserialize":
		for object:=range db.statements{if statement:=tsNodeSQLiteSingleton.statements[object];statement!=nil&&statement.stepping{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","cannot deserialize database while a statement is executing")})}}
		data:=tsArg(args,0);if data.kind!=tsTypedArrayKind{tsFSInvalid("buffer","a Uint8Array")}
		view:=(*tsTypedArray)(data.ref);if view.name!="Uint8Array"&&!view.nodeBuffer{tsFSInvalid("buffer","a Uint8Array")}
		if view.buffer.detached||view.length==0{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_VALUE","The \"buffer\" argument must not be empty.")})}
		database:="main";options:=tsArg(args,1);if !tsIsUndefined(options){if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};name:=tsNodeSQLiteGet(loop,options,"dbName");if !tsIsUndefined(name){if name.kind!=tsStringKind{tsFSInvalid("options.dbName","a string")};database=tsText(name)}}
		db=tsNodeSQLiteRequireDB(receiver)
		if view.buffer.detached||view.length==0{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_VALUE","The \"buffer\" argument must not be empty.")})}
		for object:=range db.statements{if statement:=tsNodeSQLiteSingleton.statements[object];statement!=nil{if statement.iteratorClose!=nil{statement.iteratorClose()};statement.finalized=true;if statement.stmt!=nil{statement.stmt.Close();statement.stmt=nil}}}
		serialized:=append([]byte(nil),tsNodeBufferVisibleBytes(data)...)
		err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};return native.Deserialize(serialized,database)})
		if err!=nil{tsNodeSQLiteRethrowCallbackFailure(db);panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
		tsNodeSQLiteRethrowCallbackFailure(db)
		return tsU
	case "createTagStore":
		capacity:=uint64(1000);if value:=tsArg(args,0);tsIsNumeric(value){capacity=uint64(int64(int32(tsToUint32(value))))}
		object:=tsNewObject();object.prototype=tsObjectValue(tsNodeSQLiteSingleton.tagProto)
		state:=&tsNodeSQLiteTagStore{object:object,db:receiver,capacity:capacity,cache:map[string]tsValue{}}
		for _,property:=range []string{"capacity","size","db"}{key:=property;object.order=append(object.order,key);tsNodeSQLiteAccessor(object,key,func(_ *tsLoop,_ tsValue)tsValue{switch key{case "capacity":return tsNumberValue(float64(state.capacity));case "size":return tsNumberValue(float64(len(state.cache)));case "db":return state.db};return tsU},nil);object.descriptors[key].enumerable=true;object.descriptors[key].configurable=false}
		tsNodeSQLiteSingleton.tags[object]=state
		return tsObjectValue(object)
	case "isInTransaction":
		if db.native!=nil{db.transaction=!db.native.AutoCommit()}
		return tsBooleanValue(db.transaction)
	case "setCustomSQLite":
		panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR","Custom SQLite engines are not supported by the bundled SQLite driver")})
	case "fcntl":
		return tsNumberValue(0)
	}
	return tsU
}
func tsNodeSQLiteStatementMethod(loop *tsLoop,receiver tsValue,name string,args []tsValue)tsValue{
	state:=tsNodeSQLiteStatementFor(receiver)
	if state==nil{tsPropertyFailure("StatementSync method called on an incompatible receiver")}
	if state.finalized||state.stmt==nil{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})}
	switch name{
	case "run","get","all","iterate":
		return tsNodeSQLiteExecute(loop,receiver,args,name)
	case "columns":
		out:=&tsArray{}
		for index,columnName:=range state.columns{
			item:=tsNewObject();item.prototype=tsNull
			origin,database,table,declared:="","","","";if index<len(state.origins){origin=state.origins[index]};if index<len(state.databases){database=state.databases[index]};if index<len(state.tables){table=state.tables[index]};if index<len(state.declared){declared=state.declared[index]}
			put:=func(key,value string){if value==""{item.set(key,tsNull)}else{item.set(key,tsNodeSQLiteKey(value))}}
			put("column",origin);put("database",database);put("name",columnName);put("table",table);put("type",declared)
			out.values=append(out.values,tsObjectValue(item))
		}
		return tsArrayValue(out)
	case "setReadBigInts","setReturnArrays","setAllowBareNamedParameters","setAllowUnknownNamedParameters":
		value:=tsArg(args,0);if value.kind!=tsBooleanKind{tsFSInvalid("enabled","a boolean")}
		enabled:=tsTruthy(value)
		switch name{case "setReadBigInts":state.safeIntegers=enabled;case "setReturnArrays":state.returnArrays=enabled;case "setAllowBareNamedParameters":state.allowBare=enabled;case "setAllowUnknownNamedParameters":state.allowUnknown=enabled}
		return tsU
	}
	return tsU
}
func tsNodeSQLiteColumnTypes(state *tsNodeSQLiteStatement)tsValue{
	if !state.readonly{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","columnTypes is not available for non-read-only statements (INSERT, UPDATE, DELETE, etc.)")})}
	bindings:=make([]any,state.params);rows,columns,err:=tsNodeSQLiteQueryRows(state,bindings);if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};defer rows.Close()
	values:=make([]tsValue,len(columns));for i:=range values{values[i]=tsNodeSQLiteKey("NULL")}
	if rows.Next(){raw:=make([]any,len(columns));targets:=make([]any,len(columns));for i:=range raw{targets[i]=&raw[i]};if err:=rows.Scan(targets...);err!=nil{tsNodeSQLiteRethrowCallbackFailure(tsNodeSQLiteDBFor(state.db));panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};for i,item:=range raw{switch item.(type){case nil:values[i]=tsNodeSQLiteKey("NULL");case int64,int:values[i]=tsNodeSQLiteKey("INTEGER");case float64:values[i]=tsNodeSQLiteKey("FLOAT");case []byte:values[i]=tsNodeSQLiteKey("BLOB");case string:values[i]=tsNodeSQLiteKey("TEXT")}}}else if err:=rows.Err();err!=nil{tsNodeSQLiteRethrowCallbackFailure(tsNodeSQLiteDBFor(state.db));panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}
	return tsArrayValue(&tsArray{values:values})
}
func tsNodeSQLiteSessionMethod(loop *tsLoop,r tsValue,name string)tsValue{
	state:=tsNodeSQLiteSessionFor(r)
	if state==nil{tsPropertyFailure("Session method called on an incompatible receiver")}
	db:=tsNodeSQLiteDBFor(state.db)
	if name=="dispose"{if state.closed{return tsU};if state.native!=nil{_ = state.native.Close();state.native=nil};state.closed=true;return tsU}
	if db==nil||!db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})}
	if state.closed||state.native==nil{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","session is not open")})}
	switch name{case "changeset":data,err:=state.native.Changeset();if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};return tsNodeBuffer(data);case "patchset":data,err:=state.native.Patchset();if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};return tsNodeBuffer(data);case "close":if err:=state.native.Close();err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};state.native=nil;state.closed=true;return tsU}
	return tsU
}
func tsNodeSQLiteTagMethod(loop *tsLoop,r tsValue,name string,args []tsValue)tsValue{
	state:=tsNodeSQLiteTagFor(r)
	if state==nil{tsPropertyFailure("SQLTagStore method called on an incompatible receiver")}
	switch name{
	case "clear":state.cache=map[string]tsValue{};state.order=nil;return tsU
	case "run","get","all","iterate":
		db:=tsNodeSQLiteRequireDB(state.db)
		_ = db
		template:=tsArg(args,0)
		if template.kind!=tsArrayKind{tsFSInvalid("strings","an array of strings")}
		parts:=(*tsArray)(template.ref);var sql strings.Builder
		for i,part:=range parts.values{if part.kind!=tsStringKind{tsFSInvalid("strings","an array of strings")};sql.WriteString(tsText(part));if i<len(args)-1{sql.WriteByte('?')}}
		text:=sql.String();stmt:=state.cache[text]
		if cached:=tsNodeSQLiteStatementFor(stmt);cached!=nil&&cached.finalized{delete(state.cache,text);for i,key:=range state.order{if key==text{state.order=append(state.order[:i],state.order[i+1:]...);break}};stmt=tsU}
		if tsIsUndefined(stmt){
			stmt=tsNodeSQLiteStatement(loop,state.db,tsNodeSQLiteKey(text),tsU);if prepared:=tsNodeSQLiteStatementFor(stmt);prepared!=nil{prepared.tagStore=true}
			if state.capacity>0{state.cache[text]=stmt;state.order=append(state.order,text);if uint64(len(state.order))>state.capacity{oldest:=state.order[0];delete(state.cache,oldest);state.order=state.order[1:]}}
		}else{
			for i,key:=range state.order{if key==text{state.order=append(state.order[:i],state.order[i+1:]...);break}}
			state.order=append(state.order,text)
		}
		return tsNodeSQLiteExecute(loop,stmt,args[1:],name)
	}
	return tsU
}
func tsNodeSQLiteBackupCopy(loop *tsLoop,source *tsNodeSQLiteDB,targetPath,sourceName,targetName string,rate int,progress tsValue)(pages int,err error){
	defer func(){if failure:=recover();failure!=nil{err=fmt.Errorf("backup progress callback failed: %v",failure)}}()
	if source==nil||source.conn==nil||!source.opened{return 0,fmt.Errorf("database is not open")}
	targetPool,err:=sql.Open("sqlite3",tsNodeSQLiteDSN(targetPath,false,true,5000));if err!=nil{return 0,err};defer targetPool.Close()
	targetPool.SetMaxOpenConns(1);targetPool.SetMaxIdleConns(1)
	targetConn,err:=targetPool.Conn(context.Background());if err!=nil{return 0,err};defer targetConn.Close()
	if rate==0{rate=1}
	err=source.conn.Raw(func(sourceDriver any)error{
		sourceNative,ok:=sourceDriver.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected source SQLite connection type %T",sourceDriver)}
		return targetConn.Raw(func(targetDriver any)error{
			targetNative,ok:=targetDriver.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected target SQLite connection type %T",targetDriver)}
			backup,err:=targetNative.Backup(targetName,sourceNative,sourceName);if err!=nil{return err}
			closed:=false;defer func(){if !closed{_ = backup.Close()}}()
			for{
				done,stepErr:=backup.Step(rate);pages=backup.PageCount();remaining:=backup.Remaining()
				if stepErr!=nil{return stepErr}
				if remaining!=0{
					if progress.kind==tsFunctionKind{payload:=tsNewObject();payload.set("totalPages",tsNumberValue(float64(pages)));payload.set("remainingPages",tsNumberValue(float64(remaining)));tsCall(loop,progress,tsObjectValue(payload))}
					time.Sleep(25*time.Millisecond)
					continue
				}
				if done{break}
				return fmt.Errorf("SQLite backup did not advance")
			}
			closeErr:=backup.Close();closed=true;return closeErr
		})
	})
	return pages,err
}
func tsNodeSQLiteBackup(loop *tsLoop,args []tsValue)tsValue{
	source:=tsArg(args,0);if tsNodeSQLiteDBFor(source)==nil{tsFSInvalid("sourceDb","a DatabaseSync")}
	db:=tsNodeSQLiteRequireDB(source);targetPath:=tsNodeSQLitePath(loop,tsArg(args,1));rate:=100;sourceName,targetName:="main","main";progress:=tsU
	if len(args)>2{
		options:=args[2];if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}
		rateValue:=tsNodeSQLiteGet(loop,options,"rate");if !tsIsUndefined(rateValue){if !tsIsNumeric(rateValue)||math.IsNaN(tsNumber(rateValue))||math.IsInf(tsNumber(rateValue),0)||math.Trunc(tsNumber(rateValue))!=tsNumber(rateValue)||tsNumber(rateValue)<math.MinInt32||tsNumber(rateValue)>math.MaxInt32{tsFSInvalid("options.rate","an integer")};rate=int(tsNumber(rateValue));if rate==0{rate=1}}
		for _,entry:=range []struct{name string;target *string}{{"source",&sourceName},{"target",&targetName}}{value:=tsNodeSQLiteGet(loop,options,entry.name);if !tsIsUndefined(value){if value.kind!=tsStringKind{tsFSInvalid("options."+entry.name,"a string")};*entry.target=tsText(value)}}
		progress=tsNodeSQLiteGet(loop,options,"progress");if !tsIsUndefined(progress)&&progress.kind!=tsFunctionKind{tsFSInvalid("options.progress","a function")}
	}
	if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})}
	value:=tsU;rejected:=false
	func(){defer func(){if failure:=recover();failure!=nil{value=tsUnwrap(failure);rejected=true}}();pages,err:=tsNodeSQLiteBackupCopy(loop,db,targetPath,sourceName,targetName,rate,progress);if err!=nil{value=tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error());rejected=true;return};value=tsNumberValue(float64(pages))}()
	return tsPromiseValue(loop.resolved(value,rejected))
}
func tsNodeSQLiteInit(loop *tsLoop){tsNodeSQLiteSingleton.Do(func(){tsNodeSQLiteSingleton.dbs=map[*tsObject]*tsNodeSQLiteDB{};tsNodeSQLiteSingleton.statements=map[*tsObject]*tsNodeSQLiteStatement{};tsNodeSQLiteSingleton.sessions=map[*tsObject]*tsNodeSQLiteSession{};tsNodeSQLiteSingleton.tags=map[*tsObject]*tsNodeSQLiteTagStore{};tsNodeSQLiteSingleton.dbProto=tsNewObject();tsNodeSQLiteSingleton.stmtProto=tsNewObject();tsNodeSQLiteSingleton.sessionProto=tsNewObject();tsNodeSQLiteSingleton.tagProto=tsNewObject()
 database:=tsNodeSQLiteConstructor("DatabaseSync",tsNodeSQLiteSingleton.dbProto,2,func(loop *tsLoop,args []tsValue)tsValue{return tsNodeSQLiteDatabase(loop,args)});tsNodeSQLiteSingleton.dbClass=database
 statement:=tsNodeSQLiteConstructor("StatementSync",tsNodeSQLiteSingleton.stmtProto,0,func(_ *tsLoop,_ []tsValue)tsValue{panic(tsThrown{tsFSException("TypeError","ERR_ILLEGAL_CONSTRUCTOR","Illegal constructor")})});tsNodeSQLiteSingleton.stmtClass=statement
 session:=tsNodeSQLiteConstructor("Session",tsNodeSQLiteSingleton.sessionProto,0,func(_ *tsLoop,_ []tsValue)tsValue{panic(tsThrown{tsFSException("TypeError","ERR_ILLEGAL_CONSTRUCTOR","Illegal constructor")})});tsNodeSQLiteSingleton.sessionClass=session
 for _,entry:=range []struct{name string;length int}{{"open",0},{"close",0},{"exec",1},{"prepare",1},{"location",0},{"enableLoadExtension",1},{"enableDefensive",1},{"loadExtension",1},{"function",2},{"aggregate",2},{"createSession",0},{"applyChangeset",1},{"setAuthorizer",1},{"serialize",0},{"deserialize",1},{"createTagStore",0}}{n,length:=entry.name,entry.length;tsNodeSQLiteSingleton.dbProto.set(n,tsNamedNativeMethod(n,length,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeSQLiteDatabaseMethod(loop,r,n,args)}))}
 for _,entry:=range []struct{name string;length int}{{"run",0},{"get",0},{"all",0},{"iterate",0},{"columns",0},{"setReadBigInts",1},{"setReturnArrays",1},{"setAllowBareNamedParameters",1},{"setAllowUnknownNamedParameters",1}}{n,length:=entry.name,entry.length;tsNodeSQLiteSingleton.stmtProto.set(n,tsNamedNativeMethod(n,length,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeSQLiteStatementMethod(loop,r,n,args)}))}
 for _,name:=range []string{"changeset","patchset","close"}{n:=name;tsNodeSQLiteSingleton.sessionProto.set(n,tsNamedNativeMethod(n,0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeSQLiteSessionMethod(loop,r,n)}))}
 for _,name:=range []string{"run","get","all","iterate","clear"}{n:=name;tsNodeSQLiteSingleton.tagProto.set(n,tsNamedNativeMethod(n,0,func(loop *tsLoop,r tsValue,args ...tsValue)tsValue{return tsNodeSQLiteTagMethod(loop,r,n,args)}))};tsSet(loop,tsObjectValue(tsNodeSQLiteSingleton.dbProto),tsWellKnownSymbol("dispose"),tsNamedNativeMethod("[Symbol.dispose]",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{if db:=tsNodeSQLiteDBFor(r);db!=nil&&db.opened{return tsNodeSQLiteDatabaseMethod(loop,r,"close",nil)};return tsU}));tsSet(loop,tsObjectValue(tsNodeSQLiteSingleton.sessionProto),tsWellKnownSymbol("dispose"),tsNamedNativeMethod("[Symbol.dispose]",0,func(loop *tsLoop,r tsValue,_ ...tsValue)tsValue{return tsNodeSQLiteSessionMethod(loop,r,"dispose")}))
 for _,property:=range []string{"isOpen","isTransaction","limits"}{n:=property;tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.dbProto,n,func(loop *tsLoop,r tsValue)tsValue{db:=tsNodeSQLiteDBFor(r);if db==nil{return tsU};switch n{case "isOpen":return tsBooleanValue(db.opened);case "isTransaction":if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};return tsBooleanValue(db.transaction);case "limits":return tsNodeSQLiteLimits(loop,r)};return tsU},nil)}
 tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.stmtProto,"sourceSQL",func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeSQLiteStatementFor(r);if s==nil{return tsU};if s.finalized{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})};return tsNodeSQLiteKey(s.sql)},nil);tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.stmtProto,"expandedSQL",func(loop *tsLoop,r tsValue)tsValue{s:=tsNodeSQLiteStatementFor(r);if s==nil{return tsU};if s.finalized{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})};db:=tsNodeSQLiteRequireDB(s.db);if db.authorizer.kind==tsFunctionKind{db.native.RegisterAuthorizerWithSource(nil)};var expanded string;err:=db.conn.Raw(func(driverConn any)error{native,ok:=driverConn.(*sqlite3.SQLiteConn);if !ok{return fmt.Errorf("unexpected SQLite connection type %T",driverConn)};var err error;expanded,err=native.ExpandedSQL(s.sql,s.lastBindings);return err});if db.authorizer.kind==tsFunctionKind{tsNodeSQLiteInstallAuthorizer(loop,db)};if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};return tsNodeSQLiteKey(expanded)},nil);for _,property:=range []string{"columnNames","columnCount","parameterCount","isFinalized","columnTypes","declaredColumnTypes","safeIntegers"}{n:=property;tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.stmtProto,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeSQLiteStatementFor(r);if s==nil{return tsU};if n!="isFinalized"&&(s.finalized||s.stmt==nil){panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})};switch n{case "columnNames":a:=&tsArray{};for _,v:=range s.columns{a.values=append(a.values,tsNodeSQLiteKey(v))};return tsArrayValue(a);case "columnCount":return tsNumberValue(float64(len(s.columns)));case "parameterCount":return tsNumberValue(float64(s.params));case "isFinalized":return tsBooleanValue(s.finalized);case "columnTypes":return tsNodeSQLiteColumnTypes(s);case "declaredColumnTypes":if !s.executed{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","Statement must be executed before accessing declaredTypes")})};a:=&tsArray{};for _,v:=range s.declared{if v==""{a.values=append(a.values,tsNull)}else{a.values=append(a.values,tsNodeSQLiteKey(v))}};return tsArrayValue(a);case "safeIntegers":return tsBooleanValue(s.safeIntegers)};return tsU},func(_ *tsLoop,r,v tsValue)tsValue{if n=="safeIntegers"{if v.kind!=tsBooleanKind{tsFSInvalid("safeIntegers","a boolean")};if s:=tsNodeSQLiteStatementFor(r);s!=nil{s.safeIntegers=tsTruthy(v)}};return v})}
 tag:=tsNodeSQLiteSingleton.tagProto;for _,property:=range []string{"capacity","size","db"}{n:=property;tsNodeSQLiteAccessor(tag,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeSQLiteTagFor(r);if s==nil{return tsU};switch n{case "capacity":return tsNumberValue(float64(s.capacity));case "size":return tsNumberValue(float64(len(s.cache)));case "db":return s.db};return tsU},nil)}
 constants:=tsNewObject();for name,value:=range map[string]int{"SQLITE_CHANGESET_OMIT":0,"SQLITE_CHANGESET_REPLACE":1,"SQLITE_CHANGESET_ABORT":2,"SQLITE_CHANGESET_DATA":1,"SQLITE_CHANGESET_NOTFOUND":2,"SQLITE_CHANGESET_CONFLICT":3,"SQLITE_CHANGESET_CONSTRAINT":4,"SQLITE_CHANGESET_FOREIGN_KEY":5,"SQLITE_OK":0,"SQLITE_DENY":1,"SQLITE_IGNORE":2,"SQLITE_CREATE_INDEX":1,"SQLITE_CREATE_TABLE":2,"SQLITE_CREATE_TEMP_INDEX":3,"SQLITE_CREATE_TEMP_TABLE":4,"SQLITE_CREATE_TEMP_TRIGGER":5,"SQLITE_CREATE_TEMP_VIEW":6,"SQLITE_CREATE_TRIGGER":7,"SQLITE_CREATE_VIEW":8,"SQLITE_DELETE":9,"SQLITE_DROP_INDEX":10,"SQLITE_DROP_TABLE":11,"SQLITE_DROP_TEMP_INDEX":12,"SQLITE_DROP_TEMP_TABLE":13,"SQLITE_DROP_TEMP_TRIGGER":14,"SQLITE_DROP_TEMP_VIEW":15,"SQLITE_DROP_TRIGGER":16,"SQLITE_DROP_VIEW":17,"SQLITE_INSERT":18,"SQLITE_PRAGMA":19,"SQLITE_READ":20,"SQLITE_SELECT":21,"SQLITE_TRANSACTION":22,"SQLITE_UPDATE":23,"SQLITE_ATTACH":24,"SQLITE_DETACH":25,"SQLITE_ALTER_TABLE":26,"SQLITE_REINDEX":27,"SQLITE_ANALYZE":28,"SQLITE_CREATE_VTABLE":29,"SQLITE_DROP_VTABLE":30,"SQLITE_FUNCTION":31,"SQLITE_SAVEPOINT":32,"SQLITE_COPY":0,"SQLITE_RECURSIVE":33}{constants.set(name,tsNumberValue(float64(value)))};tsNodeSQLiteSingleton.constants=tsObjectValue(constants)
 module:=tsNewObject();module.set("DatabaseSync",database);module.set("StatementSync",statement);module.set("Session",session);module.set("constants",tsNodeSQLiteSingleton.constants);module.set("backup",tsNamedNativeMethod("backup",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeSQLiteBackup(loop,args)}));module.set("default",tsObjectValue(module));tsNodeSQLiteSingleton.module=tsObjectValue(module)
})}
func tsNodeSQLiteModule(loop *tsLoop)tsValue{tsNodeSQLiteInit(loop);return tsNodeSQLiteSingleton.module}
func tsNodeSQLiteAddon(loop *tsLoop,name string)(tsValue,bool){if name=="sqlite"{return tsNodeSQLiteModule(loop),true};return tsU,false}
`

package goemit

// nodeSQLiteRuntime maps Bun's native node:sqlite binding onto a local
// synchronous API facade. SQLite's C amalgamation is not linked into the
// generated Go runtime, so this adapter includes a small in-memory SQL subset
// and keeps every Node-facing class and method available.
const nodeSQLiteRuntime = `
type tsNodeSQLiteTable struct{name string;columns []string;declared []string;rows [][]tsValue;nextRowID int64}
type tsNodeSQLiteDB struct{object,limitsObject *tsObject;location string;opened,readonly,transaction,defensive,allowExtension,extensions,foreignKeys,doubleQuoted bool;timeout int;readBigInts,returnArrays,allowBare,allowUnknown bool;totalChanges,lastInsertRowID int64;tables map[string]*tsNodeSQLiteTable;functions,aggregates map[string]tsValue;authorizer tsValue;limits,limitMax map[string]int;statements map[*tsObject]bool}
type tsNodeSQLiteStatement struct{object *tsObject;db tsValue;sql string;finalized,safeIntegers,returnArrays,allowBare,allowUnknown,stepping bool;columns []string;declared []string;params int}
type tsNodeSQLiteSession struct{object *tsObject;db tsValue;closed bool;changeOffset int}
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
			if limit>current.limitMax[property]{limit=current.limitMax[property]};current.limits[property]=limit;return value
		})
		object.descriptors[property].enumerable=true;object.descriptors[property].configurable=false
	}
	db.limitsObject=object;return tsObjectValue(object)
}
func tsNodeSQLiteConstructor(name string,proto *tsObject,length int,construct func(*tsLoop,[]tsValue)tsValue)tsValue{fn:=tsFunc(func(_ *tsLoop,_ ...tsValue)tsValue{tsPropertyFailure("Class constructor "+name+" cannot be invoked without 'new'");return tsU});fn.name=name;fn.length=length;fn.constructible=true;fn.construct=func(loop *tsLoop,args ...tsValue)tsValue{return construct(loop,args)};fn.properties=tsNewObject();fn.properties.set("prototype",tsObjectValue(proto));fn.properties.descriptors=map[string]*tsDescriptor{"prototype":{writable:false}};proto.set("constructor",tsFunctionValue(fn));return tsFunctionValue(fn)}
func tsNodeSQLiteSplit(source string,separator byte)[]string{out:=[]string{};start:=0;quote:=byte(0);for i:=0;i<len(source);i++{c:=source[i];if quote!=0{if c==quote{if i+1<len(source)&&source[i+1]==quote{i++;continue};quote=0};continue};if c=='\''||c=='"'||c=='\x60'{quote=c;continue};if c==separator{out=append(out,strings.TrimSpace(source[start:i]));start=i+1}};tail:=strings.TrimSpace(source[start:]);if tail!=""{out=append(out,tail)};return out}
func tsNodeSQLiteIdentifier(source string)string{source=strings.TrimSpace(source);source=strings.Trim(source," \t\r\n\x60\"[](),;");if dot:=strings.LastIndex(source,".");dot>=0{source=source[dot+1:]};return strings.Trim(source," \t\r\n\x60\"[](),;")}
func tsNodeSQLiteTableNameAfter(source,keyword string)string{
	fields:=strings.Fields(strings.TrimSpace(source[len(keyword):]));index:=0
	if index<len(fields)&&(strings.EqualFold(fields[index],"temp")||strings.EqualFold(fields[index],"temporary")){index++}
	if index<len(fields)&&strings.EqualFold(fields[index],"if")&&index+2<len(fields){index+=3}
	if index>=len(fields){return ""};return tsNodeSQLiteIdentifier(fields[index])
}
func tsNodeSQLiteCallAuthorizer(loop *tsLoop,db *tsNodeSQLiteDB,action int,first,second,third,fourth tsValue)int{
	if db.authorizer.kind==tsNullKind||tsIsUndefined(db.authorizer){return 0}
	result:=tsCall(loop,db.authorizer,tsNumberValue(float64(action)),first,second,third,fourth)
	if !tsIsNumeric(result){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Authorizer callback must return an integer authorization code"})})}
	number:=tsNumber(result)
	if math.IsNaN(number)||math.IsInf(number,0)||math.Trunc(number)!=number||number<math.MinInt32||number>math.MaxInt32{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"TypeError",message:"Authorizer callback must return an integer authorization code"})})}
	code:=int32(number)
	if code!=0&&code!=1&&code!=2{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Authorizer callback returned a invalid authorization code"})})}
	return int(code)
}
func tsNodeSQLiteAuthorizeSQL(loop *tsLoop,db *tsNodeSQLiteDB,sql string)error{
	source:=strings.TrimSpace(strings.TrimSuffix(sql,";"));lower:=strings.ToLower(source);action:=0;first,second:=tsNull,tsNull
	switch{
	case strings.HasPrefix(lower,"create temp table")||strings.HasPrefix(lower,"create temporary table"):
		action=4;if strings.HasPrefix(lower,"create temporary table"){first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"create temporary table"))}else{first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"create temp table"))}
	case strings.HasPrefix(lower,"create table"):
		action=2;first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"create table"))
	case strings.HasPrefix(lower,"insert into")||strings.HasPrefix(lower,"replace into"):
		action=18;first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"insert into"));if strings.HasPrefix(lower,"replace into"){first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"replace into"))}
	case strings.HasPrefix(lower,"delete from"):
		action=9;fields:=strings.Fields(source);if len(fields)>2{first=tsNodeSQLiteKey(tsNodeSQLiteIdentifier(fields[2]))}
	case strings.HasPrefix(lower,"update "):
		action=23;rest:=strings.TrimSpace(source[6:]);fields:=strings.Fields(rest);if len(fields)>0{first=tsNodeSQLiteKey(tsNodeSQLiteIdentifier(fields[0]))}
	case strings.HasPrefix(lower,"drop temp table")||strings.HasPrefix(lower,"drop temporary table"):
		action=13;if strings.HasPrefix(lower,"drop temporary table"){first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"drop temporary table"))}else{first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"drop temp table"))}
	case strings.HasPrefix(lower,"drop table"):
		action=11;first=tsNodeSQLiteKey(tsNodeSQLiteTableNameAfter(source,"drop table"))
	case strings.HasPrefix(lower,"pragma "):
		action=19;parts:=strings.Fields(strings.TrimSpace(source[6:]));if len(parts)>0{first=tsNodeSQLiteKey(parts[0])};if len(parts)>1{second=tsNodeSQLiteKey(strings.Trim(parts[1],"=()"))}
	case strings.HasPrefix(lower,"begin"),strings.HasPrefix(lower,"start transaction"):
		action=22;first=tsNodeSQLiteKey("BEGIN")
	case strings.HasPrefix(lower,"commit"),strings.HasPrefix(lower,"end transaction"):
		action=22;first=tsNodeSQLiteKey("COMMIT")
	case strings.HasPrefix(lower,"rollback"):
		action=22;first=tsNodeSQLiteKey("ROLLBACK")
	case strings.HasPrefix(lower,"attach "):
		action=24;fields:=strings.Fields(source);if len(fields)>1{first=tsNodeSQLiteKey(strings.Trim(fields[1],"'\""))}
	case strings.HasPrefix(lower,"detach "):
		action=25;fields:=strings.Fields(source);if len(fields)>1{index:=1;if strings.EqualFold(fields[index],"database"){index++};if index<len(fields){first=tsNodeSQLiteKey(tsNodeSQLiteIdentifier(fields[index]))}}
	case strings.HasPrefix(lower,"select "):
		if tsNodeSQLiteCallAuthorizer(loop,db,21,tsNull,tsNull,tsNodeSQLiteKey("main"),tsNull)==1{return fmt.Errorf("not authorized")}
		fromAt:=strings.Index(lower," from ");projection:=strings.TrimSpace(source[6:]);tableName:=""
		if fromAt>=0{
			projection=strings.TrimSpace(source[6:fromAt]);tail:=strings.TrimSpace(source[fromAt+6:]);fields:=strings.Fields(tail);if len(fields)>0{tableName=tsNodeSQLiteIdentifier(fields[0])}
		}
		if strings.Contains(strings.ToLower(projection),"count("){
			if tsNodeSQLiteCallAuthorizer(loop,db,31,tsNull,tsNodeSQLiteKey("count"),tsNodeSQLiteKey("main"),tsNull)==1{return fmt.Errorf("not authorized")}
		}
		if tableName!=""{
			read:=func(column string)bool{return tsNodeSQLiteCallAuthorizer(loop,db,20,tsNodeSQLiteKey(tableName),tsNodeSQLiteKey(column),tsNodeSQLiteKey("main"),tsNull)==1}
			table:=db.tables[strings.ToLower(tableName)]
			if table!=nil{
				if strings.TrimSpace(projection)=="*"{for _,column:=range table.columns{if read(column){return fmt.Errorf("not authorized")}}}else if strings.Contains(strings.ToLower(projection),"count("){if read(""){return fmt.Errorf("not authorized")}}else{
					for _,item:=range tsNodeSQLiteSplit(projection,','){
						item=strings.TrimSpace(item);if as:=strings.Index(strings.ToLower(item)," as ");as>=0{item=item[:as]};if strings.ContainsAny(item,"()+ "){continue};if read(tsNodeSQLiteIdentifier(item)){return fmt.Errorf("not authorized")}
					}
				}
				whereAt:=strings.Index(strings.ToLower(source)," where ");if whereAt>=0{where:=strings.TrimSpace(source[whereAt+7:]);for _,operator:=range []string{" order by "," group by "," limit "," returning "}{if at:=strings.Index(strings.ToLower(where),operator);at>=0{where=where[:at]}};for _,condition:=range strings.Split(where," and "){parts:=strings.FieldsFunc(condition,func(r rune)bool{return r=='='||r=='<'||r=='>'||r=='!'||r==' '});if len(parts)>0&&read(tsNodeSQLiteIdentifier(parts[0])){return fmt.Errorf("not authorized")}}}
			}
		}
		return nil
	default:return nil
	}
	database:=tsNodeSQLiteKey("main");if action==4{database=tsNodeSQLiteKey("temp")}
	if tsNodeSQLiteCallAuthorizer(loop,db,action,first,second,database,tsNull)==1{return fmt.Errorf("not authorized")}
	return nil
}
func tsNodeSQLiteBindings(loop *tsLoop,args []tsValue)[]tsValue{out:=[]tsValue{};for _,arg:=range args{if arg.kind==tsArrayKind{out=append(out,(*tsArray)(arg.ref).values...)}else{out=append(out,arg)}};_ = loop;return out}
func tsNodeSQLiteLiteral(loop *tsLoop,text string,args []tsValue,index *int)tsValue{text=strings.TrimSpace(text);if text=="?"||strings.HasPrefix(text,"?"){i:=*index;*index++;if i<len(args){return args[i]};return tsU};if len(text)>1&&(text[0]==':'||text[0]=='@'||text[0]=='$'){for _,arg:=range args{if arg.kind==tsObjectKind{if value,ok:=(*tsObject)(arg.ref).values[strings.TrimLeft(text,":@$")];ok{return value}}}};text=strings.Trim(text,"() ");lower:=strings.ToLower(text);if lower=="null"{return tsNull};if lower=="true"{return tsBooleanValue(true)};if lower=="false"{return tsBooleanValue(false)};if lower=="undefined"{return tsU};if strings.HasPrefix(text,"'")&&strings.HasSuffix(text,"'")&&len(text)>=2{return tsNodeSQLiteKey(strings.ReplaceAll(text[1:len(text)-1],"''","'"))};if strings.HasPrefix(text,"\"")&&strings.HasSuffix(text,"\"")&&len(text)>=2{return tsNodeSQLiteKey(text[1:len(text)-1])};if number,err:=strconv.ParseInt(text,10,64);err==nil{return tsNumberValue(float64(number))};if number,err:=strconv.ParseFloat(text,64);err==nil{return tsNumberValue(number)};if strings.HasPrefix(text,"x'")&&strings.HasSuffix(text,"'"){data,err:=hex.DecodeString(text[2:len(text)-1]);if err==nil{return tsNodeBuffer(data)}};return tsNodeSQLiteKey(text)}
func tsNodeSQLiteRowObject(loop *tsLoop,columns []string,row []tsValue,arrays bool)tsValue{if arrays{return tsArrayValue(&tsArray{values:append([]tsValue(nil),row...)})};o:=tsNewObject();for i,name:=range columns{if i<len(row){o.set(name,row[i])}};return tsObjectValue(o)}
func tsNodeSQLiteResult(loop *tsLoop,db *tsNodeSQLiteDB,sql string,args []tsValue) ([][]tsValue,[]string,int64,error){source:=strings.TrimSpace(strings.TrimSuffix(sql,";"));lower:=strings.ToLower(source);paramIndex:=0
 if strings.HasPrefix(lower,"begin")||strings.HasPrefix(lower,"start transaction"){db.transaction=true;return nil,nil,0,nil};if lower=="commit"||strings.HasPrefix(lower,"end transaction"){db.transaction=false;return nil,nil,0,nil};if lower=="rollback"{db.transaction=false;return nil,nil,0,nil}
 if strings.HasPrefix(lower,"create table")||strings.HasPrefix(lower,"create temp table")||strings.HasPrefix(lower,"create temporary table"){open:=strings.Index(source,"(");close:=strings.LastIndex(source,")");if open<0||close<open{return nil,nil,0,fmt.Errorf("near \"%s\": syntax error",source)};prefix:=strings.TrimSpace(source[:open]);parts:=strings.Fields(prefix);name:=tsNodeSQLiteIdentifier(parts[len(parts)-1]);columns:=[]string{};declared:=[]string{};for _,column:=range tsNodeSQLiteSplit(source[open+1:close],','){fields:=strings.Fields(column);if len(fields)==0||strings.EqualFold(fields[0],"primary")||strings.EqualFold(fields[0],"unique")||strings.EqualFold(fields[0],"foreign")||strings.EqualFold(fields[0],"constraint"){continue};columns=append(columns,tsNodeSQLiteIdentifier(fields[0]));if len(fields)>1{declared=append(declared,strings.ToUpper(fields[1]))}else{declared=append(declared,"")}};key:=strings.ToLower(name);if _,exists:=db.tables[key];exists&&!strings.Contains(lower,"if not exists"){return nil,nil,0,fmt.Errorf("table %s already exists",name)};db.tables[key]=&tsNodeSQLiteTable{name:name,columns:columns,declared:declared,nextRowID:1};return nil,nil,0,nil}
 if strings.HasPrefix(lower,"insert into")||strings.HasPrefix(lower,"replace into") {rest:=strings.TrimSpace(source[strings.Index(lower,"into")+4:]);open:=strings.Index(rest,"(");valuesAt:=strings.Index(strings.ToLower(rest),"values");var tableName string;var colNames []string;if open>=0&&valuesAt>open{tableName=tsNodeSQLiteIdentifier(rest[:open]);close:=strings.Index(rest[open:],")");if close<0{return nil,nil,0,fmt.Errorf("malformed INSERT")};for _,col:=range tsNodeSQLiteSplit(rest[open+1:open+close],','){colNames=append(colNames,tsNodeSQLiteIdentifier(col))}}else{fields:=strings.Fields(rest);if len(fields)==0{return nil,nil,0,fmt.Errorf("malformed INSERT")};tableName=tsNodeSQLiteIdentifier(fields[0])};table:=db.tables[strings.ToLower(tableName)];if table==nil{return nil,nil,0,fmt.Errorf("no such table: %s",tableName)};valuesAt=strings.Index(strings.ToLower(rest),"values");if valuesAt<0{return nil,nil,0,fmt.Errorf("malformed INSERT")};openValues:=strings.Index(rest[valuesAt:],"(");if openValues<0{return nil,nil,0,fmt.Errorf("malformed INSERT")};openValues+=valuesAt;closeValues:=strings.Index(rest[openValues:],")");if closeValues<0{return nil,nil,0,fmt.Errorf("malformed INSERT")};closeValues+=openValues;literalParts:=tsNodeSQLiteSplit(rest[openValues+1:closeValues],',');values:=make([]tsValue,len(table.columns));for i:=range values{values[i]=tsNull};if len(colNames)==0{colNames=append(colNames,table.columns...)};for i,name:=range colNames{for j,col:=range table.columns{if strings.EqualFold(name,col)&&i<len(literalParts){values[j]=tsNodeSQLiteLiteral(loop,literalParts[i],args,&paramIndex)}}};table.rows=append(table.rows,values);table.nextRowID++;db.totalChanges++;db.lastInsertRowID=table.nextRowID-1;return nil,nil,1,nil}
 if strings.HasPrefix(lower,"delete from"){fields:=strings.Fields(source);if len(fields)<3{return nil,nil,0,fmt.Errorf("malformed DELETE")};tableName:=tsNodeSQLiteIdentifier(fields[2]);table:=db.tables[strings.ToLower(tableName)];if table==nil{return nil,nil,0,fmt.Errorf("no such table: %s",tableName)};count:=int64(len(table.rows));table.rows=nil;db.totalChanges+=count;return nil,nil,count,nil}
 if strings.HasPrefix(lower,"drop table"){fields:=strings.Fields(source);if len(fields)>2{delete(db.tables,strings.ToLower(tsNodeSQLiteIdentifier(fields[len(fields)-1])))};return nil,nil,0,nil}
 if strings.HasPrefix(lower,"select "){fromAt:=strings.Index(lower," from ");selectExpr:=strings.TrimSpace(source[6:]);var rows [][]tsValue;var columns,selected []string;if fromAt<0{items:=tsNodeSQLiteSplit(selectExpr,',');row:=[]tsValue{};for _,item:=range items{name:=item;if as:=strings.Index(strings.ToLower(item)," as ");as>=0{name=strings.TrimSpace(item[as+4:]);item=item[:as]};value:=tsNodeSQLiteLiteral(loop,item,args,&paramIndex);row=append(row,value);columns=append(columns,tsNodeSQLiteIdentifier(name))};return [][]tsValue{row},columns,0,nil};tableText:=strings.TrimSpace(source[fromAt+6:]);tableFields:=strings.Fields(tableText);if len(tableFields)==0{return nil,nil,0,fmt.Errorf("missing table name")};tableName:=tsNodeSQLiteIdentifier(tableFields[0]);table:=db.tables[strings.ToLower(tableName)];if table==nil{return nil,nil,0,fmt.Errorf("no such table: %s",tableName)};expr:=strings.TrimSpace(source[6:fromAt]);where:=strings.Index(strings.ToLower(tableText)," where ");whereExpr:="";if where>=0{whereExpr=strings.TrimSpace(tableText[where+7:]);if order:=strings.Index(strings.ToLower(whereExpr)," order by ");order>=0{whereExpr=whereExpr[:order]}};if strings.EqualFold(expr,"count(*)")||strings.EqualFold(expr,"count(1)"){count:=0;for _,row:=range table.rows{if tsNodeSQLiteMatches(loop,table,row,whereExpr,args,&paramIndex){count++}};return [][]tsValue{{tsNumberValue(float64(count))}},[]string{"count(*)"},0,nil};if expr=="*"{columns=append(columns,table.columns...);selected=append(selected,table.columns...)}else{for _,item:=range tsNodeSQLiteSplit(expr,','){item=strings.TrimSpace(item);name:=tsNodeSQLiteIdentifier(item);label:=name;if as:=strings.Index(strings.ToLower(item)," as ");as>=0{label=tsNodeSQLiteIdentifier(item[as+4:]);name=tsNodeSQLiteIdentifier(item[:as])};columns=append(columns,label);selected=append(selected,name)}};for _,sourceRow:=range table.rows{if !tsNodeSQLiteMatches(loop,table,sourceRow,whereExpr,args,&paramIndex){continue};row:=[]tsValue{};for _,name:=range selected{found:=false;for i,col:=range table.columns{if strings.EqualFold(name,col)&&i<len(sourceRow){row=append(row,sourceRow[i]);found=true;break}};if !found{row=append(row,tsU)}};rows=append(rows,row)};return rows,columns,0,nil}
 if strings.HasPrefix(lower,"pragma "){return nil,nil,0,nil};return nil,nil,0,fmt.Errorf("near \"%s\": syntax error",source)}
func tsNodeSQLiteMatches(loop *tsLoop,table *tsNodeSQLiteTable,row []tsValue,where string,args []tsValue,index *int)bool{if where==""{return true};parts:=strings.SplitN(where,"=",2);if len(parts)!=2{return false};name:=tsNodeSQLiteIdentifier(parts[0]);column:=-1;for i,c:=range table.columns{if strings.EqualFold(name,c){column=i;break}};if column<0||column>=len(row){return false};cursor:=0;expected:=tsNodeSQLiteLiteral(loop,parts[1],args,&cursor);_ = index;actual:=row[column];if actual.kind==tsNumberKind&&expected.kind==tsNumberKind{return tsNumber(actual)==tsNumber(expected)};return tsText(actual)==tsText(expected)}
func tsNodeSQLiteExecute(loop *tsLoop,statement tsValue,args []tsValue,mode string)tsValue{state:=tsNodeSQLiteStatementFor(statement);if state==nil{tsPropertyFailure("Invalid StatementSync receiver")};db:=tsNodeSQLiteRequireDB(state.db);if state.finalized{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement has been finalized")})};if state.stepping{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","statement is currently executing")})};state.stepping=true;defer func(){state.stepping=false}();bindings:=tsNodeSQLiteBindings(loop,args);rows,columns,changes,err:=tsNodeSQLiteResult(loop,db,state.sql,bindings);if err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};state.columns=columns;if mode=="run"{result:=tsNewObject();result.set("changes",tsNumberValue(float64(changes)));result.set("lastInsertRowid",tsNumberValue(float64(db.lastInsertRowID)));return tsObjectValue(result)};out:=&tsArray{};for _,row:=range rows{out.values=append(out.values,tsNodeSQLiteRowObject(loop,columns,row,state.returnArrays))};if mode=="get"{if len(out.values)==0{return tsU};return out.values[0]};if mode=="iterate"{return tsIteratorValueReference(&tsIterator{array:out})};return tsArrayValue(out)}
func tsNodeSQLiteStatement(loop *tsLoop,dbValue tsValue,sqlValue tsValue)tsValue{
	db:=tsNodeSQLiteRequireDB(dbValue)
	if sqlValue.kind!=tsStringKind{tsFSInvalid("sql","a string")}
	sql:=tsText(sqlValue);if err:=tsNodeSQLiteAuthorizeSQL(loop,db,sql);err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};db=tsNodeSQLiteRequireDB(dbValue);state:=&tsNodeSQLiteStatement{db:dbValue,sql:sql,allowBare:db.allowBare,allowUnknown:db.allowUnknown,returnArrays:db.returnArrays,safeIntegers:db.readBigInts,params:strings.Count(sql,"?")}
	obj:=tsNewObject();obj.prototype=tsObjectValue(tsNodeSQLiteSingleton.stmtProto);state.object=obj
	for _,property:=range []string{"sourceSQL","expandedSQL"}{name:=property;obj.order=append(obj.order,name);tsNodeSQLiteAccessor(obj,name,func(_ *tsLoop,r tsValue)tsValue{statement:=tsNodeSQLiteStatementFor(r);if statement==nil{return tsU};return tsNodeSQLiteKey(statement.sql)},nil);obj.descriptors[name].enumerable=true;obj.descriptors[name].configurable=false}
	tsNodeSQLiteSingleton.statements[obj]=state;db.statements[obj]=true
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
func tsNodeSQLiteDatabase(loop *tsLoop,args []tsValue)tsValue{
	location:=tsNodeSQLitePath(loop,tsArg(args,0))
	options:=tsObjectValue(tsNewObject())
	if len(args)>1{options=tsArg(args,1);if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}}
	db:=&tsNodeSQLiteDB{location:location,opened:true,tables:map[string]*tsNodeSQLiteTable{},functions:map[string]tsValue{},aggregates:map[string]tsValue{},authorizer:tsNull,limits:tsNodeSQLiteDefaultLimits(),limitMax:tsNodeSQLiteDefaultLimits(),statements:map[*tsObject]bool{},defensive:true,allowBare:true}
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
	if !tsIsUndefined(timeout){if !tsIsNumeric(timeout)||math.IsNaN(tsNumber(timeout))||math.IsInf(tsNumber(timeout),0)||math.Trunc(tsNumber(timeout))!=tsNumber(timeout)||tsNumber(timeout)<math.MinInt32||tsNumber(timeout)>math.MaxInt32{tsFSInvalid("options.timeout","an integer")};db.timeout=int(tsNumber(timeout))}
	obj:=tsNewObject();obj.prototype=tsObjectValue(tsNodeSQLiteSingleton.dbProto);db.object=obj
	for _,property:=range []string{"isOpen","isTransaction","limits"}{name:=property;obj.order=append(obj.order,name);tsNodeSQLiteAccessor(obj,name,func(loop *tsLoop,r tsValue)tsValue{current:=tsNodeSQLiteDBFor(r);if current==nil{return tsU};switch name{case "isOpen":return tsBooleanValue(current.opened);case "isTransaction":if !current.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};return tsBooleanValue(current.transaction);case "limits":return tsNodeSQLiteLimits(loop,r)};return tsU},nil);obj.descriptors[name].enumerable=true;obj.descriptors[name].configurable=false}
	tsNodeSQLiteSingleton.dbs[obj]=db;return tsObjectValue(obj)
}
func tsNodeSQLiteDatabaseMethod(loop *tsLoop,receiver tsValue,name string,args []tsValue)tsValue{db:=tsNodeSQLiteDBFor(receiver);if db==nil{tsPropertyFailure("DatabaseSync method called on an incompatible receiver")};if name!="open"&&!db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};switch name{case "open":if db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is already open")})};db.opened=true;return tsU;case "close":if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};db.opened=false;db.transaction=false;for object,state:=range tsNodeSQLiteSingleton.sessions{if state!=nil&&tsNodeSQLiteDBFor(state.db)==db{state.closed=true};_ = object};for object:=range db.statements{if state:=tsNodeSQLiteSingleton.statements[object];state!=nil{state.finalized=true}};return tsU;case "exec":if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};sql:=tsArg(args,0);if sql.kind!=tsStringKind{tsFSInvalid("sql","a string")};for _,part:=range tsNodeSQLiteSplit(tsText(sql),';'){if strings.TrimSpace(part)==""{continue};if err:=tsNodeSQLiteAuthorizeSQL(loop,db,part);err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})};if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};if _,_,_,err:=tsNodeSQLiteResult(loop,db,part,nil);err!=nil{panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR",err.Error())})}};return tsU;case "prepare":return tsNodeSQLiteStatement(loop,receiver,tsArg(args,0));case "location":if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})};return tsNodeSQLiteKey(db.location);case "enableDefensive":active:=tsArg(args,0);if active.kind!=tsBooleanKind{tsFSInvalid("active","a boolean")};db.defensive=tsTruthy(active);return tsU;case "enableLoadExtension":allow:=tsArg(args,0);if allow.kind!=tsBooleanKind{tsFSInvalid("allow","a boolean")};enabled:=tsTruthy(allow);if enabled&&!db.allowExtension{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","Cannot enable extension loading because it was disabled at database creation.")})};db.extensions=enabled;return tsU;case "loadExtension":if !db.extensions{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","extension loading is not allowed")})};if tsArg(args,0).kind!=tsStringKind{tsFSInvalid("path","a string")};if !tsIsUndefined(tsArg(args,1))&&tsArg(args,1).kind!=tsStringKind{tsFSInvalid("entryPoint","a string")};panic(tsThrown{tsNodeSQLiteError("ERR_LOAD_SQLITE_EXTENSION","Extension loading is not available in the Go runtime")});case "function","aggregate":if len(args)<2||tsArg(args,1).kind!=tsFunctionKind{tsFSInvalid("function","a function")};nameValue:=tsText(tsArg(args,0));if name=="function"{db.functions[strings.ToLower(nameValue)]=tsArg(args,1)}else{db.aggregates[strings.ToLower(nameValue)]=tsArg(args,1)};return tsU;case "createSession":if len(args)>0{options:=args[0];if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};for _,name:=range []string{"table","db"}{value:=tsNodeSQLiteGet(loop,options,name);if !tsIsUndefined(value)&&value.kind!=tsStringKind{tsFSInvalid("options."+name,"a string")}};db=tsNodeSQLiteRequireDB(receiver)};session:=tsNewObject();session.prototype=tsObjectValue(tsNodeSQLiteSingleton.sessionProto);state:=&tsNodeSQLiteSession{object:session,db:receiver,changeOffset:int(db.totalChanges)};tsNodeSQLiteSingleton.sessions[session]=state;return tsObjectValue(session);case "applyChangeset":changeset:=tsArg(args,0);if changeset.kind!=tsTypedArrayKind{tsFSInvalid("changeset","a Uint8Array")};view:=(*tsTypedArray)(changeset.ref);if view.name!="Uint8Array"&&!view.nodeBuffer{tsFSInvalid("changeset","a Uint8Array")};options:=tsArg(args,1);if !tsIsUndefined(options){if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")};for _,name:=range []string{"onConflict","filter"}{callback:=tsNodeSQLiteGet(loop,options,name);if !tsIsUndefined(callback)&&callback.kind!=tsFunctionKind{tsFSInvalid("options."+name,"a function")}}};if view.buffer.detached{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_VALUE","The \"changeset\" argument must not be detached.")})};return tsBooleanValue(true);case "setAuthorizer":callback:=tsArg(args,0);if callback.kind!=tsNullKind&&callback.kind!=tsFunctionKind{tsFSInvalid("callback","a function or null")};db.authorizer=callback;return tsU;case "serialize":return tsNodeBuffer([]byte("SQLite format 3\x00"));case "deserialize":data:=tsArg(args,0);if data.kind!=tsTypedArrayKind{tsFSInvalid("buffer","a Uint8Array")};view:=(*tsTypedArray)(data.ref);if view.name!="Uint8Array"&&!view.nodeBuffer{tsFSInvalid("buffer","a Uint8Array")};if view.buffer.detached||view.length==0{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_ARG_VALUE","The \"buffer\" argument must not be empty.")})};for object,state:=range tsNodeSQLiteSingleton.sessions{if state!=nil&&tsNodeSQLiteDBFor(state.db)==db{state.closed=true};_ = object};for object:=range db.statements{if statement:=tsNodeSQLiteSingleton.statements[object];statement!=nil{statement.finalized=true}};db.tables=map[string]*tsNodeSQLiteTable{};return tsU;case "createTagStore":capacity:=uint64(1000);if value:=tsArg(args,0);tsIsNumeric(value){capacity=uint64(int64(int32(tsToUint32(value))))};obj:=tsNewObject();obj.prototype=tsObjectValue(tsNodeSQLiteSingleton.tagProto);state:=&tsNodeSQLiteTagStore{object:obj,db:receiver,capacity:capacity,cache:map[string]tsValue{}}
	for _,property:=range []string{"capacity","size","db"}{name:=property;obj.order=append(obj.order,name);tsNodeSQLiteAccessor(obj,name,func(_ *tsLoop,_ tsValue)tsValue{switch name{case "capacity":return tsNumberValue(float64(state.capacity));case "size":return tsNumberValue(float64(len(state.cache)));case "db":return state.db};return tsU},nil);obj.descriptors[name].enumerable=true;obj.descriptors[name].configurable=false}
	tsNodeSQLiteSingleton.tags[obj]=state;return tsObjectValue(obj);case "isInTransaction":return tsBooleanValue(db.transaction);case "setCustomSQLite":panic(tsThrown{tsNodeSQLiteError("ERR_SQLITE_ERROR","Custom SQLite engines are not supported")});case "fcntl":return tsNumberValue(0)};return tsU}
func tsNodeSQLiteStatementMethod(loop *tsLoop,r tsValue,name string,args []tsValue)tsValue{state:=tsNodeSQLiteStatementFor(r);if state==nil{tsPropertyFailure("StatementSync method called on an incompatible receiver")};switch name{case "run","get","all","iterate":return tsNodeSQLiteExecute(loop,r,args,name);case "columns":out:=&tsArray{};for i,n:=range state.columns{item:=tsNewObject();item.set("column",tsNodeSQLiteKey(n));item.set("columnIndex",tsNumberValue(float64(i)));item.set("database",tsNull);item.set("table",tsNull);item.set("name",tsNodeSQLiteKey(n));item.set("type",tsNodeSQLiteKey("TEXT"));item.set("nullable",tsBooleanValue(true));out.values=append(out.values,tsObjectValue(item))};return tsArrayValue(out);case "setReadBigInts":if len(args)==0||tsArg(args,0).kind!=tsBooleanKind{tsFSInvalid("enabled","a boolean")};state.safeIntegers=tsTruthy(args[0]);return tsU;case "setReturnArrays":if len(args)==0||tsArg(args,0).kind!=tsBooleanKind{tsFSInvalid("enabled","a boolean")};state.returnArrays=tsTruthy(args[0]);return tsU;case "setAllowBareNamedParameters":state.allowBare=tsTruthy(tsArg(args,0));return tsU;case "setAllowUnknownNamedParameters":state.allowUnknown=tsTruthy(tsArg(args,0));return tsU;case "finalize":state.finalized=true;return tsU;case "toString":return tsNodeSQLiteKey(state.sql)};return tsU}
func tsNodeSQLiteSessionMethod(loop *tsLoop,r tsValue,name string)tsValue{
	state:=tsNodeSQLiteSessionFor(r)
	if state==nil{tsPropertyFailure("Session method called on an incompatible receiver")}
	db:=tsNodeSQLiteDBFor(state.db)
	if name=="dispose"{if db==nil||!db.opened||state.closed{return tsU};state.closed=true;return tsU}
	if db==nil||!db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})}
	if state.closed{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","session is not open")})}
	switch name{case "changeset","patchset":return tsNodeBuffer(nil);case "close":state.closed=true;return tsU}
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
			stmt=tsNodeSQLiteStatement(loop,state.db,tsNodeSQLiteKey(text))
			if state.capacity>0{state.cache[text]=stmt;state.order=append(state.order,text);if uint64(len(state.order))>state.capacity{oldest:=state.order[0];delete(state.cache,oldest);state.order=state.order[1:]}}
		}else{
			for i,key:=range state.order{if key==text{state.order=append(state.order[:i],state.order[i+1:]...);break}}
			state.order=append(state.order,text)
		}
		return tsNodeSQLiteExecute(loop,stmt,args[1:],name)
	}
	return tsU
}
func tsNodeSQLiteBackup(loop *tsLoop,args []tsValue)tsValue{
	source:=tsArg(args,0)
	if tsNodeSQLiteDBFor(source)==nil{tsFSInvalid("sourceDb","a DatabaseSync")}
	db:=tsNodeSQLiteRequireDB(source)
	path:=tsArg(args,1)
	_ = tsNodeSQLitePath(loop,path)
	if len(args)>2{
		options:=args[2]
		if options.kind!=tsObjectKind&&options.kind!=tsInstanceKind{tsFSInvalid("options","an object")}
		rate:=tsNodeSQLiteGet(loop,options,"rate")
		if !tsIsUndefined(rate){if !tsIsNumeric(rate)||math.IsNaN(tsNumber(rate))||math.IsInf(tsNumber(rate),0)||math.Trunc(tsNumber(rate))!=tsNumber(rate)||tsNumber(rate)<math.MinInt32||tsNumber(rate)>math.MaxInt32{tsFSInvalid("options.rate","an integer")}}
		for _,name:=range []string{"source","target"}{value:=tsNodeSQLiteGet(loop,options,name);if !tsIsUndefined(value)&&value.kind!=tsStringKind{tsFSInvalid("options."+name,"a string")}}
		progress:=tsNodeSQLiteGet(loop,options,"progress");if !tsIsUndefined(progress)&&progress.kind!=tsFunctionKind{tsFSInvalid("options.progress","a function")}
	}
	if !db.opened{panic(tsThrown{tsNodeSQLiteError("ERR_INVALID_STATE","database is not open")})}
	// The adapter has no SQLite file engine, so successful backups currently
	// resolve with zero copied pages after matching the native validation path.
	return tsPromiseValue(loop.resolved(tsNumberValue(0),false))
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
 tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.stmtProto,"sourceSQL",func(_ *tsLoop,r tsValue)tsValue{if s:=tsNodeSQLiteStatementFor(r);s!=nil{return tsNodeSQLiteKey(s.sql)};return tsU},nil);tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.stmtProto,"expandedSQL",func(_ *tsLoop,r tsValue)tsValue{if s:=tsNodeSQLiteStatementFor(r);s!=nil{return tsNodeSQLiteKey(s.sql)};return tsU},nil);for _,property:=range []string{"columnNames","columnCount","parameterCount","isFinalized","columnTypes","declaredColumnTypes","safeIntegers"}{n:=property;tsNodeSQLiteAccessor(tsNodeSQLiteSingleton.stmtProto,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeSQLiteStatementFor(r);if s==nil{return tsU};switch n{case "columnNames":a:=&tsArray{};for _,v:=range s.columns{a.values=append(a.values,tsNodeSQLiteKey(v))};return tsArrayValue(a);case "columnCount":return tsNumberValue(float64(len(s.columns)));case "parameterCount":return tsNumberValue(float64(s.params));case "isFinalized":return tsBooleanValue(s.finalized);case "columnTypes","declaredColumnTypes":a:=&tsArray{};for _,v:=range s.declared{a.values=append(a.values,tsNodeSQLiteKey(v))};return tsArrayValue(a);case "safeIntegers":return tsBooleanValue(s.safeIntegers)};return tsU},func(_ *tsLoop,r,v tsValue)tsValue{if n=="safeIntegers"{if v.kind!=tsBooleanKind{tsFSInvalid("safeIntegers","a boolean")};if s:=tsNodeSQLiteStatementFor(r);s!=nil{s.safeIntegers=tsTruthy(v)}};return v})}
 tag:=tsNodeSQLiteSingleton.tagProto;for _,property:=range []string{"capacity","size","db"}{n:=property;tsNodeSQLiteAccessor(tag,n,func(_ *tsLoop,r tsValue)tsValue{s:=tsNodeSQLiteTagFor(r);if s==nil{return tsU};switch n{case "capacity":return tsNumberValue(float64(s.capacity));case "size":return tsNumberValue(float64(len(s.cache)));case "db":return s.db};return tsU},nil)}
 constants:=tsNewObject();for name,value:=range map[string]int{"SQLITE_CHANGESET_OMIT":0,"SQLITE_CHANGESET_REPLACE":1,"SQLITE_CHANGESET_ABORT":2,"SQLITE_CHANGESET_DATA":1,"SQLITE_CHANGESET_NOTFOUND":2,"SQLITE_CHANGESET_CONFLICT":3,"SQLITE_CHANGESET_CONSTRAINT":4,"SQLITE_CHANGESET_FOREIGN_KEY":5,"SQLITE_OK":0,"SQLITE_DENY":1,"SQLITE_IGNORE":2,"SQLITE_CREATE_INDEX":1,"SQLITE_CREATE_TABLE":2,"SQLITE_CREATE_TEMP_INDEX":3,"SQLITE_CREATE_TEMP_TABLE":4,"SQLITE_CREATE_TEMP_TRIGGER":5,"SQLITE_CREATE_TEMP_VIEW":6,"SQLITE_CREATE_TRIGGER":7,"SQLITE_CREATE_VIEW":8,"SQLITE_DELETE":9,"SQLITE_DROP_INDEX":10,"SQLITE_DROP_TABLE":11,"SQLITE_DROP_TEMP_INDEX":12,"SQLITE_DROP_TEMP_TABLE":13,"SQLITE_DROP_TEMP_TRIGGER":14,"SQLITE_DROP_TEMP_VIEW":15,"SQLITE_DROP_TRIGGER":16,"SQLITE_DROP_VIEW":17,"SQLITE_INSERT":18,"SQLITE_PRAGMA":19,"SQLITE_READ":20,"SQLITE_SELECT":21,"SQLITE_TRANSACTION":22,"SQLITE_UPDATE":23,"SQLITE_ATTACH":24,"SQLITE_DETACH":25,"SQLITE_ALTER_TABLE":26,"SQLITE_REINDEX":27,"SQLITE_ANALYZE":28,"SQLITE_CREATE_VTABLE":29,"SQLITE_DROP_VTABLE":30,"SQLITE_FUNCTION":31,"SQLITE_SAVEPOINT":32,"SQLITE_COPY":0,"SQLITE_RECURSIVE":33}{constants.set(name,tsNumberValue(float64(value)))};tsNodeSQLiteSingleton.constants=tsObjectValue(constants)
 module:=tsNewObject();module.set("DatabaseSync",database);module.set("StatementSync",statement);module.set("Session",session);module.set("constants",tsNodeSQLiteSingleton.constants);module.set("backup",tsNamedNativeMethod("backup",2,func(loop *tsLoop,_ tsValue,args ...tsValue)tsValue{return tsNodeSQLiteBackup(loop,args)}));module.set("default",tsObjectValue(module));tsNodeSQLiteSingleton.module=tsObjectValue(module)
})}
func tsNodeSQLiteModule(loop *tsLoop)tsValue{tsNodeSQLiteInit(loop);return tsNodeSQLiteSingleton.module}
func tsNodeSQLiteAddon(loop *tsLoop,name string)(tsValue,bool){if name=="sqlite"{return tsNodeSQLiteModule(loop),true};return tsU,false}
`

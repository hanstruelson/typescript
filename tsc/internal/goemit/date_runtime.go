package goemit

// Date values used by filesystem metadata store clipped milliseconds natively.
// No per-instance method closures or Go interfaces are used for these values.
const dateRuntime = dateHelpersRuntime + dateLocaleRuntime + `
func tsDateMilliseconds(value tsValue)float64{
    if value.kind==tsObjectKind{object:=(*tsObject)(value.ref);if object.nativeDate{return object.nativeTime}}
    value=tsToPrimitive(value)
    if value.kind==tsStringKind{
        return tsDateParse(tsText(value))
    };return tsNumber(value)
}
func tsDateClip(milliseconds float64)float64{if math.IsNaN(milliseconds)||math.IsInf(milliseconds,0)||math.Abs(milliseconds)>8640000000000000{return math.NaN()};milliseconds=math.Trunc(milliseconds);if milliseconds==0{return 0};return milliseconds}
func tsFSDateMillis(milliseconds float64)tsValue{out:=tsNewObject();out.nativeClass=tsNativeClass("Date");out.nativeDate=true;out.nativeTime=tsDateClip(milliseconds);return tsObjectValue(out)}
func tsFSDate(stamp time.Time)tsValue{return tsFSDateMillis(float64(stamp.Unix())*1000+float64(stamp.Nanosecond())/1e6)}
func tsDateInstant(milliseconds float64)time.Time{seconds:=math.Floor(milliseconds/1000);remainder:=milliseconds-seconds*1000;return time.Unix(int64(seconds),int64(remainder)*1000000)}
func tsDateConstruct(args []tsValue)tsValue{
    if len(args)==0{return tsFSDate(time.Now())};if len(args)==1{return tsFSDateMillis(tsDateMilliseconds(args[0]))}
    return tsFSDateMillis(tsDateFields(loop,args,false))
}
func tsInstallDateClass(class *tsClass){
    class.nativePrototype.nativeTime=math.NaN();class.nativePrototype.nativeDate=true
    for _,name:=range []string{"getTime","valueOf","toISOString","toJSON","toString","toUTCString","toDateString","toTimeString","getFullYear","getMonth","getDate","getDay","getHours","getMinutes","getSeconds","getMilliseconds","getTimezoneOffset","getUTCFullYear","getUTCMonth","getUTCDate","getUTCDay","getUTCHours","getUTCMinutes","getUTCSeconds","getUTCMilliseconds","setTime","toLocaleString","toLocaleDateString","toLocaleTimeString","getYear","setYear","toGMTString","setDate","setMonth","setFullYear","setHours","setMinutes","setSeconds","setMilliseconds","setUTCDate","setUTCMonth","setUTCFullYear","setUTCHours","setUTCMinutes","setUTCSeconds","setUTCMilliseconds"}{
        method:=name
        arity:=0;if strings.HasPrefix(method,"set"){arity=1;switch method{case "setFullYear","setUTCFullYear","setMinutes","setUTCMinutes":arity=3;case "setHours","setUTCHours":arity=4;case "setMonth","setUTCMonth","setSeconds","setUTCSeconds":arity=2}};if method=="toJSON"{arity=1}
        class.nativePrototype.set(method,tsNamedNativeMethod(method,arity,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{
            if method=="toJSON"{primitive:=tsToPrimitive(loop,receiver);if tsIsNumeric(primitive)&&(math.IsNaN(tsNumber(primitive))||math.IsInf(tsNumber(primitive),0)){return tsNull};iso:=tsGet(loop,receiver,tsStringReference(tsStringUTF8("toISOString")));if iso.kind!=tsFunctionKind{tsPropertyFailure("toISOString is not callable")};return tsCallReceiver(loop,iso,receiver)}
            if receiver.kind!=tsObjectKind||!(*tsObject)(receiver.ref).nativeDate{tsFSInvalid("this","a Date")};object:=(*tsObject)(receiver.ref);milliseconds:=object.nativeTime
            switch method{case "getTime","valueOf":return tsNumberValue(milliseconds);case "setTime":object.nativeTime=tsDateClip(tsNumber(tsToPrimitive(tsArg(args,0))));return tsNumberValue(object.nativeTime)}
            if strings.HasPrefix(method,"set"){
                utc:=strings.HasPrefix(method,"setUTC");name:=strings.TrimPrefix(strings.TrimPrefix(method,"set"),"UTC");location:=tsDateLocal();if utc{location=time.UTC}
                invalid:=math.IsNaN(milliseconds);if invalid{milliseconds=0}
                stamp:=tsDateInstant(milliseconds).In(location);if math.IsNaN(object.nativeTime){stamp=time.Date(1970,time.January,1,0,0,0,0,location)};fields:=[]float64{float64(stamp.Year()),float64(stamp.Month()-1),float64(stamp.Day()),float64(stamp.Hour()),float64(stamp.Minute()),float64(stamp.Second()),float64(stamp.Nanosecond()/1000000)}
                index:=0;switch name{case "Month":index=1;case "Date":index=2;case "Hours":index=3;case "Minutes":index=4;case "Seconds":index=5;case "Milliseconds":index=6}
                count:=1;switch name{case "FullYear":count=3;case "Month":count=2;case "Hours":count=4;case "Minutes":count=3;case "Seconds":count=2}
                if len(args)==0{object.nativeTime=math.NaN();return tsNumberValue(object.nativeTime)}
                for i:=0;i<len(args)&&i<count;i++{fields[index+i]=math.Trunc(tsNumber(tsToPrimitive(loop,args[i])))}
                if invalid&&name!="FullYear"&&name!="Year"{object.nativeTime=math.NaN();return tsNumberValue(object.nativeTime)}
                if name=="Year"&&fields[0]>=0&&fields[0]<=99{fields[0]+=1900}
                object.nativeTime=tsDateMake(fields,location);return tsNumberValue(object.nativeTime)
            }
            if math.IsNaN(milliseconds){if method=="toISOString"{panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:"Invalid time value"})})};if method=="toJSON"{return tsNull};if strings.HasPrefix(method,"to"){return tsStringReference(tsStringUTF8("Invalid Date"))};return tsNumberValue(math.NaN())}
            stamp:=tsDateInstant(milliseconds);if !strings.Contains(method,"UTC")&&method!="toISOString"&&method!="toJSON"&&method!="toGMTString"{stamp=stamp.In(tsDateLocal())}else{stamp=stamp.UTC()}
            if method=="toLocaleString"||method=="toLocaleDateString"||method=="toLocaleTimeString"{return tsDateLocale(loop,stamp,method,args)}
            switch method{
            case "toISOString","toJSON":year:=stamp.Year();prefix:=fmt.Sprintf("%04d",year);if year<0{prefix=fmt.Sprintf("-%06d",-year)}else if year>9999{prefix=fmt.Sprintf("+%06d",year)};text:=prefix+stamp.Format("-01-02T15:04:05.000Z");return tsStringReference(tsStringUTF8(text))
            case "toUTCString","toGMTString":utc:=stamp.UTC();return tsStringReference(tsStringUTF8(utc.Format("Mon, 02 Jan ")+tsDateYear(utc.Year(),4)+utc.Format(" 15:04:05 GMT")))
            case "toString":return tsStringReference(tsStringUTF8(tsDateHuman(stamp,method)))
            case "toDateString":return tsStringReference(tsStringUTF8(tsDateHuman(stamp,method)))
            case "toTimeString":return tsStringReference(tsStringUTF8(tsDateHuman(stamp,method)))
            case "getYear":return tsNumberValue(float64(stamp.Year()-1900))
            case "getFullYear","getUTCFullYear":return tsNumberValue(float64(stamp.Year()))
            case "getMonth","getUTCMonth":return tsNumberValue(float64(stamp.Month()-1))
            case "getDate","getUTCDate":return tsNumberValue(float64(stamp.Day()))
            case "getDay","getUTCDay":return tsNumberValue(float64(stamp.Weekday()))
            case "getHours","getUTCHours":return tsNumberValue(float64(stamp.Hour()))
            case "getMinutes","getUTCMinutes":return tsNumberValue(float64(stamp.Minute()))
            case "getSeconds","getUTCSeconds":return tsNumberValue(float64(stamp.Second()))
            case "getMilliseconds","getUTCMilliseconds":return tsNumberValue(float64(stamp.Nanosecond()/1000000))
            case "getTimezoneOffset":_,offset:=stamp.Zone();return tsNumberValue(float64(-offset/60))
            };return tsU
        }))
    }
    class.nativePrototype.set("toGMTString",class.nativePrototype.values["toUTCString"])
    properties:=tsInstanceProperties(class.static)
    properties.define("length");properties.extra["length"]=tsNumberValue(7)
    for _,name:=range []string{"now","parse","UTC"}{operation:=name;arity:=0;if name=="parse"{arity=1};if name=="UTC"{arity=7};properties.define(name);properties.extra[name]=tsNamedNativeMethod(name,arity,func(loop *tsLoop,receiver tsValue,args ...tsValue)tsValue{switch operation{case "now":return tsNumberValue(float64(time.Now().UnixMilli()));case "parse":return tsNumberValue(tsDateClip(tsDateMilliseconds(tsStringReference(tsStringValue(tsArg(args,0))))));default:return tsNumberValue(tsDateFields(loop,args,true))}})}
}
`

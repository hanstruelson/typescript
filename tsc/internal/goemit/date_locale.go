package goemit

// Locale patterns are development-generated data, not an engine boundary.
// The initial data set supports English Gregorian formatting. More CLDR locales,
// calendars and numbering systems must be added before full Intl compatibility.
const dateLocaleRuntime = dateLocaleDataRuntime + `
var tsDateLocalePatterns struct{sync.Mutex;tables map[string]map[string]string}
func tsDateLocalePattern(locale,key string)string{tsDateLocalePatterns.Lock();defer tsDateLocalePatterns.Unlock();if tsDateLocalePatterns.tables==nil{tsDateLocalePatterns.tables=map[string]map[string]string{}};table:=tsDateLocalePatterns.tables[locale];if table==nil{table=map[string]string{};for _,entry:=range strings.Split(tsDateLocaleData[locale],"\n"){parts:=strings.SplitN(entry,"\t",2);if len(parts)==2{table[parts[0]]=parts[1]}};tsDateLocalePatterns.tables[locale]=table};return table[key]}
func tsDateLocaleRange(message string){panic(tsThrown{tsErrorValue(&tsRuntimeError{name:"RangeError",message:message})})}
func tsDateLocaleOption(loop *tsLoop,options tsValue,key string,values []string)int{value:=tsGet(loop,options,tsStringReference(tsStringUTF8(key)));if tsIsUndefined(value){return 0};text:=tsStringValue(loop,value).String();for i,name:=range values{if text==name{return i+1}};tsDateLocaleRange("Invalid "+key+" option");return 0}
func tsDateLocaleTag(loop *tsLoop,input tsValue)string{
 locale:="en-US";selected:=false;values:=[]tsValue{};if input.kind==tsStringKind{values=append(values,input)}else if input.kind==tsArrayKind{array:=(*tsArray)(input.ref);for i:=0;i<array.length();i++{if !array.holes[i]{values=append(values,array.at(i))}}}else if input.kind==tsNullKind{tsPropertyFailure("Cannot convert null locale list to an object")}
 for _,value:=range values{if value.kind!=tsStringKind&&!(!tsIsPrimitiveValue(value)){tsPropertyFailure("Locale list elements must be strings or objects")};text:=tsStringValue(loop,value).String();if strings.Contains(text,"_"){tsDateLocaleRange("Invalid language tag")};tag,err:=language.Parse(text);if err!=nil{tsDateLocaleRange("Invalid language tag")};canonical:=tag.String();if !selected&&(canonical=="en"||strings.HasPrefix(canonical,"en-")){selected=true;locale="en-US";if strings.HasPrefix(canonical,"en-GB")||strings.HasPrefix(canonical,"en-AU")||strings.HasPrefix(canonical,"en-NZ")||strings.HasPrefix(canonical,"en-IE"){locale="en-GB"}}};return locale
}
func tsDateLocaleZone(text string)*time.Location{if text==""||text=="Local"{tsDateLocaleRange("Invalid time zone")};if strings.EqualFold(text,"UTC")||strings.EqualFold(text,"Etc/UTC")||strings.EqualFold(text,"GMT"){return time.UTC};if len(text)==6&&(text[0]=='+'||text[0]=='-')&&text[3]==':'{hours,e1:=strconv.Atoi(text[1:3]);minutes,e2:=strconv.Atoi(text[4:]);if e1==nil&&e2==nil&&hours<24&&minutes<60{offset:=hours*3600+minutes*60;if text[0]=='-'{offset=-offset};return time.FixedZone(text,offset)}};zone,err:=time.LoadLocation(text);if err!=nil{tsDateLocaleRange("Invalid time zone")};return zone}
func tsDateLocaleZoneName(stamp time.Time,style string)string{
 name,offset:=stamp.Zone();if offset==0{if style=="long"{return "Coordinated Universal Time"};if style=="shortOffset"||style=="longOffset"||style=="shortGeneric"||style=="longGeneric"{return "GMT"};return "UTC"};switch style{case "long":human:=tsDateZone(stamp);start:=strings.IndexByte(human,'(');if start>=0{return human[start+1:len(human)-1]};case "shortGeneric":for _,entry:=range []string{"EST:ET","EDT:ET","PST:PT","PDT:PT","CST:CT","CDT:CT","MST:MT","MDT:MT"}{parts:=strings.SplitN(entry,":",2);if name==parts[0]{return parts[1]}};case "longGeneric":return strings.ReplaceAll(strings.ReplaceAll(tsDateLocaleZoneName(stamp,"long")," Standard", "")," Daylight", "");case "short":if name!=""&&!strings.HasPrefix(name,"+")&&!strings.HasPrefix(name,"-"){return name}}
 sign:="+";if offset<0{sign="-";offset=-offset};if style=="longOffset"{return fmt.Sprintf("GMT%s%02d:%02d",sign,offset/3600,offset/60%60)};if offset%3600==0{return fmt.Sprintf("GMT%s%d",sign,offset/3600)};return fmt.Sprintf("GMT%s%d:%02d",sign,offset/3600,offset/60%60)
}
func tsDateLocaleRender(pattern string,stamp time.Time,cycle,zoneStyle string)string{
 var out strings.Builder;for len(pattern)>0{start:=strings.IndexByte(pattern,'{');if start<0{out.WriteString(pattern);break};out.WriteString(pattern[:start]);end:=strings.IndexByte(pattern[start:],'}')+start;if end<start{panic("Invalid date locale pattern")};token:=pattern[start+1:end];pattern=pattern[end+1:];number:=0;width:=0;if strings.HasSuffix(token,"2")&&!strings.HasPrefix(token,"fractionalSecond"){width=2;token=strings.TrimSuffix(token,"2")};switch token{
 case "year":number=stamp.Year();if number<=0{number=1-number};if width==2{number%=100}
 case "month":number=int(stamp.Month())
 case "day":number=stamp.Day()
 case "hour":number=stamp.Hour();switch cycle{case "h11":number%=12;case "h12":number%=12;if number==0{number=12};case "h24":if number==0{number=24}}
 case "minute":number=stamp.Minute()
 case "second":number=stamp.Second()
 case "monthlong":out.WriteString(stamp.Month().String());continue
 case "monthshort":out.WriteString(stamp.Month().String()[:3]);continue
 case "monthnarrow":out.WriteString(stamp.Month().String()[:1]);continue
 case "weekdaylong":out.WriteString(stamp.Weekday().String());continue
 case "weekdayshort":out.WriteString(stamp.Weekday().String()[:3]);continue
 case "weekdaynarrow":out.WriteString(stamp.Weekday().String()[:1]);continue
 case "eralong","erashort","eranarrow":era:="AD";if stamp.Year()<=0{era="BC"};if token=="eralong"{era="Anno Domini";if stamp.Year()<=0{era="Before Christ"}}else if token=="eranarrow"{era=era[:1]};out.WriteString(era);continue
 case "dayPeriod":period:="AM";if stamp.Hour()>=12{period="PM"};out.WriteString(period);continue
 case "fractionalSecond1","fractionalSecond2","fractionalSecond3":digits:=int(token[len(token)-1]-'0');text:=fmt.Sprintf("%03d",stamp.Nanosecond()/1000000);out.WriteString(text[:digits]);continue
 case "timeZoneName":out.WriteString(tsDateLocaleZoneName(stamp,zoneStyle));continue
 default:panic("Unknown date locale token: "+token)
 };if width==2{fmt.Fprintf(&out,"%02d",number)}else{out.WriteString(strconv.Itoa(number))}}
 return out.String()
}
func tsDateLocale(loop *tsLoop,stamp time.Time,method string,args []tsValue)tsValue{
 locale:=tsDateLocaleTag(loop,tsArg(args,0));options:=tsArg(args,1);if tsIsUndefined(options){options=tsObjectValue(tsNewObject())};if tsNullish(options){tsPropertyFailure("Cannot convert null options to an object")}
 // Read and validate options through the same accessor-aware property boundary.
 dateValues:=[]int{tsDateLocaleOption(loop,options,"year",[]string{"numeric","2-digit"}),tsDateLocaleOption(loop,options,"month",[]string{"numeric","2-digit","long","short","narrow"}),tsDateLocaleOption(loop,options,"day",[]string{"numeric","2-digit"}),tsDateLocaleOption(loop,options,"weekday",[]string{"long","short","narrow"}),tsDateLocaleOption(loop,options,"era",[]string{"long","short","narrow"})}
 timeValues:=[]int{tsDateLocaleOption(loop,options,"hour",[]string{"numeric","2-digit"}),tsDateLocaleOption(loop,options,"minute",[]string{"numeric","2-digit"}),tsDateLocaleOption(loop,options,"second",[]string{"numeric","2-digit"}),0,tsDateLocaleOption(loop,options,"timeZoneName",[]string{"short","long","shortOffset","longOffset","shortGeneric","longGeneric"})}
 fraction:=tsGet(loop,options,tsStringReference(tsStringUTF8("fractionalSecondDigits")));if !tsIsUndefined(fraction){number:=math.Floor(tsNumber(tsToPrimitive(loop,fraction)));if math.IsNaN(number)||number<1||number>3{tsDateLocaleRange("Invalid fractionalSecondDigits option")};timeValues[3]=int(number)}
 dateStyle:=tsDateLocaleOption(loop,options,"dateStyle",[]string{"full","long","medium","short"});timeStyle:=tsDateLocaleOption(loop,options,"timeStyle",[]string{"full","long","medium","short"});if method=="toLocaleDateString"&&timeStyle!=0||method=="toLocaleTimeString"&&dateStyle!=0{tsPropertyFailure("Invalid style for Date locale method")}
 cycleIndex:=tsDateLocaleOption(loop,options,"hourCycle",[]string{"h11","h12","h23","h24"});cycle:="h12";if locale=="en-GB"{cycle="h23"};if cycleIndex>0{cycle=[]string{"h11","h12","h23","h24"}[cycleIndex-1]};hour12:=tsGet(loop,options,tsStringReference(tsStringUTF8("hour12")));if !tsIsUndefined(hour12){cycle="h23";if tsTruthy(hour12){cycle="h12"}}
 zone:=tsGet(loop,options,tsStringReference(tsStringUTF8("timeZone")));if !tsIsUndefined(zone){stamp=stamp.In(tsDateLocaleZone(tsStringValue(loop,zone).String()))}else{stamp=stamp.In(tsDateLocal())}
 hasDate,hasTime:=false,false;for _,value:=range dateValues{if value!=0{hasDate=true}};for _,value:=range timeValues{if value!=0{hasTime=true}}
 zoneStyle:="short";if timeValues[4]>0{zoneStyle=[]string{"short","long","shortOffset","longOffset","shortGeneric","longGeneric"}[timeValues[4]-1]}
 if dateStyle!=0||timeStyle!=0{if hasDate||hasTime{tsPropertyFailure("Cannot combine style and component options")};if timeStyle==1{zoneStyle="long"};key:=fmt.Sprintf("s%d%d%s",dateStyle,timeStyle,cycle);return tsStringReference(tsStringUTF8(tsDateLocaleRender(tsDateLocalePattern(locale,key),stamp,cycle,zoneStyle)))}
 dateComponents,timeComponents:=false,false;for _,value:=range dateValues[:4]{if value!=0{dateComponents=true}};for _,value:=range timeValues[:4]{if value!=0{timeComponents=true}}
 if method=="toLocaleDateString"&&!dateComponents{dateValues[0]=1;dateValues[1]=1;dateValues[2]=1};if method=="toLocaleTimeString"&&!timeComponents{timeValues[0]=1;timeValues[1]=1;timeValues[2]=1};if method=="toLocaleString"&&!dateComponents&&!timeComponents{dateValues[0]=1;dateValues[1]=1;dateValues[2]=1;timeValues[0]=1;timeValues[1]=1;timeValues[2]=1}
 dateKey:=fmt.Sprintf("%d%d%d%d%d",dateValues[0],dateValues[1],dateValues[2],dateValues[3],dateValues[4]);timeKey:=fmt.Sprintf("%d%d%d%d%d",timeValues[0],timeValues[1],timeValues[2],timeValues[3],timeValues[4]);dateText,timeText:="","";if dateKey!="00000"{dateText=tsDateLocaleRender(tsDateLocalePattern(locale,"d"+dateKey),stamp,cycle,zoneStyle)};if timeKey!="00000"{timeText=tsDateLocaleRender(tsDateLocalePattern(locale,"t"+cycle+timeKey),stamp,cycle,zoneStyle)};join:="";if dateText!=""&&timeText!=""{join=tsDateLocalePattern(locale,"j"+dateKey)};return tsStringReference(tsStringUTF8(dateText+join+timeText))
}
`

package goemit

const dateHelpersRuntime = dateParserRuntime + `
var tsDateZones struct{sync.RWMutex;locations map[string]*time.Location}
func tsDateLocal()*time.Location{zone:=strings.TrimPrefix(os.Getenv("TZ"),":");if zone==""{return time.Local};if zone=="UTC"||zone=="Etc/UTC"||zone=="GMT"{return time.UTC};tsDateZones.RLock();location:=tsDateZones.locations[zone];tsDateZones.RUnlock();if location!=nil{return location};location,err:=time.LoadLocation(zone);if err!=nil{return time.Local};tsDateZones.Lock();if tsDateZones.locations==nil{tsDateZones.locations=map[string]*time.Location{}};tsDateZones.locations[zone]=location;tsDateZones.Unlock();return location}

// Resolve wall time with JavaScript's compatible disambiguation: earlier instant
// in a repeated interval, and advance by the gap for nonexistent wall times.
func tsDateWall(year,month,day,hour,minute,second,millisecond int,location *time.Location)time.Time{
 wall:=time.Date(year,time.Month(month+1),day,hour,minute,second,millisecond*1000000,time.UTC);if location==time.UTC{return wall}
 offsets:=[5]int{};count:=0;for _,delta:=range []int64{-172800,-86400,0,86400,172800}{_,offset:=wall.Add(time.Duration(delta)*time.Second).In(location).Zone();seen:=false;for i:=0;i<count;i++{if offsets[i]==offset{seen=true}};if !seen{offsets[count]=offset;count++}}
 var best time.Time;found:=false;var after time.Time;gap:=int64(math.MaxInt64)
 for _,offset:=range offsets[:count]{candidate:=wall.Add(-time.Duration(offset)*time.Second);local:=candidate.In(location);actual:=time.Date(local.Year(),local.Month(),local.Day(),local.Hour(),local.Minute(),local.Second(),local.Nanosecond(),time.UTC);distance:=actual.UnixMilli()-wall.UnixMilli();if distance==0{if !found||candidate.Before(best){best=candidate;found=true}}else if distance>0&&distance<gap{after=candidate;gap=distance}}
 if found{return best.In(location)};if gap!=math.MaxInt64{return after.In(location)};return time.Date(year,time.Month(month+1),day,hour,minute,second,millisecond*1000000,location)
}
func tsDateParse(text string)float64{
 parsed,ok:=tsDPParseDateISOString(text);if !ok{parsed,ok=tsDPParseDateOtherString(text)};if !ok||parsed.month<1||parsed.month>12||parsed.day<1||parsed.day>31||parsed.hour>24||parsed.min>59||parsed.sec>59||parsed.hour==24&&(parsed.min!=0||parsed.sec!=0||parsed.msec!=0){return math.NaN()}
 location:=time.UTC;if parsed.isLocal{location=tsDateLocal()}else if parsed.timeZoneOffset!=0{location=time.FixedZone("",parsed.timeZoneOffset*60)}
 stamp:=tsDateWall(parsed.year,parsed.month-1,parsed.day,parsed.hour,parsed.min,parsed.sec,parsed.msec,location);return tsDateClip(float64(stamp.Unix())*1000+float64(stamp.Nanosecond())/1e6)
}
func tsDateMake(fields []float64,location *time.Location)float64{
 for _,value:=range fields{if math.IsNaN(value)||math.IsInf(value,0){return math.NaN()}}
 year:=fields[0]+math.Floor(fields[1]/12);month:=fields[1]-math.Floor(fields[1]/12)*12;if math.Abs(year)>10000000{return math.NaN()}
 base:=time.Date(int(year),time.Month(int(month)+1),1,0,0,0,0,time.UTC)
 milliseconds:=float64(base.UnixMilli())+(fields[2]-1)*86400000+fields[3]*3600000+fields[4]*60000+fields[5]*1000+fields[6]
 if math.IsNaN(milliseconds)||math.IsInf(milliseconds,0)||math.Abs(milliseconds)>8640000000000000+172800000{return math.NaN()}
 wall:=tsDateInstant(milliseconds).UTC();stamp:=tsDateWall(wall.Year(),int(wall.Month())-1,wall.Day(),wall.Hour(),wall.Minute(),wall.Second(),wall.Nanosecond()/1000000,location);return tsDateClip(float64(stamp.UnixMilli()))
}
func tsDateFields(loop *tsLoop,args []tsValue,utc bool)float64{
 if len(args)==0{return math.NaN()};fields:=[]float64{0,0,1,0,0,0,0};for i:=0;i<len(fields)&&i<len(args);i++{fields[i]=math.Trunc(tsNumber(tsToPrimitive(loop,args[i])))}
 if fields[0]>=0&&fields[0]<=99{fields[0]+=1900};location:=tsDateLocal();if utc{location=time.UTC};return tsDateMake(fields,location)
}
func tsDateYear(year int,width int)string{if year<0{return "-"+fmt.Sprintf("%0*d",width,-year)};return fmt.Sprintf("%0*d",width,year)}
func tsDateZone(stamp time.Time)string{name,offset:=stamp.Zone();if os.Getenv("TZ")=="Australia/Lord_Howe"{name="Lord Howe Standard Time";if stamp.IsDST(){name="Lord Howe Daylight Time"}};switch name{case "UTC","GMT":name="Coordinated Universal Time";case "EST":name="Eastern Standard Time";case "EDT":name="Eastern Daylight Time";case "CST":name="Central Standard Time";case "CDT":name="Central Daylight Time";case "MST":name="Mountain Standard Time";case "MDT":name="Mountain Daylight Time";case "PST":name="Pacific Standard Time";case "PDT":name="Pacific Daylight Time"};sign:="+";if offset<0{sign="-";offset=-offset};return fmt.Sprintf("GMT%s%02d%02d (%s)",sign,offset/3600,offset/60%60,name)}
func tsDateHuman(stamp time.Time,mode string)string{date:=stamp.Format("Mon Jan 02 ")+tsDateYear(stamp.Year(),4);clock:=stamp.Format("15:04:05 ")+tsDateZone(stamp);if mode=="toDateString"{return date};if mode=="toTimeString"{return clock};return date+" "+clock}
`

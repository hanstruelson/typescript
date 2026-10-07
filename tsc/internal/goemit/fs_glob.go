package goemit

const fsGlobRuntime = `
func tsFSGlobMatch(pattern, path string) bool {
    pattern=filepath.ToSlash(pattern);path=filepath.ToSlash(path)
    if strings.Contains(pattern,"{"){
        first:=strings.IndexByte(pattern,'{');last:=strings.IndexByte(pattern[first+1:],'}')
        if last>=0{last+=first+1;for _,item:=range strings.Split(pattern[first+1:last],","){if tsFSGlobMatch(pattern[:first]+item+pattern[last+1:],path){return true}};return false}
    }
    patterns,parts:=strings.Split(pattern,"/"),strings.Split(path,"/")
    var match func(int,int)bool
    match=func(i,j int)bool{
        if i==len(patterns){return j==len(parts)}
        if patterns[i]=="**"{if match(i+1,j){return true};return j<len(parts)&&!strings.HasPrefix(parts[j],".")&&match(i,j+1)}
        if j==len(parts){return false}
        if strings.HasPrefix(parts[j],".")&&!strings.HasPrefix(patterns[i],"."){return false}
        ok,err:=filepath.Match(patterns[i],parts[j]);if err!=nil{panic(tsThrown{tsFSException("TypeError","ERR_INVALID_ARG_VALUE",err.Error())})}
        return ok&&match(i+1,j+1)
    }
    return match(0,0)
}

func tsFSGlob(loop *tsLoop, args []tsValue, asynchronous bool) tsResult {
    patterns:=[]string{};value:=tsArg(args,0);options:=tsArg(args,1)
    if value.kind==tsStringKind{patterns=append(patterns,tsText(value))}else{if value.kind!=tsArrayKind{tsFSInvalid("pattern","a string or string array")};it:=tsIterate(value);for it.next(loop){if it.value.kind!=tsStringKind{tsFSInvalid("pattern","a string or string array")};patterns=append(patterns,tsText(it.value))}}
    cwd:=tsText(tsFSOption(options,"cwd",tsStringReference(tsStringUTF8("."))))
    withTypes:=tsTruthy(tsFSProperty(options,"withFileTypes"));exclude:=tsFSProperty(options,"exclude")
    out:=&tsArray{};seen:=map[string]bool{}
    var walk func(string,string)error
    walk=func(directory,relative string)error{
        entries,err:=os.ReadDir(directory);if err!=nil{return err}
        for _,entry:=range entries{
            path:=filepath.Join(relative,entry.Name());dirent:=tsFSDirent(entry,directory,"utf8");excluded:=false
            if exclude.kind==tsFunctionKind{
                item:=tsStringReference(tsStringUTF8(path));if withTypes{item=dirent}
                if asynchronous{done:=make(chan tsResult,1);loop.post(func(){var result tsResult;defer func(){if failure:=recover();failure!=nil{result=tsResult{tsUnwrap(failure),true}};done<-result}();result.value=tsCall(loop,exclude,item)});var answer tsResult;answer=<-done;if answer.rejected{panic(tsThrown{answer.value})};excluded=tsTruthy(answer.value)}else{excluded=tsTruthy(tsCall(loop,exclude,item))}
            }else if !tsIsUndefined(exclude){if exclude.kind!=tsArrayKind{tsFSInvalid("exclude","a function or string array")};it:=tsIterate(exclude);for it.next(loop){if tsFSGlobMatch(tsText(it.value),path){excluded=true;break}}}
            if excluded{continue}
            for _,pattern:=range patterns{
                candidate:=path;if filepath.IsAbs(pattern){candidate=filepath.Join(cwd,path);candidate,_=filepath.Abs(candidate)}
                if tsFSGlobMatch(pattern,candidate)&&!seen[candidate]{seen[candidate]=true;item:=tsStringReference(tsStringUTF8(candidate));if withTypes{item=dirent};out.values=append(out.values,item)}
            }
            if entry.IsDir()&&!strings.HasPrefix(entry.Name(),"."){if err:=walk(filepath.Join(directory,entry.Name()),path);err!=nil{return err}}
        };return nil
    }
    if err:=walk(cwd,"");err!=nil{return tsResult{tsFSError(err,"glob",cwd,""),true}}
    return tsResult{tsArrayValue(out),false}
}

func tsFSGlobIterator(loop *tsLoop, args []tsValue) tsValue {
    snapshot:=tsFSPrepare(loop,args);var loaded *tsPromise;index:=0;closed:=false
    iterator:=&tsAsyncIterator{
        next:func(loop *tsLoop)*tsPromise{
            if closed{return loop.resolved(tsIteratorResult(tsU,true),false)}
            if loaded==nil{loaded=loop.start(func()(result tsResult){defer func(){if failure:=recover();failure!=nil{result=tsResult{tsUnwrap(failure),true}}}();result=tsFSGlob(loop,snapshot,true);return result})}
            result:=loop.promise()
            loop.await(tsPromiseValue(loaded),func(value tsResult){if value.rejected{closed=true;result.settle(value);return};array:=(*tsArray)(value.value.ref);if closed||index>=array.length(){closed=true;result.settle(tsResult{tsIteratorResult(tsU,true),false})}else{item:=array.at(index);index++;result.settle(tsResult{tsIteratorResult(item,false),false})}})
            return result
        },
        close:func(loop *tsLoop)*tsPromise{closed=true;return loop.resolved(tsIteratorResult(tsU,true),false)},
    }
    return tsAsyncIteratorObject(iterator)
}
`

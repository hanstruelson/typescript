package goemit

// Declared property entries are lenses into concrete fields, never copies of
// their values. The same hash table also holds dynamically added properties.
const ClassRuntime = `
type tsClass struct {builtin string;layout *tsClassLayout;construct func(...tsValue) tsValue;parent *tsClass;static tsValue}
func tsConstruct(class tsValue,args ...tsValue) tsValue {if class.kind!=tsClassKind {panic("Value is not a constructor")};return (*tsClass)(class.ref).construct(args...)}
type tsProperty struct {get func() tsValue;set func(tsValue) tsValue}
type tsProperties struct {self unsafe.Pointer;layout *tsClassLayout;prototype *tsProperties;order []string;own map[string]bool;class *tsClass;initialized bool;declared map[string]tsProperty;extra map[string]tsValue;methods map[string]tsValue}
func tsNewProperties() *tsProperties {return &tsProperties{own:make(map[string]bool),declared:make(map[string]tsProperty),extra:make(map[string]tsValue),methods:make(map[string]tsValue)}}
func(p *tsProperties)define(name string){if !p.own[name]{p.own[name]=true;p.order=append(p.order,name)}}
func(p *tsProperties)get(name string) tsValue {if property,ok:=p.declared[name];ok {return property.get()};if value,ok:=p.extra[name];ok{return value};if value,ok:=p.methods[name];ok{return value};if p.prototype!=nil {return p.prototype.get(name)};return tsU}
func(p *tsProperties)set(name string,value tsValue) tsValue {if property,ok:=p.declared[name];ok {return property.set(value)};p.define(name);p.extra[name]=value;return value}
func tsRequireThis(object tsValue) tsValue {if !tsInstanceProperties(object).initialized {panic("Cannot access this before super constructor")};return object}
func tsInstanceOf(value,class tsValue) bool {if class.kind!=tsClassKind {panic("Right operand of instanceof is not a class")};constructor:=(*tsClass)(class.ref);if constructor.builtin!=""{switch constructor.builtin{case "Map":return value.kind==tsCollectionKind&&!(*tsCollection)(value.ref).set;case "Set":return value.kind==tsCollectionKind&&(*tsCollection)(value.ref).set;case "Array":return value.kind==tsArrayKind;case "ArrayBuffer":return value.kind==tsArrayBufferKind;default:return value.kind==tsTypedArrayKind&&(*tsTypedArray)(value.ref).name==constructor.builtin}};if value.kind!=tsInstanceKind {return false};for at:=(*tsProperties)(value.ref).class;at!=nil;at=at.parent {if at==constructor {return true}};return false}
type tsClassLayout struct {name string}
`

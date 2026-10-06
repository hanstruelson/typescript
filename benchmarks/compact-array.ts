// The type is erased by Node and Bun; this emitter stores one byte per element.
const values: int8[] = new Array<int8>(4000000);
values.fill(7);
console.log(values.length + values[42]);

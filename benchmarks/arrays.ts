// The runner replaces ELEMENT_TYPE to compare native storage with TSValue storage.
const values: ELEMENT_TYPE[] = [];
for (let i: number = 0; i < 4096; i++) values.push(i);
let total: number = 0;
for (let round: number = 0; round < 512; round++) {
    values.forEach(value => { total += value; });
}
console.log(total);

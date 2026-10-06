// Each task owns its local values; no mutable data is shared between workers.
function compute(seed: number): number {
    const values: int8[] = new Array<int8>(8000000);
    values.fill(seed);
    return values.length + values[42];
}

async function main() {
    const results: Promise<number>[] = [];
    for (let seed: number = 0; seed < 8; seed++) {
        // The runner uses ordinary calls for the sequential comparison.
        results.push(go compute(seed));
    }
    const values = await Promise.all(results);
    let total: number = 0;
    for (const value of values) total += value;
    console.log(total);
}
main();

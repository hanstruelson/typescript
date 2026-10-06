import { readFile } from "node:fs/promises";

declare function delay(milliseconds: number): Promise<void>;

async function run(): Promise<void> {
    const a = 10;
    const contents = await readFile("examples/go-target/async/message.txt", "utf8");
    console.log(a, contents);
    for (let i = 0; i < 3; i++) {
        await delay(0);
        console.log(i);
    }
}

run();
console.log("sync");

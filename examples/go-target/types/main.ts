let count: number = 10;
let optional: number | null | undefined = null;
let dynamic: any = "42";

function increment(value: number): number {
    return value + 1;
}

console.log(increment(count), optional);
try {
    console.log(increment(dynamic));
} catch (error: any) {
    console.log(error.message);
}

const text = "A😀B";
console.log(text.length, text.charCodeAt(1), text.codePointAt(1));
console.log(text.slice(1, 3), "I".toLocaleLowerCase("tr"));
console.log("item42".match(/(?<digits>\d+)/)!.groups!.digits);

const { value = 10, ...rest } = { value: 3, extra: 4 };
console.log(value, rest.extra);

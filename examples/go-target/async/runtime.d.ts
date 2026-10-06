// Signatures for the text-file runtime adapter implemented by the Go target.
declare module "node:fs/promises" {
    export function readFile(path: string, encoding: "utf8"): Promise<string>;
}

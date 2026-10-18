import { runDirectBatch, runFactoryBatch, runGlobalStateCheck, runRecompileBatch } from "./batches.js";
import { getSuccessfulMatches } from "./state.js";

const factoryMatches = runFactoryBatch(32);
const directMatches = runDirectBatch(32);
const recompileMatches = runRecompileBatch(16);
const globalMatches = runGlobalStateCheck();

console.log("factory", factoryMatches, "/32");
console.log("direct", directMatches, "/32");
console.log("recompile", recompileMatches, "/16");
console.log("global state", globalMatches, "/4 calls");
console.log("shared imported state", getSuccessfulMatches(), "/83");

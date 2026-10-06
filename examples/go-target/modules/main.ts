export const B = 20;
import { A } from "./helper.js";
console.log("main", A, B);
import "./helper.js";
console.log("main complete");

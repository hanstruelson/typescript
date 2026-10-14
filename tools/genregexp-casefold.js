#!/usr/bin/env node
// Generate pinned Unicode case-fold and ECMAScript Canonicalize tables.
const fs = require("node:fs");
const path = require("node:path");

const nodeVersion = "24.19.0";
const icuVersion = "78.3";
const unicodeVersion = "17.0";
if (
  process.versions.node !== nodeVersion ||
  process.versions.icu !== icuVersion ||
  process.versions.unicode !== unicodeVersion
) {
  throw new Error(
    `expected Node ${nodeVersion}/ICU ${icuVersion}/Unicode ${unicodeVersion}; got Node ${process.versions.node}/ICU ${process.versions.icu}/Unicode ${process.versions.unicode}`,
  );
}

const root = path.resolve(__dirname, "..");
const sourcePath = path.join(root, "tools/unicode/17.0.0/CaseFolding.txt");
const outputPath = path.join(root, "tsc/internal/regexpengine/regexp_casefold_data.go");
const source = fs.readFileSync(sourcePath, "utf8");
if (!source.startsWith("# CaseFolding-17.0.0.txt")) {
  throw new Error("unexpected CaseFolding.txt version");
}

const parent = new Map();
const simpleMappings = [];
function find(codePoint) {
  if (!parent.has(codePoint)) parent.set(codePoint, codePoint);
  let root = codePoint;
  while (parent.get(root) !== root) root = parent.get(root);
  while (parent.get(codePoint) !== codePoint) {
    const next = parent.get(codePoint);
    parent.set(codePoint, root);
    codePoint = next;
  }
  return root;
}
function union(left, right) {
  const leftRoot = find(left);
  const rightRoot = find(right);
  if (leftRoot !== rightRoot) parent.set(rightRoot, leftRoot);
}

for (const [index, line] of source.split(/\r?\n/).entries()) {
  const data = line.split("#", 1)[0].trim();
  if (!data) continue;
  const fields = data.split(";").map((field) => field.trim());
  if (fields.at(-1) === "") fields.pop();
  if (fields.length !== 3) throw new Error(`malformed UCD row ${index + 1}`);
  const [codePoint, status, mapping] = fields;
  if (status !== "C" && status !== "S") continue;
  const equivalents = mapping.split(/\s+/);
  if (equivalents.length !== 1) throw new Error(`non-simple mapping at row ${index + 1}`);
  const sourceCodePoint = Number.parseInt(codePoint, 16);
  const foldedCodePoint = Number.parseInt(equivalents[0], 16);
  simpleMappings.push([sourceCodePoint, foldedCodePoint]);
  union(sourceCodePoint, foldedCodePoint);
}

const groupsByRoot = new Map();
for (const codePoint of parent.keys()) {
  const rootCodePoint = find(codePoint);
  if (!groupsByRoot.has(rootCodePoint)) groupsByRoot.set(rootCodePoint, []);
  groupsByRoot.get(rootCodePoint).push(codePoint);
}
const classes = [...groupsByRoot.values()]
  .filter((group) => group.length > 1)
  .map((group) => group.sort((a, b) => a - b));
const sharpS = classes.find((group) => group.includes(0x00df) || group.includes(0x1e9e));
if (!sharpS) throw new Error("Unicode data has no sharp-S simple-fold class");
if (!sharpS.includes(0x1df95)) sharpS.push(0x1df95);
sharpS.sort((a, b) => a - b);
if (!simpleMappings.some(([sourceCodePoint]) => sourceCodePoint === 0x1df95)) {
  simpleMappings.push([0x1df95, 0x00df]);
}
classes.sort((left, right) => left[0] - right[0]);
if (classes.length > 0xffff) throw new Error("too many fold classes for uint16 indices");

// Port Canonicalize from gen-regexp-special-case.cc. These tables let the Go
// implementation reproduce CaseFolding::IgnoreSet and
// CaseFolding::SupplementaryEquivalentsSet without ICU at runtime.
const legacyUppercase = new Map();
for (let codePoint = 0; codePoint <= 0xffff; codePoint++) {
  const uppercase = String.fromCharCode(codePoint).toUpperCase();
  // ECMAScript Canonicalize keeps the original code unit when uppercase
  // expands to more than one UTF-16 code unit.
  if (uppercase.length === 1) {
    const mapped = uppercase.charCodeAt(0);
    if (mapped !== codePoint) legacyUppercase.set(codePoint, mapped);
  }
}

function canonicalizeLegacy(codePoint) {
  const uppercase = legacyUppercase.get(codePoint) ?? codePoint;
  return codePoint >= 0x80 && uppercase < 0x80 ? codePoint : uppercase;
}

const ignoreSet = new Set();
const specialAdd = new Set();
const supplementaryEquivalents = new Set();
for (const group of classes) {
  const bmp = group.filter((codePoint) => codePoint <= 0xffff);
  const hasSupplementary = group.some((codePoint) => codePoint >= 0x10000);
  for (const codePoint of bmp) {
    if (hasSupplementary) supplementaryEquivalents.add(codePoint);
    const canonical = canonicalizeLegacy(codePoint);
    let hasMatchingCanonical = false;
    let hasNonMatchingCanonical = false;
    for (const equivalent of bmp) {
      if (equivalent === codePoint) continue;
      if (canonicalizeLegacy(equivalent) === canonical) {
        hasMatchingCanonical = true;
      } else {
        hasNonMatchingCanonical = true;
      }
    }
    if (hasNonMatchingCanonical) {
      if (hasMatchingCanonical) specialAdd.add(codePoint);
      else ignoreSet.add(codePoint);
    }
  }
}

// The C++ generator checks that each simple-fold class contains at most one
// non-trivial ECMAScript class after IgnoreSet members are removed.
for (const codePoint of specialAdd) {
  const group = classes.find((candidate) => candidate.includes(codePoint));
  const canonical = canonicalizeLegacy(codePoint);
  for (const equivalent of group) {
    if (equivalent > 0xffff || ignoreSet.has(equivalent)) continue;
    if (canonicalizeLegacy(equivalent) !== canonical) {
      throw new Error(`case-fold class violates IgnoreSet invariant at U+${codePoint.toString(16)}`);
    }
  }
}

const lines = [
  "// Code generated by tools/genregexp-casefold.js; DO NOT EDIT.",
  "// Simple folds: Unicode CaseFolding-17.0.0.txt, statuses C and S.",
  "// Legacy uppercase: ECMAScript toUpperCase from Node 24.19.0 / ICU 78.3.",
  "// Non-Unicode IgnoreSet and supplementary equivalents follow V8's generator.",
  "// U+1DF95 follows the pinned V8 source's documented supplementary fold.",
  "// Unicode data is subject to the Unicode Terms of Use.",
  "package regexpengine",
  "",
  'const regexpCaseFoldDataVersion = "17.0.0"',
  'const regexpLegacyCanonicalizeDataVersion = "17.0.0"',
  "",
  "var regexpCaseFoldClasses = [][]rune{",
];
for (const group of classes) {
  lines.push(`\t{${group.map((codePoint) => `0x${codePoint.toString(16)}`).join(", ")}},`);
}
lines.push("}", "", "var regexpCaseFoldClassByRune = map[rune]uint16{");
for (let index = 0; index < classes.length; index++) {
  const pairs = classes[index]
    .map((codePoint) => `0x${codePoint.toString(16)}: ${index}`)
    .join(", ");
  lines.push(`\t${pairs},`);
}
lines.push("}", "", "var regexpSimpleCaseFoldTo = map[rune]rune{");
for (const [sourceCodePoint, foldedCodePoint] of simpleMappings) {
  lines.push(`\t0x${sourceCodePoint.toString(16)}: 0x${foldedCodePoint.toString(16)},`);
}
lines.push("}", "", "var regexpLegacyUppercase = map[rune]rune{");
for (const [codePoint, mapped] of legacyUppercase) {
  lines.push(`\t0x${codePoint.toString(16)}: 0x${mapped.toString(16)},`);
}
lines.push("}", "", "var regexpLegacyIgnoreSet = map[rune]bool{");
for (const codePoint of [...ignoreSet].sort((a, b) => a - b)) {
  lines.push(`\t0x${codePoint.toString(16)}: true,`);
}
lines.push("}", "", "var regexpSupplementaryEquivalentsSet = map[uint16]bool{");
for (const codePoint of [...supplementaryEquivalents].sort((a, b) => a - b)) {
  lines.push(`\t0x${codePoint.toString(16)}: true,`);
}
lines.push("}", "");
fs.writeFileSync(outputPath, lines.join("\n"));

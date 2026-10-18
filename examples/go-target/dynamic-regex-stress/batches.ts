import { makeDynamicMatcher, replaceMatcher } from "./factory.js";
import { recordMatch } from "./state.js";

export function runFactoryBatch(count: number): number {
  let matches = 0;
  for (let i = 0; i < count; i++) {
    const matcher = makeDynamicMatcher("^/user/" + i + "/([0-9]+)$", "i");
    const result = matcher.exec("/USER/" + i + "/42");
    const matched = result !== null && result[1] === "42";
    if (matched) matches++;
    recordMatch(matched);
  }
  return matches;
}

export function runDirectBatch(count: number): number {
  let matches = 0;
  for (let i = 0; i < count; i++) {
    const matcher = new RegExp("^tag-" + i + ":[a-z]+$", "i");
    const matched = matcher.test("TAG-" + i + ":value");
    if (matched) matches++;
    recordMatch(matched);
  }
  return matches;
}

export function runRecompileBatch(count: number): number {
  const matcher = makeDynamicMatcher("^seed$", "");
  let matches = 0;
  for (let i = 0; i < count; i++) {
    replaceMatcher(matcher, "^key-" + i + "=[a-z]+$", "i");
    const matched = matcher.test("KEY-" + i + "=value");
    if (matched) matches++;
    recordMatch(matched);
  }
  return matches;
}

export function runGlobalStateCheck(): number {
  const matcher = makeDynamicMatcher("id=[0-9]+", "g");
  let matches = 0;
  const input = "id=1 id=2";
  for (let i = 0; i < 4; i++) {
    const matched = matcher.test(input);
    if (matched) matches++;
    recordMatch(matched);
  }
  return matches;
}

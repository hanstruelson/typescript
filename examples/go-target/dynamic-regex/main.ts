function createRouteMatcher(pattern: string): RegExp {
  return new RegExp(pattern, "i");
}

const prefix = "^/users/";
const suffix = "[0-9]+/?$";
const matcher = createRouteMatcher(prefix + suffix);
const match = matcher.exec("/USERS/42");
console.log(match ? match[0] : "no match");

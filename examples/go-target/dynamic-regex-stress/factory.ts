export function makeDynamicMatcher(source: string, flags: string): RegExp {
  return new RegExp(source, flags);
}

export function replaceMatcher(matcher: RegExp, source: string, flags: string): RegExp {
  matcher.compile(source, flags);
  return matcher;
}

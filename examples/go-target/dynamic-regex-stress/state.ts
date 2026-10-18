let successfulMatches = 0;

export function recordMatch(matched: boolean): void {
  if (matched) successfulMatches++;
}

export function getSuccessfulMatches(): number {
  return successfulMatches;
}
